//go:build !release
// +build !release

package errors

import "runtime"

func callers() *stack {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(4, pcs[:])
	var st stack = pcs[0:n]

	return &st
}
