package view

import (
	"github.com/robkau/go-raytrace/lib/colors"
	"github.com/robkau/go-raytrace/lib/geom"
	"github.com/robkau/go-raytrace/lib/materials"
	"github.com/robkau/go-raytrace/lib/patterns"
	"github.com/robkau/go-raytrace/lib/shapes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math"
	"strconv"
	"testing"
)

func Test_NewWorld(t *testing.T) {
	w := NewWorld()

	assert.Len(t, w.objects, 0)
	assert.Len(t, w.pointLights, 0)
}

func Test_DefaultWorld(t *testing.T) {
	w := defaultWorld()
	l := shapes.NewPointLight(geom.NewPoint(-10, 10, -10), colors.White())
	sA := shapes.NewSphere()
	m := sA.GetMaterial()
	m.Color = colors.NewColor(0.8, 1.0, 0.6)
	m.Diffuse = 0.7
	m.Specular = 0.2
	sA.SetMaterial(m)
	sB := shapes.NewSphere()
	sB.SetTransform(geom.Scale(0.5, 0.5, 0.5))

	assert.Len(t, w.objects, 2)
	assert.Equal(t, w.objects[0].GetMaterial(), sA.GetMaterial())
	assert.Equal(t, w.objects[1].GetMaterial(), sB.GetMaterial())
	assert.Len(t, w.pointLights, 1)
	assert.Contains(t, w.pointLights, l)
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
	c := i.Compute(r, shapes.NewIntersections())
	cs := w.ShadeHit(c, 0).RoundTo(5)

	assert.Equal(t, colors.NewColor(0.38066, 0.47583, 0.2855), cs)
}

func Test_Shading_Intersection_Inside(t *testing.T) {
	w := defaultWorld()
	w.pointLights[0] = shapes.NewPointLight(geom.NewPoint(0, 0.25, 0), colors.White())
	r := geom.RayWith(geom.NewPoint(0, 0, 0), geom.NewVector(0, 0, 1))
	s := w.objects[1]

	i := shapes.NewIntersection(0.5, s)
	c := i.Compute(r, shapes.NewIntersections())
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
	w.objects[0].SetMaterial(m)
	m = w.objects[1].GetMaterial()
	m.Ambient = 1
	w.objects[1].SetMaterial(m)

	r := geom.RayWith(geom.NewPoint(0, 0, 0.75), geom.NewVector(0, 0, -1))

	c := w.ColorAt(r, 0)

	assert.Equal(t, w.objects[1].GetMaterial().Color, c)
}

func Test_IsShadowed(t *testing.T) {
	w := defaultWorld()
	lightPosition := geom.NewPoint(-10, -10, -10)

	type args struct {
		p      geom.Tuple
		expect bool
	}

	tests := []args{
		{geom.NewPoint(-10, -10, 10), false},
		{geom.NewPoint(10, 10, 10), true},
		{geom.NewPoint(-20, -20, -20), false},
		{geom.NewPoint(-5, -5, -5), false},
	}

	for ti, tt := range tests {
		t.Run(t.Name()+strconv.Itoa(ti), func(t *testing.T) {
			v := w.IsShadowed(lightPosition, tt.p)
			require.Equal(t, tt.expect, v)
		})
	}
}

func Test_ShadeHit_HasShadow(t *testing.T) {
	w := NewWorld()
	w.AddPointLight(shapes.NewPointLight(geom.NewPoint(0, 0, -10), colors.White()))
	s1 := shapes.NewSphere()
	w.AddObject(s1)
	s2 := shapes.NewSphere()
	s2.SetTransform(geom.Translate(0, 0, 10))
	w.AddObject(s2)
	r := geom.RayWith(geom.NewPoint(0, 0, 5), geom.NewVector(0, 0, 1))
	i := shapes.NewIntersection(4, s2)
	comps := i.Compute(r, shapes.NewIntersections())
	c := w.ShadeHit(comps, 0)

	assert.Equal(t, colors.NewColor(0.1, 0.1, 0.1), c)
}

func Test_ShadeHit_TransparentMaterial(t *testing.T) {
	w := defaultWorld()

	floor := shapes.NewPlane()
	floor.SetTransform(geom.Translate(0, -1, 0))
	m := floor.GetMaterial()
	m.Transparency = 0.5
	m.RefractiveIndex = 1.5
	floor.SetMaterial(m)
	w.AddObject(floor)

	ball := shapes.NewSphere()
	ball.SetTransform(geom.Translate(0, -3.5, -0.5))
	m = ball.GetMaterial()
	m.Ambient = 0.5
	m.Color = colors.NewColor(1, 0, 0)
	ball.SetMaterial(m)
	w.AddObject(ball)

	r := geom.RayWith(geom.NewPoint(0, 0, -3), geom.NewVector(0, -math.Sqrt(2)/2, math.Sqrt(2)/2))
	xs := shapes.NewIntersections(
		shapes.NewIntersection(math.Sqrt(2), floor),
	)

	comps := xs.I[0].Compute(r, xs)
	c := w.ShadeHit(comps, 5)

	assert.Equal(t, colors.NewColor(0.93643, 0.68643, 0.68643), c.RoundTo(5))
}

