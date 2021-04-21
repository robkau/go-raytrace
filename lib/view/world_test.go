package view

import (
	"github.com/stretchr/testify/assert"
	"go-raytrace/lib/colors"
	"go-raytrace/lib/geom"
	"go-raytrace/lib/shapes"
	"testing"
)

func Test_NewWorld(t *testing.T) {
	w := NewWorld()

	assert.Len(t, w.objects, 0)
	assert.Len(t, w.lightSources, 0)
}

func Test_DefaultWorld(t *testing.T) {
	w := defaultWorld()
	l := shapes.NewPointLight(geom.NewPoint(-10, 10, -10), colors.White())
	sA := shapes.NewSphere()
	sA.M.Color = colors.NewColor(0.8, 1.0, 0.6)
	sA.M.Diffuse = 0.7
	sA.M.Specular = 0.2
	var sB shapes.Shape = shapes.NewSphere()
	sB = sB.SetTransform(geom.Scale(0.5, 0.5, 0.5))

	assert.Len(t, w.objects, 2)
	assert.Contains(t, w.objects, sA)
	assert.Contains(t, w.objects, sB)
	assert.Len(t, w.lightSources, 1)
	assert.Contains(t, w.lightSources, l)
}

func Test_World_Ray_Intersect(t *testing.T) {
	w := defaultWorld()
	r := geom.RayWith(geom.NewPoint(0, 0, -5), geom.NewVector(0, 0, 1))

	xs := w.Intersect(r)

	assert.Len(t, xs.I, 4)
	assert.Equal(t, 4.0, xs.I[0].T)
	assert.Equal(t, 4.5, xs.I[1].T)
	assert.Equal(t, 5.5, xs.I[2].T)
	assert.Equal(t, 6.0, xs.I[3].T)
}

func Test_Shading_Intersection(t *testing.T) {
	w := defaultWorld()
	r := geom.RayWith(geom.NewPoint(0, 0, -5), geom.NewVector(0, 0, 1))
	s := w.objects[0]

	i := shapes.NewIntersection(4, s)
	c := i.Compute(r)
	cs := w.ShadeHit(c).RoundTo(5)

	assert.Equal(t, colors.NewColor(0.38066, 0.47583, 0.2855), cs)
}

func Test_Shading_Intersection_Inside(t *testing.T) {
	w := defaultWorld()
	w.lightSources[0] = shapes.NewPointLight(geom.NewPoint(0, 0.25, 0), colors.White())
	r := geom.RayWith(geom.NewPoint(0, 0, 0), geom.NewVector(0, 0, 1))
	s := w.objects[1]

	i := shapes.NewIntersection(0.5, s)
	c := i.Compute(r)
	cs := w.ShadeHit(c).RoundTo(5)

	assert.Equal(t, colors.NewColor(0.90498, 0.90498, 0.90498), cs)
}

func Test_RayMissed_Color(t *testing.T) {
	w := defaultWorld()
	r := geom.RayWith(geom.NewPoint(0, 0, -5), geom.NewVector(0, 1, 0))

	c := w.ColorAt(r)

	assert.Equal(t, colors.Black(), c)
}

func Test_RayHit_Color(t *testing.T) {
	w := defaultWorld()
	r := geom.RayWith(geom.NewPoint(0, 0, -5), geom.NewVector(0, 0, 1))

	c := w.ColorAt(r).RoundTo(5)

	assert.Equal(t, colors.NewColor(0.38066, 0.47583, 0.2855), c)
}

func Test_RayIntersectionBehind_Color(t *testing.T) {
	w := defaultWorld()
	m := w.objects[0].GetMaterial()
	m.Ambient = 1
	w.objects[0] = w.objects[0].SetMaterial(m)
	m = w.objects[1].GetMaterial()
	m.Ambient = 1
	w.objects[1] = w.objects[1].SetMaterial(m)

	r := geom.RayWith(geom.NewPoint(0, 0, 0.75), geom.NewVector(0, 0, -1))

	c := w.ColorAt(r)

	assert.Equal(t, w.objects[1].GetMaterial().Color, c)
}

func Test_NoShadow(t *testing.T) {
	w := defaultWorld()
	p := geom.NewPoint(0, 10, 0)

	assert.False(t, w.IsShadowed(p))
}

func Test_NoShadow_LightInBetween(t *testing.T) {
	w := defaultWorld()
	p := geom.NewPoint(10, -10, 10)

	assert.True(t, w.IsShadowed(p))
}

func Test_NoShadow_ObjectBehindLight(t *testing.T) {
	w := defaultWorld()
	p := geom.NewPoint(-20, 20, -20)

	assert.False(t, w.IsShadowed(p))
}

func Test_NoShadow_ObjectBehindPoint(t *testing.T) {
	w := defaultWorld()
	p := geom.NewPoint(-2, 2, -2)

	assert.False(t, w.IsShadowed(p))
}

func Test_ShadeHit_HasShadow(t *testing.T) {
	w := NewWorld()
	w.AddLight(shapes.NewPointLight(geom.NewPoint(0, 0, -10), colors.White()))
	s1 := shapes.NewSphere()
	w.AddObject(s1)
	s2 := shapes.NewSphereWith(geom.Translate(0, 0, 10))
	w.AddObject(s2)
	r := geom.RayWith(geom.NewPoint(0, 0, 5), geom.NewVector(0, 0, 1))
	i := shapes.NewIntersection(4, s2)
	comps := i.Compute(r)
	c := w.ShadeHit(comps)

	assert.Equal(t, colors.NewColor(0.1, 0.1, 0.1), c)
}
