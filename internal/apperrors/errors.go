package apperrors

import "errors"

//package:user
var (
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidEmail      = errors.New("email cannot be empty")
	ErrUserAlreadyExists = errors.New("user already exists")
)
