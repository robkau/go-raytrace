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

func (c *Cube) Bounds() Bounds {
	return newBounds(geom.NewPoint(-1, -1, -1), geom.NewPoint(1, 1, 1)).TransformTo(c.t)
}

func (c *Cube) NormalAt(p geom.Tuple) geom.Tuple {
	return NormalAt(c, p, c.LocalNormalAt)
}

func (c *Cube) LocalNormalAt(p geom.Tuple) geom.Tuple {
	maxC := math.Max(math.Max(math.Abs(p.X), math.Abs(p.Y)), math.Abs(p.Z))

	if maxC == math.Abs(p.X) {
		return geom.NewVector(p.X, 0, 0)
	} else if maxC == math.Abs(p.Y) {
		return geom.NewVector(0, p.Y, 0)
	} else {
		return geom.NewVector(0, 0, p.Z)
	}
}

func (c *Cube) Intersect(r geom.Ray) Intersections {
	return Intersect(r, c.t, c.LocalIntersect)
}

func (c *Cube) LocalIntersect(r geom.Ray) Intersections {
	// bounds for a unit cube because we are in local space
	tMin, tMax := intersectsCube(r, newBounds(geom.NewPoint(-1, -1, -1), geom.NewPoint(1, 1, 1)))
	if tMin > tMax {
		return NewIntersections()
	}
	return NewIntersections(
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

func intersectsCube(r geom.Ray, b Bounds) (tMin float64, tMax float64) {
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
