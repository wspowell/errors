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
//
// Example:

type Causer interface {
	~uint
}

// Error is an instance of a sentinel error. If the error was created via New() then the error is
// considered and inline error and is not comparable with another Error.
//
// Error should not be used with "==". Instead, use Error.As().
type Error[T Causer] struct {
	Cause T
}

// New error instance.
//
// This should be called when the application creates a new error.
// This new Error is not (intended to be) comparable with other Errors and therefore cannot be
// used as a sentinel error. A Sentinel may be used and compared using Error.Sentinel().
// A context is accepted in order to pass through API level feature flags.
func New[T Causer](cause T) Error[T] {
	return Error[T]{
		Cause: cause,
	}
}

// Error string representation.
//
// Satisfies golang's Error() string interface.
func (self Error[T]) Error() string {
	if asStringer, ok := any(self.Cause).(fmt.Stringer); ok {
		return asStringer.String()
	}

	if self.Cause == 0 {
		return fmt.Sprintf("%T(Ok)", self.Cause)
	}

	return fmt.Sprintf("%T(%d)", self.Cause, self.Cause)
}

// IsNone returns true if the Error is zero value.
//
// It is recommended to use Error along with Result and instead use Result.IsOk().
func (self Error[T]) IsOk() bool {
	return self.Cause == 0
}

func (self Error[T]) IsErr() bool {
	return self.Cause != 0
}

func (self Error[T]) IsSome() bool {
	return self.Cause != 0
}

func (self Error[T]) IsNone() bool {
	return self.Cause == 0
}

// Ok or no error present.
//
// Used for the success case of a function.
func Ok[T Causer]() Error[T] {
	return Error[T]{}
}
