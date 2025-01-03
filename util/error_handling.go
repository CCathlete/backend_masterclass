package u

import "fmt"

func Must(val any, err error) any {
	if err != nil {
		panic(err)
	}

	return val
}

type PropagatedError *error

func WrapError(err PropagatedError, msg string) {
	if *err == nil {
		*err = fmt.Errorf("%s: %w", msg, *err)
	} else {
		*err = fmt.Errorf("%s", msg)
	}
}

func NewPropagatedError() (err PropagatedError) {
	err = new(error)

	return
}

func GetMessage(err PropagatedError) string {
	return fmt.Sprintf("%v", *err)
}
