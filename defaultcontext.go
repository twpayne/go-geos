package geos

// DefaultContext is the default context.
var DefaultContext = NewContext()

// Clone clones g into c.
func Clone(g *Geom) *Geom {
	return DefaultContext.Clone(g)
}

// NewGeomFromBounds returns a new polygon populated with bounds.
func NewGeomFromBounds(minX, minY, maxX, maxY float64) *Geom {
	return DefaultContext.NewGeomFromBounds(minX, minY, maxX, maxY)
}

// NewCollection returns a new collection.
func NewCollection(typeID TypeID, geoms []*Geom) *Geom {
	return DefaultContext.NewCollection(typeID, geoms)
}

// NewCoordSeq returns a new CoordSeq.
func NewCoordSeq(size, dims int) *CoordSeq {
	return DefaultContext.NewCoordSeq(size, dims)
}

// NewCoordSeqFromCoords returns a new CoordSeq populated with coords.
func NewCoordSeqFromCoords(coords [][]float64) *CoordSeq {
	return DefaultContext.NewCoordSeqFromCoords(coords)
}

// NewEmptyCollection returns a new empty collection.
func NewEmptyCollection(typeID TypeID) *Geom {
	return DefaultContext.NewEmptyCollection(typeID)
}

// NewEmptyLineString returns a new empty line string.
func NewEmptyLineString() *Geom {
	return DefaultContext.NewEmptyLineString()
}

// NewEmptyPoint returns a new empty point.
func NewEmptyPoint() *Geom {
	return DefaultContext.NewEmptyPoint()
}

// NewEmptyPolygon returns a new empty polygon.
func NewEmptyPolygon() *Geom {
	return DefaultContext.NewEmptyPolygon()
}

// NewGEOMFromGeoJSON parses a geometry in GeoJSON format from GeoJSON.
func NewGeomFromGeoJSON(geoJSON string) (*Geom, error) {
	return DefaultContext.NewGeomFromGeoJSON(geoJSON)
}

// NewGeomFromWKB parses a geometry in WKB format from wkb.
func NewGeomFromWKB(wkb []byte) (*Geom, error) {
	return DefaultContext.NewGeomFromWKB(wkb)
}

// NewGeomFromWKT parses a geometry in WKT format from wkt.
func NewGeomFromWKT(wkt string) (*Geom, error) {
	return DefaultContext.NewGeomFromWKT(wkt)
}

// NewLinearRing returns a new linear ring populated with coords.
func NewLinearRing(coords [][]float64) *Geom {
	return DefaultContext.NewLinearRing(coords)
}

// NewLineString returns a new line string populated with coords.
func NewLineString(coords [][]float64) *Geom {
	return DefaultContext.NewLineString(coords)
}

// NewPoint returns a new point populated with coord.
func NewPoint(coord []float64) *Geom {
	return DefaultContext.NewPoint(coord)
}

// NewPointFromXY returns a new point with x and y.
func NewPointFromXY(x, y float64) *Geom {
	return DefaultContext.NewPointFromXY(x, y)
}

// NewPolygon returns a new point populated with coordss.
func NewPolygon(coordss [][][]float64) *Geom {
	return DefaultContext.NewPolygon(coordss)
}

// Polygonize returns a set of geometries which contains linework that
// represents the edges of a planar graph.
func Polygonize(geoms []*Geom) *Geom {
	return DefaultContext.Polygonize(geoms)
}

// PolygonizeValid returns a set of polygons which contains linework that
// represents the edges of a planar graph.
func PolygonizeValid(geoms []*Geom) *Geom {
	return DefaultContext.PolygonizeValid(geoms)
}
