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

// PrecisionRules.
const (
	// The output is always valid. Collapsed geometry elements
	// (including both polygons and lines) are removed.
	PrecisionRuleValidOutput PrecisionRule = C.GEOS_PREC_VALID_OUTPUT
	// Precision reduction is performed pointwise. Output geometry may
	// be invalid due to collapse or self-intersection.
	PrecisionRulePointwise PrecisionRule = C.GEOS_PREC_NO_TOPO
	// Like the default mode, except that collapsed linear geometry
	// elements are preserved. Collapsed polygonal input elements are
	// removed.
	PrecisionRuleKeepCollapsed PrecisionRule = C.GEOS_PREC_KEEP_COLLAPSED
)

func PrecisionRulesOr(flags []PrecisionRule) C.int {
	result := 0
	for _, rule := range flags {
		result |= int(rule)
	}
	return C.int(result)
}

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
