package geos

// #include <stdlib.h>
// #include "go-geos.h"
import "C"

import (
	"runtime"
	"sync"
	"unsafe"
)

// A Context is a context.
type Context struct {
	sync.Mutex
	handle              C.GEOSContextHandle_t
	geoJSONReader       *C.struct_GEOSGeoJSONReader_t
	geoJSONWriter       *C.struct_GEOSGeoJSONWriter_t
	wkbReader           *C.struct_GEOSWKBReader_t
	wkbWriter           *C.struct_GEOSWKBWriter_t
	wktReader           *C.struct_GEOSWKTReader_t
	wktWriter           *C.struct_GEOSWKTWriter_t
	err                 error
	geomFinalizeFunc    func(*Geom)
	strTreeFinalizeFunc func(*STRtree)
}

// A ContextOption sets an option on a Context.
type ContextOption func(*Context)

// WithGeomFinalizeFunc sets a function to be called just before a geometry is
// finalized. This is typically used to log the geometry to help debug geometry
// leaks.
func WithGeomFinalizeFunc(geomFinalizeFunc func(*Geom)) ContextOption {
	return func(c *Context) {
		c.geomFinalizeFunc = geomFinalizeFunc
	}
}

// WithSTRtreeFinalizeFunc sets a function to be called just before an STRtree
// is finalized. This is typically used to log the STRtree to help debug STRtree
// leaks.
func WithSTRtreeFinalizeFunc(strTreeFinalizeFunc func(*STRtree)) ContextOption {
	return func(c *Context) {
		c.strTreeFinalizeFunc = strTreeFinalizeFunc
	}
}

// NewContext returns a new Context.
func NewContext(options ...ContextOption) *Context {
	c := &Context{
		handle: C.GEOS_init_r(),
	}
	runtime.SetFinalizer(c, (*Context).finish)
	// FIXME in GitHub Actions, golangci-lint complains about the following line saying:
	// Error: dupSubExpr: suspicious identical LHS and RHS for `==` operator (gocritic)
	// As the line does not contain an `==` operator, disable gocritic on this line.
	//nolint:gocritic
	C.GEOSContext_setErrorMessageHandler_r(c.handle, C.GEOSMessageHandler_r(C.c_errorMessageHandler), unsafe.Pointer(&c.err))
	for _, option := range options {
		option(c)
	}
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
	c.Lock()
	defer c.Unlock()
	cGeoms := make([]*C.GEOSGeometry, len(geoms))
	for i, geom := range geoms {
		cGeoms[i] = geom.geom
	}
	g := c.newNonNilGeom(C.GEOSGeom_createCollection_r(c.handle, C.int(typeID), &cGeoms[0], C.uint(len(geoms))), nil)
	for _, geom := range geoms {
		geom.parent = g
	}
	return g
}

// NewCoordSeq returns a new CoordSeq.
func (c *Context) NewCoordSeq(size, dims int) *CoordSeq {
	c.Lock()
	defer c.Unlock()
	return c.newNonNilCoordSeq(C.GEOSCoordSeq_create_r(c.handle, C.uint(size), C.uint(dims)))
}

// NewCoordSeqFromCoords returns a new CoordSeq populated with coords.
func (c *Context) NewCoordSeqFromCoords(coords [][]float64) *CoordSeq {
	c.Lock()
	defer c.Unlock()
	return c.newNonNilCoordSeq(c.newGEOSCoordSeqFromCoords(coords))
}

// NewEmptyCollection returns a new empty collection.
func (c *Context) NewEmptyCollection(typeID TypeID) *Geom {
	c.Lock()
	defer c.Unlock()
	return c.newNonNilGeom(C.GEOSGeom_createEmptyCollection_r(c.handle, C.int(typeID)), nil)
}

// NewEmptyLineString returns a new empty line string.
func (c *Context) NewEmptyLineString() *Geom {
	c.Lock()
	defer c.Unlock()
	return c.newNonNilGeom(C.GEOSGeom_createEmptyLineString_r(c.handle), nil)
}

// NewEmptyPoint returns a new empty point.
func (c *Context) NewEmptyPoint() *Geom {
	c.Lock()
	defer c.Unlock()
	return c.newNonNilGeom(C.GEOSGeom_createEmptyPoint_r(c.handle), nil)
}

// NewEmptyPolygon returns a new empty polygon.
func (c *Context) NewEmptyPolygon() *Geom {
	c.Lock()
	defer c.Unlock()
	return c.newNonNilGeom(C.GEOSGeom_createEmptyPolygon_r(c.handle), nil)
}

