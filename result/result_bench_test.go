package result_test

import (

	//nolint:depguard // reason: importing "errors" for test comparison
	goerrors "errors"
	"fmt"
	"testing"

	"github.com/wspowell/errors/result"
)

//nolint:gochecknoglobals // reason: storage to prevent benchmarks from optimizing away calls
var (
	resGLOBAL   int
	errGLOBAL   error
	errorGLOBAL FailureError
)

func resultOkInt() result.Result[int, FailureError] {
	return result.Ok[int, FailureError](1)
}

func errorOkInt() (int, error) {
	return 1, nil
}

func resultErrInt() result.Result[int, FailureError] {
	return result.Err[int](errFailureError)
}

func errorErrInt() (int, error) {
	//nolint:goerr113 // reason: wrapping is not the focus of this comparison.
	return 0, goerrors.New("failure")
}

func BenchmarkResultOk(b *testing.B) {
	var res result.Result[int, FailureError]

	for i := 0; i < b.N; i++ {
		for k := 0; k < 10000; k++ {
			res = resultOkInt()
		}
	}

	b.StopTimer()

	resGLOBAL, errorGLOBAL = res.Result()
}

func BenchmarkResultErr(b *testing.B) {
	var res result.Result[int, FailureError]

	for i := 0; i < b.N; i++ {
		for k := 0; k < 10000; k++ {
			res = resultErrInt()
		}
	}

	b.StopTimer()

	resGLOBAL, errorGLOBAL = res.Result()
}

func BenchmarkResultOkResult(b *testing.B) {
	var res int
	var err FailureError

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

func BenchmarkResultErrResult(b *testing.B) {
	var res int
	var err FailureError

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

func goerrorErr1() (int, error) {
	//nolint:goerr113 // reason: wrapping is not the focus of this comparison.
	return 0, goerrors.New("failure")
}

func goerrorErr2() (int, error) {
	value, err := goerrorErr1()
	if err != nil {
		//nolint:wrapcheck // reason: this is wrapped, not sure why the linter is mad
		return 0, fmt.Errorf("wrap 1: %w", err)
	}

	return value, nil
}

func goerrorErr3() (int, error) {
	value, err := goerrorErr2()
	if err != nil {
		//nolint:wrapcheck // reason: this is wrapped, not sure why the linter is mad
		return 0, fmt.Errorf("wrap 2: %w", err)
	}

	return value, nil
}

func BenchmarkGoerrorErrThreeCallsDeep(b *testing.B) {
	var res int
	var err error

	for i := 0; i < b.N; i++ {
		for k := 0; k < 10000; k++ {
			res, err = goerrorErr3()
		}
	}

	b.StopTimer()

	resGLOBAL = res
	errGLOBAL = err
}

func resultErr1() result.Result[int, FailureError] {
	return result.Err[int](errFailureError)
}

func resultErr2() result.Result[int, FailureError] {
	return resultErr1()
}

func resultErr3() result.Result[int, FailureError] {
	return resultErr2()
}

func BenchmarkResultErrThreeCallsDeep(b *testing.B) {
	var res int
	var err FailureError

	for i := 0; i < b.N; i++ {
		for k := 0; k < 10000; k++ {
			res, err = resultErr3().Result()
		}
	}

	b.StopTimer()

	resGLOBAL = res
	errorGLOBAL = err
}

func goerrorOk1() (int, error) {
	return 0, nil
}

func goerrorOk2() (int, error) {
	value, err := goerrorOk1()
	if err != nil {
		//nolint:wrapcheck // reason: this is wrapped, not sure why the linter is mad
		return 0, fmt.Errorf("wrap 1: %w", err)
	}

	return value, nil
}

func goerrorOk3() (int, error) {
	value, err := goerrorOk2()
	if err != nil {
		//nolint:wrapcheck // reason: this is wrapped, not sure why the linter is mad
		return 0, fmt.Errorf("wrap 2: %w", err)
	}

	return value, nil
}

func BenchmarkGoerrorOkThreeCallsDeep(b *testing.B) {
	var res int
	var err error

	for i := 0; i < b.N; i++ {
		for k := 0; k < 10000; k++ {
			res, err = goerrorOk3()
		}
	}

	b.StopTimer()

	resGLOBAL = res
	errGLOBAL = err
}

func resultOk1() result.Result[int, FailureError] {
	return result.Ok[int, FailureError](10)
}

func resultOk2() result.Result[int, FailureError] {
	return resultOk1()
}

func resultOk3() result.Result[int, FailureError] {
	return resultOk2()
}

func BenchmarkResultOkThreeCallsDeep(b *testing.B) {
	var res int
	var err FailureError

	for i := 0; i < b.N; i++ {
		for k := 0; k < 10000; k++ {
			res, err = resultOk3().Result()
		}
	}

	b.StopTimer()

	resGLOBAL = res
	errorGLOBAL = err
}
