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

## New issues introduced in go1.20

go1.20 introduced the concept of being able to wrap multiple errors and check for those via `errors.Is()` and `Unwrap() []error`. This changes the `error` data type from a linked list to a tree. This is even less performant than the previous error handling and should be avoided at all costs.

# Re-imagining error handling

Goals:
* Simple
* Consistent
* Self-documenting
* Fast
* Readable

To meet our goals, we must throw out everything we know about golang error handling. Type `error` and package `errors` are deprecated with this approach.

## New error handling style

```
type ExampleError uint

const (
	ExampleErrorInternalFailure = ExampleErr(iota + 1)
	ExampleErrorOtherFailure
)

func ExampleFunc() ExampleErr {
	return errors.New(ExampleErrorInternalFailure)
}

// ...
// Handle the error using a switch.
switch err := ExampleFunc(); err.Cause {
case ExampleErrorInternalFailure:
	// Handle the internal failure
case ExampleErrorOtherFailure:
	// Handle the other failure
default:
	// Optional branch to handle the success case.
}

// Otherwise ok.
```

There are several takeaways from this approach:
1) Instead of returning `error` we return `errors.Error[ExampleError]` which is an error typed with an enum. That means we already have self-documenting code.
2) Handling the error uses `switch` with allows linters to inform us when a case is not handled. See: https://golangci-lint.run/usage/linters/#exhaustive
3) `ExampleErr` is just a `uint` instead of an interface, which is a pointer to a value. So this is fast (faster than `error` actually), because there is no heap allocation required.

Treating errors a enums of `uint` rather than strings means that we can focus on the actual error and leave the human readable messaging for later. Lots of error strings are created in golang only to possibly be thrown away later after the error is handled. This new approach simplifies the process.

## Adding human readable messages to errors

There is already a pattern that exists for this. Simply implement the `fmt.Stringer` interface on your cause type:
```
func (self ExampleFnError) String() string {
	switch self {
	case ExampleErrorInternalFailure:
		return "Internal failure"
	case ExampleErrorOtherFailure:
		return "Other failure"
	}

	return "Ok"
}
```

# Extending Error

The Error type can be extended by embedding it into your own struct. Error only provides a baseline and likely needs more information included with each error, but that is to be decided on a per project basis.
```
type MyError[T Causer] struct {
	Error[T]
	Message string // Add a dynamic string to the error.
}
```

# Benchmarks

Take all benchmarks with a bucket of salt.

```
go test -bench=. -benchmem -count=1 -parallel 8 ./
goos: linux
goarch: amd64
pkg: github.com/wspowell/errors
cpu: AMD Ryzen 9 4900HS with Radeon Graphics         
BenchmarkErrorsNew-8            1000000000               0.2402 ns/op          0 B/op          0 allocs/op
BenchmarkErrorsNewFn-8          1000000000               0.5498 ns/op          0 B/op          0 allocs/op
BenchmarkErrorsOk-8             1000000000               0.2390 ns/op          0 B/op          0 allocs/op
BenchmarkErrorsHandle-8         100000000               11.65 ns/op            0 B/op          0 allocs/op
BenchmarkGoerrorsNew-8          39808287                30.02 ns/op           16 B/op          1 allocs/op
BenchmarkGoerrorsWrap-8          5792497               194.0 ns/op            36 B/op          2 allocs/op
BenchmarkGoErrorsHandle-8       834593935                1.604 ns/op           0 B/op          0 allocs/op
BenchmarkGoerrorsWrapIs-8          10000            178858 ns/op               0 B/op          0 allocs/op

```
