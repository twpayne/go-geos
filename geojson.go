package geos

// MarshalJSON implements encoding/json.Marshaler.MarshalJSON.
func (g *Geom) MarshalJSON() ([]byte, error) {
	return []byte(g.ToGeoJSON(0)), nil
}

// UnmarshalJSON implements encoding/json.Unmarshaler.UnmarshalJSON.
func (g *Geom) UnmarshalJSON(data []byte) error {
	context := g.context
	if context == nil {
		context = DefaultContext
	}
	geom, err := context.NewGeomFromGeoJSON(string(data))
	if err != nil {
		return err
	}
	g.Destroy()
	*g = *geom
	return nil
}
