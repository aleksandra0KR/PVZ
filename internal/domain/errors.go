package domain

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrMissingSecret      = errors.New("SECRET_KEY is not set")
	ErrInitAuthConfig     = errors.New("failed to init auth config")
	ErrReceptionNotFound  = errors.New("no active reception found for this PVZ")
)
