package geos

// defaultContext is the default context.
var defaultContext = NewContext()

// Clone clones g into c.
func Clone(g *Geom) *Geom {
	return defaultContext.Clone(g)
}

// NewGeomFromBounds returns a new polygon populated with bounds.
func NewGeomFromBounds(bounds *Bounds) *Geom {
	return defaultContext.NewGeomFromBounds(bounds)
}

// NewCollection returns a new collection.
func NewCollection(typeID GeometryTypeID, geoms []*Geom) *Geom {
	return defaultContext.NewCollection(typeID, geoms)
}

// NewCoordSeq returns a new CoordSeq.
func NewCoordSeq(size, dims int) *CoordSeq {
	return defaultContext.NewCoordSeq(size, dims)
}

// NewCoordSeqFromCoords returns a new CoordSeq populated with coords.
func NewCoordSeqFromCoords(coords [][]float64) *CoordSeq {
	return defaultContext.NewCoordSeqFromCoords(coords)
}

// NewEmptyCollection returns a new empty collection.
func NewEmptyCollection(typeID GeometryTypeID) *Geom {
	return defaultContext.NewEmptyCollection(typeID)
}

// NewEmptyLineString returns a new empty line string.
func NewEmptyLineString() *Geom {
	return defaultContext.NewEmptyLineString()
}

// NewEmptyPoint returns a new empty point.
func NewEmptyPoint() *Geom {
	return defaultContext.NewEmptyPoint()
}

// NewEmptyPolygon returns a new empty polygon.
func NewEmptyPolygon() *Geom {
	return defaultContext.NewEmptyPolygon()
}

// NewGeomFromWKB parses a geometry in WKB format from wkb.
func NewGeomFromWKB(wkb []byte) (*Geom, error) {
	return defaultContext.NewGeomFromWKB(wkb)
}

// NewGeomFromWKT parses a geometry in WKT format from wkt.
func NewGeomFromWKT(wkt string) (*Geom, error) {
	return defaultContext.NewGeomFromWKT(wkt)
}

// NewLinearRing returns a new linear ring populated with coords.
func NewLinearRing(coords [][]float64) *Geom {
	return defaultContext.NewLinearRing(coords)
}

// NewLineString returns a new line string populated with coords.
func NewLineString(coords [][]float64) *Geom {
	return defaultContext.NewLineString(coords)
}

// NewPoint returns a new point populated with coord.
func NewPoint(coord []float64) *Geom {
	return defaultContext.NewPoint(coord)
}

// NewPolygon returns a new point populated with coordss.
func NewPolygon(coordss [][][]float64) *Geom {
	return defaultContext.NewPolygon(coordss)
}
