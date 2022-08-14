package errors_test

import (
	"context"
	//nolint:depguard // reason: importing "errors" for bench comparison
	goerrors "errors"
	"fmt"
	"testing"

	"github.com/wspowell/errors"
)

//nolint:gochecknoglobals // reason: storage to prevent benchmarks from optimizing away calls
var (
	errGLOBAL    error
	errorGLOBAL  errors.Error
	outputGLOBAL string
)

const errSentinel = "test"
const errSentinelFmt = "%s"

func errorFn(ctx context.Context) errors.Error {
	return errors.New(ctx, errSentinel)
}

func BenchmarkErrorsSentinelNew(b *testing.B) {
	ctx := context.Background()
	var err errors.Error
	for i := 0; i < b.N; i++ {
		err = errors.New(ctx, errSentinel)
	}

	// Ensure that the compiler is not optimizing away the call.
	b.StopTimer()
	errorGLOBAL = err
}

func BenchmarkErrorsSentinelNewFn(b *testing.B) {
	ctx := context.Background()
	var err errors.Error
	for i := 0; i < b.N; i++ {
		err = errorFn(ctx)
	}

	// Ensure that the compiler is not optimizing away the call.
	b.StopTimer()
	errorGLOBAL = err
}

func BenchmarkErrorsSentinelNewValues(b *testing.B) {
	ctx := context.Background()
	var err errors.Error
	for i := 0; i < b.N; i++ {
		err = errors.New(ctx, errSentinelFmt, "test")
	}

	// Ensure that the compiler is not optimizing away the call.
	b.StopTimer()
	errorGLOBAL = err
}

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
		err = errors.New(ctx, "test")
	}

	// Ensure that the compiler is not optimizing away the call.
	b.StopTimer()
	errorGLOBAL = err
}

func BenchmarkErrorsNewValues(b *testing.B) {
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
		//nolint:goerr113 // reason: error created for test
		err = goerrors.New("test")
	}

	// Ensure that the compiler is not optimizing away the call.
	b.StopTimer()
	errGLOBAL = err
}

func BenchmarkGoerrorsNewValues(b *testing.B) {
	var err error
	for i := 0; i < b.N; i++ {
		//nolint:goerr113 // reason: error created for test
		err = fmt.Errorf("%s", "test")
	}

	// Ensure that the compiler is not optimizing away the call.
	b.StopTimer()
	errGLOBAL = err
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

func BenchmarkErrorError(b *testing.B) {
	err := errors.New(context.Background(), "test")
	var output string
	for i := 0; i < b.N; i++ {
		output = err.Error()
	}

	// Ensure that the compiler is not optimizing away the call.
	b.StopTimer()
	outputGLOBAL = output
}

func BenchmarkErrorErrorWithStackTrace(b *testing.B) {
	ctx := errors.WithStackTrace(context.Background())
	err := errors.New(ctx, "test")
	var output string
	for i := 0; i < b.N; i++ {
		output = err.Error()
	}

	// Ensure that the compiler is not optimizing away the call.
	b.StopTimer()
	outputGLOBAL = output
}

func BenchmarkErrorFormat(b *testing.B) {
	err := errors.New(context.Background(), "test")
	var output string
	for i := 0; i < b.N; i++ {
		output = fmt.Sprintf("%+v", err)
	}

	// Ensure that the compiler is not optimizing away the call.
	b.StopTimer()
	outputGLOBAL = output
}

func BenchmarkErrorFormatWithStackTrace(b *testing.B) {
	ctx := errors.WithStackTrace(context.Background())
	err := errors.New(ctx, "test")
	var output string
	for i := 0; i < b.N; i++ {
		output = fmt.Sprintf("%+v", err)
	}

	// Ensure that the compiler is not optimizing away the call.
	b.StopTimer()
	outputGLOBAL = output
}
