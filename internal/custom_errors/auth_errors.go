package custom_errors

import "errors"

var (
	ErrInvalidCredentials = errors.New("email or password is incorrect")
	ErrInvalidEmail       = errors.New("email is incorrect")
)
