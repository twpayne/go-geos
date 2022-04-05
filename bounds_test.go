package geos

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBounds(t *testing.T) {
	b := NewBounds(1, 2, 3, 4)
	assert.True(t, b.Contains(b))
	assert.True(t, b.Contains(NewBounds(1.5, 2.5, 2.5, 3.5)))
	assert.False(t, b.Contains(NewBounds(1.5, 0.5, 2.5, 1.5)))
	assert.True(t, b.ContainsPoint(2, 3))
	assert.False(t, b.ContainsPoint(2, 1))
	assert.True(t, b.Equals(NewBounds(1, 2, 3, 4)))
	assert.Equal(t, "POLYGON ((1.0000000000000000 2.0000000000000000, 3.0000000000000000 2.0000000000000000, 3.0000000000000000 4.0000000000000000, 1.0000000000000000 4.0000000000000000, 1.0000000000000000 2.0000000000000000))", b.Geom().ToWKT())
	assert.False(t, b.IsEmpty())
	assert.Equal(t, 2.0, b.Height())
	assert.True(t, b.Intersects(b))
	assert.True(t, b.Intersects(NewBounds(1.5, 2.5, 2.5, 3.5)))
	assert.True(t, b.Intersects(NewBounds(1.5, 0.5, 2.5, 3.5)))
	assert.False(t, b.Intersects(NewBounds(1.5, 0.5, 2.5, 1.5)))
	assert.False(t, b.IsPoint())
	assert.Equal(t, 2.0, b.Width())
}

func TestBoundsEmpty(t *testing.T) {
	b := NewBoundsEmpty()
	assert.False(t, b.Contains(b))
	assert.False(t, b.Contains(NewBoundsEmpty()))
	assert.False(t, b.ContainsPoint(0, 0))
	//nolint:gocritic
	assert.True(t, b.Equals(b))
	assert.Equal(t, "POINT EMPTY", b.Geom().ToWKT())
	assert.True(t, b.IsEmpty())
	assert.False(t, b.Intersects(b))
	assert.False(t, b.IsPoint())
}

func TestBoundsPoint(t *testing.T) {
	b := NewBounds(0, 0, 0, 0)
	assert.True(t, b.Contains(b))
	assert.False(t, b.Contains(NewBounds(1, 2, 3, 4)))
	assert.True(t, b.ContainsPoint(0, 0))
	assert.False(t, b.ContainsPoint(1, 2))
	//nolint:gocritic
	assert.True(t, b.Equals(b))
	assert.False(t, b.Equals(NewBounds(1, 2, 3, 4)))
	assert.False(t, b.Equals(NewBoundsEmpty()))
	assert.Equal(t, "POINT (0.0000000000000000 0.0000000000000000)", b.Geom().ToWKT())
	assert.False(t, b.IsEmpty())
	assert.Equal(t, 0.0, b.Height())
	assert.True(t, b.IsPoint())
	assert.Equal(t, 0.0, b.Width())
}
