package geometry

import (
	"database/sql/driver"
	"encoding/hex"
	"fmt"

	"github.com/twpayne/go-geos"
)

// Scan implements database/sql.Scanner.
func (g *Geometry) Scan(src interface{}) error {
	switch src := src.(type) {
	case nil:
		g.Geom = nil
		return nil
	case []byte:
		return g.scanWKB(src)
	case string:
		wkb, err := hex.DecodeString(src)
		if err != nil {
			return err
		}
		return g.scanWKB(wkb)
	default:
		return fmt.Errorf("want nil, []byte, or string, got %T", src)
	}
}

func (g *Geometry) scanWKB(wkb []byte) error {
	if len(wkb) == 0 {
		g.Geom = geos.NewEmptyPoint()
		return nil
	}
	geom, err := geos.NewGeomFromWKB(wkb)
	if err != nil {
		return err
	}
	g.Geom = geom
	return nil
}

// Value implements database/sql/driver.Value.
func (g Geometry) Value() (driver.Value, error) {
	if g.Geom == nil {
		return nil, nil //nolint:nilnil
	}
	return hex.EncodeToString(g.Geom.ToEWKBWithSRID()), nil
}
