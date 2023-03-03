package geos_test

import (
	"math"
	"runtime"
	"strconv"
	"sync"
	"testing"

	"github.com/alecthomas/assert/v2"

	"github.com/twpayne/go-geos"
)

func TestGeometryConstructors(t *testing.T) {
	for _, tc := range []struct {
		name        string
		newGeomFunc func(*geos.Context) *geos.Geom
		expectedWKT string
	}{
		{
			name: "NewCollection_MultiPoint_empty",
			newGeomFunc: func(c *geos.Context) *geos.Geom {
				return c.NewCollection(geos.TypeIDMultiPoint, nil)
			},
			expectedWKT: "MULTIPOINT EMPTY",
		},
		{
			name: "NewCollection_MultiPoint_one",
			newGeomFunc: func(c *geos.Context) *geos.Geom {
				return c.NewCollection(geos.TypeIDMultiPoint, []*geos.Geom{
					c.NewPoint([]float64{0, 1}),
				})
			},
			expectedWKT: "MULTIPOINT ((0 1))",
		},
		{
			name: "NewCollection_MultiPoint_many",
			newGeomFunc: func(c *geos.Context) *geos.Geom {
				return c.NewCollection(geos.TypeIDMultiPoint, []*geos.Geom{
					c.NewPoint([]float64{0, 1}),
					c.NewPoint([]float64{2, 3}),
					c.NewPoint([]float64{4, 5}),
				})
			},
			expectedWKT: "MULTIPOINT ((0 1), (2 3), (4 5))",
		},
		{
			name: "NewCollection_MultiLineString_empty",
			newGeomFunc: func(c *geos.Context) *geos.Geom {
				return c.NewCollection(geos.TypeIDMultiLineString, nil)
			},
			expectedWKT: "MULTILINESTRING EMPTY",
		},
		{
			name: "NewCollection_MultiPolygon_empty",
			newGeomFunc: func(c *geos.Context) *geos.Geom {
				return c.NewCollection(geos.TypeIDMultiPolygon, nil)
			},
			expectedWKT: "MULTIPOLYGON EMPTY",
		},
		{
			name: "NewCollection_GeometryCollection_empty",
			newGeomFunc: func(c *geos.Context) *geos.Geom {
				return c.NewCollection(geos.TypeIDGeometryCollection, nil)
			},
			expectedWKT: "GEOMETRYCOLLECTION EMPTY",
		},
		{
			name: "NewCollection_GeometryCollection_many",
			newGeomFunc: func(c *geos.Context) *geos.Geom {
				return c.NewCollection(geos.TypeIDGeometryCollection, []*geos.Geom{
					c.NewPoint([]float64{0, 1}),
					c.NewLineString([][]float64{{2, 3}, {4, 5}}),
					c.NewCollection(geos.TypeIDMultiPoint, []*geos.Geom{
						c.NewPoint([]float64{6, 7}),
					}),
				})
			},
			expectedWKT: "GEOMETRYCOLLECTION (POINT (0 1), LINESTRING (2 3, 4 5), MULTIPOINT (6 7))",
		},
		{
			name: "NewEmptyCollection_MultiPoint",
			newGeomFunc: func(c *geos.Context) *geos.Geom {
				return c.NewEmptyCollection(geos.TypeIDMultiPoint)
			},
			expectedWKT: "MULTIPOINT EMPTY",
		},
		{
			name: "NewEmptyCollection_MultiLineString",
			newGeomFunc: func(c *geos.Context) *geos.Geom {
				return c.NewEmptyCollection(geos.TypeIDMultiLineString)
			},
			expectedWKT: "MULTILINESTRING EMPTY",
		},
		{
			name: "NewEmptyCollection_MultiPolygon",
			newGeomFunc: func(c *geos.Context) *geos.Geom {
				return c.NewEmptyCollection(geos.TypeIDMultiPolygon)
			},
			expectedWKT: "MULTIPOLYGON EMPTY",
		},
		{
			name: "NewEmptyCollection_GeometryCollection",
			newGeomFunc: func(c *geos.Context) *geos.Geom {
				return c.NewEmptyCollection(geos.TypeIDGeometryCollection)
			},
			expectedWKT: "GEOMETRYCOLLECTION EMPTY",
		},
		{
			name: "NewEmptyPoint",
			newGeomFunc: func(c *geos.Context) *geos.Geom {
				return c.NewEmptyPoint()
			},
			expectedWKT: "POINT EMPTY",
		},
		{
			name: "NewGeomFromBounds_polygon",
			newGeomFunc: func(c *geos.Context) *geos.Geom {
				return c.NewGeomFromBounds(&geos.Bounds{MinX: 0, MinY: 1, MaxX: 2, MaxY: 3})
			},
			expectedWKT: "POLYGON ((0 1, 2 1, 2 3, 0 3, 0 1))",
		},
		{
			name: "NewGeomFromBounds_empty",
			newGeomFunc: func(c *geos.Context) *geos.Geom {
				return c.NewGeomFromBounds(&geos.Bounds{MinX: math.Inf(1), MinY: math.Inf(1), MaxX: math.Inf(-1), MaxY: math.Inf(-1)})
			},
			expectedWKT: "POINT EMPTY",
		},
		{
			name: "NewGeomFromBounds_point",
			newGeomFunc: func(c *geos.Context) *geos.Geom {
				return c.NewGeomFromBounds(&geos.Bounds{MinX: 0, MinY: 1, MaxX: 0, MaxY: 1})
			},
			expectedWKT: "POINT (0 1)",
		},
		{
			name: "NewPoint",
			newGeomFunc: func(c *geos.Context) *geos.Geom {
				return c.NewPoint([]float64{1, 2})
			},
			expectedWKT: "POINT (1 2)",
		},
		{
			name: "NewLinearRing",
			newGeomFunc: func(c *geos.Context) *geos.Geom {
				return c.NewLinearRing([][]float64{{1, 2}, {3, 4}, {5, 6}, {1, 2}})
			},
			expectedWKT: "LINEARRING (1 2, 3 4, 5 6, 1 2)",
		},
		{
			name: "NewEmptyLineString",
			newGeomFunc: func(c *geos.Context) *geos.Geom {
				return c.NewEmptyLineString()
			},
			expectedWKT: "LINESTRING EMPTY",
		},
		{
			name: "NewLineString",
			newGeomFunc: func(c *geos.Context) *geos.Geom {
				return c.NewLineString([][]float64{{1, 2}, {3, 4}})
			},
			expectedWKT: "LINESTRING (1 2, 3 4)",
		},
		{
			name: "NewEmptyPolygon",
			newGeomFunc: func(c *geos.Context) *geos.Geom {
				return c.NewEmptyPolygon()
			},
			expectedWKT: "POLYGON EMPTY",
		},
		{
			name: "NewPolygon_empty",
			newGeomFunc: func(c *geos.Context) *geos.Geom {
				return c.NewPolygon(nil)
			},
			expectedWKT: "POLYGON EMPTY",
		},
		{
			name: "NewPolygon",
			newGeomFunc: func(c *geos.Context) *geos.Geom {
				return c.NewPolygon([][][]float64{{{0, 0}, {1, 1}, {0, 1}, {0, 0}}})
			},
			expectedWKT: "POLYGON ((0 0, 1 1, 0 1, 0 0))",
		},
		{
			name: "NewPolygon_with_hole",
			newGeomFunc: func(c *geos.Context) *geos.Geom {
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
			c := geos.NewContext()
			g := tc.newGeomFunc(c)
			assert.NotZero(t, g)
			expectedGeom, err := c.NewGeomFromWKT(tc.expectedWKT)
			assert.NoError(t, err)
			assert.Equal(t, expectedGeom.TypeID(), g.TypeID())
			assert.True(t, g.Equals(expectedGeom))
		})
	}
}

