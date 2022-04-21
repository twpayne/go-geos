package geos

import (
	"math"
	"runtime"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGeometryConstructors(t *testing.T) {
	for _, tc := range []struct {
		name        string
		newGeomFunc func(*Context) *Geom
		expectedWKT string
	}{
		{
			name: "NewCollection_MultiPoint_empty",
			newGeomFunc: func(c *Context) *Geom {
				return c.NewCollection(MultiPointTypeID, nil)
			},
			expectedWKT: "MULTIPOINT EMPTY",
		},
		{
			name: "NewCollection_MultiPoint_one",
			newGeomFunc: func(c *Context) *Geom {
				return c.NewCollection(MultiPointTypeID, []*Geom{
					c.NewPoint([]float64{0, 1}),
				})
			},
			expectedWKT: "MULTIPOINT ((0 1))",
		},
		{
			name: "NewCollection_MultiPoint_many",
			newGeomFunc: func(c *Context) *Geom {
				return c.NewCollection(MultiPointTypeID, []*Geom{
					c.NewPoint([]float64{0, 1}),
					c.NewPoint([]float64{2, 3}),
					c.NewPoint([]float64{4, 5}),
				})
			},
			expectedWKT: "MULTIPOINT ((0 1), (2 3), (4 5))",
		},
		{
			name: "NewCollection_MultiLineString_empty",
			newGeomFunc: func(c *Context) *Geom {
				return c.NewCollection(MultiLineStringTypeID, nil)
			},
			expectedWKT: "MULTILINESTRING EMPTY",
		},
		{
			name: "NewCollection_MultiPolygon_empty",
			newGeomFunc: func(c *Context) *Geom {
				return c.NewCollection(MultiPolygonTypeID, nil)
			},
			expectedWKT: "MULTIPOLYGON EMPTY",
		},
		{
			name: "NewCollection_GeometryCollection_empty",
			newGeomFunc: func(c *Context) *Geom {
				return c.NewCollection(GeometryCollectionTypeID, nil)
			},
			expectedWKT: "GEOMETRYCOLLECTION EMPTY",
		},
		{
			name: "NewCollection_GeometryCollection_many",
			newGeomFunc: func(c *Context) *Geom {
				return c.NewCollection(GeometryCollectionTypeID, []*Geom{
					c.NewPoint([]float64{0, 1}),
					c.NewLineString([][]float64{{2, 3}, {4, 5}}),
					c.NewCollection(MultiPointTypeID, []*Geom{
						c.NewPoint([]float64{6, 7}),
					}),
				})
			},
			expectedWKT: "GEOMETRYCOLLECTION (POINT (0 1), LINESTRING (2 3, 4 5), MULTIPOINT (6 7))",
		},
		{
			name: "NewEmptyCollection_MultiPoint",
			newGeomFunc: func(c *Context) *Geom {
				return c.NewEmptyCollection(MultiPointTypeID)
			},
			expectedWKT: "MULTIPOINT EMPTY",
		},
		{
			name: "NewEmptyCollection_MultiLineString",
			newGeomFunc: func(c *Context) *Geom {
				return c.NewEmptyCollection(MultiLineStringTypeID)
			},
			expectedWKT: "MULTILINESTRING EMPTY",
		},
		{
			name: "NewEmptyCollection_MultiPolygon",
			newGeomFunc: func(c *Context) *Geom {
				return c.NewEmptyCollection(MultiPolygonTypeID)
			},
			expectedWKT: "MULTIPOLYGON EMPTY",
		},
		{
			name: "NewEmptyCollection_GeometryCollection",
			newGeomFunc: func(c *Context) *Geom {
				return c.NewEmptyCollection(GeometryCollectionTypeID)
			},
			expectedWKT: "GEOMETRYCOLLECTION EMPTY",
		},
		{
			name: "NewEmptyPoint",
			newGeomFunc: func(c *Context) *Geom {
				return c.NewEmptyPoint()
			},
			expectedWKT: "POINT EMPTY",
		},
		{
			name: "NewGeomFromBounds_polygon",
			newGeomFunc: func(c *Context) *Geom {
				return c.NewGeomFromBounds(&Bounds{MinX: 0, MinY: 1, MaxX: 2, MaxY: 3})
			},
			expectedWKT: "POLYGON ((0 1, 2 1, 2 3, 0 3, 0 1))",
		},
		{
			name: "NewGeomFromBounds_empty",
			newGeomFunc: func(c *Context) *Geom {
				return c.NewGeomFromBounds(&Bounds{MinX: math.Inf(1), MinY: math.Inf(1), MaxX: math.Inf(-1), MaxY: math.Inf(-1)})
			},
			expectedWKT: "POINT EMPTY",
		},
		{
			name: "NewGeomFromBounds_point",
			newGeomFunc: func(c *Context) *Geom {
				return c.NewGeomFromBounds(&Bounds{MinX: 0, MinY: 1, MaxX: 0, MaxY: 1})
			},
			expectedWKT: "POINT (0 1)",
		},
		{
			name: "NewPoint",
			newGeomFunc: func(c *Context) *Geom {
				return c.NewPoint([]float64{1, 2})
			},
			expectedWKT: "POINT (1 2)",
		},
		{
			name: "NewLinearRing",
			newGeomFunc: func(c *Context) *Geom {
				return c.NewLinearRing([][]float64{{1, 2}, {3, 4}, {5, 6}, {1, 2}})
			},
			expectedWKT: "LINEARRING (1 2, 3 4, 5 6, 1 2)",
		},
		{
			name: "NewEmptyLineString",
			newGeomFunc: func(c *Context) *Geom {
				return c.NewEmptyLineString()
			},
			expectedWKT: "LINESTRING EMPTY",
		},
		{
			name: "NewLineString",
			newGeomFunc: func(c *Context) *Geom {
				return c.NewLineString([][]float64{{1, 2}, {3, 4}})
			},
			expectedWKT: "LINESTRING (1 2, 3 4)",
		},
		{
			name: "NewEmptyPolygon",
			newGeomFunc: func(c *Context) *Geom {
				return c.NewEmptyPolygon()
			},
			expectedWKT: "POLYGON EMPTY",
		},
		{
			name: "NewPolygon_empty",
			newGeomFunc: func(c *Context) *Geom {
				return c.NewPolygon(nil)
			},
			expectedWKT: "POLYGON EMPTY",
		},
		{
			name: "NewPolygon",
			newGeomFunc: func(c *Context) *Geom {
				return c.NewPolygon([][][]float64{{{0, 0}, {1, 1}, {0, 1}, {0, 0}}})
			},
			expectedWKT: "POLYGON ((0 0, 1 1, 0 1, 0 0))",
		},
		{
			name: "NewPolygon_with_hole",
			newGeomFunc: func(c *Context) *Geom {
				return c.NewPolygon([][][]float64{
					{{0, 0}, {3, 0}, {3, 3}, {0, 3}, {0, 0}},
					{{1, 1}, {1, 2}, {2, 2}, {2, 1}, {1, 1}},
				})
			},
			expectedWKT: "POLYGON ((0 0, 3 0, 3 3, 0 3, 0 0), (1 1, 1 2, 2 2, 2 1, 1 1))",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			defer runtime.GC() // Exercise finalizers.
			c := NewContext()
			g := tc.newGeomFunc(c)
			require.NotNil(t, g)
			expectedGeom, err := c.NewGeomFromWKT(tc.expectedWKT)
			require.NoError(t, err)
			assert.Equal(t, expectedGeom.TypeID(), g.TypeID())
			assert.True(t, g.Equals(expectedGeom))
		})
	}
}

