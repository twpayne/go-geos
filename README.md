# geos

Package `geos` provides an interface to [GEOS](https://trac.osgeo.org/geos).

## Features

* Idiomatic Go API.

* Concurrency-safe. `geos` uses GEOS's threadsafe `*_r` functions under the
  hood, with locking to ensure safety, even when used across multiple
  goroutines. For best performance, use one `geos.Context` per goroutine.

* Automatic finalization of GEOS objects. GEOS objects are finalized by Go's
  garbage collector when they are no longer referenced. To force finalization of
  all objects, call `runtime.GC()`.

## Exceptions

`geos` uses the stable C GEOS bindings. These bindings catch exceptions from the
underlying C++ code and convert them to a return code. For normal geometry
operations, `geos` `panic`s whenever it encounters a return code indicating an
error, rather than returning an `error`. This behaviour is similar to slice
access in Go (out-of-bounds accesses `panic`) and keeps the API easy to use.
When parsing WKB and WKT, errors are expected so an `error` is returned.