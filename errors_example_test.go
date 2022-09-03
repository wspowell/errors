package errors_test

import (
	"fmt"

	"github.com/wspowell/errors"
)

const (
	ErrSome = "WHOOPS"
)

func OkFunc() errors.Error[string] {
	return errors.None[string]()
}

func ErrFunc() errors.Error[string] {
	return errors.New(ErrSome)
}

func ExampleNone() {
	if err := OkFunc(); err.IsSome() {
		// Handle error.
		fmt.Println(err)

		return
	}

	fmt.Println("no error")

	// Output:
	// no error
}

func ExampleNew() {
	if err := ErrFunc(); err.IsSome() {
		// Handle error.
		fmt.Println(err)

		return
	}

	fmt.Println("no error")

	// Output:
	// WHOOPS
}

type EnumErr string

const (
	ErrEnumErr1 = EnumErr("enum error 1")
	ErrEnumErr2 = EnumErr("enum error 2")
)

func EnumOkFn() errors.Error[EnumErr] {
	return errors.None[EnumErr]()
}

func EnumErrFn() errors.Error[EnumErr] {
	return errors.New(ErrEnumErr2)
}

func ExampleEnumErrFn() {
	switch err := EnumErrFn(); err.Into() {
	case ErrEnumErr1:
		// Handle error.
		fmt.Println(err)

		return
	case ErrEnumErr2:
		// Handle error.
		fmt.Println(err)

		return
	}

	fmt.Println("no error")

	// Output:
	// enum error 2
}

func ExampleEnumOkFn() {
	switch err := EnumOkFn(); err.Into() {
	case ErrEnumErr1:
		// Handle error.
		fmt.Println(err)

		return
	case ErrEnumErr2:
		// Handle error.
		fmt.Println(err)

		return
	}

	fmt.Println("no error")

	// Output:
	// no error
}
