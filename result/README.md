# Result

A pattern of return a value/error pair that is inspired by other modern languages. Result allows a function to return a singular value instead of a (value, error) that is commonly found in golang. This creates simpler usage and ensures better API patterns by discouraging a multiple value return, such as (value, flag, error).

# Benchmarks/Performance

goos: linux
goarch: amd64
pkg: github.com/wspowell/errors/result
cpu: AMD Ryzen 9 4900HS with Radeon Graphics         
BenchmarkResultOkResult-8               1000000000               0.5100 ns/op          0 B/op          0 allocs/op
BenchmarkResultOkIsOk-8                 1000000000               0.5086 ns/op          0 B/op          0 allocs/op
BenchmarkResultErrResult-8              1000000000               0.2647 ns/op          0 B/op          0 allocs/op
BenchmarkResultErr-8                    1000000000               0.2509 ns/op          0 B/op          0 allocs/op
BenchmarkGoerrorOk-8                    1000000000               0.2616 ns/op          0 B/op          0 allocs/op
BenchmarkGoerrorErr-8                   39030572                30.39 ns/op           16 B/op          1 allocs/op
BenchmarkGoerrorErrThreeCallsDeep-8      2603626               443.3 ns/op           120 B/op          5 allocs/op