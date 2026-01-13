package main

import (
	"errors"
	"fmt"
)

type FooError struct {
	Err error
}

func (e *FooError) Error() string {
	return fmt.Sprintf("%s", "foo error: " + e.Err.Error())
}

// after Go 1.13
func wrapErrorV2(err error) error {
	return fmt.Errorf("foo v2 error: %w", err)
}

func wrapErrorV3(err error) error {
	return fmt.Errorf("foo v3 error: %v", err)
}

func main() {
	wrappedErr := &FooError{Err: fmt.Errorf("original error")}
	fmt.Println(wrappedErr)

	wrappedErrV2 := wrapErrorV2(wrappedErr)
	fmt.Println(wrappedErrV2)
	var errW *FooError
	if errors.As(wrappedErrV2, &errW) {
		fmt.Println("wrappedErrV2 contains wrappedErr")
	}
	innerErr := errors.Unwrap(wrappedErrV2)
	fmt.Println("Unwrapped error:", innerErr)

	wrappedErrV3 := wrapErrorV3(wrappedErr)
	fmt.Println(wrappedErrV3)
	if errors.Is(wrappedErrV3, wrappedErr) {
		fmt.Println("wrappedErrV3 contains wrappedErr")
	}
}

// wrap vs sentinel errors

var ErrNotFound = errors.New("not found")

func wrapVsSentinel() {
	err := errors.New("sentinel error")

	if errors.Is(err, ErrNotFound) {
		fmt.Println("Error is ErrNotFound")
	}
}
