package geometry_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/twpayne/go-geos"
	"github.com/twpayne/go-geos/geometry"
)

func TestGob(t *testing.T) {
	g := geometry.NewGeometry(geos.NewPoint([]float64{1, 2}))
	data, err := g.GobEncode()
	require.NoError(t, err)
	var geom geometry.Geometry
	require.NoError(t, geom.GobDecode(data))
	assert.True(t, g.Geom.Equals(geom.Geom))
}