func TestFinalizeFunc(t *testing.T) {
	var wg sync.WaitGroup
	finalizeHookCalled := false
	wg.Add(1)
	c := geos.NewContext(geos.WithGeomFinalizeFunc(func(g *geos.Geom) {
		defer wg.Done()
		finalizeHookCalled = true
	}))
	_ = c.NewPoint([]float64{0, 0})
	runtime.GC()
	wg.Wait()
	assert.True(t, finalizeHookCalled)
}

func TestMultipleContexts(t *testing.T) {
	c1, c2 := geos.NewContext(), geos.NewContext()
	g1s, g2s := []*geos.Geom{}, []*geos.Geom{}
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
	c := geos.NewContext()
	assert.Equal(t, nil, c.NewPoints(nil))
	gs := c.NewPoints([][]float64{{1, 2}, {3, 4}})
	assert.Equal(t, 2, len(gs))
	assert.True(t, gs[0].Equals(mustNewGeomFromWKT(t, c, "POINT (1 2)")))
	assert.True(t, gs[1].Equals(mustNewGeomFromWKT(t, c, "POINT (3 4)")))
}

func TestPolygonize(t *testing.T) {
	for _, tc := range []struct {
		name             string
		geomWKTs         []string
		expectedWKT      string
		expectedValidWKT string
	}{
		{
			name:             "empty",
			expectedWKT:      "GEOMETRYCOLLECTION EMPTY",
			expectedValidWKT: "GEOMETRYCOLLECTION EMPTY",
		},
		{
			name: "simple",
			geomWKTs: []string{
				"LINESTRING (0 0,1 0,1 1)",
				"LINESTRING (1 1,0 1,0 0)",
			},
			expectedWKT:      "GEOMETRYCOLLECTION (POLYGON ((0 0,1 0,1 1,0 1,0 0)))",
			expectedValidWKT: "POLYGON ((0 0,1 0,1 1,0 1,0 0))",
		},
		{
			name: "extra_linestring",
			geomWKTs: []string{
				"LINESTRING (0 0,1 0,1 1)",
				"LINESTRING (1 1,0 1,0 0)",
				"LINESTRING (0 0,0 -1)",
			},
			expectedWKT:      "GEOMETRYCOLLECTION (POLYGON ((0 0,1 0,1 1,0 1,0 0)))",
			expectedValidWKT: "POLYGON ((0 0,1 0,1 1,0 1,0 0))",
		},
		{
			name: "two_polygons",
			geomWKTs: []string{
				"LINESTRING (0 0,1 0,1 1)",
				"LINESTRING (1 1,0 1,0 0)",
				"LINESTRING (2 2,3 2,3 3)",
				"LINESTRING (3 3 2 3,2 2)",
			},
			expectedWKT:      "GEOMETRYCOLLECTION (POLYGON ((0 0,1 0,1 1,0 1,0 0)),POLYGON ((2 2,3 2,3 3,2 3,2 2)))",
			expectedValidWKT: "MULTIPOLYGON (((0 0,1 0,1 1,0 1,0 0)),((2 2,3 2,3 3,2 3,2 2)))",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := geos.NewContext()
			geoms := make([]*geos.Geom, 0, len(tc.geomWKTs))
			for _, geomWKT := range tc.geomWKTs {
				geom := mustNewGeomFromWKT(t, c, geomWKT)
				geoms = append(geoms, geom)
			}
			assert.Equal(t, mustNewGeomFromWKT(t, c, tc.expectedWKT), c.Polygonize(geoms))
			assert.Equal(t, mustNewGeomFromWKT(t, c, tc.expectedValidWKT), c.PolygonizeValid(geoms))
		})
	}
}

