package geos

// #include "go-geos.h"
import "C"

import (
	"errors"
	"fmt"
	"runtime/cgo"
	"unsafe"
)

//export goGEOSInterrupt
func goGEOSInterrupt(userdata unsafe.Pointer) C.int {
	fn, ok := (*cgo.Handle)(userdata).Value().(func() bool)
	if !ok || !fn() {
		return 0
	}
	return 1
}

// SetInterruptCallback sets fn to be called periodically during operations on
// c, which are interrupted if it returns true. fn must not use c. It returns a
// function that clears fn, or an error wrapping [errors.ErrUnsupported] if GEOS
// is older than 3.14.
func (c *Context) SetInterruptCallback(fn func() bool) (func(), error) {
	// checkptr rejects a cgo.Handle converted to unsafe.Pointer, so pass its address.
	hp := new(cgo.Handle)
	*hp = cgo.NewHandle(fn)
	c.mutex.Lock()
	ok := C.c_GEOSContext_setInterruptCallback_r(c.cHandle, unsafe.Pointer(hp)) != 0
	c.mutex.Unlock()
	if !ok {
		hp.Delete()
		return nil, fmt.Errorf("SetInterruptCallback: GEOS %d.%d.%d: %w", VersionMajor, VersionMinor, VersionPatch, errors.ErrUnsupported)
	}
	return func() {
		c.mutex.Lock()
		C.c_GEOSContext_setInterruptCallback_r(c.cHandle, nil)
		c.mutex.Unlock()
		hp.Delete()
	}, nil
}
