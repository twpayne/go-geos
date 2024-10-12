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
	if geos.VersionCompare(3, 12, 0) >= 0 {
		assert.True(t, unitSquare.ContainsXY(0.5, 0.5))
	}
	assert.False(t, unitSquare.ContainsXY(2, 2))
	assert.False(t, unitSquare.CoveredBy(middleSquare))
	assert.True(t, unitSquare.Covers(middleSquare))
	assert.False(t, unitSquare.Crosses(middleSquare))
	assert.False(t, unitSquare.Disjoint(middleSquare))
	assert.False(t, unitSquare.DistanceWithin(mustNewGeomFromWKT(t, c, "POINT (1.5 0.5)"), 0.1))
	assert.True(t, unitSquare.Intersects(middleSquare))
	if geos.VersionCompare(3, 12, 0) >= 0 {
		assert.True(t, unitSquare.IntersectsXY(0.5, 0.5))
		assert.False(t, unitSquare.IntersectsXY(2, 2))
	}
	assert.Equal(t, [][]float64{{1, 1}, {2, 2}}, unitSquare.NearestPoints(mustNewGeomFromWKT(t, c, "POINT (2 2)")).ToCoords())
	assert.False(t, unitSquare.Overlaps(middleSquare))
	assert.False(t, unitSquare.Touches(middleSquare))
	assert.False(t, unitSquare.Within(middleSquare))
}
