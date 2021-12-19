package errors

import (
	"fmt"
)

var (
	ErrPanic = New("recovered panic")
)

func Recover(errPanic any) error {
	if errPanic != nil {
		// nolint:wrapcheck // reason: passthrough for handling panic error
		return fmt.Errorf("%w: %v%+v", ErrPanic, errPanic, callers(callersSkipPanic))
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
