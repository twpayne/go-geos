//go:generate go run ./internal/cmds/execute-template -data geommethods.yaml -output geommethods.go geommethods.go.tmpl

package geos

// #include "go-geos.h"
import "C"

import (
	"runtime"
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

// Clone clones g into c.
func (c *Context) Clone(g *Geom) *Geom {
	if g.context == c {
		return g.Clone()
	}
	// FIXME use a more intelligent method than a WKB roundtrip (although a WKB
	// roundtrip might actually be quite fast if the cgo overhead is
	// significant)
	clone, err := c.NewGeomFromWKB(g.ToWKB())
	if err != nil {
		panic(err)
	}
	return clone
}

// NewCollection returns a new collection.
func (c *Context) NewCollection(typeID TypeID, geoms []*Geom) *Geom {
	if len(geoms) == 0 {
		return c.NewEmptyCollection(typeID)
	}
	c.mutex.Lock()
	defer c.mutex.Unlock()
	cGeoms := make([]*C.GEOSGeometry, len(geoms))
	for i, geom := range geoms {
		cGeoms[i] = geom.cGeom
	}
	g := c.newNonNilGeom(C.GEOSGeom_createCollection_r(c.cHandle, C.int(typeID), &cGeoms[0], C.uint(len(geoms))), nil)
	for _, geom := range geoms {
		geom.parent = g
	}
	return g
}

// NewEmptyCollection returns a new empty collection.
func (c *Context) NewEmptyCollection(typeID TypeID) *Geom {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.newNonNilGeom(C.GEOSGeom_createEmptyCollection_r(c.cHandle, C.int(typeID)), nil)
}

// NewEmptyLineString returns a new empty line string.
func (c *Context) NewEmptyLineString() *Geom {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.newNonNilGeom(C.GEOSGeom_createEmptyLineString_r(c.cHandle), nil)
}

// NewEmptyPoint returns a new empty point.
func (c *Context) NewEmptyPoint() *Geom {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.newNonNilGeom(C.GEOSGeom_createEmptyPoint_r(c.cHandle), nil)
}

// NewEmptyPolygon returns a new empty polygon.
func (c *Context) NewEmptyPolygon() *Geom {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.newNonNilGeom(C.GEOSGeom_createEmptyPolygon_r(c.cHandle), nil)
}

// NewGeomFromBounds returns a new polygon constructed from bounds.
func (c *Context) NewGeomFromBounds(minX, minY, maxX, maxY float64) *Geom {
	var typeID C.int
	geom := C.c_newGEOSGeomFromBounds_r(c.cHandle, &typeID, C.double(minX), C.double(minY), C.double(maxX), C.double(maxY))
	if geom == nil {
		panic(c.err)
	}
	g := &Geom{
		context:       c,
		cGeom:         geom,
		typeID:        TypeID(typeID),
		numGeometries: 1,
	}
	runtime.AddCleanup(g, func(cGeom *C.struct_GEOSGeom_t) {
		C.GEOSGeom_destroy_r(c.cHandle, cGeom)
	}, g.cGeom)
	return g
}

// NewLinearRing returns a new linear ring populated with coords.
func (c *Context) NewLinearRing(coords [][]float64) *Geom {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	s := c.newGEOSCoordSeqFromCoords(coords)
	return c.newNonNilGeom(C.GEOSGeom_createLinearRing_r(c.cHandle, s), nil)
}

// NewLineString returns a new line string populated with coords.
func (c *Context) NewLineString(coords [][]float64) *Geom {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	s := c.newGEOSCoordSeqFromCoords(coords)
	return c.newNonNilGeom(C.GEOSGeom_createLineString_r(c.cHandle, s), nil)
}

// NewPoint returns a new point populated with coord.
func (c *Context) NewPoint(coord []float64) *Geom {
	s := c.newGEOSCoordSeqFromCoords([][]float64{coord})
	return c.newNonNilGeom(C.GEOSGeom_createPoint_r(c.cHandle, s), nil)
}

// NewPointFromXY returns a new point with a x and y.
func (c *Context) NewPointFromXY(x, y float64) *Geom {
	return c.newNonNilGeom(C.GEOSGeom_createPointFromXY_r(c.cHandle, C.double(x), C.double(y)), nil)
}

// NewPoints returns a new slice of points populated from coords.
func (c *Context) NewPoints(coords [][]float64) []*Geom {
	if coords == nil {
		return nil
	}
	geoms := make([]*Geom, len(coords))
	for i := range geoms {
		geoms[i] = c.NewPoint(coords[i])
	}
	return geoms
}

// NewPolygon returns a new point populated with coordss.
func (c *Context) NewPolygon(coordss [][][]float64) *Geom {
	if len(coordss) == 0 {
		return c.NewEmptyPolygon()
	}
	var (
		shell      *C.struct_GEOSGeom_t
		holesSlice []*C.struct_GEOSGeom_t
	)
	defer func() {
		if v := recover(); v != nil {
			C.GEOSGeom_destroy_r(c.cHandle, shell)
			for _, hole := range holesSlice {
				C.GEOSGeom_destroy_r(c.cHandle, hole)
			}
			panic(v)
		}
	}()
	shell = C.GEOSGeom_createLinearRing_r(c.cHandle, c.newGEOSCoordSeqFromCoords(coordss[0]))
	if shell == nil {
		panic(c.err)
	}
	var holes **C.struct_GEOSGeom_t
	nholes := len(coordss) - 1
	if nholes > 0 {
		holesSlice = make([]*C.struct_GEOSGeom_t, nholes)
		for i := range holesSlice {
			hole := C.GEOSGeom_createLinearRing_r(c.cHandle, c.newGEOSCoordSeqFromCoords(coordss[i+1]))
			if hole == nil {
				panic(c.err)
			}
			holesSlice[i] = hole
		}
		holes = (**C.struct_GEOSGeom_t)(unsafe.Pointer(&holesSlice[0]))
	}
	return c.newNonNilGeom(C.GEOSGeom_createPolygon_r(c.cHandle, shell, holes, C.uint(nholes)), nil)
}

// Bounds returns g's bounds.
func (g *Geom) Bounds() *Box2D {
	bounds := NewBox2DEmpty()
	g.context.mutex.Lock()
	defer g.context.mutex.Unlock()
	C.c_GEOSGeomBounds_r(g.context.cHandle, g.cGeom, (*C.double)(&bounds.MinX), (*C.double)(&bounds.MinY), (*C.double)(&bounds.MaxX), (*C.double)(&bounds.MaxY))
	return bounds
}

// MakeValidWithParams returns a new valid geometry using the MakeValidMethods and MakeValidCollapsed parameters.
func (g *Geom) MakeValidWithParams(method MakeValidMethod, collapse MakeValidCollapsed) *Geom {
	g.context.mutex.Lock()
	defer g.context.mutex.Unlock()
	cRes := C.c_GEOSMakeValidWithParams_r(g.context.cHandle, g.cGeom, (C.enum_GEOSMakeValidMethods)(method), C.int(collapse))
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
	g.context.mutex.Lock()
	defer g.context.mutex.Unlock()
	s := C.GEOSGeom_getCoordSeq_r(g.context.cHandle, g.cGeom)
	// Don't set a finalizer as coordSeq is owned by g and will be finalized when g is
	// finalized.
	coordSeq := g.context.newCoordSeq(s, false)
	if coordSeq == nil {
		return nil
	}
	coordSeq.parent = g
	return coordSeq
}

// ExteriorRing returns the exterior ring.
func (g *Geom) ExteriorRing() *Geom {
	g.context.mutex.Lock()
	defer g.context.mutex.Unlock()
	return g.context.newNonNilGeom(C.GEOSGetExteriorRing_r(g.context.cHandle, g.cGeom), g)
}

// Geometry returns the nth geometry of g.
func (g *Geom) Geometry(n int) *Geom {
	g.context.mutex.Lock()
	defer g.context.mutex.Unlock()
	if n < 0 || g.numGeometries <= n {
		panic(errIndexOutOfRange)
	}
	return g.context.newNonNilGeom(C.GEOSGetGeometryN_r(g.context.cHandle, g.cGeom, C.int(n)), g)
}

// InteriorRing returns the nth interior ring.
func (g *Geom) InteriorRing(n int) *Geom {
	g.context.mutex.Lock()
	defer g.context.mutex.Unlock()
	if n < 0 || g.numInteriorRings <= n {
		panic(errIndexOutOfRange)
	}
	return g.context.newNonNilGeom(C.GEOSGetInteriorRingN_r(g.context.cHandle, g.cGeom, C.int(n)), g)
}

// IsValidReason returns the reason that g is invalid.
func (g *Geom) IsValidReason() string {
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
	g.context.mutex.Lock()
	defer g.context.mutex.Unlock()
	if C.GEOSNormalize_r(g.context.cHandle, g.cGeom) != 0 {
		panic(g.context.err)
	}
	return g
}

// NumCoordinates returns the number of coordinates in g.
func (g *Geom) NumCoordinates() int {
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
	return g.numGeometries
}

// NumInteriorRings returns the number of interior rings in g.
func (g *Geom) NumInteriorRings() int {
	return g.numInteriorRings
}

// NumPoints returns the number of points in g.
func (g *Geom) NumPoints() int {
	return g.numPoints
}

// Point returns the g's nth point.
func (g *Geom) Point(n int) *Geom {
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
	g.context.mutex.Lock()
	defer g.context.mutex.Unlock()
	var cCuts, cDangles, cInvalidRings *C.struct_GEOSGeom_t
	cGeom := C.GEOSPolygonize_full_r(g.context.cHandle, g.cGeom, &cCuts, &cDangles, &cInvalidRings) //nolint:gocritic
	geom = g.context.newNonNilGeom(cGeom, nil)
	cuts = g.context.newGeom(cCuts, nil)
	dangles = g.context.newGeom(cDangles, nil)
	invalidRings = g.context.newGeom(cInvalidRings, nil)
	return
}

// Precision returns g's precision.
func (g *Geom) Precision() float64 {
	g.context.mutex.Lock()
	defer g.context.mutex.Unlock()
	return float64(C.GEOSGeom_getPrecision_r(g.context.cHandle, g.cGeom))
}

// RelatePattern returns if the DE9IM pattern for g and other matches pat.
func (g *Geom) RelatePattern(other *Geom, pat string) bool {
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
	g.context.mutex.Lock()
	defer g.context.mutex.Unlock()
	C.GEOSSetSRID_r(g.context.cHandle, g.cGeom, C.int(srid))
	return g
}

// SetUserData sets g's userdata and returns g.
func (g *Geom) SetUserData(userdata uintptr) *Geom {
	g.context.mutex.Lock()
	defer g.context.mutex.Unlock()
	C.c_GEOSGeom_setUserData_r(g.context.cHandle, g.cGeom, C.uintptr_t(userdata))
	return g
}

// String returns g in WKT format.
func (g *Geom) String() string {
	return g.ToWKT()
}

// ToGeoJSON returns g in GeoJSON format.
func (g *Geom) ToGeoJSON(indent int) string {
	g.context.mutex.Lock()
	defer g.context.mutex.Unlock()
	if g.context.geoJSONWriter == nil {
		g.context.geoJSONWriter = g.context.NewGeoJSONWriter()
	}
	return g.context.geoJSONWriter.WriteGeometry(g, indent)
}

// ToWKB returns g in WKB format.
func (g *Geom) ToWKB() []byte {
	g.context.mutex.Lock()
	defer g.context.mutex.Unlock()
	if g.context.wkbWriter == nil {
		g.context.wkbWriter = g.context.NewWKBWriter(
			WithFlavor(WKBFlavorExtended),
			WithIncludeSRID(true),
		)
	}
	return g.context.wkbWriter.Write(g)
}

// ToWKT returns g in WKT format.
func (g *Geom) ToWKT() string {
	g.context.mutex.Lock()
	defer g.context.mutex.Unlock()
	if g.context.wktWriter == nil {
		g.context.wktWriter = g.context.NewWKTWriter()
	}
	return g.context.wktWriter.Write(g)
}

// Type returns g's type.
func (g *Geom) Type() string {
	g.context.mutex.Lock()
	defer g.context.mutex.Unlock()
	typeCStr := C.GEOSGeomType_r(g.context.cHandle, g.cGeom)
	if typeCStr == nil {
		panic(g.context.err)
	}
	defer C.GEOSFree_r(g.context.cHandle, unsafe.Pointer(typeCStr))
	return C.GoString(typeCStr)
}

// TypeID returns g's geometry type id.
func (g *Geom) TypeID() TypeID {
	return g.typeID
}

// UserData returns g's userdata.
func (g *Geom) UserData() uintptr {
	g.context.mutex.Lock()
	defer g.context.mutex.Unlock()
	return uintptr(C.c_GEOSGeom_getUserData_r(g.context.cHandle, g.cGeom))
}
