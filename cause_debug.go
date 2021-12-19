//go:build !release
// +build !release

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

	// Debug only.
	callStack *stack
}

func newCause(format string, values ...any) *cause {
	return &cause{
		format: format,
		values: values,

		callStack: callers(callersSkipFunction),
	}
}

func newCauseWithWrappedError(fromErr error, toErr error) *cause {
	return &cause{
		format:  "%s",
		values:  []any{toErr},
		fromErr: fromErr,
		toErr:   toErr,

		callStack: callers(callersSkipFunction),
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

		if state.Flag('+') {
			if self.toErr != nil && self.fromErr != nil {
				fmt.Fprintf(state, " -> %s", self.fromErr)
			}

			// Print stack trace.
			self.callStack.Format(state, verb)

			if self.toErr != nil && self.fromErr != nil {
				fmt.Fprintf(state, "\n\n")

				// nolint:errorlint // reason: type conversion, not an error check.
				if formatter, ok := self.fromErr.(fmt.Formatter); ok {
					formatter.Format(state, verb)
				} else {
					fmt.Fprintf(state, "%s\n(no stack trace available)", self.fromErr)
				}
			}
		} else if state.Flag('#') {
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
