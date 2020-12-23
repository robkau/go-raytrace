package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_PlaneNormal_Constant(t *testing.T) {
	p := newPlane()
	n1 := p.localNormalAt(newPoint(0, 0, 0))
	n2 := p.localNormalAt(newPoint(10, 0, -10))
	n3 := p.localNormalAt(newPoint(-5, 10, 150))

	assert.Equal(t, newVector(0, 1, 0), n1)
	assert.Equal(t, newVector(0, 1, 0), n2)
	assert.Equal(t, newVector(0, 1, 0), n3)
}

func Test_Plane_Intersect_Parallel_Ray(t *testing.T) {
	p := newPlane()
	r := rayWith(newPoint(0, 10, 0), newVector(0, 0, 1))

	xs := p.localIntersect(r)

	assert.Len(t, xs.i, 0)
}

func Test_Plane_Intersect_Coplanar_Ray(t *testing.T) {
	p := newPlane()
	r := rayWith(newPoint(0, 0, 0), newVector(0, 0, 1))

	xs := p.localIntersect(r)

	assert.Len(t, xs.i, 0)
}

func Test_Plane_Intersect_Ray_Above(t *testing.T) {
	p := newPlane()
	r := rayWith(newPoint(0, 1, 0), newVector(0, -1, 0))

	xs := p.localIntersect(r)

	assert.Len(t, xs.i, 1)
	assert.Equal(t, 1.0, xs.i[0].t)
	assert.Equal(t, p, xs.i[0].o)
}

func Test_Plane_Intersect_Ray_Below(t *testing.T) {
	p := newPlane()
	r := rayWith(newPoint(0, -1, 0), newVector(0, 1, 0))

	xs := p.localIntersect(r)

	assert.Len(t, xs.i, 1)
	assert.Equal(t, 1.0, xs.i[0].t)
	assert.Equal(t, p, xs.i[0].o)
}
