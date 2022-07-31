package result_test

import (
	"testing"

	"github.com/wspowell/errors"
	"github.com/wspowell/errors/result"
)

type Error = errors.Error

// nolint:gochecknoglobals // reason: storage to prevent benchmarks from optimizing away calls
var (
	resGLOBAL   int
	errGLOBAL   error
	errorGLOBAL Error
)

func resultOkInt() result.Result[int] {
	return result.Ok(1)
}

func BenchmarkResultOkInt(b *testing.B) {
	var res int
	var err Error

	for i := 0; i < b.N; i++ {
		for k := 0; k < 10000; k++ {
			res, err = resultOkInt().Result()
		}
	}

	b.StopTimer()

	resGLOBAL = res
	errorGLOBAL = err
}

func resultErrInt() result.Result[int] {
	return result.Err[int](errErrorFailure)
}

func BenchmarkResultErrInt(b *testing.B) {
	var res int
	var err Error

	for i := 0; i < b.N; i++ {
		for k := 0; k < 10000; k++ {
			res, err = resultErrInt().Result()
		}
	}

	b.StopTimer()

	resGLOBAL = res
	errorGLOBAL = err
}

func errorOkInt() (int, error) {
	return 1, nil
}

func BenchmarkGoerrorOkInt(b *testing.B) {
	var res int
	var err error

	for i := 0; i < b.N; i++ {
		for k := 0; k < 10000; k++ {
			res, err = errorOkInt()
		}
	}

	b.StopTimer()

	resGLOBAL = res
	errGLOBAL = err
}

func errorErrInt() (int, error) {
	return 0, errGoFailure
}

func BenchmarkGoerrorErrInt(b *testing.B) {
	var res int
	var err error

	for i := 0; i < b.N; i++ {
		for k := 0; k < 10000; k++ {
			res, err = errorErrInt()
		}
	}

	b.StopTimer()

	resGLOBAL = res
	errGLOBAL = err
}
