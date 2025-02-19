package geos

// #include "go-geos.h"
import "C"

import (
	"runtime"
	"runtime/cgo"
	"unsafe"
)

// An STRtree is an R-tree spatial index structure for two dimensional data.
type STRtree struct {
	context     *Context
	cSTRtree    *C.struct_GEOSSTRtree_t
	itemToValue map[unsafe.Pointer]any
	valueToItem map[any]unsafe.Pointer
}

// NewSTRtree returns a new STRtree.
func (c *Context) NewSTRtree(nodeCapacity int) *STRtree {
	c.Lock()
	defer c.Unlock()
	t := &STRtree{
		context:     c,
		cSTRtree:    C.GEOSSTRtree_create_r(c.cHandle, C.size_t(nodeCapacity)),
		itemToValue: make(map[unsafe.Pointer]any),
		valueToItem: make(map[any]unsafe.Pointer),
	}
	runtime.AddCleanup(t, func(cSTRtree *C.struct_GEOSSTRtree_t) {
		C.GEOSSTRtree_destroy_r(c.cHandle, cSTRtree)
	}, t.cSTRtree)
	return t
}

// Insert inserts value with geometry g.
func (t *STRtree) Insert(g *Geom, value any) error {
	if g.context != t.context {
		panic(errContextMismatch)
	}
	t.context.Lock()
	defer t.context.Unlock()
	if _, ok := t.valueToItem[value]; ok {
		return errDuplicateValue
	}
	item := C.calloc(1, C.size_t(unsafe.Sizeof(uintptr(0))))
	t.itemToValue[item] = value
	t.valueToItem[value] = item
	C.GEOSSTRtree_insert_r(t.context.cHandle, t.cSTRtree, g.cGeom, item)
	return nil
}

// Iterate calls f for every value in the t.
func (t *STRtree) Iterate(callback func(any)) {
	handle := cgo.NewHandle(func(item unsafe.Pointer) {
		callback(t.itemToValue[item])
	})
	defer handle.Delete()
	t.context.Lock()
	defer t.context.Unlock()
	C.GEOSSTRtree_iterate_r(
		t.context.cHandle,
		t.cSTRtree,
		(*[0]byte)(C.c_GEOSSTRtree_query_callback), // FIXME understand why the cast to *[0]byte is needed
		unsafe.Pointer(&handle),                    //nolint:gocritic
	)
}

// Nearest returns the nearest item in t to value.
func (t *STRtree) Nearest(value any, valueEnvelope *Geom, geomfn func(any) *Geom) any {
	if t.context != valueEnvelope.context {
		panic(errContextMismatch)
	}
	handle := cgo.NewHandle(func(item1, item2 unsafe.Pointer, distance *C.double) C.int {
		geom1 := geomfn(t.itemToValue[item1])
		if geom1 == nil {
			return 0
		}
		geom2 := geomfn(t.itemToValue[item2])
		if geom2 == nil {
			return 0
		}
		return C.GEOSDistance_r(t.context.cHandle, geom1.cGeom, geom2.cGeom, distance)
	})
	defer handle.Delete()
	t.context.Lock()
	defer t.context.Unlock()
	nearestItem := C.GEOSSTRtree_nearest_generic_r(
		t.context.cHandle,
		t.cSTRtree,
		t.valueToItem[value],
		valueEnvelope.cGeom,
		(*[0]byte)(C.c_GEOSSTRtree_distance_callback), // FIXME understand why the cast to *[0]byte is needed
		unsafe.Pointer(&handle),                       //nolint:gocritic
	)
	return t.itemToValue[nearestItem]
}

// Query calls f with each value that intersects g.
func (t *STRtree) Query(g *Geom, callback func(any)) {
	handle := cgo.NewHandle(func(elem unsafe.Pointer) {
		callback(t.itemToValue[elem])
	})
	defer handle.Delete()
	t.context.Lock()
	defer t.context.Unlock()
	C.GEOSSTRtree_query_r(
		t.context.cHandle,
		t.cSTRtree,
		g.cGeom,
		(*[0]byte)(C.c_GEOSSTRtree_query_callback), // FIXME understand why the cast to *[0]byte is needed
		unsafe.Pointer(&handle),                    //nolint:gocritic
	)
}

// Remove removes value with geometry g from t.
func (t *STRtree) Remove(g *Geom, value any) bool {
	if g.context != t.context {
		panic(errContextMismatch)
	}
	item := t.valueToItem[value]
	t.context.Lock()
	defer t.context.Unlock()
	switch C.GEOSSTRtree_remove_r(t.context.cHandle, t.cSTRtree, g.cGeom, item) {
	case 0:
		return false
	case 1:
		delete(t.valueToItem, value)
		delete(t.itemToValue, item)
		C.free(item)
		return true
	default:
		panic(t.context.err)
	}
}

//export go_GEOSSTRtree_distance_callback
func go_GEOSSTRtree_distance_callback(item1, item2 unsafe.Pointer, distance *C.double, userdata unsafe.Pointer) C.int {
	handle := *(*cgo.Handle)(userdata)
	return handle.Value().(func(unsafe.Pointer, unsafe.Pointer, *C.double) C.int)(item1, item2, distance) //nolint:forcetypeassert
}

//export go_GEOSSTRtree_query_callback
func go_GEOSSTRtree_query_callback(elem, userdata unsafe.Pointer) {
	handle := *(*cgo.Handle)(userdata)
	handle.Value().(func(unsafe.Pointer))(elem) //nolint:forcetypeassert
}
