package geos_test

import (
	"math"
	"runtime"
	"strconv"
	"testing"

	"github.com/alecthomas/assert/v2"

	"github.com/twpayne/go-geos"
)

func TestGeometryMethods(t *testing.T) {
	for _, tc := range []struct {
		name                                   string
		wkt                                    string
		expectedBounds                         *geos.Box2D
		expectedEmpty                          bool
		expectedEnvelopeWKT                    string
		expectedNumCoordinates                 int
		expectedNumGeometries                  int
		expectedSRID                           int
		expectedType                           string
		expectedTypeID                         geos.TypeID
		expectedArea                           float64
		expectedLength                         float64
		expectedValidWKT                       string
		expectedValidWKTStructureKeepCollapsed string
	}{
		{
			name:                   "point",
			wkt:                    "POINT (0.0000000000000000 0.0000000000000000)",
			expectedBounds:         &geos.Box2D{MinX: 0, MinY: 0, MaxX: 0, MaxY: 0},
			expectedEmpty:          false,
			expectedEnvelopeWKT:    "POINT (0 0)",
			expectedNumCoordinates: 1,
			expectedNumGeometries:  1,
			expectedSRID:           0,
			expectedType:           "Point",
			expectedTypeID:         geos.TypeIDPoint,
			expectedArea:           0,
			expectedLength:         0,
			expectedValidWKT:       "POINT (0.0000000000000000 0.0000000000000000)",
		},
		{
			name:                  "point_empty",
			wkt:                   "POINT EMPTY",
			expectedBounds:        &geos.Box2D{MinX: math.Inf(1), MinY: math.Inf(1), MaxX: math.Inf(-1), MaxY: math.Inf(-1)},
			expectedEmpty:         true,
			expectedEnvelopeWKT:   "POINT EMPTY",
			expectedNumGeometries: 1,
			expectedSRID:          0,
			expectedType:          "Point",
			expectedTypeID:        geos.TypeIDPoint,
			expectedArea:          0,
			expectedLength:        0,
			expectedValidWKT:      "POINT EMPTY",
		},
		{
			name:                   "linestring",
			wkt:                    "LINESTRING (0.0000000000000000 0.0000000000000000, 1.0000000000000000 1.0000000000000000)",
			expectedBounds:         &geos.Box2D{MinX: 0, MinY: 0, MaxX: 1, MaxY: 1},
			expectedEmpty:          false,
			expectedEnvelopeWKT:    "POLYGON ((0 0, 1 0, 1 1, 0 1, 0 0))",
			expectedNumCoordinates: 2,
			expectedNumGeometries:  1,
			expectedSRID:           0,
			expectedType:           "LineString",
			expectedTypeID:         geos.TypeIDLineString,
			expectedArea:           0,
			expectedLength:         math.Sqrt(2),
			expectedValidWKT:       "LINESTRING (0.0000000000000000 0.0000000000000000, 1.0000000000000000 1.0000000000000000)",
		},
		{
			name:                  "linestring_empty",
			wkt:                   "LINESTRING EMPTY",
			expectedBounds:        &geos.Box2D{MinX: math.Inf(1), MinY: math.Inf(1), MaxX: math.Inf(-1), MaxY: math.Inf(-1)},
			expectedEmpty:         true,
			expectedEnvelopeWKT:   "POLYGON EMPTY",
			expectedNumGeometries: 1,
			expectedSRID:          0,
			expectedType:          "LineString",
			expectedTypeID:        geos.TypeIDLineString,
			expectedArea:          0,
			expectedLength:        0,
			expectedValidWKT:      "LINESTRING EMPTY",
		},
		{
			name:                   "polygon",
			wkt:                    "POLYGON ((0 0, 1 0, 1 1, 0 0))",
			expectedBounds:         &geos.Box2D{MinX: 0, MinY: 0, MaxX: 1, MaxY: 1},
			expectedEmpty:          false,
			expectedEnvelopeWKT:    "POLYGON ((0 0, 1 0, 1 1, 0 1, 0 0))",
			expectedNumCoordinates: 4,
			expectedNumGeometries:  1,
			expectedSRID:           0,
			expectedType:           "Polygon",
			expectedTypeID:         geos.TypeIDPolygon,
			expectedArea:           0.5,
			expectedLength:         math.Sqrt(2) + 2,
			expectedValidWKT:       "POLYGON ((0 0, 1 0, 1 1, 0 0))",
		},
		{
			name:                                   "polygon_empty",
			wkt:                                    "POLYGON EMPTY",
			expectedBounds:                         &geos.Box2D{MinX: math.Inf(1), MinY: math.Inf(1), MaxX: math.Inf(-1), MaxY: math.Inf(-1)},
			expectedEmpty:                          true,
			expectedEnvelopeWKT:                    "POLYGON EMPTY",
			expectedNumGeometries:                  1,
			expectedSRID:                           0,
			expectedType:                           "Polygon",
			expectedTypeID:                         geos.TypeIDPolygon,
			expectedArea:                           0,
			expectedLength:                         0,
			expectedValidWKT:                       "POLYGON EMPTY",
			expectedValidWKTStructureKeepCollapsed: "LINESTRING EMPTY",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := geos.NewContext()
			g := mustNewGeomFromWKT(t, c, tc.wkt)
			v := mustNewGeomFromWKT(t, c, tc.expectedValidWKT)
			assert.Equal(t, tc.expectedBounds, g.Bounds())
			assert.Equal(t, tc.expectedEmpty, g.IsEmpty())
			expectedEnvelope := mustNewGeomFromWKT(t, c, tc.expectedEnvelopeWKT)
			assert.True(t, expectedEnvelope.Equals(g.Envelope()))
			assert.Equal(t, tc.expectedNumCoordinates, g.NumCoordinates())
			assert.Equal(t, tc.expectedNumGeometries, g.NumGeometries())
			assert.True(t, g.IsSimple())
			assert.Equal(t, tc.expectedSRID, g.SRID())
			assert.Equal(t, tc.expectedType, g.Type())
			assert.Equal(t, tc.expectedTypeID, g.TypeID())
			assert.True(t, g.Clone().Equals(g))
			assert.Equal(t, !g.IsEmpty(), g.Contains(g))
			//nolint:gocritic
			assert.True(t, g.Equals(g))
			assert.True(t, g.Geometry(0).Equals(g))
			assert.Equal(t, !g.IsEmpty(), g.Intersects(g))
			assert.True(t, g.IsValid())
			assert.Equal(t, "Valid Geometry", g.IsValidReason())
			g.SetSRID(4326)
			assert.Equal(t, 4326, g.SRID())
			assert.Equal(t, tc.expectedArea, g.Area())
			assert.Equal(t, tc.expectedLength, g.Length())
			assert.Equal(t, v, g.MakeValidWithParams(geos.MakeValidLinework, geos.MakeValidDiscardCollapsed))
			assert.Equal(t, v, g.MakeValidWithParams(geos.MakeValidStructure, geos.MakeValidDiscardCollapsed))
			assert.Equal(t, v, g.MakeValidWithParams(geos.MakeValidLinework, geos.MakeValidKeepCollapsed))
			var expectedValidStructureKeepCollapsed *geos.Geom
			if tc.expectedValidWKTStructureKeepCollapsed != "" {
				expectedValidStructureKeepCollapsed = mustNewGeomFromWKT(t, c, tc.expectedValidWKTStructureKeepCollapsed)
			} else {
				expectedValidStructureKeepCollapsed = v
			}
			assert.Equal(t, expectedValidStructureKeepCollapsed, g.MakeValidWithParams(geos.MakeValidStructure, geos.MakeValidKeepCollapsed))
		})
	}
}

