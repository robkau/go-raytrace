package shapes

import (
	"github.com/robkau/go-raytrace/lib/geom"
	"math"
)

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

func NewEmptyBoundingBox() *BoundingBox {
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

func (b *BoundingBox) Transform(t *geom.X4Matrix) {
	p1 := b.Min
	p2 := geom.NewPoint(b.Min.X, b.Min.Y, b.Max.Z)
	p3 := geom.NewPoint(b.Min.X, b.Max.Y, b.Min.Z)
	p4 := geom.NewPoint(b.Min.X, b.Max.Y, b.Max.Z)
	p5 := geom.NewPoint(b.Max.X, b.Min.Y, b.Min.Z)
	p6 := geom.NewPoint(b.Max.X, b.Min.Y, b.Max.Z)
	p7 := geom.NewPoint(b.Max.X, b.Max.Y, b.Min.Z)
	p8 := b.Max

	newB := NewEmptyBoundingBox()
	newB.Add(
		t.MulTuple(p1),
		t.MulTuple(p2),
		t.MulTuple(p3),
		t.MulTuple(p4),
		t.MulTuple(p5),
		t.MulTuple(p6),
		t.MulTuple(p7),
		t.MulTuple(p8),
	)
	b.Min = newB.Min
	b.Max = newB.Max
}

func (b *BoundingBox) AddBoundingBoxes(boxes ...*BoundingBox) {
	for _, box := range boxes {
		b.Add(box.Min)
		b.Add(box.Max)
	}
}

func (b *BoundingBox) Contains(p geom.Tuple) bool {
	return p.X >= b.Min.X && p.X <= b.Max.X &&
		p.Y >= b.Min.Y && p.Y <= b.Max.Y &&
		p.Z >= b.Min.Z && p.Z <= b.Max.Z
}

func (b *BoundingBox) ContainsBox(box *BoundingBox) bool {
	return b.Contains(box.Min) && b.Contains(box.Max)
}

func (b *BoundingBox) Intersect(r geom.Ray) bool {
	tMin, tMax := intersectsCube(r, b)
	if tMin > tMax {
		return false
	}
	return true
}

func (b *BoundingBox) SplitBounds() (left, right *BoundingBox) {
	dx := b.Max.X - b.Min.X
	dy := b.Max.Y - b.Min.Y
	dz := b.Max.Z - b.Min.Z

	greatest := math.Max(math.Max(dx, dy), dz)

	x0 := b.Min.X
	y0 := b.Min.Y
	z0 := b.Min.Z
	x1 := b.Max.X
	y1 := b.Max.Y
	z1 := b.Max.Z

	if greatest == dx {
		x0 = x0 + dx/2.0
		x1 = x0
	} else if greatest == dy {
		y0 = y0 + dy/2.0
		y1 = y0
	} else {
		z0 = z0 + dz/2.0
		z1 = z0
	}

	midMin := geom.NewPoint(x0, y0, z0)
	midMax := geom.NewPoint(x1, y1, z1)

	left = NewBoundingBox(b.Min, midMax)
	right = NewBoundingBox(midMin, b.Max)
	return
}

func (b *BoundingBox) Center() geom.Tuple {
	return b.Min.Add(b.Max).Div(2)
}
