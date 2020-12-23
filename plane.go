package main

import "math"

type plane struct {
	t x4Matrix
	m material
}

func newPlane() plane {
	return plane{
		t: newIdentityMatrixX4(),
		m: newMaterial(),
	}
}

func newPlaneWith(t x4Matrix) plane {
	p := newPlane()
	p.t = t
	return p
}

func (p plane) normalAt(pt tuple) tuple {
	return normalAt(pt, p.t, p.localNormalAt)
}

func (p plane) localNormalAt(pt tuple) tuple {
	return newVector(0, 1, 0)
}

func (p plane) intersect(r ray) intersections {
	return intersect(r, p.t, p.localIntersect)
}

func (p plane) localIntersect(r ray) intersections {
	if math.Abs(r.direction.y) < floatComparisonEpsilon {
		return newIntersections()
	}
	t := -r.origin.y / r.direction.y
	return newIntersections(newIntersection(t, p))
}

func (p plane) getTransform() x4Matrix {
	return p.t
}

func (p plane) setTransform(m x4Matrix) shape {
	p.t = m
	return p
}

func (p plane) getMaterial() material {
	return p.m
}

func (p plane) setMaterial(m material) shape {
	p.m = m
	return p
}
