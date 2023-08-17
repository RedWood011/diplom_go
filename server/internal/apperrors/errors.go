package apperrors

import "errors"

var (
	ErrUserExists = errors.New("user already exists")
	ErrNotFound   = errors.New("user not found")

	ErrAuth          = errors.New("invalid login or password")
	ErrUserNoDeleted = errors.New("user is not deleted")
)
