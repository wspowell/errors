package errors

import (
	"fmt"
)

const ErrNone = ""

// Standard error is an alias for a string typed error.
type Standard = Error[string]

// Error defined as a message of type string whose actual error
// type is type T. T is a type whose underlying type is string.
//
// Errors may be generic as Error[string] or may have an explicit
// type. Usually, this is an enum type that allows Error[T] to
// be self documenting and is used in conjunction with Into() and
// the "exhaustive" linter.
type Error[T ~string] string

// New error instance.
//
// Creates a new error instance that is of error type T.
// Remember, the type T, not the error instance, is what
// determines what kind of error this represents.
func New[T ~string](format T, values ...any) Error[T] {
	return newError(format, values...)
}

// None error instance.
//
// Creates an empty Error that IsNone().
func None[T ~string]() Error[T] {
	return Error[T](ErrNone)
}

// newError instance.
func newError[T ~string](format T, values ...any) Error[T] {
	var err string
	if len(values) != 0 {
		// Do not call fmt.Sprintf() if not necessary.
		// Major performance improvement.
		// Only necessary if there are any values.
		err = fmt.Sprintf(string(format), values...)
	} else {
		err = string(format)
	}

	return Error[T](err)
}

// IsNone return true when this error instance represents the
// absence of an error.
func (self Error[T]) IsNone() bool {
	return self == ErrNone
}

// IsSome return true when this error instance represents the
// presence of an error.
//
// Should be used when T is string or when checking if an
// error exists. Should not be used when checking errors when
// T is not string.
func (self Error[T]) IsSome() bool {
	return self != ErrNone
}

// Into the error type.
//
// Should be used when T is not a string in conjunction with switch.
func (self Error[T]) Into() T {
	return T(self)
}
