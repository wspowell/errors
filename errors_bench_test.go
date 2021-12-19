package errors_test

import (

	// nolint:depguard // reason: importing "errors" for bench comparison
	goerrors "errors"
	"fmt"
	"io"
	"testing"

	"github.com/wspowell/errors"
)

var errWrapped = errors.New("wrapped")

func Benchmark_errors_New(b *testing.B) {
	var err error
	for i := 0; i < b.N; i++ {
		// nolint:errcheck // reason: error created for test
		err = errors.New("code", "test")
	}

	// Ensure that the compiler is not optimizing away the call.
	b.StopTimer()
	fmt.Fprintf(io.Discard, "%s", err)
}

func Benchmark_goerrors_New(b *testing.B) {
	var err error
	for i := 0; i < b.N; i++ {
		// nolint:errcheck,goerr113,govet // reason: error created for test
		err = goerrors.New("test")
	}

	// Ensure that the compiler is not optimizing away the call.
	b.StopTimer()
	fmt.Fprintf(io.Discard, "%s", err)
}

func Benchmark_goerrors_Wrap(b *testing.B) {
	// nolint:goerr113 // reason: error created for test
	e := goerrors.New("test")
	var err error
	for i := 0; i < b.N; i++ {
		// nolint:govet // reason: error created for test
		err = fmt.Errorf("%w", e)
	}

	// Ensure that the compiler is not optimizing away the call.
	b.StopTimer()
	fmt.Fprintf(io.Discard, "%s", err)
}

func Benchmark_errors_Wrap_cause(b *testing.B) {
	e := errors.New("code", "test")
	var err error
	for i := 0; i < b.N; i++ {
		// nolint:errcheck // reason: error created for test
		err = errors.Wrap(e, errWrapped)
	}

	// Ensure that the compiler is not optimizing away the call.
	b.StopTimer()
	fmt.Fprintf(io.Discard, "%s", err)
}

func Benchmark_errors_Wrap_goerror(b *testing.B) {
	// nolint:goerr113 // reason: error created for test
	e := goerrors.New("test")
	var err error
	for i := 0; i < b.N; i++ {
		// nolint:errcheck // reason: error created for test
		err = errors.Wrap(e, errWrapped)
	}

	// Ensure that the compiler is not optimizing away the call.
	b.StopTimer()
	fmt.Fprintf(io.Discard, "%s", err)
}