// NewGeomFromBounds returns a new polygon constructed from bounds.
func (c *Context) NewGeomFromBounds(bounds *Bounds) *Geom {
	var typeID C.int
	geom := C.c_newGEOSGeomFromBounds_r(c.handle, &typeID, C.double(bounds.MinX), C.double(bounds.MinY), C.double(bounds.MaxX), C.double(bounds.MaxY))
	if geom == nil {
		panic(c.err)
	}
	g := &Geom{
		context:       c,
		geom:          geom,
		typeID:        TypeID(typeID),
		numGeometries: 1,
	}
	runtime.SetFinalizer(g, (*Geom).finalize)
	return g
}

// NewGeomFromGeoJSON returns a new geometry in JSON format from json.
func (c *Context) NewGeomFromGeoJSON(geoJSON string) (*Geom, error) {
	c.Lock()
	defer c.Unlock()
	c.err = nil
	if c.geoJSONReader == nil {
		c.geoJSONReader = C.GEOSGeoJSONReader_create_r(c.handle)
	}
	geoJSONCStr := C.CString(geoJSON)
	defer C.free(unsafe.Pointer(geoJSONCStr))
	return c.newGeom(C.GEOSGeoJSONReader_readGeometry_r(c.handle, c.geoJSONReader, geoJSONCStr), nil), c.err
}

// NewGeomFromWKB parses a geometry in WKB format from wkb.
func (c *Context) NewGeomFromWKB(wkb []byte) (*Geom, error) {
	c.Lock()
	defer c.Unlock()
	c.err = nil
	if c.wkbReader == nil {
		c.wkbReader = C.GEOSWKBReader_create_r(c.handle)
	}
	wkbCBuf := C.CBytes(wkb)
	defer C.free(wkbCBuf)
	return c.newGeom(C.GEOSWKBReader_read_r(c.handle, c.wkbReader, (*C.uchar)(wkbCBuf), C.ulong(len(wkb))), nil), c.err
}

// NewGeomFromWKT parses a geometry in WKT format from wkt.
func (c *Context) NewGeomFromWKT(wkt string) (*Geom, error) {
	c.Lock()
	defer c.Unlock()
	c.err = nil
	if c.wktReader == nil {
		c.wktReader = C.GEOSWKTReader_create_r(c.handle)
	}
	wktCStr := C.CString(wkt)
	defer C.free(unsafe.Pointer(wktCStr))
	return c.newGeom(C.GEOSWKTReader_read_r(c.handle, c.wktReader, wktCStr), nil), c.err
}

// NewLinearRing returns a new linear ring populated with coords.
func (c *Context) NewLinearRing(coords [][]float64) *Geom {
	c.Lock()
	defer c.Unlock()
	s := c.newGEOSCoordSeqFromCoords(coords)
	return c.newNonNilGeom(C.GEOSGeom_createLinearRing_r(c.handle, s), nil)
}

// NewLineString returns a new line string populated with coords.
func (c *Context) NewLineString(coords [][]float64) *Geom {
	c.Lock()
	defer c.Unlock()
	s := c.newGEOSCoordSeqFromCoords(coords)
	return c.newNonNilGeom(C.GEOSGeom_createLineString_r(c.handle, s), nil)
}

// NewPoint returns a new point populated with coord.
func (c *Context) NewPoint(coord []float64) *Geom {
	s := c.newGEOSCoordSeqFromCoords([][]float64{coord})
	return c.newNonNilGeom(C.GEOSGeom_createPoint_r(c.handle, s), nil)
}

// NewPointFromXY returns a new point with a x and y.
func (c *Context) NewPointFromXY(x, y float64) *Geom {
	return c.newNonNilGeom(C.GEOSGeom_createPointFromXY_r(c.handle, C.double(x), C.double(y)), nil)
}