func TestGeomMethods(t *testing.T) {
	defer runtime.GC() // Exercise finalizers.
	c := geos.NewContext()
	unitSquare := mustNewGeomFromWKT(t, c, "POLYGON ((0 0, 1 0, 1 1, 0 1, 0 0))")
	northSouthLine := mustNewGeomFromWKT(t, c, "LINESTRING (0.5 0, 0.5 1)")
	eastWestLine := mustNewGeomFromWKT(t, c, "LINESTRING (0 0.5, 1 0.5)")
	northWestSquare := mustNewGeomFromWKT(t, c, "POLYGON ((0 0.5, 0 1, 0.5 1, 0.5 0.5, 0 0.5))")
	southEastSquare := mustNewGeomFromWKT(t, c, "POLYGON ((0.5 0, 0.5 0.5, 1 0.5, 1 0, 0.5 0))")
	middleSquare := mustNewGeomFromWKT(t, c, "POLYGON ((0.25 0.25, 0.25 0.75, 0.75 0.75, 0.75 0.25, 0.25 0.25))")
	assert.True(t, unitSquare.Equals(unitSquare.ConvexHull()))
	assert.True(t, unitSquare.Contains(middleSquare))
	assert.False(t, unitSquare.Contains(mustNewGeomFromWKT(t, c, "POINT (-0.5 -0.5)")))
	assert.False(t, unitSquare.CoveredBy(middleSquare))
	assert.True(t, middleSquare.CoveredBy(unitSquare))
	assert.True(t, unitSquare.Covers(middleSquare))
	assert.False(t, middleSquare.Covers(unitSquare))
	assert.True(t, northSouthLine.Crosses(eastWestLine))
	assert.False(t, northSouthLine.Crosses(mustNewGeomFromWKT(t, c, "LINESTRING (0 0, 0 1)")))
	assert.False(t, northSouthLine.Disjoint(eastWestLine))
	assert.True(t, southEastSquare.Disjoint(mustNewGeomFromWKT(t, c, "LINESTRING (0 0, 0 1)")))
	assert.Equal(t, unitSquare.Distance(unitSquare), 0.)
	assert.Equal(t, unitSquare.Distance(mustNewGeomFromWKT(t, c, "POLYGON ((2 0, 3 0, 3 1, 2 1, 2 0))")), 1.)
	assert.Equal(t, unitSquare.DistanceIndexed(mustNewGeomFromWKT(t, c, "POLYGON ((2 0, 3 0, 3 1, 2 1, 2 0))")), 1.)
	assert.True(t, unitSquare.DistanceWithin(mustNewGeomFromWKT(t, c, "POINT (2 2)"), 2))
	assert.False(t, unitSquare.DistanceWithin(mustNewGeomFromWKT(t, c, "POINT (2 2)"), 1))
	assert.True(t, middleSquare.Equals(unitSquare.Intersection(middleSquare)))
	assert.True(t, unitSquare.EqualsExact(unitSquare, 0.125))
	assert.Equal(t, unitSquare.FrechetDistance(unitSquare), 0.)
	assert.Equal(t, mustNewGeomFromWKT(t, c, "LINESTRING (0 1, 0 0)").FrechetDistance(mustNewGeomFromWKT(t, c, "LINESTRING (0 0, 0 1)")), 1.)
	assert.Equal(t, unitSquare.FrechetDistance(mustNewGeomFromWKT(t, c, "LINESTRING (0 0, 0 1)")), 1.)
	assert.Equal(t, unitSquare.FrechetDistanceDensify(unitSquare, 0.1), 0.)
	assert.Equal(t, unitSquare.HausdorffDistance(unitSquare), 0.)
	assert.Equal(t, unitSquare.HausdorffDistance(mustNewGeomFromWKT(t, c, "LINESTRING (0 0, 0 1)")), 1.)
	assert.Equal(t, unitSquare.HausdorffDistanceDensify(mustNewGeomFromWKT(t, c, "LINESTRING (0 0, 0 1)"), 0.01), 1.)
	assert.Equal(t, eastWestLine.ProjectNormalized(mustNewGeomFromWKT(t, c, "Point (0.5 0.5)")), 0.5)
	assert.Equal(t, eastWestLine.Project(mustNewGeomFromWKT(t, c, "Point (0.5 0.5)")), 0.5)
	assert.True(t, northSouthLine.Intersects(eastWestLine))
	assert.False(t, southEastSquare.Intersects(mustNewGeomFromWKT(t, c, "LINESTRING (0 0, 0 1)")))
	assert.Equal(t, [][]float64{{1, 1}, {2, 2}}, unitSquare.NearestPoints(mustNewGeomFromWKT(t, c, "POLYGON ((2 2, 3 2, 3 3, 2 3, 2 2))")))
	assert.Equal(t, nil, unitSquare.NearestPoints(mustNewGeomFromWKT(t, c, "GEOMETRYCOLLECTION EMPTY")))
	assert.True(t, middleSquare.Overlaps(southEastSquare))
	assert.False(t, northWestSquare.Overlaps(southEastSquare))
	assert.True(t, eastWestLine.Touches(southEastSquare))
	assert.False(t, southEastSquare.Touches(mustNewGeomFromWKT(t, c, "LINESTRING (0 0, 0 1)")))
	assert.True(t, middleSquare.Within(unitSquare))
	assert.False(t, unitSquare.Within(middleSquare))
	assert.Equal(t, 1.0, northSouthLine.Buffer(0.5, 4).MinimumWidth().Length())
	assert.Equal(t, 3, northSouthLine.Densify(0.5).NumPoints())
}

