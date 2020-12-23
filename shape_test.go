package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math"
	"testing"
)

func Test_GetDefaultTransformation(t *testing.T) {
	s := newTestShape()
	require.Equal(t, newIdentityMatrixX4(), s.getTransform())
}

func Test_SetTransformation(t *testing.T) {
	var s shape = newTestShape()
	s = s.setTransform(translate(2, 3, 4))
	require.Equal(t, translate(2, 3, 4), s.getTransform())
}

func Test_GetDefaultMaterial(t *testing.T) {
	s := newTestShape()
	require.Equal(t, newMaterial(), s.getMaterial())
}

func Test_SetMaterial(t *testing.T) {
	var s shape = newTestShape()
	m := s.getMaterial()
	m.ambient = 1
	s = s.setMaterial(m)
	require.Equal(t, m, s.getMaterial())
}

func Test_IntersectScaledShape(t *testing.T) {
	r := rayWith(newPoint(0, 0, -5), newVector(0, 0, 1))
	var ts = newTestShape()
	var s shape = ts
	s = ts.setTransform(scale(2, 2, 2))

	_ = s.intersect(r)

	assert.Equal(t, newPoint(0, 0, -2.5), ts.savedRay.origin)
	assert.Equal(t, newVector(0, 0, 0.5), ts.savedRay.direction)
}

func Test_IntersectTranslatedShape(t *testing.T) {
	r := rayWith(newPoint(0, 0, -5), newVector(0, 0, 1))
	var ts = newTestShape()
	var s shape = ts
	s = ts.setTransform(translate(5, 0, 0))

	_ = s.intersect(r)

	assert.Equal(t, newPoint(-5, 0, -5), ts.savedRay.origin)
	assert.Equal(t, newVector(0, 0, 1), ts.savedRay.direction)
}

func Test_NormalTranslatedShape(t *testing.T) {
	var ts = newTestShape()
	var s shape = ts
	s = ts.setTransform(translate(0, 1, 0))

	n := s.normalAt(newPoint(0, 1.70711, -0.70711)).roundTo(5)

	assert.Equal(t, newVector(0, 0.70711, -0.70711), n)
}

func Test_NormalTransformedShape(t *testing.T) {
	var ts = newTestShape()
	var s shape = ts
	s = ts.setTransform(scale(1, 0.5, 1).mulX4Matrix(rotateZ(math.Pi / 5)))

	n := s.normalAt(newPoint(0, math.Sqrt2/2, -math.Sqrt2/2)).roundTo(5)

	assert.Equal(t, newVector(0, 0.97014, -0.24254), n)
}

type testShape struct {
	t        x4Matrix
	m        material
	savedRay ray
}

func newTestShape() *testShape {
	return &testShape{
		t: newIdentityMatrixX4(),
		m: newMaterial(),
	}
}

func (t *testShape) intersect(r ray) intersections {
	lr := r.transform(t.t.invert())
	t.savedRay = lr
	return intersections{}
}

func (t *testShape) normalAt(p tuple) tuple {
	localPoint := t.t.invert().mulTuple(p)
	localNormal := newVector(localPoint.x, localPoint.y, localPoint.z)
	worldNormal := t.t.invert().transpose().mulTuple(localNormal)
	worldNormal.c = vector
	return worldNormal.normalize()
}

func (t *testShape) getTransform() x4Matrix {
	return t.t
}

func (t *testShape) setTransform(m x4Matrix) shape {
	t.t = m
	return t // make interface happy
}

func (t *testShape) getMaterial() material {
	return t.m
}

func (t *testShape) setMaterial(m material) shape {
	t.m = m
	return t // make interface happy
}
