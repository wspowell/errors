package result_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wspowell/errors"
	"github.com/wspowell/errors/result"
)

func TestThenOk(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	value := 1
	res := result.Ok[int, Error](value)

	res2 := result.Then(ctx, res, func(ctx context.Context, v int) result.Result[int, Error] {
		return result.Ok[int, Error](v + 1)
	})

	assert.True(t, res2.IsOk())
	assert.Equal(t, 2, res2.ValueOr(0))
}

func TestThenErr(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	res := result.Err[int](errors.New(context.Background(), errErrorFailure))

	res2 := result.Then(ctx, res, func(ctx context.Context, v int) result.Result[int, Error] {
		return result.Ok[int, Error](v + 1)
	})

	assert.False(t, res2.IsOk())
	assert.Equal(t, 0, res2.ValueOr(0))
}

func TestWhenThenOk(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	value := 1
	res := result.Ok[int, Error](value)

	res2 := result.When[int, float64, Error](res).Then(ctx, func(ctx context.Context, v int) result.Result[float64, Error] {
		return result.Ok[float64, Error](float64(v + 1.0))
	})

	assert.True(t, res2.IsOk())
	assert.Equal(t, float64(2), res2.ValueOr(0))
}

func TestWhenThenErr(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	res := result.Err[int](errors.New(context.Background(), errErrorFailure))

	res2 := result.When[int, float64, Error](res).Then(ctx, func(ctx context.Context, v int) result.Result[float64, Error] {
		return result.Ok[float64, Error](float64(v + 1.0))
	})

	assert.False(t, res2.IsOk())
	assert.Equal(t, float64(0), res2.ValueOr(0))
	assert.Equal(t, errErrorFailure, res.Error().Into())
}
