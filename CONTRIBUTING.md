# Contributing

`go-geos` uses `libgeos`'s stable C API,
[`geos_c.h`](http://libgeos.org/doxygen/geos__c_8h.html). Only the thread-safe
`*_r` functions are used.

## Adding methods to `*Geom`

Wherever possible, `go-geos` uses code generation to generate wrappers for
`*Geom` methods. The generated code is in
[`geommethods.go`](https://github.com/twpayne/go-geos/blob/master/geommethods.go).

There are five parts to this:

* [`geommethods.yaml`](https://github.com/twpayne/go-geos/blob/master/geommethods.yaml)
  contains the high-level definitions of the methods.
* [`geommethods.go.tmpl`](https://github.com/twpayne/go-geos/blob/master/geommethods.go.tmpl)
  is a `text/template` template that is executed with the data from
  `geommethods.yaml`.
* [`internal/cmds/execute-template/`](https://github.com/twpayne/go-geos/tree/master/internal/cmds/execute-template)
  executes a template with data and includes custom template functions.
* [`go generate`](https://go.dev/blog/generate) runs
  `internal/cmds/execute-template/` with `geommethods.yaml` and
  `geommethods.go.tmpl` as inputs and writes `geommethods.go`.
* [`geom_test.go`](https://github.com/twpayne/go-geos/blob/master/geom_test.go)
  contains unit tests to ensure that the method is wrapped correctly.

Adding a method to `*Geom` consists of one or more steps, depending on how
similar the method is to existing methods:

1. In simple cases, adding a few lines to `geommethods.yaml` and running `go
   generate` is sufficient. You will need to add a test to `geom_test.go`.
2. For more complex cases, you might have to modify or extend
   `geommethods.go.tmpl`.
3. If you need to add or modify a template function, you will need to modify
   `internal/cmds/execute-template/`.
   
## Maintaining backwards compatibility

`go-geos` supports all `libgeos` versions used in [Ubuntu LTS
releases](https://ubuntu.com/about/release-cycle). These are:

| Ubuntu Version | `libgeos` version |
| -------------- | ----------------- |
| 18.04 LTS      | 3.6.2             |
| 20.04 LTS      | 3.8.0             |
| 22.04 LTS      | 3.10.2            |

As `libgeos` is under active development, bugs are fixed and new features are
added over time. This causes problems when old versions might behave incorrectly
or miss newly-added features. In these cases:

* In general, it is the user's responsibility to ensure that they are using a
  sufficiently recent version of `libgeos` for their needs. `go-geos` can
  forward incorrect results from `libgeos` and behave in an undefined manner
  (including crashing the program) when missing features are invoked.
* For tests, use the `versionEqualToOrGreaterThan` function to only run the test
  when the functionality is both present and working correctly. It's OK to skip
  tests entirely if the functionality is buggy or not present.
* For new features, you will need to add stubs in
  [`geos.c`](https://github.com/twpayne/go-geos/blob/master/geos.c) and
  [`geos.h`](https://github.com/twpayne/go-geos/blob/master/geos.h) to provide
  the function when it is not provided. This is a necessity to permit `go-geos`
  to compile with earlier versions of `libgeos`. For an example of this, see how
  `GEOSDistanceIndexed_r` is stubbed.
  

## C code formatting

`go-geos` uses [`clang-format`](https://clang.llvm.org/docs/ClangFormat.html) to
format C code. You can run this with:

```console
$ clang-format -i *.c *.h
```