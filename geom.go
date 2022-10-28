//go:generate go run ./internal/cmds/execute-template -data geommethods.yaml -output geommethods.go geommethods.go.tmpl

package geos

// #include "geos.h"
import "C"

import (
	"unsafe"
)

// A Geom is a geometry.
type Geom struct {
	context          *Context
	geom             *C.struct_GEOSGeom_t
	parent           *Geom
	typeID           GeometryTypeID
	numGeometries    int
	numInteriorRings int
	numPoints        int
}

// Destroy destroys g and releases all resources it holds.
func (g *Geom) Destroy() {
	// Protect against Destroy being called more than once.
	if g == nil || g.context == nil {
		return
	}
	if g.parent == nil {
		g.context.Lock()
		defer g.context.Unlock()
		C.GEOSGeom_destroy_r(g.context.handle, g.geom)
	}
	*g = Geom{} // Clear all references.
}

// Bounds returns g's bounds.
func (g *Geom) Bounds() *Bounds {
	g.mustNotBeDestroyed()
	bounds := NewBoundsEmpty()
	g.context.Lock()
	defer g.context.Unlock()
	C.c_GEOSGeomBounds_r(g.context.handle, g.geom, (*C.double)(&bounds.MinX), (*C.double)(&bounds.MinY), (*C.double)(&bounds.MaxX), (*C.double)(&bounds.MaxY))
	return bounds
}

// CoordSeq returns g's coordinate sequence.
func (g *Geom) CoordSeq() *CoordSeq {
	g.mustNotBeDestroyed()
	g.context.Lock()
	defer g.context.Unlock()
	s := C.GEOSGeom_getCoordSeq_r(g.context.handle, g.geom)
	// Don't set a finalizer as coordSeq is owned by g and will be finalized when g is
	// finalized.
	coordSeq := g.context.newCoordSeq(s, nil)
	coordSeq.parent = g
	return coordSeq
}

// ExteriorRing returns the exterior ring.
func (g *Geom) ExteriorRing() *Geom {
	g.mustNotBeDestroyed()
	g.context.Lock()
	defer g.context.Unlock()
	return g.context.newNonNilGeom(C.GEOSGetExteriorRing_r(g.context.handle, g.geom), g)
}

// Geometry returns the nth geometry of g.
func (g *Geom) Geometry(n int) *Geom {
	g.mustNotBeDestroyed()
	g.context.Lock()
	defer g.context.Unlock()
	if n < 0 || g.numGeometries <= n {
		panic(errIndexOutOfRange)
	}
	return g.context.newNonNilGeom(C.GEOSGetGeometryN_r(g.context.handle, g.geom, C.int(n)), g)
}

// InteriorRing returns the nth interior ring.
func (g *Geom) InteriorRing(n int) *Geom {
	g.mustNotBeDestroyed()
	g.context.Lock()
	defer g.context.Unlock()
	if n < 0 || g.numInteriorRings <= n {
		panic(errIndexOutOfRange)
	}
	return g.context.newNonNilGeom(C.GEOSGetInteriorRingN_r(g.context.handle, g.geom, C.int(n)), g)
}

// Buffer buffers a geometry.
func (g *Geom) Buffer(width float64, quadsegs int) *Geom {
	g.mustNotBeDestroyed()
	g.context.Lock()
	defer g.context.Unlock()
	return g.context.newNonNilGeom(C.GEOSBuffer_r(g.context.handle, g.geom, C.double(width), C.int(quadsegs)), g)
}

// Densify densifies a geometry using a given distance tolerance.
func (g *Geom) Densify(tolerance float64) *Geom {
	g.mustNotBeDestroyed()
	g.context.Lock()
	defer g.context.Unlock()
	return g.context.newNonNilGeom(C.GEOSDensify_r(g.context.handle, g.geom, C.double(tolerance)), g)
}

// Intersection returns the intersection between g and other.
func (g *Geom) Intersection(other *Geom) *Geom {
	g.mustNotBeDestroyed()
	g.context.Lock()
	defer g.context.Unlock()
	if other.context != g.context {
		other.context.Lock()
		defer other.context.Unlock()
	}
	return g.context.newNonNilGeom(C.GEOSIntersection_r(g.context.handle, g.geom, other.geom), nil)
}

// IsValidReason returns the reason that g is invalid.
func (g *Geom) IsValidReason() string {
	g.mustNotBeDestroyed()
	g.context.Lock()
	defer g.context.Unlock()
	reason := C.GEOSisValidReason_r(g.context.handle, g.geom)
	if reason == nil {
		panic(g.context.err)
	}
	defer C.GEOSFree_r(g.context.handle, unsafe.Pointer(reason))
	return C.GoString(reason)
}

// NearestPoints returns the nearest coordinates of g and other. If the nearest
// coordinates do not exist (e.g., when either geom is empty), it returns nil.
func (g *Geom) NearestPoints(other *Geom) [][]float64 {
	g.mustNotBeDestroyed()
	g.context.Lock()
	defer g.context.Unlock()
	if other.context != g.context {
		other.context.Lock()
		defer other.context.Unlock()
	}
	s := C.GEOSNearestPoints_r(g.context.handle, g.geom, other.geom)
	if s == nil {
		return nil
	}
	defer C.GEOSCoordSeq_destroy_r(g.context.handle, s)
	return g.context.newCoordsFromGEOSCoordSeq(s)
}

