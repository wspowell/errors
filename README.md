# errors

Drop-in replacement for golang `errors` package.

# Features

## Internal Codes

Each error has an internal code (string) that is used to trace an error to its origin. Each internal code should be unique and can be in any pattern the developer desires. It could be numbers, "28108", or it could be animal names, "dolphin". The internal code is for the benefit of the developer.

For example:
`"[123] something failed"`

## Stack Traces

All errors provide stack traces. This makes it very easy to see exactly where an error ocurred and what called it.

```
[fooE] concrete error
====[fooE]====
github.com/wspowell/errors.fooE
	/workspaces/errors/errors_debug_internal_test.go:21
github.com/wspowell/errors.Test_error_Format.func1
	/workspaces/errors/errors_debug_internal_test.go:119
testing.tRunner
	/usr/local/go/src/testing/testing.go:1259
runtime.goexit
	/usr/local/go/src/runtime/asm_amd64.s:1581
```

Note: release builds do not include stack traces for performance. However, all `panic`s generate stack traces regardless of build mode. After all, your application should never panic under normal conditions.

## Simpler Wrapping

With all of the features `errors` provides, handling errors becomes easier. 

It is common to wrap an error with a new message the developer understands, such as:
```
if err != nil {
    return fmt.Errorf("something went wrong: %w", err)
}
```

However, if you wrap this enough it gets to the point where your error messages are long are hard to read. Really, the whole idea was to give the developer some context as to where the error came from. The internal code does this for you. Because of this, `errors` provides no such error wrapping and instead only propagates existing errors until they are handled.

Propagating an application error:
```
var ErrBroken = errors.New("app-1", "this thing is broken")
const icFooBrokenFailure = "foo-app-1"

...

func foo() error {
    ...
    if err != nil {
        return errors.Propagate(icFooBrokenFailure, err)
    }
}
```

### Error Conversion

Part of error wrapping is sending back an error we need to check on. But what if there was an `io.EOF` error inside a function that handles JSON. In order for your caller to catch this error and handle it, they would have to know the implementation details of your function to know to catch `io.EOF` and handle it as invalid JSON. Unfortunately, nothing in standard golang allows *propagating an error* and at the same time return a *new* error.

For example, this is impossible in standard golang:
```
var ErrInvalidJson = errors.New("json value is not valid")

...

if errors.Is(err, io.EOF) {
    return fmt.Errorf("%w: %w", ErrInvalidJson, err)
}
```

Since you cannot wrap two errors at the same time, the only recourse you have is wrap the application error and print the err string:
```
var ErrInvalidJson = errors.New("json value is not valid")

...

if errors.Is(err, io.EOF) {
    return fmt.Errorf("%w: %s", ErrInvalidJson, err)
}
```

This is not ideal because now we have lost all details on the original error. `errors` provides a new function `Convert()` that takes a "from" `error` and wraps it with a "to" `error` while retaining the original error. This allows both errors to be formatted properly according to its own `Format()`. Additionally, since `errors` record internal codes and stack traces for all `errors` it creates, this information can be printed as well.

Example converted `error` formatted from an internal test. The first stack trace is the converted error and the next stack trace (after the newline) is the original error stack trace.
```
[fooE] concrete error
====[fooE]====
github.com/wspowell/errors.fooE
	/workspaces/errors/errors_debug_internal_test.go:21
github.com/wspowell/errors.Test_error_Format.func1
	/workspaces/errors/errors_debug_internal_test.go:119
testing.tRunner
	/usr/local/go/src/testing/testing.go:1259
runtime.goexit
	/usr/local/go/src/runtime/asm_amd64.s:1581

[fooD][fooC] whoops: this is bad
====[fooC]====
github.com/wspowell/errors.fooC
	/workspaces/errors/errors_debug_internal_test.go:13
github.com/wspowell/errors.fooD
	/workspaces/errors/errors_debug_internal_test.go:17
github.com/wspowell/errors.fooE
	/workspaces/errors/errors_debug_internal_test.go:21
github.com/wspowell/errors.Test_error_Format.func1
	/workspaces/errors/errors_debug_internal_test.go:119
testing.tRunner
	/usr/local/go/src/testing/testing.go:1259
runtime.goexit
	/usr/local/go/src/runtime/asm_amd64.s:1581
====[fooD]====
github.com/wspowell/errors.fooD
	/workspaces/errors/errors_debug_internal_test.go:17
github.com/wspowell/errors.fooE
	/workspaces/errors/errors_debug_internal_test.go:21
github.com/wspowell/errors.Test_error_Format.func1
	/workspaces/errors/errors_debug_internal_test.go:119
testing.tRunner
	/usr/local/go/src/testing/testing.go:1259
runtime.goexit
	/usr/local/go/src/runtime/asm_amd64.s:1581
```

## Error Formatting

