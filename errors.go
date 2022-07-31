package errors

import (
	"context"
	"fmt"
)

// New error.
// This should be called when the application creates a brand new error.
// If an error has been received from an external function, use Wrap().
func New(ctx context.Context, format string, values ...any) Error {
	return newError(ctx, callersSkipError, false, format, values...)
}

// Error keeps the context of the error to print later.
//
// Error must be kept comparable.
type Error struct {
	// format to print the error string.
	err string

	// callStack only when isDebug.
	callStack *stack
}

func newError(ctx context.Context, callersSkipCount int, isPanic bool, format string, values ...any) Error {
	var callStack *stack

	if isPanic || shouldPrintStackTrace(ctx) {
		callStack = callers(callersSkipCount)
	}

	var err string
	if len(values) == 0 {
		// Do not call fmt.Sprintf() if not necessary.
		// Major performance improvement.
		err = format
	} else {
		err = fmt.Sprintf(format, values...)
	}

	return Error{
		err: err,
		//formatValues: formatValues,
		callStack: callStack,
	}
}

func (self Error) String() string {
	if self.callStack == nil {
		// Do not call fmt.Sprintf() if not necessary.
		// Major performance improvement.
		return self.err
	}

	return fmt.Sprintf("%s", self)
}

func (self Error) Format(state fmt.State, verb rune) {
	fmt.Fprintf(state, self.err)

	if self.callStack != nil {
		// Print stack trace.
		fmt.Fprintf(state, "%+v", self.callStack)
	}
}

func (self Error) IsNone() bool {
	// nolint:exhaustruct // reason: zero value for Error is desired
	return self == Error{}
}
