package errors_test

import (
	goerrors "errors"
	"fmt"
	"testing"

	"github.com/wspowell/errors"
)

var errConverted = errors.New("DISCRETE", "test")

func Benchmark_errors_New(b *testing.B) {
	for i := 0; i < b.N; i++ {
		errors.New("code", "test")
	}
}

func Benchmark_goerrors_New(b *testing.B) {
	for i := 0; i < b.N; i++ {
		goerrors.New("test")
	}
}

func Benchmark_errors_Propagate_cause(b *testing.B) {
	err := errors.New("code", "test")
	for i := 0; i < b.N; i++ {
		errors.Propagate("propagate", err)
	}
}

func Benchmark_errors_Propagate_goerror(b *testing.B) {
	err := goerrors.New("test")
	for i := 0; i < b.N; i++ {
		errors.Propagate("propagate", err)
	}
}

func Benchmark_goerrors_Wrap(b *testing.B) {
	err := goerrors.New("test")
	for i := 0; i < b.N; i++ {
		fmt.Errorf("%w", err)
	}
}

func Benchmark_errors_Convert_cause(b *testing.B) {
	err := errors.New("code", "test")
	for i := 0; i < b.N; i++ {
		errors.Convert("convert", err, errConverted)
	}
}

func Benchmark_errors_Convert_goerror(b *testing.B) {
	err := goerrors.New("test")
	for i := 0; i < b.N; i++ {
		errors.Convert("convert", err, errConverted)
	}
}