func TestPolygonizeMultiContext(t *testing.T) {
	c1 := geos.NewContext()
	c2 := geos.NewContext()
	for i := 0; i < 4; i++ {
		assert.Equal(t,
			mustNewGeomFromWKT(t, c1, "GEOMETRYCOLLECTION (POLYGON ((0 0,1 0,1 1,0 1,0 0)))"),
			c1.Polygonize([]*geos.Geom{
				mustNewGeomFromWKT(t, c1, "LINESTRING (0 0,1 0)"),
				mustNewGeomFromWKT(t, c2, "LINESTRING (1 0,1 1)"),
				mustNewGeomFromWKT(t, c1, "LINESTRING (1 1,0 1)"),
				mustNewGeomFromWKT(t, c2, "LINESTRING (0 1,0 0)"),
			}),
		)
	}
}

func TestSegmentIntersection(t *testing.T) {
	for i, tc := range []struct {
		ax0, ay0, ax1, ay1 float64
		bx0, by0, bx1, by1 float64
		cx, cy             float64
		intersects         bool
	}{
		{
			ax0:        0,
			ay0:        0,
			ax1:        1,
			ay1:        1,
			bx0:        0,
			by0:        1,
			bx1:        1,
			by1:        0,
			cx:         0.5,
			cy:         0.5,
			intersects: true,
		},
		{
			ax0: 0,
			ay0: 0,
			ax1: 1,
			ay1: 0,
			bx0: 0,
			by0: 1,
			bx1: 1,
			by1: 1,
		},
		{
			ax0:        0,
			ay0:        0,
			ax1:        1,
			ay1:        0,
			bx0:        0,
			by0:        0,
			bx1:        1,
			by1:        0,
			intersects: true,
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			actualCX, actualCY, actualIntersects := geos.NewContext().SegmentIntersection(tc.ax0, tc.ay0, tc.ax1, tc.ay1, tc.bx0, tc.by0, tc.ax1, tc.by1)
			assert.Equal(t, tc.cx, actualCX)
			assert.Equal(t, tc.cy, actualCY)
			assert.Equal(t, tc.intersects, actualIntersects)
		})
	}
}
