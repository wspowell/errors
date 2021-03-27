package errors

import (
	"fmt"

	pkgerrors "github.com/pkg/errors"
)

// New error.
// This should be called when the application creates a brand new error.
// If an error has been received from an external function or is propogating an error, use Wrap().
func New(internalCode string, format string, values ...interface{}) error {
	return newInternalError(internalCode, pkgerrors.New(fmt.Sprintf(format, values...)))
}

// Wrap an existing error.
// The internalCode should be a unique code to allow developers to easily identify the source of an issue.
func Wrap(internalCode string, err error) error {
	return newInternalError(internalCode, pkgerrors.WithStack(err))
}

// Is reports whether any error in err's chain matches target.
// See: https://golang.org/pkg/errors/#Is
func Is(err, target error) bool {
	return pkgerrors.Is(err, target)
}

// As finds the first error in err's chain that matches target,
// and if so, sets target to that error value and returns true.
// Otherwise, it returns false.
// See: https://golang.org/pkg/errors/#As
func As(err error, target interface{}) bool {
	return pkgerrors.As(err, target)
}

// Unwrap returns the result of calling the Unwrap method on err,
// if err's type contains an Unwrap method returning error.
// Otherwise, Unwrap returns nil.
// See: https://golang.org/pkg/errors/#Unwrap
func Unwrap(err error) error {
	return pkgerrors.Unwrap(err)
}

// Cause of the error.
// This returns the very first error encoutered whether that was
// a new application error or an external error.
func Cause(err error) error {
	var next error

	for err != nil {
		if next = Unwrap(err); next == nil {
			break
		}
		err = next
	}

	return err
}
