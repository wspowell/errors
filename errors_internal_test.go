package errors

import (
	"fmt"
	"testing"
)

var (
	errCause  = New("cause: %s", "error")
	errGolang = fmt.Errorf("golang: %s", "error")

	errWrapped   = New("wrapped")
	errRewrapped = New("rewrapped")
)

func returnCause() error {
	return errCause
}

func returnCauseWrapped() error {
	return Wrap(errCause, errWrapped)
}

func returnCauseWrappedTwice() error {
	return Wrap(returnCauseWrapped(), errRewrapped)
}

func returnGolang() error {
	// nolint:wrapcheck // reason: not wrapped for testing
	return errGolang
}

func returnGolangWrapped() error {
	return Wrap(errGolang, errWrapped)
}

func returnGolangWrappedTwice() error {
	return Wrap(returnGolangWrapped(), errRewrapped)
}

func Test_New(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		about  string
		format string
		values []any
	}{
		{
			about:  "it creates a new internal error with format",
			format: "whoops",
		},
		{
			about:  "it creates a new internal error with format and values",
			format: "whoops: %s",
			values: []any{"bad"},
		},
		{
			about:  "it creates a new internal error with empty internal code",
			format: "whoops",
			values: []any{"bad"},
		},
	}
	for index := range testCases {
		testCase := testCases[index]
		t.Run(testCase.about, func(t *testing.T) {
			t.Parallel()

			err := New(testCase.format, testCase.values...)

			if err == nil {
				t.Errorf("created internal error is nil")

				return
			}

			expectedErrorString := fmt.Sprintf(testCase.format, testCase.values...)
			if err.Error() != expectedErrorString {
				t.Errorf("expected error '%v', got '%v'", expectedErrorString, err)

				return
			}
		})
	}
}

func Test_Wrap(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		about               string
		fromErr             error
		toErr               error
		expectedErrorString string
	}{
		{
			about:               "it creates a new converted error with the given cause",
			fromErr:             New("whoops"),
			toErr:               errCause,
			expectedErrorString: "cause: error",
		},
		{
			about:               "it creates a new converted error with the given error",
			fromErr:             fmt.Errorf("whoops"), // nolint:goerr113 // reason: error created for test
			toErr:               errCause,
			expectedErrorString: "cause: error",
		},
		{
			about:               "it creates a new converted error when error is nil",
			fromErr:             nil,
			toErr:               errCause,
			expectedErrorString: "cause: error",
		},
		{
			about:               "it passes back original error when discrete error is nil",
			fromErr:             New("whoops"),
			toErr:               nil,
			expectedErrorString: "whoops",
		},
		// Golang errors.
		{
			about:               "it creates a new golang wrapped error with the given cause",
			fromErr:             New("whoops"),
			toErr:               errGolang,
			expectedErrorString: "golang: error",
		},
		{
			about:               "it creates a new golang wrapped error with the given error",
			fromErr:             fmt.Errorf("whoops"), // nolint:goerr113 // reason: error created for test
			toErr:               errGolang,
			expectedErrorString: "golang: error",
		},
		{
			about:               "it creates a new golang wrapped error when internal code is empty",
			fromErr:             fmt.Errorf("whoops"), // nolint:goerr113 // reason: error created for test
			toErr:               errGolang,
			expectedErrorString: "golang: error",
		},
		{
			about:               "it creates a new golang wrapped error when error is nil",
			fromErr:             nil,
			toErr:               errGolang,
			expectedErrorString: "golang: error",
		},
		{
			about:               "it passes back original golang error when wrapped error is nil",
			fromErr:             fmt.Errorf("whoops"), // nolint:goerr113 // reason: error created for test
			toErr:               nil,
			expectedErrorString: "whoops",
		},
	}
	for index := range testCases {
		testCase := testCases[index]
		t.Run(testCase.about, func(t *testing.T) {
			t.Parallel()

			err := Wrap(testCase.fromErr, testCase.toErr)

			if err == nil {
				t.Errorf("created internal error is nil")

				return
			}

			if err.Error() != testCase.expectedErrorString {
				t.Errorf("expected error '%v', got '%v'", testCase.expectedErrorString, err.Error())

				return
			}

			if testCase.toErr != nil && !Is(err, testCase.toErr) {
				t.Errorf("expected error to be converted error, but is not")

				return
			}

			if testCase.toErr != nil && Is(err, testCase.fromErr) {
				t.Errorf("expected error to no longer be original error, but is")

				return
			}
		})
	}
}

