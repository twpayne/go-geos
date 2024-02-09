package geometry

import geos "github.com/twpayne/go-geos"

// initialStringBufferSize is the initial size of strings.Buffers used for
// building GeoJSON and KML representations.
const initialStringBufferSize = 1024

// A Geometry is a geometry.
type Geometry struct {
	*geos.Geom
}

// Must panics with err if err is non-nil, otherwise it returns g.
func Must(g *Geometry, err error) *Geometry {
	if err != nil {
		panic(err)
	}
	return g
}

// NewGeometry returns a new Geometry using geom.
func NewGeometry(geom *geos.Geom) *Geometry {
	return &Geometry{Geom: geom}
}

// Bounds returns g's bounds.
func (g *Geometry) Bounds() *geos.Box2D {
	return g.Geom.Bounds()
}

// Destroy destroys g's geom.
func (g *Geometry) Destroy() {
	g.Geom.Destroy()
	g.Geom = nil
}

// SetSRID sets g's SRID.
func (g *Geometry) SetSRID(srid int) *Geometry {
	g.Geom.SetSRID(srid)
	return g
}
