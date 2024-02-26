package geos_test

import (
	"runtime"
	"testing"

	"github.com/alecthomas/assert/v2"

	"github.com/twpayne/go-geos"
)

func TestBufferWithParams(t *testing.T) {
	defer runtime.GC() // Exercise finalizers.
	c := geos.NewContext()
	p := c.NewBufferParams()
	assert.NotZero(t, p)
	assert.NotZero(t, p.SetJoinStyle(geos.BufJoinStyleMitre))
	assert.NotZero(t, p.SetEndCapStyle(geos.BufCapStyleSquare))
	assert.NotZero(t, p.SetMitreLimit(1))
	assert.NotZero(t, p.SetQuadrantSegments(1))
	assert.NotZero(t, p.SetSingleSided(true))
	g := c.NewLineString([][]float64{{0, 0}, {1, 0}}).BufferWithParams(p, 1)
	assert.NotZero(t, g)
	assert.Equal(t, geos.TypeIDPolygon, g.TypeID())
	assert.Equal(t, [][]float64{{1, 0}, {0, 0}, {0, 1}, {1, 1}, {1, 0}}, g.ExteriorRing().CoordSeq().ToCoords())
}
