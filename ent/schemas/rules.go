package schemas

import (
	"errors"
	"regexp"
)

var nameRegex = regexp.MustCompile(`^[A-ZÁÉÍÓÚÑ][a-záéíóúñ]+(?:\s[A-ZÁÉÍÓÚÑ][a-záéíóúñ]+)*$`)
var emailRegex = regexp.MustCompile(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`)

func ValidateEmail(s string)  error{
	if !emailRegex.MatchString(s) {
		return errors.New("The text input is not an email")
	}
	return nil
}
func ValidateName(s string) error{
	if !nameRegex.MatchString(s) {
		return errors.New("The name format is invalid (should start with capital letters and no numbers")
	}
	return nil
}