package geos

// #include "go-geos.h"
import "C"

import "runtime"

// A PrepGeom is a prepared geometry.
type PrepGeom struct {
	owner     *Geom
	cPrepGeom *C.struct_GEOSPrepGeom_t
}

// Prepare prepares g.
func (g *Geom) Prepare() *PrepGeom {
	g.context.mutex.Lock()
	defer g.context.mutex.Unlock()
	cPrepGeom := C.GEOSPrepare_r(g.context.cHandle, g.cGeom)
	prepGeom := &PrepGeom{
		owner:     g,
		cPrepGeom: cPrepGeom,
	}
	g.context.ref()
	runtime.AddCleanup(prepGeom, g.context.destroyPrepGeom, cPrepGeom)
	return prepGeom
}

// Contains returns if pg contains g.
func (pg *PrepGeom) Contains(g *Geom) bool {
	pg.owner.context.mutex.Lock()
	defer pg.owner.context.mutex.Unlock()
	if g.context != pg.owner.context {
		g.context.mutex.Lock()
		defer g.context.mutex.Unlock()
	}
	switch C.GEOSPreparedContains_r(pg.owner.context.cHandle, pg.cPrepGeom, g.cGeom) {
	case 0:
		return false
	case 1:
		return true
	default:
		panic(pg.owner.context.err)
	}
}

// ContainsProperly returns if pg contains g properly.
func (pg *PrepGeom) ContainsProperly(g *Geom) bool {
	pg.owner.context.mutex.Lock()
	defer pg.owner.context.mutex.Unlock()
	if g.context != pg.owner.context {
		g.context.mutex.Lock()
		defer g.context.mutex.Unlock()
	}
	switch C.GEOSPreparedContainsProperly_r(pg.owner.context.cHandle, pg.cPrepGeom, g.cGeom) {
	case 0:
		return false
	case 1:
		return true
	default:
		panic(pg.owner.context.err)
	}
}

// ContainsXY returns if pg contains the point (x, y).
func (pg *PrepGeom) ContainsXY(x, y float64) bool {
	pg.owner.context.mutex.Lock()
	defer pg.owner.context.mutex.Unlock()
	switch C.GEOSPreparedContainsXY_r(pg.owner.context.cHandle, pg.cPrepGeom, C.double(x), C.double(y)) {
	case 0:
		return false
	case 1:
		return true
	default:
		panic(pg.owner.context.err)
	}
}

// CoveredBy returns if pg is covered by g.
func (pg *PrepGeom) CoveredBy(g *Geom) bool {
	pg.owner.context.mutex.Lock()
	defer pg.owner.context.mutex.Unlock()
	if g.context != pg.owner.context {
		g.context.mutex.Lock()
		defer g.context.mutex.Unlock()
	}
	switch C.GEOSPreparedCoveredBy_r(pg.owner.context.cHandle, pg.cPrepGeom, g.cGeom) {
	case 0:
		return false
	case 1:
		return true
	default:
		panic(pg.owner.context.err)
	}
}

// Covers returns if pg covers g.
func (pg *PrepGeom) Covers(g *Geom) bool {
	pg.owner.context.mutex.Lock()
	defer pg.owner.context.mutex.Unlock()
	if g.context != pg.owner.context {
		g.context.mutex.Lock()
		defer g.context.mutex.Unlock()
	}
	switch C.GEOSPreparedCovers_r(pg.owner.context.cHandle, pg.cPrepGeom, g.cGeom) {
	case 0:
		return false
	case 1:
		return true
	default:
		panic(pg.owner.context.err)
	}
}

// Crosses returns if pg crosses g.
func (pg *PrepGeom) Crosses(g *Geom) bool {
	pg.owner.context.mutex.Lock()
	defer pg.owner.context.mutex.Unlock()
	if g.context != pg.owner.context {
		g.context.mutex.Lock()
		defer g.context.mutex.Unlock()
	}
	switch C.GEOSPreparedCrosses_r(pg.owner.context.cHandle, pg.cPrepGeom, g.cGeom) {
	case 0:
		return false
	case 1:
		return true
	default:
		panic(pg.owner.context.err)
	}
}

// Disjoint returns if pg is disjoint from g.
func (pg *PrepGeom) Disjoint(g *Geom) bool {
	pg.owner.context.mutex.Lock()
	defer pg.owner.context.mutex.Unlock()
	if g.context != pg.owner.context {
		g.context.mutex.Lock()
		defer g.context.mutex.Unlock()
	}
	switch C.GEOSPreparedDisjoint_r(pg.owner.context.cHandle, pg.cPrepGeom, g.cGeom) {
	case 0:
		return false
	case 1:
		return true
	default:
		panic(pg.owner.context.err)
	}
}

