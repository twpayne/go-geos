package geos_test

import (
	"fmt"

	"github.com/twpayne/go-geos"
)

func ExamplePrepGeom() {
	geom, err := geos.NewGeomFromWKT("POLYGON ((189 115, 200 170, 130 170, 35 242, 156 215, 210 290, 274 256, 360 190, 267 215, 300 50, 200 60, 189 115))")
	if err != nil {
		panic(err)
	}
	prepGeom := geom.Prepare()
	point := geos.NewPointFromXY(190, 200)
	if prepGeom.Intersects(point) {
		fmt.Println("intersects")
	}
	// Output: intersects
}
