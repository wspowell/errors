package result

import "github.com/wspowell/errors"

type Error = errors.Error

// Ok result. Used upon success.
func Ok[T any](result T) Result[T] {
	return Result[T]{
		value: result,
		err:   errors.Error{},
	}
}

// Err result. Used upon failure.
func Err[T any](err Error) Result[T] {
	return Result[T]{
		err: err,
	}
}

// result_ of an operation.
//
// Designed to replace the return pattern of (value, error). result_ is either a value or an error.
// Ideally, the internal err would be a non-interface type, but kept as error for backwards compatibility.
// The internal err is designed to be a singular error and not a linked list of errors. This removes a
// ton of complexity and uncertainty in the error chain and error usage/lifecycle.
type Result[T any] struct {
	value T
	err   Error
}

// Result decomposes into the basic (T, error) return value.
//
// Useful when decomposing into variables for custom evaluation.
func (self Result[T]) Result() (T, Error) {
	return self.value, self.err
}

// IsOk then return true, false otherwise.
func (self Result[T]) IsOk() bool {
	return self.err.IsNone()
}

// Value of the Ok result.
//
// Note: If called on an error result, this will be the zero value of T.
func (self Result[T]) Value() T {
	return self.value
}

// Error of the result.
func (self Result[T]) Error() Error {
	return self.err
}

// ValueOr default value if not an Ok result.
func (self Result[T]) ValueOr(defaultValue T) T {
	if self.IsOk() {
		return self.value
	}

	return defaultValue
}

// ValueOrPanic if not an Ok result.
func (self Result[T]) ValueOrPanic() T {
	if self.IsOk() {
		return self.value
	}

	panic(self.err)
}