func TestPointMethods(t *testing.T) {
	for _, tc := range []struct {
		name                   string
		wkt                    string
		expectedCoordSeqCoords [][]float64
		expectedX              float64
		expectedY              float64
	}{
		{
			name:                   "point",
			wkt:                    "POINT (1 2)",
			expectedCoordSeqCoords: [][]float64{{1, 2}},
			expectedX:              1,
			expectedY:              2,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			defer runtime.GC() // Exercise finalizers.
			c := geos.NewContext()
			g := mustNewGeomFromWKT(t, c, tc.wkt)
			assert.Equal(t, tc.expectedCoordSeqCoords, g.CoordSeq().ToCoords())
			assert.Equal(t, tc.expectedX, g.X())
			assert.Equal(t, tc.expectedY, g.Y())
		})
	}
}

func TestLineStringMethods(t *testing.T) {
	for _, tc := range []struct {
		name                   string
		wkt                    string
		expectedClosed         bool
		expectedCoordSeqCoords [][]float64
		expectedPointWKTs      []string
		expectedRing           bool
	}{
		{
			name:                   "linestring",
			wkt:                    "LINESTRING (0 1, 2 3)",
			expectedClosed:         false,
			expectedCoordSeqCoords: [][]float64{{0, 1}, {2, 3}},
			expectedPointWKTs: []string{
				"POINT (0 1)",
				"POINT (2 3)",
			},
			expectedRing: false,
		},
		{
			name:                   "linearring",
			wkt:                    "LINEARRING (0 0, 1 0, 1 1, 0 0)",
			expectedClosed:         true,
			expectedCoordSeqCoords: [][]float64{{0, 0}, {1, 0}, {1, 1}, {0, 0}},
			expectedPointWKTs: []string{
				"POINT (0 0)",
				"POINT (1 0)",
				"POINT (1 1)",
				"POINT (0 0)",
			},
			expectedRing: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			defer runtime.GC() // Exercise finalizers.
			c := geos.NewContext()
			g := mustNewGeomFromWKT(t, c, tc.wkt)
			assert.Equal(t, tc.expectedClosed, g.IsClosed())
			assert.Equal(t, len(tc.expectedPointWKTs), g.NumPoints())
			for i, expectedPointWKT := range tc.expectedPointWKTs {
				expectedPoint := mustNewGeomFromWKT(t, c, expectedPointWKT)
				assert.True(t, expectedPoint.Equals(g.Point(i)))
			}
			assert.Equal(t, tc.expectedRing, g.IsRing())
			assert.Equal(t, tc.expectedCoordSeqCoords, g.CoordSeq().ToCoords())
		})
	}
}

