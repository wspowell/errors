package errors

import (
	// nolint:depguard // reason: import to shadow the package
	"errors"
)

// New error.
// This should be called when the application creates a brand new error.
// If an error has been received from an external function, use Wrap().
func New(format string, values ...any) error {
	return newCause(format, values...)
}

// Wrap an existing error with a new error.
// Calls to As(), Is(), InternalCode(), Cause(), and Unwrap() will only refer to the new error.
// Original error used for formatting/printing only.
// This should be used when returning errors from an external package. This makes the new error as the cause.
// This can also be used to return a new discret error so that the caller is not required
//   to know the underlying implementation. For example, a function may use the io package and therefore handle
//   io.EOF and other errors, but the function could wrap those in a new error that could be checked by callers.
func Wrap(fromErr error, toErr error) error {
	if toErr == nil {
		return newCauseWithWrappedError(nil, fromErr)
	}
	if fromErr == nil {
		return newCauseWithWrappedError(nil, toErr)
	}

	return newCauseWithWrappedError(fromErr, toErr)
}

// Is reports whether any error in err's chain matches target.
// See: https://golang.org/pkg/errors/#Is
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// As finds the first error in err's chain that matches target,
// and if so, sets target to that error value and returns true.
// Otherwise, it returns false.
// See: https://golang.org/pkg/errors/#As
func As(err error, target any) bool {
	return errors.As(err, target)
}

// Unwrap returns the result of calling the Unwrap method on err,
// if err's type contains an Unwrap method returning error.
// Otherwise, Unwrap returns nil.
// See: https://golang.org/pkg/errors/#Unwrap
func Unwrap(err error) error {
	// nolint:wrapcheck // reason: passthrough for shadowing errors package
	return errors.Unwrap(err)
}

// Cause of the error.
// This returns the first "cause" error encountered.
// If there is no "cause" error, then it returns the very first error in the Unwrap() chain.
func Cause(err error) error {
	causeErr := &cause{}
	if errors.As(err, &causeErr) {
		if causeErr.toErr != nil {
			return causeErr.toErr
		}
		return causeErr
	}

	return findFirstError(err)
}

func findFirstError(err error) error {
	if nextErr := errors.Unwrap(err); nextErr != nil {
		return findFirstError(nextErr)
	}
	return err
}
