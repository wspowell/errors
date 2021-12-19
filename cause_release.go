//go:build release
// +build release

package errors

import (
	"fmt"
)

type cause struct {
	format string
	values []any

	// Optional. May be present if an error was wrapped.
	fromErr error
	// Optional. May be present if an error was wrapped.
	toErr error
}

func newCause(format string, values ...any) *cause {
	return &cause{
		format: format,
		values: values,
	}
}

func newCauseWithWrappedError(fromErr error, toErr error) *cause {
	return &cause{
		format:  "%s",
		values:  []any{toErr},
		fromErr: fromErr,
		toErr:   toErr,
	}
}

func (self *cause) Unwrap() error {
	if self.toErr != nil {
		return self.toErr
	}

	return nil
}

func (self *cause) Error() string {
	return fmt.Sprintf(self.format, self.values...)
}

func (self *cause) Format(state fmt.State, verb rune) {
	switch verb {
	case 's':
		fmt.Fprintf(state, "%s", self.Error())
	case 'v':
		self.Format(state, 's')

		if state.Flag('+') || state.Flag('#') {
			if self.toErr != nil && self.fromErr != nil {
				// Print the wrapped error string.
				fmt.Fprintf(state, " -> ")

				// nolint:errorlint // reason: type conversion, not an error check.
				if formatter, ok := self.fromErr.(fmt.Formatter); ok {
					formatter.Format(state, verb)
				} else {
					fmt.Fprintf(state, "%s", self.fromErr)
				}
			}
		}
	default:
		self.Format(state, 's')
	}
}