func TestPolygonMethods(t *testing.T) {
	defer runtime.GC() // Exercise finalizers.
	c := geos.NewContext()
	polygon := mustNewGeomFromWKT(t, c, "POLYGON ((0 0, 3 0, 3 3, 0 3, 0 0), (1 1, 1 2, 2 2, 2 1, 1 1))")
	assert.Equal(t, nil, polygon.CoordSeq())
	assert.Equal(t, 1, polygon.NumInteriorRings())
	exteriorRing := polygon.ExteriorRing()
	expectedExteriorRing := mustNewGeomFromWKT(t, c, "LINEARRING (0 0, 3 0, 3 3, 0 3, 0 0)")
	assert.True(t, expectedExteriorRing.Equals(exteriorRing))
	assert.NotEqual(t, nil, exteriorRing.CoordSeq())
	interiorRing := polygon.InteriorRing(0)
	expectedInteriorRing := mustNewGeomFromWKT(t, c, "LINEARRING (1 1, 1 2, 2 2, 2 1, 1 1)")
	assert.True(t, expectedInteriorRing.Equals(interiorRing))
	assert.NotEqual(t, nil, interiorRing.CoordSeq())
}

func TestPolygonUnion(t *testing.T) {
	for _, tc := range []struct {
		name        string
		wkt         string
		expectedWKT string
	}{
		{
			name:        "only_one",
			wkt:         "GEOMETRYCOLLECTION (POLYGON ((0 0,1 0,1 1,0 1,0 0)))",
			expectedWKT: "POLYGON ((0 0,1 0,1 1,0 1,0 0))",
		},
		{
			name:        "two_identical_polygons",
			wkt:         "GEOMETRYCOLLECTION (POLYGON ((0 0,1 0,1 1,0 1,0 0)), POLYGON ((0 0,1 0,1 1,0 1,0 0)))",
			expectedWKT: "POLYGON ((0 0,1 0,1 1,0 1,0 0))",
		},
		{
			name:        "two_disjoint_polygons",
			wkt:         "GEOMETRYCOLLECTION (POLYGON ((0 0,1 0,1 1,0 1,0 0)), POLYGON ((10 10,11 10,11 11,10 11,10 10)))",
			expectedWKT: "MULTIPOLYGON(((0 0,1 0,1 1,0 1,0 0)), ((10 10,11 10,11 11,10 11,10 10)))",
		},
		{
			name:        "two_intersecting_polygons",
			wkt:         "GEOMETRYCOLLECTION (POLYGON ((0 0,10 0,10 10,0 10,0 0)), POLYGON ((5 0,15 0,15 10,5 10,5 0)))",
			expectedWKT: "POLYGON((0 0,5 0,10 0,15 0,15 10,10 10,5 10,0 10,0 0))",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			defer runtime.GC() // Exercise finalizers.
			c := geos.NewContext()
			polygons := mustNewGeomFromWKT(t, c, tc.wkt)
			expectedUnion := mustNewGeomFromWKT(t, c, tc.expectedWKT)

			// Test different union methods
			union1 := polygons.UnaryUnion()
			assert.True(t, expectedUnion.Equals(union1))

			if geos.VersionCompare(3, 12, 0) >= 0 {
				union2 := polygons.DisjointSubsetUnion()
				assert.True(t, expectedUnion.Equals(union2))
			}
		})
	}
}

