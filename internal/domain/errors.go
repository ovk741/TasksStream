package domain

import "errors"

var (
	ErrNotFound           = errors.New("not found")
	ErrInvalidInput       = errors.New("invalid input")
	ErrConflict           = errors.New("conflict")
	ErrInternal           = errors.New("internal error")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrForbidden          = errors.New("forbidden")
)
