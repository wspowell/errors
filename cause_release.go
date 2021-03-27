// +build release

package errors

import (
	"fmt"
)

type cause struct {
	internalCode string
	format       string
	values       []interface{}

	// Optional. May be present if an error was wrapped.
	err error
}

func newCause(internalCode string, format string, values ...interface{}) *cause {
	return &cause{
		internalCode: internalCode,
		format:       format,
		values:       values,
		err:          nil,
	}
}

func newCauseWithError(internalCode string, err error) *cause {
	return &cause{
		internalCode: internalCode,
		format:       "%s",
		values:       []interface{}{err},
		err:          err,
	}
}

func (self *cause) Unwrap() error {
	return self.err
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
	default:
		self.Format(state, 's')
	}
}
