package geometry

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/twpayne/go-geos"
)

var (
	_ json.Marshaler   = &Geometry{}
	_ json.Unmarshaler = &Geometry{}
)

func TestGeoJSON(t *testing.T) {
	for i, tc := range []struct {
		geom       *Geometry
		geoJSONStr string
	}{
		{
			geom:       NewGeometry(geos.NewPoint([]float64{1, 2})),
			geoJSONStr: `{"type":"Point","coordinates":[1,2]}`,
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			actualGeoJSON, err := tc.geom.MarshalJSON()
			require.NoError(t, err)
			assert.Equal(t, tc.geoJSONStr, string(actualGeoJSON))

			var geom Geometry
			require.NoError(t, geom.UnmarshalJSON([]byte(tc.geoJSONStr)))
			assert.True(t, tc.geom.Equals(geom.Geom))
		})
	}
}
