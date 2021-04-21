package shapes

import (
	"github.com/stretchr/testify/assert"
	"go-raytrace/lib/geom"
	"testing"
)

func Test_IntersectionHoldsSphere(t *testing.T) {
	s := NewSphere()
	i := NewIntersection(3.5, s)

	assert.Equal(t, s, i.o)
	assert.Equal(t, 3.5, i.T)

}

func Test_NewIntersections(t *testing.T) {
	i := NewIntersection(1.0, NewSphere())
	ii := NewIntersection(2.5, NewSphere())

	is := NewIntersections(i, ii)

	assert.Len(t, is.I, 2)
	assert.Equal(t, i, is.I[0])
	assert.Equal(t, ii, is.I[1])
}

func Test_NewIntersections_Sorted(t *testing.T) {
	i := NewIntersection(1.0, NewSphere())
	ii := NewIntersection(-2.5, NewSphere())
	iii := NewIntersection(0, NewSphere())

	is := NewIntersections(i, ii, iii)

	assert.Len(t, is.I, 3)
	assert.Equal(t, ii, is.I[0])
	assert.Equal(t, iii, is.I[1])
	assert.Equal(t, i, is.I[2])
}

func Test_Hit_AllPositive(t *testing.T) {
	s := NewSphere()
	i1 := NewIntersection(1, s)
	i2 := NewIntersection(2, s)
	xs := NewIntersections(i1, i2)

	i, ok := xs.Hit()

	assert.True(t, ok)
	assert.Equal(t, i1, i)
}

func Test_Hit_SomePositive(t *testing.T) {
	s := NewSphere()
	i1 := NewIntersection(-1, s)
	i2 := NewIntersection(1, s)
	xs := NewIntersections(i1, i2)

	i, ok := xs.Hit()

	assert.True(t, ok)
	assert.Equal(t, i2, i)
}

func Test_Hit_SomePositive_NotOrdered(t *testing.T) {
	s := NewSphere()
	i1 := NewIntersection(5, s)
	i2 := NewIntersection(7, s)
	i3 := NewIntersection(-3, s)
	i4 := NewIntersection(2, s)
	xs := NewIntersections(i1, i2, i3, i4)

	i, ok := xs.Hit()

	assert.True(t, ok)
	assert.Equal(t, i4, i)
}

func Test_Hit_NonePositive(t *testing.T) {
	s := NewSphere()
	i1 := NewIntersection(-2, s)
	i2 := NewIntersection(-1, s)
	xs := NewIntersections(i1, i2)

	_, ok := xs.Hit()

	assert.False(t, ok)
}

func Test_Compute(t *testing.T) {
	r := geom.RayWith(geom.NewPoint(0, 0, -5), geom.NewVector(0, 0, 1))
	s := NewSphere()

	i := NewIntersection(4, s)
	c := i.Compute(r)

	assert.Equal(t, i.T, c.t)
	assert.Equal(t, i.o, c.Object)
	assert.Equal(t, geom.NewPoint(0, 0, -1), c.point)
	assert.Equal(t, geom.NewVector(0, 0, -1), c.Eyev)
	assert.Equal(t, geom.NewVector(0, 0, -1), c.Normalv)

}

func Test_Compute_Outside(t *testing.T) {
	r := geom.RayWith(geom.NewPoint(0, 0, -5), geom.NewVector(0, 0, 1))
	s := NewSphere()

	i := NewIntersection(4, s)
	c := i.Compute(r)

	assert.False(t, c.inside)
}

func Test_Compute_Inside(t *testing.T) {
	r := geom.RayWith(geom.NewPoint(0, 0, 0), geom.NewVector(0, 0, 1))
	s := NewSphere()

	i := NewIntersection(1, s)
	c := i.Compute(r)

	assert.Equal(t, geom.NewPoint(0, 0, 1), c.point)
	assert.Equal(t, geom.NewVector(0, 0, -1), c.Eyev)
	assert.True(t, c.inside)
	assert.Equal(t, geom.NewVector(0, 0, -1), c.Normalv)
}
