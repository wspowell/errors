package errors

import (
	"fmt"
	"strings"
)

// internalError implements error and provides internal code tracking and stack traces.
type internalError struct {
	internalCode string
	err          error
}

// newInternalError creates a new internal error.
// This wraps the given error.
func newInternalError(internalCode string, err error) *internalError {
	return &internalError{
		internalCode: internalCode,
		err:          err,
	}
}

// Error string of the origin error.
// This does not include the internal code.
func (self *internalError) Error() string {
	return self.err.Error()
}

// internalCodeStack starting with first (cause) internal code to most recent.
func (self *internalError) internalCodeStack() string {
	var stack string

	recurseErrorStack(self, func(err error) {
		if asInternalError, ok := err.(*internalError); ok {
			if stack == "" {
				stack = asInternalError.internalCode
			} else {
				stack = asInternalError.internalCode + "," + stack
			}
		}
	})

	return stack
}

// first internal error
func (self *internalError) first() *internalError {
	var firstError *internalError

	recurseErrorStack(self, func(err error) {
		if asInternalError, ok := err.(*internalError); ok {
			firstError = asInternalError
		}
	})

	return firstError
}

// Format the internal error for different situations.
//   * %+v - Print error with internal code stack and stack trace
//   * %#v - Print error with internal code stack
//   * %v  - Print error with first internal code
func (self *internalError) Format(s fmt.State, verb rune) {
	// Determine what to prepend to the error string, if anything.
	switch verb {
	case 'v':
		if s.Flag('+') {
			// Print every internal code given, plus stack trace.
			fmt.Fprintf(s, "[%s] ", self.internalCodeStack())
		} else if s.Flag('#') {
			// Print every internal code given.
			fmt.Fprintf(s, "[%s] ", self.internalCodeStack())
		} else {
			// Print only the origin internal code.
			stack := self.internalCodeStack()
			fmt.Fprintf(s, "[%s] ", strings.Split(stack, ",")[0])
		}
	}

	// Call Format for wrapped errors.
	self.first().err.(fmt.Formatter).Format(s, verb)
}

// Unwrap to get the underlying error.
func (self *internalError) Unwrap() error {
	// We do not care about the internal error so unwrap it and return that value.
	return Unwrap(self.err)
}
