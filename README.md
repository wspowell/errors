# errors

Replacement for golang `errors` package.

# Benchmarks

Take all benchmarks with a bucket of salt.

```
go test -bench=. -benchmem -count=1 -parallel 8 ./...
goos: linux
goarch: amd64
pkg: github.com/wspowell/errors
cpu: Intel(R) Core(TM) i7-8700 CPU @ 3.20GHz
BenchmarkErr-12                 654749272                1.720 ns/op           0 B/op          0 allocs/op
BenchmarkGoerrorsNew-12         44907604                26.61 ns/op           16 B/op          1 allocs/op
BenchmarkErrFunc-12             650940289                1.783 ns/op           0 B/op          0 allocs/op
BenchmarkGoerrorsFunc-12        43059499                27.05 ns/op           16 B/op          1 allocs/op
BenchmarkErrFormat-12           21598695                50.30 ns/op            4 B/op          1 allocs/op
BenchmarkGoerrorsFormat-12      13807819                77.40 ns/op           20 B/op          2 allocs/op
BenchmarkErrorInto-12           1000000000               1.110 ns/op           0 B/op          0 allocs/op
BenchmarkGoerrorsWrap-12         7857031               159.0 ns/op            36 B/op          2 allocs/op
PASS
ok      github.com/wspowell/errors      10.068s
goos: linux
goarch: amd64
pkg: github.com/wspowell/errors/result
cpu: Intel(R) Core(TM) i7-8700 CPU @ 3.20GHz
BenchmarkResultOk-12                              101038             11289 ns/op               0 B/op          0 allocs/op
BenchmarkResultErr-12                              48570             24485 ns/op               0 B/op          0 allocs/op
BenchmarkResultOkResult-12                         59952             20081 ns/op               0 B/op          0 allocs/op
BenchmarkGoerrorOk-12                             345360              3345 ns/op               0 B/op          0 allocs/op
BenchmarkResultErrResult-12                        35138             33459 ns/op               0 B/op          0 allocs/op
BenchmarkGoerrorErr-12                              4263            267203 ns/op          160000 B/op      10000 allocs/op
BenchmarkGoerrorErrThreeCallsDeep-12                 399           2924053 ns/op         1200521 B/op      50000 allocs/op
BenchmarkResultErrThreeCallsDeep-12                32073             37071 ns/op               0 B/op          0 allocs/op
BenchmarkGoerrorOkThreeCallsDeep-12                52590             22194 ns/op               0 B/op          0 allocs/op
BenchmarkResultOkThreeCallsDeep-12                 66266             17705 ns/op               0 B/op          0 allocs/op
PASS
ok      github.com/wspowell/errors/result       13.792s
```
