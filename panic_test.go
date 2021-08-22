package errors_test

import (
	"testing"

	"github.com/wspowell/errors"

	"github.com/stretchr/testify/assert"
)

const (
	expectedCatchStackTrace   = "[PANIC] recovered panic: expected\ngithub.com/wspowell/errors_test.panicFunction"
	expectedRecoverStackTrace = "[PANIC] recovered panic: expected\ngithub.com/wspowell/errors_test.recoverPanicAsError"
)

func panicFunction() {
	panic("expected")
}

func doThing() {
	panicFunction()
}

func recoverPanicAsError() (err error) {
	defer func() {
		err = errors.Recover(recover())
	}()

	panic("expected")
}

func doThingAgain() error {
	return recoverPanicAsError()
}

func Test_Catch_closure(t *testing.T) {
	t.Parallel()

	defer func() {
		if err := recover(); err != nil {
			assert.Fail(t, "expected recovered panic")
		}
	}()

	err := errors.Catch(func() {
		panicFunction()
	})

	assert.ErrorIs(t, err, errors.Panic)
	assert.Contains(t, err.Error(), expectedCatchStackTrace)
}

func Test_Catch_function(t *testing.T) {
	t.Parallel()

	defer func() {
		if err := recover(); err != nil {
			assert.Fail(t, "expected recovered panic")
		}
	}()

	err := errors.Catch(doThing)

	assert.ErrorIs(t, err, errors.Panic)
	assert.Contains(t, err.Error(), expectedCatchStackTrace)
}

func Test_Recover_same_function(t *testing.T) {
	t.Parallel()

	defer func() {
		if err := recover(); err != nil {
			assert.Fail(t, "expected recovered panic")
		}
	}()

	err := recoverPanicAsError()
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), expectedRecoverStackTrace)
}

func Test_Recover_nested(t *testing.T) {
	t.Parallel()

	defer func() {
		if err := recover(); err != nil {
			assert.Fail(t, "expected recovered panic")
		}
	}()

	err := doThingAgain()
	assert.ErrorIs(t, err, errors.Panic)
	assert.Contains(t, err.Error(), expectedRecoverStackTrace)
}
