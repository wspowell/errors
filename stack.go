// Copied (and modified) from: github.com/pkg/errors@v0.9.1/stack.go
package errors

import (
	"fmt"
	"io"
	"path"
	"runtime"
	"strconv"
)

const (
	callersSkipPanic = 6
	callersSkipError = 4
)

// frame represents a program counter inside a stack frame.
// For historical reasons if frame is interpreted as a uintptr
// its value represents the program counter + 1.
type frame uintptr

// pc returns the program counter for this frame;
// multiple frames may have the same PC value.
func (f frame) pc() uintptr { return uintptr(f) - 1 }

// file returns the full path to the file that contains the
// function for this Frame's pc.
func (f frame) file() string {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return "unknown"
	}
	file, _ := fn.FileLine(f.pc())

	return file
}

// line returns the line number of source code of the
// function for this Frame's pc.
func (f frame) line() int {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return 0
	}
	_, line := fn.FileLine(f.pc())

	return line
}

// name returns the name of this function, if known.
func (f frame) name() string {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return "unknown"
	}

	return fn.Name()
}

// Format formats the frame according to the fmt.Formatter interface.
//
//    %s    source file
//    %d    source line
//    %v    equivalent to %s:%d
//
// Format accepts flags that alter the printing of some verbs, as follows:
//
//    %+s   function name and path of source file relative to the compile time
//          GOPATH separated by \n\t (<funcname>\n\t<path>)
//    %+v   equivalent to %+s:%d
func (f frame) Format(state fmt.State, verb rune) {
	switch verb {
	case 's':
		switch {
		case state.Flag('+'):
			io.WriteString(state, f.name()) //nolint:errcheck // reason: no real action to take
			io.WriteString(state, "\n\t")   //nolint:errcheck // reason: no real action to take
			io.WriteString(state, f.file()) //nolint:errcheck // reason: no real action to take
		default:
			io.WriteString(state, path.Base(f.file())) //nolint:errcheck // reason: no real action to take
		}
	case 'd':
		io.WriteString(state, strconv.Itoa(f.line())) //nolint:errcheck // reason: no real action to take
	case 'v':
		f.Format(state, 's')
		io.WriteString(state, ":") //nolint:errcheck // reason: no real action to take
		f.Format(state, 'd')
	}
}

// stack represents a stack of program counters.
type stack []uintptr

func (s *stack) Format(state fmt.State, verb rune) {
	switch verb {
	case 'v':
		switch {
		case state.Flag('+'):
			for _, pc := range *s {
				f := frame(pc)
				fmt.Fprintf(state, "\n%+v", f)
			}
		}
	}
}

func callers(skip int) *stack {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(skip, pcs[:])
	var st stack = pcs[0:n]

	return &st
}
