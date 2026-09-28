package geos

// #include "go-geos.h"
import "C"

import (
	"runtime/cgo"
	"unsafe"
)

//export go_contextInterruptCallback
func go_contextInterruptCallback(userdata unsafe.Pointer) C.int {
	fn, ok := (*cgo.Handle)(userdata).Value().(func() bool)
	if !ok || !fn() {
		return 0
	}
	return 1
}

// SetInterruptCallback sets fn to be called periodically during operations on
// c, which are interrupted if it returns true. fn must not use c. It returns a
// function that clears fn.
func (c *Context) SetInterruptCallback(fn func() bool) func() {
	handle := cgo.NewHandle(fn)
	c.mutex.Lock()
	defer c.mutex.Unlock()
	C.GEOSContext_setInterruptCallback_r(c.cHandle, (*C.GEOSContextInterruptCallback)(C.c_contextInterruptCallback), unsafe.Pointer(&handle)) //nolint:gocritic
	return func() {
		c.mutex.Lock()
		defer c.mutex.Unlock()
		C.GEOSContext_setInterruptCallback_r(c.cHandle, nil, nil)
		handle.Delete()
	}
}
