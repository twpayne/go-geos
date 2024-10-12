package geos

// #include "go-geos.h"
import "C"

import "runtime"

// A PrepGeom is a prepared geometry.
type PrepGeom struct {
	parent *Geom
	pgeom  *C.struct_GEOSPrepGeom_t
}

// Prepare prepares g.
func (g *Geom) Prepare() *PrepGeom {
	g.context.Lock()
	defer g.context.Unlock()
	pg := &PrepGeom{
		parent: g,
		pgeom:  C.GEOSPrepare_r(g.context.handle, g.geom),
	}
	runtime.SetFinalizer(pg, (*PrepGeom).Destroy)
	return pg
}

// Contains returns if pg contains g.
func (pg *PrepGeom) Contains(g *Geom) bool {
	pg.parent.context.Lock()
	defer pg.parent.context.Unlock()
	if g.context != pg.parent.context {
		g.context.Lock()
		defer g.context.Unlock()
	}
	switch C.GEOSPreparedContains_r(pg.parent.context.handle, pg.pgeom, g.geom) {
	case 0:
		return false
	case 1:
		return true
	default:
		panic(pg.parent.context.err)
	}
}

// ContainsProperly returns if pg contains g properly.
func (pg *PrepGeom) ContainsProperly(g *Geom) bool {
	pg.parent.context.Lock()
	defer pg.parent.context.Unlock()
	if g.context != pg.parent.context {
		g.context.Lock()
		defer g.context.Unlock()
	}
	switch C.GEOSPreparedContainsProperly_r(pg.parent.context.handle, pg.pgeom, g.geom) {
	case 0:
		return false
	case 1:
		return true
	default:
		panic(pg.parent.context.err)
	}
}

// ContainsXY returns if pg contains the point (x, y).
func (pg *PrepGeom) ContainsXY(x, y float64) bool {
	pg.parent.context.Lock()
	defer pg.parent.context.Unlock()
	switch C.GEOSPreparedContainsXY_r(pg.parent.context.handle, pg.pgeom, C.double(x), C.double(y)) {
	case 0:
		return false
	case 1:
		return true
	default:
		panic(pg.parent.context.err)
	}
}

// CoveredBy returns if pg is covered by g.
func (pg *PrepGeom) CoveredBy(g *Geom) bool {
	pg.parent.context.Lock()
	defer pg.parent.context.Unlock()
	if g.context != pg.parent.context {
		g.context.Lock()
		defer g.context.Unlock()
	}
	switch C.GEOSPreparedCoveredBy_r(pg.parent.context.handle, pg.pgeom, g.geom) {
	case 0:
		return false
	case 1:
		return true
	default:
		panic(pg.parent.context.err)
	}
}

// Covers returns if pg covers g.
func (pg *PrepGeom) Covers(g *Geom) bool {
	pg.parent.context.Lock()
	defer pg.parent.context.Unlock()
	if g.context != pg.parent.context {
		g.context.Lock()
		defer g.context.Unlock()
	}
	switch C.GEOSPreparedCovers_r(pg.parent.context.handle, pg.pgeom, g.geom) {
	case 0:
		return false
	case 1:
		return true
	default:
		panic(pg.parent.context.err)
	}
}

// Crosses returns if pg crosses g.
func (pg *PrepGeom) Crosses(g *Geom) bool {
	pg.parent.context.Lock()
	defer pg.parent.context.Unlock()
	if g.context != pg.parent.context {
		g.context.Lock()
		defer g.context.Unlock()
	}
	switch C.GEOSPreparedCrosses_r(pg.parent.context.handle, pg.pgeom, g.geom) {
	case 0:
		return false
	case 1:
		return true
	default:
		panic(pg.parent.context.err)
	}
}

// Destroy destroys pg and all resources associated with s.
func (pg *PrepGeom) Destroy() {
	if pg == nil || pg.parent == nil || pg.parent.context == nil {
		return
	}
	pg.parent.context.Lock()
	defer pg.parent.context.Unlock()
	C.GEOSPreparedGeom_destroy_r(pg.parent.context.handle, pg.pgeom)
	*pg = PrepGeom{} // Clear all references.
}

