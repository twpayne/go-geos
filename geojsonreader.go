package geos

// #include <stdlib.h>
// #include "go-geos.h"
import "C"

import (
	"runtime"
	"unsafe"
)

// A GeoJSONReader reads GeoJSON.
type GeoJSONReader struct {
	context        *Context
	cGeoJSONReader *C.struct_GEOSGeoJSONReader_t
}

// NewGeoJSONReader returns a new GeoJSONReader.
func (c *Context) NewGeoJSONReader() *GeoJSONReader {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	cGeoJSONReader := C.GEOSGeoJSONReader_create_r(c.cHandle)
	geoJSONReader := &GeoJSONReader{
		context:        c,
		cGeoJSONReader: cGeoJSONReader,
	}
	c.ref()
	runtime.AddCleanup(c, c.destroyGeoJSONReader, cGeoJSONReader)
	return geoJSONReader
}

// ReadGeometry reads a geometry from geoJSON.
func (r *GeoJSONReader) ReadGeometry(geoJSON string) (*Geom, error) {
	r.context.mutex.Lock()
	defer r.context.mutex.Unlock()
	geoJSONCStr := C.CString(geoJSON)
	defer C.free(unsafe.Pointer(geoJSONCStr))
	r.context.err = nil
	return r.context.newGeom(C.GEOSGeoJSONReader_readGeometry_r(r.context.cHandle, r.cGeoJSONReader, geoJSONCStr), nil), r.context.err
}

func (c *Context) destroyGeoJSONReader(cGeoJSONReader *C.struct_GEOSGeoJSONReader_t) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	C.GEOSGeoJSONReader_destroy_r(c.cHandle, cGeoJSONReader)
	c.unref()
}
