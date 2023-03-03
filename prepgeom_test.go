package geos_test

import (
	"runtime"
	"testing"

	"github.com/alecthomas/assert/v2"

	"github.com/twpayne/go-geos"
)

func TestPrepGeom(t *testing.T) {
	defer runtime.GC() // Exercise finalizers.
	c := geos.NewContext()
	unitSquare := mustNewGeomFromWKT(t, c, "POLYGON ((0 0, 0 1, 1 1, 1 0, 0 0))").Prepare()
	middleSquare := mustNewGeomFromWKT(t, c, "POLYGON ((0.25 0.25, 0.25 0.75, 0.75 0.75, 0.75 0.25, 0.25 0.25))")
	assert.True(t, unitSquare.Contains(middleSquare))
	assert.True(t, unitSquare.ContainsProperly(middleSquare))
	assert.False(t, unitSquare.CoveredBy(middleSquare))
	assert.True(t, unitSquare.Covers(middleSquare))
	assert.False(t, unitSquare.Crosses(middleSquare))
	assert.False(t, unitSquare.Disjoint(middleSquare))
	assert.True(t, unitSquare.Intersects(middleSquare))
	assert.False(t, unitSquare.Overlaps(middleSquare))
	assert.False(t, unitSquare.Touches(middleSquare))
	assert.False(t, unitSquare.Within(middleSquare))
}