func TestGeometryPanics(t *testing.T) {
	defer runtime.GC() // Exercise finalizers.
	c := geos.NewContext()
	assert.Panics(t, func() { c.NewEmptyLineString().Point(-1) })
	assert.Panics(t, func() { c.NewEmptyLineString().Point(0) })
	assert.NotPanics(t, func() { c.NewEmptyPolygon().ExteriorRing() })
	assert.Panics(t, func() { c.NewEmptyPolygon().InteriorRing(-1) })
	assert.Panics(t, func() { c.NewEmptyPolygon().InteriorRing(0) })
}

func TestBinaryMethods(t *testing.T) {
	defer runtime.GC() // Exercise finalizers.
	c := geos.NewContext()
	multiPoint1 := mustNewGeomFromWKT(t, c, "MULTIPOINT (0 0,1 1)")
	multiPoint2 := mustNewGeomFromWKT(t, c, "MULTIPOINT (1 1,2 2)")
	difference := multiPoint1.Difference(multiPoint2)
	assert.True(t, mustNewGeomFromWKT(t, c, "POINT (0 0)").Equals(difference))
}

func TestGeomInterpolate(t *testing.T) {
	defer runtime.GC() // Exercise finalizers.
	c := geos.NewContext()

	lineString := mustNewGeomFromWKT(t, c, "LINESTRING (0 0,1 0)")
	assert.True(t, mustNewGeomFromWKT(t, c, "POINT (0.5 0)").Equals(lineString.Interpolate(0.5)))

	point := mustNewGeomFromWKT(t, c, "POINT (0 0)")
	assert.Equal(t, nil, point.Interpolate(0.5))
}

func TestGeomPolygonizeFull(t *testing.T) {
	for _, tc := range []struct {
		name                    string
		wkt                     string
		expectedWKT             string
		expectedCutsWKT         string
		expectedDanglesWKT      string
		expectedInvalidRingsWKT string
	}{
		{
			name:                    "empty",
			wkt:                     "GEOMETRYCOLLECTION EMPTY",
			expectedWKT:             "GEOMETRYCOLLECTION EMPTY",
			expectedCutsWKT:         "GEOMETRYCOLLECTION EMPTY",
			expectedDanglesWKT:      "GEOMETRYCOLLECTION EMPTY",
			expectedInvalidRingsWKT: "GEOMETRYCOLLECTION EMPTY",
		},
		{
			name:                    "simple",
			wkt:                     "MULTILINESTRING ((0 0,1 0,1 1),(1 1,0 1,0 0))",
			expectedWKT:             "GEOMETRYCOLLECTION (POLYGON ((0 0,1 0,1 1,0 1,0 0)))",
			expectedCutsWKT:         "GEOMETRYCOLLECTION EMPTY",
			expectedDanglesWKT:      "GEOMETRYCOLLECTION EMPTY",
			expectedInvalidRingsWKT: "GEOMETRYCOLLECTION EMPTY",
		},
		{
			name:                    "dangle",
			wkt:                     "MULTILINESTRING ((0 0,1 0,1 1),(1 1,0 1,0 0),(0 0,0 -1))",
			expectedWKT:             "GEOMETRYCOLLECTION (POLYGON ((0 0,1 0,1 1,0 1,0 0)))",
			expectedCutsWKT:         "GEOMETRYCOLLECTION EMPTY",
			expectedDanglesWKT:      "GEOMETRYCOLLECTION (LINESTRING (0 0,0 -1))",
			expectedInvalidRingsWKT: "GEOMETRYCOLLECTION EMPTY",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := geos.NewContext()
			g := mustNewGeomFromWKT(t, c, tc.wkt)
			actual, cuts, dangles, invalidRings := g.PolygonizeFull()
			assert.Equal(t, mustNewGeomFromWKT(t, c, tc.expectedWKT), actual)
			assert.Equal(t, mustNewGeomFromWKT(t, c, tc.expectedCutsWKT), cuts)
			assert.Equal(t, mustNewGeomFromWKT(t, c, tc.expectedDanglesWKT), dangles)
			assert.Equal(t, mustNewGeomFromWKT(t, c, tc.expectedInvalidRingsWKT), invalidRings)
		})
	}
}

