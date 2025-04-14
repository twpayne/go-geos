package geometry

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/twpayne/go-geos"
)

var (
	geojsonType = map[geos.TypeID]string{
		geos.TypeIDPoint:              "Point",
		geos.TypeIDLineString:         "LineString",
		geos.TypeIDPolygon:            "Polygon",
		geos.TypeIDMultiPoint:         "MultiPoint",
		geos.TypeIDMultiLineString:    "MultiLineString",
		geos.TypeIDMultiPolygon:       "MultiPolygon",
		geos.TypeIDGeometryCollection: "GeometryCollection",
	}

	errUnsupportedEmptyGeometry = errors.New("unsupported empty geometry")
)

// NewGeometryFromGeoJSON returns a new Geometry parsed from GeoJSON.
func NewGeometryFromGeoJSON(geoJSON []byte) (*Geometry, error) {
	g := &Geometry{}
	if err := g.UnmarshalJSON(geoJSON); err != nil {
		return nil, err
	}
	return g, nil
}

// AsGeoJSON returns the GeoJSON representation of g.
func (g *Geometry) AsGeoJSON() ([]byte, error) {
	return g.MarshalJSON()
}

// MarshalJSON implements encoding/json.Marshaler.
func (g *Geometry) MarshalJSON() ([]byte, error) {
	sb := &strings.Builder{}
	sb.Grow(initialStringBufferSize)
	if err := geojsonWriteGeom(sb, g.Geom); err != nil {
		return nil, err
	}
	return []byte(sb.String()), nil
}

// UnmarshalJSON implements encoding/json.Unmarshaler.
func (g *Geometry) UnmarshalJSON(data []byte) error {
	var geoJSON struct {
		Type        string            `json:"type"`
		Coordinates json.RawMessage   `json:"coordinates"`
		Geometries  []json.RawMessage `json:"geometries"`
	}
	if err := json.Unmarshal(data, &geoJSON); err != nil {
		return err
	}
	switch geoJSON.Type {
	case "Point":
		var coordinates []float64
		if err := json.Unmarshal(geoJSON.Coordinates, &coordinates); err != nil {
			return err
		}
		g.Geom = geos.NewPoint(coordinates)
		return nil
	case "LineString":
		var coordinates [][]float64
		if err := json.Unmarshal(geoJSON.Coordinates, &coordinates); err != nil {
			return err
		}
		g.Geom = geos.NewLineString(coordinates)
		return nil
	case "Polygon":
		var coordinates [][][]float64
		if err := json.Unmarshal(geoJSON.Coordinates, &coordinates); err != nil {
			return err
		}
		g.Geom = geos.NewPolygon(coordinates)
		return nil
	case "MultiPoint":
		var coordinates [][]float64
		if err := json.Unmarshal(geoJSON.Coordinates, &coordinates); err != nil {
			return err
		}
		geoms := make([]*geos.Geom, len(coordinates))
		for i, pointCoord := range coordinates {
			geoms[i] = geos.NewPoint(pointCoord)
		}
		g.Geom = geos.NewCollection(geos.TypeIDMultiPoint, geoms)
		return nil
	case "MultiLineString":
		var coordinates [][][]float64
		if err := json.Unmarshal(geoJSON.Coordinates, &coordinates); err != nil {
			return err
		}
		geoms := make([]*geos.Geom, len(coordinates))
		for i, lineStringCoords := range coordinates {
			geoms[i] = geos.NewLineString(lineStringCoords)
		}
		g.Geom = geos.NewCollection(geos.TypeIDMultiLineString, geoms)
		return nil
	case "MultiPolygon":
		var coordinates [][][][]float64
		if err := json.Unmarshal(geoJSON.Coordinates, &coordinates); err != nil {
			return err
		}
		geoms := make([]*geos.Geom, len(coordinates))
		for i, polygonCoords := range coordinates {
			geoms[i] = geos.NewPolygon(polygonCoords)
		}
		g.Geom = geos.NewCollection(geos.TypeIDMultiPolygon, geoms)
		return nil
	case "MultiGeometry":
		fallthrough // FIXME handle MultiGeometry
	default:
		return fmt.Errorf("unsupported type: %s", geoJSON.Type)
	}
}

func geojsonWriteCoordinates(sb *strings.Builder, geom *geos.Geom) error {
	for i, coord := range geom.CoordSeq().ToCoords() {
		if i != 0 {
			if err := sb.WriteByte(','); err != nil {
				return err
			}
		}
		if err := sb.WriteByte('['); err != nil {
			return err
		}
		for j, ord := range coord {
			if j != 0 {
				if err := sb.WriteByte(','); err != nil {
					return err
				}
			}
			if _, err := sb.WriteString(strconv.FormatFloat(ord, 'f', -1, 64)); err != nil {
				return err
			}
		}
		if err := sb.WriteByte(']'); err != nil {
			return err
		}
	}
	return nil
}

