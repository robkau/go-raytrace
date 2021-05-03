package shapes

import (
	"go-raytrace/lib/geom"
	"math"
)

type Plane struct {
	t          geom.X4Matrix
	m          Material
	id         string
	shadowless bool
}

func NewPlane() Plane {
	return Plane{
		t:  geom.NewIdentityMatrixX4(),
		m:  NewMaterial(),
		id: newId(),
	}
}

func NewPlaneWith(t geom.X4Matrix) Plane {
	p := NewPlane()
	p.t = t
	return p
}

func (p Plane) Id() string {
	return p.id
}

func (p Plane) NormalAt(pt geom.Tuple) geom.Tuple {
	return NormalAt(pt, p.t, p.LocalNormalAt)
}

func (p Plane) LocalNormalAt(pt geom.Tuple) geom.Tuple {
	return geom.NewVector(0, 1, 0)
}

func (p Plane) Intersect(r geom.Ray) Intersections {
	return Intersect(r, p.t, p.LocalIntersect)
}

func (p Plane) LocalIntersect(r geom.Ray) Intersections {
	if math.Abs(r.Direction.Y) < geom.FloatComparisonEpsilon {
		return NewIntersections()
	}
	t := -r.Origin.Y / r.Direction.Y
	return NewIntersections(NewIntersection(t, p))
}

func (p Plane) GetTransform() geom.X4Matrix {
	return p.t
}

func (p Plane) SetTransform(m geom.X4Matrix) Shape {
	p.t = m
	return p
}

func (p Plane) GetMaterial() Material {
	return p.m
}

func (p Plane) SetMaterial(m Material) Shape {
	p.m = m
	return p
}

func (p Plane) GetShadowless() bool {
	return p.shadowless
}

func (p Plane) SetShadowless(s bool) Shape {
	p.shadowless = s
	return p
}
