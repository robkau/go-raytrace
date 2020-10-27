package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_NewSphere_DefaultTransform(t *testing.T) {
	s := newSphere()

	assert.Equal(t, newIdentityMatrixX4(), s.t)
}

func Test_Sphere_SetTransform(t *testing.T) {
	s := newSphere()
	tr := translate(2, 3, 4)
	s = s.setTransform(tr)

	assert.Equal(t, tr, s.t)
}

func Test_RayIntersectSphere(t *testing.T) {
	s := newSphere()
	r := rayWith(newPoint(0, 0, -5), newVector(0, 0, 1))

	xs := s.intersect(r)

	assert.Len(t, xs.i, 2)
	assert.Equal(t, 4.0, xs.i[0].t)
	assert.Equal(t, 6.0, xs.i[1].t)
}

func Test_RayIntersectSphere_Tangent(t *testing.T) {
	s := newSphere()
	r := rayWith(newPoint(0, 1, -5), newVector(0, 0, 1))

	xs := s.intersect(r)

	assert.Len(t, xs.i, 2)
	assert.Equal(t, 5.0, xs.i[0].t)
	assert.Equal(t, 5.0, xs.i[1].t)
}

func Test_RayIntersectSphere_Miss(t *testing.T) {
	s := newSphere()
	r := rayWith(newPoint(0, 2, -5), newVector(0, 0, 1))

	xs := s.intersect(r)

	assert.Len(t, xs.i, 0)
}

func Test_RayIntersectSphere_FromInside(t *testing.T) {
	s := newSphere()
	r := rayWith(newPoint(0, 0, 0), newVector(0, 0, 1))

	xs := s.intersect(r)

	assert.Len(t, xs.i, 2)
	assert.Equal(t, -1.0, xs.i[0].t)
	assert.Equal(t, 1.0, xs.i[1].t)
}

func Test_RayIntersectSphere_Behind(t *testing.T) {
	s := newSphere()
	r := rayWith(newPoint(0, 0, 5), newVector(0, 0, 1))

	xs := s.intersect(r)

	assert.Len(t, xs.i, 2)
	assert.Equal(t, -6.0, xs.i[0].t)
	assert.Equal(t, -4.0, xs.i[1].t)
}

func Test_RayIntersectSphere_ObjectSet(t *testing.T) {
	s := newSphere()
	r := rayWith(newPoint(0, 0, -5), newVector(0, 0, 1))

	xs := s.intersect(r)

	assert.Len(t, xs.i, 2)
	assert.Equal(t, s, xs.i[0].o)
	assert.Equal(t, s, xs.i[1].o)
}

func Test_ScaledSphere_Intersect_Ray(t *testing.T) {
	s := newSphereWith(scale(2, 2, 2))
	r := rayWith(newPoint(0, 0, -5), newVector(0, 0, 1))

	xs := s.intersect(r)

	assert.Len(t, xs.i, 2)
	assert.Equal(t, 3.0, xs.i[0].t)
	assert.Equal(t, 7.0, xs.i[1].t)
}

func Test_TranslatedSphere_Intersect_Ray(t *testing.T) {
	s := newSphereWith(translate(5, 0, 0))
	r := rayWith(newPoint(0, 0, -5), newVector(0, 0, 1))

	xs := s.intersect(r)

	assert.Len(t, xs.i, 0)
}