// NewPoints returns a new slice of points populated from coords.
func (c *Context) NewPoints(coords [][]float64) []*Geom {
	if coords == nil {
		return nil
	}
	geoms := make([]*Geom, 0, len(coords))
	for _, coord := range coords {
		geom := c.NewPoint(coord)
		geoms = append(geoms, geom)
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
			C.GEOSGeom_destroy_r(c.handle, shell)
			for _, hole := range holesSlice {
				C.GEOSGeom_destroy_r(c.handle, hole)
			}
			panic(v)
		}
	}()
	shell = C.GEOSGeom_createLinearRing_r(c.handle, c.newGEOSCoordSeqFromCoords(coordss[0]))
	if shell == nil {
		panic(c.err)
	}
	var holes **C.struct_GEOSGeom_t
	nholes := len(coordss) - 1
	if nholes > 0 {
		holesSlice = make([]*C.struct_GEOSGeom_t, 0, nholes)
		for i := 0; i < nholes; i++ {
			hole := C.GEOSGeom_createLinearRing_r(c.handle, c.newGEOSCoordSeqFromCoords(coordss[i+1]))
			if hole == nil {
				panic(c.err)
			}
			holesSlice = append(holesSlice, hole)
		}
		holes = (**C.struct_GEOSGeom_t)(unsafe.Pointer(&holesSlice[0]))
	}
	return c.newNonNilGeom(C.GEOSGeom_createPolygon_r(c.handle, shell, holes, C.uint(nholes)), nil)
}

// NewSTRtree returns a new STRtree.
func (c *Context) NewSTRtree(nodeCapacity int) *STRtree {
	c.Lock()
	defer c.Unlock()
	t := &STRtree{
		context:     c,
		strTree:     C.GEOSSTRtree_create_r(c.handle, C.size_t(nodeCapacity)),
		itemToValue: make(map[unsafe.Pointer]any),
		valueToItem: make(map[any]unsafe.Pointer),
	}
	runtime.SetFinalizer(t, (*STRtree).finalize)
	return t
}

// OrientationIndex returns the orientation index from A to B and then to P.
func (c *Context) OrientationIndex(Ax, Ay, Bx, By, Px, Py float64) int { //nolint:gocritic
	c.Lock()
	defer c.Unlock()
	return int(C.GEOSOrientationIndex_r(c.handle, C.double(Ax), C.double(Ay), C.double(Bx), C.double(By), C.double(Px), C.double(Py)))
}

// Polygonize returns a set of geometries which contains linework that
// represents the edges of a planar graph.
func (c *Context) Polygonize(geoms []*Geom) *Geom {
	c.Lock()
	defer c.Unlock()
	cGeoms, unlockFunc := c.cGeomsLocked(geoms)
	defer unlockFunc()
	return c.newNonNilGeom(C.GEOSPolygonize_r(c.handle, cGeoms, C.uint(len(geoms))), nil)
}

// PolygonizeValid returns a set of polygons which contains linework that
// represents the edges of a planar graph.
func (c *Context) PolygonizeValid(geoms []*Geom) *Geom {
	c.Lock()
	defer c.Unlock()
	cGeoms, unlockFunc := c.cGeomsLocked(geoms)
	defer unlockFunc()
	return c.newNonNilGeom(C.GEOSPolygonize_valid_r(c.handle, cGeoms, C.uint(len(geoms))), nil)
}

// RelatePatternMatch returns if two DE9IM patterns are consistent.
func (c *Context) RelatePatternMatch(mat, pat string) bool {
	matCStr := C.CString(mat)
	defer C.free(unsafe.Pointer(matCStr))
	patCStr := C.CString(pat)
	defer C.free(unsafe.Pointer(patCStr))
	c.Lock()
	defer c.Unlock()
	switch C.GEOSRelatePatternMatch_r(c.handle, matCStr, patCStr) {
	case 0:
		return false
	case 1:
		return true
	default:
		panic(c.err)
	}
}

