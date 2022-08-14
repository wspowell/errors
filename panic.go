package errors

import "context"

func Recover(ctx context.Context, errPanic any) Error[string] {
	if errPanic != nil {
		ctx = WithStackTrace(ctx)

		return newError(ctx, callersSkipPanic, "recovered panic: %s", errPanic)
	}

	//nolint:exhaustruct // reason: zero value for Error is desired
	return None[string]()
}

// Catch potential panics that occur in a function call.
// Panics should never occur, therefore stack traces print regardless of build mode (release or debug).
//nolint:nonamedreturns // reason: need named return to alter the Error return in the defer
func Catch(ctx context.Context, fn func(ctx context.Context)) (err Error[string]) {
	defer func(ctx context.Context) {
		err = Recover(ctx, recover())
	}(ctx)

	fn(ctx)

	//nolint:exhaustruct // reason: zero value for Error is desired
	return None[string]()
}
