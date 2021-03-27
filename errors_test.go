package errors

import (
	"fmt"
	"strings"
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
		err          string
	}{
		{
			about:        "it creates a new internal error with the given error",
			internalCode: "ER1001",
			err:          "whoops",
		},
		{
			about:        "it creates a new internal error when internal code is empty",
			internalCode: "",
			err:          "whoops",
		},
		{
			about:        "it creates a new internal error when err is nil",
			internalCode: "ER1001",
			err:          "",
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.about, func(t *testing.T) {
			err := New(testCase.internalCode, testCase.err)

			if err == nil {
				t.Errorf("created internal error is nil")
				return
			}

			var expectedFeatureCode string

			if testCase.internalCode == "" {
				expectedFeatureCode = ""
			} else {
				expectedFeatureCode = testCase.internalCode
			}

			expectedErrorString := testCase.err

			if err.(*internalError).internalCode != expectedFeatureCode {
				t.Errorf("expected internal code '%v', got '%v'", testCase.err, err)
				return
			}

			if err.Error() != expectedErrorString {
				t.Errorf("expected error '%v', got '%v'", testCase.err, err)
				return
			}
		})
	}
}

func Test_Wrap(t *testing.T) {
	testCases := []struct {
		about        string
		internalCode string
		err          error
	}{
		{
			about:        "it creates a new internal error with the given error",
			internalCode: "ER1001",
			err:          fmt.Errorf("whoops"),
		},
		{
			about:        "it creates a new internal error when internal code is empty",
			internalCode: "",
			err:          fmt.Errorf("whoops"),
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.about, func(t *testing.T) {
			err := Wrap(testCase.internalCode, testCase.err)

			if err == nil {
				t.Errorf("created internal error is nil")
				return
			}

			var expectedFeatureCode string
			var expectedErrorString string

			if testCase.internalCode == "" {
				expectedFeatureCode = ""
			} else {
				expectedFeatureCode = testCase.internalCode
			}

			if testCase.err == nil {
				expectedErrorString = "ERROR"
			} else {
				expectedErrorString = testCase.err.Error()
			}

			if err.(*internalError).internalCode != expectedFeatureCode {
				t.Errorf("expected internal code '%v', got '%v'", testCase.err.Error(), err.Error())
				return
			}

			if err.Error() != expectedErrorString {
				t.Errorf("expected error '%v', got '%v'", expectedErrorString, err.Error())
				return
			}
		})
	}
}

func Test_internalError_Format(t *testing.T) {

	testCases := []struct {
		about               string
		errorFunc           func() error
		formatString        string
		expectedErrorString string
	}{
		// All errors are internalError.
		{
			about:               "all errors internalErrr - it prints error as string",
			errorFunc:           fooB,
			formatString:        "%s",
			expectedErrorString: "whoops: this is bad",
		},
		{
			about:               "all errors internalErrr - it prints error as value",
			errorFunc:           fooB,
			formatString:        "%v",
			expectedErrorString: "[fooA] whoops: this is bad",
		},
		{
			about:               "all errors internalErrr - it prints error as value with internal code stack",
			errorFunc:           fooB,
			formatString:        "%#v",
			expectedErrorString: "[fooA,fooB] whoops: this is bad",
		},
		{
			about:               "all errors internalErrr - it prints error as value with internal code stack and with stack trace",
			errorFunc:           fooB,
			formatString:        "%+v",
			expectedErrorString: "[fooA,fooB] whoops: this is bad\ngithub.com/wspowell/errors.New",
		},
		// Root error is golang error.
		{
			about:               "cause error is golang error - it prints error as string",
			errorFunc:           fooB2,
			formatString:        "%s",
			expectedErrorString: "whoops: this is bad",
		},
		{
			about:               "cause error is golang error - it prints error as value",
			errorFunc:           fooB2,
			formatString:        "%v",
			expectedErrorString: "[fooB] whoops: this is bad",
		},
		{
			about:               "cause error is golang error - it prints error as value with internal code stack",
			errorFunc:           fooB2,
			formatString:        "%#v",
			expectedErrorString: "[fooB] whoops: this is bad",
		},
		{
			about:               "cause error is golang error - it prints error as value with internal code stack and with stack trace",
			errorFunc:           fooB2,
			formatString:        "%+v",
			expectedErrorString: "[fooB] whoops: this is bad\ngithub.com/wspowell/errors.Wrap",
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.about, func(t *testing.T) {
			err := testCase.errorFunc()

			actual := fmt.Sprintf(testCase.formatString, err)

			if testCase.formatString == "%+v" {
				if !strings.HasPrefix(actual, testCase.expectedErrorString) {
					t.Errorf("expected %s, but got %s", testCase.expectedErrorString, actual)
				}
			} else {
				if actual != testCase.expectedErrorString {
					t.Errorf("expected %s, but got %s", testCase.expectedErrorString, actual)
				}
			}
		})
	}
}

func Test_Cause(t *testing.T) {

	testCases := []struct {
		about                string
		errorFunc            func() error
		expectedCause        error
		expectedInternalCode string
	}{
		{
			about:                "it shows internalError as cause",
			errorFunc:            fooB,
			expectedCause:        errFooA,
			expectedInternalCode: "fooA",
		},
		{
			about:                "it shows golang error as cause",
			errorFunc:            fooB2,
			expectedCause:        errFooA2,
			expectedInternalCode: "fooB",
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.about, func(t *testing.T) {
			err := testCase.errorFunc()
			actual := Cause(err)
			if actual != testCase.expectedCause {
				t.Errorf("expected %s, but got %s", testCase.expectedCause, actual)
			}
		})
	}
}

func Test_InternalCode(t *testing.T) {
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
				if err.(*internalError).Unwrap() != Unwrap(err) {
					t.Errorf("unwraps are not equal")
				}
			}
		})
	}
}

func Test_Is(t *testing.T) {

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

	testCases := []struct {
		about         string
		errorFunc     func() error
		asErr         error
		expectedError error
	}{
		{
			about:         "it is an internalError",
			errorFunc:     fooB,
			asErr:         &internalError{},
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
