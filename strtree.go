package geos

// #include <stdlib.h>
// #include "go-geos.h"
import "C"

import (
	"runtime"
	"runtime/cgo"
	"unsafe"
)

// An STRtree is an R-tree spatial index structure for two dimensional data.
//
// WARNING The Go bindings to STRtree are currently broken. Do not use them.
type STRtree struct {
	context      *Context
	cSTRtree     *C.struct_GEOSSTRtree_t
	valueHandles map[any]*cgo.Handle
}

// NewSTRtree returns a new STRtree.
func (c *Context) NewSTRtree(nodeCapacity int) *STRtree {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	cSTRtree := C.GEOSSTRtree_create_r(c.cHandle, C.size_t(nodeCapacity))
	strTree := &STRtree{
		context:      c,
		cSTRtree:     cSTRtree,
		valueHandles: make(map[any]*cgo.Handle),
	}
	c.ref()
	runtime.AddCleanup(strTree, c.destroySTRtree, cSTRtree)
	runtime.AddCleanup(strTree, destoryValueHandles, strTree.valueHandles)
	return strTree
}

// Insert inserts value with geometry g.
func (t *STRtree) Insert(g *Geom, value any) error {
	if g.context != t.context {
		panic(errContextMismatch)
	}
	t.context.mutex.Lock()
	defer t.context.mutex.Unlock()
	if _, ok := t.valueHandles[value]; ok {
		return errDuplicateValue
	}
	valueHandle := cgo.NewHandle(value)
	t.valueHandles[value] = &valueHandle
	// FIXME golangci-lint complains about the following line saying:
	// dupSubExpr: suspicious identical LHS and RHS for `==` operator (gocritic)
	// As the line does not contain an `==` operator, disable gocritic on this
	// line.
	//nolint:gocritic
	C.GEOSSTRtree_insert_r(t.context.cHandle, t.cSTRtree, g.cGeom, unsafe.Pointer(&valueHandle))
	return nil
}

// Iterate calls f for every value in the t.
func (t *STRtree) Iterate(callback func(any)) {
	callbackHandle := cgo.NewHandle(func(item unsafe.Pointer) {
		valueHandle := (*cgo.Handle)(item)
		callback(valueHandle.Value())
	})
	defer callbackHandle.Delete()
	t.context.mutex.Lock()
	defer t.context.mutex.Unlock()
	C.GEOSSTRtree_iterate_r(
		t.context.cHandle,
		t.cSTRtree,
		(*[0]byte)(C.c_GEOSSTRtree_query_callback),
		unsafe.Pointer(&callbackHandle), //nolint:gocritic
	)
}

// Nearest returns the nearest geometry to geom in t.
//
// WARNING Nearest is currently broken and always panics with a segmentation
// fault.
func (t *STRtree) Nearest(geom *Geom) *Geom {
	t.context.mutex.Lock()
	defer t.context.mutex.Unlock()
	nearestGeom := C.GEOSSTRtree_nearest_r(t.context.cHandle, t.cSTRtree, geom.cGeom)
	if nearestGeom == nil {
		return nil
	}
	return t.context.newGeom(nearestGeom, nil)
}

// NearestGeneric returns the nearest value to value.
//
// WARNING NearestGeneric is currently broken and always panics with a
// segmentation fault.
func (t *STRtree) NearestGeneric(value any, valueEnvelope *Geom, distanceFunc func(any, any) float64) any {
	if t.context != valueEnvelope.context {
		panic(errContextMismatch)
	}
	callbackHandle := cgo.NewHandle(func(item1, item2 unsafe.Pointer, distance *C.double) C.int {
		// FIXME neither item1 not item2 should be nil, but in practice at least
		// one of them often is
		value1 := (*cgo.Handle)(item1).Value()
		value2 := (*cgo.Handle)(item2).Value()
		*distance = C.double(distanceFunc(value1, value2))
		return 1
	})
	defer callbackHandle.Delete()
	t.context.mutex.Lock()
	defer t.context.mutex.Unlock()
	nearestItem := C.GEOSSTRtree_nearest_generic_r(
		t.context.cHandle,
		t.cSTRtree,
		unsafe.Pointer(t.valueHandles[value]),
		valueEnvelope.cGeom,
		(*[0]byte)(C.c_GEOSSTRtree_distance_callback),
		unsafe.Pointer(&callbackHandle), //nolint:gocritic
	)
	if nearestItem == nil {
		return nil
	}
	nearestItemHandle := (*cgo.Handle)(nearestItem)
	return nearestItemHandle.Value()
}

// Query calls f with each value that intersects g.
func (t *STRtree) Query(g *Geom, callback func(any)) {
	callbackHandle := cgo.NewHandle(func(item unsafe.Pointer) {
		valueHandle := (*cgo.Handle)(item)
		callback(valueHandle.Value())
	})
	defer callbackHandle.Delete()
	t.context.mutex.Lock()
	defer t.context.mutex.Unlock()
	C.GEOSSTRtree_query_r(
		t.context.cHandle,
		t.cSTRtree,
		g.cGeom,
		(*[0]byte)(C.c_GEOSSTRtree_query_callback),
		unsafe.Pointer(&callbackHandle), //nolint:gocritic
	)
}

// Remove removes value with geometry g from t.
func (t *STRtree) Remove(g *Geom, value any) bool {
	if g.context != t.context {
		panic(errContextMismatch)
	}
	valueHandle := t.valueHandles[value]
	t.context.mutex.Lock()
	defer t.context.mutex.Unlock()
	switch C.GEOSSTRtree_remove_r(t.context.cHandle, t.cSTRtree, g.cGeom, unsafe.Pointer(valueHandle)) {
	case 0:
		return false
	case 1:
		delete(t.valueHandles, value)
		valueHandle.Delete()
		return true
	default:
		panic(t.context.err)
	}
}

func (c *Context) destroySTRtree(cSTRtree *C.struct_GEOSSTRtree_t) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	C.GEOSSTRtree_destroy_r(c.cHandle, cSTRtree)
	c.unref()
}

func destoryValueHandles(valueHandles map[any]*cgo.Handle) {
	for _, valueHandle := range valueHandles {
		valueHandle.Delete()
	}
}

//export go_GEOSSTRtree_distance_callback
func go_GEOSSTRtree_distance_callback(item1, item2 unsafe.Pointer, distance *C.double, userdata unsafe.Pointer) C.int {
	distanceCallbackHandle := (*cgo.Handle)(userdata)
	distanceCallback := distanceCallbackHandle.Value().(func(unsafe.Pointer, unsafe.Pointer, *C.double) C.int) //nolint:forcetypeassert,revive
	return distanceCallback(item1, item2, distance)
}

//export go_GEOSSTRtree_query_callback
func go_GEOSSTRtree_query_callback(item, userdata unsafe.Pointer) {
	callbackHandle := (*cgo.Handle)(userdata)
	callback := callbackHandle.Value().(func(unsafe.Pointer)) //nolint:forcetypeassert,revive
	callback(item)
}
