package errors

import (
	"fmt"
	"runtime"
)

const (
	icPanic = "PANIC"
)

var (
	ErrPanic = New(icPanic, "recovered panic")
)

func Recover(errPanic interface{}) error {
	if errPanic != nil {
		// nolint:wrapcheck // reason: passthrough for handling panic error
		return fmt.Errorf("%w: %v%+v", ErrPanic, errPanic, panicCallers())
	}

	return nil
}

// Catch potential panics that occur in a function call.
// Panics should never occur, therefore stack traces print regardless of build mode (release or debug).
func Catch(fn func()) (err error) {
	defer func() {
		err = Recover(recover())
	}()

	fn()

	return nil
}

func panicCallers() *stack {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(5, pcs[:])
	var st stack = pcs[0:n]

	return &st
}
