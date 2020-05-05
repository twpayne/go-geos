package geometry

import (
	"database/sql/driver"
	"fmt"

	"github.com/twpayne/go-geos"
)

// Scan implements sql.Scanner.
func (g *Geometry) Scan(src interface{}) error {
	switch src := src.(type) {
	case nil:
		g.Geom = nil
		return nil
	case []byte:
		if len(src) == 0 {
			g.Geom = geos.NewEmptyPoint()
			return nil
		}
		var err error
		g.Geom, err = geos.NewGeomFromWKB(src)
		return err
	default:
		return fmt.Errorf("want []byte, got %T", src)
	}
}

// Value implements driver.Value.
func (g *Geometry) Value() (driver.Value, error) {
	if g.Geom == nil {
		return nil, nil
	}
	return g.Geom.ToWKB(), nil
}
