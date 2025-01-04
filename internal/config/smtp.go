package config

import (
	"crypto/tls"
	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
	"os"
	"strconv"
)

var SMTPDialer *gomail.Dialer

func SMTPConnect() {
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	email := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMTP_PASSWORD")

	portInt, err := strconv.Atoi(port)
	if err != nil {
		panic(err)
	}

	// create a new dialer
	dialer := gomail.NewDialer(host, portInt, email, password)
	dialer.TLSConfig = &tls.Config{ServerName: host}

	// test the connection
	s, err := dialer.Dial()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	SMTPDialer = dialer

	logrus.Info("SMTP Connected")
}
