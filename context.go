package geos

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
	cHandle            C.GEOSContextHandle_t
	ewkbWithSRIDWriter *WKBWriter
	geoJSONReader      *GeoJSONReader
	geoJSONWriter      *GeoJSONWriter
	wkbReader          *WKBReader
	wkbWriter          *WKBWriter
	wktReader          *WKTReader
	wktWriter          *WKTWriter
	err                error
}

// A ContextOption sets an option on a Context.
type ContextOption func(*Context)

// NewContext returns a new Context.
func NewContext(options ...ContextOption) *Context {
	c := &Context{
		cHandle: C.GEOS_init_r(),
	}
	runtime.AddCleanup(c, func(cHandle C.GEOSContextHandle_t) {
		C.finishGEOS_r(cHandle)
	}, c.cHandle)
	// FIXME in GitHub Actions, golangci-lint complains about the following line saying:
	// Error: dupSubExpr: suspicious identical LHS and RHS for `==` operator (gocritic)
	// As the line does not contain an `==` operator, disable gocritic on this line.
	//nolint:gocritic
	C.GEOSContext_setErrorMessageHandler_r(c.cHandle, C.GEOSMessageHandler_r(C.c_errorMessageHandler), unsafe.Pointer(&c.err))
	for _, option := range options {
		option(c)
	}
	return c
}

// NewGeomFromGeoJSON returns a new geometry in JSON format from json.
func (c *Context) NewGeomFromGeoJSON(geoJSON string) (*Geom, error) {
	c.Lock()
	defer c.Unlock()
	c.err = nil
	if c.geoJSONReader == nil {
		c.geoJSONReader = c.NewGeoJSONReader()
	}
	return c.geoJSONReader.ReadGeometry(geoJSON)
}

// NewGeomFromWKB parses a geometry in WKB format from wkb.
func (c *Context) NewGeomFromWKB(wkb []byte) (*Geom, error) {
	c.Lock()
	defer c.Unlock()
	c.err = nil
	if c.wkbReader == nil {
		c.wkbReader = c.NewWKBReader()
	}
	return c.wkbReader.Read(wkb)
}

// NewGeomFromWKT parses a geometry in WKT format from wkt.
func (c *Context) NewGeomFromWKT(wkt string) (*Geom, error) {
	c.Lock()
	defer c.Unlock()
	c.err = nil
	if c.wktReader == nil {
		c.wktReader = c.NewWKTReader()
	}
	return c.wktReader.Read(wkt)
}

// OrientationIndex returns the orientation index from A to B and then to P.
func (c *Context) OrientationIndex(Ax, Ay, Bx, By, Px, Py float64) int { //nolint:gocritic
	c.Lock()
	defer c.Unlock()
	return int(C.GEOSOrientationIndex_r(c.cHandle, C.double(Ax), C.double(Ay), C.double(Bx), C.double(By), C.double(Px), C.double(Py)))
}

// Polygonize returns a set of geometries which contains linework that
// represents the edges of a planar graph.
func (c *Context) Polygonize(geoms []*Geom) *Geom {
	c.Lock()
	defer c.Unlock()
	cGeoms, unlockFunc := c.cGeomsLocked(geoms)
	defer unlockFunc()
	return c.newNonNilGeom(C.GEOSPolygonize_r(c.cHandle, cGeoms, C.uint(len(geoms))), nil)
}

// PolygonizeValid returns a set of polygons which contains linework that
// represents the edges of a planar graph.
func (c *Context) PolygonizeValid(geoms []*Geom) *Geom {
	c.Lock()
	defer c.Unlock()
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
	c.Lock()
	defer c.Unlock()
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
func (c *Context) SegmentIntersection(ax0, ay0, ax1, ay1, bx0, by0, bx1, by1 float64) (float64, float64, bool) {
	c.Lock()
	defer c.Unlock()
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
			geom.context.Lock()
			uniqueContexts[geom.context] = struct{}{}
			extraContexts = append(extraContexts, geom.context)
		}
		cGeoms[i] = geom.cGeom
	}
	return &cGeoms[0], func() {
		for i := len(extraContexts) - 1; i >= 0; i-- {
			extraContexts[i].Unlock()
		}
	}
}

