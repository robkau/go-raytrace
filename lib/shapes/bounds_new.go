package shapes

import "github.com/robkau/go-raytrace/lib/geom"

type BoundingBox struct {
	Min geom.Tuple
	Max geom.Tuple
}

func NewBoundingBox(min, max geom.Tuple) *BoundingBox {
	if min == geom.ZeroPoint() && max == geom.ZeroPoint() {
		return &BoundingBox{
			Min: geom.PosInfPoint(),
			Max: geom.NegInfPoint(),
		}
	}

	return &BoundingBox{
		Min: min,
		Max: max,
	}
}
