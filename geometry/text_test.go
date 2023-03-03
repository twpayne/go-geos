package geometry_test

import (
	"encoding"
	"strconv"
	"testing"

	"github.com/alecthomas/assert/v2"

	"github.com/twpayne/go-geos"
	"github.com/twpayne/go-geos/geometry"
)

var (
	_ encoding.TextMarshaler   = &geometry.Geometry{}
	_ encoding.TextUnmarshaler = &geometry.Geometry{}
)

func TestText(t *testing.T) {
	for i, tc := range []struct {
		geom    *geometry.Geometry
		textStr string
	}{
		{
			geom:    geometry.NewGeometry(geos.NewPoint([]float64{1, 2})),
			textStr: "POINT (1.0000000000000000 2.0000000000000000)",
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			text, err := tc.geom.MarshalText()
			assert.NoError(t, err)
			assert.Equal(t, tc.textStr, string(text))

			var geom geometry.Geometry
			assert.NoError(t, geom.UnmarshalText([]byte(tc.textStr)))
			assert.True(t, tc.geom.Equals(geom.Geom))
		})
	}
}