func Test_Cause(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		about         string
		errorFunc     func() error
		expectedCause error
	}{
		{
			about:         "cause",
			errorFunc:     returnCause,
			expectedCause: errCause,
		},
		{
			about:         "wrapped cause",
			errorFunc:     returnCauseWrapped,
			expectedCause: errWrapped,
		},
		{
			about:         "golang error",
			errorFunc:     returnGolang,
			expectedCause: errGolang,
		},
		{
			about:         "wrapped golang error",
			errorFunc:     returnGolangWrapped,
			expectedCause: errWrapped,
		},
	}
	for index := range testCases {
		testCase := testCases[index]
		t.Run(testCase.about, func(t *testing.T) {
			t.Parallel()

			err := testCase.errorFunc()
			actual := Cause(err)
			// nolint:errorlint,goerr113 // reason: test should check the exact value, not find one in the error chain
			if actual != testCase.expectedCause {
				t.Errorf("expected '%s' (%p), but got '%s' (%p)", testCase.expectedCause, testCase.expectedCause, actual, actual)
			}
		})
	}
}

func Test_Unwrap(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		about         string
		errorFunc     func() error
		expectedError error
	}{
		{
			about:         "it unwraps cause",
			errorFunc:     returnCause,
			expectedError: nil,
		},
		{
			about:         "it unwraps golang error",
			errorFunc:     returnGolang,
			expectedError: nil,
		},
		// Wrapped errors.
		{
			about:         "it unwraps wrapped cause",
			errorFunc:     returnCauseWrapped,
			expectedError: errWrapped,
		},
		{
			about:         "it unwraps wrapped golang error",
			errorFunc:     returnGolangWrapped,
			expectedError: errWrapped,
		},
	}
	for index := range testCases {
		testCase := testCases[index]
		t.Run(testCase.about, func(t *testing.T) {
			t.Parallel()

			err := testCase.errorFunc()
			actual := Unwrap(err)
			// nolint:errorlint,goerr113 // reason: test should check the exact value, not find one in the error chain
			if actual != testCase.expectedError {
				t.Errorf("expected '%s', but got '%s'", testCase.expectedError, actual)
			}

			// nolint:errorlint,goerr113 // reason: test should check the exact value, not find one in the error chain
			if errCause == testCase.expectedError {
				if err.(*cause).Unwrap() != Unwrap(err) {
					t.Errorf("unwraps are not equal")
				}
			}
		})
	}
}

func Test_Is(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		about         string
		errorFunc     func() error
		expectedError error
	}{
		{
			about:         "it is a cause error",
			errorFunc:     returnCause,
			expectedError: errCause,
		},
		{
			about:         "it is a golang error",
			errorFunc:     returnGolang,
			expectedError: errGolang,
		},
		{
			about:         "it is a wrapped cause error",
			errorFunc:     returnCauseWrapped,
			expectedError: errWrapped,
		},
		{
			about:         "it is a golang error",
			errorFunc:     returnGolangWrapped,
			expectedError: errWrapped,
		},
	}
	for index := range testCases {
		testCase := testCases[index]
		t.Run(testCase.about, func(t *testing.T) {
			t.Parallel()

			err := testCase.errorFunc()
			if !Is(err, testCase.expectedError) {
				t.Errorf("expected '%s' %p, but got '%s' %p", testCase.expectedError, testCase.expectedError, err, err)
			}
		})
	}
}

func Test_As(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		about         string
		errorFunc     func() error
		asErr         error
		expectedError error
	}{
		{
			about:         "it is a cause",
			errorFunc:     returnCause,
			asErr:         &cause{},
			expectedError: errCause,
		},
		{
			about:         "it is a golang error",
			errorFunc:     returnGolang,
			asErr:         fmt.Errorf(""), // nolint:goerr113 // reason: error created for test
			expectedError: errGolang,
		},
		{
			about:         "it is a wrapped cause",
			errorFunc:     returnCauseWrapped,
			asErr:         &cause{},
			expectedError: errWrapped,
		},
		{
			about:         "it is a wrapped golang error",
			errorFunc:     returnGolangWrapped,
			asErr:         fmt.Errorf(""), // nolint:goerr113 // reason: error created for test
			expectedError: errWrapped,
		},
	}
	for index := range testCases {
		testCase := testCases[index]
		t.Run(testCase.about, func(t *testing.T) {
			t.Parallel()

			err := testCase.errorFunc()

			if !As(err, &testCase.asErr) {
				t.Errorf("not the error expected")
			}
			if !Is(testCase.asErr, testCase.expectedError) {
				t.Errorf("expected '%s' %p, but got '%s' %p", testCase.expectedError, testCase.expectedError, testCase.asErr, testCase.asErr)
			}
		})
	}
}
