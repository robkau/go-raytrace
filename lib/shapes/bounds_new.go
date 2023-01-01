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

func EmptyBoundingBox() *BoundingBox {
	return &BoundingBox{
		Min: geom.PosInfPoint(),
		Max: geom.NegInfPoint(),
	}
}

func (b *BoundingBox) Add(points ...geom.Tuple) {
	for _, p := range points {
		if p.X < b.Min.X {
			b.Min.X = p.X
		}
		if p.Y < b.Min.Y {
			b.Min.Y = p.Y
		}
		if p.Z < b.Min.Z {
			b.Min.Z = p.Z
		}
		if p.X > b.Max.X {
			b.Max.X = p.X
		}
		if p.Y > b.Max.Y {
			b.Max.Y = p.Y
		}
		if p.Z > b.Max.Z {
			b.Max.Z = p.Z
		}
	}
}
