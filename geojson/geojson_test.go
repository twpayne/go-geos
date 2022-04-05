package geojson

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twpayne/go-geos"
	"github.com/twpayne/go-geos/geometry"
)

func TestGeoJSON(t *testing.T) {
	for i, tc := range []struct {
		feat       *Feature
		geoJSONStr string
	}{
		{
			feat: &Feature{
				ID:       "testID",
				Geometry: *geometry.NewGeometry(geos.NewPoint([]float64{1, 2})),
				Properties: map[string]interface{}{
					"key": "value",
				},
			},
			geoJSONStr: `{"id":"testID","type":"Feature","geometry":{"type":"Point","coordinates":[1,2]},"properties":{"key":"value"}}`,
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			actualGeoJSON, err := tc.feat.MarshalJSON()
			require.NoError(t, err)
			assert.Equal(t, tc.geoJSONStr, string(actualGeoJSON))

			var feat Feature
			require.NoError(t, feat.UnmarshalJSON([]byte(tc.geoJSONStr)))
			assert.True(t, tc.feat.Geometry.Equals(feat.Geometry.Geom))
		})
	}
}
