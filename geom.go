//go:generate go tool execute-template -data geommethods.yaml -output geommethods.go geommethods.go.tmpl

package geos

// #include <stdlib.h>
// #include "go-geos.h"
import "C"

import (
	"unsafe"
)

// A Geom is a geometry.
type Geom struct {
	context          *Context
	cGeom            *C.struct_GEOSGeom_t
	parent           *Geom
	typeID           TypeID
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
		g.context.mutex.Lock()
		defer g.context.mutex.Unlock()
		C.GEOSGeom_destroy_r(g.context.cHandle, g.cGeom)
	}
	*g = Geom{} // Clear all references.
}

// Bounds returns g's bounds.
func (g *Geom) Bounds() *Box2D {
	g.mustNotBeDestroyed()
	bounds := NewBox2DEmpty()
	g.context.mutex.Lock()
	defer g.context.mutex.Unlock()
	C.c_GEOSGeomBounds_r(g.context.cHandle, g.cGeom, (*C.double)(&bounds.MinX), (*C.double)(&bounds.MinY), (*C.double)(&bounds.MaxX), (*C.double)(&bounds.MaxY))
	return bounds
}

// MakeValidWithParams returns a new valid geometry using the MakeValidMethods and MakeValidCollapsed parameters.
func (g *Geom) MakeValidWithParams(method MakeValidMethod, collapse MakeValidCollapsed) *Geom {
	g.mustNotBeDestroyed()
	g.context.mutex.Lock()
	defer g.context.mutex.Unlock()
	cRes := C.c_GEOSMakeValidWithParams_r(g.context.cHandle, g.cGeom, C.enum_GEOSMakeValidMethods(method), C.int(collapse))
	return g.context.newGeom(cRes, nil)
}

// BufferWithParams returns g buffered with bufferParams.
func (g *Geom) BufferWithParams(bufferParams *BufferParams, width float64) *Geom {
	g.context.mutex.Lock()
	defer g.context.mutex.Unlock()
	if bufferParams.context != g.context {
		bufferParams.context.mutex.Lock()
		defer bufferParams.context.mutex.Unlock()
	}
	return g.context.newNonNilGeom(C.GEOSBufferWithParams_r(g.context.cHandle, g.cGeom, bufferParams.cBufParams, C.double(width)), nil)
}

func (g *Geom) ClipByBox2D(box2d *Box2D) *Geom {
	return g.ClipByRect(box2d.MinX, box2d.MinY, box2d.MaxX, box2d.MaxY)
}

// CoordSeq returns g's coordinate sequence.
func (g *Geom) CoordSeq() *CoordSeq {
	g.mustNotBeDestroyed()
	g.context.mutex.Lock()
	defer g.context.mutex.Unlock()
	cCoordSeq := C.GEOSGeom_getCoordSeq_r(g.context.cHandle, g.cGeom)
	// Don't set a finalizer as coordSeq is owned by g and will be finalized when g is
	// finalized.
	coordSeq := g.context.newCoordSeqInternal(cCoordSeq, nil)
	if coordSeq == nil {
		return nil
	}
	coordSeq.parent = g
	return coordSeq
}

// ExteriorRing returns the exterior ring.
func (g *Geom) ExteriorRing() *Geom {
	g.mustNotBeDestroyed()
	g.context.mutex.Lock()
	defer g.context.mutex.Unlock()
	return g.context.newNonNilGeom(C.GEOSGetExteriorRing_r(g.context.cHandle, g.cGeom), g)
}

// Geometry returns the nth geometry of g.
func (g *Geom) Geometry(n int) *Geom {
	g.mustNotBeDestroyed()
	g.context.mutex.Lock()
	defer g.context.mutex.Unlock()
	if n < 0 || g.numGeometries <= n {
		panic(errIndexOutOfRange)
	}
	return g.context.newNonNilGeom(C.GEOSGetGeometryN_r(g.context.cHandle, g.cGeom, C.int(n)), g)
}

// InteriorRing returns the nth interior ring.
func (g *Geom) InteriorRing(n int) *Geom {
	g.mustNotBeDestroyed()
	g.context.mutex.Lock()
	defer g.context.mutex.Unlock()
	if n < 0 || g.numInteriorRings <= n {
		panic(errIndexOutOfRange)
	}
	return g.context.newNonNilGeom(C.GEOSGetInteriorRingN_r(g.context.cHandle, g.cGeom, C.int(n)), g)
}

