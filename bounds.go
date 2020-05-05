package geos

import (
	"fmt"
	"math"

	"github.com/llgcode/draw2d"
)

// A Bounds is a two-dimensional bounds.
type Bounds struct {
	MinX float64
	MinY float64
	MaxX float64
	MaxY float64
}

// NewBoundsEmpty returns a new empty bounds.
func NewBoundsEmpty() *Bounds {
	return &Bounds{
		MinX: math.Inf(1),
		MinY: math.Inf(1),
		MaxX: math.Inf(-1),
		MaxY: math.Inf(-1),
	}
}

// Contains returns true if b contains other.
func (b *Bounds) Contains(other *Bounds) bool {
	return other.MinX >= b.MinX && other.MinY >= b.MinY && other.MaxX <= b.MaxX && other.MaxY <= b.MaxY
}

// ContainsPoint returns true if b contains the point at x, y.
func (b *Bounds) ContainsPoint(x, y float64) bool {
	return b.MinX <= x && x <= b.MaxX && b.MinY <= y && y <= b.MaxY
}

// Draw draws b on gc.
func (b *Bounds) Draw(gc draw2d.GraphicContext) {
	gc.MoveTo(b.MinX, b.MinY)
	gc.LineTo(b.MinX, b.MaxY)
	gc.LineTo(b.MaxX, b.MaxY)
	gc.LineTo(b.MaxX, b.MinY)
	gc.LineTo(b.MinX, b.MinY)
	gc.Stroke()
}

// Equals returns true if b equals other.
func (b *Bounds) Equals(other *Bounds) bool {
	return b.MinX == other.MinX && b.MinY == other.MinY && b.MaxX == other.MaxX && b.MaxY == other.MaxY
}

// Geom returns b as a Geom.
func (b *Bounds) Geom() *Geom {
	return defaultContext.NewGeomFromBounds(b)
}

// IsEmpty returns true if b is empty.
func (b *Bounds) IsEmpty() bool {
	return b.MinX > b.MaxX || b.MinY > b.MaxY
}

// Height returns the height of b.
func (b *Bounds) Height() float64 {
	return b.MaxY - b.MinY
}

// Intersects returns true if b intersects other.
func (b *Bounds) Intersects(other *Bounds) bool {
	return !(other.MinX > b.MaxX || other.MinY > b.MaxY || other.MaxX < b.MinX || other.MaxY < b.MinY)
}

// IsPoint returns true if b is a point.
func (b *Bounds) IsPoint() bool {
	return b.MinX == b.MaxX && b.MinY == b.MaxY
}

func (b *Bounds) String() string {
	return fmt.Sprintf("[%f %f %f %f]", b.MinX, b.MinY, b.MaxX, b.MaxY)
}

// Width returns the width of b.
func (b *Bounds) Width() float64 {
	return b.MaxX - b.MinX
}
