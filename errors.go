package errors

import (
	"fmt"
)

// Cause of an error.
//
// This is the reason for the error. It should be based on a typed enum that provides the
// ability to use a switch to handle the specific Cause of and Error. Pairing this with
// golangci-lint (exhaustive) allows developers to always be sure that all error cases
// have been handled properly.
//
// Causes allow a function to self document all non-success cases they may return.
// This should be provided in documentation to a consumer of the function and allows
// the consumer to understand how to properly utilize the function without ever being
// required to peer into its implementation.
//
// Enum types should follow a pattern of: type <Name>Error uint
// Enum values should follow a pattern of: <Name>Error<FailureCase>
// Enum values MUST start at 1 and can utilize: <Name>Error(iota + 1)
type Causer interface {
	~uint
}

// Error instance whose Cause is T.
//
// An Error is only storage for context for the Cause that triggered the error.
type Error[T Causer] struct {
	Cause T
}

// New Error instance of a given Cause.
func New[T Causer](cause T) Error[T] {
	return Error[T]{
		Cause: cause,
	}
}

// Error string representation.
//
// Satisfies golang's Error() string interface.
//
// For best performance, implement Stringer for the Cause type.
// Not implementing Stringer will require reflection via fmt.Sprintf().
func (self Error[T]) Error() string {
	if asStringer, ok := any(self.Cause).(fmt.Stringer); ok {
		return asStringer.String()
	}

	if self.Cause == 0 {
		return fmt.Sprintf("%T(Ok)", self.Cause)
	}

	return fmt.Sprintf("%T(%d)", self.Cause, self.Cause)
}

// IsOk returns true if this is an Ok Error instance.
//
// See: Ok()
func (self Error[T]) IsOk() bool {
	return self.Cause == 0
}

// IsErr returns true if this is a non-zero Error instance.
//
// See: New()
func (self Error[T]) IsErr() bool {
	return self.Cause != 0
}

// Ok or no error present.
//
// Used for the success case of a function.
func Ok[T Causer]() Error[T] {
	return Error[T]{}
}