Format verbs and flags
* '%s' - prints only the error message, no internal code, no stack trace
    * ie `error doing the thing`
* '%v' - prints the error message with the internal code of the cause error
    * ie `[aaaa] error doing the thing`
* '%#v' - prints the error message with all internal codes in the error chain, left most is latest and right most is the cause
    * ie `[app-1][1111][aaaa] error doing the thing`
* '%+v' - prints the error message with all internal codes in the error chain along with the stack trace
    * ie ```
[fooD][fooC] whoops: this is bad
====[fooC]====
github.com/wspowell/errors.fooC
	/workspaces/errors/errors_debug_internal_test.go:13
github.com/wspowell/errors.fooD
	/workspaces/errors/errors_debug_internal_test.go:17
github.com/wspowell/errors.Test_error_Format.func1
	/workspaces/errors/errors_debug_internal_test.go:119
testing.tRunner
	/usr/local/go/src/testing/testing.go:1259
runtime.goexit
	/usr/local/go/src/runtime/asm_amd64.s:1581
====[fooD]====
github.com/wspowell/errors.fooD
	/workspaces/errors/errors_debug_internal_test.go:17
github.com/wspowell/errors.Test_error_Format.func1
	/workspaces/errors/errors_debug_internal_test.go:119
testing.tRunner
	/usr/local/go/src/testing/testing.go:1259
runtime.goexit
	/usr/local/go/src/runtime/asm_amd64.s:1581
```

## Catching Panics

Standard practice for catching panics is to use `recover()`. This is clumsy and does not actually provide you with an `error` or useful information unless you go the extra step to grab that info. Instead of duplicating that whole `defer` block each time you need it, just call `errors.Catch()`.

`Catch()` will `defer` the `recover()` and generate an `errors.ErrPanic` for any panic that occurs. You can run this for functions or closures. That makes it easy to run either potentially `panic`-y code in a function or run a goroutine function and catch any `panic`s.

Catch function panics:
```
err := errors.Catch(doThing)

...

// goroutine

go func() {
    if err := errors.Catch(doThing); err != nil {
        // Handle the panic by perhaps logging it.
    }
}
```

Catch closures panics:
```
var result *myStruct
err := error.Catch(func() {
    result = panicyFunction()
})
if err != nil {
    // Handle the panic.
}
```

# Best Practices

1. Always create new errors using `errors.New()` as exported `var`s that can be checked by a consumer.
2. Always keep internal codes unique.
3. Always `errors.Propagate()` all errors, even new ones.
    * This allows new errors to be generic and used from multiple code locations.
    * Errors not `Propagate()`d do not gain the context of where they came from.
4. Always `errors.Convert()` external errors.
    * This keeps all error causes local to the application.
5. Make all new errors have human recognizable internal codes, while propagated errors can use more cryptic internal codes.

# Benchmarks

Take all benchmarks with a bucket of salt.

Debug
```
go test -bench=. -benchmem -count=1 -parallel 8

goos: linux
goarch: amd64
pkg: github.com/wspowell/errors
cpu: AMD Ryzen 9 4900HS with Radeon Graphics         
Benchmark_errors_New-8                   1993179               604.9 ns/op           280 B/op          2 allocs/op
Benchmark_goerrors_New-8                1000000000               0.2399 ns/op          0 B/op          0 allocs/op
Benchmark_errors_Propagate_cause-8       1268781               945.8 ns/op           432 B/op          5 allocs/op
Benchmark_errors_Propagate_goerror-8     1000000              1012 ns/op             496 B/op          6 allocs/op
Benchmark_goerrors_Wrap-8                6335365               185.9 ns/op            36 B/op          2 allocs/op
Benchmark_errors_Convert_cause-8         1548940               775.9 ns/op           392 B/op          4 allocs/op
Benchmark_errors_Convert_goerror-8       1553916               769.2 ns/op           392 B/op          4 allocs/op
```

Release
```
go test -bench=. -benchmem -count=1 -parallel 8 -tags release

goos: linux
goarch: amd64
pkg: github.com/wspowell/errors
cpu: AMD Ryzen 9 4900HS with Radeon Graphics         
Benchmark_errors_New-8                  1000000000               0.2392 ns/op          0 B/op          0 allocs/op
Benchmark_goerrors_New-8                1000000000               0.2466 ns/op          0 B/op          0 allocs/op
Benchmark_errors_Propagate_cause-8       6388952               193.2 ns/op           136 B/op          3 allocs/op
Benchmark_errors_Propagate_goerror-8     4371741               275.8 ns/op           216 B/op          4 allocs/op
Benchmark_goerrors_Wrap-8                6524692               184.3 ns/op            36 B/op          2 allocs/op
Benchmark_errors_Convert_cause-8        13631049                87.24 ns/op          112 B/op          2 allocs/op
Benchmark_errors_Convert_goerror-8      14596249                85.61 ns/op          112 B/op          2 allocs/op
```