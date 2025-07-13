package geos

// #include "go-geos.h"
import "C"

import (
	"runtime"
	"unsafe"
)

// A CoordSeq is a coordinate sequence.
type CoordSeq struct {
	context    *Context
	s          *C.struct_GEOSCoordSeq_t
	owner      *Geom
	dimensions int
	size       int
}

// NewCoordSeq returns a new CoordSeq.
func (c *Context) NewCoordSeq(size, dims int) *CoordSeq {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.newNonNilCoordSeq(C.GEOSCoordSeq_create_r(c.cHandle, C.uint(size), C.uint(dims)))
}

// NewCoordSeqFromCoords returns a new CoordSeq populated with coords.
func (c *Context) NewCoordSeqFromCoords(coords [][]float64) *CoordSeq {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.newNonNilCoordSeq(c.newGEOSCoordSeqFromCoords(coords))
}

// Clone returns a clone of s.
func (s *CoordSeq) Clone() *CoordSeq {
	s.context.mutex.Lock()
	defer s.context.mutex.Unlock()
	return s.context.newNonNilCoordSeq(C.GEOSCoordSeq_clone_r(s.context.cHandle, s.s))
}

// Dimensions returns the dimensions of s.
func (s *CoordSeq) Dimensions() int {
	return s.dimensions
}

// IsCCW returns if s is counter-clockwise.
func (s *CoordSeq) IsCCW() bool {
	s.context.mutex.Lock()
	defer s.context.mutex.Unlock()
	var cIsCCW C.char
	switch C.GEOSCoordSeq_isCCW_r(s.context.cHandle, s.s, &cIsCCW) {
	case 1:
		return cIsCCW != 0
	default:
		panic(s.context.err)
	}
}

// Ordinate returns the idx-th dim coordinate of s.
func (s *CoordSeq) Ordinate(idx, dim int) float64 {
	s.context.mutex.Lock()
	defer s.context.mutex.Unlock()
	if idx < 0 || s.size <= idx {
		panic(errIndexOutOfRange)
	}
	if dim < 0 || s.dimensions <= dim {
		panic(errDimensionOutOfRange)
	}
	var value float64
	if C.GEOSCoordSeq_getOrdinate_r(s.context.cHandle, s.s, C.uint(idx), C.uint(dim), (*C.double)(&value)) == 0 {
		panic(s.context.err)
	}
	return value
}

// SetOrdinate sets the idx-th dim coordinate of s to val.
func (s *CoordSeq) SetOrdinate(idx, dim int, val float64) {
	s.context.mutex.Lock()
	defer s.context.mutex.Unlock()
	if idx < 0 || s.size <= idx {
		panic(errIndexOutOfRange)
	}
	if dim < 0 || s.dimensions <= dim {
		panic(errDimensionOutOfRange)
	}
	if C.GEOSCoordSeq_setOrdinate_r(s.context.cHandle, s.s, C.uint(idx), C.uint(dim), C.double(val)) == 0 {
		panic(s.context.err)
	}
}

// SetX sets the idx-th X coordinate of s to val.
func (s *CoordSeq) SetX(idx int, val float64) {
	s.context.mutex.Lock()
	defer s.context.mutex.Unlock()
	if idx < 0 || s.size <= idx {
		panic(errIndexOutOfRange)
	}
	if s.dimensions == 0 {
		panic(errDimensionOutOfRange)
	}
	if C.GEOSCoordSeq_setX_r(s.context.cHandle, s.s, C.uint(idx), C.double(val)) == 0 {
		panic(s.context.err)
	}
}

// SetY sets the idx-th Y coordinate of s to val.
func (s *CoordSeq) SetY(idx int, val float64) {
	s.context.mutex.Lock()
	defer s.context.mutex.Unlock()
	if idx < 0 || s.size <= idx {
		panic(errIndexOutOfRange)
	}
	if s.dimensions < 2 {
		panic(errDimensionOutOfRange)
	}
	if C.GEOSCoordSeq_setY_r(s.context.cHandle, s.s, C.uint(idx), C.double(val)) == 0 {
		panic(s.context.err)
	}
}

// SetZ sets the idx-th Z coordinate of s to val.
func (s *CoordSeq) SetZ(idx int, val float64) {
	s.context.mutex.Lock()
	defer s.context.mutex.Unlock()
	if idx < 0 || s.size <= idx {
		panic(errIndexOutOfRange)
	}
	if s.dimensions < 3 {
		panic(errDimensionOutOfRange)
	}
	if C.GEOSCoordSeq_setZ_r(s.context.cHandle, s.s, C.uint(idx), C.double(val)) == 0 {
		panic(s.context.err)
	}
}

// Size returns the size of s.
func (s *CoordSeq) Size() int {
	return s.size
}

// ToCoords returns s as a [][]float64.
func (s *CoordSeq) ToCoords() [][]float64 {
	s.context.mutex.Lock()
	defer s.context.mutex.Unlock()
	if s.size == 0 || s.dimensions == 0 {
		return nil
	}
	flatCoords := make([]float64, s.size*s.dimensions)
	var hasZ C.int
	if s.dimensions > 2 {
		hasZ = 1
	}
	var hasM C.int
	if s.dimensions > 3 {
		hasM = 1
	}
	if C.GEOSCoordSeq_copyToBuffer_r(s.context.cHandle, s.s, (*C.double)(&flatCoords[0]), hasZ, hasM) == 0 {
		panic(s.context.err)
	}
	coords := make([][]float64, s.size)
	j := 0
	for i := range s.size {
		coords[i] = flatCoords[j : j+s.dimensions : j+s.dimensions]
		j += s.dimensions
	}
	return coords
}

