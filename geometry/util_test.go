package geometry_test

import (
	"testing"

	"github.com/alecthomas/assert/v2"

	"github.com/twpayne/go-geos"
	"github.com/twpayne/go-geos/geometry"
)

func mustNewGeometryFromWKT(t *testing.T, wkt string) *geometry.Geometry {
	t.Helper()
	geom, err := geos.NewGeomFromWKT(wkt)
	assert.NoError(t, err)
	return &geometry.Geometry{Geom: geom}
}
