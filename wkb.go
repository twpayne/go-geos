package geos

// #include "go-geos.h"
import "C"

import (
	"runtime"
	"unsafe"
)

// A WKBFlavor is a flavor of WKB.
type WKBFlavor int

// WKB flavors.
const (
	WKBFlavorExtended WKBFlavor = C.GEOS_WKB_EXTENDED
	WKBFlavorISO      WKBFlavor = C.GEOS_WKB_ISO
)

// A WKBReader reads geometries in WKB format.
type WKBReader struct {
	context    *Context
	cWKBReader *C.struct_GEOSWKBReader_t
}

// A WKBWriter writes geometries in WKB format.
type WKBWriter struct {
	context    *Context
	cWKBWriter *C.struct_GEOSWKBWriter_t
}

// A WKBWriterOption sets an option on a WKBWriter.
type WKBWriterOption func(*WKBWriter)

// NewWKBReader returns a new WKBReader.
func (c *Context) NewWKBReader() *WKBReader {
	r := &WKBReader{
		context:    c,
		cWKBReader: C.GEOSWKBReader_create_r(c.cHandle),
	}
	runtime.AddCleanup(r, func(cWKBReader *C.struct_GEOSWKBReader_t) {
		C.GEOSWKBReader_destroy_r(c.cHandle, cWKBReader)
	}, r.cWKBReader)
	return r
}

// Read returns the geometry represented by wkb.
func (r *WKBReader) Read(wkb []byte) (*Geom, error) {
	cWKB := C.CBytes(wkb)
	defer C.free(cWKB)
	return r.context.newGeom(C.GEOSWKBReader_read_r(r.context.cHandle, r.cWKBReader, (*C.uchar)(cWKB), C.ulong(len(wkb))), nil), r.context.err
}

// NewWKBWriter returns a new WKBWriter.
func (c *Context) NewWKBWriter(options ...WKBWriterOption) *WKBWriter {
	w := &WKBWriter{
		context:    c,
		cWKBWriter: C.GEOSWKBWriter_create_r(c.cHandle),
	}
	runtime.AddCleanup(w, func(cWKBWriter *C.struct_GEOSWKBWriter_t) {
		C.GEOSWKBWriter_destroy_r(c.cHandle, cWKBWriter)
	}, w.cWKBWriter)
	for _, option := range options {
		option(w)
	}
	return w
}

// WithFlavor sets the WKB flavor.
func WithFlavor(flavor WKBFlavor) WKBWriterOption {
	return func(wkbWriter *WKBWriter) {
		C.GEOSWKBWriter_setFlavor_r(wkbWriter.context.cHandle, wkbWriter.cWKBWriter, C.int(flavor))
	}
}

// WithIncludeSRID sets whether to include the SRID.
func WithIncludeSRID(includeSRID bool) WKBWriterOption {
	return func(wkbWriter *WKBWriter) {
		C.GEOSWKBWriter_setIncludeSRID_r(wkbWriter.context.cHandle, wkbWriter.cWKBWriter, toInt[C.char](includeSRID))
	}
}

// Write returns the WKB representation of g.
func (w *WKBWriter) Write(g *Geom) []byte {
	var cSize C.size_t
	cWKB := C.GEOSWKBWriter_write_r(g.context.cHandle, w.cWKBWriter, g.cGeom, &cSize)
	defer C.GEOSFree_r(g.context.cHandle, unsafe.Pointer(cWKB))
	return C.GoBytes(unsafe.Pointer(cWKB), C.int(cSize))
}
