package apperrors

import "errors"

var (
	ErrAuth = errors.New("invalid login or password")
)
