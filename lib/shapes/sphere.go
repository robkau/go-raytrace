package shapes

import (
	"go-raytrace/lib/geom"
	"math"
)

type Sphere struct {
	baseShape
}

func NewSphere() *Sphere {
	return &Sphere{
		baseShape: newBaseShape(),
	}
}

func (s *Sphere) NormalAt(p geom.Tuple) geom.Tuple {
	return NormalAt(s, p, s.LocalNormalAt)
}

func (s *Sphere) LocalNormalAt(p geom.Tuple) geom.Tuple {
	return p.Sub(geom.NewPoint(0, 0, 0))
}

func (s *Sphere) Intersect(r geom.Ray) Intersections {
	return Intersect(r, s.t, s.LocalIntersect)
}

func (s *Sphere) LocalIntersect(r geom.Ray) Intersections {
	sr := r.Origin.Sub(geom.NewPoint(0, 0, 0))
	a := r.Direction.Dot(r.Direction)
	b := 2 * r.Direction.Dot(sr)
	c := sr.Dot(sr) - 1

	d := b*b - 4*a*c

	if d < 0 {
		return NewIntersections()
	}

	return NewIntersections(
		Intersection{
			T: (-b - math.Sqrt(d)) / (2 * a),
			O: s,
		},
		Intersection{
			T: (-b + math.Sqrt(d)) / (2 * a),
			O: s,
		},
	)
}
