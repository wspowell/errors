package errors

import (
	goerrors "errors"
	"fmt"
	"testing"
)

func Benchmark_errors_New(b *testing.B) {
	for i := 0; i < b.N; i++ {
		New("code", "test")
	}
}

func Benchmark_goerrors_New(b *testing.B) {
	for i := 0; i < b.N; i++ {
		goerrors.New("test")
	}
}

func Benchmark_errors_Wrap_cause(b *testing.B) {
	err := New("code", "test")
	for i := 0; i < b.N; i++ {
		Wrap("wrap", err)
	}
}

func Benchmark_errors_Wrap_goerror(b *testing.B) {
	err := goerrors.New("test")
	for i := 0; i < b.N; i++ {
		Wrap("wrap", err)
	}
}

func Benchmark_goerrors_Wrap(b *testing.B) {
	err := goerrors.New("test")
	for i := 0; i < b.N; i++ {
		fmt.Errorf("%w", err)
	}
}
