//go:build !release
// +build !release

package errors

import (
	"fmt"
	"strings"
	"testing"
)

func fooC() error {
	return New("fooC", "whoops: %s", "this is bad")
}

func fooD() error {
	return Wrap("fooD", fooC())
}

func Test_error_Format(t *testing.T) {
	t.Parallel()

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
			expectedErrorString: "[fooB][fooA] whoops: this is bad",
		},
		{
			about:               "all errors internalErrr - it prints error as value with internal code stack and with stack trace",
			errorFunc:           fooB,
			formatString:        "%+v",
			expectedErrorString: "[fooB][fooA] whoops: this is bad\ngithub.com/wspowell/errors.init",
		},
		{
			about:               "all errors internalErrr - it prints error as value with internal code stack and with stack trace nested error",
			errorFunc:           fooD,
			formatString:        "%+v",
			expectedErrorString: "[fooD][fooC] whoops: this is bad\ngithub.com/wspowell/errors.fooC",
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
			expectedErrorString: "[fooB] whoops: this is bad\ngithub.com/wspowell/errors.fooB2",
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