// IsValidReason returns the reason that g is invalid.
func (g *Geom) IsValidReason() string {
	g.mustNotBeDestroyed()
	g.context.mutex.Lock()
	defer g.context.mutex.Unlock()
	reason := C.GEOSisValidReason_r(g.context.cHandle, g.cGeom)
	if reason == nil {
		panic(g.context.err)
	}
	defer C.GEOSFree_r(g.context.cHandle, unsafe.Pointer(reason))
	return C.GoString(reason)
}

// NearestPoints returns the nearest coordinates of g and other. If the nearest
// coordinates do not exist (e.g., when either geom is empty), it returns nil.
func (g *Geom) NearestPoints(other *Geom) [][]float64 {
	g.mustNotBeDestroyed()
	g.context.mutex.Lock()
	defer g.context.mutex.Unlock()
	s := C.GEOSNearestPoints_r(g.context.cHandle, g.cGeom, other.cGeom)
	if s == nil {
		return nil
	}
	defer C.GEOSCoordSeq_destroy_r(g.context.cHandle, s)
	return g.context.newCoordsFromGEOSCoordSeq(s)
}

func (g *Geom) Normalize() *Geom {
	g.mustNotBeDestroyed()
	g.context.mutex.Lock()
	defer g.context.mutex.Unlock()
	if C.GEOSNormalize_r(g.context.cHandle, g.cGeom) != 0 {
		panic(g.context.err)
	}
	return g
}

// NumCoordinates returns the number of coordinates in g.
func (g *Geom) NumCoordinates() int {
	g.mustNotBeDestroyed()
	g.context.mutex.Lock()
	defer g.context.mutex.Unlock()
	numCoordinates := C.GEOSGetNumCoordinates_r(g.context.cHandle, g.cGeom)
	if numCoordinates == -1 {
		panic(g.context.err)
	}
	return int(numCoordinates)
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
	g.context.mutex.Lock()
	defer g.context.mutex.Unlock()
	if n < 0 || g.numPoints <= n {
		panic(errIndexOutOfRange)
	}
	return g.context.newNonNilGeom(C.GEOSGeomGetPointN_r(g.context.cHandle, g.cGeom, C.int(n)), nil)
}

// PolygonizeFull returns a set of geometries which contains linework that
// represents the edge of a planar graph.
func (g *Geom) PolygonizeFull() (geom, cuts, dangles, invalidRings *Geom) {
	g.mustNotBeDestroyed()
	g.context.mutex.Lock()
	defer g.context.mutex.Unlock()
	var cCuts, cDangles, cInvalidRings *C.struct_GEOSGeom_t
	cGeom := C.GEOSPolygonize_full_r(g.context.cHandle, g.cGeom, &cCuts, &cDangles, &cInvalidRings) //nolint:gocritic
	geom = g.context.newNonNilGeom(cGeom, nil)
	cuts = g.context.newGeom(cCuts, nil)
	dangles = g.context.newGeom(cDangles, nil)
	invalidRings = g.context.newGeom(cInvalidRings, nil)
	return geom, cuts, dangles, invalidRings
}

// Precision returns g's precision.
func (g *Geom) Precision() float64 {
	g.mustNotBeDestroyed()
	g.context.mutex.Lock()
	defer g.context.mutex.Unlock()
	return float64(C.GEOSGeom_getPrecision_r(g.context.cHandle, g.cGeom))
}

// RelatePattern returns if the DE9IM pattern for g and other matches pat.
func (g *Geom) RelatePattern(other *Geom, pat string) bool {
	g.mustNotBeDestroyed()
	patCStr := C.CString(pat)
	defer C.free(unsafe.Pointer(patCStr))
	g.context.mutex.Lock()
	defer g.context.mutex.Unlock()
	switch C.GEOSRelatePattern_r(g.context.cHandle, g.cGeom, other.cGeom, patCStr) {
	case 0:
		return false
	case 1:
		return true
	default:
		panic(g.context.err)
	}
}

// SRID returns g's SRID.
func (g *Geom) SRID() int {
	g.mustNotBeDestroyed()
	g.context.mutex.Lock()
	defer g.context.mutex.Unlock()
	srid := C.GEOSGetSRID_r(g.context.cHandle, g.cGeom)
	// geos_c.h states that GEOSGetSRID_r "Return 0 on exception" but 0 is also
	// returned if the SRID is not set, so we can't rely on it to propagate
	// exceptions.
	return int(srid)
}

// SetSRID sets g's SRID to srid.
func (g *Geom) SetSRID(srid int) *Geom {
	g.mustNotBeDestroyed()
	g.context.mutex.Lock()
	defer g.context.mutex.Unlock()
	C.GEOSSetSRID_r(g.context.cHandle, g.cGeom, C.int(srid))
	return g
}

