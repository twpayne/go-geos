package geos

// #include "go-geos.h"
import "C"

// A BufferParams contains parameters for BufferWithParams.
type BufferParams struct {
	context      *Context
	bufferParams *C.struct_GEOSBufParams_t
}

// Destroy destroys all resources associated with p.
func (p *BufferParams) Destroy() {
	// Protect against Destroy being called more than once.
	if p == nil || p.context == nil {
		return
	}
	p.context.Lock()
	defer p.context.Unlock()
	C.GEOSBufferParams_destroy_r(p.context.handle, p.bufferParams)
	*p = BufferParams{} // Clear all references.
}

// SetEndCapStyle sets p's end cap style.
func (p *BufferParams) SetEndCapStyle(style BufCapStyle) *BufferParams {
	p.context.Lock()
	defer p.context.Unlock()
	if C.GEOSBufferParams_setEndCapStyle_r(p.context.handle, p.bufferParams, C.int(style)) != 1 {
		panic(p.context.err)
	}
	return p
}

// SetJoinStyle sets p's join style.
func (p *BufferParams) SetJoinStyle(style BufJoinStyle) *BufferParams {
	p.context.Lock()
	defer p.context.Unlock()
	if C.GEOSBufferParams_setJoinStyle_r(p.context.handle, p.bufferParams, C.int(style)) != 1 {
		panic(p.context.err)
	}
	return p
}

// SetMitreLimit sets p's mitre limit.
func (p *BufferParams) SetMitreLimit(mitreLimit float64) *BufferParams {
	p.context.Lock()
	defer p.context.Unlock()
	if C.GEOSBufferParams_setMitreLimit_r(p.context.handle, p.bufferParams, C.double(mitreLimit)) != 1 {
		panic(p.context.err)
	}
	return p
}

// SetQuadrantSegments sets the number of segments to stroke each quadrant of
// circular arcs.
func (p *BufferParams) SetQuadrantSegments(quadSegs int) *BufferParams {
	p.context.Lock()
	defer p.context.Unlock()
	if C.GEOSBufferParams_setQuadrantSegments_r(p.context.handle, p.bufferParams, C.int(quadSegs)) != 1 {
		panic(p.context.err)
	}
	return p
}

// SetSingleSided sets whether the computed buffer should be single sided.
func (p *BufferParams) SetSingleSided(singleSided bool) *BufferParams {
	p.context.Lock()
	defer p.context.Unlock()
	if C.GEOSBufferParams_setSingleSided_r(p.context.handle, p.bufferParams, C.int(intFromBool(singleSided))) != 1 {
		panic(p.context.err)
	}
	return p
}

func (p *BufferParams) finalize() {
	if p.context == nil {
		return
	}
	p.Destroy()
}

func intFromBool(b bool) int {
	if b {
		return 1
	}
	return 0
}
