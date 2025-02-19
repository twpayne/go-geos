package geos

// #include "go-geos.h"
import "C"

// A BufferParams contains parameters for BufferWithParams.
type BufferParams struct {
	context    *Context
	cBufParams *C.struct_GEOSBufParams_t
}

// NewBufferParams returns a new BufferParams.
func (c *Context) NewBufferParams() *BufferParams {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	cBufferParams := C.GEOSBufferParams_create_r(c.cHandle)
	if cBufferParams == nil {
		panic(c.err)
	}
	return c.newBufParams(cBufferParams)
}

// SetEndCapStyle sets p's end cap style.
func (p *BufferParams) SetEndCapStyle(style BufCapStyle) *BufferParams {
	p.context.mutex.Lock()
	defer p.context.mutex.Unlock()
	if C.GEOSBufferParams_setEndCapStyle_r(p.context.cHandle, p.cBufParams, C.int(style)) != 1 {
		panic(p.context.err)
	}
	return p
}

// SetJoinStyle sets p's join style.
func (p *BufferParams) SetJoinStyle(style BufJoinStyle) *BufferParams {
	p.context.mutex.Lock()
	defer p.context.mutex.Unlock()
	if C.GEOSBufferParams_setJoinStyle_r(p.context.cHandle, p.cBufParams, C.int(style)) != 1 {
		panic(p.context.err)
	}
	return p
}

// SetMitreLimit sets p's mitre limit.
func (p *BufferParams) SetMitreLimit(mitreLimit float64) *BufferParams {
	p.context.mutex.Lock()
	defer p.context.mutex.Unlock()
	if C.GEOSBufferParams_setMitreLimit_r(p.context.cHandle, p.cBufParams, C.double(mitreLimit)) != 1 {
		panic(p.context.err)
	}
	return p
}

// SetQuadrantSegments sets the number of segments to stroke each quadrant of
// circular arcs.
func (p *BufferParams) SetQuadrantSegments(quadSegs int) *BufferParams {
	p.context.mutex.Lock()
	defer p.context.mutex.Unlock()
	if C.GEOSBufferParams_setQuadrantSegments_r(p.context.cHandle, p.cBufParams, C.int(quadSegs)) != 1 {
		panic(p.context.err)
	}
	return p
}

// SetSingleSided sets whether the computed buffer should be single sided.
func (p *BufferParams) SetSingleSided(singleSided bool) *BufferParams {
	p.context.mutex.Lock()
	defer p.context.mutex.Unlock()
	if C.GEOSBufferParams_setSingleSided_r(p.context.cHandle, p.cBufParams, toInt[C.int](singleSided)) != 1 {
		panic(p.context.err)
	}
	return p
}
