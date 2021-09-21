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
	return Propagate("fooB", fooA())
}

func fooA2() error {
	return errFooA2
}

func fooB2() error {
	return Propagate("fooB", fooA2())
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
	for index := range testCases {
		testCase := testCases[index]
		t.Run(testCase.about, func(t *testing.T) {
			t.Parallel()

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

func Test_Propagate(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		about                string
		internalCode         string
		err                  error
		expectedInternalCode string
		expectedErrorString  string
	}{
		{
			about:                "it creates a new propagated error with the given cause",
			internalCode:         "ER1001",
			err:                  New("ER1000", "whoops"),
			expectedInternalCode: "ER1000",
			expectedErrorString:  "whoops",
		},
		{
			about:                "it creates a new propagated error with the given error",
			internalCode:         "ER1001",
			err:                  fmt.Errorf("whoops"),
			expectedInternalCode: "ER1001",
			expectedErrorString:  "whoops",
		},
		{
			about:                "it creates a new propagated error when internal code is empty",
			internalCode:         "",
			err:                  fmt.Errorf("whoops"),
			expectedInternalCode: "",
			expectedErrorString:  "whoops",
		},
		{
			about:                "it creates a new propagated error when error is nil",
			internalCode:         "ER1001",
			err:                  nil,
			expectedInternalCode: "ER1001",
			expectedErrorString:  "%!s(<nil>)",
		},
	}
	for index := range testCases {
		testCase := testCases[index]
		t.Run(testCase.about, func(t *testing.T) {
			t.Parallel()

			err := Propagate(testCase.internalCode, testCase.err)

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

var (
	DiscreteErr       = New("DISCRETE", "concrete error")
	DiscreteGolangErr = fmt.Errorf("concrete golang error")
)

func Test_Convert(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		about                string
		fromErr              error
		toErr                error
		expectedInternalCode string
		expectedErrorString  string
	}{
		{
			about:                "it creates a new converted error with the given cause",
			fromErr:              New("ER1000", "whoops"),
			toErr:                DiscreteErr,
			expectedInternalCode: "DISCRETE",
			expectedErrorString:  "concrete error",
		},
		{
			about:                "it creates a new converted error with the given error",
			fromErr:              fmt.Errorf("whoops"),
			toErr:                DiscreteErr,
			expectedInternalCode: "DISCRETE",
			expectedErrorString:  "concrete error",
		},
		{
			about:                "it creates a new converted error when error is nil",
			fromErr:              nil,
			toErr:                DiscreteErr,
			expectedInternalCode: "DISCRETE",
			expectedErrorString:  "concrete error",
		},
		{
			about:                "it passes back original error when discrete error is nil",
			fromErr:              New("ER1000", "whoops"),
			toErr:                nil,
			expectedInternalCode: "ER1000",
			expectedErrorString:  "whoops",
		},
		// Golang errors.
		{
			about:                "it creates a new golang wrapped error with the given cause",
			fromErr:              New("ER1000", "whoops"),
			toErr:                DiscreteGolangErr,
			expectedInternalCode: "ER9999",
			expectedErrorString:  "concrete golang error",
		},
		{
			about:                "it creates a new golang wrapped error with the given error",
			fromErr:              fmt.Errorf("whoops"),
			toErr:                DiscreteGolangErr,
			expectedInternalCode: "ER9999",
			expectedErrorString:  "concrete golang error",
		},
		{
			about:                "it creates a new golang wrapped error when internal code is empty",
			fromErr:              fmt.Errorf("whoops"),
			toErr:                DiscreteGolangErr,
			expectedInternalCode: "ER9999",
			expectedErrorString:  "concrete golang error",
		},
		{
			about:                "it creates a new golang wrapped error when error is nil",
			fromErr:              nil,
			toErr:                DiscreteGolangErr,
			expectedInternalCode: "ER9999",
			expectedErrorString:  "concrete golang error",
		},
		{
			about:                "it passes back original golang error when discrete error is nil",
			fromErr:              fmt.Errorf("whoops"),
			toErr:                nil,
			expectedInternalCode: "ER9999",
			expectedErrorString:  "whoops",
		},
	}
	for index := range testCases {
		testCase := testCases[index]
		t.Run(testCase.about, func(t *testing.T) {
			t.Parallel()

			err := Convert("ER9999", testCase.fromErr, testCase.toErr)

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
			about:         "it shows internalError as cause",
			errorFunc:     fooB,
			expectedCause: errFooA,
		},
		{
			about:         "it shows golang error as cause",
			errorFunc:     fooB2,
			expectedCause: errFooA2,
		},
		{
			about:         "it shows golang error as cause when no cause error",
			errorFunc:     fooA2,
			expectedCause: errFooA2,
		},
		// Propagated errors.
		{
			about: "it shows propagated error as cause",
			errorFunc: func() error {
				return Propagate("ER9999", fooB())
			},
			expectedCause: errFooA,
		},
		{
			about: "it shows propagated golang error as cause",
			errorFunc: func() error {
				return Propagate("ER9999", fooB2())
			},
			expectedCause: errFooA2,
		},
		{
			about: "it shows propagated golang error as cause when no cause error",
			errorFunc: func() error {
				return Propagate("ER9999", fooA2())
			},
			expectedCause: errFooA2,
		},
		// Converted errors.
		{
			about: "it shows converted error as cause",
			errorFunc: func() error {
				return Convert("ER9999", fooB(), DiscreteErr)
			},
			expectedCause: DiscreteErr,
		},
		{
			about: "it shows converted golang error as cause",
			errorFunc: func() error {
				return Convert("ER9999", fooB2(), DiscreteErr)
			},
			expectedCause: DiscreteErr,
		},
		{
			about: "it shows converted golang error as cause when no cause error",
			errorFunc: func() error {
				return Convert("ER9999", fooA2(), DiscreteGolangErr)
			},
			expectedCause: DiscreteGolangErr,
		},
	}
	for index := range testCases {
		testCase := testCases[index]
		t.Run(testCase.about, func(t *testing.T) {
			t.Parallel()

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
		// Propagated errors.
		{
			about: "it shows Propagated error internal code where cause is internalError",
			errorFunc: func() error {
				return Propagate("ER9999", fooB())
			},
			expectedInternalCode: "fooA",
		},
		{
			about: "it shows error internal code where cause is golang error",
			errorFunc: func() error {
				return Propagate("ER9999", fooB2())
			},
			expectedInternalCode: "fooB",
		},
		// Converted errors.
		{
			about: "it shows Propagated error internal code where cause is internalError",
			errorFunc: func() error {
				return Convert("ER9999", fooB(), DiscreteErr)
			},
			expectedInternalCode: "DISCRETE",
		},
		{
			about: "it shows error internal code where cause is golang error",
			errorFunc: func() error {
				return Convert("ER9999", fooB2(), DiscreteErr)
			},
			expectedInternalCode: "DISCRETE",
		},
	}
	for index := range testCases {
		testCase := testCases[index]
		t.Run(testCase.about, func(t *testing.T) {
			t.Parallel()

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

	errA := fooB()
	errA2 := fooB2()

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
		// Propagated errors.
		{
			about: "it unwraps propagated internalError",
			errorFunc: func() error {
				return Propagate("ER9999", errA)
			},
			expectedError: errA,
		},
		{
			about: "it unwraps propagated golang error",
			errorFunc: func() error {
				return Propagate("ER9999", errA2)
			},
			expectedError: errA2,
		},
		// Converted errors.
		{
			about: "it unwraps converted internalError",
			errorFunc: func() error {
				return Convert("ER9999", fooB(), DiscreteErr)
			},
			expectedError: DiscreteErr,
		},
		{
			about: "it unwraps converted golang error",
			errorFunc: func() error {
				return Convert("ER9999", fooB2(), DiscreteErr)
			},
			expectedError: DiscreteErr,
		},
	}
	for index := range testCases {
		testCase := testCases[index]
		t.Run(testCase.about, func(t *testing.T) {
			t.Parallel()

			err := testCase.errorFunc()
			actual := Unwrap(err)
			if actual != testCase.expectedError {
				t.Errorf("expected '%+v', but got '%+v'", testCase.expectedError, actual)
			}

			if testCase.expectedError == errFooA {
				if err.(*propagated).Unwrap() != Unwrap(err) {
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
	for index := range testCases {
		testCase := testCases[index]
		t.Run(testCase.about, func(t *testing.T) {
			t.Parallel()

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
			asErr:         &propagated{},
			expectedError: errFooA,
		},
		{
			about:         "it is a golang error",
			errorFunc:     fooB2,
			asErr:         fmt.Errorf(""),
			expectedError: errFooA2,
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
				t.Errorf("expected '%s', but got '%s'", testCase.expectedError, testCase.asErr)
			}
		})
	}
}
