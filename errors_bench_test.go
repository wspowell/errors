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
	errorGLOBAL  uint64
	outputGLOBAL string
)

const errSentinel = "test"
const errSentinelFmt = "%s"

func goerrorFn() error {
	//nolint:goerr113 // reason: do not wrap error created for benchmark
	return goerrors.New(errSentinel)
}

func BenchmarkEnumErr(b *testing.B) {
	var err EnumErr
	for i := 0; i < b.N; i++ {
		err = ErrEnumErr1
	}

	// Ensure that the compiler is not optimizing away the call.
	b.StopTimer()
	errorGLOBAL = uint64(err)
}

func BenchmarkErrorsMessageErr(b *testing.B) {
	var err errors.Message[DetailedErr]
	for i := 0; i < b.N; i++ {
		err = errors.NewMessage(ErrDetailedWhoops, "whoops")
	}

	// Ensure that the compiler is not optimizing away the call.
	b.StopTimer()
	errorGLOBAL = uint64(err.Error)
}

func BenchmarkEnumOk(b *testing.B) {
	var err EnumErr
	for i := 0; i < b.N; i++ {
		err = ErrEnumNone
	}

	// Ensure that the compiler is not optimizing away the call.
	b.StopTimer()
	errorGLOBAL = uint64(err)
}

func BenchmarkErrorsMessageOk(b *testing.B) {
	var err errors.Message[DetailedErr]
	for i := 0; i < b.N; i++ {
		err = errors.Ok[DetailedErr]()
	}

	// Ensure that the compiler is not optimizing away the call.
	b.StopTimer()
	errorGLOBAL = uint64(err.Error)
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

func BenchmarkEnumErrFunc(b *testing.B) {
	var err EnumErr
	for i := 0; i < b.N; i++ {
		err = EnumErrFn()
	}

	// Ensure that the compiler is not optimizing away the call.
	b.StopTimer()
	errorGLOBAL = uint64(err)
}

func BenchmarkErrorsMessageErrFunc(b *testing.B) {
	var err errors.Message[DetailedErr]
	for i := 0; i < b.N; i++ {
		err = ErrDetailedErrFn()
	}

	// Ensure that the compiler is not optimizing away the call.
	b.StopTimer()
	errorGLOBAL = uint64(err.Error)
}

func BenchmarkEnumOkFunc(b *testing.B) {
	var err EnumErr
	for i := 0; i < b.N; i++ {
		err = EnumOkFn()
	}

	// Ensure that the compiler is not optimizing away the call.
	b.StopTimer()
	errorGLOBAL = uint64(err)
}

func BenchmarkErrorsMessageOkFunc(b *testing.B) {
	var err errors.Message[DetailedErr]
	for i := 0; i < b.N; i++ {
		err = ErrDetailedOkFn()
	}

	// Ensure that the compiler is not optimizing away the call.
	b.StopTimer()
	errorGLOBAL = uint64(err.Error)
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

func BenchmarkErrorsMessageErrFormat(b *testing.B) {
	var err errors.Message[DetailedErr]
	for i := 0; i < b.N; i++ {
		err = errors.NewMessage(ErrDetailedWhoops, errSentinelFmt, "test")
	}

	// Ensure that the compiler is not optimizing away the call.
	b.StopTimer()
	errorGLOBAL = uint64(err.Error)
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

func BenchmarkErrorsString(b *testing.B) {
	err := errors.NewMessage(ErrDetailedWhoops, "whoops")
	var output string
	for i := 0; i < b.N; i++ {
		output = err.String()
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
