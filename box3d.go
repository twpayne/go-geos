package geos

import (
	"fmt"
	"math"
)

// A Box3D is a three-dimensional bounds.
type Box3D struct {
	MinX float64
	MinY float64
	MinZ float64
	MaxX float64
	MaxY float64
	MaxZ float64
}

// NewBox3D returns a new bounds.
func NewBox3D(minX, minY, minZ, maxX, maxY, maxZ float64) *Box3D {
	return &Box3D{
		MinX: minX,
		MinY: minY,
		MinZ: minZ,
		MaxX: maxX,
		MaxY: maxY,
		MaxZ: maxZ,
	}
}

// NewBox3DEmpty returns a new empty bounds.
func NewBox3DEmpty() *Box3D {
	return &Box3D{
		MinX: math.Inf(1),
		MinY: math.Inf(1),
		MinZ: math.Inf(1),
		MaxX: math.Inf(-1),
		MaxY: math.Inf(-1),
		MaxZ: math.Inf(-1),
	}
}

func (b *Box3D) String() string {
	return fmt.Sprintf("[%f %f %f %f %f %f]", b.MinX, b.MinY, b.MinZ, b.MaxX, b.MaxY, b.MaxX)
}
