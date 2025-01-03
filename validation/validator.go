package validation

import (
	u "backend-masterclass/util"
	"fmt"
	"regexp"
)

var (
	// We build a regexp object out of the pattern and alias it's MatchString method.
	isValidUsername = regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString

	isValidEmail = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`).MatchString

	isValidFullName = regexp.MustCompile(`^[a-zA-Z\s+$]`).MatchString
)

type PropagatedError = u.PropagatedError

var WrapError = u.WrapError

func IsValidLength(
	s string,
	minLength, maxLength int,
	err PropagatedError,
) (ok bool) {

	len := len(s)
	ok = len >= minLength && len <= maxLength
	if !ok {
		WrapError(err, fmt.Sprintf("string length must be between %d and %d",
			minLength, maxLength))
	}

	return
}

func ValidateUsername(s string, err PropagatedError) (ok bool) {

	ok = IsValidLength(s, 3, 25, err) && isValidUsername(s)
	if !ok {
		WrapError(err, "invalid username")
		return
	}

	return
}

func ValidateFullName(s string, err PropagatedError) (ok bool) {

	ok = IsValidLength(s, 3, 50, err) && isValidFullName(s)
	if !ok {
		WrapError(err, "invalid full name")
		return
	}

	return
}

func ValidatePassword(s string, err PropagatedError) (ok bool) {

	ok = IsValidLength(s, 8, 64, err)
	if !ok {
		WrapError(err, "invalid password")
		return
	}

	return
}

func ValidateEmail(s string, err PropagatedError) (ok bool) {

	ok = IsValidLength(s, 5, 255, err) && isValidEmail(s)
	if !ok {
		WrapError(err, "invalid email")
		return
	}

	return
}

func ValidateCurrency(s string, err PropagatedError) (ok bool) {

	ok = u.IsValidCurrency(s)
	if !ok {
		WrapError(err, "invalid currency")
		return
	}

	return
}
