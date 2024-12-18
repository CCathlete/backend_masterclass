package u

import "fmt"

type ErrorInstructions struct {
	WhatToDo     int
	Wrappers     []string
	WrapperIndex int
	CallerName   string
}

type HandledResult struct {
	ErrNotNil  bool
	WrappedErr error
	// Additional return values.
}

const (
	Wrap = 1
)

func Must(val any, err error) any {
	if err != nil {
		panic(err)
	}

	return val
}

// Provides error handling logic.
// Curretly only wraps.
func HandleError(
	i ErrorInstructions,
	err error,
) (result HandledResult) {

	switch i.WhatToDo {
	case Wrap:
		if err != nil {
			result.ErrNotNil = true
			result.WrappedErr = fmt.Errorf("%s: %s: %w", i.CallerName,
				i.Wrappers[i.WrapperIndex], err)
		}

	}
	return
}
