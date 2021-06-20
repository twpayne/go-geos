package geometry

import (
	"database/sql"
	"database/sql/driver"
)

var (
	_ driver.Value = &Geometry{}
	_ sql.Scanner  = &Geometry{}
)
