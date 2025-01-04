package services

import (
	"authentication/internal/entities"
	"authentication/internal/models"
	"context"
)

type AuthService interface {
	Create(admin *entities.User) error
	LoginCheck(input models.LoginInput) (string, string, error)
	GetByID(id uint) (*entities.User, error)
	SendResetPassword(c context.Context, email string) error
	ResetPassword(c context.Context, input models.ResetPasswordInput) error
}
