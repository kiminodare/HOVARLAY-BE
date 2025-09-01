package utils

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

func FormatValidationErrors(err error) []string {
	var messages []string
	if err == nil {
		return messages
	}

	var validate validator.ValidationErrors
	if errors.As(err, &validate) {
		for _, fe := range validate {
			messages = append(messages, formatFieldError(fe))
		}
		return messages
	}

	messages = append(messages, err.Error())
	return messages
}

func formatFieldError(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", fe.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s", fe.Field(), fe.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s", fe.Field(), fe.Param())
	case "uuid4":
		return fmt.Sprintf("%s must be a valid UUID", fe.Field())
	default:
		return fmt.Sprintf("%s failed validation: %s", fe.Field(), fe.Tag())
	}
}
