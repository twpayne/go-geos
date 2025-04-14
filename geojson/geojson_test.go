package geojson_test

import (
	"strconv"
	"testing"

	"github.com/alecthomas/assert/v2"

	"github.com/twpayne/go-geos"
	"github.com/twpayne/go-geos/geojson"
	"github.com/twpayne/go-geos/geometry"
)

func TestFeature(t *testing.T) {
	for i, tc := range []struct {
		feat       *geojson.Feature
		geoJSONStr string
	}{
		{
			feat: &geojson.Feature{
				ID:       "testID",
				Geometry: *geometry.NewGeometry(geos.NewPoint([]float64{1, 2})),
				Properties: map[string]any{
					"key": "value",
				},
			},
			geoJSONStr: `{"id":"testID","type":"Feature","geometry":{"type":"Point","coordinates":[1,2]},"properties":{"key":"value"}}`,
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			actualGeoJSON, err := tc.feat.MarshalJSON()
			assert.NoError(t, err)
			assert.Equal(t, tc.geoJSONStr, string(actualGeoJSON))

			var feat geojson.Feature
			assert.NoError(t, feat.UnmarshalJSON([]byte(tc.geoJSONStr)))
			assert.True(t, tc.feat.Geometry.Equals(feat.Geometry.Geom))
		})
	}
}

func TestFeatureCollection(t *testing.T) {
	for i, tc := range []struct {
		featureCollection geojson.FeatureCollection
		geoJSONStr        string
	}{
		{
			featureCollection: geojson.FeatureCollection{},
			geoJSONStr:        `{"type":"FeatureCollection","features":[]}`,
		},
		{
			featureCollection: geojson.FeatureCollection{
				{
					ID:       "point",
					Geometry: *geometry.NewGeometry(geos.NewPoint([]float64{1, 2})),
					Properties: map[string]any{
						"key": "value",
					},
				},
			},
			geoJSONStr: `{"type":"FeatureCollection","features":[{"id":"point","type":"Feature","geometry":{"type":"Point","coordinates":[1,2]},"properties":{"key":"value"}}]}`,
		},
		{
			featureCollection: geojson.FeatureCollection{
				{
					ID:       "point",
					Geometry: *geometry.NewGeometry(geos.NewPoint([]float64{1, 2})),
					Properties: map[string]any{
						"key": "value",
					},
				},
				{
					ID:       "linestring",
					Geometry: *geometry.NewGeometry(geos.NewLineString([][]float64{{1, 2}, {3, 4}})),
					Properties: map[string]any{
						"key": "value",
					},
				},
			},
			geoJSONStr: `{"type":"FeatureCollection","features":[{"id":"point","type":"Feature","geometry":{"type":"Point","coordinates":[1,2]},"properties":{"key":"value"}},{"id":"linestring","type":"Feature","geometry":{"type":"LineString","coordinates":[[1,2],[3,4]]},"properties":{"key":"value"}}]}`,
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			actualGeoJSON, err := tc.featureCollection.MarshalJSON()
			assert.NoError(t, err)
			assert.Equal(t, tc.geoJSONStr, string(actualGeoJSON))

			var featureCollection geojson.FeatureCollection
			assert.NoError(t, featureCollection.UnmarshalJSON([]byte(tc.geoJSONStr)))
			assert.Equal(t, tc.featureCollection, featureCollection)
		})
	}
}
