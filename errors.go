package errors

import (

	// nolint:depguard // reason: import to shadow the package
	goerrors "errors"
)

// New error.
// This should be called when the application creates a brand new error.
// If an error has been received from an external function or is propogating an error, use Propagate().
func New(internalCode string, format string, values ...interface{}) error {
	return newCause(internalCode, format, values...)
}

// Propagate an existing error.
// The internalCode should be a unique code to allow developers to easily identify the source of an issue.
func Propagate(internalCode string, err error) error {
	causeErr := &cause{}
	if As(err, &causeErr) {
		return newPropagated(internalCode, err)
	}

	return newCauseWithError(internalCode, err)
}

// Convert an existing error to a new error.
// Calls to As(), Is(), InternalCode(), Cause(), and Unwrap() will only refer to the new error.
// Original error used for formatting/printing only.
// This may be useful when a concrete error should be returned from a function to enable consumers
//   to detect this error with As() or Is().
func Convert(internalCode string, fromErr error, toErr error) error {
	if toErr == nil {
		return Propagate(internalCode, fromErr)
	}
	if fromErr == nil {
		return Propagate(internalCode, toErr)
	}

	return newCauseWithErrorConversion(internalCode, fromErr, toErr)
}

// Is reports whether any error in err's chain matches target.
// See: https://golang.org/pkg/errors/#Is
func Is(err, target error) bool {
	return goerrors.Is(err, target)
}

// As finds the first error in err's chain that matches target,
// and if so, sets target to that error value and returns true.
// Otherwise, it returns false.
// See: https://golang.org/pkg/errors/#As
func As(err error, target interface{}) bool {
	return goerrors.As(err, target)
}

// Unwrap returns the result of calling the Unwrap method on err,
// if err's type contains an Unwrap method returning error.
// Otherwise, Unwrap returns nil.
// See: https://golang.org/pkg/errors/#Unwrap
func Unwrap(err error) error {
	// nolint:wrapcheck // reason: passthrough for shadowing errors package
	return goerrors.Unwrap(err)
}

// InternalCode of the first error created or converted.
// If err does not have an internal code then return empty string.
func InternalCode(err error) string {
	var internalCode string

	recurseErrorStack(err, func(err error) {
		// nolint:errorlint // reason: recurse over each error explicitly
		if asCause, ok := err.(*cause); ok {
			internalCode = asCause.internalCode
		}
	})

	return internalCode
}

// Cause of the error.
// This returns the very first error encoutered whether that was
// a new application error or an external error.
func Cause(err error) error {
	var firstErr error
	var causeErr *cause

	recurseErrorStack(err, func(err error) {
		// nolint:errorlint // reason: recurse over each error explicitly
		if asCause, ok := err.(*cause); ok {
			if causeErr == nil || causeErr.toErr == nil {
				causeErr = asCause
			}
		}
		firstErr = err
	})

	if causeErr != nil {
		if causeErr.toErr != nil {
			return causeErr.toErr
		}
		if causeErr.fromErr != nil {
			return causeErr.fromErr
		}

		return causeErr
	}

	return firstErr
}

func recurseErrorStack(err error, processFn func(error)) {
	var next error
	for err != nil {
		processFn(err)

		if next = Unwrap(err); next == nil {
			break
		}
		err = next
	}
}
