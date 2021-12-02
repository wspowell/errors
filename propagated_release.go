//go:build release
// +build release

package errors

import "fmt"

// propagated implements error and provides internal code tracking and stack traces (debug only).
type propagated struct {
	internalCode string
	err          error
}

// newPropagated creates a new propagated error.
// This wraps the given error with a new internal code.
func newPropagated(internalCode string, err error) *propagated {
	return &propagated{
		internalCode: internalCode,
		err:          err,
	}
}

// Format the propagated error for different situations.
//   * %+v - Print error with internal code stack and stack trace
//   * %#v - Print error with internal code stack
//   * %v  - Print error with first internal code
func (self *propagated) Format(state fmt.State, verb rune) {
	// Determine what to prepend to the error string, if anything.
	switch verb {
	case 'v':
		if state.Flag('+') {
			// Print internal codes.
			fmt.Fprintf(state, "[%s]", self.internalCode)
		} else if state.Flag('#') {
			// Print every internal code given.
			fmt.Fprintf(state, "[%s]", self.internalCode)
		}
	}

	// Call Format for propagated errors.
	// nolint:errorlint // reason: type conversion, not an error check.
	if formatter, ok := self.err.(fmt.Formatter); ok {
		formatter.Format(state, verb)
	}
}
