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
BenchmarkErr-12                 686074262                1.676 ns/op           0 B/op          0 allocs/op
BenchmarkGoerrorsNew-12         44342046                26.80 ns/op           16 B/op          1 allocs/op
BenchmarkErrFunc-12             655313150                1.772 ns/op           0 B/op          0 allocs/op
BenchmarkGoerrorsFunc-12        42196820                26.93 ns/op           16 B/op          1 allocs/op
BenchmarkErrFormat-12           23672893                49.89 ns/op            4 B/op          1 allocs/op
BenchmarkGoerrorsFormat-12      15810512                75.96 ns/op           20 B/op          2 allocs/op
BenchmarkErrorInto-12           1000000000               1.105 ns/op           0 B/op          0 allocs/op
BenchmarkGoerrorsWrap-12         7936600               166.1 ns/op            36 B/op          2 allocs/op
PASS
ok      github.com/wspowell/errors      11.275s

goos: linux
goarch: amd64
pkg: github.com/wspowell/errors/result
cpu: Intel(R) Core(TM) i7-8700 CPU @ 3.20GHz
BenchmarkResultOk-12              100574             11109 ns/op               0 B/op          0 allocs/op
BenchmarkResultErr-12              47370             24312 ns/op               0 B/op          0 allocs/op
BenchmarkOkResult-12               59538             20152 ns/op               0 B/op          0 allocs/op
BenchmarkGoerrorOk-12             352614              3323 ns/op               0 B/op          0 allocs/op
BenchmarkErrResult-12              37786             31357 ns/op               0 B/op          0 allocs/op
BenchmarkGoerrorErr-12              4060            274782 ns/op          160000 B/op      10000 allocs/op
PASS
ok      github.com/wspowell/errors/result       7.918s
```
