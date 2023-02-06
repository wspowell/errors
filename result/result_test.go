package result_test

import (
	goerrors "errors"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wspowell/errors"
	"github.com/wspowell/errors/result"
)

type FailureError uint64

const (
	errFailureNone FailureError = iota
	errFailureError
)

func (self FailureError) String() string {
	return [...]string{
		"none",
		"failure",
	}[self]
}

func TestOk(t *testing.T) {
	t.Parallel()

	expectedValue := 1
	res := result.Ok[int, error](expectedValue)

	actualValue, actualErr := res.Result()
	assert.Equal(t, expectedValue, actualValue)
	assert.Equal(t, nil, actualErr)

	assert.True(t, res.IsOk())
	assert.Equal(t, 1, res.Value())
	assert.Equal(t, 1, res.ValueOr(0))
	assert.Equal(t, 1, res.ValueOrPanic())
}

func TestErr(t *testing.T) {
	t.Parallel()

	expectedErr := goerrors.New("test")
	res := result.Err[int](expectedErr)

	val, err := res.Result()
	assert.Equal(t, 0, val)
	assert.Equal(t, expectedErr, err)

	assert.False(t, res.IsOk())
	assert.Equal(t, 0, res.Value())
	assert.Equal(t, 0, res.ValueOr(0))
	assert.Panics(t, func() { res.ValueOrPanic() })
	assert.Equal(t, expectedErr, res.Error())
}

func TestResultIntPointer(t *testing.T) {
	t.Parallel()

	value := 1
	res := result.Ok[*int, error](&value)

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
	res := result.Ok[S, error](s)

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
	res := result.Ok[*S, error](s)

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

	resultChannel := make(chan result.Result[int, error], len(testCases))

	for _, testCase := range testCases {
		go func(testCase int) {
			defer waitGroup.Done()
			r := result.Ok[int, error](testCase * 2)
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

	expectedErr := goerrors.New("test")
	resultChannel := make(chan result.Result[int, error], len(testCases))

	for _, testCase := range testCases {
		go func(_ int) {
			defer waitGroup.Done()
			r := result.Err[int](expectedErr)
			resultChannel <- r
		}(testCase)
	}

	waitGroup.Wait()
	close(resultChannel)

	for res := range resultChannel {
		assert.Equal(t, expectedErr, res.Error())
	}
}

type TestError uint

const (
	TestErrorOne = TestError(iota + 1)
	TestErrorTwo
)

func TestErrorsOk(t *testing.T) {
	t.Parallel()

	expectedValue := 1
	res := result.Ok[int, errors.Error[TestError]](expectedValue)

	actualValue, actualErr := res.Result()
	assert.Equal(t, expectedValue, actualValue)
	assert.Equal(t, errors.Ok[TestError](), actualErr)

	assert.True(t, res.IsOk())
	assert.Equal(t, 1, res.Value())
	assert.Equal(t, 1, res.ValueOr(0))
	assert.Equal(t, 1, res.ValueOrPanic())
}

func TestErrorsErr(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New(TestErrorOne)
	res := result.Err[int](expectedErr)

	val, err := res.Result()
	assert.Equal(t, 0, val)
	assert.Equal(t, expectedErr, err)

	assert.False(t, res.IsOk())
	assert.Equal(t, 0, res.Value())
	assert.Equal(t, 0, res.ValueOr(0))
	assert.Panics(t, func() { res.ValueOrPanic() })
	assert.Equal(t, expectedErr, res.Error())
}
