package geometry

import (
	"encoding"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/twpayne/go-geos"
)

var (
	_ encoding.TextMarshaler   = &Geometry{}
	_ encoding.TextUnmarshaler = &Geometry{}
)

func TestText(t *testing.T) {
	for i, tc := range []struct {
		geom    *Geometry
		textStr string
	}{
		{
			geom:    NewGeometry(geos.NewPoint([]float64{1, 2})),
			textStr: "POINT (1.0000000000000000 2.0000000000000000)",
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			text, err := tc.geom.MarshalText()
			require.NoError(t, err)
			assert.Equal(t, tc.textStr, string(text))

			var geom Geometry
			require.NoError(t, geom.UnmarshalText([]byte(tc.textStr)))
			assert.True(t, tc.geom.Equals(geom.Geom))
		})
	}
}
