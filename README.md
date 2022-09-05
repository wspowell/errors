# errors

Overhaul of golang error handling. This package proposes a new pattern of handling errors that utilizes existing golang patterns to better effect.

## Current golang error handling patterns

Returning and handling a new error:
```
func foo() error {
    return fmt.Errorf("failure")
}
...
if err := foo(); err != nil {
    // Handle error
}
```

Returning and handling a wrapped error:
```
func foo() error {
    return fmt.Errorf("failure: %w", external.Call())
}
...
if err := foo(); err != nil {
    if errors.Is(err, external.Err) {
        // Handle external error.
    } else {
        // Handle foo error.
    }
}
```

## Issues with golang error handling

Golang error handling is vague. The only standard convention is to check `if err != nil` and then handle it. But once you start wanting to defined sentinel errors or wrapped errors from some other function, it becomes fairly vague as to the correct course of action. Do we check for some other error? Are there errors bubbled up from lower in the call stack that we should handle? How would we know? How do we know the proper way to handle those errors? 

When creating our own errors, do we create errors for "not found" cases? That would mean we might return (nil, ErrNotFound), but how would the caller know that the returned error is not actually an error but a different type of non-success to handle?

Unfortunately, golang error handling leaves much to be desired and has left golang devs scrambling to figure this out for themselves.

# Re-imagining error handling

Goals:
* Simple
* Consistent
* Self-documenting
* Fast
* Readable

To meet our goals, we must throw out everything we know about golang error handling. Type `error` and package `errors` are deprecated with this approach.

## Basic error handling

```
type ExampleErr uint64

const (
	ErrNone ExampleErr = iota
	ErrFailed
)

func ErrFunc() ExampleErr {
	return ErrFailed
}

func ExampleOkFunc() {
	switch err := OkFunc(); err {
	case ErrFailed:
		// Handle error.
		fmt.Println(err)
	case ErrNone:
		fmt.Println("no error")
	}
}
```

There are several takeaways from this approach:
1) Instead of returning `error` we return `ExampleErr` which is an enum. That means we already have self-documenting code.
2) Handling the error uses `switch` with allows linters to inform us when a case is not handled.
3) `ExampleErr` is just a `uint64` instead of an interface, which is a pointer to a value. So this is fast (faster than `error` actually).

Treating errors a enums of `uint64` rather than strings means that we can focus on the actual error and leave the human readable messaging for later. Lots of error strings are created in golang only to possibly be thrown away later after the error is handled. This new approach simplifies the process.

## Adding human readable messages to errors

There is already a pattern that exists for this. Simply attach this function to your error enum:
```
func (self ExampleErr) String() string {
	return [...]string{
		"ok",
		"failure",
	}[self]
}
```


# Benchmarks

Take all benchmarks with a bucket of salt.

```
go test -bench=. -benchmem -count=1 -parallel 8 ./...
goos: linux
goarch: amd64
pkg: github.com/wspowell/errors
cpu: Intel(R) Core(TM) i7-8700 CPU @ 3.20GHz
BenchmarkEnumErr-12                     1000000000               0.2288 ns/op          0 B/op          0 allocs/op
BenchmarkErrorsMessageErr-12            628586955                1.863 ns/op           0 B/op          0 allocs/op
BenchmarkEnumOk-12                      1000000000               0.2241 ns/op          0 B/op          0 allocs/op
BenchmarkErrorsMessageOk-12             1000000000               0.4447 ns/op          0 B/op          0 allocs/op
BenchmarkGoerrorsNew-12                 43195971                27.28 ns/op           16 B/op          1 allocs/op
BenchmarkEnumErrFunc-12                 1000000000               0.4543 ns/op          0 B/op          0 allocs/op
BenchmarkErrorsMessageErrFunc-12        18642760                64.08 ns/op           24 B/op          1 allocs/op
BenchmarkEnumOkFunc-12                  1000000000               0.2236 ns/op          0 B/op          0 allocs/op
BenchmarkErrorsMessageOkFunc-12         1000000000               0.4444 ns/op          0 B/op          0 allocs/op
BenchmarkGoerrorsFunc-12                42288238                26.94 ns/op           16 B/op          1 allocs/op
BenchmarkErrorsMessageErrFormat-12      22842732                49.64 ns/op            4 B/op          1 allocs/op
BenchmarkGoerrorsFormat-12              15409854                75.69 ns/op           20 B/op          2 allocs/op
BenchmarkErrorsString-12                 8896356               131.4 ns/op            72 B/op          2 allocs/op
BenchmarkGoerrorsWrap-12                 7978118               163.7 ns/op            36 B/op          2 allocs/op
PASS
ok      github.com/wspowell/errors      12.502s
goos: linux
goarch: amd64
pkg: github.com/wspowell/errors/result
cpu: Intel(R) Core(TM) i7-8700 CPU @ 3.20GHz
BenchmarkResultOk-12                               98908             11267 ns/op               0 B/op          0 allocs/op
BenchmarkResultErr-12                             108148             11002 ns/op               0 B/op          0 allocs/op
BenchmarkResultOkResult-12                         64800             17879 ns/op               0 B/op          0 allocs/op
BenchmarkGoerrorOk-12                             263226              4475 ns/op               0 B/op          0 allocs/op
BenchmarkResultErrResult-12                        67218             17663 ns/op               0 B/op          0 allocs/op
BenchmarkGoerrorErr-12                              4310            272167 ns/op          160000 B/op      10000 allocs/op
BenchmarkGoerrorErrThreeCallsDeep-12                 397           3111556 ns/op         1200525 B/op      50000 allocs/op
BenchmarkResultErrThreeCallsDeep-12                66409             17617 ns/op               0 B/op          0 allocs/op
BenchmarkGoerrorOkThreeCallsDeep-12                58506             20019 ns/op               0 B/op          0 allocs/op
BenchmarkResultOkThreeCallsDeep-12                 66783             17869 ns/op               0 B/op          0 allocs/op
PASS
ok      github.com/wspowell/errors/result       13.340s
```