func Test_ShadeHit_TransparentAndReflectiveMaterial(t *testing.T) {
	w := defaultWorld()

	floor := shapes.NewPlane()
	floor.SetTransform(geom.Translate(0, -1, 0))
	m := floor.GetMaterial()
	m.Reflective = 0.5
	m.Transparency = 0.5
	m.RefractiveIndex = 1.5
	floor.SetMaterial(m)
	w.AddObject(floor)

	ball := shapes.NewSphere()
	ball.SetTransform(geom.Translate(0, -3.5, -0.5))
	m = ball.GetMaterial()
	m.Ambient = 0.5
	m.Color = colors.NewColor(1, 0, 0)
	ball.SetMaterial(m)
	w.AddObject(ball)

	r := geom.RayWith(geom.NewPoint(0, 0, -3), geom.NewVector(0, -math.Sqrt(2)/2, math.Sqrt(2)/2))
	xs := shapes.NewIntersections(
		shapes.NewIntersection(math.Sqrt(2), floor),
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
	comps := i.Compute(r, shapes.NewIntersections())
	c := w.ReflectedColor(comps, 4)

	// no reflection
	assert.Equal(t, colors.NewColor(0.0, 0.0, 0.0), c)
}

func Test_ReflectiveMaterial_ReflectedColor(t *testing.T) {
	w := defaultWorld()
	r := geom.RayWith(geom.NewPoint(0, 0, -3), geom.NewVector(0, -math.Sqrt(2)/2, math.Sqrt(2)/2))
	s := shapes.NewPlane()
	s.SetTransform(geom.Translate(0, -1, 0))
	m := s.GetMaterial()
	m.Reflective = 0.5
	s.SetMaterial(m)
	w.AddObject(s)
	i := shapes.NewIntersection(math.Sqrt(2), s)
	comps := i.Compute(r, shapes.NewIntersections())
	c := w.ReflectedColor(comps, 4)

	// reflected
	assert.Equal(t, colors.NewColor(0.19033, 0.23791, 0.14275), c.RoundTo(5))
}

func Test_ReflectiveMaterial_ReflectedColor_ShadeHit(t *testing.T) {
	w := defaultWorld()
	r := geom.RayWith(geom.NewPoint(0, 0, -3), geom.NewVector(0, -math.Sqrt(2)/2, math.Sqrt(2)/2))
	s := shapes.NewPlane()
	s.SetTransform(geom.Translate(0, -1, 0))
	m := s.GetMaterial()
	m.Reflective = 0.5
	s.SetMaterial(m)
	w.AddObject(s)
	i := shapes.NewIntersection(math.Sqrt(2), s)
	comps := i.Compute(r, shapes.NewIntersections())
	c := w.ShadeHit(comps, 4)

	// reflected
	assert.Equal(t, colors.NewColor(0.87676, 0.92434, 0.82917), c.RoundTo(5))
}

func Test_ReflectiveMaterial_ReflectedColor_ShadeHit_NoReflectionsRemaining(t *testing.T) {
	w := defaultWorld()
	r := geom.RayWith(geom.NewPoint(0, 0, -3), geom.NewVector(0, -math.Sqrt(2)/2, math.Sqrt(2)/2))
	s := shapes.NewPlane()
	s.SetTransform(geom.Translate(0, -1, 0))
	m := s.GetMaterial()
	m.Reflective = 0.5
	s.SetMaterial(m)
	w.AddObject(s)
	i := shapes.NewIntersection(math.Sqrt(2), s)
	comps := i.Compute(r, shapes.NewIntersections())
	c := w.ReflectedColor(comps, 0)

	// not reflected
	assert.Equal(t, colors.NewColor(0, 0, 0), c)
}

func Test_MutuallyReflecting_InfiniteRecursion(t *testing.T) {
	w := NewWorld()
	w.AddPointLight(shapes.NewPointLight(geom.NewPoint(0, 0, 0), colors.NewColor(1, 1, 1)))

	lower := shapes.NewPlane()
	lower.SetTransform(geom.Translate(0, -1, 0))
	m := lower.GetMaterial()
	m.Reflective = 1
	lower.SetMaterial(m)
	w.AddObject(lower)

	upper := shapes.NewPlane()
	upper.SetTransform(geom.Translate(0, 1, 0))
	m = upper.GetMaterial()
	m.Reflective = 1
	upper.SetMaterial(m)
	w.AddObject(upper)

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
	s.SetMaterial(m)
	w.objects[0] = s

	r := geom.RayWith(geom.NewPoint(0, 0, -5), geom.NewVector(0, 0, 1))

	xs := shapes.NewIntersections(
		shapes.NewIntersection(4, s),
		shapes.NewIntersection(6, s),
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
	s.SetMaterial(m)
	w.objects[0] = s

	r := geom.RayWith(geom.NewPoint(0, 0, math.Sqrt(2)/2), geom.NewVector(0, 1, 0))

	xs := shapes.NewIntersections(
		shapes.NewIntersection(-math.Sqrt(2)/2, s),
		shapes.NewIntersection(math.Sqrt(2)/2, s),
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
	sA.SetMaterial(m)
	w.objects[0] = sA
	sB := w.objects[1]
	m = sB.GetMaterial()
	m.Transparency = 1
	m.RefractiveIndex = 1.5
	sB.SetMaterial(m)
	w.objects[1] = sB

	r := geom.RayWith(geom.NewPoint(0, 0, 0.1), geom.NewVector(0, 1, 0))

	xs := shapes.NewIntersections(
		shapes.NewIntersection(-0.9899, sA),
		shapes.NewIntersection(-0.4899, sB),
		shapes.NewIntersection(0.4899, sB),
		shapes.NewIntersection(0.9899, sA),
	)

	comps := xs.I[2].Compute(r, xs)
	c := w.RefractedColor(comps, 5)

	assert.Equal(t, colors.NewColor(0, 0.99888, 0.04722), c.RoundTo(5))
}

func Test_Shadowless_NotShadeOthers(t *testing.T) {
	w := NewWorld()
	w.AddPointLight(shapes.NewPointLight(geom.NewPoint(0, 0, -10), colors.White()))
	var s1 shapes.Shape = shapes.NewSphere()
	s1.SetShadowless(true)
	w.AddObject(s1)
	s2 := shapes.NewSphere()
	s2.SetTransform(geom.Translate(0, 0, 10))
	w.AddObject(s2)
	r := geom.RayWith(geom.NewPoint(0, 0, 5), geom.NewVector(0, 0, 1))
	i := shapes.NewIntersection(4, s2)
	comps := i.Compute(r, shapes.NewIntersections())
	c := w.ShadeHit(comps, 0)

	assert.True(t, colors.NewColor(1.9, 1.9, 1.9).Equal(c))
}

func Test_UnShaded_NotShadedByOthers(t *testing.T) {
	w := NewWorld()
	w.AddPointLight(shapes.NewPointLight(geom.NewPoint(0, 0, -10), colors.White()))
	s1 := shapes.NewSphere()
	w.AddObject(s1)
	s2 := shapes.NewSphere()
	s2.SetTransform(geom.Translate(0, 0, 10))
	s2.SetShaded(false)
	w.AddObject(s2)
	r := geom.RayWith(geom.NewPoint(0, 0, 5), geom.NewVector(0, 0, 1))
	i := shapes.NewIntersection(4, s2)
	comps := i.Compute(r, shapes.NewIntersections())
	c := w.ShadeHit(comps, 0)

	assert.True(t, colors.NewColor(1.9, 1.9, 1.9).Equal(c))
}

func Test_PointLight_PassesIntensity(t *testing.T) {
	w := defaultWorld()
	require.Len(t, w.pointLights, 1)
	l := w.pointLights[0]

	type args struct {
		p      geom.Tuple
		expect float64
	}

	tests := []args{
		{geom.NewPoint(0, 1.0001, 0), 1.0},
		{geom.NewPoint(-1.0001, 0, 0), 1.0},
		{geom.NewPoint(0, 0, -1.0001), 1.0},
		{geom.NewPoint(0, 0, 1.0001), 0.0},
		{geom.NewPoint(1.0001, 0, 0), 0.0},
		{geom.NewPoint(0, -1.0001, 0), 0.0},
		{geom.NewPoint(0, 0, 0), 0.0},
	}

	for ti, tt := range tests {
		t.Run(t.Name()+strconv.Itoa(ti), func(t *testing.T) {
			intensity := IntensityAt(l, tt.p, w)
			require.Equal(t, tt.expect, intensity)
		})
	}
}

func Test_AreaLight_PassesIntensity(t *testing.T) {
	w := defaultWorld()

	corner := geom.NewPoint(-0.5, -0.5, -5)
	v1 := geom.NewVector(1, 0, 0)
	v2 := geom.NewVector(0, 1, 0)
	light := shapes.NewAreaLight(corner, v1, 2, v2, 2, colors.White(), shapes.NewJitterSequence(0.5))

	type args struct {
		p      geom.Tuple
		expect float64
	}

	tests := []args{
		{geom.NewPoint(0, 0, 2), 0.0},
		{geom.NewPoint(1, -1, 2), 0.25},
		{geom.NewPoint(1.5, 0, 2), 0.5},
		{geom.NewPoint(1.25, 1.25, 3), 0.75},
		{geom.NewPoint(0, 0, -2), 1.0},
	}

	for ti, tt := range tests {
		t.Run(t.Name()+strconv.Itoa(ti), func(t *testing.T) {
			intensity := IntensityAtAreaLight(light, tt.p, w)
			require.Equal(t, tt.expect, intensity)
		})
	}
}

func Test_AreaLightJittered_PassesIntensity(t *testing.T) {
	w := defaultWorld()

	corner := geom.NewPoint(-0.5, -0.5, -5)
	v1 := geom.NewVector(1, 0, 0)
	v2 := geom.NewVector(0, 1, 0)

	type args struct {
		p      geom.Tuple
		expect float64
	}

	tests := []args{
		{geom.NewPoint(0, 0, 2), 0.0},
		{geom.NewPoint(1, -1, 2), 0.5},
		{geom.NewPoint(1.5, 0, 2), 0.75}, // todo fix me , one extra hitting?
		{geom.NewPoint(1.25, 1.25, 3), 0.75},
		{geom.NewPoint(0, 0, -2), 1.0},
	}

	for ti, tt := range tests {
		t.Run(t.Name()+strconv.Itoa(ti), func(t *testing.T) {
			light := shapes.NewAreaLight(corner, v1, 2, v2, 2, colors.White(), shapes.NewJitterSequence(0.7, 0.3, 0.9, 0.1, 0.5))
			intensity := IntensityAtAreaLight(light, tt.p, w)
			require.Equal(t, tt.expect, intensity)
		})
	}
}

func Test_AreaLight_JitteredPoints(t *testing.T) {
	corner := geom.NewPoint(0, 0, 0)
	v1 := geom.NewVector(2, 0, 0)
	v2 := geom.NewVector(0, 0, 1)
	light := shapes.NewAreaLight(corner, v1, 4, v2, 2, colors.White(), shapes.NewJitterSequence(0.3, 0.7))

	type args struct {
		u      int
		v      int
		expect geom.Tuple
	}

	tests := []args{
		{0, 0, geom.NewPoint(0.15, 0, 0.35)},
		{1, 0, geom.NewPoint(0.65, 0, 0.35)},
		{0, 1, geom.NewPoint(0.15, 0, 0.85)},
		{2, 0, geom.NewPoint(1.15, 0, 0.35)},
		{3, 1, geom.NewPoint(1.65, 0, 0.85)},
	}

	for ti, tt := range tests {
		t.Run(t.Name()+strconv.Itoa(ti), func(t *testing.T) {
			pt := light.PointOnLight(tt.u, tt.v)
			require.Equal(t, tt.expect, pt)
		})
	}
}

func Test_Lighting_UsesIntensity(t *testing.T) {
	w := defaultWorld()
	require.Len(t, w.pointLights, 1)
	w.pointLights[0] = shapes.PointLight{
		Position:  geom.NewPoint(0, 0, -10),
		Intensity: colors.NewColor(1, 1, 1),
	}

	s := w.objects[0]
	m := materials.NewMaterial()
	m.Ambient = 0.1
	m.Diffuse = 0.9
	m.Specular = 0
	m.Color = colors.NewColor(1, 1, 1)
	s.SetMaterial(m)

	pt := geom.NewPoint(0, 0, -1)
	eyeV := geom.NewVector(0, 0, -1)
	normalV := geom.NewVector(0, 0, -1)

	type args struct {
		intensity float64
		expect    colors.Color
	}

	tests := []args{
		{1.0, colors.NewColor(1, 1, 1)},
		{0.5, colors.NewColor(0.55, 0.55, 0.55)},
		{0.0, colors.NewColor(0.1, 0.1, 0.1)},
	}

	for ti, tt := range tests {
		t.Run(t.Name()+strconv.Itoa(ti), func(t *testing.T) {
			c := shapes.Lighting(s.GetMaterial(), s, w.pointLights[0], pt, eyeV, normalV, tt.intensity)
			require.Equal(t, tt.expect, c)
		})
	}

}
