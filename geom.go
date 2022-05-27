package geos

// #include "geos.h"
import "C"

import (
	"runtime"
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

// Clone returns a clone of g.
func (g *Geom) Clone() *Geom {
	g.mustNotBeDestroyed()
	g.context.Lock()
	defer g.context.Unlock()
	return g.context.newNonNilGeom(C.GEOSGeom_clone_r(g.context.handle, g.geom), nil)
}

// ConvexHull returns g's convex hull.
func (g *Geom) ConvexHull() *Geom {
	g.mustNotBeDestroyed()
	g.context.Lock()
	defer g.context.Unlock()
	return g.context.newNonNilGeom(C.GEOSConvexHull_r(g.context.handle, g.geom), nil)
}

// Contains returns true if g contains other.
func (g *Geom) Contains(other *Geom) bool {
	g.mustNotBeDestroyed()
	g.context.Lock()
	defer g.context.Unlock()
	if other.context != g.context {
		other.context.Lock()
		defer other.context.Unlock()
	}
	switch C.GEOSContains_r(g.context.handle, g.geom, other.geom) {
	case 0:
		return false
	case 1:
		return true
	default:
		panic(g.context.err)
	}
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

// CoveredBy returns true if g is covered by other.
func (g *Geom) CoveredBy(other *Geom) bool {
	g.mustNotBeDestroyed()
	g.context.Lock()
	defer g.context.Unlock()
	if other.context != g.context {
		other.context.Lock()
		defer other.context.Unlock()
	}
	switch C.GEOSCoveredBy_r(g.context.handle, g.geom, other.geom) {
	case 0:
		return false
	case 1:
		return true
	default:
		panic(g.context.err)
	}
}

// Covers returns true if g covers other.
func (g *Geom) Covers(other *Geom) bool {
	g.mustNotBeDestroyed()
	g.context.Lock()
	defer g.context.Unlock()
	if other.context != g.context {
		other.context.Lock()
		defer other.context.Unlock()
	}
	switch C.GEOSCovers_r(g.context.handle, g.geom, other.geom) {
	case 0:
		return false
	case 1:
		return true
	default:
		panic(g.context.err)
	}
}

// Crosses returns true if g crosses other.
func (g *Geom) Crosses(other *Geom) bool {
	g.mustNotBeDestroyed()
	g.context.Lock()
	defer g.context.Unlock()
	if other.context != g.context {
		other.context.Lock()
		defer other.context.Unlock()
	}
	switch C.GEOSCrosses_r(g.context.handle, g.geom, other.geom) {
	case 0:
		return false
	case 1:
		return true
	default:
		panic(g.context.err)
	}
}

// Disjoint returns true if g is disjoint from other.
func (g *Geom) Disjoint(other *Geom) bool {
	g.mustNotBeDestroyed()
	g.context.Lock()
	defer g.context.Unlock()
	if other.context != g.context {
		other.context.Lock()
		defer other.context.Unlock()
	}
	switch C.GEOSDisjoint_r(g.context.handle, g.geom, other.geom) {
	case 0:
		return false
	case 1:
		return true
	default:
		panic(g.context.err)
	}
}

// Distance returns the distance between the closest points on g and other.
func (g *Geom) Distance(other *Geom) float64 {
	g.mustNotBeDestroyed()
	g.context.Lock()
	defer g.context.Unlock()
	if other.context != g.context {
		other.context.Lock()
		defer other.context.Unlock()
	}
	var distance float64
	if C.GEOSDistance_r(g.context.handle, g.geom, other.geom, (*C.double)(&distance)) == 0 {
		panic(g.context.err)
	}
	return distance
}

// Envelope returns the envelope of g.
func (g *Geom) Envelope() *Geom {
	g.mustNotBeDestroyed()
	g.context.Lock()
	defer g.context.Unlock()
	return g.context.newNonNilGeom(C.GEOSEnvelope_r(g.context.handle, g.geom), nil)
}

// Equals returns true if g equals other.
func (g *Geom) Equals(other *Geom) bool {
	g.mustNotBeDestroyed()
	g.context.Lock()
	defer g.context.Unlock()
	if other.context != g.context {
		other.context.Lock()
		defer other.context.Unlock()
	}
	switch C.GEOSEquals_r(g.context.handle, g.geom, other.geom) {
	case 0:
		return false
	case 1:
		return true
	default:
		panic(g.context.err)
	}
}

// EqualsExact returns true if g equals other exactly.
func (g *Geom) EqualsExact(other *Geom, tolerance float64) bool {
	g.mustNotBeDestroyed()
	g.context.Lock()
	defer g.context.Unlock()
	if other.context != g.context {
		other.context.Lock()
		defer other.context.Unlock()
	}
	switch C.GEOSEqualsExact_r(g.context.handle, g.geom, other.geom, C.double(tolerance)) {
	case 0:
		return false
	case 1:
		return true
	default:
		panic(g.context.err)
	}
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

// Buffer() returns a buffer polygon with width w.
func (g *Geom) Buffer(width float64, quadsegs int) *Geom {
	g.mustNotBeDestroyed()
	g.context.Lock()
	defer g.context.Unlock()

	return g.context.newNonNilGeom(C.GEOSBuffer_r(g.context.handle, g.geom, C.double(width), C.int(quadsegs)), g)
}

// UnaryUnion() returns the union of all components of a single geometry. Usually used to convert a collection into the smallest set of polygons that cover the same area.
func (g *Geom) UnaryUnion() *Geom {
	g.mustNotBeDestroyed()
	g.context.Lock()
	defer g.context.Unlock()

	return g.context.newNonNilGeom(C.GEOSUnaryUnion_r(g.context.handle, g.geom), g)
}

// Densifies a geometry using a given distance tolerance.
// Additional vertices will be added to every line segment that is greater this tolerance; these vertices will evenly subdivide that segment.
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

// Intersects returns true if g intersects other.
func (g *Geom) Intersects(other *Geom) bool {
	g.mustNotBeDestroyed()
	g.context.Lock()
	defer g.context.Unlock()
	if other.context != g.context {
		other.context.Lock()
		defer other.context.Unlock()
	}
	switch C.GEOSIntersects_r(g.context.handle, g.geom, other.geom) {
	case 0:
		return false
	case 1:
		return true
	default:
		panic(g.context.err)
	}
}

// IsClosed returns true if g is closed.
func (g *Geom) IsClosed() bool {
	g.mustNotBeDestroyed()
	g.context.Lock()
	defer g.context.Unlock()
	switch C.GEOSisClosed_r(g.context.handle, g.geom) {
	case 0:
		return false
	case 1:
		return true
	default:
		panic(g.context.err)
	}
}

// IsEmpty returns true if g is empty.
func (g *Geom) IsEmpty() bool {
	g.mustNotBeDestroyed()
	g.context.Lock()
	defer g.context.Unlock()
	switch C.GEOSisEmpty_r(g.context.handle, g.geom) {
	case 0:
		return false
	case 1:
		return true
	default:
		panic(g.context.err)
	}
}

// IsRing returns true if g is a ring.
func (g *Geom) IsRing() bool {
	g.mustNotBeDestroyed()
	g.context.Lock()
	defer g.context.Unlock()
	switch C.GEOSisRing_r(g.context.handle, g.geom) {
	case 0:
		return false
	case 1:
		return true
	default:
		panic(g.context.err)
	}
}

// IsSimple returns true if g is simple.
func (g *Geom) IsSimple() bool {
	g.mustNotBeDestroyed()
	g.context.Lock()
	defer g.context.Unlock()
	switch C.GEOSisSimple_r(g.context.handle, g.geom) {
	case 0:
		return false
	case 1:
		return true
	default:
		panic(g.context.err)
	}
}

// IsValid returns true if g is valid.
func (g *Geom) IsValid() bool {
	g.mustNotBeDestroyed()
	g.context.Lock()
	defer g.context.Unlock()
	switch C.GEOSisValid_r(g.context.handle, g.geom) {
	case 0:
		return false
	case 1:
		return true
	default:
		panic(g.context.err)
	}
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

// Overlaps returns true if g overlaps other.
func (g *Geom) Overlaps(other *Geom) bool {
	g.mustNotBeDestroyed()
	g.context.Lock()
	defer g.context.Unlock()
	if other.context != g.context {
		other.context.Lock()
		defer other.context.Unlock()
	}
	switch C.GEOSOverlaps_r(g.context.handle, g.geom, other.geom) {
	case 0:
		return false
	case 1:
		return true
	default:
		panic(g.context.err)
	}
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

// Touches returns true if g touches other.
func (g *Geom) Touches(other *Geom) bool {
	g.mustNotBeDestroyed()
	g.context.Lock()
	defer g.context.Unlock()
	if other.context != g.context {
		other.context.Lock()
		defer other.context.Unlock()
	}
	switch C.GEOSTouches_r(g.context.handle, g.geom, other.geom) {
	case 0:
		return false
	case 1:
		return true
	default:
		panic(g.context.err)
	}
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

// Within returns if g is within other.
func (g *Geom) Within(other *Geom) bool {
	g.mustNotBeDestroyed()
	g.context.Lock()
	defer g.context.Unlock()
	if other.context != g.context {
		other.context.Lock()
		defer other.context.Unlock()
	}
	switch C.GEOSWithin_r(g.context.handle, g.geom, other.geom) {
	case 0:
		return false
	case 1:
		return true
	default:
		panic(g.context.err)
	}
}

// X returns g's X coordinate.
func (g *Geom) X() float64 {
	g.mustNotBeDestroyed()
	g.context.Lock()
	defer g.context.Unlock()
	var value float64
	if C.GEOSGeomGetX_r(g.context.handle, g.geom, (*C.double)(&value)) == -1 {
		panic(g.context.err)
	}
	return value
}

// Y returns g's Y coordinate.
func (g *Geom) Y() float64 {
	g.mustNotBeDestroyed()
	g.context.Lock()
	defer g.context.Unlock()
	var value float64
	if C.GEOSGeomGetY_r(g.context.handle, g.geom, (*C.double)(&value)) == -1 {
		panic(g.context.err)
	}
	return value
}

// Calculate the area of a geometry.
func (g *Geom) Area() float64 {
	g.mustNotBeDestroyed()
	g.context.Lock()
	defer g.context.Unlock()
	var area float64
	if C.GEOSArea_r(g.context.handle, g.geom, (*C.double)(&area)) == 0 {
		panic(g.context.err)
	}
	return area
}

// Calculate the length of a geometry.
func (g *Geom) Length() float64 {
	g.mustNotBeDestroyed()
	g.context.Lock()
	defer g.context.Unlock()
	var area float64
	if C.GEOSLength_r(g.context.handle, g.geom, (*C.double)(&area)) == 0 {
		panic(g.context.err)
	}
	return area
}

// Returns a linestring geometry which represents the minimum diameter of the geometry.
func (g *Geom) MinimumWidth() *Geom {
	g.mustNotBeDestroyed()
	g.context.Lock()
	defer g.context.Unlock()

	return g.context.newNonNilGeom(C.GEOSMinimumWidth_r(g.context.handle, g.geom), g)
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

// Implements UnmarshalJSON interface
func (g *Geom) UnmarshalJSON(b []byte) error {
	c := defaultContext
	c.Lock()
	defer c.Unlock()
	c.err = nil
	if c.jsonReader == nil {
		c.jsonReader = C.GEOSGeoJSONReader_create_r(c.handle)
	}
	jsonCStr := C.CString(string(b))
	defer C.GEOSFree_r(c.handle, unsafe.Pointer(jsonCStr))
	newgeom := C.GEOSGeoJSONReader_readGeometry_r(c.handle, c.jsonReader, jsonCStr)
	var (
		typeID           C.int
		numGeometries    C.int
		numPoints        C.int
		numInteriorRings C.int
	)
	if C.c_GEOSGeomGetInfo_r(c.handle, newgeom, &typeID, &numGeometries, &numPoints, &numInteriorRings) == 0 {
		panic(c.err)
	}

	g.context = c
	g.geom = newgeom
	g.parent = nil
	g.typeID = GeometryTypeID(typeID)
	g.numGeometries = int(numGeometries)
	g.numInteriorRings = int(numInteriorRings)
	g.numPoints = int(numPoints)

	runtime.SetFinalizer(g, (*Geom).finalize)
	return c.err
}
