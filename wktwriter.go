package geos

// #include "go-geos.h"
import "C"

import (
	"runtime"
	"unsafe"
)

// A WKTWriter writes geometries in WKT (Well Known Text).
type WKTWriter struct {
	context    *Context
	cWKTWriter *C.struct_GEOSWKTWriter_t
}

// NewWKTWriter returns a new WKTWriter.
func (c *Context) NewWKTWriter() *WKTWriter {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	cWKTWriter := C.GEOSWKTWriter_create_r(c.cHandle)
	wktWriter := &WKTWriter{
		context:    c,
		cWKTWriter: cWKTWriter,
	}
	c.ref()
	runtime.AddCleanup(wktWriter, c.destroyWKTWriter, cWKTWriter)
	return wktWriter
}

// Write returns the WKT representation of g.
func (w *WKTWriter) Write(g *Geom) string {
	w.context.mutex.Lock()
	defer w.context.mutex.Unlock()
	cWKTStr := C.GEOSWKTWriter_write_r(w.context.cHandle, w.cWKTWriter, g.cGeom)
	defer C.GEOSFree_r(g.context.cHandle, unsafe.Pointer(cWKTStr))
	return C.GoString(cWKTStr)
}

func (c *Context) destroyWKTWriter(cWKTWriter *C.struct_GEOSWKTWriter_t) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	C.GEOSWKTWriter_destroy_r(c.cHandle, cWKTWriter)
	c.unref()
}
