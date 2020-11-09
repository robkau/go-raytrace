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
	return sphere{
		t: t,
	}
}

func (s sphere) setTransform(t x4Matrix) sphere {
	s.t = t
	return s
}

func (s sphere) normalAt(p tuple) tuple {
	objectPoint := s.t.invert().mulTuple(p)
	objectNormal := objectPoint.sub(newPoint(0, 0, 0))
	worldNormal := s.t.invert().transpose().mulTuple(objectNormal)
	worldNormal.c = vector
	return worldNormal.normalize()
}

func (s sphere) intersect(r ray) intersections {
	// apply inverse sphere transformation onto ray
	r = r.transform(s.t.invert())

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
