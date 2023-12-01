package geos_test

import (
	"fmt"
	"testing"

	"github.com/alecthomas/assert/v2"

	"github.com/twpayne/go-geos"
)

func mustNewGeomFromWKT(t *testing.T, c *geos.Context, wkt string) *geos.Geom {
	t.Helper()
	geom, err := c.NewGeomFromWKT(wkt)
	if err != nil {
		err = fmt.Errorf("%s: %w", wkt, err)
	}
	assert.NoError(t, err)
	assert.True(t, geom.IsValid())
	return geom
}

func newInvalidGeomFromWKT(t *testing.T, c *geos.Context, wkt string) *geos.Geom {
	t.Helper()
	geom, err := c.NewGeomFromWKT(wkt)
	if err != nil {
		err = fmt.Errorf("%s: %w", wkt, err)
	}
	assert.NoError(t, err)
	assert.False(t, geom.IsValid())
	return geom
}