func TestNewGeomFromGeoJSON(t *testing.T) {
	for i, tc := range []struct {
		geoJSON     string
		expectedWKT string
	}{
		{
			geoJSON:     `{"type":"Point","coordinates":[1,2]}`,
			expectedWKT: "POINT (1 2)",
		},
		{
			geoJSON:     `{"type":"LineString","coordinates":[[1,2],[3,4]]}`,
			expectedWKT: "LINESTRING (1 2, 3 4)",
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			defer runtime.GC() // Exercise finalizers.
			context := geos.NewContext()
			actualGeom, err := context.NewGeomFromGeoJSON(tc.geoJSON)
			assert.NoError(t, err)
			assert.True(t, mustNewGeomFromWKT(t, context, tc.expectedWKT).Equals(actualGeom))
		})
	}
}

func TestNewGeomFromGeoJSONError(t *testing.T) {
	_, err := geos.NewContext().NewGeomFromGeoJSON(`{"type":`)
	assert.Error(t, err)
}

func TestGeomNearestPointsAliasing(t *testing.T) {
	c := geos.NewContext()
	geom1 := mustNewGeomFromWKT(t, c, "POINT (0 1)")
	geom2 := mustNewGeomFromWKT(t, c, "POINT (2 3)")
	points := geom1.NearestPoints(geom2)
	points[0] = append(points[0], 4)
	assert.Equal(t, []float64{2, 3}, points[1])
}

func TestGeomToJSON(t *testing.T) {
	geom := mustNewGeomFromWKT(t, geos.NewContext(), "POINT (1 2)")
	assert.Equal(t, `{"type":"Point","coordinates":[1.0,2.0]}`, geom.ToGeoJSON(-1))
}

func TestWKBError(t *testing.T) {
	_, err := geos.NewContext().NewGeomFromWKT("POINT (0 0")
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "ParseException: Expected word but encountered end of stream")
}

func TestWKTError(t *testing.T) {
	_, err := geos.NewContext().NewGeomFromWKB(nil)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "ParseException: Unexpected EOF parsing WKB")
}

func TestWKXRoundTrip(t *testing.T) {
	for _, tc := range []struct {
		name       string
		wkt        string
		wktPre3_12 string
	}{
		{
			name:       "point",
			wkt:        "POINT (0 0)",
			wktPre3_12: "POINT (0.0000000000000000 0.0000000000000000)",
		},
		{
			name:       "line_string",
			wkt:        "LINESTRING (0 0, 1 0)",
			wktPre3_12: "LINESTRING (0.0000000000000000 0.0000000000000000, 1.0000000000000000 0.0000000000000000)",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			defer runtime.GC() // Exercise finalizers.
			c := geos.NewContext()
			wkt := tc.wkt
			if geos.VersionCompare(3, 12, 0) < 0 {
				wkt = tc.wktPre3_12
			}
			wktGeom := mustNewGeomFromWKT(t, c, wkt)
			assert.Equal(t, wkt, wktGeom.ToWKT())
			wkbGeom, err := c.NewGeomFromWKB(wktGeom.ToWKB())
			assert.NoError(t, err)
			assert.Equal(t, wkt, wkbGeom.ToWKT())
			ewkbWithSRIDGeom, err := c.NewGeomFromWKB(wktGeom.ToEWKBWithSRID())
			assert.NoError(t, err)
			assert.Equal(t, wkt, ewkbWithSRIDGeom.ToWKT())
		})
	}
}

