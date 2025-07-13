package geos

// #include <stdlib.h>
// #include "go-geos.h"
import "C"

import (
	"runtime"
	"unsafe"
)

// A GeoJSONWriter writes geometries as GeoJSON.
type GeoJSONWriter struct {
	context        *Context
	cGeoJSONWriter *C.struct_GEOSGeoJSONWriter_t
}

// NewGeoJSONWriter returns a new GeoJSONWriter.
func (c *Context) NewGeoJSONWriter() *GeoJSONWriter {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	cGeoJSONWriter := C.GEOSGeoJSONWriter_create_r(c.cHandle)
	geoJSONWriter := &GeoJSONWriter{
		context:        c,
		cGeoJSONWriter: cGeoJSONWriter,
	}
	c.ref()
	runtime.AddCleanup(c, c.destroyGeoJSONWriter, cGeoJSONWriter)
	return geoJSONWriter
}

// WriteGeometry returns the GeoJSON representation of g.
func (w *GeoJSONWriter) WriteGeometry(g *Geom, indent int) string {
	w.context.mutex.Lock()
	defer w.context.mutex.Unlock()
	cGeoJSONStr := C.GEOSGeoJSONWriter_writeGeometry_r(w.context.cHandle, w.cGeoJSONWriter, g.cGeom, C.int(indent))
	defer C.GEOSFree_r(g.context.cHandle, unsafe.Pointer(cGeoJSONStr))
	return C.GoString(cGeoJSONStr)
}

func (c *Context) destroyGeoJSONWriter(cGeoJSONWriter *C.struct_GEOSGeoJSONWriter_t) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	C.GEOSGeoJSONWriter_destroy_r(c.cHandle, cGeoJSONWriter)
	c.unref()
}
