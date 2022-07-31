# errors

Replacement for golang `errors` package.

# Benchmarks

Take all benchmarks with a bucket of salt.

```
go test -bench=. -benchmem -count=1 -parallel 8 ./...

goos: linux
goarch: amd64
pkg: github.com/wspowell/errors
cpu: AMD Ryzen 9 4900HS with Radeon Graphics         
BenchmarkErrorsNew-8                    301735152                3.890 ns/op           0 B/op          0 allocs/op
BenchmarkErrorsWithStackTrace-8          1283198               929.4 ns/op           284 B/op          3 allocs/op
BenchmarkErrorsNewFmt-8                 16286179                69.42 ns/op            4 B/op          1 allocs/op
BenchmarkGoerrorsNew-8                  46670034                25.19 ns/op           16 B/op          1 allocs/op
BenchmarkGoerrorsWrap-8                  6644190               169.2 ns/op            36 B/op          2 allocs/op
BenchmarkErrorString-8                  1000000000               0.7221 ns/op          0 B/op          0 allocs/op
BenchmarkErrorStringWithStackTrace-8      416150              3012 ns/op             400 B/op          9 allocs/op
PASS
ok      github.com/wspowell/errors      11.682s

goos: linux
goarch: amd64
pkg: github.com/wspowell/errors/result
cpu: AMD Ryzen 9 4900HS with Radeon Graphics         
BenchmarkResultOkInt-8            328060              3646 ns/op               0 B/op          0 allocs/op
BenchmarkResultErrInt-8           326503              3643 ns/op               0 B/op          0 allocs/op
BenchmarkGoerrorOkInt-8           244785              4864 ns/op               0 B/op          0 allocs/op
BenchmarkGoerrorErrInt-8          501597              2488 ns/op               0 B/op          0 allocs/op
PASS
ok      github.com/wspowell/errors/result       4.987s
```
