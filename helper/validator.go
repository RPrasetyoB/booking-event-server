package helper

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

func GetErrorMessage(err error) string {
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return "Validation error: " + err.Error()
	}

	var errorMsgs []string
	for _, fieldErr := range validationErrors {
		message := GetMessage(fieldErr)
		if len(message) > 0 {
			message = strings.ToLower(string(message[0])) + message[1:]
		}
		errorMsgs = append(errorMsgs, message)
	}

	return strings.Join(errorMsgs, " | ")
}

func GetMessage(fieldErr validator.FieldError) string {
	switch fieldErr.Tag() {
	case "required":
		return fieldErr.Field() + " is required"
	case "min":
		return fieldErr.Field() + " must be at least " + fieldErr.Param() + " characters long"
	case "max":
		return fieldErr.Field() + " must be at most " + fieldErr.Param() + " characters long"
	case "numeric":
		return fieldErr.Field() + " must be a valid number"
	case "email":
		return fieldErr.Field() + " must be a valid email address"
	case "oneof":
		return fieldErr.Field() + " must be 'hr' or 'vendor'"

	default:
		return fieldErr.Error()
	}
}
