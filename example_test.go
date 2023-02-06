package errors_test

import (
	"fmt"

	"github.com/wspowell/errors"
)

type ExampleFnError uint

const (
	ExampleFnErrorMyBad = ExampleFnError(iota + 1)
	ExampleFnErrorInternalFailure
)

func (self ExampleFnError) String() string {
	switch self {
	case ExampleFnErrorMyBad:
		return "MyBad"
	case ExampleFnErrorInternalFailure:
		return "InternalFailure"
	}

	return "Ok"
}

func exampleFn(option int) errors.Error[ExampleFnError] {
	switch option {
	case 1:
		return errors.New(ExampleFnErrorMyBad)
	case 2:
		return errors.New(ExampleFnErrorInternalFailure)
	default:
		return errors.Ok[ExampleFnError]()
	}
}

func ExampleOk() {
	switch err := exampleFn(0); err.Cause {
	case ExampleFnErrorMyBad:
		fmt.Println(err)
	case ExampleFnErrorInternalFailure:
		fmt.Println(err)
	default:
		fmt.Println("Ok")
	}

	// Output:
	// Ok
}

func ExampleExampleFnErrorMyBad() {
	switch err := exampleFn(1); err.Cause {
	case ExampleFnErrorMyBad:
		fmt.Println(err)
	case ExampleFnErrorInternalFailure:
		fmt.Println(err)
	default:
		fmt.Println("Ok")
	}

	// Output:
	// MyBad
}

func ExampleExampleFnErrorInternalFailure() {
	switch err := exampleFn(2); err.Cause {
	case ExampleFnErrorMyBad:
		fmt.Println(err)
	case ExampleFnErrorInternalFailure:
		fmt.Println(err)
	default:
		fmt.Println("Ok")
	}

	// Output:
	// InternalFailure
}
