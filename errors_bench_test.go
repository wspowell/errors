package errors_test

import (
	goerrors "errors"
	"fmt"
	"testing"

	"github.com/wspowell/errors"
)

//nolint:gochecknoglobals // reason: storage to prevent benchmarks from optimizing away calls
var (
	goErrGLOBAL  error
	errorGLOBAL  errors.Error[TestError]
	outputGLOBAL string
)

func errorFn() errors.Error[TestError] {
	return errors.New(TestErrorInternalFailure)
}

func goErrorFn() error {
	return goerrors.New("InternalFailure")
}

func BenchmarkErrorsNew(b *testing.B) {
	var err errors.Error[TestError]
	for i := 0; i < b.N; i++ {
		err = errors.New(TestErrorInternalFailure)
	}

	// Ensure that the compiler is not optimizing away the call.
	b.StopTimer()
	errorGLOBAL = err
}

func BenchmarkErrorsNewFn(b *testing.B) {
	var err errors.Error[TestError]
	for i := 0; i < b.N; i++ {
		err = errorFn()
	}

	// Ensure that the compiler is not optimizing away the call.
	b.StopTimer()
	errorGLOBAL = err
}

func BenchmarkErrorsOk(b *testing.B) {
	var err errors.Error[TestError]
	for i := 0; i < b.N; i++ {
		err = errors.Ok[TestError]()
	}

	// Ensure that the compiler is not optimizing away the call.
	b.StopTimer()
	errorGLOBAL = err
}

func BenchmarkErrorsHandle(b *testing.B) {
	var errString string
	err := errors.New(TestErrorInternalFailure)
	for i := 0; i < b.N; i++ {
		switch err.Cause {
		case TestErrorMyBad:
			errString = err.Error()
		case TestErrorInternalFailure:
			errString = err.Error()
		}
	}

	// Ensure that the compiler is not optimizing away the call.
	b.StopTimer()
	outputGLOBAL = errString
}

func BenchmarkGoerrorsNew(b *testing.B) {
	var err error
	for i := 0; i < b.N; i++ {
		//nolint:goerr113 // reason: error created for test
		err = goerrors.New("test")
	}

	// Ensure that the compiler is not optimizing away the call.
	b.StopTimer()
	goErrGLOBAL = err
}

func BenchmarkGoerrorsWrap(b *testing.B) {
	//nolint:goerr113 // reason: error created for test
	err := goerrors.New("test")
	for i := 0; i < b.N; i++ {
		err = fmt.Errorf("%w", err)
	}

	// Ensure that the compiler is not optimizing away the call.
	b.StopTimer()
	goErrGLOBAL = err
}

// *** DO NOT RUN. CONSUMES MASSIVE AMOUNTS OF MEMORY. ***
// func BenchmarkGoerrorsWrapGo120(b *testing.B) {
// 	//nolint:goerr113 // reason: error created for test
// 	var err error
// 	errOther := goerrors.New("test")
// 	for i := 0; i < b.N; i++ {
// 		err = fmt.Errorf("%w %w", err, errOther)
// 	}

// 	// Ensure that the compiler is not optimizing away the call.
// 	b.StopTimer()
// 	goErrGLOBAL = err
// }

func BenchmarkGoErrorsHandle(b *testing.B) {
	var errString string
	err := goerrors.New("test")
	for i := 0; i < b.N; i++ {
		if err != nil {
			errString = err.Error()
		}
	}

	// Ensure that the compiler is not optimizing away the call.
	b.StopTimer()
	outputGLOBAL = errString
}

func BenchmarkGoerrorsWrapIs(b *testing.B) {
	//nolint:goerr113 // reason: error created for test
	err := goerrors.New("test")
	for i := 0; i < b.N; i++ {
		err = fmt.Errorf("%w", err)
	}

	sentinelErr := goerrors.New("other")

	b.ResetTimer()
	var errString string
	for i := 0; i < b.N; i++ {
		if goerrors.Is(err, sentinelErr) {
			errString = err.Error()
		}
	}

	// Ensure that the compiler is not optimizing away the call.
	b.StopTimer()
	outputGLOBAL = errString
}

// *** DO NOT RUN. CONSUMES MASSIVE AMOUNTS OF MEMORY. ***
// func BenchmarkGoerrorsWrapGo120Is(b *testing.B) {
// 	//nolint:goerr113 // reason: error created for test
// 	var err error
// 	errOther := goerrors.New("test")
// 	for i := 0; i < b.N; i++ {
// 		err = fmt.Errorf("%w %w", err, errOther)
// 	}

// 	sentinelErr := goerrors.New("other")

// 	b.ResetTimer()
// 	var errString string
// 	for i := 0; i < b.N; i++ {
// 		if goerrors.Is(err, sentinelErr) {
// 			errString = err.Error()
// 		}
// 	}

// 	// Ensure that the compiler is not optimizing away the call.
// 	b.StopTimer()
// 	outputGLOBAL = errString
// }
