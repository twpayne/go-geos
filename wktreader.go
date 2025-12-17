package geos

// #include <stdlib.h>
// #include "go-geos.h"
import "C"

import (
	"runtime"
	"unsafe"
)

// A WKTReader reads geometries from WKT (Well Known Text).
type WKTReader struct {
	context    *Context
	cWKTReader *C.struct_GEOSWKTReader_t
}

// NewWKTReader returns a new WKTReader.
func (c *Context) NewWKTReader() *WKTReader {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	cWKTReader := C.GEOSWKTReader_create_r(c.cHandle)
	wktReader := &WKTReader{
		context:    c,
		cWKTReader: cWKTReader,
	}
	c.ref()
	runtime.AddCleanup(wktReader, c.destroyWKTReader, cWKTReader)
	return wktReader
}

// Read reads a geometry from wkt.
func (r *WKTReader) Read(wkt string) (*Geom, error) {
	r.context.mutex.Lock()
	defer r.context.mutex.Unlock()
	wktCStr := C.CString(wkt)
	defer C.free(unsafe.Pointer(wktCStr))
	r.context.err = nil
	return r.context.newGeom(C.GEOSWKTReader_read_r(r.context.cHandle, r.cWKTReader, wktCStr), nil), r.context.err
}

func (c *Context) destroyWKTReader(cWKTReader *C.struct_GEOSWKTReader_t) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	C.GEOSWKTReader_destroy_r(c.cHandle, cWKTReader)
	c.unref()
}
