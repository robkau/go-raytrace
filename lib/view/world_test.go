package view

import (
	"github.com/stretchr/testify/assert"
	"go-raytrace/lib/colors"
	"go-raytrace/lib/geom"
	"go-raytrace/lib/patterns"
	"go-raytrace/lib/shapes"
	"math"
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
	assert.Equal(t, w.objects[0].GetMaterial(), sA.GetMaterial())
	assert.Equal(t, w.objects[1].GetMaterial(), sB.GetMaterial())
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
	c := i.Compute(r, shapes.NoIntersections)
	cs := w.ShadeHit(c, 0).RoundTo(5)

	assert.Equal(t, colors.NewColor(0.38066, 0.47583, 0.2855), cs)
}

func Test_Shading_Intersection_Inside(t *testing.T) {
	w := defaultWorld()
	w.lightSources[0] = shapes.NewPointLight(geom.NewPoint(0, 0.25, 0), colors.White())
	r := geom.RayWith(geom.NewPoint(0, 0, 0), geom.NewVector(0, 0, 1))
	s := w.objects[1]

	i := shapes.NewIntersection(0.5, s)
	c := i.Compute(r, shapes.NoIntersections)
	cs := w.ShadeHit(c, 0).RoundTo(5)

	assert.Equal(t, colors.NewColor(0.90498, 0.90498, 0.90498), cs)
}

func Test_RayMissed_Color(t *testing.T) {
	w := defaultWorld()
	r := geom.RayWith(geom.NewPoint(0, 0, -5), geom.NewVector(0, 1, 0))

	c := w.ColorAt(r, 0)

	assert.Equal(t, colors.Black(), c)
}

func Test_RayHit_Color(t *testing.T) {
	w := defaultWorld()
	r := geom.RayWith(geom.NewPoint(0, 0, -5), geom.NewVector(0, 0, 1))

	c := w.ColorAt(r, 0).RoundTo(5)

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

	c := w.ColorAt(r, 0)

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
	comps := i.Compute(r, shapes.NoIntersections)
	c := w.ShadeHit(comps, 0)

	assert.Equal(t, colors.NewColor(0.1, 0.1, 0.1), c)
}

func Test_ShadeHit_TransparentMaterial(t *testing.T) {
	w := defaultWorld()

	floor := shapes.NewPlaneWith(geom.Translate(0, -1, 0))
	m := floor.GetMaterial()
	m.Transparency = 0.5
	m.RefractiveIndex = 1.5
	floorShape := floor.SetMaterial(m)
	w.AddObject(floorShape)

	ball := shapes.NewSphereWith(geom.Translate(0, -3.5, -0.5))
	m = floor.GetMaterial()
	m.Ambient = 0.5
	m.Color = colors.NewColor(1, 0, 0)
	ballShape := ball.SetMaterial(m)
	w.AddObject(ballShape)

	r := geom.RayWith(geom.NewPoint(0, 0, -3), geom.NewVector(0, -math.Sqrt(2)/2, math.Sqrt(2)/2))
	xs := shapes.NewIntersections(
		shapes.NewIntersection(math.Sqrt(2), floorShape),
	)

	comps := xs.I[0].Compute(r, xs)
	c := w.ShadeHit(comps, 5)

	assert.Equal(t, colors.NewColor(0.93643, 0.68643, 0.68643), c.RoundTo(5))
}

