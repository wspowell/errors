// +build release

package errors

import (
	"fmt"
	"testing"
)

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
			expectedErrorString: "[fooB][fooA] whoops: this is bad",
		},
		{
			about:               "all errors internalErrr - it prints error as value with internal code stack and with no stack trace",
			errorFunc:           fooB,
			formatString:        "%+v",
			expectedErrorString: "[fooB][fooA] whoops: this is bad",
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
	for _, testCase := range testCases {
		t.Run(testCase.about, func(t *testing.T) {
			err := testCase.errorFunc()

			actual := fmt.Sprintf(testCase.formatString, err)

			if actual != testCase.expectedErrorString {
				t.Errorf("expected %s, but got %s", testCase.expectedErrorString, actual)
			}
		})
	}
}