// DistanceWithin returns if pg is within dist g.
func (pg *PrepGeom) DistanceWithin(g *Geom, dist float64) bool {
	pg.owner.context.mutex.Lock()
	defer pg.owner.context.mutex.Unlock()
	if g.context != pg.owner.context {
		g.context.mutex.Lock()
		defer g.context.mutex.Unlock()
	}
	switch C.GEOSPreparedDistanceWithin_r(pg.owner.context.cHandle, pg.cPrepGeom, g.cGeom, C.double(dist)) {
	case 0:
		return false
	case 1:
		return true
	default:
		panic(pg.owner.context.err)
	}
}

// Intersects returns if pg contains g.
func (pg *PrepGeom) Intersects(g *Geom) bool {
	pg.owner.context.mutex.Lock()
	defer pg.owner.context.mutex.Unlock()
	if g.context != pg.owner.context {
		g.context.mutex.Lock()
		defer g.context.mutex.Unlock()
	}
	switch C.GEOSPreparedIntersects_r(pg.owner.context.cHandle, pg.cPrepGeom, g.cGeom) {
	case 0:
		return false
	case 1:
		return true
	default:
		panic(pg.owner.context.err)
	}
}

// IntersectsXY returns if pg intersects the point at (x, y).
func (pg *PrepGeom) IntersectsXY(x, y float64) bool {
	pg.owner.context.mutex.Lock()
	defer pg.owner.context.mutex.Unlock()
	switch C.GEOSPreparedIntersectsXY_r(pg.owner.context.cHandle, pg.cPrepGeom, C.double(x), C.double(y)) {
	case 0:
		return false
	case 1:
		return true
	default:
		panic(pg.owner.context.err)
	}
}

// NearestPoints returns if pg overlaps g.
func (pg *PrepGeom) NearestPoints(g *Geom) *CoordSeq {
	pg.owner.context.mutex.Lock()
	defer pg.owner.context.mutex.Unlock()
	if g.context != pg.owner.context {
		g.context.mutex.Lock()
		defer g.context.mutex.Unlock()
	}
	return pg.owner.context.newNonNilCoordSeq(C.GEOSPreparedNearestPoints_r(pg.owner.context.cHandle, pg.cPrepGeom, g.cGeom))
}

// Overlaps returns if pg overlaps g.
func (pg *PrepGeom) Overlaps(g *Geom) bool {
	pg.owner.context.mutex.Lock()
	defer pg.owner.context.mutex.Unlock()
	if g.context != pg.owner.context {
		g.context.mutex.Lock()
		defer g.context.mutex.Unlock()
	}
	switch C.GEOSPreparedOverlaps_r(pg.owner.context.cHandle, pg.cPrepGeom, g.cGeom) {
	case 0:
		return false
	case 1:
		return true
	default:
		panic(pg.owner.context.err)
	}
}

// Touches returns if pg contains g.
func (pg *PrepGeom) Touches(g *Geom) bool {
	pg.owner.context.mutex.Lock()
	defer pg.owner.context.mutex.Unlock()
	if g.context != pg.owner.context {
		g.context.mutex.Lock()
		defer g.context.mutex.Unlock()
	}
	switch C.GEOSPreparedTouches_r(pg.owner.context.cHandle, pg.cPrepGeom, g.cGeom) {
	case 0:
		return false
	case 1:
		return true
	default:
		panic(pg.owner.context.err)
	}
}

// Within returns if pg is within g.
func (pg *PrepGeom) Within(g *Geom) bool {
	pg.owner.context.mutex.Lock()
	defer pg.owner.context.mutex.Unlock()
	if g.context != pg.owner.context {
		g.context.mutex.Lock()
		defer g.context.mutex.Unlock()
	}
	switch C.GEOSPreparedWithin_r(pg.owner.context.cHandle, pg.cPrepGeom, g.cGeom) {
	case 0:
		return false
	case 1:
		return true
	default:
		panic(pg.owner.context.err)
	}
}

func (c *Context) destroyPrepGeom(cPrepGeom *C.struct_GEOSPrepGeom_t) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	C.GEOSPreparedGeom_destroy_r(c.cHandle, cPrepGeom)
	c.unref()
}
