package errors

import (
	"context"
	"fmt"
	"io"
)

const None = ""

var ErrNone = Error{err: None}

// New error instance.
//
// This should be called when the application creates a new error.
// This new Error is not (intended to be) comparable with other Errors and therefore cannot be
// used as a sentinel error. A Sentinel may be used and compared using Error.Sentinel().
// A context is accepted in order to pass through API level feature flags.
func New[T ~string](ctx context.Context, format T, values ...any) Error {
	return newError(ctx, callersSkipError, string(format), values...)
}

// Error is an instance of a sentinel error. If the error was created via New() then the error is
// considered and inline error and is not comparable with another Error.
//
// Error should not be used with "==". Instead, use Error.As().
type Error struct {
	// format to print the error string.
	err string

	// callStack only when isDebug.
	callStack *stack
}

func newError(ctx context.Context, callersSkipCount int, err string, values ...any) Error {
	if len(values) != 0 {
		// Do not call fmt.Sprintf() if not necessary.
		// Major performance improvement.
		// Only necessary if there are any values.
		err = fmt.Sprintf(err, values...)
	}

	var callStack *stack
	if shouldPrintStackTrace(ctx) {
		callStack = callers(callersSkipCount)
	}

	return Error{
		err:       err,
		callStack: callStack,
	}
}

// Error string, ignoring any call stack.
func (self Error) Error() string {
	return self.err
}

func (self Error) Format(state fmt.State, verb rune) {
	if _, err := io.WriteString(state, self.err); err != nil {
		fmt.Fprint(state, "<failed formatting error>")
	}

	if verb == 'v' && self.callStack != nil && state.Flag('+') {
		fmt.Fprintf(state, "%+v", self.callStack)
	}
}

// IsNone returns true if the Error is zero value.
//
// It is recommended to use Error along with Result and instead use Result.IsOk().
func (self Error) IsNone() bool {
	switch As[string](self) {
	case None:
		return true
	default:
		return false
	}
}

// As the sentinel type for the error.
//
// This allows Error to be used in a switch in conjunction with linter "exhaustive".
//
// For example:
//
//  type myError errors.Sentinel
//
//  const (
//    errOne   = myError("one")
//    errTwo   = myError("two")
//    errThree = myError("three")
//  )
//
//  func multipleErrors(ctx context.Context, err int) errors.Error {
//    switch err {
//    case 1:
//      return errors.New(ctx, errOne)
//    case 2:
//      return errors.New(ctx, errTwo)
//    case 3:
//      return errors.New(ctx, errThree)
//    default:
//      return errors.ErrNone
//    }
//  }
//
//  func TestSentinelEnum(t *testing.T) {
//    t.Parallel()
//
//    err := multipleErrors(context.Background(), 1)
//    //nolint:exhaustive // reason: expected lint error for testing
//    switch errors.As[myError](err) {
//    case errOne:
//      // Expected case.
//    default:
//      assert.Fail(t, "expected errOne")
//    }
//  }
func As[T ~string](err Error) T {
	return T(err.err)
}
