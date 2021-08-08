package shapes

import (
	"github.com/robkau/go-raytrace/lib/geom"
	"math"
)

type Bounds struct {
	Min geom.Tuple
	Max geom.Tuple
	set bool
}

func newBounds(points ...geom.Tuple) Bounds {
	b := Bounds{}
	return b.Add(points...)
}

func (b Bounds) Add(points ...geom.Tuple) Bounds {
	for _, p := range points {
		if !b.set {
			b.set = true
			b.Min = p
			b.Max = p
		}
		b.Min.X = math.Min(b.Min.X, p.X)
		b.Min.Y = math.Min(b.Min.Y, p.Y)
		b.Min.Z = math.Min(b.Min.Z, p.Z)
		b.Max.X = math.Max(b.Max.X, p.X)
		b.Max.Y = math.Max(b.Max.Y, p.Y)
		b.Max.Z = math.Max(b.Max.Z, p.Z)
	}
	return b
}

func (b Bounds) Center() geom.Tuple {
	centerX := (b.Min.X + b.Max.X) / 2
	centerY := (b.Min.Y + b.Max.Y) / 2
	centerZ := (b.Min.Z + b.Max.Z) / 2
	return geom.NewPoint(centerX, centerY, centerZ)
}

func (b Bounds) TransformTo(t *geom.X4Matrix) Bounds {
	// get corners of the bounding cube
	// lowercase negative direction, uppercase positive direction
	bxyz := geom.NewPoint(b.Min.X, b.Min.Y, b.Min.Z)
	bxyZ := geom.NewPoint(b.Min.X, b.Min.Y, b.Max.Z)
	bxYz := geom.NewPoint(b.Min.X, b.Max.Y, b.Min.Z)
	bxYZ := geom.NewPoint(b.Min.X, b.Max.Y, b.Max.Z)
	bXyz := geom.NewPoint(b.Max.X, b.Min.Y, b.Min.Z)
	bXyZ := geom.NewPoint(b.Max.X, b.Min.Y, b.Max.Z)
	bXYz := geom.NewPoint(b.Max.X, b.Max.Y, b.Min.Z)
	bXYZ := geom.NewPoint(b.Max.X, b.Max.Y, b.Max.Z)

	// transform to new space by multiplying each corner by childs transformation matrix
	bxyzT := t.MulTuple(bxyz)
	bxyZT := t.MulTuple(bxyZ)
	bxYzT := t.MulTuple(bxYz)
	bxYZT := t.MulTuple(bxYZ)
	bXyzT := t.MulTuple(bXyz)
	bXyZT := t.MulTuple(bXyZ)
	bXYzT := t.MulTuple(bXYz)
	bXYZT := t.MulTuple(bXYZ)

	// get new bounds - add each corner
	return newBounds(
		bxyzT,
		bxyZT,
		bxYzT,
		bxYZT,
		bXyzT,
		bXyZT,
		bXYzT,
		bXYZT,
	)
}

func (b Bounds) Intersects(r geom.Ray) bool {
	tMin, tMax := intersectsCube(r, b)
	if tMin > tMax {
		return false
	}
	return true
}
