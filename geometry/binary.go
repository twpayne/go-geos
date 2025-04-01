package geometry

import "github.com/twpayne/go-geos"

// NewGeometryFromWKB returns a new Geometry from wkb.
func NewGeometryFromWKB(wkb []byte) (*Geometry, error) {
	geom, err := geos.NewGeomFromWKB(wkb)
	if err != nil {
		return nil, err
	}
	return &Geometry{Geom: geom}, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (g *Geometry) MarshalBinary() ([]byte, error) {
	return g.ToEWKBWithSRID(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (g *Geometry) UnmarshalBinary(data []byte) error {
	geom, err := geos.NewGeomFromWKB(data)
	if err != nil {
		return err
	}
	g.Geom = geom
	return nil
}
