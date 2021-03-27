// +build release

package errors

import "fmt"

// wrapped implements error and provides internal code tracking and stack traces.
type wrapped struct {
	internalCode string
	err          error
}

// newInternalError creates a new internal error.
// This wraps the given error.
func newWrapped(internalCode string, err error) *wrapped {
	return &wrapped{
		internalCode: internalCode,
		err:          err,
	}
}

// Format the internal error for different situations.
//   * %+v - Print error with internal code stack and stack trace
//   * %#v - Print error with internal code stack
//   * %v  - Print error with first internal code
func (self *wrapped) Format(state fmt.State, verb rune) {
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

	// Call Format for wrapped errors.
	self.err.(fmt.Formatter).Format(state, verb)
}
