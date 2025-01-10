package geometry_test

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"runtime"
	"strings"
	"testing"

	"github.com/alecthomas/assert/v2"

	"github.com/twpayne/go-geos"
	"github.com/twpayne/go-geos/geometry"
)

var (
	_ driver.Value     = &geometry.Geometry{}
	_ gob.GobEncoder   = &geometry.Geometry{}
	_ gob.GobDecoder   = &geometry.Geometry{}
	_ json.Marshaler   = &geometry.Geometry{}
	_ json.Unmarshaler = &geometry.Geometry{}
	_ sql.Scanner      = &geometry.Geometry{}
	_ xml.Marshaler    = &geometry.Geometry{}
)

func TestGeometry(t *testing.T) {
	for _, tc := range []struct {
		name                 string
		geometry             *geometry.Geometry
		skipGeoJSON          bool
		expectedGeoJSONError bool
		expectedKML          string
	}{
		{
			name:                 "point_empty",
			geometry:             mustNewGeometryFromWKT(t, "POINT EMPTY"),
			expectedGeoJSONError: true,
		},
		{
			name:        "point",
			geometry:    mustNewGeometryFromWKT(t, "POINT (0 1)"),
			expectedKML: "<Point><coordinates>0,1</coordinates></Point>",
		},
		{
			name:                 "linestring_empty",
			geometry:             mustNewGeometryFromWKT(t, "LINESTRING EMPTY"),
			expectedGeoJSONError: true,
			expectedKML:          "<LineString></LineString>",
		},
		{
			name:        "linestring",
			geometry:    mustNewGeometryFromWKT(t, "LINESTRING (0 1, 2 3)"),
			expectedKML: "<LineString><coordinates>0,1 2,3</coordinates></LineString>",
		},
		{
			name:                 "linearring_empty",
			geometry:             mustNewGeometryFromWKT(t, "LINEARRING EMPTY"),
			expectedGeoJSONError: true,
			expectedKML:          "<LinearRing></LinearRing>",
		},
		{
			name:                 "linearring",
			geometry:             mustNewGeometryFromWKT(t, "LINEARRING (0 0, 1 0, 1 1, 0 0)"),
			expectedGeoJSONError: true,
			expectedKML:          "<LinearRing><coordinates>0,0 1,0 1,1 0,0</coordinates></LinearRing>",
		},
		{
			name:                 "polygon_empty",
			geometry:             mustNewGeometryFromWKT(t, "POLYGON EMPTY"),
			expectedGeoJSONError: true,
			expectedKML:          "<Polygon></Polygon>",
		},
		{
			name:        "polygon",
			geometry:    mustNewGeometryFromWKT(t, "POLYGON ((0 0, 1 0, 1 1, 0 0))"),
			expectedKML: "<Polygon><outerBoundaryIs><LinearRing><coordinates>0,0 1,0 1,1 0,0</coordinates></LinearRing></outerBoundaryIs></Polygon>",
		},
		{
			name:     "polygon_interior_rings",
			geometry: mustNewGeometryFromWKT(t, "POLYGON ((0 0, 3 0, 3 3, 0 3, 0 0), (1 1, 1 2, 2 2, 2 1, 1 1))"),
			expectedKML: "" +
				"<Polygon>" +
				"<outerBoundaryIs><LinearRing><coordinates>0,0 3,0 3,3 0,3 0,0</coordinates></LinearRing></outerBoundaryIs>" +
				"<innerBoundaryIs><LinearRing><coordinates>1,1 1,2 2,2 2,1 1,1</coordinates></LinearRing></innerBoundaryIs>" +
				"</Polygon>",
		},
		{
			name:        "multipoint_empty",
			geometry:    mustNewGeometryFromWKT(t, "MULTIPOINT EMPTY"),
			skipGeoJSON: true, // FIXME
			expectedKML: "<MultiGeometry></MultiGeometry>",
		},
		{
			name:     "multipoint",
			geometry: mustNewGeometryFromWKT(t, "MULTIPOINT (0 1, 2 3)"),
			expectedKML: "" +
				"<MultiGeometry>" +
				"<Point><coordinates>0,1</coordinates></Point>" +
				"<Point><coordinates>2,3</coordinates></Point>" +
				"</MultiGeometry>",
		},
		{
			name:        "multilinestring_empty",
			geometry:    mustNewGeometryFromWKT(t, "MULTILINESTRING EMPTY"),
			expectedKML: "<MultiGeometry></MultiGeometry>",
		},
		{
			name:     "multilinestring",
			geometry: mustNewGeometryFromWKT(t, "MULTILINESTRING ((0 1, 2 3), (4 5, 6 7))"),
			expectedKML: "" +
				"<MultiGeometry>" +
				"<LineString><coordinates>0,1 2,3</coordinates></LineString>" +
				"<LineString><coordinates>4,5 6,7</coordinates></LineString>" +
				"</MultiGeometry>",
		},
		{
			name:        "multipolygon_empty",
			geometry:    mustNewGeometryFromWKT(t, "MULTIPOLYGON EMPTY"),
			expectedKML: "<MultiGeometry></MultiGeometry>",
		},
		{
			name:     "multipolygon",
			geometry: mustNewGeometryFromWKT(t, "MULTIPOLYGON (((-1 -1, 0 -1, 0 0, -1 -1)), ((0 0, 3 0, 3 3, 0 3, 0 0), (1 1, 1 2, 2 2, 2 1, 1 1)))"),
			expectedKML: "" +
				"<MultiGeometry>" +
				"<Polygon>" +
				"<outerBoundaryIs><LinearRing><coordinates>-1,-1 0,-1 0,0 -1,-1</coordinates></LinearRing></outerBoundaryIs>" +
				"</Polygon>" +
				"<Polygon>" +
				"<outerBoundaryIs><LinearRing><coordinates>0,0 3,0 3,3 0,3 0,0</coordinates></LinearRing></outerBoundaryIs>" +
				"<innerBoundaryIs><LinearRing><coordinates>1,1 1,2 2,2 2,1 1,1</coordinates></LinearRing></innerBoundaryIs>" +
				"</Polygon>" +
				"</MultiGeometry>",
		},
		{
			name:        "geometrycollection_empty",
			geometry:    mustNewGeometryFromWKT(t, "GEOMETRYCOLLECTION EMPTY"),
			skipGeoJSON: true, // FIXME
			expectedKML: "<MultiGeometry></MultiGeometry>",
		},
		// FIXME geometrycollection
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Run("gob", func(t *testing.T) {
				defer runtime.GC() // Exercise finalizers.
				data := &bytes.Buffer{}
				assert.NoError(t, gob.NewEncoder(data).Encode(tc.geometry))
				var actualG geometry.Geometry
				assert.NoError(t, gob.NewDecoder(data).Decode(&actualG))
				assert.True(t, actualG.Equals(tc.geometry.Geom))
			})
			t.Run("geojson", func(t *testing.T) {
				defer runtime.GC() // Exercise finalizers.
				if tc.skipGeoJSON {
					t.Skip()
				}
				geoJSON, err := tc.geometry.AsGeoJSON()
				if tc.expectedGeoJSONError {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					actualG, err := geometry.NewGeometryFromGeoJSON(geoJSON)
					assert.NoError(t, err)
					assert.True(t, actualG.Equals(tc.geometry.Geom))
				}
			})
			if tc.expectedKML != "" {
				t.Run("kml", func(t *testing.T) {
					defer runtime.GC() // Exercise finalizers.
					data := &strings.Builder{}
					assert.NoError(t, xml.NewEncoder(data).Encode(tc.geometry))
					assert.Equal(t, tc.expectedKML, data.String())
				})
			}
			t.Run("sql", func(t *testing.T) {
				defer runtime.GC() // Exercise finalizers.
				value, err := tc.geometry.Value()
				assert.NoError(t, err)
				var actualG geometry.Geometry
				assert.NoError(t, actualG.Scan(value))
				assert.True(t, actualG.Equals(tc.geometry.Geom))
			})
		})
	}
}

func TestNewGeometry(t *testing.T) {
	expected := geometry.NewGeometry(geos.NewPoint([]float64{1, 2}))

	actual, err := geometry.NewGeometryFromGeoJSON([]byte(`{"type":"Point","coordinates":[1,2]}`))
	assert.NoError(t, err)
	assert.Equal(t, expected, actual, assert.Exclude[*geos.Context]())

	wkb, err := hex.DecodeString("0101000000000000000000f03f0000000000000040")
	assert.NoError(t, err)
	actual, err = geometry.NewGeometryFromWKB(wkb)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)

	actual, err = geometry.NewGeometryFromWKT("POINT (1 2)")
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}
