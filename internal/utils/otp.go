package utils

import (
	"authentication/internal/config"
	"bytes"
	"context"
	"crypto/rand"
	"errors"
	"github.com/go-redis/redis/v8"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
	"html/template"
	"math/big"
	"os"
	"strconv"
	"time"
)

var (
	OTPLength  int
	otpCharSet string
	otpExp     int
	resetLink  string
)

func init() {
	var err error
	OTPLength, err = strconv.Atoi(os.Getenv("OTP_LENGTH"))
	if err != nil {
		OTPLength = 10
	}

	otpCharSet = os.Getenv("OTP_SECRET")
	if otpCharSet == "" {
		otpCharSet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	}

	otpExp, err = strconv.Atoi(os.Getenv("OTP_MINUTES_LIFESPAN"))
	if err != nil {
		otpExp = 10
	}
	resetLink = os.Getenv("LINK_RESET_PASSWORD")

}

func GenerateOTP() string {
	result := make([]byte, OTPLength)
	charsetLength := big.NewInt(int64(len(otpCharSet)))

	for i := range result {
		num, _ := rand.Int(rand.Reader, charsetLength)
		result[i] = otpCharSet[num.Int64()]
	}

	return string(result)
}

func AddOTPtoRedis(c context.Context, otp string, email string) error {
	key := otkKeyPrefix + email

	data, _ := bcrypt.GenerateFromPassword([]byte(otp), 10)

	err := config.RedisClient.Set(c, key, data, time.Duration(otpExp)*time.Minute).Err()
	if err != nil {
		return err
	}

	return nil
}

func SendOTP(otp string, recipient string) error {
	sender := os.Getenv("SMTP_EMAIL")
	body, err := generateEmailBody(recipient, otp)
	if err != nil {
		return err
	}

	message := gomail.NewMessage()
	message.SetHeader("From", sender)
	message.SetHeader("To", recipient)
	message.SetHeader("Subject", "Reservo password reset")
	message.SetBody("text/html", body)
	message.Embed(logoPath)
	message.Embed(facebookPath)
	message.Embed(instagramPath)
	message.Embed(linkedinPath)
	message.Embed(lockPath)
	message.Embed(twitterPath)

	if err := config.SMTPDialer.DialAndSend(message); err != nil {
		return err
	}

	return nil
}

func generateEmailBody(recipient string, otp string) (string, error) {
	t, _ := template.ParseFiles(templatePath)
	var body bytes.Buffer
	err := t.Execute(&body, struct {
		Link  string
		Email string
		Otp   string
	}{
		Link:  resetLink,
		Email: recipient,
		Otp:   otp,
	})

	if err != nil {
		return "", err
	}
	return body.String(), nil
}

func VerifyOTP(otp string, email string, c context.Context) (error, bool) {
	key := otkKeyPrefix + email

	// get the value for the key
	value, err := config.RedisClient.Get(c, key).Result()
	if err != nil {
		// the following states that the key was not found
		if errors.Is(err, redis.Nil) {
			return errors.New("otp expired / incorrect email"), false
		}

		return err, true
	}

	// compare received otp's hash with value in redis
	err = bcrypt.CompareHashAndPassword([]byte(value), []byte(otp))
	if err != nil {
		return errors.New("incorrect otp"), false
	}

	// delete redis key to prevent abuse of otp
	err = config.RedisClient.Del(c, key).Err()
	if err != nil {
		return err, true
	}

	return nil, false
}
