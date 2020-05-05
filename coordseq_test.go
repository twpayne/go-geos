package geos

import (
	"fmt"
	"math"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCoordSeqMethods(t *testing.T) {
	defer runtime.GC() // Exercise finalizers.
	c := NewContext()
	s := c.NewCoordSeq(2, 3)
	assert.Equal(t, 2, s.Size())
	assert.Equal(t, 3, s.Dimensions())
	assert.Equal(t, 0.0, s.X(0))
	assert.Equal(t, 0.0, s.Y(0))
	assert.True(t, math.IsNaN(s.Z(0)))
	s.SetX(1, 1)
	s.SetY(1, 2)
	s.SetZ(1, 3)
	assert.Equal(t, 1.0, s.X(1))
	assert.Equal(t, 2.0, s.Y(1))
	assert.Equal(t, 3.0, s.Z(1))
	assert.Equal(t, 1.0, s.Ordinate(1, 0))
	assert.Equal(t, 2.0, s.Ordinate(1, 1))
	assert.Equal(t, 3.0, s.Ordinate(1, 2))
	clone := s.Clone()
	assert.Equal(t, 1.0, clone.X(1))
	assert.Equal(t, 2.0, clone.Y(1))
	clone.SetOrdinate(0, 0, -1.0)
	clone.SetOrdinate(0, 1, -2.0)
	assert.Equal(t, -1.0, clone.X(0))
	assert.Equal(t, -2.0, clone.Y(0))

	// GEOS version 3.8.0 distributed with homebrew on macOS has a bug where
	// GEOSCoordSeq_clone_r does not clone the dimension correctly. The original
	// has dimension 3 but the returned clone has dimension 2. As we do not use
	// three dimensional geometries, we are not affected by this bug, so skip
	// the tests that fail because of the bug.
	if clone.Dimensions() == 2 {
		t.Skip("skipping tests in buggy GEOS library")
	}
	require.Equal(t, 3, clone.Dimensions())
	assert.Equal(t, 3.0, clone.Z(1))
	clone.SetOrdinate(0, 2, -3.0)
	assert.Equal(t, -3.0, clone.Z(0))
}

func TestCoordSeqPanics(t *testing.T) {
	c := NewContext()
	s := c.NewCoordSeq(1, 2)

	assert.Panics(t, func() { s.X(-1) })
	assert.NotPanics(t, func() { s.X(0) })
	assert.Panics(t, func() { s.X(1) })

	assert.Panics(t, func() { s.Y(-1) })
	assert.NotPanics(t, func() { s.Y(0) })
	assert.Panics(t, func() { s.Y(1) })

	assert.Panics(t, func() { s.Z(-1) })
	assert.Panics(t, func() { s.Z(0) })
	assert.Panics(t, func() { s.Z(1) })

	assert.Panics(t, func() { s.SetX(-1, 0) })
	assert.NotPanics(t, func() { s.SetX(0, 0) })
	assert.Panics(t, func() { s.SetX(1, 0) })

	assert.Panics(t, func() { s.SetY(-1, 0) })
	assert.NotPanics(t, func() { s.SetY(0, 0) })
	assert.Panics(t, func() { s.SetY(1, 0) })

	assert.Panics(t, func() { s.SetZ(-1, 0) })
	assert.Panics(t, func() { s.SetZ(0, 0) })
	assert.Panics(t, func() { s.SetZ(1, 0) })

	for idx := -1; idx <= 1; idx++ {
		for dim := -1; dim <= 4; dim++ {
			t.Run(fmt.Sprintf("idx_%d_dim_%d", idx, dim), func(t *testing.T) {
				if idx == 0 && 0 <= dim && dim < 2 {
					assert.NotPanics(t, func() { s.Ordinate(idx, dim) })
					assert.NotPanics(t, func() { s.SetOrdinate(idx, dim, 0) })
				} else {
					assert.Panics(t, func() { s.Ordinate(idx, dim) })
					assert.Panics(t, func() { s.SetOrdinate(idx, dim, 0) })
				}
			})
		}
	}
}

func TestCoordSeqCoordsMethods(t *testing.T) {
	for _, tc := range []struct {
		name   string
		coords [][]float64
	}{
		{
			name:   "point_2d",
			coords: [][]float64{{1, 2}},
		},
		{
			name:   "point_3d",
			coords: [][]float64{{1, 2, 3}},
		},
		{
			name:   "linestring_2d",
			coords: [][]float64{{1, 2}, {3, 4}},
		},
		{
			name:   "linestring_3d",
			coords: [][]float64{{1, 2, 3}, {4, 5, 6}},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			defer runtime.GC() // Exercise finalizers.
			c := NewContext()
			s := c.NewCoordSeqFromCoords(tc.coords)
			assert.Equal(t, tc.coords, s.ToCoords())
		})
	}
}
