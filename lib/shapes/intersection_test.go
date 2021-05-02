package shapes

import (
	"github.com/stretchr/testify/assert"
	"go-raytrace/lib/geom"
	"math"
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
	c := i.Compute(r, NoIntersections)

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
	c := i.Compute(r, NoIntersections)

	assert.False(t, c.inside)
}

func Test_Compute_Inside(t *testing.T) {
	r := geom.RayWith(geom.NewPoint(0, 0, 0), geom.NewVector(0, 0, 1))
	s := NewSphere()

	i := NewIntersection(1, s)
	c := i.Compute(r, NoIntersections)

	assert.Equal(t, geom.NewPoint(0, 0, 1), c.point)
	assert.Equal(t, geom.NewVector(0, 0, -1), c.Eyev)
	assert.True(t, c.inside)
	assert.Equal(t, geom.NewVector(0, 0, -1), c.Normalv)
}

func Test_Compute_UnderPoint(t *testing.T) {
	r := geom.RayWith(geom.NewPoint(0, 0, -5), geom.NewVector(0, 0, 1))
	s := NewGlassSphere()
	ss := s.SetTransform(geom.Translate(0, 0, 1))

	i := NewIntersection(5, ss)
	c := i.Compute(r, NewIntersections(i))

	assert.Greater(t, c.UnderPoint.Z, geom.FloatComparisonEpsilon/2)
	assert.Less(t, c.point.Z, c.UnderPoint.Z)
}

func Test_Compute_ReflectionVector(t *testing.T) {
	r := geom.RayWith(geom.NewPoint(0, 1, -1), geom.NewVector(0, -math.Sqrt(2)/2, math.Sqrt(2)/2))
	s := NewPlane()

	i := NewIntersection(math.Sqrt(2), s)
	c := i.Compute(r, NoIntersections)

	assert.Equal(t, i.T, c.t)
	assert.Equal(t, i.o, c.Object)
	assert.Equal(t, geom.NewVector(0, math.Sqrt(2)/2, math.Sqrt(2)/2), c.Reflectv)
}

func Test_Compute_RefractionScene(t *testing.T) {
	sA := NewGlassSphere()
	sA.SetTransform(geom.Scale(2, 2, 2))
	m := sA.GetMaterial()
	m.RefractiveIndex = 1.5
	ssA := sA.SetMaterial(m)

	sB := NewGlassSphere()
	sB.SetTransform(geom.Translate(0, 0, -0.25))
	m = sB.GetMaterial()
	m.RefractiveIndex = 2
	ssB := sB.SetMaterial(m)

	sC := NewGlassSphere()
	sC.SetTransform(geom.Translate(0, 0, 0.25))
	m = sC.GetMaterial()
	m.RefractiveIndex = 2.5
	ssC := sC.SetMaterial(m)

	r := geom.RayWith(geom.NewPoint(0, 0, -4), geom.NewVector(0, 0, 1))
	xs := Intersections{I: []Intersection{
		{2, ssA},
		{2.75, ssB},
		{3.25, ssC},
		{4.75, ssB},
		{5.25, ssC},
		{6, ssA},
	}}

	type expected struct {
		n1 float64
		n2 float64
	}
	e := []expected{
		{1, 1.5},
		{1.5, 2},
		{2, 2.5},
		{2.5, 2.5},
		{2.5, 1.5},
		{1.5, 1},
	}

	for i := 0; i < len(e); i++ {
		c := xs.I[i].Compute(r, xs)
		assert.Equal(t, e[i].n1, c.N1)
		assert.Equal(t, e[i].n2, c.N2)

	}
}
