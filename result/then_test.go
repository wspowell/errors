package result_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wspowell/errors/result"
)

func TestThenOk(t *testing.T) {
	t.Parallel()

	value := 1
	res := result.Ok[int, FailureError](value)

	res2 := result.Then(res, func(v int) result.Result[int, FailureError] {
		return result.Ok[int, FailureError](v + 1)
	})

	assert.True(t, res2.IsOk())
	assert.Equal(t, 2, res2.ValueOr(0))
}

func TestThenErr(t *testing.T) {
	t.Parallel()

	res := result.Err[int](errFailureError)

	res2 := result.Then(res, func(v int) result.Result[int, FailureError] {
		return result.Ok[int, FailureError](v + 1)
	})

	assert.False(t, res2.IsOk())
	assert.Equal(t, 0, res2.ValueOr(0))
}

func TestWhenThenOk(t *testing.T) {
	t.Parallel()

	value := 1
	res := result.Ok[int, FailureError](value)

	res2 := result.When[int, float64, FailureError](res).Then(func(v int) result.Result[float64, FailureError] {
		return result.Ok[float64, FailureError](float64(v + 1.0))
	})

	assert.True(t, res2.IsOk())
	assert.Equal(t, float64(2), res2.ValueOr(0))
}

func TestWhenThenErr(t *testing.T) {
	t.Parallel()

	res := result.Err[int](errFailureError)

	res2 := result.When[int, float64, FailureError](res).Then(func(v int) result.Result[float64, FailureError] {
		return result.Ok[float64, FailureError](float64(v + 1.0))
	})

	assert.False(t, res2.IsOk())
	assert.Equal(t, float64(0), res2.ValueOr(0))
	assert.Equal(t, errFailureError, res.Error())
}
