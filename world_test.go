package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_NewWorld(t *testing.T) {
	w := newWorld()

	assert.Len(t, w.objects, 0)
	assert.Len(t, w.lightSources, 0)
}

func Test_DefaultWorld(t *testing.T) {
	w := defaultWorld()
	l := newPointLight(newPoint(-10, 10, -10), color{1, 1, 1})
	sA := newSphere()
	sA.m.color = color{0.8, 1.0, 0.6}
	sA.m.diffuse = 0.7
	sA.m.specular = 0.2
	sB := newSphere()
	sB = sB.setTransform(scale(0.5, 0.5, 0.5))

	assert.Len(t, w.objects, 2)
	assert.Contains(t, w.objects, sA)
	assert.Contains(t, w.objects, sB)
	assert.Len(t, w.lightSources, 1)
	assert.Contains(t, w.lightSources, l)
}

func Test_World_Ray_Intersect(t *testing.T) {
	w := defaultWorld()
	r := rayWith(newPoint(0, 0, -5), newVector(0, 0, 1))

	xs := w.intersect(r)

	assert.Len(t, xs.i, 4)
	assert.Equal(t, 4.0, xs.i[0].t)
	assert.Equal(t, 4.5, xs.i[1].t)
	assert.Equal(t, 5.5, xs.i[2].t)
	assert.Equal(t, 6.0, xs.i[3].t)
}

func Test_Shading_Intersection(t *testing.T) {
	w := defaultWorld()
	r := rayWith(newPoint(0, 0, -5), newVector(0, 0, 1))
	s := w.objects[0]

	i := newIntersection(4, s)
	c := i.compute(r)
	cs := w.shadeHit(c).roundTo(5)

	assert.Equal(t, color{0.38066, 0.47583, 0.2855}, cs)
}

func Test_Shading_Intersection_Inside(t *testing.T) {
	w := defaultWorld()
	w.lightSources[0] = newPointLight(newPoint(0, 0.25, 0), color{1, 1, 1})
	r := rayWith(newPoint(0, 0, 0), newVector(0, 0, 1))
	s := w.objects[1]

	i := newIntersection(0.5, s)
	c := i.compute(r)
	cs := w.shadeHit(c).roundTo(5)

	assert.Equal(t, color{0.90498, 0.90498, 0.90498}, cs)
}

func Test_RayMissed_Color(t *testing.T) {
	w := defaultWorld()
	r := rayWith(newPoint(0, 0, -5), newVector(0, 1, 0))

	c := w.colorAt(r)

	assert.Equal(t, color{0, 0, 0}, c)
}

func Test_RayHit_Color(t *testing.T) {
	w := defaultWorld()
	r := rayWith(newPoint(0, 0, -5), newVector(0, 0, 1))

	c := w.colorAt(r).roundTo(5)

	assert.Equal(t, color{0.38066, 0.47583, 0.2855}, c)
}

func Test_RayIntersectionBehind_Color(t *testing.T) {
	w := defaultWorld()
	w.objects[0].m.ambient = 1
	w.objects[1].m.ambient = 1

	r := rayWith(newPoint(0, 0, 0.75), newVector(0, 0, -1))

	c := w.colorAt(r)

	assert.Equal(t, w.objects[1].m.color, c)
}

func Test_NoShadow(t *testing.T) {
	w := defaultWorld()
	p := newPoint(0, 10, 0)

	assert.False(t, w.isShadowed(p))
}

func Test_NoShadow_LightInBetween(t *testing.T) {
	w := defaultWorld()
	p := newPoint(10, -10, 10)

	assert.True(t, w.isShadowed(p))
}

func Test_NoShadow_ObjectBehindLight(t *testing.T) {
	w := defaultWorld()
	p := newPoint(-20, 20, -20)

	assert.False(t, w.isShadowed(p))
}

func Test_NoShadow_ObjectBehindPoint(t *testing.T) {
	w := defaultWorld()
	p := newPoint(-2, 2, -2)

	assert.False(t, w.isShadowed(p))
}

func Test_ShadeHit_HasShadow(t *testing.T) {
	w := newWorld()
	w.addLight(newPointLight(newPoint(0, 0, -10), color{1, 1, 1}))
	s1 := newSphere()
	w.addObject(s1)
	s2 := newSphereWith(translate(0, 0, 10))
	w.addObject(s2)
	r := rayWith(newPoint(0, 0, 5), newVector(0, 0, 1))
	i := newIntersection(4, s2)
	comps := i.compute(r)
	c := w.shadeHit(comps)

	assert.Equal(t, color{0.1, 0.1, 0.1}, c)
}
