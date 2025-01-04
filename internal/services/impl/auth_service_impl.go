package impl

import (
	"authentication/internal/custom_errors"
	"authentication/internal/entities"
	"authentication/internal/models"
	"authentication/internal/repositories"
	"authentication/internal/services"
	"authentication/internal/utils"
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	userRepo repositories.UserRepository
}

func NewAuthService(adminRepo repositories.UserRepository) services.AuthService {
	return &authService{userRepo: adminRepo}
}

func (s *authService) Create(admin *entities.User) error {

	existingUser, err := s.userRepo.GetByEmail(admin.Email)
	if err == nil && existingUser != nil {
		return custom_errors.ErrAdminEmailExists
	}

	if err := s.prepareAdmin(admin); err != nil {
		return err
	}

	return s.userRepo.Create(admin)
}

func (s *authService) GetByID(id uint) (*entities.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, custom_errors.ErrAdminNotFound
	}
	return user, nil
}

func (s *authService) LoginCheck(data models.LoginInput) (string, string, error) {
	var err error
	var u *entities.User

	if utils.IsEmail(data.Email) {
		u, err = s.userRepo.GetByEmail(data.Email)
	}

	if err != nil {
		return "", "", custom_errors.ErrAdminNotFound
	}

	err = s.verifyPassword(data.Password, u.Password)
	if err != nil && errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return "", "", custom_errors.ErrInvalidCredentials
	}

	accessToken, refreshToken, err := utils.GenerateTokens(u.ID)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *authService) SendResetPassword(c context.Context, email string) error {
	if !utils.IsEmail(email) {
		return custom_errors.ErrInvalidEmail
	}

	_, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return custom_errors.ErrAdminNotFound
	}

	otp := utils.GenerateOTP()
	err = utils.AddOTPtoRedis(c, otp, email)
	if err != nil {
		return err
	}

	err = utils.SendOTP(otp, email)
	if err != nil {
		return err
	}
	return nil
}

func (s *authService) ResetPassword(c context.Context, input models.ResetPasswordInput) error {
	err, _ := utils.VerifyOTP(input.OTP, input.Email, c)
	if err != nil {
		return custom_errors.ErrIncorrectPasswordOTP
	}

	user, err := s.userRepo.GetByEmail(input.Email)
	if err != nil {
		return custom_errors.ErrAdminNotFound
	}

	user.Password = input.NewPassword
	if err := s.prepareAdmin(user); err != nil {
		return err
	}

	return s.userRepo.Update(user)
}

func (s *authService) verifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (s *authService) prepareAdmin(user *entities.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return nil
}
