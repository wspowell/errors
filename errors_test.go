package errors_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wspowell/errors"
)

const (
	expectedError           = "WHOOPS"
	expectedErrorStackTrace = "WHOOPS\ngithub.com/wspowell/errors_test.errFunc"
)

func errFunc(ctx context.Context) errors.Error {
	return errors.New(ctx, "%s", "WHOOPS")
}

func TestNew(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	err := errFunc(ctx)

	assert.Contains(t, err.String(), expectedError)
}

func TestNewWithStackTrace(t *testing.T) {
	t.Parallel()

	ctx := errors.WithStackTrace(context.Background())
	err := errFunc(ctx)

	assert.Contains(t, err.String(), expectedErrorStackTrace)
}
