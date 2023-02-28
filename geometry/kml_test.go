package geometry_test

import (
	"encoding/xml"

	"github.com/twpayne/go-geos/geometry"
)

var _ xml.Marshaler = &geometry.Geometry{}
