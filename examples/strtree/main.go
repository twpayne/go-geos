package main

import (
	"fmt"
	"math/rand/v2"
	"os"

	"github.com/twpayne/go-geos"
)

type Item struct {
	id   int
	geom *geos.Geom
}

func NewRandomItem(_range float64) *Item {
	return &Item{
		id:   rand.Int(),
		geom: geos.NewPointFromXY(_range*rand.Float64(), _range*rand.Float64()),
	}
}

func run() error {
	const (
		NumItems = 10000
		Range    = 100
	)

	strTree := geos.NewSTRtree(10)

	items := make([]*Item, NumItems)
	for i := range items {
		item := NewRandomItem(Range)
		items[i] = item
		strTree.Insert(item.geom, item)
	}

	randomItem := NewRandomItem(Range)
	nearestItem := strTree.NearestGeneric(randomItem, randomItem.geom, func(item1, item2 any) float64 {
		geom1 := item1.(*Item).geom
		geom2 := item2.(*Item).geom
		return geom1.Distance(geom2)
	})
	fmt.Printf(" Random point: %+v\n", randomItem)
	fmt.Printf("Nearest Point: %+v\n", nearestItem)

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
