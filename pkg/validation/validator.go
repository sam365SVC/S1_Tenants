package validation

import (
	"regexp"
	"unicode"

	"github.com/go-playground/validator/v10"
)

var Validator *validator.Validate

var nameValidatorRegex=regexp.MustCompile(`^[a-zA-ZáéíóúÁÉÍÓÚñÑ\s]+$`)
func isPasswordSecure(fl validator.FieldLevel)bool{
	password:=fl.Field().String()

	var(
		haveUpper=false
		haveLower=false
		haveNumber=false
		haveSpecial=false
		haveSpace=false
	)

	for _, letter := range  password{
		switch{
			case unicode.IsUpper(letter):
				haveUpper=true
			case unicode.IsLower(letter):
				haveLower=true
			case unicode.IsNumber(letter):
				haveNumber=true
			case unicode.IsPunct(letter)||unicode.IsSymbol(letter):
				haveSpecial=true
			case unicode.IsSpace(letter):
				haveSpace=true
		}
	}
	return haveUpper&&haveLower&&haveNumber&&haveSpecial&&!haveSpace
}

func InitValidator()  {
	Validator=validator.New()

	Validator.RegisterValidation("is_name",func(fl validator.FieldLevel) bool {
		return nameValidatorRegex.MatchString(fl.Field().String())
	})
	Validator.RegisterValidation("secure_password",isPasswordSecure)
}