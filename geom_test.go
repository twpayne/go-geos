package geos

import (
	"encoding/json"
	"math"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGeometryMethods(t *testing.T) {
	for _, tc := range []struct {
		name                  string
		wkt                   string
		expectedBounds        *Bounds
		expectedEmpty         bool
		expectedEnvelopeWKT   string
		expectedNumGeometries int
		expectedSRID          int
		expectedType          string
		expectedTypeID        GeometryTypeID
		expectedLength        float64
		expectedArea          float64
	}{
		{
			name:                  "point",
			wkt:                   "POINT (0.0000000000000000 0.0000000000000000)",
			expectedBounds:        &Bounds{MinX: 0, MinY: 0, MaxX: 0, MaxY: 0},
			expectedEmpty:         false,
			expectedEnvelopeWKT:   "POINT (0 0)",
			expectedNumGeometries: 1,
			expectedSRID:          0,
			expectedType:          "Point",
			expectedTypeID:        PointTypeID,
			expectedLength:        0,
			expectedArea:          0,
		},
		{
			name:                  "point_empty",
			wkt:                   "POINT EMPTY",
			expectedBounds:        &Bounds{MinX: math.Inf(1), MinY: math.Inf(1), MaxX: math.Inf(-1), MaxY: math.Inf(-1)},
			expectedEmpty:         true,
			expectedEnvelopeWKT:   "POINT EMPTY",
			expectedNumGeometries: 1,
			expectedSRID:          0,
			expectedType:          "Point",
			expectedTypeID:        PointTypeID,
			expectedLength:        0,
			expectedArea:          0,
		},
		{
			name:                  "linestring",
			wkt:                   "LINESTRING (0.0000000000000000 0.0000000000000000, 1.0000000000000000 1.0000000000000000)",
			expectedBounds:        &Bounds{MinX: 0, MinY: 0, MaxX: 1, MaxY: 1},
			expectedEmpty:         false,
			expectedEnvelopeWKT:   "POLYGON ((0 0, 1 0, 1 1, 0 1, 0 0))",
			expectedNumGeometries: 1,
			expectedSRID:          0,
			expectedType:          "LineString",
			expectedTypeID:        LineStringTypeID,
			expectedLength:        math.Sqrt(2),
			expectedArea:          0,
		},
		{
			name:                  "linestring_empty",
			wkt:                   "LINESTRING EMPTY",
			expectedBounds:        &Bounds{MinX: math.Inf(1), MinY: math.Inf(1), MaxX: math.Inf(-1), MaxY: math.Inf(-1)},
			expectedEmpty:         true,
			expectedEnvelopeWKT:   "POLYGON EMPTY",
			expectedNumGeometries: 1,
			expectedSRID:          0,
			expectedType:          "LineString",
			expectedTypeID:        LineStringTypeID,
			expectedLength:        0,
			expectedArea:          0,
		},
		{
			name:                  "polygon",
			wkt:                   "POLYGON ((0 0, 1 0, 1 1, 0 0))",
			expectedBounds:        &Bounds{MinX: 0, MinY: 0, MaxX: 1, MaxY: 1},
			expectedEmpty:         false,
			expectedEnvelopeWKT:   "POLYGON ((0 0, 1 0, 1 1, 0 1, 0 0))",
			expectedNumGeometries: 1,
			expectedSRID:          0,
			expectedType:          "Polygon",
			expectedTypeID:        PolygonTypeID,
			expectedLength:        math.Sqrt(2) + 2,
			expectedArea:          0.5,
		},
		{
			name:                  "polygon_empty",
			wkt:                   "POLYGON EMPTY",
			expectedBounds:        &Bounds{MinX: math.Inf(1), MinY: math.Inf(1), MaxX: math.Inf(-1), MaxY: math.Inf(-1)},
			expectedEmpty:         true,
			expectedEnvelopeWKT:   "POLYGON EMPTY",
			expectedNumGeometries: 1,
			expectedSRID:          0,
			expectedType:          "Polygon",
			expectedTypeID:        PolygonTypeID,
			expectedLength:        0,
			expectedArea:          0,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := NewContext()
			g := mustNewGeomFromWKT(t, c, tc.wkt)
			assert.Equal(t, tc.expectedBounds, g.Bounds())
			assert.Equal(t, tc.expectedEmpty, g.IsEmpty())
			expectedEnvelope := mustNewGeomFromWKT(t, c, tc.expectedEnvelopeWKT)
			assert.True(t, expectedEnvelope.Equals(g.Envelope()))
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
			assert.Equal(t, tc.expectedLength, g.Length())
			assert.Equal(t, tc.expectedArea, g.Area())
		})
	}
}

