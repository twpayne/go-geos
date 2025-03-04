package geos

// #include "go-geos.h"
import "C"

import (
	"runtime"
	"unsafe"
)

// A GeoJSONReader reads geometries in GeoJSON format.
type GeoJSONReader struct {
	context        *Context
	cGeoJSONReader *C.struct_GEOSGeoJSONReader_t
}

// A GeoJSONWriter writes geometries in GeoJSON format.
type GeoJSONWriter struct {
	context        *Context
	cGeoJSONWriter *C.struct_GEOSGeoJSONWriter_t
}

// NewGeoJSONReader returns a new GeoJSONReader.
func (c *Context) NewGeoJSONReader() *GeoJSONReader {
	r := &GeoJSONReader{
		context:        c,
		cGeoJSONReader: C.GEOSGeoJSONReader_create_r(c.cHandle),
	}
	runtime.AddCleanup(r, func(cGeoJSONReader *C.struct_GEOSGeoJSONReader_t) {
		C.GEOSGeoJSONReader_destroy_r(c.cHandle, cGeoJSONReader)
	}, r.cGeoJSONReader)
	return r
}

// ReadGeometry returns the geometry represented by geoJSON.
func (r *GeoJSONReader) ReadGeometry(geoJSON string) (*Geom, error) {
	cGeoJSON := C.CString(geoJSON)
	defer C.free(unsafe.Pointer(cGeoJSON))
	return r.context.newGeom(C.GEOSGeoJSONReader_readGeometry_r(r.context.cHandle, r.cGeoJSONReader, cGeoJSON), nil), r.context.err
}

// NewGeoJSONWriter returns a new GeoJSONWriter.
func (c *Context) NewGeoJSONWriter() *GeoJSONWriter {
	w := &GeoJSONWriter{
		context:        c,
		cGeoJSONWriter: C.GEOSGeoJSONWriter_create_r(c.cHandle),
	}
	runtime.AddCleanup(w, func(cGeoJSONWriter *C.struct_GEOSGeoJSONWriter_t) {
		C.GEOSGeoJSONWriter_destroy_r(c.cHandle, cGeoJSONWriter)
	}, w.cGeoJSONWriter)
	return w
}

// WriteGeometry returns the GeoJSON representation of g.
func (w *GeoJSONWriter) WriteGeometry(g *Geom, indent int) string {
	cGeoJSON := C.GEOSGeoJSONWriter_writeGeometry_r(g.context.cHandle, w.cGeoJSONWriter, g.cGeom, C.int(indent))
	defer C.GEOSFree_r(g.context.cHandle, unsafe.Pointer(cGeoJSON))
	return C.GoString(cGeoJSON)
}
