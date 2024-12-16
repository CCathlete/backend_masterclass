package util

func Must(val any, err error) any {
	if err != nil {
		panic(err)
	}

	return val
}
