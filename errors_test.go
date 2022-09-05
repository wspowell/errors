package errors_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wspowell/errors"
)

func TestNewMessageError(t *testing.T) {
	t.Parallel()

	err := errors.NewMessage(ErrDetailedWhoops, "whoops: %s", "testing")

	assert.Equal(t, ErrDetailedWhoops, err.Error)
	assert.Equal(t, "whoops: testing", err.Message)
	assert.Equal(t, "errors_test.DetailedErr(1, whoops: testing)", err.String())
}

func TestOk(t *testing.T) {
	t.Parallel()

	err := errors.Ok[DetailedErr]()

	assert.Equal(t, ErrDetailedNone, err.Error)
	assert.Equal(t, "", err.Message)
	assert.Equal(t, "errors_test.DetailedErr(Ok)", err.String())
	assert.Panics(t, func() { errors.NewMessage[DetailedErr](0, "") })
}
