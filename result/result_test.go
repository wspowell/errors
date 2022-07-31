package result_test

import (
	"context"
	// nolint:depguard // reason: importing "errors" for test comparison
	goerrors "errors"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wspowell/errors"
	"github.com/wspowell/errors/result"
)

var (
	errErrorFailure = errors.New(context.Background(), "failure")
	errGoFailure    = goerrors.New("failure")
)

func TestOk(t *testing.T) {
	t.Parallel()

	value := 1
	res := result.Ok(value)

	val, err := res.Result()
	assert.Equal(t, value, val)
	assert.Equal(t, errors.Error{}, err)

	assert.True(t, res.IsOk())
	assert.Equal(t, 1, res.Value())
	assert.Equal(t, 1, res.ValueOr(0))
	assert.Equal(t, 1, res.ValueOrPanic())
}

func TestErr(t *testing.T) {
	t.Parallel()

	res := result.Err[int](errErrorFailure)

	val, err := res.Result()
	assert.Equal(t, 0, val)
	assert.Equal(t, errErrorFailure, err)

	assert.False(t, res.IsOk())
	assert.Equal(t, 0, res.Value())
	assert.Equal(t, 0, res.ValueOr(0))
	assert.Panics(t, func() { res.ValueOrPanic() })
	assert.Equal(t, "failure", res.Error().String())
	assert.Equal(t, res.Error(), errErrorFailure)
}

func TestErrPanicsWithError(t *testing.T) {
	t.Parallel()

	res := result.Err[struct{}](errErrorFailure)

	defer func() {
		if err := errors.Recover(context.Background(), recover()); err.IsNone() {
			t.Error("Result should panic")
		} else if res.Error() != errErrorFailure {
			t.Error("Result should panic with the error")
		}
	}()

	res.ValueOrPanic()
}

func TestResultIntPointer(t *testing.T) {
	t.Parallel()

	value := 1
	res := result.Ok(&value)

	expected := 1
	assert.True(t, res.IsOk())
	assert.Equal(t, &expected, res.Value())
	assert.Equal(t, &expected, res.ValueOr(nil))
	assert.Equal(t, &expected, res.ValueOrPanic())
}

func TestResultStruct(t *testing.T) {
	t.Parallel()

	type S struct {
		A int
		B string
	}
	s := S{A: 1, B: "b"}
	res := result.Ok(s)

	assert.True(t, res.IsOk())
	assert.Equal(t, S{1, "b"}, res.Value())
	assert.Equal(t, S{1, "b"}, res.ValueOr(S{}))
	assert.Equal(t, S{1, "b"}, res.ValueOrPanic())
}

func TestResultStructPointer(t *testing.T) {
	t.Parallel()

	type S struct {
		A int
		B string
	}
	s := &S{A: 1, B: "b"}
	res := result.Ok(s)

	assert.True(t, res.IsOk())
	assert.Equal(t, &S{1, "b"}, res.Value())
	assert.Equal(t, &S{1, "b"}, res.ValueOr(&S{}))
	assert.Equal(t, &S{1, "b"}, res.ValueOrPanic())
}

func TestResultOkPassingThroughChannel(t *testing.T) {
	t.Parallel()

	testCases := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	var waitGroup sync.WaitGroup
	waitGroup.Add(len(testCases))

	resultChannel := make(chan result.Result[int], len(testCases))

	for _, testCase := range testCases {
		go func(testCase int) {
			defer waitGroup.Done()
			r := result.Ok(testCase * 2)
			resultChannel <- r
		}(testCase)
	}

	waitGroup.Wait()
	close(resultChannel)

	for res := range resultChannel {
		assert.NotEqual(t, 0, res.Value())
	}
}

func TestResultErrPassingThroughChannel(t *testing.T) {
	t.Parallel()

	testCases := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	var waitGroup sync.WaitGroup
	waitGroup.Add(len(testCases))

	resultChannel := make(chan result.Result[int], len(testCases))

	for _, testCase := range testCases {
		go func(_ int) {
			defer waitGroup.Done()
			r := result.Err[int](errErrorFailure)
			resultChannel <- r
		}(testCase)
	}

	waitGroup.Wait()
	close(resultChannel)

	for res := range resultChannel {
		assert.Equal(t, res.Error(), errErrorFailure)
	}
}
