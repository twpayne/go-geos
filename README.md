# go-geos

[![PkgGoDev](https://pkg.go.dev/badge/github.com/twpayne/go-geos)](https://pkg.go.dev/github.com/twpayne/go-geos)

Package `go-geos` provides an interface to [GEOS](https://libgeos.org).

## Install

```console
$ go get github.com/twpayne/go-geos
```

You must also install the GEOS development headers and libraries. These are
typically in the package `libgeos-dev` on Debian-like systems, `geos-devel` on
RedHat-like systems, and `geos` in Homebrew.

## Features

* Fluent Go API.

* Low-level `Context`, `CoordSeq`, `Geom`, `PrepGeom`, and `STRtree` types
  provide access to all GEOS methods.

* High-level `geometry.Geometry` type implements all GEOS functionality and
  many standard Go interfaces:

  * `database/sql/driver.Valuer` and `database/sql.Scanner` (WKB) for PostGIS
     database integration.
  * `encoding/json.Marshaler` and `encoding/json.Unmarshaler` (GeoJSON).
  * `encoding/xml.Marshaler` (KML).
  * `encoding.BinaryMarshaler` and `encoding.BinaryUnmarshaler` (WKB).
  * `encoding.TextMarshaler` and `encoding.TextUnmarshaler` (WKT).
  * `encoding/gob.GobEncoder` and `encoding/gob.GobDecoder` (GOB).

  See the [PostGIS example](examples/postgis/README.md) for a demonstration of
  the use of these interfaces.

* Concurrency-safe. `go-geos` uses GEOS's threadsafe `*_r` functions under the
  hood, with locking to ensure safety, even when used across multiple
  goroutines. For best performance, use one `geos.Context` per goroutine.

* Caching of geometry properties to avoid cgo overhead.

* Optimized GeoJSON encoder.

* Automatic finalization of GEOS objects.

## Memory management

`go-geos` objects live mostly on the C heap. `go-geos` sets finalizers on the
objects it creates that free the associated C memory. However, the C heap is not
visible to the Go runtime. The can result in significant memory pressure as
memory is consumed by large, non-finalized geometries, of which the Go runtime
is unaware. Consequently, if it is known that a geometry will no longer be used,
it should be explicitly freed by calling its `Destroy()` method. Periodic calls
to `runtime.GC()` can also help, but the Go runtime makes no guarantees about
when or if finalizers will be called.

You can set a function to be called whenever a geometry's finalizer is invoked
with the `WithGeomFinalizeFunc` option to `NewContext()`. This can be helpful
for tracking down geometry leaks.

For more information, see the [documentation for
`runtime.SetFinalizer()`](https://pkg.go.dev/runtime#SetFinalizer) and [this
thread on
`golang-nuts`](https://groups.google.com/g/golang-nuts/c/XnV16PxXBfA/m/W8VEzIvHBAAJ).

## Errors, exceptions, and panics

`go-geos` uses the stable GEOS C bindings. These bindings catch exceptions from
the underlying C++ code and convert them to an integer return code. For normal
geometry operations, `go-geos` panics whenever it encounters a GEOS return code
indicating an error, rather than returning an `error`. Such panics will not
occur if `go-geos` is used correctly. Panics will occur for invalid API calls,
out-of-bounds access, or operations on invalid geometries. This behavior is
similar to slice access in Go (out-of-bounds accesses panic) and keeps the API
fluent. When parsing data, errors are expected so an `error` is returned.

## Comparison with `github.com/twpayne/go-geom`

[`github.com/twpayne/go-geom`](https://github.com/twpayne/go-geom) is a pure Go
library providing similar functionality to `go-geos`. The major differences are:

* `go-geos` uses [GEOS](https://libgeos.org), which is an extremely mature
  library with a rich feature set.
* `go-geos` uses cgo, with all the disadvantages that that entails, notably
  expensive function call overhead, more complex memory management and trickier
  cross-compilation.
* `go-geom` uses a cache-friendly coordinate layout which is generally faster
  than GEOS for many operations.

`go-geos` is a good fit if your program is short-lived (meaning you can ignore
memory management), or you require the battle-tested geometry functions provided
by GEOS and are willing to handle memory management manually. `go-geom` is
recommended for long-running processes with less stringent geometry function
requirements.

## License

MIT