package geos_test

import (
	"testing"

	"github.com/alecthomas/assert/v2"

	"github.com/twpayne/go-geos"
)

func TestGeomDestroy(t *testing.T) {
	c := geos.NewContext()
	g := mustNewGeomFromWKT(t, c, "POLYGON ((0 0, 1 0, 1 1, 0 1, 0 0))")
	g.Destroy()
	g.Destroy()
}

func TestGeomDestroySubGeometry(t *testing.T) {
	c := geos.NewContext()
	g := mustNewGeomFromWKT(t, c, "POLYGON ((0 0, 10 0, 10 10, 0 10, 0 0), (1 1, 2 1, 2 2, 1 2, 1 1))")
	ring := g.InteriorRing(0)
	ring.Destroy()
	assert.Equal(t, 5, ring.NumPoints())
	assert.Equal(t, 1, g.NumInteriorRings())
}

func TestGeomDestroyOwnedGeometry(t *testing.T) {
	c := geos.NewContext()
	g := mustNewGeomFromWKT(t, c, "POINT (1 2)")
	collection := c.NewCollection(geos.TypeIDGeometryCollection, []*geos.Geom{g})
	g.Destroy()
	assert.Equal(t, "GEOMETRYCOLLECTION (POINT (1 2))", collection.ToWKT())
}
