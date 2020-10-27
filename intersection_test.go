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
