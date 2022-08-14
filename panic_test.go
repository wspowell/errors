package errors_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wspowell/errors"
)

const (
	expectedCatchStackTrace   = "recovered panic: WHOOPS\ngithub.com/wspowell/errors_test.panicFunction"
	expectedRecoverStackTrace = "recovered panic: WHOOPS\ngithub.com/wspowell/errors_test.recoverPanicAsError"
)

func panicFunction() {
	panic("WHOOPS")
}

func doThing(ctx context.Context) {
	panicFunction()
}

//nolint:nonamedreturns // reason: need named return to alter the Error return in the defer
func recoverPanicAsError() (err errors.Error) {
	defer func() {
		err = errors.Recover(context.Background(), recover())
	}()

	panic("WHOOPS")
}

func doThingAgain() errors.Error {
	return recoverPanicAsError()
}

func TestCatchClosure(t *testing.T) {
	t.Parallel()

	defer func() {
		if err := recover(); err != nil {
			assert.Fail(t, "expected recovered panic")
		}
	}()

	err := errors.Catch(context.Background(), func(ctx context.Context) {
		panicFunction()
	})

	assert.Contains(t, fmt.Sprintf("%+v", err), expectedCatchStackTrace)
}

func TestCatchFunction(t *testing.T) {
	t.Parallel()

	defer func() {
		if err := recover(); err != nil {
			assert.Fail(t, "expected recovered panic")
		}
	}()

	err := errors.Catch(context.Background(), doThing)

	assert.Contains(t, fmt.Sprintf("%+v", err), expectedCatchStackTrace)
}

func TestRecoverSameFunction(t *testing.T) {
	t.Parallel()

	defer func() {
		if err := recover(); err != nil {
			assert.Fail(t, "expected recovered panic")
		}
	}()

	err := recoverPanicAsError()
	assert.NotNil(t, err)
	assert.Contains(t, fmt.Sprintf("%+v", err), expectedRecoverStackTrace)
}

func TestRecoverNested(t *testing.T) {
	t.Parallel()

	defer func() {
		if err := recover(); err != nil {
			assert.Fail(t, "expected recovered panic")
		}
	}()

	err := doThingAgain()
	assert.Contains(t, fmt.Sprintf("%+v", err), expectedRecoverStackTrace)
}

func TestRecoverNone(t *testing.T) {
	t.Parallel()

	err := errors.Recover(context.Background(), nil)

	if !err.IsNone() {
		t.Error("expected error to be None")
	}

	assert.True(t, err.IsNone())
}
