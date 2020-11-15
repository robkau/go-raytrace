package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_IntersectionHoldsSphere(t *testing.T) {
	s := newSphere()
	i := newIntersection(3.5, s)

	assert.Equal(t, s, i.o)
	assert.Equal(t, 3.5, i.t)

}

func Test_NewIntersections(t *testing.T) {
	i := newIntersection(1.0, newSphere())
	ii := newIntersection(2.5, newSphere())

	is := newIntersections(i, ii)

	assert.Len(t, is.i, 2)
	assert.Equal(t, i, is.i[0])
	assert.Equal(t, ii, is.i[1])
}

func Test_NewIntersections_Sorted(t *testing.T) {
	i := newIntersection(1.0, newSphere())
	ii := newIntersection(-2.5, newSphere())
	iii := newIntersection(0, newSphere())

	is := newIntersections(i, ii, iii)

	assert.Len(t, is.i, 3)
	assert.Equal(t, ii, is.i[0])
	assert.Equal(t, iii, is.i[1])
	assert.Equal(t, i, is.i[2])
}

func Test_Hit_AllPositive(t *testing.T) {
	s := newSphere()
	i1 := newIntersection(1, s)
	i2 := newIntersection(2, s)
	xs := newIntersections(i1, i2)

	i, ok := xs.hit()

	assert.True(t, ok)
	assert.Equal(t, i1, i)
}

func Test_Hit_SomePositive(t *testing.T) {
	s := newSphere()
	i1 := newIntersection(-1, s)
	i2 := newIntersection(1, s)
	xs := newIntersections(i1, i2)

	i, ok := xs.hit()

	assert.True(t, ok)
	assert.Equal(t, i2, i)
}

func Test_Hit_SomePositive_NotOrdered(t *testing.T) {
	s := newSphere()
	i1 := newIntersection(5, s)
	i2 := newIntersection(7, s)
	i3 := newIntersection(-3, s)
	i4 := newIntersection(2, s)
	xs := newIntersections(i1, i2, i3, i4)

	i, ok := xs.hit()

	assert.True(t, ok)
	assert.Equal(t, i4, i)
}

func Test_Hit_NonePositive(t *testing.T) {
	s := newSphere()
	i1 := newIntersection(-2, s)
	i2 := newIntersection(-1, s)
	xs := newIntersections(i1, i2)

	_, ok := xs.hit()

	assert.False(t, ok)
}

func Test_Compute(t *testing.T) {
	r := rayWith(newPoint(0, 0, -5), newVector(0, 0, 1))
	s := newSphere()

	i := newIntersection(4, s)
	c := i.compute(r)

	assert.Equal(t, i.t, c.t)
	assert.Equal(t, i.o, c.object)
	assert.Equal(t, newPoint(0, 0, -1), c.point)
	assert.Equal(t, newVector(0, 0, -1), c.eyev)
	assert.Equal(t, newVector(0, 0, -1), c.normalv)

}

func Test_Compute_Outside(t *testing.T) {
	r := rayWith(newPoint(0, 0, -5), newVector(0, 0, 1))
	s := newSphere()

	i := newIntersection(4, s)
	c := i.compute(r)

	assert.False(t, c.inside)
}

func Test_Compute_Inside(t *testing.T) {
	r := rayWith(newPoint(0, 0, 0), newVector(0, 0, 1))
	s := newSphere()

	i := newIntersection(1, s)
	c := i.compute(r)

	assert.Equal(t, newPoint(0, 0, 1), c.point)
	assert.Equal(t, newVector(0, 0, -1), c.eyev)
	assert.True(t, c.inside)
	assert.Equal(t, newVector(0, 0, -1), c.normalv)
}
