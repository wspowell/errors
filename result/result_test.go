package result_test

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wspowell/errors"
	"github.com/wspowell/errors/result"
)

const (
	errErrorFailure = "failure"
)

func TestOk(t *testing.T) {
	t.Parallel()

	value := 1
	res := result.Ok[int, errors.Standard](value)

	val, err := res.Result()
	assert.Equal(t, value, val)
	assert.Equal(t, errors.None[string](), err)

	assert.True(t, res.IsOk())
	assert.Equal(t, 1, res.Value())
	assert.Equal(t, 1, res.ValueOr(0))
	assert.Equal(t, 1, res.ValueOrPanic())
}

func TestErr(t *testing.T) {
	t.Parallel()

	res := result.Err[int](errors.New(errErrorFailure))

	val, err := res.Result()
	assert.Equal(t, 0, val)
	assert.Equal(t, errErrorFailure, err.Into())

	assert.False(t, res.IsOk())
	assert.Equal(t, 0, res.Value())
	assert.Equal(t, 0, res.ValueOr(0))
	assert.Panics(t, func() { res.ValueOrPanic() })
	assert.Equal(t, errErrorFailure, res.Error().Into())
}

func TestResultIntPointer(t *testing.T) {
	t.Parallel()

	value := 1
	res := result.Ok[*int, errors.Standard](&value)

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
	res := result.Ok[S, errors.Standard](s)

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
	res := result.Ok[*S, errors.Standard](s)

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

	resultChannel := make(chan result.Result[int, errors.Standard], len(testCases))

	for _, testCase := range testCases {
		go func(testCase int) {
			defer waitGroup.Done()
			r := result.Ok[int, errors.Standard](testCase * 2)
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

	resultChannel := make(chan result.Result[int, errors.Standard], len(testCases))

	for _, testCase := range testCases {
		go func(_ int) {
			defer waitGroup.Done()
			r := result.Err[int](errors.New(errErrorFailure))
			resultChannel <- r
		}(testCase)
	}

	waitGroup.Wait()
	close(resultChannel)

	for res := range resultChannel {
		assert.Equal(t, errErrorFailure, res.Error().Into())
	}
}
