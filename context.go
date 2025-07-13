package geos

// #include <stdlib.h>
// #include "go-geos.h"
import "C"

import (
	"runtime"
	"sync"
	"sync/atomic"
	"unsafe"
)

// A Context is a context.
type Context struct {
	mutex              sync.Mutex
	cHandle            C.GEOSContextHandle_t
	refCount           *atomic.Int64
	ewkbWithSRIDWriter func() *WKBWriter
	geoJSONReader      func() *GeoJSONReader
	geoJSONWriter      func() *GeoJSONWriter
	wkbWriter          func() *WKBWriter
	wkbReader          func() *WKBReader
	wktReader          func() *WKTReader
	wktWriter          func() *WKTWriter
	err                error
}

// NewContext returns a new Context.
func NewContext() *Context {
	cHandle := C.GEOS_init_r()
	var refCount atomic.Int64
	c := &Context{
		cHandle:  cHandle,
		refCount: &refCount,
	}
	c.ref()
	runtime.AddCleanup(c, func(cHandle C.GEOSContextHandle_t) {
		// Inline unref here so that the cleanup function does not hold a
		// reference to c.
		if refCount.Add(-1) == 0 {
			C.finishGEOS_r(cHandle)
		}
	}, cHandle)
	c.ewkbWithSRIDWriter = sync.OnceValue(func() *WKBWriter {
		return c.NewWKBWriter(
			WithWKBWriterFlavor(WKBFlavorExtended),
			WithWKBWriterIncludeSRID(true),
		)
	})
	c.geoJSONReader = sync.OnceValue(func() *GeoJSONReader {
		return c.NewGeoJSONReader()
	})
	c.geoJSONWriter = sync.OnceValue(func() *GeoJSONWriter {
		return c.NewGeoJSONWriter()
	})
	c.wkbReader = sync.OnceValue(func() *WKBReader {
		return c.NewWKBReader()
	})
	c.wkbWriter = sync.OnceValue(func() *WKBWriter {
		return c.NewWKBWriter()
	})
	c.wktReader = sync.OnceValue(func() *WKTReader {
		return c.NewWKTReader()
	})
	c.wktWriter = sync.OnceValue(func() *WKTWriter {
		return c.NewWKTWriter()
	})
	// FIXME golangci-lint complains about the following line saying: Error:
	// dupSubExpr: suspicious identical LHS and RHS for `==` operator (gocritic)
	// As the line does not contain an `==` operator, disable gocritic on this
	// line.
	//nolint:gocritic
	C.GEOSContext_setErrorMessageHandler_r(c.cHandle, C.GEOSMessageHandler_r(C.c_errorMessageHandler), unsafe.Pointer(&c.err))
	return c
}

// Clone clones g into c.
func (c *Context) Clone(g *Geom) *Geom {
	if g.context == c {
		return g.Clone()
	}
	// FIXME use a more intelligent method than a WKB roundtrip (although a WKB
	// roundtrip might actually be quite fast if the cgo overhead is
	// significant)
	clone, err := c.NewGeomFromWKB(g.ToEWKBWithSRID())
	if err != nil {
		panic(err)
	}
	return clone
}

// NewGeomFromGeoJSON returns a new geometry in JSON format from json.
func (c *Context) NewGeomFromGeoJSON(geoJSON string) (*Geom, error) {
	return c.geoJSONReader().ReadGeometry(geoJSON)
}

// NewGeomFromWKB parses a geometry in WKB format from wkb.
func (c *Context) NewGeomFromWKB(wkb []byte) (*Geom, error) {
	return c.wkbReader().Read(wkb)
}

// NewGeomFromWKT parses a geometry in WKT format from wkt.
func (c *Context) NewGeomFromWKT(wkt string) (*Geom, error) {
	return c.wktReader().Read(wkt)
}

// OrientationIndex returns the orientation index from A to B and then to P.
func (c *Context) OrientationIndex(ax, ay, bx, by, px, py float64) int {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return int(C.GEOSOrientationIndex_r(c.cHandle, C.double(ax), C.double(ay), C.double(bx), C.double(by), C.double(px), C.double(py)))
}

// Polygonize returns a set of geometries which contains linework that
// represents the edges of a planar graph.
func (c *Context) Polygonize(geoms []*Geom) *Geom {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	cGeoms, unlockFunc := c.cGeomsLocked(geoms)
	defer unlockFunc()
	return c.newNonNilGeom(C.GEOSPolygonize_r(c.cHandle, cGeoms, C.uint(len(geoms))), nil)
}

// PolygonizeValid returns a set of polygons which contains linework that
// represents the edges of a planar graph.
func (c *Context) PolygonizeValid(geoms []*Geom) *Geom {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	cGeoms, unlockFunc := c.cGeomsLocked(geoms)
	defer unlockFunc()
	return c.newNonNilGeom(C.GEOSPolygonize_valid_r(c.cHandle, cGeoms, C.uint(len(geoms))), nil)
}

// RelatePatternMatch returns if two DE9IM patterns are consistent.
func (c *Context) RelatePatternMatch(mat, pat string) bool {
	matCStr := C.CString(mat)
	defer C.free(unsafe.Pointer(matCStr))
	patCStr := C.CString(pat)
	defer C.free(unsafe.Pointer(patCStr))
	c.mutex.Lock()
	defer c.mutex.Unlock()
	switch C.GEOSRelatePatternMatch_r(c.cHandle, matCStr, patCStr) {
	case 0:
		return false
	case 1:
		return true
	default:
		panic(c.err)
	}
}

// SegmentIntersection returns the coordinate where two lines intersect.
func (c *Context) SegmentIntersection(ax0, ay0, ax1, ay1, bx0, by0, bx1, by1 float64) (x, y float64, intersection bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	var cx, cy float64
	switch C.GEOSSegmentIntersection_r(c.cHandle,
		C.double(ax0), C.double(ay0), C.double(ax1), C.double(ay1),
		C.double(bx0), C.double(by0), C.double(bx1), C.double(by1),
		(*C.double)(&cx), (*C.double)(&cy)) {
	case 1:
		return cx, cy, true
	case -1:
		return 0, 0, false
	default:
		panic(c.err)
	}
}

func (c *Context) cGeomsLocked(geoms []*Geom) (**C.struct_GEOSGeom_t, func()) {
	if len(geoms) == 0 {
		return nil, func() {}
	}
	uniqueContexts := map[*Context]struct{}{c: {}}
	var extraContexts []*Context
	cGeoms := make([]*C.struct_GEOSGeom_t, len(geoms))
	for i := range cGeoms {
		geom := geoms[i]
		if _, ok := uniqueContexts[geom.context]; !ok {
			geom.context.mutex.Lock()
			uniqueContexts[geom.context] = struct{}{}
			extraContexts = append(extraContexts, geom.context)
		}
		cGeoms[i] = geom.cGeom
	}
	return &cGeoms[0], func() {
		for i := len(extraContexts) - 1; i >= 0; i-- {
			extraContexts[i].mutex.Unlock()
		}
	}
}

// ref increases c's reference count by 1.
func (c *Context) ref() {
	c.refCount.Add(1)
}

// unref decreases c's reference count by 1 and finishes c if its reference
// count becomes zero.
func (c *Context) unref() {
	if c.refCount.Add(-1) == 0 {
		C.finishGEOS_r(c.cHandle)
	}
}

//export go_errorMessageHandler
func go_errorMessageHandler(message *C.char, userdata unsafe.Pointer) {
	errP := (*error)(userdata)
	*errP = Error(C.GoString(message))
}