func geojsonWriteCoordinatesArray(sb *strings.Builder, geom *geos.Geom) error {
	if err := sb.WriteByte('['); err != nil {
		return err
	}
	if err := geojsonWriteCoordinates(sb, geom); err != nil {
		return err
	}
	return sb.WriteByte(']')
}

func geojsonWriteGeom(sb *strings.Builder, geom *geos.Geom) error {
	if geom == nil {
		_, err := sb.WriteString("null")
		return err
	}
	typ, ok := geojsonType[geom.TypeID()]
	if !ok {
		return fmt.Errorf("unsupported type: %s", geom.Type())
	}
	if _, err := sb.WriteString(`{"type":"` + typ + `"`); err != nil {
		return err
	}
	//nolint:exhaustive
	switch geom.TypeID() {
	case geos.TypeIDPoint:
		if geom.IsEmpty() {
			return errUnsupportedEmptyGeometry
		}
		if _, err := sb.WriteString(`,"coordinates":`); err != nil {
			return err
		}
		if err := geojsonWriteCoordinates(sb, geom); err != nil {
			return err
		}
	case geos.TypeIDLineString:
		if geom.IsEmpty() {
			return errUnsupportedEmptyGeometry
		}
		if _, err := sb.WriteString(`,"coordinates":`); err != nil {
			return err
		}
		if err := geojsonWriteCoordinatesArray(sb, geom); err != nil {
			return err
		}
	case geos.TypeIDPolygon:
		if geom.IsEmpty() {
			return errUnsupportedEmptyGeometry
		}
		if _, err := sb.WriteString(`,"coordinates":`); err != nil {
			return err
		}
		if err := geojsonWritePolygonCoordinates(sb, geom); err != nil {
			return err
		}
	case geos.TypeIDMultiPoint:
		if _, err := sb.WriteString(`,"coordinates":[`); err != nil {
			return err
		}
		for i, n := 0, geom.NumGeometries(); i < n; i++ {
			if i != 0 {
				if err := sb.WriteByte(','); err != nil {
					return err
				}
			}
			if err := geojsonWriteCoordinates(sb, geom.Geometry(i)); err != nil {
				return err
			}
		}
		if err := sb.WriteByte(']'); err != nil {
			return err
		}
	case geos.TypeIDMultiLineString:
		if _, err := sb.WriteString(`,"coordinates":[`); err != nil {
			return err
		}
		for i, n := 0, geom.NumGeometries(); i < n; i++ {
			if i != 0 {
				if err := sb.WriteByte(','); err != nil {
					return err
				}
			}
			if err := geojsonWriteCoordinatesArray(sb, geom.Geometry(i)); err != nil {
				return err
			}
		}
		if err := sb.WriteByte(']'); err != nil {
			return err
		}
	case geos.TypeIDMultiPolygon:
		if _, err := sb.WriteString(`,"coordinates":[`); err != nil {
			return err
		}
		for i, n := 0, geom.NumGeometries(); i < n; i++ {
			if i != 0 {
				if err := sb.WriteByte(','); err != nil {
					return err
				}
			}
			if err := geojsonWritePolygonCoordinates(sb, geom.Geometry(i)); err != nil {
				return err
			}
		}
		if err := sb.WriteByte(']'); err != nil {
			return err
		}
	case geos.TypeIDGeometryCollection:
		if _, err := sb.WriteString(`,"geometries":[`); err != nil {
			return err
		}
		for i, n := 0, geom.NumGeometries(); i < n; i++ {
			if err := geojsonWriteGeom(sb, geom.Geometry(i)); err != nil {
				return err
			}
		}
		if err := sb.WriteByte(']'); err != nil {
			return err
		}
	}
	return sb.WriteByte('}')
}

func geojsonWritePolygonCoordinates(sb *strings.Builder, geom *geos.Geom) error {
	if err := sb.WriteByte('['); err != nil {
		return err
	}
	if err := geojsonWriteCoordinatesArray(sb, geom.ExteriorRing()); err != nil {
		return err
	}
	for i, n := 0, geom.NumInteriorRings(); i < n; i++ {
		if err := sb.WriteByte(','); err != nil {
			return err
		}
		if err := geojsonWriteCoordinatesArray(sb, geom.InteriorRing(i)); err != nil {
			return err
		}
	}
	return sb.WriteByte(']')
}
