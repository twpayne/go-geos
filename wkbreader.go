package geos

// #include <stdlib.h>
// #include "go-geos.h"
import "C"

import (
	"runtime"
)

// A WKBReader reads geometries from WKB (Well Known Binary).
type WKBReader struct {
	context    *Context
	cWKBReader *C.struct_GEOSWKBReader_t
}

// NewWKBReader returns a new WKBReader.
func (c *Context) NewWKBReader() *WKBReader {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	cWKBReader := C.GEOSWKBReader_create_r(c.cHandle)
	wkbReader := &WKBReader{
		context:    c,
		cWKBReader: cWKBReader,
	}
	c.ref()
	runtime.AddCleanup(c, c.destroyWKBReader, cWKBReader)
	return wkbReader
}

// Read reads a geometry from wkb.
func (r *WKBReader) Read(wkb []byte) (*Geom, error) {
	r.context.mutex.Lock()
	defer r.context.mutex.Unlock()
	wkbCBuf := C.CBytes(wkb)
	defer C.free(wkbCBuf)
	r.context.err = nil
	return r.context.newGeom(C.GEOSWKBReader_read_r(r.context.cHandle, r.cWKBReader, (*C.uchar)(wkbCBuf), C.ulong(len(wkb))), nil), r.context.err
}

func (c *Context) destroyWKBReader(cWKBReader *C.struct_GEOSWKBReader_t) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	C.GEOSWKBReader_destroy_r(c.cHandle, cWKBReader)
	c.unref()
}
