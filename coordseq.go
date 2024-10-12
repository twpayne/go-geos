package geos

// #include "go-geos.h"
import "C"

// A CoordSeq is a coordinate sequence.
type CoordSeq struct {
	context    *Context
	s          *C.struct_GEOSCoordSeq_t
	parent     *Geom
	dimensions int
	size       int
}

// Clone returns a clone of s.
func (s *CoordSeq) Clone() *CoordSeq {
	s.context.Lock()
	defer s.context.Unlock()
	return s.context.newNonNilCoordSeq(C.GEOSCoordSeq_clone_r(s.context.handle, s.s))
}

// Destroy destroys s and all resources associated with s.
func (s *CoordSeq) Destroy() {
	if s == nil || s.context == nil {
		return
	}
	s.context.Lock()
	defer s.context.Unlock()
	C.GEOSCoordSeq_destroy_r(s.context.handle, s.s)
	*s = CoordSeq{} // Clear all references.
}

// Dimensions returns the dimensions of s.
func (s *CoordSeq) Dimensions() int {
	return s.dimensions
}

// IsCCW returns if s is counter-clockwise.
func (s *CoordSeq) IsCCW() bool {
	s.context.Lock()
	defer s.context.Unlock()
	var cIsCCW C.char
	switch C.GEOSCoordSeq_isCCW_r(s.context.handle, s.s, &cIsCCW) {
	case 1:
		return cIsCCW != 0
	default:
		panic(s.context.err)
	}
}

// Ordinate returns the idx-th dim coordinate of s.
func (s *CoordSeq) Ordinate(idx, dim int) float64 {
	s.context.Lock()
	defer s.context.Unlock()
	if idx < 0 || s.size <= idx {
		panic(errIndexOutOfRange)
	}
	if dim < 0 || s.dimensions <= dim {
		panic(errDimensionOutOfRange)
	}
	var value float64
	if C.GEOSCoordSeq_getOrdinate_r(s.context.handle, s.s, C.uint(idx), C.uint(dim), (*C.double)(&value)) == 0 {
		panic(s.context.err)
	}
	return value
}

// SetOrdinate sets the idx-th dim coordinate of s to val.
func (s *CoordSeq) SetOrdinate(idx, dim int, val float64) {
	s.context.Lock()
	defer s.context.Unlock()
	if idx < 0 || s.size <= idx {
		panic(errIndexOutOfRange)
	}
	if dim < 0 || s.dimensions <= dim {
		panic(errDimensionOutOfRange)
	}
	if C.GEOSCoordSeq_setOrdinate_r(s.context.handle, s.s, C.uint(idx), C.uint(dim), C.double(val)) == 0 {
		panic(s.context.err)
	}
}

// SetX sets the idx-th X coordinate of s to val.
func (s *CoordSeq) SetX(idx int, val float64) {
	s.context.Lock()
	defer s.context.Unlock()
	if idx < 0 || s.size <= idx {
		panic(errIndexOutOfRange)
	}
	if s.dimensions == 0 {
		panic(errDimensionOutOfRange)
	}
	if C.GEOSCoordSeq_setX_r(s.context.handle, s.s, C.uint(idx), C.double(val)) == 0 {
		panic(s.context.err)
	}
}

// SetY sets the idx-th Y coordinate of s to val.
func (s *CoordSeq) SetY(idx int, val float64) {
	s.context.Lock()
	defer s.context.Unlock()
	if idx < 0 || s.size <= idx {
		panic(errIndexOutOfRange)
	}
	if s.dimensions < 2 {
		panic(errDimensionOutOfRange)
	}
	if C.GEOSCoordSeq_setY_r(s.context.handle, s.s, C.uint(idx), C.double(val)) == 0 {
		panic(s.context.err)
	}
}

// SetZ sets the idx-th Z coordinate of s to val.
func (s *CoordSeq) SetZ(idx int, val float64) {
	s.context.Lock()
	defer s.context.Unlock()
	if idx < 0 || s.size <= idx {
		panic(errIndexOutOfRange)
	}
	if s.dimensions < 3 {
		panic(errDimensionOutOfRange)
	}
	if C.GEOSCoordSeq_setZ_r(s.context.handle, s.s, C.uint(idx), C.double(val)) == 0 {
		panic(s.context.err)
	}
}

// Size returns the size of s.
func (s *CoordSeq) Size() int {
	return s.size
}

// ToCoords returns s as a [][]float64.
func (s *CoordSeq) ToCoords() [][]float64 {
	s.context.Lock()
	defer s.context.Unlock()
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
	if C.GEOSCoordSeq_copyToBuffer_r(s.context.handle, s.s, (*C.double)(&flatCoords[0]), hasZ, hasM) == 0 {
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
	s.context.Lock()
	defer s.context.Unlock()
	if idx < 0 || s.size <= idx {
		panic(errIndexOutOfRange)
	}
	if s.dimensions == 0 {
		panic(errDimensionOutOfRange)
	}
	var val float64
	if C.GEOSCoordSeq_getX_r(s.context.handle, s.s, C.uint(idx), (*C.double)(&val)) == 0 {
		panic(s.context.err)
	}
	return val
}

// Y returns the idx-th Y coordinate of s.
func (s *CoordSeq) Y(idx int) float64 {
	s.context.Lock()
	defer s.context.Unlock()
	if idx < 0 || s.size <= idx {
		panic(errIndexOutOfRange)
	}
	if s.dimensions < 2 {
		panic(errDimensionOutOfRange)
	}
	var val float64
	if C.GEOSCoordSeq_getY_r(s.context.handle, s.s, C.uint(idx), (*C.double)(&val)) == 0 {
		panic(s.context.err)
	}
	return val
}

// Z returns the idx-th Z coordinate of s.
func (s *CoordSeq) Z(idx int) float64 {
	s.context.Lock()
	defer s.context.Unlock()
	if idx < 0 || s.size <= idx {
		panic(errIndexOutOfRange)
	}
	if s.dimensions < 3 {
		panic(errDimensionOutOfRange)
	}
	var val float64
	if C.GEOSCoordSeq_getZ_r(s.context.handle, s.s, C.uint(idx), (*C.double)(&val)) == 0 {
		panic(s.context.err)
	}
	return val
}
