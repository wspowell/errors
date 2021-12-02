//go:build release
// +build release

package errors

import (
	"fmt"
)

type cause struct {
	internalCode string
	format       string
	values       []interface{}

	// Optional. May be present if an error was propagated or converted.
	fromErr error
	// Optional. May be present if an error was converted.
	toErr error
}

func newCause(internalCode string, format string, values ...interface{}) *cause {
	return &cause{
		internalCode: internalCode,
		format:       format,
		values:       values,
	}
}

func newCauseWithError(internalCode string, err error) *cause {
	return &cause{
		internalCode: internalCode,
		format:       "%s",
		values:       []interface{}{err},
		fromErr:      err,
	}
}

func newCauseWithErrorConversion(internalCode string, fromErr error, toErr error) *cause {
	return &cause{
		internalCode: internalCode,
		format:       "%s",
		values:       []interface{}{toErr},
		fromErr:      fromErr,
		toErr:        toErr,
	}
}

func (self *cause) Unwrap() error {
	if self.toErr != nil {
		return self.toErr
	}

	return self.fromErr
}

func (self *cause) Error() string {
	return fmt.Sprintf(self.format, self.values...)
}

func (self *cause) Format(state fmt.State, verb rune) {
	switch verb {
	case 's':
		fmt.Fprintf(state, self.format, self.values...)
	case 'v':
		fmt.Fprintf(state, "[%s] ", self.internalCode)
		self.Format(state, 's')

		if self.toErr != nil {
			// Print the converted errors.
			if state.Flag('+') {
				fmt.Fprintf(state, "\n\n")
			} else if state.Flag('#') {
				fmt.Fprintf(state, " -> ")
			} else {
				// Do not print the converted from error.
				break
			}

			// nolint:errorlint // reason: type conversion, not an error check.
			if formatter, ok := self.fromErr.(fmt.Formatter); ok {
				formatter.Format(state, verb)
			}
		}
	default:
		self.Format(state, 's')
	}
}