// SegmentIntersection returns the coordinate where two lines intersect.
func (c *Context) SegmentIntersection(ax0, ay0, ax1, ay1, bx0, by0, bx1, by1 float64) (float64, float64, bool) {
	c.Lock()
	defer c.Unlock()
	var cx, cy float64
	switch C.GEOSSegmentIntersection_r(c.handle,
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
	cGeoms := make([]*C.struct_GEOSGeom_t, 0, len(geoms))
	for _, geom := range geoms {
		geom.mustNotBeDestroyed()
		if _, ok := uniqueContexts[geom.context]; !ok {
			geom.context.Lock()
			uniqueContexts[geom.context] = struct{}{}
			extraContexts = append(extraContexts, geom.context)
		}
		cGeoms = append(cGeoms, geom.geom)
	}
	return &cGeoms[0], func() {
		for i := len(extraContexts) - 1; i >= 0; i-- {
			extraContexts[i].Unlock()
		}
	}
}

func (c *Context) finish() {
	c.Lock()
	defer c.Unlock()
	if c.geoJSONReader != nil {
		C.GEOSGeoJSONReader_destroy_r(c.handle, c.geoJSONReader)
	}
	if c.geoJSONWriter != nil {
		C.GEOSGeoJSONWriter_destroy_r(c.handle, c.geoJSONWriter)
	}
	if c.wkbReader != nil {
		C.GEOSWKBReader_destroy_r(c.handle, c.wkbReader)
	}
	if c.wkbWriter != nil {
		C.GEOSWKBWriter_destroy_r(c.handle, c.wkbWriter)
	}
	if c.wktReader != nil {
		C.GEOSWKTReader_destroy_r(c.handle, c.wktReader)
	}
	if c.wktWriter != nil {
		C.GEOSWKTWriter_destroy_r(c.handle, c.wktWriter)
	}
	C.finishGEOS_r(c.handle)
}

func (c *Context) newCoordSeq(gs *C.struct_GEOSCoordSeq_t, finalizer func(*CoordSeq)) *CoordSeq {
	if gs == nil {
		return nil
	}
	var (
		dimensions C.uint
		size       C.uint
	)
	if C.GEOSCoordSeq_getDimensions_r(c.handle, gs, &dimensions) == 0 {
		panic(c.err)
	}
	if C.GEOSCoordSeq_getSize_r(c.handle, gs, &size) == 0 {
		panic(c.err)
	}
	s := &CoordSeq{
		context:    c,
		s:          gs,
		dimensions: int(dimensions),
		size:       int(size),
	}
	if finalizer != nil {
		runtime.SetFinalizer(s, finalizer)
	}
	return s
}

func (c *Context) newCoordsFromGEOSCoordSeq(s *C.struct_GEOSCoordSeq_t) [][]float64 {
	var dimensions C.uint
	if C.GEOSCoordSeq_getDimensions_r(c.handle, s, &dimensions) == 0 {
		panic(c.err)
	}

	var size C.uint
	if C.GEOSCoordSeq_getSize_r(c.handle, s, &size) == 0 {
		panic(c.err)
	}

	var hasZ C.int
	if dimensions > 2 {
		hasZ = 1
	}

	var hasM C.int
	if dimensions > 3 {
		hasM = 1
	}

	flatCoords := make([]float64, size*dimensions)
	if C.GEOSCoordSeq_copyToBuffer_r(c.handle, s, (*C.double)(&flatCoords[0]), hasZ, hasM) == 0 {
		panic(c.err)
	}
	coords := make([][]float64, 0, size)
	for i := 0; i < int(size); i++ {
		coord := flatCoords[i*int(dimensions) : (i+1)*int(dimensions)]
		coords = append(coords, coord)
	}
	return coords
}

func (c *Context) newGEOSCoordSeqFromCoords(coords [][]float64) *C.struct_GEOSCoordSeq_t {
	var hasZ C.int
	if len(coords[0]) > 2 {
		hasZ = 1
	}

	var hasM C.int
	if len(coords[0]) > 3 {
		hasM = 1
	}

	flatCoords := make([]float64, 0, len(coords)*len(coords[0]))
	for _, coord := range coords {
		flatCoords = append(flatCoords, coord...)
	}
	return C.GEOSCoordSeq_copyFromBuffer_r(c.handle, (*C.double)(unsafe.Pointer(&flatCoords[0])), C.uint(len(coords)), hasZ, hasM)
}

func (c *Context) newGeom(geom *C.struct_GEOSGeom_t, parent *Geom) *Geom {
	if geom == nil {
		return nil
	}
	var (
		typeID           C.int
		numGeometries    C.int
		numPoints        C.int
		numInteriorRings C.int
	)
	if C.c_GEOSGeomGetInfo_r(c.handle, geom, &typeID, &numGeometries, &numPoints, &numInteriorRings) == 0 {
		panic(c.err)
	}
	g := &Geom{
		context:          c,
		geom:             geom,
		parent:           parent,
		typeID:           TypeID(typeID),
		numGeometries:    int(numGeometries),
		numInteriorRings: int(numInteriorRings),
		numPoints:        int(numPoints),
	}
	runtime.SetFinalizer(g, (*Geom).finalize)
	return g
}

func (c *Context) newNonNilCoordSeq(s *C.struct_GEOSCoordSeq_t) *CoordSeq {
	if s == nil {
		panic(c.err)
	}
	return c.newCoordSeq(s, (*CoordSeq).destroy)
}

func (c *Context) newNonNilGeom(geom *C.struct_GEOSGeom_t, parent *Geom) *Geom {
	if geom == nil {
		panic(c.err)
	}
	return c.newGeom(geom, parent)
}

//export go_errorMessageHandler
func go_errorMessageHandler(message *C.char, userdata unsafe.Pointer) {
	errP := (*error)(userdata)
	*errP = Error(C.GoString(message))
}
