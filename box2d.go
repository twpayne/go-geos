package geos

import (
	"fmt"
	"math"
)

// A Box2D is a two-dimensional bounds.
type Box2D struct {
	MinX float64
	MinY float64
	MaxX float64
	MaxY float64
}

// NewBox2D returns a new bounds.
func NewBox2D(minX, minY, maxX, maxY float64) *Box2D {
	return &Box2D{
		MinX: minX,
		MinY: minY,
		MaxX: maxX,
		MaxY: maxY,
	}
}

// NewBox2DEmpty returns a new empty bounds.
func NewBox2DEmpty() *Box2D {
	return &Box2D{
		MinX: math.Inf(1),
		MinY: math.Inf(1),
		MaxX: math.Inf(-1),
		MaxY: math.Inf(-1),
	}
}

// Contains returns true if b contains other.
func (b *Box2D) Contains(other *Box2D) bool {
	if b.IsEmpty() || other.IsEmpty() {
		return false
	}
	return other.MinX >= b.MinX && other.MinY >= b.MinY && other.MaxX <= b.MaxX && other.MaxY <= b.MaxY
}

// ContainsPoint returns true if b contains the point at x, y.
func (b *Box2D) ContainsPoint(x, y float64) bool {
	return b.MinX <= x && x <= b.MaxX && b.MinY <= y && y <= b.MaxY
}

// ContextGeom returns b as a Geom.
func (b *Box2D) ContextGeom(context *Context) *Geom {
	return context.NewGeomFromBounds(b.MinX, b.MinY, b.MaxX, b.MaxY)
}

// Equals returns true if b equals other.
func (b *Box2D) Equals(other *Box2D) bool {
	return b.MinX == other.MinX && b.MinY == other.MinY && b.MaxX == other.MaxX && b.MaxY == other.MaxY
}

// Geom returns b as a Geom.
func (b *Box2D) Geom() *Geom {
	return b.ContextGeom(DefaultContext)
}

// IsEmpty returns true if b is empty.
func (b *Box2D) IsEmpty() bool {
	return b.MinX > b.MaxX || b.MinY > b.MaxY
}

// Height returns the height of b.
func (b *Box2D) Height() float64 {
	return b.MaxY - b.MinY
}

// Intersects returns true if b intersects other.
func (b *Box2D) Intersects(other *Box2D) bool {
	return !(other.MinX > b.MaxX || other.MinY > b.MaxY || other.MaxX < b.MinX || other.MaxY < b.MinY)
}

// IsPoint returns true if b is a point.
func (b *Box2D) IsPoint() bool {
	return b.MinX == b.MaxX && b.MinY == b.MaxY
}

func (b *Box2D) String() string {
	return fmt.Sprintf("[%f %f %f %f]", b.MinX, b.MinY, b.MaxX, b.MaxY)
}

// Width returns the width of b.
func (b *Box2D) Width() float64 {
	return b.MaxX - b.MinX
}
