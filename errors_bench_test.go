package errors_test

import (
	"context"
	// nolint:depguard // reason: importing "errors" for bench comparison
	goerrors "errors"
	"fmt"
	"testing"

	"github.com/wspowell/errors"
)

// nolint:gochecknoglobals // reason: storage to prevent benchmarks from optimizing away calls
var (
	errGLOBAL    error
	errorGLOBAL  errors.Error
	outputGLOBAL string
)

func BenchmarkErrorsNew(b *testing.B) {
	ctx := context.Background()
	var err errors.Error
	for i := 0; i < b.N; i++ {
		err = errors.New(ctx, "test")
	}

	// Ensure that the compiler is not optimizing away the call.
	b.StopTimer()
	errorGLOBAL = err
}

func BenchmarkErrorsWithStackTrace(b *testing.B) {
	ctx := errors.WithStackTrace(context.Background())
	var err errors.Error
	for i := 0; i < b.N; i++ {
		err = errors.New(ctx, "%s", "test")
	}

	// Ensure that the compiler is not optimizing away the call.
	b.StopTimer()
	errorGLOBAL = err
}

func BenchmarkErrorsNewFmt(b *testing.B) {
	ctx := context.Background()
	var err errors.Error
	for i := 0; i < b.N; i++ {
		err = errors.New(ctx, "%s", "test")
	}

	// Ensure that the compiler is not optimizing away the call.
	b.StopTimer()
	errorGLOBAL = err
}

func BenchmarkGoerrorsNew(b *testing.B) {
	var err error
	for i := 0; i < b.N; i++ {
		// nolint:goerr113 // reason: error created for test
		err = goerrors.New("test")
	}

	// Ensure that the compiler is not optimizing away the call.
	b.StopTimer()
	errGLOBAL = err
}

func BenchmarkGoerrorsWrap(b *testing.B) {
	// nolint:goerr113 // reason: error created for test
	err := goerrors.New("test")
	for i := 0; i < b.N; i++ {
		err = fmt.Errorf("%w", err)
	}

	// Ensure that the compiler is not optimizing away the call.
	b.StopTimer()
	errGLOBAL = err
}

func BenchmarkErrorString(b *testing.B) {
	ctx := context.Background()
	err := errors.New(ctx, "test")
	var output string
	for i := 0; i < b.N; i++ {
		output = err.String()
	}

	// Ensure that the compiler is not optimizing away the call.
	b.StopTimer()
	outputGLOBAL = output
}

func BenchmarkErrorStringWithStackTrace(b *testing.B) {
	ctx := errors.WithStackTrace(context.Background())
	err := errors.New(ctx, "test")
	var output string
	for i := 0; i < b.N; i++ {
		output = err.String()
	}

	// Ensure that the compiler is not optimizing away the call.
	b.StopTimer()
	outputGLOBAL = output
}
