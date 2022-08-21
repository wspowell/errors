package result_test

import (

	//nolint:depguard // reason: importing "errors" for test comparison
	goerrors "errors"
	"testing"

	"github.com/wspowell/errors"
	"github.com/wspowell/errors/result"
)

//nolint:gochecknoglobals // reason: storage to prevent benchmarks from optimizing away calls
var (
	resGLOBAL   int
	errGLOBAL   error
	errorGLOBAL errors.Error[string]
)

func resultOkInt() result.Result[int, errors.Error[string]] {
	return result.Ok[int, errors.Error[string]](1)
}

func errorOkInt() (int, error) {
	return 1, nil
}

func resultErrInt() result.Result[int, errors.Error[string]] {
	return result.Err[int](errors.Some(errErrorFailure))
}

func errorErrInt() (int, error) {
	//nolint:goerr113 // reason: wrapping is not the focus of this comparison.
	return 0, goerrors.New("failure")
}

func BenchmarkResultOk(b *testing.B) {
	var res result.Result[int, errors.Error[string]]

	for i := 0; i < b.N; i++ {
		for k := 0; k < 10000; k++ {
			res = resultOkInt()
		}
	}

	b.StopTimer()

	resGLOBAL, errorGLOBAL = res.Result()
}

func BenchmarkResultErr(b *testing.B) {
	var res result.Result[int, errors.Error[string]]

	for i := 0; i < b.N; i++ {
		for k := 0; k < 10000; k++ {
			res = resultErrInt()
		}
	}

	b.StopTimer()

	resGLOBAL, errorGLOBAL = res.Result()
}

func BenchmarkOkResult(b *testing.B) {
	var res int
	var err errors.Error[string]

	for i := 0; i < b.N; i++ {
		for k := 0; k < 10000; k++ {
			res, err = resultOkInt().Result()
		}
	}

	b.StopTimer()

	resGLOBAL = res
	errorGLOBAL = err
}

func BenchmarkGoerrorOk(b *testing.B) {
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

func BenchmarkErrResult(b *testing.B) {
	var res int
	var err errors.Error[string]

	for i := 0; i < b.N; i++ {
		for k := 0; k < 10000; k++ {
			res, err = resultErrInt().Result()
		}
	}

	b.StopTimer()

	resGLOBAL = res
	errorGLOBAL = err
}

func BenchmarkGoerrorErr(b *testing.B) {
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