func Test_ShadeHit_TransparentAndReflectiveMaterial(t *testing.T) {
	w := defaultWorld()

	floor := shapes.NewPlaneWith(geom.Translate(0, -1, 0))
	m := floor.GetMaterial()
	m.Reflective = 0.5
	m.Transparency = 0.5
	m.RefractiveIndex = 1.5
	floorShape := floor.SetMaterial(m)
	w.AddObject(floorShape)

	ball := shapes.NewSphereWith(geom.Translate(0, -3.5, -0.5))
	m = floor.GetMaterial()
	m.Ambient = 0.5
	m.Color = colors.NewColor(1, 0, 0)
	ballShape := ball.SetMaterial(m)
	w.AddObject(ballShape)

	r := geom.RayWith(geom.NewPoint(0, 0, -3), geom.NewVector(0, -math.Sqrt(2)/2, math.Sqrt(2)/2))
	xs := shapes.NewIntersections(
		shapes.NewIntersection(math.Sqrt(2), floorShape),
	)

	comps := xs.I[0].Compute(r, xs)
	c := w.ShadeHit(comps, 5)

	assert.Equal(t, colors.NewColor(0.93392, 0.69643, 0.69243), c.RoundTo(5))
}

func Test_NonReflectiveMaterial_ReflectedColor(t *testing.T) {
	w := defaultWorld()
	r := geom.RayWith(geom.NewPoint(0, 0, 0), geom.NewVector(0, 0, 1))
	s := w.objects[1]
	m := s.GetMaterial()
	m.Ambient = 1
	s.SetMaterial(m)
	i := shapes.NewIntersection(1, s)
	comps := i.Compute(r, shapes.NoIntersections)
	c := w.ReflectedColor(comps, 4)

	// no reflection
	assert.Equal(t, colors.NewColor(0.0, 0.0, 0.0), c)
}

func Test_ReflectiveMaterial_ReflectedColor(t *testing.T) {
	w := defaultWorld()
	r := geom.RayWith(geom.NewPoint(0, 0, -3), geom.NewVector(0, -math.Sqrt(2)/2, math.Sqrt(2)/2))
	s := shapes.NewPlaneWith(geom.Translate(0, -1, 0))
	m := s.GetMaterial()
	m.Reflective = 0.5
	ss := s.SetMaterial(m)
	w.AddObject(ss)
	i := shapes.NewIntersection(math.Sqrt(2), ss)
	comps := i.Compute(r, shapes.NoIntersections)
	c := w.ReflectedColor(comps, 4)

	// reflected
	assert.Equal(t, colors.NewColor(0.19033, 0.23791, 0.14275), c.RoundTo(5))
}

func Test_ReflectiveMaterial_ReflectedColor_ShadeHit(t *testing.T) {
	w := defaultWorld()
	r := geom.RayWith(geom.NewPoint(0, 0, -3), geom.NewVector(0, -math.Sqrt(2)/2, math.Sqrt(2)/2))
	s := shapes.NewPlaneWith(geom.Translate(0, -1, 0))
	m := s.GetMaterial()
	m.Reflective = 0.5
	ss := s.SetMaterial(m)
	w.AddObject(ss)
	i := shapes.NewIntersection(math.Sqrt(2), ss)
	comps := i.Compute(r, shapes.NoIntersections)
	c := w.ShadeHit(comps, 4)

	// reflected
	assert.Equal(t, colors.NewColor(0.87676, 0.92434, 0.82917), c.RoundTo(5))
}

func Test_ReflectiveMaterial_ReflectedColor_ShadeHit_NoReflectionsRemaining(t *testing.T) {
	w := defaultWorld()
	r := geom.RayWith(geom.NewPoint(0, 0, -3), geom.NewVector(0, -math.Sqrt(2)/2, math.Sqrt(2)/2))
	s := shapes.NewPlaneWith(geom.Translate(0, -1, 0))
	m := s.GetMaterial()
	m.Reflective = 0.5
	ss := s.SetMaterial(m)
	w.AddObject(ss)
	i := shapes.NewIntersection(math.Sqrt(2), ss)
	comps := i.Compute(r, shapes.NoIntersections)
	c := w.ReflectedColor(comps, 0)

	// not reflected
	assert.Equal(t, colors.NewColor(0, 0, 0), c)
}

