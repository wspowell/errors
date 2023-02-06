package errors_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wspowell/errors"
)

type TestError uint

const (
	TestErrorMyBad = TestError(iota + 1)
	TestErrorInternalFailure
)

func (self TestError) String() string {
	switch self {
	case TestErrorMyBad:
		return "MyBad"
	case TestErrorInternalFailure:
		return "InternalFailure"
	}

	return "Ok"
}

func runFn(option int) errors.Error[TestError] {
	switch option {
	case 1:
		return errors.New(TestErrorMyBad)
	case 2:
		return errors.New(TestErrorInternalFailure)
	default:
		return errors.Ok[TestError]()
	}
}

func TestNewError(t *testing.T) {
	err := errors.New(TestErrorInternalFailure)
	assert.Implements(t, (*error)(nil), err)
	assert.True(t, err.IsErr())
	assert.False(t, err.IsOk())
	assert.Equal(t, TestError(2), err.Cause)
	assert.NotEqual(t, NoStringerError(2), err.Cause)
	assert.Equal(t, "InternalFailure", err.Error())
}

func TestOk(t *testing.T) {
	err := errors.Ok[TestError]()
	assert.Implements(t, (*error)(nil), err)
	assert.False(t, err.IsErr())
	assert.True(t, err.IsOk())
	assert.Equal(t, TestError(0), err.Cause)
	assert.NotEqual(t, NoStringerError(2), err.Cause)
	assert.Equal(t, "Ok", err.Error())
}

type NoStringerError uint

const (
	NoStringerErrorMyBad = NoStringerError(iota + 1)
	NoStringerErrorInternalFailure
)

func TestNoStringerError(t *testing.T) {
	err := errors.New(NoStringerErrorInternalFailure)
	assert.Implements(t, (*error)(nil), err)
	assert.True(t, err.IsErr())
	assert.False(t, err.IsOk())
	assert.Equal(t, NoStringerError(2), err.Cause)
	assert.NotEqual(t, TestError(2), err.Cause)
	assert.Equal(t, "errors_test.NoStringerError(2)", err.Error())
}

func TestNoStringerOk(t *testing.T) {
	err := errors.Ok[NoStringerError]()
	assert.Implements(t, (*error)(nil), err)
	assert.False(t, err.IsErr())
	assert.True(t, err.IsOk())
	assert.Equal(t, NoStringerError(0), err.Cause)
	assert.NotEqual(t, TestError(0), err.Cause)
	assert.Equal(t, "errors_test.NoStringerError(Ok)", err.Error())
}
