package Middlewares

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func UsernameValid(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	isValid := regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(username)
	return isValid
}
