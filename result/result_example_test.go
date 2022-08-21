package result_test

import (
	"fmt"

	"github.com/wspowell/errors"
	"github.com/wspowell/errors/result"
)

type Result[T any, E result.Optional] struct {
	result.Result[T, E]
}

func Ok[T any, E result.Optional](value T) Result[T, E] {
	return Result[T, E]{result.Ok[T, E](value)}
}

func Err[T any, E result.Optional](err E) Result[T, E] {
	return Result[T, E]{result.Err[T](err)}
}

func ExampleResult() {
	exampleFn := func(success bool) Result[int, errors.Error[string]] {
		if success {
			return Ok[int, errors.Error[string]](1)
		}

		return Err[int](errors.Some(errErrorFailure))
	}

	// Success result.
	fmt.Println(exampleFn(true).IsOk())
	fmt.Println(exampleFn(true).Error().Into())
	fmt.Println(exampleFn(true).Value())
	fmt.Println(exampleFn(true).ValueOr(2))
	fmt.Println(exampleFn(true).ValueOrPanic())

	// Error result.
	fmt.Println(exampleFn(false).IsOk())
	fmt.Println(exampleFn(false).Error().Into())
	fmt.Println(exampleFn(false).Value())
	fmt.Println(exampleFn(false).ValueOr(2))
	//fmt.Println(exampleFn(false).ValueOrPanic()) // Will panic

	// Output:
	// true
	//
	// 1
	// 1
	// 1
	// false
	// failure
	// 0
	// 2
}
