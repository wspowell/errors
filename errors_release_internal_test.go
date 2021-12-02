//go:build release
// +build release

package errors

import (
	"fmt"
	"testing"
)

func fooC() error {
	return New("fooC", "whoops: %s", "this is bad")
}

func fooD() error {
	return Propagate("fooD", fooC())
}

func fooE() error {
	return Convert("fooE", fooD(), ErrDiscrete)
}

func Test_internalError_Format(t *testing.T) {
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
			about:               "all errors internalErrr - it prints error as value with internal code stack and with no stack trace",
			errorFunc:           fooB,
			formatString:        "%+v",
			expectedErrorString: "[fooB][fooA] whoops: this is bad",
		},
		{
			about:               "all errors internalErrr - it prints converted error as string",
			errorFunc:           fooE,
			formatString:        "%s",
			expectedErrorString: "concrete error",
		},
		{
			about:               "all errors internalErrr - it prints converted error as value",
			errorFunc:           fooE,
			formatString:        "%v",
			expectedErrorString: "[fooE] concrete error",
		},
		{
			about:               "all errors internalErrr - it prints converted error as value with internal code stack",
			errorFunc:           fooE,
			formatString:        "%#v",
			expectedErrorString: "[fooE] concrete error -> [fooD][fooC] whoops: this is bad",
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
			about:               "cause error is golang error - it prints error as value with internal code stack and with no stack trace",
			errorFunc:           fooB2,
			formatString:        "%+v",
			expectedErrorString: "[fooB] whoops: this is bad",
		},
	}
	for index := range testCases {
		testCase := testCases[index]
		t.Run(testCase.about, func(t *testing.T) {
			t.Parallel()

			err := testCase.errorFunc()

			actual := fmt.Sprintf(testCase.formatString, err)

			if actual != testCase.expectedErrorString {
				t.Errorf("expected %s, but got %s", testCase.expectedErrorString, actual)
			}
		})
	}
}
