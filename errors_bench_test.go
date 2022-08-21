package errors_test

import (

	//nolint:depguard // reason: importing "errors" for bench comparison
	goerrors "errors"
	"fmt"
	"testing"

	"github.com/wspowell/errors"
)

//nolint:gochecknoglobals // reason: storage to prevent benchmarks from optimizing away calls
var (
	errGLOBAL    error
	errorGLOBAL  errors.Error[string]
	outputGLOBAL string
)

const errSentinel = "test"
const errSentinelFmt = "%s"

func errorFn() errors.Error[string] {
	return errors.Some(errSentinel)
}
func goerrorFn() error {
	//nolint:goerr113 // reason: do not wrap error created for benchmark
	return goerrors.New(errSentinel)
}

func BenchmarkErr(b *testing.B) {
	var err errors.Error[string]
	for i := 0; i < b.N; i++ {
		err = errors.Some(errSentinel)
	}

	// Ensure that the compiler is not optimizing away the call.
	b.StopTimer()
	errorGLOBAL = err
}
func BenchmarkGoerrorsNew(b *testing.B) {
	var err error
	for i := 0; i < b.N; i++ {
		//nolint:goerr113 // reason: error created for test
		err = goerrors.New(errSentinel)
	}

	// Ensure that the compiler is not optimizing away the call.
	b.StopTimer()
	errGLOBAL = err
}

func BenchmarkErrFunc(b *testing.B) {
	var err errors.Error[string]
	for i := 0; i < b.N; i++ {
		err = errorFn()
	}

	// Ensure that the compiler is not optimizing away the call.
	b.StopTimer()
	errorGLOBAL = err
}

func BenchmarkGoerrorsFunc(b *testing.B) {
	var err error
	for i := 0; i < b.N; i++ {
		err = goerrorFn()
	}

	// Ensure that the compiler is not optimizing away the call.
	b.StopTimer()
	errGLOBAL = err
}

func BenchmarkErrFormat(b *testing.B) {
	var err errors.Error[string]
	for i := 0; i < b.N; i++ {
		err = errors.Some(errSentinelFmt, "test")
	}

	// Ensure that the compiler is not optimizing away the call.
	b.StopTimer()
	errorGLOBAL = err
}

func BenchmarkGoerrorsFormat(b *testing.B) {
	var err error
	for i := 0; i < b.N; i++ {
		//nolint:goerr113 // reason: error created for test
		err = fmt.Errorf(errSentinelFmt, "test")
	}

	// Ensure that the compiler is not optimizing away the call.
	b.StopTimer()
	errGLOBAL = err
}

func BenchmarkErrorInto(b *testing.B) {
	err := errors.Some("test")
	var output string
	for i := 0; i < b.N; i++ {
		output = err.Into()
	}

	// Ensure that the compiler is not optimizing away the call.
	b.StopTimer()
	outputGLOBAL = output
}

func BenchmarkGoerrorsWrap(b *testing.B) {
	//nolint:goerr113 // reason: error created for test
	err := goerrors.New("test")
	for i := 0; i < b.N; i++ {
		err = fmt.Errorf("%w", err)
	}

	// Ensure that the compiler is not optimizing away the call.
	b.StopTimer()
	errGLOBAL = err
}