func TestEWKBWithSRIDRoundTrip(t *testing.T) {
	c := geos.NewContext()
	for _, tc := range []struct {
		name string
		geom *geos.Geom
	}{
		{
			name: "point",
			geom: mustNewGeomFromWKT(t, c, "POINT (0 0)").SetSRID(4326),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			newG, err := c.NewGeomFromWKB(tc.geom.ToEWKBWithSRID())
			assert.NoError(t, err)
			assert.True(t, newG.Equals(tc.geom))
			assert.Equal(t, tc.geom.SRID(), newG.SRID())
		})
	}
}

func TestGeomRelate(t *testing.T) {
	c := geos.NewContext()
	g1 := mustNewGeomFromWKT(t, c, "POINT (0 0)")
	g2 := mustNewGeomFromWKT(t, c, "LINESTRING (0 0,1 0)")
	assert.Equal(t, "F0FFFF102", g1.Relate(g2))
}

func TestSetPrecision(t *testing.T) {
	g1 := mustNewGeomFromWKT(t, geos.NewContext(), "POINT (1 2)")
	g2 := g1.SetPrecision(1, geos.PrecisionRulePointwise)
	assert.Equal(t, 0., g1.Precision())
	assert.Equal(t, 1., g2.Precision())
}

func TestUserData(t *testing.T) {
	g := mustNewGeomFromWKT(t, geos.NewContext(), "POINT (0 0)")
	assert.Equal(t, uintptr(0), g.UserData())
	assert.Equal(t, g, g.SetUserData(1))
	assert.Equal(t, uintptr(1), g.UserData())
}

