package result

import "fmt"

type Optional interface {
	IsNone() bool
}

// Ok result. Used upon success.
func Ok[T any, E Optional](result T) Result[T, E] {
	return Result[T, E]{
		value: result,
	}
}

// Err result. Used upon failure.
func Err[T any, E Optional](err E) Result[T, E] {
	return Result[T, E]{
		err: err,
	}
}

// result_ of an operation.
//
// Designed to replace the return pattern of (value, error). result_ is either a value or an error.
// Ideally, the internal err would be a non-interface type, but kept as error for backwards compatibility.
// The internal err is designed to be a singular error and not a linked list of errors. This removes a
// ton of complexity and uncertainty in the error chain and error usage/lifecycle.
type Result[T any, E Optional] struct {
	value T
	err   E
}

// Result decomposes into the basic (T, error) return value.
//
// Useful when decomposing into variables for custom evaluation.
func (self Result[T, E]) Result() (T, E) {
	return self.value, self.err
}

// IsOk then return true, false otherwise.
func (self Result[T, E]) IsOk() bool {
	return self.err.IsNone()
}

// Value of the Ok result.
//
// Note: If called on an error result, this will be the zero value of T.
func (self Result[T, E]) Value() T {
	return self.value
}

// Error of the result.
//
// The Error value can never be `nil`, so it is recommended to use Result.IsOk().
// Another alternative is to call Result.Error().IsNone().
func (self Result[T, E]) Error() E {
	return self.err
}

// ValueOr default value if not an Ok result.
func (self Result[T, E]) ValueOr(defaultValue T) T {
	if self.IsOk() {
		return self.value
	}

	return defaultValue
}

// ValueOrPanic if not an Ok result.
//
// It is recommended to only call this during app initialization.
// Otherwise, use Result.ValueOr().
func (self Result[T, E]) ValueOrPanic() T {
	if self.IsOk() {
		return self.value
	}

	panic(fmt.Sprintf("%+v", self.err))
}
