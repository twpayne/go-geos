# geos

[![PkgGoDev](https://pkg.go.dev/badge/github.com/twpayne/go-geos)](https://pkg.go.dev/github.com/twpayne/go-geos)

Package `geos` provides an interface to [GEOS](https://trac.osgeo.org/geos).

## Features

* Idiomatic Go API.

* `geometry.Geometry` type implements all GEOS functionality and standard Go
  interfaces:

  *  `sql/driver.Valuer` and `sql.Scanner` (EWKB) for PostGIS database integration.
  *  `json.Marshaler` and `json.Unmarshaler` (GeoJSON).
  *  `xml.Marshaler` (KML).
  *  `gob.GobEncoder` and `gob.GobDecoder` (GOB).

* Concurrency-safe. `geos` uses GEOS's threadsafe `*_r` functions under the
  hood, with locking to ensure safety, even when used across multiple
  goroutines. For best performance, use one `geos.Context` per goroutine.

* Caching of some geometry properties to avoid cgo overhead.

* Automatic finalization of GEOS objects.

## Exceptions

`geos` uses the stable C GEOS bindings. These bindings catch exceptions from the
underlying C++ code and convert them to a return code. For normal geometry
operations, `geos` `panic`s whenever it encounters a return code indicating an
error, rather than returning an `error`. This behavior is similar to slice
access in Go (out-of-bounds accesses `panic`) and keeps the API easy to use.
When parsing WKB and WKT, errors are expected so an `error` is returned.

## Memory management

`geos` objects live mostly on the C heap, and `geos` sets finalizers on the
objects it creates that free the associated C memory. However, the C heap is not
visible to the Go runtime. The can result in significant memory pressure as
memory is consumed by large, non-finalized geometries, of which the Go runtime
is unaware. Consequently, if it is known that a geometry will no longer be used,
it should be explicitly freed by calling its `Destroy()` method. Periodic calls
to `runtime.GC()` can also help, but the Go runtime makes no guarantees about
when or if finalizers will be called.

For more information, see the [documentation for
`runtime.SetFinalizer()`](https://pkg.go.dev/runtime#SetFinalizer) and [this
thread on
`golang-nuts`](https://groups.google.com/g/golang-nuts/c/XnV16PxXBfA/m/W8VEzIvHBAAJ).

## Licence

MIT