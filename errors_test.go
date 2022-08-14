package errors_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wspowell/errors"
)

const (
	errExpected = "WHOOPS"
)

func errFunc(ctx context.Context) errors.Error[string] {
	return errors.New(ctx, errExpected)
}

func TestNone(t *testing.T) {
	t.Parallel()

	assert.True(t, errors.None[string]().IsNone())
	assert.True(t, errors.New(context.Background(), "").IsNone())
	assert.False(t, errors.New(context.Background(), errExpected).IsNone())
}

func TestNew(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	err := errFunc(ctx)

	assert.Equal(t, err.Error(), "WHOOPS")
	assert.Equal(t, errExpected, err.Into())
}

func TestNewWithStackTrace(t *testing.T) {
	t.Parallel()

	ctx := errors.WithStackTrace(context.Background())
	err := errFunc(ctx)

	stackTraceFragment := "WHOOPS\ngithub.com/wspowell/errors_test.errFunc"
	assert.Equal(t, err.Error(), "WHOOPS")
	assert.Contains(t, fmt.Sprintf("%+v", err), stackTraceFragment)
	assert.Equal(t, errExpected, err.Into())
}

type myError string

const (
	errOne   = myError("one")
	errTwo   = myError("two")
	errThree = myError("three")
)

func multipleErrors(ctx context.Context, err int) errors.Error[myError] {
	switch err {
	case 1:
		return errors.New(ctx, errOne)
	case 2:
		return errors.New(ctx, errTwo)
	case 3:
		return errors.New(ctx, errThree)
	default:
		return errors.None[myError]()
	}
}

func TestSentinelEnum(t *testing.T) {
	t.Parallel()

	err := multipleErrors(context.Background(), 1)
	//nolint:exhaustive // reason: expected lint error for testing
	switch err.Into() {
	case errOne:
		// Expected case.
	default:
		assert.Fail(t, "expected errOne")
	}
}
