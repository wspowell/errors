//go:build release
// +build release

package errors

import (
	"fmt"
	"testing"
)

func Test_error_Format(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		about               string
		errorFunc           func() error
		formatString        string
		expectedErrorString string
	}{
		{
			about:               "cause error - it prints error as string",
			errorFunc:           returnCause,
			formatString:        "%s",
			expectedErrorString: "cause: error",
		},
		{
			about:               "cause error - it prints error as value",
			errorFunc:           returnCause,
			formatString:        "%v",
			expectedErrorString: "cause: error",
		},
		{
			about:               "cause error - it prints error as value",
			errorFunc:           returnCause,
			formatString:        "%#v",
			expectedErrorString: "cause: error",
		},
		{
			about:               "cause error - it prints error as value and stack trace",
			errorFunc:           returnCause,
			formatString:        "%+v",
			expectedErrorString: "cause: error",
		},
		{
			about:               "wrapped cause error - it prints wrapped error as string",
			errorFunc:           returnCauseWrapped,
			formatString:        "%s",
			expectedErrorString: "wrapped",
		},
		{
			about:               "wrapped cause error - it prints wrapped error as value",
			errorFunc:           returnCauseWrapped,
			formatString:        "%v",
			expectedErrorString: "wrapped",
		},
		{
			about:               "wrapped cause error - it prints wrapped error as value with wrapped error",
			errorFunc:           returnCauseWrapped,
			formatString:        "%#v",
			expectedErrorString: "wrapped -> cause: error",
		},
		{
			about:               "wrapped cause error - it prints wrapped error as value with wrapped error and stack trace",
			errorFunc:           returnCauseWrapped,
			formatString:        "%+v",
			expectedErrorString: "wrapped -> cause: error",
		},
		{
			about:               "rewrapped cause error - it prints wrapped error as string",
			errorFunc:           returnCauseWrappedTwice,
			formatString:        "%s",
			expectedErrorString: "rewrapped",
		},
		{
			about:               "rewrapped cause error - it prints wrapped error as value",
			errorFunc:           returnCauseWrappedTwice,
			formatString:        "%v",
			expectedErrorString: "rewrapped",
		},
		{
			about:               "rewrapped cause error - it prints wrapped error as value with wrapped error",
			errorFunc:           returnCauseWrappedTwice,
			formatString:        "%#v",
			expectedErrorString: "rewrapped -> wrapped -> cause: error",
		},
		{
			about:               "rewrapped cause error - it prints wrapped error as value with wrapped error and stack trace",
			errorFunc:           returnCauseWrappedTwice,
			formatString:        "%+v",
			expectedErrorString: "rewrapped -> wrapped -> cause: error",
		},
		// golang errors.
		{
			about:               "wrapped golang error - it prints error as string",
			errorFunc:           returnGolangWrapped,
			formatString:        "%s",
			expectedErrorString: "wrapped",
		},
		{
			about:               "wrapped golang error - it prints error as value",
			errorFunc:           returnGolangWrapped,
			formatString:        "%v",
			expectedErrorString: "wrapped",
		},
		{
			about:               "wrapped golang error - it prints error as value with wrapped error",
			errorFunc:           returnGolangWrapped,
			formatString:        "%#v",
			expectedErrorString: "wrapped -> golang: error",
		},
		{
			about:               "wrapped golang error - it prints error as value with wrapped error with stack trace",
			errorFunc:           returnGolangWrapped,
			formatString:        "%+v",
			expectedErrorString: "wrapped -> golang: error",
		},
		{
			about:               "rewrapped golang error - it prints error as string",
			errorFunc:           returnGolangWrappedTwice,
			formatString:        "%s",
			expectedErrorString: "rewrapped",
		},
		{
			about:               "rewrapped golang error - it prints error as value",
			errorFunc:           returnGolangWrappedTwice,
			formatString:        "%v",
			expectedErrorString: "rewrapped",
		},
		{
			about:               "rewrapped golang error - it prints error as value with wrapped error",
			errorFunc:           returnGolangWrappedTwice,
			formatString:        "%#v",
			expectedErrorString: "rewrapped -> wrapped -> golang: error",
		},
		{
			about:               "rewrapped golang error - it prints error as value with wrapped error with stack trace",
			errorFunc:           returnGolangWrappedTwice,
			formatString:        "%+v",
			expectedErrorString: "rewrapped -> wrapped -> golang: error",
		},
	}
	for index := range testCases {
		testCase := testCases[index]
		t.Run(testCase.about, func(t *testing.T) {
			t.Parallel()

			err := testCase.errorFunc()

			actual := fmt.Sprintf(testCase.formatString, err)

			if testCase.formatString == "%+v" {
				if actual != testCase.expectedErrorString {
					t.Errorf("expected '%s', but got '%s'", testCase.expectedErrorString, actual)
				}
			} else {
				if actual != testCase.expectedErrorString {
					t.Errorf("expected '%s', but got '%s'", testCase.expectedErrorString, actual)
				}
			}
		})
	}
}
