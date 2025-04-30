package utils

import "fmt"

func Must[T any](v T, err error) T {
	if err != nil {
		panic(fmt.Sprintf("must has encountered error: %v", err.Error()))
	}
	return v
}
