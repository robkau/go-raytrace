package shapes

import (
	"github.com/robkau/go-raytrace/lib/geom"
	"math"
)

type Cone struct {
	baseShape

	maximum float64
	minimum float64
	capped  bool
}

func NewCone(min, max float64, capped bool) *Cone {
	return &Cone{
		baseShape: newBaseShape(),
		minimum:   min,
		maximum:   max,
		capped:    capped,
	}
}

func NewUnitCone(capped bool) *Cone {
	return NewCone(-1, 0, capped)
}

func NewInfiniteCone() *Cone {
	return NewCone(math.Inf(-1), math.Inf(1), false)
}

// BoundsOf is for untransformed shape
func (c *Cone) BoundsOf() *BoundingBox {
	a := math.Abs(c.minimum)
	b := math.Abs(c.maximum)
	limit := math.Max(a, b)

	return NewBoundingBox(geom.NewPoint(-limit, c.minimum, -limit),
		geom.NewPoint(limit, c.maximum, limit))
}

func (c *Cone) NormalAt(p geom.Tuple, _ Intersection) geom.Tuple {
	return NormalAt(c, p, c.LocalNormalAt, Intersection{})
}

func (c *Cone) LocalNormalAt(p geom.Tuple, _ Intersection) geom.Tuple {
	d := p.X*p.X + p.Z*p.Z

	if d < 1 && p.Y >= c.maximum-geom.FloatComparisonEpsilon {
		return geom.NewVector(0, 1, 0)
	}

	if d < 1 && p.Y <= c.minimum+geom.FloatComparisonEpsilon {
		return geom.NewVector(0, -1, 0)
	}

	y := math.Sqrt(p.X*p.X + p.Z*p.Z)
	if p.Y > 0 {
		y = -y
	}

	return geom.NewVector(p.X, y, p.Z)
}

func (c *Cone) Intersect(r geom.Ray) *Intersections {
	return Intersect(r, c.t, c.LocalIntersect)
}

func (c *Cone) LocalIntersect(r geom.Ray) *Intersections {
	xs := NewIntersections()

	a := r.Direction.X*r.Direction.X - r.Direction.Y*r.Direction.Y + r.Direction.Z*r.Direction.Z
	b := 2*r.Origin.X*r.Direction.X - 2*r.Origin.Y*r.Direction.Y + 2*r.Origin.Z*r.Direction.Z
	cc := r.Origin.X*r.Origin.X - r.Origin.Y*r.Origin.Y + r.Origin.Z*r.Origin.Z

	if geom.AlmostEqual(0, a) {
		// ray parallel to one of cone sides

		if b != 0 {
			// it hits the other side
			xs.Add(NewIntersection(-cc/(2*b), c))
		}
	} else {
		// check sides
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

func (c *Cone) intersectCaps(r geom.Ray, xs *Intersections) *Intersections {
	if !c.capped || geom.AlmostEqual(r.Direction.Y, 0) {
		return xs
	}

	// check lower cap
	t := (c.minimum - r.Origin.Y) / r.Direction.Y
	if checkConeCap(r, t, c.minimum) {
		xs.Add(NewIntersection(t, c))
	}

	// check upper cap
	t = (c.maximum - r.Origin.Y) / r.Direction.Y
	if checkConeCap(r, t, c.maximum) {
		xs.Add(NewIntersection(t, c))
	}

	return xs
}

func checkConeCap(r geom.Ray, t float64, y float64) bool {
	x := r.Origin.X + t*r.Direction.X
	z := r.Origin.Z + t*r.Direction.Z

	// radius of cap is always equal to current y value
	return (x*x + z*z) <= math.Abs(y)
}
