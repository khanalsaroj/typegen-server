package domain

import "errors"

var (
	ErrNotFound           = errors.New("resource not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrForbidden          = errors.New("forbidden")
	ErrBadRequest         = errors.New("bad request")
	ErrConflict           = errors.New("resource already exists")
	ErrInternalServer     = errors.New("internal server error")
	ErrInvalidToken       = errors.New("invalid or expired token")
)
