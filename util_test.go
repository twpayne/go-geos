package geos

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func mustNewGeomFromWKT(t *testing.T, c *Context, wkt string) *Geom {
	t.Helper()
	geom, err := c.NewGeomFromWKT(wkt)
	require.NoError(t, err)
	require.True(t, geom.IsValid())
	return geom
}