func TestFinalizeFunc(t *testing.T) {
	var wg sync.WaitGroup
	finalizeHookCalled := false
	wg.Add(1)
	c := NewContext(WithGeomFinalizeFunc(func(g *Geom) {
		defer wg.Done()
		finalizeHookCalled = true
	}))
	_ = c.NewPoint([]float64{0, 0})
	runtime.GC()
	wg.Wait()
	assert.True(t, finalizeHookCalled)
}

func TestMultipleContexts(t *testing.T) {
	c1, c2 := NewContext(), NewContext()
	g1s, g2s := []*Geom{}, []*Geom{}
	for _, wkt := range []string{
		"POINT (0 0)",
		"LINESTRING (0 0, 0 1)",
		"POLYGON ((0 0, 1 0, 1 1, 0 0))",
	} {
		g1 := mustNewGeomFromWKT(t, c1, wkt)
		g1s = append(g1s, g1)
		g2 := mustNewGeomFromWKT(t, c2, wkt)
		g2s = append(g2s, g2)
	}
	for _, g1 := range g1s {
		for _, g2 := range g2s {
			assert.Equal(t, g1.Contains(g2), g2.Contains(g1))
			assert.Equal(t, g1.Equals(g2), g2.Equals(g1))
			assert.Equal(t, g1.Intersects(g2), g2.Intersects(g1))
			g2CloneInC1 := c1.Clone(g2)
			assert.Equal(t, g1.Contains(g2CloneInC1), g2CloneInC1.Contains(g1))
			assert.Equal(t, g1.Equals(g2CloneInC1), g2CloneInC1.Equals(g1))
			assert.Equal(t, g1.Intersects(g2CloneInC1), g2CloneInC1.Intersects(g1))
		}
	}
}

func TestNewPoints(t *testing.T) {
	c := NewContext()
	assert.Nil(t, c.NewPoints(nil))
	gs := c.NewPoints([][]float64{{1, 2}, {3, 4}})
	assert.Len(t, gs, 2)
	assert.True(t, gs[0].Equals(mustNewGeomFromWKT(t, c, "POINT (1 2)")))
	assert.True(t, gs[1].Equals(mustNewGeomFromWKT(t, c, "POINT (3 4)")))
}
