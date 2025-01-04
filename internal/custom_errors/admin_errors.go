package custom_errors

import "errors"

var (
	ErrAdminEmailExists     = errors.New("admin with this email already exists")
	ErrAdminNotFound        = errors.New("admin with such email dont exist")
	ErrIncorrectPasswordOTP = errors.New("incorrect OTP password")
)
