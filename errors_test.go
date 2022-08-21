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
	assert.True(t, errors.Some("").IsNone())
	assert.False(t, errors.Some(ErrSome).IsNone())
	assert.False(t, errors.Some(ErrEnumErr1).IsNone())
	assert.False(t, errors.Some("whoops").IsNone())
}

func TestIsSome(t *testing.T) {
	t.Parallel()

	assert.False(t, errors.None[string]().IsSome())
	assert.False(t, errors.None[EnumErr]().IsSome())
	assert.False(t, errors.Some("").IsSome())
	assert.True(t, errors.Some(ErrSome).IsSome())
	assert.True(t, errors.Some(ErrEnumErr1).IsSome())
	assert.True(t, errors.Some("whoops").IsSome())
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

	assert.Equal(t, "error: WHOOPS", errors.Some("error: %s", "WHOOPS").Into())
}
