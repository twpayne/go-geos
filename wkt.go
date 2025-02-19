package geos

// #include "go-geos.h"
import "C"

import (
	"runtime"
	"unsafe"
)

// A WKTReader reads geometries in WKT format.
type WKTReader struct {
	context    *Context
	cWKTReader *C.struct_GEOSWKTReader_t
}

// A WKTWriter writes geometries in WKT format.
type WKTWriter struct {
	context    *Context
	cWKTWriter *C.struct_GEOSWKTWriter_t
}

// NewWKTReader returns a new WKTReader.
func (c *Context) NewWKTReader() *WKTReader {
	r := &WKTReader{
		context:    c,
		cWKTReader: C.GEOSWKTReader_create_r(c.cHandle),
	}
	runtime.AddCleanup(r, func(cWKTReader *C.struct_GEOSWKTReader_t) {
		C.GEOSWKTReader_destroy_r(c.cHandle, cWKTReader)
	}, r.cWKTReader)
	return r
}

// Read returns the geometry represented by wkt.
func (r *WKTReader) Read(wkt string) (*Geom, error) {
	cWKT := C.CString(wkt)
	defer C.free(unsafe.Pointer(cWKT))
	return r.context.newGeom(C.GEOSWKTReader_read_r(r.context.cHandle, r.cWKTReader, cWKT), nil), r.context.err
}

// NewWKTWriter returns a new WKTWriter.
func (c *Context) NewWKTWriter() *WKTWriter {
	w := &WKTWriter{
		context:    c,
		cWKTWriter: C.GEOSWKTWriter_create_r(c.cHandle),
	}
	runtime.AddCleanup(w, func(cWKTWriter *C.struct_GEOSWKTWriter_t) {
		C.GEOSWKTWriter_destroy_r(c.cHandle, cWKTWriter)
	}, w.cWKTWriter)
	return w
}

// Write returns the WKT representation of g.
func (w *WKTWriter) Write(g *Geom) string {
	cWKT := C.GEOSWKTWriter_write_r(g.context.cHandle, w.cWKTWriter, g.cGeom)
	defer C.GEOSFree_r(g.context.cHandle, unsafe.Pointer(cWKT))
	return C.GoString(cWKT)
}
