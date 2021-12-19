# errors

Replacement for golang `errors` package.

# Features

## Stack Traces

All errors provide stack traces. This makes it very easy to see exactly where an error ocurred and what called it. Stack traces only print when using formatting code '%+v'.

```
cause error
github.com/wspowell/errors.returnCauseError
	/workspaces/errors/errors_debug_internal_test.go:21
github.com/wspowell/errors.Test_error_Format.func1
	/workspaces/errors/errors_debug_internal_test.go:119
testing.tRunner
	/usr/local/go/src/testing/testing.go:1259
runtime.goexit
	/usr/local/go/src/runtime/asm_amd64.s:1581
```

Note: release builds do not include stack traces for performance. However, all `panic` errors generate stack traces regardless of build mode. After all, your application should never panic under normal conditions.

## Wrapping

It is common to wrap an error with a new message the developer understands, such as:
```
if err != nil {
    return fmt.Errorf("something went wrong: %w", err)
}
```

However, if you wrap this enough it gets to the point where your error messages are long are hard to read. Really, the whole idea was to give the developer some context as to where the error came from. Stack traces embedded in errors can help solve this in conjunction with trace IDs on logs. Additionally, since only one cause error is possible while using `errors`, the cause should be easily identifiable.

For external errors, instead of wrapping an error using '%w', Wrap() embeds the external error and returns a new error that can checked using Is() and As(). This way the original error is not lost and can be viewed using formatting code '%#v' or '%+v'.

Wrapping an external error:
```
var ErrBroken = errors.New("this thing is broken")

...

func foo() error {
    ...
    if err := externalCall(); err != nil {
        return errors.Wrap(err, ErrBroken)
    }
	...
}
```

### Benefit of Wrap() over golang wrapping

Part of error wrapping is sending back an error we need to check on. But what if there was an `io.EOF` error inside a function that handles JSON. In order for your caller to catch this error and handle it, they would have to know the implementation details of your function to know to catch `io.EOF` and handle it as invalid JSON. Unfortunately, nothing in standard golang allows *wrapping an error* and at the same time return a *new* error.

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

This is not ideal because now we have lost all details on the original error. `errors` provides a new function `Wrap()` that takes a "from" `error` and wraps it with a "to" `error` in order to retain the original error. This allows both errors to be formatted properly according to its own `Format()`. Additionally, since `errors` record internal codes and stack traces for all `errors` it creates, this information can be printed as well.

Format strings are saved internally and processed only when `Format()` is invoked. This is because formatting values into strings requires some overhead and not every error will be printed. This is compared to `fmt.Errorf()` which immediately invokes `Format()` even if the error is never printed. 

## Error Formatting

Format verbs and flags
* '%s' - prints only the error message, no internal code, no stack trace
    * ie `rewrapped`
* '%v' - prints the error message with the internal code of the cause error
    * ie `rewrapped`
* '%#v' - prints the error message with all wrapped error in the error chain, left most is latest and right most is the cause
    * ie `rewrapped -> wrapped -> cause`
* '%+v' - (release) identical to '%#v'
    * ie `rewrapped -> wrapped -> cause`
* '%+v' - (debug) prints the error message with all internal codes in the error chain along with the stack trace
    * ie ```
rewrapped -> wrapped
github.com/wspowell/errors.returnGolangWrappedTwice
	/workspaces/errors/errors_internal_test.go:38
github.com/wspowell/errors.Test_error_Format.func1
	/workspaces/errors/errors_debug_internal_test.go:149
testing.tRunner
	/usr/local/go/src/testing/testing.go:1410
runtime.goexit
	/usr/local/go/src/runtime/asm_amd64.s:1571

wrapped -> golang: error
github.com/wspowell/errors.returnGolangWrapped
	/workspaces/errors/errors_internal_test.go:34
github.com/wspowell/errors.returnGolangWrappedTwice
	/workspaces/errors/errors_internal_test.go:38
github.com/wspowell/errors.Test_error_Format.func1
	/workspaces/errors/errors_debug_internal_test.go:149
testing.tRunner
	/usr/local/go/src/testing/testing.go:1410
runtime.goexit
	/usr/local/go/src/runtime/asm_amd64.s:1571

golang: error
(no stack trace available)
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
2. Always `errors.Wrap()` external/third party errors.
	* The cause of an error should always be within the application.
	* Never use '%w'.
3. Should `return err` if the error originated from within the same application.
	* Except if a function wanted to create a new error for consumer error handling.

# Benchmarks

Take all benchmarks with a bucket of salt.

Debug
```
go test -bench=. -benchmem -count=1 -parallel 8

goos: linux
goarch: amd64
pkg: github.com/wspowell/errors
cpu: AMD Ryzen 9 4900HS with Radeon Graphics         
Benchmark_errors_New-8                   1409652               862.6 ns/op           376 B/op          4 allocs/op
Benchmark_goerrors_New-8                44891362                25.39 ns/op           16 B/op          1 allocs/op
Benchmark_goerrors_Wrap-8                7561659               161.3 ns/op            36 B/op          2 allocs/op
Benchmark_errors_Wrap_cause-8            1232857               970.9 ns/op           376 B/op          4 allocs/op
Benchmark_errors_Wrap_goerror-8          1242308               983.0 ns/op           376 B/op          4 allocs/op
```

Release
```
go test -bench=. -benchmem -count=1 -parallel 8 -tags release

goos: linux
goarch: amd64
pkg: github.com/wspowell/errors
cpu: AMD Ryzen 9 4900HS with Radeon Graphics         
Benchmark_errors_New-8                  15377871                79.15 ns/op           96 B/op          2 allocs/op
Benchmark_goerrors_New-8                48064021                25.40 ns/op           16 B/op          1 allocs/op
Benchmark_goerrors_Wrap-8                7481930               159.6 ns/op            36 B/op          2 allocs/op
Benchmark_errors_Wrap_cause-8           14293320                79.63 ns/op           96 B/op          2 allocs/op
Benchmark_errors_Wrap_goerror-8         15050116                81.20 ns/op           96 B/op          2 allocs/op
PASS
ok      github.com/wspowell/errors      6.446s
```