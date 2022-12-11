package geometry

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"

	geos "github.com/twpayne/go-geos"
)

var (
	kmlPointStartElement           = xml.StartElement{Name: xml.Name{Local: "Point"}}
	kmlLineStringStartElement      = xml.StartElement{Name: xml.Name{Local: "LineString"}}
	kmlLinearRingStartElement      = xml.StartElement{Name: xml.Name{Local: "LinearRing"}}
	kmlPolygonStartElement         = xml.StartElement{Name: xml.Name{Local: "Polygon"}}
	kmlMultiGeometryStartElement   = xml.StartElement{Name: xml.Name{Local: "MultiGeometry"}}
	kmlCoordinatesStartElement     = xml.StartElement{Name: xml.Name{Local: "coordinates"}}
	kmlInnerBoundaryIsStartElement = xml.StartElement{Name: xml.Name{Local: "innerBoundaryIs"}}
	kmlOuterBoundaryIsStartElement = xml.StartElement{Name: xml.Name{Local: "outerBoundaryIs"}}
)

// MarshalXML implements encoding/xml.Marshaler.
func (g *Geometry) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return kmlEncodeGeom(e, g.Geom)
}

func kmlEncodeCoords(e *xml.Encoder, startElement xml.StartElement, geom *geos.Geom) error {
	if err := e.EncodeToken(startElement); err != nil {
		return err
	}
	if coords := geom.CoordSeq().ToCoords(); coords != nil {
		if err := e.EncodeToken(kmlCoordinatesStartElement); err != nil {
			return err
		}
		sb := &strings.Builder{}
		sb.Grow(initialStringBufferSize)
		for i, coord := range coords {
			if i != 0 {
				if err := sb.WriteByte(' '); err != nil {
					return err
				}
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
		}
		if err := e.EncodeToken(xml.CharData(sb.String())); err != nil {
			return err
		}
		if err := e.EncodeToken(kmlCoordinatesStartElement.End()); err != nil {
			return err
		}
	}
	return e.EncodeToken(startElement.End())
}

func kmlEncodeGeom(e *xml.Encoder, geom *geos.Geom) error {
	switch geom.TypeID() {
	case geos.TypeIDPoint:
		return kmlEncodeCoords(e, kmlPointStartElement, geom)
	case geos.TypeIDLineString:
		return kmlEncodeCoords(e, kmlLineStringStartElement, geom)
	case geos.TypeIDLinearRing:
		return kmlEncodeCoords(e, kmlLinearRingStartElement, geom)
	case geos.TypeIDPolygon:
		return kmlEncodePolygon(e, geom)
	case geos.TypeIDMultiPoint:
		fallthrough
	case geos.TypeIDMultiLineString:
		fallthrough
	case geos.TypeIDMultiPolygon:
		fallthrough
	case geos.TypeIDGeometryCollection:
		return kmlEncodeMultiGeometry(e, geom)
	default:
		return fmt.Errorf("unsupported type: %s", geom.Type())
	}
}

func kmlEncodeLinearRing(e *xml.Encoder, startElement xml.StartElement, geom *geos.Geom) error {
	if err := e.EncodeToken(startElement); err != nil {
		return err
	}
	if err := kmlEncodeCoords(e, kmlLinearRingStartElement, geom); err != nil {
		return err
	}
	return e.EncodeToken(startElement.End())
}

func kmlEncodeMultiGeometry(e *xml.Encoder, geom *geos.Geom) error {
	if err := e.EncodeToken(kmlMultiGeometryStartElement); err != nil {
		return err
	}
	for i, n := 0, geom.NumGeometries(); i < n; i++ {
		if err := kmlEncodeGeom(e, geom.Geometry(i)); err != nil {
			return err
		}
	}
	return e.EncodeToken(kmlMultiGeometryStartElement.End())
}

func kmlEncodePolygon(e *xml.Encoder, geom *geos.Geom) error {
	if err := e.EncodeToken(kmlPolygonStartElement); err != nil {
		return err
	}
	if !geom.IsEmpty() {
		if err := kmlEncodeLinearRing(e, kmlOuterBoundaryIsStartElement, geom.ExteriorRing()); err != nil {
			return err
		}
		for i, n := 0, geom.NumInteriorRings(); i < n; i++ {
			if err := kmlEncodeLinearRing(e, kmlInnerBoundaryIsStartElement, geom.InteriorRing(i)); err != nil {
				return err
			}
		}
	}
	return e.EncodeToken(kmlPolygonStartElement.End())
}
