package utils

import (
	"errors"
	"github.com/kiminodare/HOVARLAY-BE/ent/generated"
)

// Custom errors
var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidData        = errors.New("invalid data")
	ErrInvalidCredentials = errors.New("invalid email or password")
)

// MapEntError converts Ent/DB errors to custom errors
func MapEntError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case generated.IsConstraintError(err):
		return ErrEmailAlreadyExists
	case generated.IsNotFound(err):
		return ErrUserNotFound
	case generated.IsValidationError(err):
		return ErrInvalidData
	default:
		return err
	}
}
