package geos_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/twpayne/go-geos"
)

func mustNewGeomFromWKT(t *testing.T, c *geos.Context, wkt string) *geos.Geom {
	t.Helper()
	geom, err := c.NewGeomFromWKT(wkt)
	require.NoError(t, err)
	require.True(t, geom.IsValid())
	return geom
}
