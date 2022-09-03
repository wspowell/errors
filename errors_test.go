package errors_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wspowell/errors"
)

func TestIsNone(t *testing.T) {
	t.Parallel()

	assert.True(t, errors.None[string]().IsNone())
	assert.True(t, errors.None[EnumErr]().IsNone())
	assert.True(t, errors.New("").IsNone())
	assert.False(t, errors.New(ErrSome).IsNone())
	assert.False(t, errors.New(ErrEnumErr1).IsNone())
	assert.False(t, errors.New("whoops").IsNone())
}

func TestIsSome(t *testing.T) {
	t.Parallel()

	assert.False(t, errors.None[string]().IsSome())
	assert.False(t, errors.None[EnumErr]().IsSome())
	assert.False(t, errors.New("").IsSome())
	assert.True(t, errors.New(ErrSome).IsSome())
	assert.True(t, errors.New(ErrEnumErr1).IsSome())
	assert.True(t, errors.New("whoops").IsSome())
}

func TestErr(t *testing.T) {
	t.Parallel()

	err := ErrFunc()

	assert.Equal(t, err.Into(), "WHOOPS")
	assert.Equal(t, ErrSome, err.Into())
	assert.Equal(t, err, ErrFunc())
}

func TestErrFormat(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "error: WHOOPS", errors.New("error: %s", "WHOOPS").Into())
}