// SetUserData sets g's userdata and returns g.
func (g *Geom) SetUserData(userdata uintptr) *Geom {
	g.mustNotBeDestroyed()
	g.context.mutex.Lock()
	defer g.context.mutex.Unlock()
	C.c_GEOSGeom_setUserData_r(g.context.cHandle, g.cGeom, C.uintptr_t(userdata))
	return g
}

// String returns g in WKT format.
func (g *Geom) String() string {
	g.mustNotBeDestroyed()
	return g.ToWKT()
}

// ToEWKBWithSRID returns g in Extended WKB format with its SRID.
func (g *Geom) ToEWKBWithSRID() []byte {
	g.mustNotBeDestroyed()
	g.context.mutex.Lock()
	defer g.context.mutex.Unlock()
	if g.context.cEWKBWithSRIDWriter == nil {
		g.context.cEWKBWithSRIDWriter = C.GEOSWKBWriter_create_r(g.context.cHandle)
		C.GEOSWKBWriter_setFlavor_r(g.context.cHandle, g.context.cEWKBWithSRIDWriter, C.GEOS_WKB_EXTENDED)
		C.GEOSWKBWriter_setIncludeSRID_r(g.context.cHandle, g.context.cEWKBWithSRIDWriter, 1)
	}
	var size C.size_t
	cEWKBBuf := C.GEOSWKBWriter_write_r(g.context.cHandle, g.context.cEWKBWithSRIDWriter, g.cGeom, &size)
	defer C.GEOSFree_r(g.context.cHandle, unsafe.Pointer(cEWKBBuf))
	return C.GoBytes(unsafe.Pointer(cEWKBBuf), C.int(size))
}

// ToGeoJSON returns g in GeoJSON format.
func (g *Geom) ToGeoJSON(indent int) string {
	g.mustNotBeDestroyed()
	g.context.mutex.Lock()
	defer g.context.mutex.Unlock()
	if g.context.cGeoJSONWriter == nil {
		g.context.cGeoJSONWriter = C.GEOSGeoJSONWriter_create_r(g.context.cHandle)
	}
	cGeoJSONStr := C.GEOSGeoJSONWriter_writeGeometry_r(g.context.cHandle, g.context.cGeoJSONWriter, g.cGeom, C.int(indent))
	defer C.GEOSFree_r(g.context.cHandle, unsafe.Pointer(cGeoJSONStr))
	return C.GoString(cGeoJSONStr)
}

// ToWKB returns g in WKB format.
func (g *Geom) ToWKB() []byte {
	g.mustNotBeDestroyed()
	g.context.mutex.Lock()
	defer g.context.mutex.Unlock()
	if g.context.cWKBWriter == nil {
		g.context.cWKBWriter = C.GEOSWKBWriter_create_r(g.context.cHandle)
	}
	var size C.size_t
	cWKBBuf := C.GEOSWKBWriter_write_r(g.context.cHandle, g.context.cWKBWriter, g.cGeom, &size)
	defer C.GEOSFree_r(g.context.cHandle, unsafe.Pointer(cWKBBuf))
	return C.GoBytes(unsafe.Pointer(cWKBBuf), C.int(size))
}

// ToWKT returns g in WKT format.
func (g *Geom) ToWKT() string {
	g.mustNotBeDestroyed()
	g.context.mutex.Lock()
	defer g.context.mutex.Unlock()
	if g.context.cWKTWriter == nil {
		g.context.cWKTWriter = C.GEOSWKTWriter_create_r(g.context.cHandle)
	}
	cWKTStr := C.GEOSWKTWriter_write_r(g.context.cHandle, g.context.cWKTWriter, g.cGeom)
	defer C.GEOSFree_r(g.context.cHandle, unsafe.Pointer(cWKTStr))
	return C.GoString(cWKTStr)
}

// Type returns g's type.
func (g *Geom) Type() string {
	g.mustNotBeDestroyed()
	g.context.mutex.Lock()
	defer g.context.mutex.Unlock()
	cTypeStr := C.GEOSGeomType_r(g.context.cHandle, g.cGeom)
	if cTypeStr == nil {
		panic(g.context.err)
	}
	defer C.GEOSFree_r(g.context.cHandle, unsafe.Pointer(cTypeStr))
	return C.GoString(cTypeStr)
}

// TypeID returns g's geometry type id.
func (g *Geom) TypeID() TypeID {
	g.mustNotBeDestroyed()
	return g.typeID
}

// UserData returns g's userdata.
func (g *Geom) UserData() uintptr {
	g.mustNotBeDestroyed()
	g.context.mutex.Lock()
	defer g.context.mutex.Unlock()
	return uintptr(C.c_GEOSGeom_getUserData_r(g.context.cHandle, g.cGeom))
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
