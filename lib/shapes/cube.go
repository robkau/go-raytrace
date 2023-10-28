package shapes

import (
	"github.com/robkau/go-raytrace/lib/geom"
	"math"
)

type Cube struct {
	baseShape
}

func NewCube() *Cube {
	return &Cube{
		baseShape: newBaseShape(),
	}
}

// BoundsOf is for untransformed shape
func (c *Cube) BoundsOf() *BoundingBox {
	return NewBoundingBox(geom.NewPoint(-1, -1, -1), geom.NewPoint(1, 1, 1))
}

func (c *Cube) NormalAt(p geom.Tuple, _ Intersection) geom.Tuple {
	return NormalAt(c, p, c.LocalNormalAt, Intersection{})
}

func (c *Cube) LocalNormalAt(p geom.Tuple, _ Intersection) geom.Tuple {
	maxC := math.Max(math.Max(math.Abs(p.X), math.Abs(p.Y)), math.Abs(p.Z))

	if maxC == math.Abs(p.X) {
		return geom.NewVector(p.X, 0, 0)
	} else if maxC == math.Abs(p.Y) {
		return geom.NewVector(0, p.Y, 0)
	} else {
		return geom.NewVector(0, 0, p.Z)
	}
}

func (c *Cube) Intersect(r geom.Ray, i *Intersections) {
	Intersect(r, i, c.t, c.LocalIntersect)
}

func (c *Cube) LocalIntersect(r geom.Ray, i *Intersections) {
	// bounds for a unit cube because we are in local space
	tMin, tMax := intersectsCube(r, c.BoundsOf())
	if tMin > tMax {
		return
	}
	i.Add(
		Intersection{
			T: tMin,
			O: c,
		},
		Intersection{
			T: tMax,
			O: c,
		},
	)
}

func intersectsCube(r geom.Ray, b *BoundingBox) (tMin float64, tMax float64) {
	xtMin, xtMax := checkAxis(r.Origin.X, r.Direction.X, b.Min.X, b.Max.X)
	ytMin, ytMax := checkAxis(r.Origin.Y, r.Direction.Y, b.Min.Y, b.Max.Y)
	ztMin, ztMax := checkAxis(r.Origin.Z, r.Direction.Z, b.Min.Z, b.Max.Z)

	tMin = math.Max(math.Max(xtMin, ytMin), ztMin)
	tMax = math.Min(math.Min(xtMax, ytMax), ztMax)

	return tMin, tMax
}

func checkAxis(origin, direction, axisMin, axisMax float64) (tMin float64, tMax float64) {
	tMinNumerator := axisMin - origin
	tMaxNumerator := axisMax - origin

	if math.Abs(direction) >= geom.FloatComparisonEpsilon {
		tMin = tMinNumerator / direction
		tMax = tMaxNumerator / direction
	} else {
		tMin = tMinNumerator * math.Inf(1)
		tMax = tMaxNumerator * math.Inf(1)
	}

	if tMin > tMax {
		// swap
		return tMax, tMin
	}
	// no swap
	return
}
