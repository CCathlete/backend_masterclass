package api

import (
	u "backend-masterclass/util"

	"github.com/go-playground/validator/v10"
)

var validCurrency validator.Func = func(fl validator.FieldLevel) bool {
	// we're useing the Field() method to get the field value.
	// Field() returns a reflect.Value, we convert it to an interface and then use an assertion to validate that it's a string.
	currency, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	// We check if the currency is one of the valid options.
	return u.IsValidCurrency(currency)
}