// Disjoint returns if pg is disjoint from g.
func (pg *PrepGeom) Disjoint(g *Geom) bool {
	pg.parent.context.Lock()
	defer pg.parent.context.Unlock()
	if g.context != pg.parent.context {
		g.context.Lock()
		defer g.context.Unlock()
	}
	switch C.GEOSPreparedDisjoint_r(pg.parent.context.handle, pg.pgeom, g.geom) {
	case 0:
		return false
	case 1:
		return true
	default:
		panic(pg.parent.context.err)
	}
}

// DistanceWithin returns if pg is within dist g.
func (pg *PrepGeom) DistanceWithin(g *Geom, dist float64) bool {
	pg.parent.context.Lock()
	defer pg.parent.context.Unlock()
	if g.context != pg.parent.context {
		g.context.Lock()
		defer g.context.Unlock()
	}
	switch C.GEOSPreparedDistanceWithin_r(pg.parent.context.handle, pg.pgeom, g.geom, C.double(dist)) {
	case 0:
		return false
	case 1:
		return true
	default:
		panic(pg.parent.context.err)
	}
}

// Intersects returns if pg contains g.
func (pg *PrepGeom) Intersects(g *Geom) bool {
	pg.parent.context.Lock()
	defer pg.parent.context.Unlock()
	if g.context != pg.parent.context {
		g.context.Lock()
		defer g.context.Unlock()
	}
	switch C.GEOSPreparedIntersects_r(pg.parent.context.handle, pg.pgeom, g.geom) {
	case 0:
		return false
	case 1:
		return true
	default:
		panic(pg.parent.context.err)
	}
}

// IntersectsXY returns if pg intersects the point at (x, y).
func (pg *PrepGeom) IntersectsXY(x, y float64) bool {
	pg.parent.context.Lock()
	defer pg.parent.context.Unlock()
	switch C.GEOSPreparedIntersectsXY_r(pg.parent.context.handle, pg.pgeom, C.double(x), C.double(y)) {
	case 0:
		return false
	case 1:
		return true
	default:
		panic(pg.parent.context.err)
	}
}

// NearestPoints returns if pg overlaps g.
func (pg *PrepGeom) NearestPoints(g *Geom) *CoordSeq {
	pg.parent.context.Lock()
	defer pg.parent.context.Unlock()
	if g.context != pg.parent.context {
		g.context.Lock()
		defer g.context.Unlock()
	}
	return pg.parent.context.newNonNilCoordSeq(C.GEOSPreparedNearestPoints_r(pg.parent.context.handle, pg.pgeom, g.geom))
}

// Overlaps returns if pg overlaps g.
func (pg *PrepGeom) Overlaps(g *Geom) bool {
	pg.parent.context.Lock()
	defer pg.parent.context.Unlock()
	if g.context != pg.parent.context {
		g.context.Lock()
		defer g.context.Unlock()
	}
	switch C.GEOSPreparedOverlaps_r(pg.parent.context.handle, pg.pgeom, g.geom) {
	case 0:
		return false
	case 1:
		return true
	default:
		panic(pg.parent.context.err)
	}
}

// Touches returns if pg contains g.
func (pg *PrepGeom) Touches(g *Geom) bool {
	pg.parent.context.Lock()
	defer pg.parent.context.Unlock()
	if g.context != pg.parent.context {
		g.context.Lock()
		defer g.context.Unlock()
	}
	switch C.GEOSPreparedTouches_r(pg.parent.context.handle, pg.pgeom, g.geom) {
	case 0:
		return false
	case 1:
		return true
	default:
		panic(pg.parent.context.err)
	}
}

// Within returns if pg is within g.
func (pg *PrepGeom) Within(g *Geom) bool {
	pg.parent.context.Lock()
	defer pg.parent.context.Unlock()
	if g.context != pg.parent.context {
		g.context.Lock()
		defer g.context.Unlock()
	}
	switch C.GEOSPreparedWithin_r(pg.parent.context.handle, pg.pgeom, g.geom) {
	case 0:
		return false
	case 1:
		return true
	default:
		panic(pg.parent.context.err)
	}
}
