package geos

// #include "go-geos.h"
import "C"

import "runtime"

// A BufParams contains parameters for BufferWithParams.
type BufParams struct {
	context    *Context
	cBufParams *C.struct_GEOSBufParams_t
}

// NewBufParams returns a new BufParams.
func (c *Context) NewBufParams() *BufParams {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	cBufParams := C.GEOSBufferParams_create_r(c.cHandle)
	if cBufParams == nil {
		panic(c.err)
	}
	bufParams := &BufParams{
		context:    c,
		cBufParams: cBufParams,
	}
	c.ref()
	runtime.AddCleanup(bufParams, c.destroyBufParams, cBufParams)
	return bufParams
}

// SetEndCapStyle sets p's end cap style.
func (p *BufParams) SetEndCapStyle(style BufCapStyle) *BufParams {
	p.context.mutex.Lock()
	defer p.context.mutex.Unlock()
	if C.GEOSBufferParams_setEndCapStyle_r(p.context.cHandle, p.cBufParams, C.int(style)) != 1 {
		panic(p.context.err)
	}
	return p
}

// SetJoinStyle sets p's join style.
func (p *BufParams) SetJoinStyle(style BufJoinStyle) *BufParams {
	p.context.mutex.Lock()
	defer p.context.mutex.Unlock()
	if C.GEOSBufferParams_setJoinStyle_r(p.context.cHandle, p.cBufParams, C.int(style)) != 1 {
		panic(p.context.err)
	}
	return p
}

// SetMitreLimit sets p's mitre limit.
func (p *BufParams) SetMitreLimit(mitreLimit float64) *BufParams {
	p.context.mutex.Lock()
	defer p.context.mutex.Unlock()
	if C.GEOSBufferParams_setMitreLimit_r(p.context.cHandle, p.cBufParams, C.double(mitreLimit)) != 1 {
		panic(p.context.err)
	}
	return p
}

// SetQuadrantSegments sets the number of segments to stroke each quadrant of
// circular arcs.
func (p *BufParams) SetQuadrantSegments(quadSegs int) *BufParams {
	p.context.mutex.Lock()
	defer p.context.mutex.Unlock()
	if C.GEOSBufferParams_setQuadrantSegments_r(p.context.cHandle, p.cBufParams, C.int(quadSegs)) != 1 {
		panic(p.context.err)
	}
	return p
}

// SetSingleSided sets whether the computed buffer should be single sided.
func (p *BufParams) SetSingleSided(singleSided bool) *BufParams {
	p.context.mutex.Lock()
	defer p.context.mutex.Unlock()
	if C.GEOSBufferParams_setSingleSided_r(p.context.cHandle, p.cBufParams, toInt[C.int](singleSided)) != 1 {
		panic(p.context.err)
	}
	return p
}

func (c *Context) destroyBufParams(cBufParams *C.struct_GEOSBufParams_t) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	C.GEOSBufferParams_destroy_r(c.cHandle, cBufParams)
	c.unref()
}
