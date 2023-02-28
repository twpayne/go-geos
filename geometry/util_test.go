package geometry_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/twpayne/go-geos"
	"github.com/twpayne/go-geos/geometry"
)

func mustNewGeometryFromWKT(t *testing.T, wkt string) *geometry.Geometry {
	t.Helper()
	geom, err := geos.NewGeomFromWKT(wkt)
	require.NoError(t, err)
	return &geometry.Geometry{Geom: geom}
}
