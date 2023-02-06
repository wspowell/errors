package result_test

import (
	goerrors "errors"
	"fmt"
	"testing"

	"github.com/wspowell/errors"
	"github.com/wspowell/errors/result"
)

//nolint:gochecknoglobals // reason: storage to prevent benchmarks from optimizing away calls
var (
	resGLOBAL   int
	errGLOBAL   error
	errorGLOBAL errors.Error[TestError]
	isOkGLOBAL  bool
)

func resultOk() result.Result[int, errors.Error[TestError]] {
	return result.Ok[int, errors.Error[TestError]](1)
}

func resultErr() result.Result[int, errors.Error[TestError]] {
	return result.Err[int](errors.New(TestErrorOne))
}

func BenchmarkResultOkResult(b *testing.B) {
	var res int
	var err errors.Error[TestError]

	for i := 0; i < b.N; i++ {
		res, err = resultOk().Result()
	}

	b.StopTimer()

	resGLOBAL = res
	errorGLOBAL = err
}

func BenchmarkResultOkIsOk(b *testing.B) {
	var isOk bool

	for i := 0; i < b.N; i++ {
		isOk = resultOk().IsOk()
	}

	b.StopTimer()

	isOkGLOBAL = isOk
}

func BenchmarkResultErrResult(b *testing.B) {
	var res int
	var err errors.Error[TestError]

	for i := 0; i < b.N; i++ {
		res, err = resultErr().Result()
	}

	b.StopTimer()

	resGLOBAL = res
	errorGLOBAL = err
}

func BenchmarkResultErr(b *testing.B) {
	var res result.Result[int, errors.Error[TestError]]

	for i := 0; i < b.N; i++ {
		res = resultErr()
	}

	b.StopTimer()

	resGLOBAL, errorGLOBAL = res.Result()
}

func errorOk() (int, error) {
	return 1, nil
}

func BenchmarkGoerrorOk(b *testing.B) {
	var res int
	var err error

	for i := 0; i < b.N; i++ {
		res, err = errorOk()
	}

	b.StopTimer()

	resGLOBAL = res
	errGLOBAL = err
}

func BenchmarkGoerrorErr(b *testing.B) {
	var res int
	var err error

	for i := 0; i < b.N; i++ {
		res, err = error1()
	}

	b.StopTimer()

	resGLOBAL = res
	errGLOBAL = err
}

func error1() (int, error) {
	//nolint:goerr113 // reason: wrapping is not the focus of this comparison.
	return 0, goerrors.New("failure")
}

func error2() (int, error) {
	value, err := error1()
	if err != nil {
		//nolint:wrapcheck // reason: this is wrapped, not sure why the linter is mad
		return 0, fmt.Errorf("wrap 1: %w", err)
	}

	return value, nil
}

func error3() (int, error) {
	//nolint:goerr113 // reason: wrapping is not the focus of this comparison.
	value, err := error2()
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
		res, err = error3()
	}

	b.StopTimer()

	resGLOBAL = res
	errGLOBAL = err
}

func errorErr1() (int, error) {
	//nolint:goerr113 // reason: wrapping is not the focus of this comparison.
	return 0, goerrors.New("failure")
}

func errorErr2() (int, error) {
	value, err := errorErr1()
	if err != nil {
		//nolint:wrapcheck // reason: this is wrapped, not sure why the linter is mad
		return 0, fmt.Errorf("wrap 1: %w", err)
	}

	return value, nil
}

func errorErr3() (int, error) {
	//nolint:goerr113 // reason: wrapping is not the focus of this comparison.
	value, err := errorErr2()
	if err != nil {
		//nolint:wrapcheck // reason: this is wrapped, not sure why the linter is mad
		return 0, fmt.Errorf("wrap 2: %w", err)
	}

	return value, nil
}
