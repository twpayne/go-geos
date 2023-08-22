package geos_test

import (
	"testing"

	"github.com/alecthomas/assert/v2"

	"github.com/twpayne/go-geos"
)

func mustNewGeomFromWKT(t *testing.T, c *geos.Context, wkt string) *geos.Geom {
	t.Helper()
	geom, err := c.NewGeomFromWKT(wkt)
	assert.NoError(t, err)
	assert.True(t, geom.IsValid())
	return geom
}

func newInvalidGeomFromWKT(t *testing.T, c *geos.Context, wkt string) *geos.Geom {
	t.Helper()
	geom, err := c.NewGeomFromWKT(wkt)
	assert.NoError(t, err)
	assert.False(t, geom.IsValid())
	return geom
}
