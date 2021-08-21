package errors

import (
	"fmt"
	"testing"
)

var (
	errFooA  = New("fooA", "whoops: %s", "this is bad")
	errFooA2 = fmt.Errorf("whoops: %s", "this is bad")
)

func fooA() error {
	return errFooA
}

func fooB() error {
	return Wrap("fooB", fooA())
}

func fooA2() error {
	return errFooA2
}

func fooB2() error {
	return Wrap("fooB", fooA2())
}

func Test_New(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		about        string
		internalCode string
		format       string
		values       []interface{}
	}{
		{
			about:        "it creates a new internal error with format",
			internalCode: "ER1001",
			format:       "whoops",
		},
		{
			about:        "it creates a new internal error with format and values",
			internalCode: "ER1001",
			format:       "whoops: %s",
			values:       []interface{}{"bad"},
		},
		{
			about:        "it creates a new internal error with empty internal code",
			internalCode: "",
			format:       "whoops",
			values:       []interface{}{"bad"},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.about, func(t *testing.T) {
			err := New(testCase.internalCode, testCase.format, testCase.values...)

			if err == nil {
				t.Errorf("created internal error is nil")
				return
			}

			expectedErrorString := fmt.Sprintf(testCase.format, testCase.values...)

			if err.(*cause).internalCode != testCase.internalCode {
				t.Errorf("expected internal code '%v', got '%v'", expectedErrorString, err)
				return
			}

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
		about                string
		internalCode         string
		err                  error
		expectedInternalCode string
		expectedErrorString  string
	}{
		{
			about:                "it creates a new wrapped error with the given cause",
			internalCode:         "ER1001",
			err:                  New("ER1000", "whoops"),
			expectedInternalCode: "ER1000",
			expectedErrorString:  "whoops",
		},
		{
			about:                "it creates a new wrapped error with the given error",
			internalCode:         "ER1001",
			err:                  fmt.Errorf("whoops"),
			expectedInternalCode: "ER1001",
			expectedErrorString:  "whoops",
		},
		{
			about:                "it creates a new wrapped error when internal code is empty",
			internalCode:         "",
			err:                  fmt.Errorf("whoops"),
			expectedInternalCode: "",
			expectedErrorString:  "whoops",
		},
		{
			about:                "it creates a new wrapped error when error is nil",
			internalCode:         "ER1001",
			err:                  nil,
			expectedInternalCode: "ER1001",
			expectedErrorString:  "%!s(<nil>)",
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.about, func(t *testing.T) {
			err := Wrap(testCase.internalCode, testCase.err)

			if err == nil {
				t.Errorf("created internal error is nil")
				return
			}

			if InternalCode(err) != testCase.expectedInternalCode {
				t.Errorf("expected internal code '%v', got '%v'", testCase.expectedInternalCode, InternalCode(err))
				return
			}

			if err.Error() != testCase.expectedErrorString {
				t.Errorf("expected error '%v', got '%v'", testCase.expectedErrorString, err.Error())
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
			about:         "it shows internalError as cause",
			errorFunc:     fooB,
			expectedCause: errFooA,
		},
		{
			about:         "it shows golang error as cause",
			errorFunc:     fooB2,
			expectedCause: errFooA2,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.about, func(t *testing.T) {
			err := testCase.errorFunc()
			actual := Cause(err)
			if actual != testCase.expectedCause {
				t.Errorf("expected %s (%p), but got %s (%p)", testCase.expectedCause, testCase.expectedCause, actual, actual)
			}
		})
	}
}

func Test_InternalCode(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		about                string
		errorFunc            func() error
		expectedInternalCode string
	}{
		{
			about:                "it shows internal code where cause is internalError",
			errorFunc:            fooB,
			expectedInternalCode: "fooA",
		},
		{
			about:                "it shows internal code where cause is golang error",
			errorFunc:            fooB2,
			expectedInternalCode: "fooB",
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.about, func(t *testing.T) {
			err := testCase.errorFunc()
			actual := InternalCode(err)
			if actual != testCase.expectedInternalCode {
				t.Errorf("expected %s, but got %s", testCase.expectedInternalCode, actual)
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
			about:         "it unwraps internalError",
			errorFunc:     fooB,
			expectedError: errFooA,
		},
		{
			about:         "it unwraps golang error",
			errorFunc:     fooB2,
			expectedError: errFooA2,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.about, func(t *testing.T) {
			err := testCase.errorFunc()
			actual := Unwrap(err)
			if actual != testCase.expectedError {
				t.Errorf("expected '%s', but got '%s'", testCase.expectedError, actual)
			}

			if testCase.expectedError == errFooA {
				if err.(*wrapped).Unwrap() != Unwrap(err) {
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
			about:         "it is internalError",
			errorFunc:     fooB,
			expectedError: errFooA,
		},
		{
			about:         "it is golang error",
			errorFunc:     fooB2,
			expectedError: errFooA2,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.about, func(t *testing.T) {
			err := testCase.errorFunc()
			if !Is(err, testCase.expectedError) {
				t.Errorf("expected '%s', but got '%s'", testCase.expectedError, err)
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
			about:         "it is an internalError",
			errorFunc:     fooB,
			asErr:         &wrapped{},
			expectedError: errFooA,
		},
		{
			about:         "it is a golang error",
			errorFunc:     fooB2,
			asErr:         fmt.Errorf(""),
			expectedError: errFooA2,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.about, func(t *testing.T) {
			err := testCase.errorFunc()

			if !As(err, &testCase.asErr) {
				t.Errorf("not the error expected")
			}
			if !Is(testCase.asErr, testCase.expectedError) {
				t.Errorf("expected '%s', but got '%s'", testCase.expectedError, testCase.asErr)
			}
		})
	}
}
