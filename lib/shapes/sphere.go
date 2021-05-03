package shapes

import (
	"go-raytrace/lib/geom"
	"math"
)

type Sphere struct {
	t          geom.X4Matrix
	M          Material
	id         string
	shadowless bool
}

func NewSphere() Sphere {
	return Sphere{
		t:  geom.NewIdentityMatrixX4(),
		M:  NewMaterial(),
		id: newId(),
	}
}

func NewSphereWith(t geom.X4Matrix) Sphere {
	s := NewSphere()
	s.t = t
	return s
}

func NewGlassSphere() Sphere {
	s := NewSphere()

	m := NewMaterial()
	m.Transparency = 1
	m.RefractiveIndex = 1.5
	s.M = m
	return s
}

func (s Sphere) Id() string {
	return s.id
}

func (s Sphere) NormalAt(p geom.Tuple) geom.Tuple {
	return NormalAt(p, s.t, s.LocalNormalAt)
}

func (s Sphere) LocalNormalAt(p geom.Tuple) geom.Tuple {
	return p.Sub(geom.NewPoint(0, 0, 0))
}

func (s Sphere) Intersect(r geom.Ray) Intersections {
	return Intersect(r, s.t, s.LocalIntersect)
}

func (s Sphere) LocalIntersect(r geom.Ray) Intersections {
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

func (s Sphere) GetTransform() geom.X4Matrix {
	return s.t
}

func (s Sphere) SetTransform(m geom.X4Matrix) Shape {
	s.t = m
	return s
}

func (s Sphere) GetMaterial() Material {
	return s.M
}

func (s Sphere) SetMaterial(m Material) Shape {
	s.M = m
	return s
}

func (s Sphere) GetShadowless() bool {
	return s.shadowless
}

func (s Sphere) SetShadowless(ss bool) Shape {
	s.shadowless = ss
	return s
}