func Test_MutuallyReflecting_InfiniteRecursion(t *testing.T) {
	w := NewWorld()
	w.AddLight(shapes.NewPointLight(geom.NewPoint(0, 0, 0), colors.NewColor(1, 1, 1)))

	lower := shapes.NewPlaneWith(geom.Translate(0, -1, 0))
	m := lower.GetMaterial()
	m.Reflective = 1
	lowerS := lower.SetMaterial(m)
	w.AddObject(lowerS)

	upper := shapes.NewPlaneWith(geom.Translate(0, 1, 0))
	m = upper.GetMaterial()
	m.Reflective = 1
	upperS := upper.SetMaterial(m)
	w.AddObject(upperS)

	r := geom.RayWith(geom.NewPoint(0, 0, 0), geom.NewVector(0, 1, 0))

	// just testing that the call finishes
	assert.NotPanics(t, func() { w.ColorAt(r, 5) })
}

func Test_RefractedColor_OpaqueSurface(t *testing.T) {
	w := defaultWorld()
	s := w.objects[0]
	r := geom.RayWith(geom.NewPoint(0, 0, -5), geom.NewVector(0, 0, 1))

	xs := shapes.NewIntersections(
		shapes.NewIntersection(4, s),
		shapes.NewIntersection(6, s),
	)

	comps := xs.I[0].Compute(r, xs)
	c := w.RefractedColor(comps, 5)

	assert.Equal(t, colors.NewColor(0, 0, 0), c)
}

func Test_RefractedColor_AtMaxRecursionDepth(t *testing.T) {
	w := defaultWorld()
	s := w.objects[0]
	m := s.GetMaterial()
	m.Transparency = 1
	m.RefractiveIndex = 1.5
	ss := s.SetMaterial(m)
	w.objects[0] = ss

	r := geom.RayWith(geom.NewPoint(0, 0, -5), geom.NewVector(0, 0, 1))

	xs := shapes.NewIntersections(
		shapes.NewIntersection(4, ss),
		shapes.NewIntersection(6, ss),
	)

	comps := xs.I[0].Compute(r, xs)
	c := w.RefractedColor(comps, 0)

	assert.Equal(t, colors.NewColor(0, 0, 0), c)
}

func Test_RefractedColor_TotalInternalRefraction(t *testing.T) {
	w := defaultWorld()
	s := w.objects[0]
	m := s.GetMaterial()
	m.Transparency = 1
	m.RefractiveIndex = 1.5
	ss := s.SetMaterial(m)
	w.objects[0] = ss

	r := geom.RayWith(geom.NewPoint(0, 0, math.Sqrt(2)/2), geom.NewVector(0, 1, 0))

	xs := shapes.NewIntersections(
		shapes.NewIntersection(-math.Sqrt(2)/2, ss),
		shapes.NewIntersection(math.Sqrt(2)/2, ss),
	)

	comps := xs.I[1].Compute(r, xs)
	c := w.RefractedColor(comps, 5)

	assert.Equal(t, colors.NewColor(0, 0, 0), c)
}

func Test_RefractedColor_Ray(t *testing.T) {
	w := defaultWorld()
	sA := w.objects[0]
	m := sA.GetMaterial()
	m.Ambient = 1
	m.Pattern = patterns.NewPositionAsColorPattern()
	ssA := sA.SetMaterial(m)
	w.objects[0] = ssA
	sB := w.objects[1]
	m = sB.GetMaterial()
	m.Transparency = 1
	m.RefractiveIndex = 1.5
	ssB := sB.SetMaterial(m)
	w.objects[1] = ssB

	r := geom.RayWith(geom.NewPoint(0, 0, 0.1), geom.NewVector(0, 1, 0))

	xs := shapes.NewIntersections(
		shapes.NewIntersection(-0.9899, ssA),
		shapes.NewIntersection(-0.4899, ssB),
		shapes.NewIntersection(0.4899, ssB),
		shapes.NewIntersection(0.9899, ssA),
	)

	comps := xs.I[2].Compute(r, xs)
	c := w.RefractedColor(comps, 5)

	assert.Equal(t, colors.NewColor(0, 0.99888, 0.04722), c.RoundTo(5))
}