func TestMakeValid(t *testing.T) {
	for _, tc := range []struct {
		name                        string
		wkt                         string
		expectedWktLineworkDiscard  string
		expectedWktStructureDiscard string
		expectedWktLineworkKeep     string
		expectedWktStructureKeep    string
	}{
		{
			name:                        "LINESTRING",
			wkt:                         `LINESTRING(0 0, 0 0)`,
			expectedWktLineworkDiscard:  `POINT (0.0 0.0)`,
			expectedWktStructureDiscard: `LINESTRING EMPTY`,
			expectedWktLineworkKeep:     `POINT (0.0 0.0)`,
			expectedWktStructureKeep:    `POINT (0.0 0.0)`,
		},
		{
			name: "MULTIPOLYGON",
			wkt: `MULTIPOLYGON(((91 50,79 22,51 10,23 22,11 50,23 78,51 90,79 78,91 50)),
			((91 100,79 72,51 60,23 72,11 100,23 128,51 140,79 128,91 100)),
			((91 150,79 122,51 110,23 122,11 150,23 178,51 190,79 178,91 150)),
			((141 50,129 22,101 10,73 22,61 50,73 78,101 90,129 78,141 50)),
			((141 100,129 72,101 60,73 72,61 100,73 128,101 140,129 128,141 100)),
			((141 150,129 122,101 110,73 122,61 150,73 178,101 190,129 178,141 150)))`,
			expectedWktLineworkDiscard:  `MultiPolygon (((51 110, 68.5 117.5, 61 100, 68.5 82.5, 51 90, 23 78, 21.7142857142857153 75, 11 100, 21.7142857142857153 125, 23 122, 51 110)),((51 190, 76 179.28571428571427759, 73 178, 61 150, 68.5 132.5, 51 140, 23 128, 21.7142857142857153 125, 11 150, 23 178, 51 190)),((141 150, 130.28571428571427759 125, 129 128, 101 140, 83.5 132.5, 91 150, 79 178, 76 179.28571428571427759, 101 190, 129 178, 141 150)),((129 78, 101 90, 83.5 82.5, 91 100, 83.5 117.5, 101 110, 129 122, 130.28571428571427759 125, 141 100, 130.28571428571427759 75, 129 78)),((101 10, 76 20.7142857142857153, 79 22, 91 50, 83.5 67.5, 101 60, 129 72, 130.28571428571427759 75, 141 50, 129 22, 101 10)),((11 50, 21.7142857142857153 75, 23 72, 51 60, 68.5 67.5, 61 50, 73 22, 76 20.7142857142857153, 51 10, 23 22, 11 50)),((83.5 82.5, 80.2857142857142918 75, 79 78, 76 79.2857142857142918, 83.5 82.5)),((83.5 67.5, 76 70.7142857142857082, 79 72, 80.2857142857142918 75, 83.5 67.5)),((68.5 67.5, 71.7142857142857082 75, 73 72, 76 70.7142857142857082, 68.5 67.5)),((68.5 82.5, 76 79.2857142857142918, 73 78, 71.7142857142857082 75, 68.5 82.5)),((83.5 132.5, 80.2857142857142918 125, 79 128, 76 129.28571428571427759, 83.5 132.5)),((83.5 117.5, 76 120.7142857142857082, 79 122, 80.2857142857142918 125, 83.5 117.5)),((68.5 117.5, 71.7142857142857082 125, 73 122, 76 120.7142857142857082, 68.5 117.5)),((68.5 132.5, 76 129.28571428571427759, 73 128, 71.7142857142857082 125, 68.5 132.5)))`,
			expectedWktStructureDiscard: `Polygon ((23 22, 11 50, 21.7142857142857153 75, 11 100, 21.7142857142857153 125, 11 150, 23 178, 51 190, 76 179.28571428571427759, 101 190, 129 178, 141 150, 130.28571428571427759 125, 141 100, 130.28571428571427759 75, 141 50, 129 22, 101 10, 76 20.7142857142857153, 51 10, 23 22))`,
			expectedWktLineworkKeep:     `MultiPolygon (((51 110, 68.5 117.5, 61 100, 68.5 82.5, 51 90, 23 78, 21.7142857142857153 75, 11 100, 21.7142857142857153 125, 23 122, 51 110)),((51 190, 76 179.28571428571427759, 73 178, 61 150, 68.5 132.5, 51 140, 23 128, 21.7142857142857153 125, 11 150, 23 178, 51 190)),((141 150, 130.28571428571427759 125, 129 128, 101 140, 83.5 132.5, 91 150, 79 178, 76 179.28571428571427759, 101 190, 129 178, 141 150)),((129 78, 101 90, 83.5 82.5, 91 100, 83.5 117.5, 101 110, 129 122, 130.28571428571427759 125, 141 100, 130.28571428571427759 75, 129 78)),((101 10, 76 20.7142857142857153, 79 22, 91 50, 83.5 67.5, 101 60, 129 72, 130.28571428571427759 75, 141 50, 129 22, 101 10)),((11 50, 21.7142857142857153 75, 23 72, 51 60, 68.5 67.5, 61 50, 73 22, 76 20.7142857142857153, 51 10, 23 22, 11 50)),((83.5 82.5, 80.2857142857142918 75, 79 78, 76 79.2857142857142918, 83.5 82.5)),((83.5 67.5, 76 70.7142857142857082, 79 72, 80.2857142857142918 75, 83.5 67.5)),((68.5 67.5, 71.7142857142857082 75, 73 72, 76 70.7142857142857082, 68.5 67.5)),((68.5 82.5, 76 79.2857142857142918, 73 78, 71.7142857142857082 75, 68.5 82.5)),((83.5 132.5, 80.2857142857142918 125, 79 128, 76 129.28571428571427759, 83.5 132.5)),((83.5 117.5, 76 120.7142857142857082, 79 122, 80.2857142857142918 125, 83.5 117.5)),((68.5 117.5, 71.7142857142857082 125, 73 122, 76 120.7142857142857082, 68.5 117.5)),((68.5 132.5, 76 129.28571428571427759, 73 128, 71.7142857142857082 125, 68.5 132.5)))`,
			expectedWktStructureKeep:    `Polygon ((23 22, 11 50, 21.7142857142857153 75, 11 100, 21.7142857142857153 125, 11 150, 23 178, 51 190, 76 179.28571428571427759, 101 190, 129 178, 141 150, 130.28571428571427759 125, 141 100, 130.28571428571427759 75, 141 50, 129 22, 101 10, 76 20.7142857142857153, 51 10, 23 22))`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := geos.NewContext()
			g := newInvalidGeomFromWKT(t, c, tc.wkt)
			v1 := mustNewGeomFromWKT(t, c, tc.expectedWktLineworkDiscard)
			v2 := mustNewGeomFromWKT(t, c, tc.expectedWktStructureDiscard)
			v3 := mustNewGeomFromWKT(t, c, tc.expectedWktLineworkKeep)
			v4 := mustNewGeomFromWKT(t, c, tc.expectedWktStructureKeep)
			assert.Equal(t, v1, g.MakeValidWithParams(geos.MakeValidLinework, geos.MakeValidDiscardCollapsed))
			assert.Equal(t, v2, g.MakeValidWithParams(geos.MakeValidStructure, geos.MakeValidDiscardCollapsed))
			assert.Equal(t, v3, g.MakeValidWithParams(geos.MakeValidLinework, geos.MakeValidKeepCollapsed))
			assert.Equal(t, v4, g.MakeValidWithParams(geos.MakeValidStructure, geos.MakeValidKeepCollapsed))
		})
	}
}
