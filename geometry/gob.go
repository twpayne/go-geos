package geometry

import "github.com/twpayne/go-geos"

// GobDecode implements encoding/gob.GobDecoder.
func (g *Geometry) GobDecode(data []byte) error {
	if len(data) == 0 {
		g.Geom = geos.NewEmptyPoint()
		return nil
	}
	var err error
	g.Geom, err = geos.NewGeomFromWKB(data)
	return err
}

// GobEncode implements encoding/gob.GobEncoder.
func (g *Geometry) GobEncode() ([]byte, error) {
	if g.Geom == nil {
		return nil, nil
	}
	return g.Geom.ToEWKBWithSRID(), nil
}