func TestGeomMethods(t *testing.T) {
	defer runtime.GC() // Exercise finalizers.
	c := NewContext()
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
	assert.True(t, middleSquare.Equals(unitSquare.Intersection(middleSquare)))
	assert.True(t, unitSquare.EqualsExact(unitSquare, 0.125))
	assert.True(t, northSouthLine.Intersects(eastWestLine))
	assert.False(t, southEastSquare.Intersects(mustNewGeomFromWKT(t, c, "LINESTRING (0 0, 0 1)")))
	assert.Equal(t, [][]float64{{1, 1}, {2, 2}}, unitSquare.NearestPoints(mustNewGeomFromWKT(t, c, "POLYGON ((2 2, 3 2, 3 3, 2 3, 2 2))")))
	assert.Nil(t, unitSquare.NearestPoints(mustNewGeomFromWKT(t, c, "GEOMETRYCOLLECTION EMPTY")))
	assert.True(t, middleSquare.Overlaps(southEastSquare))
	assert.False(t, northWestSquare.Overlaps(southEastSquare))
	assert.True(t, eastWestLine.Touches(southEastSquare))
	assert.False(t, southEastSquare.Touches(mustNewGeomFromWKT(t, c, "LINESTRING (0 0, 0 1)")))
	assert.True(t, middleSquare.Within(unitSquare))
	assert.False(t, unitSquare.Within(middleSquare))
	assert.Equal(t, 3, northSouthLine.Densify(0.5).NumPoints())
	assert.Equal(t, 1.0, northSouthLine.Buffer(0.5, 4).MinimumWidth().Length())
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
			c := NewContext()
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
			c := NewContext()
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
	c := NewContext()
	polygon := mustNewGeomFromWKT(t, c, "POLYGON ((0 0, 3 0, 3 3, 0 3, 0 0), (1 1, 1 2, 2 2, 2 1, 1 1))")
	assert.Equal(t, 1, polygon.NumInteriorRings())
	expectedOuterRing := mustNewGeomFromWKT(t, c, "LINEARRING (0 0, 3 0, 3 3, 0 3, 0 0)")
	assert.True(t, expectedOuterRing.Equals(polygon.ExteriorRing()))
	expectedInnerRing := mustNewGeomFromWKT(t, c, "LINEARRING (1 1, 1 2, 2 2, 2 1, 1 1)")
	assert.True(t, expectedInnerRing.Equals(polygon.InteriorRing(0)))
}

func TestGeometryPanics(t *testing.T) {
	defer runtime.GC() // Exercise finalizers.
	c := NewContext()
	assert.Panics(t, func() { c.NewEmptyLineString().Point(-1) })
	assert.Panics(t, func() { c.NewEmptyLineString().Point(0) })
	assert.NotPanics(t, func() { c.NewEmptyPolygon().ExteriorRing() })
	assert.Panics(t, func() { c.NewEmptyPolygon().InteriorRing(-1) })
	assert.Panics(t, func() { c.NewEmptyPolygon().InteriorRing(0) })
	assert.NotPanics(t, func() {
		g := NewEmptyPoint()
		g.Destroy()
		g.Destroy()
	})
}

func TestWKBError(t *testing.T) {
	_, err := NewContext().NewGeomFromWKT("POINT (0 0")
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "ParseException: Expected word but encountered end of stream")
}

func TestWKTError(t *testing.T) {
	_, err := NewContext().NewGeomFromWKB(nil)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "ParseException: Unexpected EOF parsing WKB")
}

func TestJsonError(t *testing.T) {
	_, err := NewContext().NewGeomFromJSON("invalidjson")
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "ParseException: Error parsing JSON: '[json.exception.parse_error.101] parse error at line 1, column 1: syntax error while parsing value - invalid literal; last read: 'i''")
}

func TestWKXRoundTrip(t *testing.T) {
	for _, tc := range []struct {
		name string
		wkt  string
	}{
		{
			name: "point",
			wkt:  "POINT (0.0000000000000000 0.0000000000000000)",
		},
		{
			name: "line_string",
			wkt:  "LINESTRING (0.0000000000000000 0.0000000000000000, 1.0000000000000000 0.0000000000000000)",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			defer runtime.GC() // Exercise finalizers.
			c := NewContext()
			g := mustNewGeomFromWKT(t, c, tc.wkt)
			assert.Equal(t, tc.wkt, g.ToWKT())
			newG, err := c.NewGeomFromWKB(g.ToWKB())
			require.NoError(t, err)
			assert.Equal(t, tc.wkt, newG.ToWKT())
		})
	}
}
func TestFromJSON(t *testing.T) {
	for _, tc := range []struct {
		name string
		json string
		wkt  string
	}{
		{
			name: "point",
			json: `{ "type": "Point","coordinates": [30, 10]}`,
			wkt:  "POINT (30.0000000000000000 10.0000000000000000)",
		},
		{
			name: "line_string",
			json: `{ "type": "LineString","coordinates": [[30, 10], [10, 30]]}`,
			wkt:  "LINESTRING (30.0000000000000000 10.0000000000000000, 10.0000000000000000 30.0000000000000000)",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			defer runtime.GC() // Exercise finalizers.
			c := NewContext()
			jsonG, err := c.NewGeomFromJSON(tc.json)
			require.NoError(t, err)
			assert.Equal(t, tc.wkt, jsonG.ToWKT())
		})
	}
}

func TestUnmarshalJSON(t *testing.T) {
	var g *Geom
	jsonData := []byte(`{ "type": "Point","coordinates": [30, 10]}`)
	err := json.Unmarshal(jsonData, &g)
	require.NoError(t, err)
	assert.Equal(t, "POINT (30.0000000000000000 10.0000000000000000)", g.ToWKT())
}
