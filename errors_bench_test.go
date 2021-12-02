package errors_test

import (

	// nolint:depguard // reason: importing "errors" for bench comparison
	goerrors "errors"
	"fmt"
	"testing"

	"github.com/wspowell/errors"
)

var errConverted = errors.New("DISCRETE", "test")

func Benchmark_errors_New(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// nolint:errcheck // reason: error created for test
		errors.New("code", "test")
	}
}

func Benchmark_goerrors_New(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// nolint:errcheck,goerr113,govet // reason: error created for test
		goerrors.New("test")
	}
}

func Benchmark_errors_Propagate_cause(b *testing.B) {
	err := errors.New("code", "test")
	for i := 0; i < b.N; i++ {
		// nolint:errcheck // reason: error created for test
		errors.Propagate("propagate", err)
	}
}

func Benchmark_errors_Propagate_goerror(b *testing.B) {
	// nolint:goerr113 // reason: error created for test
	err := goerrors.New("test")
	for i := 0; i < b.N; i++ {
		// nolint:errcheck // reason: error created for test
		errors.Propagate("propagate", err)
	}
}

func Benchmark_goerrors_Wrap(b *testing.B) {
	// nolint:goerr113 // reason: error created for test
	err := goerrors.New("test")
	for i := 0; i < b.N; i++ {
		// nolint:govet // reason: error created for test
		fmt.Errorf("%w", err)
	}
}

func Benchmark_errors_Convert_cause(b *testing.B) {
	err := errors.New("code", "test")
	for i := 0; i < b.N; i++ {
		// nolint:errcheck // reason: error created for test
		errors.Convert("convert", err, errConverted)
	}
}

func Benchmark_errors_Convert_goerror(b *testing.B) {
	// nolint:goerr113 // reason: error created for test
	err := goerrors.New("test")
	for i := 0; i < b.N; i++ {
		// nolint:errcheck // reason: error created for test
		errors.Convert("convert", err, errConverted)
	}
}
