package main

import "math"

type sphere struct {
	t x4Matrix
	m material
}

func newSphere() sphere {
	return sphere{
		t: newIdentityMatrixX4(),
		m: newMaterial(),
	}
}

func newSphereWith(t x4Matrix) sphere {
	s := newSphere()
	s.t = t
	return s
}

func (s sphere) normalAt(p tuple) tuple {
	return normalAt(p, s.t, s.localNormalAt)
}

func (s sphere) localNormalAt(p tuple) tuple {
	return p.sub(newPoint(0, 0, 0))
}

func (s sphere) intersect(r ray) intersections {
	return intersect(r, s.t, s.localIntersect)
}

func (s sphere) localIntersect(r ray) intersections {
	sr := r.origin.sub(newPoint(0, 0, 0))
	a := r.direction.dot(r.direction)
	b := 2 * r.direction.dot(sr)
	c := sr.dot(sr) - 1

	d := b*b - 4*a*c

	if d < 0 {
		return newIntersections()
	}

	return newIntersections(
		intersection{
			t: (-b - math.Sqrt(d)) / (2 * a),
			o: s,
		},
		intersection{
			t: (-b + math.Sqrt(d)) / (2 * a),
			o: s,
		},
	)
}

func (s sphere) getTransform() x4Matrix {
	return s.t
}

func (s sphere) setTransform(m x4Matrix) shape {
	s.t = m
	return s
}

func (s sphere) getMaterial() material {
	return s.m
}

func (s sphere) setMaterial(m material) shape {
	s.m = m
	return s
}
