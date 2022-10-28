// Package geos provides an interface to GEOS. See https://trac.osgeo.org/geos/.
package geos

// #cgo LDFLAGS: -lgeos_c
// #include "geos.h"
import "C"

// Version.
const (
	VersionMajor = C.GEOS_VERSION_MAJOR
	VersionMinor = C.GEOS_VERSION_MINOR
	VersionPatch = C.GEOS_VERSION_PATCH
)

// A GeometryTypeID is a geometry type id.
type GeometryTypeID int

// Geometry type ids.
const (
	PointTypeID              GeometryTypeID = C.GEOS_POINT
	LineStringTypeID         GeometryTypeID = C.GEOS_LINESTRING
	LinearRingTypeID         GeometryTypeID = C.GEOS_LINEARRING
	PolygonTypeID            GeometryTypeID = C.GEOS_POLYGON
	MultiPointTypeID         GeometryTypeID = C.GEOS_MULTIPOINT
	MultiLineStringTypeID    GeometryTypeID = C.GEOS_MULTILINESTRING
	MultiPolygonTypeID       GeometryTypeID = C.GEOS_MULTIPOLYGON
	GeometryCollectionTypeID GeometryTypeID = C.GEOS_GEOMETRYCOLLECTION
)

// An Error is an error returned by GEOS.
type Error string

func (e Error) Error() string {
	return string(e)
}

var (
	errDimensionOutOfRange = Error("dimension out of range")
	errIndexOutOfRange     = Error("index out of range")
)

// versionEqualOrGreaterThan returns true if the GEOS version is at least
// major.minor.patch.
func versionEqualOrGreaterThan(major, minor, patch int) bool {
	switch {
	case VersionMajor > major:
		return true
	case VersionMajor < major:
		return false
	}
	switch {
	case VersionMinor > minor:
		return true
	case VersionMinor < minor:
		return false
	}
	return VersionPatch >= patch
}
