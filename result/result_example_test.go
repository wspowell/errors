package result_test

import (
	"context"
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
	exampleFn := func(ctx context.Context, success bool) Result[int, Error] {
		if success {
			return Ok[int, Error](1)
		}

		return Err[int](errors.New(ctx, errErrorFailure))
	}

	ctx := context.Background()

	// Success result.
	fmt.Println(exampleFn(ctx, true).IsOk())
	fmt.Println(exampleFn(ctx, true).Error().Error())
	fmt.Println(exampleFn(ctx, true).Value())
	fmt.Println(exampleFn(ctx, true).ValueOr(2))
	fmt.Println(exampleFn(ctx, true).ValueOrPanic())

	// Error result.
	fmt.Println(exampleFn(ctx, false).IsOk())
	fmt.Println(exampleFn(ctx, false).Error().Error())
	fmt.Println(exampleFn(ctx, false).Value())
	fmt.Println(exampleFn(ctx, false).ValueOr(2))
	//fmt.Println(exampleFn(ctx, false).ValueOrPanic()) // Will panic

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
