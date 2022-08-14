package result_test

import (
	"context"
	//nolint:depguard // reason: importing "errors" for test comparison
	goerrors "errors"
	"testing"

	"github.com/wspowell/errors"
	"github.com/wspowell/errors/result"
)

type Error = errors.Error[string]

//nolint:gochecknoglobals // reason: storage to prevent benchmarks from optimizing away calls
var (
	resGLOBAL   int
	errGLOBAL   error
	errorGLOBAL Error
)

func resultOkInt() result.Result[int, Error] {
	return result.Ok[int, Error](1)
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

func resultErrInt(ctx context.Context) result.Result[int, Error] {
	return result.Err[int](errors.New(ctx, errErrorFailure))
}

func BenchmarkResultErrIntResult(b *testing.B) {
	var res int
	var err Error

	ctx := context.Background()

	for i := 0; i < b.N; i++ {
		for k := 0; k < 10000; k++ {
			res, err = resultErrInt(ctx).Result()
		}
	}

	b.StopTimer()

	resGLOBAL = res
	errorGLOBAL = err
}
func BenchmarkResultErrInt(b *testing.B) {
	var res result.Result[int, Error]

	ctx := context.Background()

	for i := 0; i < b.N; i++ {
		for k := 0; k < 10000; k++ {
			res = resultErrInt(ctx)
		}
	}

	b.StopTimer()

	resGLOBAL, errorGLOBAL = res.Result()
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
	//nolint:goerr113 // reason: wrapping is not the focus of this comparison.
	return 0, goerrors.New("failure")
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
