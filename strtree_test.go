package geos_test

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/twpayne/go-geos"
)

func TestSTRtree(t *testing.T) {
	defer runtime.GC() // Exercise finalizers.
	c := geos.NewContext()

	tree := c.NewSTRtree(4)

	allItems := func() map[any]struct{} {
		result := make(map[any]struct{})
		tree.Iterate(func(value any) {
			result[value] = struct{}{}
		})
		return result
	}
	assert.Equal(t, map[any]struct{}{}, allItems())

	g1 := mustNewGeomFromWKT(t, c, "POINT (0 0)")
	assert.NoError(t, tree.Insert(g1, 1))
	assert.Equal(t, map[any]struct{}{
		1: {},
	}, allItems())

	g2 := mustNewGeomFromWKT(t, c, "POINT (0 2)")
	assert.NoError(t, tree.Insert(g2, 2))
	assert.Equal(t, map[any]struct{}{
		1: {},
		2: {},
	}, allItems())

	items := make(map[any]struct{})
	tree.Query(mustNewGeomFromWKT(t, c, "POLYGON ((-1 -1,1 -1,1 1,-1 1,-1 -1))"), func(value any) {
		items[value] = struct{}{}
	})
	assert.Equal(t, map[any]struct{}{
		1: {},
	}, items)

	require.True(t, tree.Remove(g1, 1))
	if false {
		// Items removed with GEOSSTRtree_remove_r are still returned by
		// STRtree.Iterate. See https://github.com/libgeos/geos/issues/833.
		assert.Equal(t, map[any]struct{}{
			2: {},
		}, allItems())
	}

	items2 := make(map[any]struct{})
	tree.Query(mustNewGeomFromWKT(t, c, "POLYGON ((-1 -1,1 -1,1 1,-1 1,-1 -1))"), func(value any) {
		items2[value] = struct{}{}
	})
	assert.Equal(t, map[any]struct{}{}, items2)
}

func TestSTRtreeNearest(t *testing.T) {
	defer runtime.GC() // Exercise finalizers.
	c := geos.NewContext()

	tree := c.NewSTRtree(8)
	g1 := mustNewGeomFromWKT(t, c, "POINT (0 1)")
	assert.NoError(t, tree.Insert(g1, g1))
	g2 := mustNewGeomFromWKT(t, c, "POINT (0 2)")
	assert.NoError(t, tree.Insert(g2, g2))
	g4 := mustNewGeomFromWKT(t, c, "POINT (0 4)")
	assert.NoError(t, tree.Insert(g4, g4))

	assert.Equal(t, g2, tree.Nearest(g1, g1, func(value any) *geos.Geom {
		return value.(*geos.Geom) //nolint:forcetypeassert
	}))
	assert.Equal(t, g1, tree.Nearest(g2, g2, func(value any) *geos.Geom {
		return value.(*geos.Geom) //nolint:forcetypeassert
	}))
	assert.Equal(t, g2, tree.Nearest(g4, g4, func(value any) *geos.Geom {
		return value.(*geos.Geom) //nolint:forcetypeassert
	}))
}

func TestSTRtreeLoad(t *testing.T) {
	defer runtime.GC() // Exercise finalizers.
	c := geos.NewContext()

	points := make(map[[2]int]*geos.Geom, 256*256)
	for x := 0; x < 256; x++ {
		for y := 0; y < 256; y++ {
			value := [2]int{x, y}
			points[value] = c.NewPoint([]float64{float64(x), float64(y)})
		}
	}

	tree := c.NewSTRtree(8)
	for value, geom := range points {
		require.NoError(t, tree.Insert(geom, value))
	}

	items := make(map[[2]int]struct{})
	tree.Query(mustNewGeomFromWKT(t, c, "POLYGON ((0 0,256 0,256 256,0 256,0 0))"), func(v any) {
		value, ok := v.([2]int)
		require.True(t, ok)
		items[value] = struct{}{}
	})
	assert.Equal(t, 256*256, len(items))

	for x := 0; x < 256; x++ {
		for y := 0; y < 256; y++ {
			if (x+y)%2 == 0 {
				value := [2]int{x, y}
				require.True(t, tree.Remove(points[value], value))
			}
		}
	}

	runtime.GC()

	itemsAfterRemove := make(map[[2]int]struct{})
	tree.Query(mustNewGeomFromWKT(t, c, "POLYGON ((0 0,256 0,256 256,0 256,0 0))"), func(value any) {
		array, ok := value.([2]int)
		require.True(t, ok)
		itemsAfterRemove[array] = struct{}{}
	})
	assert.Equal(t, 256*256/2, len(itemsAfterRemove))
}
