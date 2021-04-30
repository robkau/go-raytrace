package shapes

import (
	"github.com/stretchr/testify/assert"
	"go-raytrace/lib/colors"
	"go-raytrace/lib/geom"
	"go-raytrace/lib/patterns"
	"math"
	"testing"
)

func Test_DefaultMaterial(t *testing.T) {
	m := NewMaterial()

	assert.Equal(t, colors.White(), m.Color)
	assert.Equal(t, 0.1, m.Ambient)
	assert.Equal(t, 0.9, m.Diffuse)
	assert.Equal(t, 0.9, m.Specular)
	assert.Equal(t, 200.0, m.Shininess)
	assert.Equal(t, 0.0, m.Reflective)
}

func Test_Light_Eye_Inline(t *testing.T) {
	m := NewMaterial()
	pos := geom.NewPoint(0, 0, 0)
	eyev := geom.NewVector(0, 0, -1)
	nv := geom.NewVector(0, 0, -1)
	light := NewPointLight(geom.NewPoint(0, 0, -10), colors.White())

	r := Lighting(m, NewSphere(), light, pos, eyev, nv, false)

	assert.Equal(t, colors.NewColor(1.9, 1.9, 1.9), r)
}

func Test_Light_Eye_Offset45(t *testing.T) {
	m := NewMaterial()
	pos := geom.NewPoint(0, 0, 0)
	eyev := geom.NewVector(0, math.Sqrt2/2, -math.Sqrt2/2)
	nv := geom.NewVector(0, 0, -1)
	light := NewPointLight(geom.NewPoint(0, 0, -10), colors.White())

	r := Lighting(m, NewSphere(), light, pos, eyev, nv, false)

	assert.Equal(t, colors.White(), r)
}

func Test_Light_Offset45_Eye(t *testing.T) {
	m := NewMaterial()
	pos := geom.NewPoint(0, 0, 0)
	eyev := geom.NewVector(0, 0, -1)
	nv := geom.NewVector(0, 0, -1)
	light := NewPointLight(geom.NewPoint(0, 10, -10), colors.White())

	r := Lighting(m, NewSphere(), light, pos, eyev, nv, false).RoundTo(4)

	assert.Equal(t, colors.NewColor(0.7364, 0.7364, 0.7364), r)
}

func Test_Light_Eye_Reflected(t *testing.T) {
	m := NewMaterial()
	pos := geom.NewPoint(0, 0, 0)
	eyev := geom.NewVector(0, -math.Sqrt2/2, -math.Sqrt2/2)
	nv := geom.NewVector(0, 0, -1)
	light := NewPointLight(geom.NewPoint(0, 10, -10), colors.White())

	r := Lighting(m, NewSphere(), light, pos, eyev, nv, false).RoundTo(4)

	assert.Equal(t, colors.NewColor(1.6364, 1.6364, 1.6364), r)
}

func Test_Eye_LightBehindSurface(t *testing.T) {
	m := NewMaterial()
	pos := geom.NewPoint(0, 0, 0)
	eyev := geom.NewVector(0, 0, -1)
	nv := geom.NewVector(0, 0, -1)
	light := NewPointLight(geom.NewPoint(0, 0, 10), colors.White())

	r := Lighting(m, NewSphere(), light, pos, eyev, nv, false)

	assert.Equal(t, colors.NewColor(0.1, 0.1, 0.1), r)
}

func Test_Eye_SurfaceDhaded(t *testing.T) {
	m := NewMaterial()
	pos := geom.NewPoint(0, 0, 0)
	eyev := geom.NewVector(0, 0, -1)
	normalv := geom.NewVector(0, 0, -1)
	light := NewPointLight(geom.NewPoint(0, 0, -10), colors.White())
	shaded := true

	result := Lighting(m, NewSphere(), light, pos, eyev, normalv, shaded)

	assert.Equal(t, colors.NewColor(0.1, 0.1, 0.1), result)
}

func Test_HitShouldOffsetPoint(t *testing.T) {
	r := geom.RayWith(geom.NewPoint(0, 0, -5), geom.NewVector(0, 0, 1))
	shape := NewSphereWith(geom.Translate(0, 0, 1))
	i := NewIntersection(5, shape)
	comps := i.Compute(r)

	assert.True(t, comps.OverPoint.Z < -geom.FloatComparisonEpsilon/2)
	assert.True(t, comps.point.Z > comps.OverPoint.Z)
}

func Test_Lighting_WithPattern(t *testing.T) {
	m := NewMaterial()
	sp := patterns.NewStripePattern(patterns.NewSolidColorPattern(colors.White()), patterns.NewSolidColorPattern(colors.Black()))
	m.Pattern = sp
	m.Ambient = 1
	m.Diffuse = 0
	m.Specular = 0
	eyeV := geom.NewVector(0, 0, -1)
	normalV := geom.NewVector(0, 0, -1)
	light := NewPointLight(geom.NewPoint(0, 0, -10), colors.White())

	c1 := Lighting(m, NewSphere(), light, geom.NewPoint(0.9, 0, 0), eyeV, normalV, false)
	c2 := Lighting(m, NewSphere(), light, geom.NewPoint(1.1, 0, 0), eyeV, normalV, false)

	assert.Equal(t, colors.White(), c1)
	assert.Equal(t, colors.Black(), c2)
}
