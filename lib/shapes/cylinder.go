package shapes

import (
	"github.com/robkau/go-raytrace/lib/geom"
	"math"
)

type Cylinder struct {
	baseShape

	maximum float64
	minimum float64
	capped  bool
}

func NewCylinder(min, max float64, capped bool) *Cylinder {
	return &Cylinder{
		baseShape: newBaseShape(),
		minimum:   min,
		maximum:   max,
		capped:    capped,
	}
}

func NewInfiniteCylinder() *Cylinder {
	return NewCylinder(math.Inf(-1), math.Inf(1), false)
}

func (c *Cylinder) Bounds() Bounds {
	return newBounds(geom.NewPoint(-1, c.minimum, -1), geom.NewPoint(1, c.maximum, 1)).TransformTo(c.t)
}

func (c *Cylinder) NormalAt(p geom.Tuple, _ Intersection) geom.Tuple {
	return NormalAt(c, p, c.LocalNormalAt, Intersection{})
}

func (c *Cylinder) LocalNormalAt(p geom.Tuple, _ Intersection) geom.Tuple {
	d := p.X*p.X + p.Z*p.Z

	if d < 1 && p.Y >= c.maximum-geom.FloatComparisonEpsilon {
		return geom.NewVector(0, 1, 0)
	}

	if d < 1 && p.Y <= c.minimum+geom.FloatComparisonEpsilon {
		return geom.NewVector(0, -1, 0)
	}

	return geom.NewVector(p.X, 0, p.Z)
}

func (c *Cylinder) Intersect(r geom.Ray) Intersections {
	return Intersect(r, c.t, c.LocalIntersect)
}

func (c *Cylinder) LocalIntersect(r geom.Ray) Intersections {
	xs := NewIntersections()

	// check intersection with cylinder walls, if needed
	a := r.Direction.X*r.Direction.X + r.Direction.Z*r.Direction.Z
	if !geom.AlmostEqual(0, a) {
		// ray is not parallel up
		b := 2*r.Origin.X*r.Direction.X + 2*r.Origin.Z*r.Direction.Z
		cc := r.Origin.X*r.Origin.X + r.Origin.Z*r.Origin.Z - 1

		disc := b*b - 4*a*cc

		if disc < 0 {
			return NewIntersections()
		}

		t0 := (-b - math.Sqrt(disc)) / (2 * a)
		t1 := (-b + math.Sqrt(disc)) / (2 * a)

		if t0 > t1 {
			t0, t1 = t1, t0
		}

		y0 := r.Origin.Y + t0*r.Direction.Y
		if c.minimum < y0 && y0 < c.maximum {
			xs.Add(NewIntersection(t0, c))
		}

		y1 := r.Origin.Y + t1*r.Direction.Y
		if c.minimum < y1 && y1 < c.maximum {
			xs.Add(NewIntersection(t1, c))
		}
	}

	// check intersection with caps
	xs = c.intersectCaps(r, xs)

	return xs
}

func (c *Cylinder) intersectCaps(r geom.Ray, xs Intersections) Intersections {
	if !c.capped || geom.AlmostEqual(r.Direction.Y, 0) {
		return xs
	}

	// check lower cap
	t := (c.minimum - r.Origin.Y) / r.Direction.Y
	if checkCylinderCap(r, t) {
		xs.Add(NewIntersection(t, c))
	}

	// check upper cap
	t = (c.maximum - r.Origin.Y) / r.Direction.Y
	if checkCylinderCap(r, t) {
		xs.Add(NewIntersection(t, c))
	}

	return xs
}

func checkCylinderCap(r geom.Ray, t float64) bool {
	x := r.Origin.X + t*r.Direction.X
	z := r.Origin.Z + t*r.Direction.Z

	return (x*x + z*z) <= 1
}
