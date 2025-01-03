package validation

import (
	"fmt"
	"regexp"
)

var (
	// We build a regexp object out of the pattern and alias it's MatchString method.
	isValidUsername = regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString
)

type PropagatedError *error

func WrapError(err PropagatedError, msg string) {
	if *err == nil {
		*err = fmt.Errorf("%s: %w", msg, *err)
	} else {
		*err = fmt.Errorf("%s", msg)
	}
}

func ValidateString(
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

	ok = ValidateString(s, 3, 25, err) && isValidUsername(s)
	if !ok {
		WrapError(err, "invalid username")
		return
	}

	return
}
