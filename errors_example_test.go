package errors_test

import (
	"fmt"

	"github.com/wspowell/errors"
)

type ExampleErr uint64

const (
	ErrNone ExampleErr = iota
	ErrFailed
)

func OkFunc() ExampleErr {
	return ErrNone
}

func ErrFunc() ExampleErr {
	return ErrFailed
}

func ExampleOkFunc() {
	switch err := OkFunc(); err {
	case ErrFailed:
		// Handle error.
		fmt.Println(err)
	case ErrNone:
		fmt.Println("no error")
	}

	// Output:
	// no error
}

func ExampleErrFunc() {
	switch err := ErrFunc(); err {
	case ErrFailed:
		// Handle error.
		fmt.Println(err)
	case ErrNone:
		fmt.Println("no error")
	}

	// Output:
	// 1
}

type DetailedErr uint64

const (
	ErrDetailedNone DetailedErr = iota
	ErrDetailedWhoops
)

func ErrDetailedErrFn() errors.Message[DetailedErr] {
	return errors.NewMessage(ErrDetailedWhoops, "WHOOPS: %s", "specific details")
}

func ErrDetailedOkFn() errors.Message[DetailedErr] {
	return errors.Ok[DetailedErr]()
}

func ExampleErrDetailedErrFn() {
	switch err := ErrDetailedErrFn(); err.Error {
	case ErrDetailedWhoops:
		// Handle error.
		fmt.Println(err)
	case ErrDetailedNone:
		fmt.Println(err)
	}

	// Output:
	// errors_test.DetailedErr(1, WHOOPS: specific details)
}

func ExampleErrDetailedOkFn() {
	switch err := ErrDetailedOkFn(); err.Error {
	case ErrDetailedWhoops:
		// Handle error.
		fmt.Println(err)
	case ErrDetailedNone:
		fmt.Println(err)
	}

	// Output:
	// errors_test.DetailedErr(Ok)
}

type EnumErr uint64

const (
	ErrEnumNone EnumErr = iota
	ErrEnumErr1
	ErrEnumErr2
)

func (self EnumErr) String() string {
	return [...]string{
		"no error",
		"enum error 1",
		"enum error 2",
	}[self]
}

func EnumOkFn() EnumErr {
	return ErrEnumNone
}

func EnumErrFn() EnumErr {
	return ErrEnumErr2
}

func ExampleEnumErr() {
	switch err := EnumErrFn(); err {
	case ErrEnumErr1:
		// Handle error.
		fmt.Println(err)
	case ErrEnumErr2:
		// Handle error.
		fmt.Println(err)
	case ErrEnumNone:
		// No error
		fmt.Println("no error")
	}

	// Output:
	// enum error 2
}

type HandledError uint64

const (
	ErrHandledErrorOk HandledError = iota
	ErrHandledErrorNew
	ErrHandledErrorExternal
)

func HandleExternalErrors() errors.Message[HandledError] {
	switch err := ErrDetailedErrFn(); err.Error {
	case ErrDetailedWhoops:
		// Handle error.
		return errors.NewMessage(ErrHandledErrorExternal, err.Message)
	case ErrDetailedNone:
		// No error
	}

	return errors.Ok[HandledError]()
}

func ExampleHandleExternalErrors() {
	switch err := HandleExternalErrors(); err.Error {
	case ErrHandledErrorNew:
		fmt.Println(err)
	case ErrHandledErrorExternal:
		fmt.Println(err)
	case ErrHandledErrorOk:
		fmt.Println("ok")
	}

	// Output:
	// errors_test.HandledError(2, WHOOPS: specific details)
}
