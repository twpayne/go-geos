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

// A TypeID is a geometry type id.
type TypeID int

// Geometry type ids.
const (
	TypeIDPoint              TypeID = C.GEOS_POINT
	TypeIDLineString         TypeID = C.GEOS_LINESTRING
	TypeIDLinearRing         TypeID = C.GEOS_LINEARRING
	TypeIDPolygon            TypeID = C.GEOS_POLYGON
	TypeIDMultiPoint         TypeID = C.GEOS_MULTIPOINT
	TypeIDMultiLineString    TypeID = C.GEOS_MULTILINESTRING
	TypeIDMultiPolygon       TypeID = C.GEOS_MULTIPOLYGON
	TypeIDGeometryCollection TypeID = C.GEOS_GEOMETRYCOLLECTION
)

type BufCapStyle int

// Buffer cap styles.
const (
	BufCapStyleRound  BufCapStyle = C.GEOSBUF_CAP_ROUND
	BufCapStyleFlat   BufCapStyle = C.GEOSBUF_CAP_FLAT
	BufCapStyleSquare BufCapStyle = C.GEOSBUF_CAP_SQUARE
)

type BufJoinStyle int

// Buffer join styles.
const (
	BufJoinStyleRound BufJoinStyle = C.GEOSBUF_JOIN_ROUND
	BufJoinStyleMitre BufJoinStyle = C.GEOSBUF_JOIN_MITRE
	BufJoinStyleBevel BufJoinStyle = C.GEOSBUF_JOIN_BEVEL
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

type PrecisionRule int

// Precision rules.
const (
	PrecisionRuleValidOutput   PrecisionRule = C.GEOS_PREC_VALID_OUTPUT
	PrecisionRulePointwise     PrecisionRule = C.GEOS_PREC_NO_TOPO
	PrecisionRuleKeepCollapsed PrecisionRule = C.GEOS_PREC_KEEP_COLLAPSED
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