// NumGeometries returns the number of geometries in g.
func (g *Geom) NumGeometries() int {
	g.mustNotBeDestroyed()
	return g.numGeometries
}

// NumInteriorRings returns the number of interior rings in g.
func (g *Geom) NumInteriorRings() int {
	g.mustNotBeDestroyed()
	return g.numInteriorRings
}

// NumPoints returns the number of points in g.
func (g *Geom) NumPoints() int {
	g.mustNotBeDestroyed()
	return g.numPoints
}

// Point returns the g's nth point.
func (g *Geom) Point(n int) *Geom {
	g.mustNotBeDestroyed()
	g.context.Lock()
	defer g.context.Unlock()
	if n < 0 || g.numPoints <= n {
		panic(errIndexOutOfRange)
	}
	return g.context.newNonNilGeom(C.GEOSGeomGetPointN_r(g.context.handle, g.geom, C.int(n)), nil)
}

// SRID returns g's SRID.
func (g *Geom) SRID() int {
	g.mustNotBeDestroyed()
	g.context.Lock()
	defer g.context.Unlock()
	srid := C.GEOSGetSRID_r(g.context.handle, g.geom)
	// geos_c.h states that GEOSGetSRID_r "Return 0 on exception" but 0 is also
	// returned if the SRID is not set, so we can't rely on it to propagate
	// exceptions.
	return int(srid)
}

// SetSRID sets g's SRID to srid.
func (g *Geom) SetSRID(srid int) *Geom {
	g.mustNotBeDestroyed()
	g.context.Lock()
	defer g.context.Unlock()
	C.GEOSSetSRID_r(g.context.handle, g.geom, C.int(srid))
	return g
}

// String returns g in WKT format.
func (g *Geom) String() string {
	g.mustNotBeDestroyed()
	return g.ToWKT()
}

// ToGeoJSON returns g in GeoJSON format.
func (g *Geom) ToGeoJSON(indent int) string {
	g.mustNotBeDestroyed()
	g.context.Lock()
	defer g.context.Unlock()
	if g.context.geoJSONWriter == nil {
		g.context.geoJSONWriter = C.GEOSGeoJSONWriter_create_r(g.context.handle)
	}
	geoJSONCStr := C.GEOSGeoJSONWriter_writeGeometry_r(g.context.handle, g.context.geoJSONWriter, g.geom, C.int(indent))
	defer C.GEOSFree_r(g.context.handle, unsafe.Pointer(geoJSONCStr))
	return C.GoString(geoJSONCStr)
}

// ToWKB returns g in WKB format.
func (g *Geom) ToWKB() []byte {
	g.mustNotBeDestroyed()
	g.context.Lock()
	defer g.context.Unlock()
	if g.context.wkbWriter == nil {
		g.context.wkbWriter = C.GEOSWKBWriter_create_r(g.context.handle)
	}
	var size C.size_t
	wkbCBuf := C.GEOSWKBWriter_write_r(g.context.handle, g.context.wkbWriter, g.geom, &size)
	defer C.GEOSFree_r(g.context.handle, unsafe.Pointer(wkbCBuf))
	return C.GoBytes(unsafe.Pointer(wkbCBuf), C.int(size))
}

// ToWKT returns g in WKT format.
func (g *Geom) ToWKT() string {
	g.mustNotBeDestroyed()
	g.context.Lock()
	defer g.context.Unlock()
	if g.context.wktWriter == nil {
		g.context.wktWriter = C.GEOSWKTWriter_create_r(g.context.handle)
	}
	wktCStr := C.GEOSWKTWriter_write_r(g.context.handle, g.context.wktWriter, g.geom)
	defer C.GEOSFree_r(g.context.handle, unsafe.Pointer(wktCStr))
	return C.GoString(wktCStr)
}

// Type returns g's type.
func (g *Geom) Type() string {
	g.mustNotBeDestroyed()
	g.context.Lock()
	defer g.context.Unlock()
	typeCStr := C.GEOSGeomType_r(g.context.handle, g.geom)
	if typeCStr == nil {
		panic(g.context.err)
	}
	defer C.GEOSFree_r(g.context.handle, unsafe.Pointer(typeCStr))
	return C.GoString(typeCStr)
}

// TypeID returns g's geometry type id.
func (g *Geom) TypeID() GeometryTypeID {
	g.mustNotBeDestroyed()
	return g.typeID
}

func (g *Geom) finalize() {
	if g.context == nil {
		return
	}
	if g.context.geomFinalizeFunc != nil {
		g.context.geomFinalizeFunc(g)
	}
	g.Destroy()
}

func (g *Geom) mustNotBeDestroyed() {
	if g.context == nil {
		panic("destroyed Geom")
	}
}