// X returns the idx-th X coordinate of s.
func (s *CoordSeq) X(idx int) float64 {
	s.context.mutex.Lock()
	defer s.context.mutex.Unlock()
	if idx < 0 || s.size <= idx {
		panic(errIndexOutOfRange)
	}
	if s.dimensions == 0 {
		panic(errDimensionOutOfRange)
	}
	var val float64
	if C.GEOSCoordSeq_getX_r(s.context.cHandle, s.s, C.uint(idx), (*C.double)(&val)) == 0 {
		panic(s.context.err)
	}
	return val
}

// Y returns the idx-th Y coordinate of s.
func (s *CoordSeq) Y(idx int) float64 {
	s.context.mutex.Lock()
	defer s.context.mutex.Unlock()
	if idx < 0 || s.size <= idx {
		panic(errIndexOutOfRange)
	}
	if s.dimensions < 2 {
		panic(errDimensionOutOfRange)
	}
	var val float64
	if C.GEOSCoordSeq_getY_r(s.context.cHandle, s.s, C.uint(idx), (*C.double)(&val)) == 0 {
		panic(s.context.err)
	}
	return val
}

// Z returns the idx-th Z coordinate of s.
func (s *CoordSeq) Z(idx int) float64 {
	s.context.mutex.Lock()
	defer s.context.mutex.Unlock()
	if idx < 0 || s.size <= idx {
		panic(errIndexOutOfRange)
	}
	if s.dimensions < 3 {
		panic(errDimensionOutOfRange)
	}
	var val float64
	if C.GEOSCoordSeq_getZ_r(s.context.cHandle, s.s, C.uint(idx), (*C.double)(&val)) == 0 {
		panic(s.context.err)
	}
	return val
}

func (c *Context) newCoordSeqInternal(cCoordSeq *C.struct_GEOSCoordSeq_t, owner *Geom) *CoordSeq {
	if cCoordSeq == nil {
		return nil
	}
	var (
		dimensions C.uint
		size       C.uint
	)
	if C.GEOSCoordSeq_getDimensions_r(c.cHandle, cCoordSeq, &dimensions) == 0 {
		panic(c.err)
	}
	if C.GEOSCoordSeq_getSize_r(c.cHandle, cCoordSeq, &size) == 0 {
		panic(c.err)
	}
	coordSeq := &CoordSeq{
		context:    c,
		s:          cCoordSeq,
		owner:      owner,
		dimensions: int(dimensions),
		size:       int(size),
	}
	if owner == nil {
		c.ref()
		runtime.AddCleanup(coordSeq, c.destroyCoordSeq, cCoordSeq)
	}
	return coordSeq
}

func (c *Context) newCoordsFromGEOSCoordSeq(cCoordSeq *C.struct_GEOSCoordSeq_t) [][]float64 {
	var dimensions C.uint
	if C.GEOSCoordSeq_getDimensions_r(c.cHandle, cCoordSeq, &dimensions) == 0 {
		panic(c.err)
	}

	var size C.uint
	if C.GEOSCoordSeq_getSize_r(c.cHandle, cCoordSeq, &size) == 0 {
		panic(c.err)
	}

	var hasZ C.int
	if dimensions > 2 {
		hasZ = 1
	}

	var hasM C.int
	if dimensions > 3 {
		hasM = 1
	}

	flatCoords := make([]float64, size*dimensions)
	if C.GEOSCoordSeq_copyToBuffer_r(c.cHandle, cCoordSeq, (*C.double)(&flatCoords[0]), hasZ, hasM) == 0 {
		panic(c.err)
	}
	coords := make([][]float64, size)
	for i := range coords {
		coord := flatCoords[i*int(dimensions) : (i+1)*int(dimensions) : (i+1)*int(dimensions)]
		coords[i] = coord
	}
	return coords
}

func (c *Context) newGEOSCoordSeqFromCoords(coords [][]float64) *C.struct_GEOSCoordSeq_t {
	var hasZ C.int
	if len(coords[0]) > 2 {
		hasZ = 1
	}

	var hasM C.int
	if len(coords[0]) > 3 {
		hasM = 1
	}

	dimensions := len(coords[0])
	flatCoords := make([]float64, len(coords)*dimensions)
	for i, coord := range coords {
		copy(flatCoords[i*dimensions:(i+1)*dimensions], coord)
	}
	return C.GEOSCoordSeq_copyFromBuffer_r(c.cHandle, (*C.double)(unsafe.Pointer(&flatCoords[0])), C.uint(len(coords)), hasZ, hasM)
}

func (c *Context) newNonNilCoordSeq(cCoordSeq *C.struct_GEOSCoordSeq_t) *CoordSeq {
	if cCoordSeq == nil {
		panic(c.err)
	}
	return c.newCoordSeqInternal(cCoordSeq, nil)
}

func (c *Context) destroyCoordSeq(cCoordSeq *C.struct_GEOSCoordSeq_t) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	C.GEOSCoordSeq_destroy_r(c.cHandle, cCoordSeq)
	c.unref()
}