func (c *Context) newBufParams(p *C.struct_GEOSBufParams_t) *BufferParams {
	if p == nil {
		return nil
	}
	b := &BufferParams{
		context:    c,
		cBufParams: p,
	}
	runtime.AddCleanup(b, func(cBufParams *C.struct_GEOSBufParams_t) {
		C.GEOSBufferParams_destroy_r(c.cHandle, cBufParams)
	}, b.cBufParams)
	return b
}

func (c *Context) newCoordSeq(gs *C.struct_GEOSCoordSeq_t, addCleanup bool) *CoordSeq {
	if gs == nil {
		return nil
	}
	var (
		dimensions C.uint
		size       C.uint
	)
	if C.GEOSCoordSeq_getDimensions_r(c.cHandle, gs, &dimensions) == 0 {
		panic(c.err)
	}
	if C.GEOSCoordSeq_getSize_r(c.cHandle, gs, &size) == 0 {
		panic(c.err)
	}
	s := &CoordSeq{
		context:    c,
		cCoordSeq:  gs,
		dimensions: int(dimensions),
		size:       int(size),
	}
	if addCleanup {
		runtime.AddCleanup(s, func(cCoordSeq *C.struct_GEOSCoordSeq_t) {
			C.GEOSCoordSeq_destroy_r(s.context.cHandle, cCoordSeq)
		}, s.cCoordSeq)
	}
	return s
}

func (c *Context) newCoordsFromGEOSCoordSeq(s *C.struct_GEOSCoordSeq_t) [][]float64 {
	var dimensions C.uint
	if C.GEOSCoordSeq_getDimensions_r(c.cHandle, s, &dimensions) == 0 {
		panic(c.err)
	}

	var size C.uint
	if C.GEOSCoordSeq_getSize_r(c.cHandle, s, &size) == 0 {
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
	if C.GEOSCoordSeq_copyToBuffer_r(c.cHandle, s, (*C.double)(&flatCoords[0]), hasZ, hasM) == 0 {
		panic(c.err)
	}
	coords := make([][]float64, size)
	for i := range coords {
		coord := flatCoords[i*int(dimensions) : (i+1)*int(dimensions) : (i+1)*int(dimensions)]
		coords[i] = coord
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

	dimensions := len(coords[0])
	flatCoords := make([]float64, len(coords)*dimensions)
	for i, coord := range coords {
		copy(flatCoords[i*dimensions:(i+1)*dimensions], coord)
	}
	return C.GEOSCoordSeq_copyFromBuffer_r(c.cHandle, (*C.double)(unsafe.Pointer(&flatCoords[0])), C.uint(len(coords)), hasZ, hasM)
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
	if C.c_GEOSGeomGetInfo_r(c.cHandle, geom, &typeID, &numGeometries, &numPoints, &numInteriorRings) == 0 {
		panic(c.err)
	}
	g := &Geom{
		context:          c,
		cGeom:            geom,
		parent:           parent,
		typeID:           TypeID(typeID),
		numGeometries:    int(numGeometries),
		numInteriorRings: int(numInteriorRings),
		numPoints:        int(numPoints),
	}
	runtime.AddCleanup(g, func(cGeom *C.struct_GEOSGeom_t) {
		C.GEOSGeom_destroy_r(g.context.cHandle, cGeom)
	}, g.cGeom)
	return g
}

func (c *Context) newNonNilBufferParams(p *C.struct_GEOSBufParams_t) *BufferParams {
	if p == nil {
		panic(c.err)
	}
	return c.newBufParams(p)
}

func (c *Context) newNonNilCoordSeq(s *C.struct_GEOSCoordSeq_t) *CoordSeq {
	if s == nil {
		panic(c.err)
	}
	return c.newCoordSeq(s, true)
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
