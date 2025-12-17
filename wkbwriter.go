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

// A WKBWriter writes geometries as WKB (Well Known Binary).
type WKBWriter struct {
	context    *Context
	cWKBWriter *C.struct_GEOSWKBWriter_t
}

// A WKBWriterOption sets an option on a WKBWriter.
type WKBWriterOption func(*WKBWriter)

// WithWKBWriterFlavor sets the WKB flavor.
func WithWKBWriterFlavor(flavor WKBFlavor) WKBWriterOption {
	return func(w *WKBWriter) {
		C.GEOSWKBWriter_setFlavor_r(w.context.cHandle, w.cWKBWriter, C.int(flavor))
	}
}

// WithWKBWriterIncludeSRID sets whether to include the SRID.
func WithWKBWriterIncludeSRID(includeSRID bool) WKBWriterOption {
	return func(w *WKBWriter) {
		C.GEOSWKBWriter_setIncludeSRID_r(w.context.cHandle, w.cWKBWriter, toInt[C.char](includeSRID))
	}
}

// NewWKBWriter returns a new WKBWriter with the given options.
func (c *Context) NewWKBWriter(options ...WKBWriterOption) *WKBWriter {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	cWKBWriter := C.GEOSWKBWriter_create_r(c.cHandle)
	wkbWriter := &WKBWriter{
		context:    c,
		cWKBWriter: cWKBWriter,
	}
	c.ref()
	runtime.AddCleanup(wkbWriter, c.destroyWKBWriter, cWKBWriter)
	for _, option := range options {
		option(wkbWriter)
	}
	return wkbWriter
}

// Write returns the WKB representation of g.
func (w *WKBWriter) Write(g *Geom) []byte {
	w.context.mutex.Lock()
	defer w.context.mutex.Unlock()
	var size C.size_t
	cWKBBuf := C.GEOSWKBWriter_write_r(g.context.cHandle, w.cWKBWriter, g.cGeom, &size)
	defer C.GEOSFree_r(g.context.cHandle, unsafe.Pointer(cWKBBuf))
	return C.GoBytes(unsafe.Pointer(cWKBBuf), C.int(size))
}

func (c *Context) destroyWKBWriter(cWKBWriter *C.struct_GEOSWKBWriter_t) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	C.GEOSWKBWriter_destroy_r(c.cHandle, cWKBWriter)
	c.unref()
}
