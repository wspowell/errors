package errors

import (
	"fmt"
	"strconv"
)

// Ok error.
//
// By convention, any error value that is 0 is considered "none" or "no error".
// Therefore, the zero value of T must be 0.
func Ok[T ~uint64]() Message[T] {
	return Message[T]{}
}

type Message[T ~uint64] struct {
	Error   T
	Message string
}

func NewMessage[T ~uint64](err T, message string, values ...any) Message[T] {
	if err == 0 {
		panic("By convention, error instances must not be 0. Use a 'None'/'Ok' value for iota 0 and errors.Ok() for success returns.")
	}

	if len(values) != 0 {
		// Do not call fmt.Sprintf() if not necessary.
		// Major performance improvement.
		// Only necessary if there are any values.
		message = fmt.Sprintf(message, values...)
	}

	return Message[T]{
		Error:   err,
		Message: message,
	}
}

// Into the error type.
//
// Should be used when T is not a string in conjunction with switch.
func (self Message[T]) String() string {
	if self.Error == 0 {
		return fmt.Sprintf("%T", self.Error) + "(Ok)"
	}

	return fmt.Sprintf("%T", self.Error) + "(" + strconv.FormatUint(uint64(self.Error), 10) + ", " + self.Message + ")"
}
