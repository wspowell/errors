package result_test

import (
	"fmt"

	"github.com/wspowell/errors/result"
)

func ExampleResult() {
	exampleFn := func(success bool) result.Result[int, FailureError] {
		if success {
			return result.Ok[int, FailureError](1)
		}

		return result.Err[int](errFailureError)
	}

	// Success result.
	fmt.Println(exampleFn(true).IsOk())
	fmt.Println(exampleFn(true).Error())
	fmt.Println(exampleFn(true).Value())
	fmt.Println(exampleFn(true).ValueOr(2))
	fmt.Println(exampleFn(true).ValueOrPanic())

	// Error result.
	fmt.Println(exampleFn(false).IsOk())
	fmt.Println(exampleFn(false).Error())
	fmt.Println(exampleFn(false).Value())
	fmt.Println(exampleFn(false).ValueOr(2))
	//fmt.Println(exampleFn(false).ValueOrPanic()) // Will panic

	// Output:
	// true
	// none
	// 1
	// 1
	// 1
	// false
	// failure
	// 0
	// 2
}
