package shapes

import (
	"github.com/robkau/go-raytrace/lib/colors"
	"github.com/robkau/go-raytrace/lib/geom"
	"github.com/robkau/go-raytrace/lib/materials"
	"github.com/robkau/go-raytrace/lib/patterns"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math"
	"strconv"
	"testing"
)

func Test_Light_Eye_Inline(t *testing.T) {
	m := materials.NewMaterial()
	pos := geom.ZeroPoint()
	eyev := geom.NewVector(0, 0, -1)
	nv := geom.NewVector(0, 0, -1)
	light := NewPointLight(geom.NewPoint(0, 0, -10), colors.White())

	r := Lighting(m, NewSphere(), light, pos, eyev, nv, 1.0)

	assert.True(t, colors.NewColor(1.9, 1.9, 1.9).Equal(r))
}

func Test_Light_Eye_Offset45(t *testing.T) {
	m := materials.NewMaterial()
	pos := geom.ZeroPoint()
	eyev := geom.NewVector(0, math.Sqrt2/2, -math.Sqrt2/2)
	nv := geom.NewVector(0, 0, -1)
	light := NewPointLight(geom.NewPoint(0, 0, -10), colors.White())

	r := Lighting(m, NewSphere(), light, pos, eyev, nv, 1.0)

	assert.Equal(t, colors.White(), r)
}

func Test_Light_Offset45_Eye(t *testing.T) {
	m := materials.NewMaterial()
	pos := geom.ZeroPoint()
	eyev := geom.NewVector(0, 0, -1)
	nv := geom.NewVector(0, 0, -1)
	light := NewPointLight(geom.NewPoint(0, 10, -10), colors.White())

	r := Lighting(m, NewSphere(), light, pos, eyev, nv, 1.0).RoundTo(4)

	assert.Equal(t, colors.NewColor(0.7364, 0.7364, 0.7364), r)
}

func Test_Light_Eye_Reflected(t *testing.T) {
	m := materials.NewMaterial()
	pos := geom.ZeroPoint()
	eyev := geom.NewVector(0, -math.Sqrt2/2, -math.Sqrt2/2)
	nv := geom.NewVector(0, 0, -1)
	light := NewPointLight(geom.NewPoint(0, 10, -10), colors.White())

	r := Lighting(m, NewSphere(), light, pos, eyev, nv, 1.0).RoundTo(4)

	assert.Equal(t, colors.NewColor(1.6364, 1.6364, 1.6364), r)
}

func Test_Eye_LightBehindSurface(t *testing.T) {
	m := materials.NewMaterial()
	pos := geom.ZeroPoint()
	eyev := geom.NewVector(0, 0, -1)
	nv := geom.NewVector(0, 0, -1)
	light := NewPointLight(geom.NewPoint(0, 0, 10), colors.White())

	r := Lighting(m, NewSphere(), light, pos, eyev, nv, 1.0)

	assert.Equal(t, colors.NewColor(0.1, 0.1, 0.1), r)
}

func Test_Eye_SurfaceShaded(t *testing.T) {
	m := materials.NewMaterial()
	pos := geom.ZeroPoint()
	eyev := geom.NewVector(0, 0, -1)
	normalv := geom.NewVector(0, 0, -1)
	light := NewPointLight(geom.NewPoint(0, 0, -10), colors.White())

	result := Lighting(m, NewSphere(), light, pos, eyev, normalv, 0.0)

	assert.Equal(t, colors.NewColor(0.1, 0.1, 0.1), result)
}

func Test_HitShouldOffsetPoint(t *testing.T) {
	r := geom.RayWith(geom.NewPoint(0, 0, -5), geom.NewVector(0, 0, 1))
	shape := NewSphere()
	shape.SetTransform(geom.Translate(0, 0, 1))
	i := NewIntersection(5, shape)
	comps := i.Compute(r, NewIntersections())

	assert.True(t, comps.OverPoint.Z < -geom.FloatComparisonEpsilon/2)
	assert.True(t, comps.point.Z > comps.OverPoint.Z)
}

func Test_Lighting_WithPattern(t *testing.T) {
	m := materials.NewMaterial()
	sp := patterns.NewStripePattern(patterns.NewSolidColorPattern(colors.White()), patterns.NewSolidColorPattern(colors.Black()))
	m.Pattern = sp
	m.Ambient = 1
	m.Diffuse = 0
	m.Specular = 0
	eyeV := geom.NewVector(0, 0, -1)
	normalV := geom.NewVector(0, 0, -1)
	light := NewPointLight(geom.NewPoint(0, 0, -10), colors.White())

	c1 := Lighting(m, NewSphere(), light, geom.NewPoint(0.9, 0, 0), eyeV, normalV, 1.0)
	c2 := Lighting(m, NewSphere(), light, geom.NewPoint(1.1, 0, 0), eyeV, normalV, 1.0)

	assert.Equal(t, colors.White(), c1)
	assert.Equal(t, colors.Black(), c2)
}

func Test_Lighting_SamplesAreaLight(t *testing.T) {
	corner := geom.NewPoint(-0.5, -0.5, -5)
	v1 := geom.NewVector(1, 0, 0)
	v2 := geom.NewVector(0, 1, 0)
	light := NewAreaLight(corner, v1, 2, v2, 2, colors.White(), NewJitterSequence(0.5))

	shape := NewSphere()
	m := materials.NewMaterial()
	m.Ambient = 0.1
	m.Diffuse = 0.9
	m.Specular = 0
	m.Color = colors.White()
	shape.SetMaterial(m)
	eye := geom.NewPoint(0, 0, -5)

	type args struct {
		p      geom.Tuple
		expect colors.Color
	}

	tests := []args{
		{geom.NewPoint(0, 0, -1), colors.NewColor(0.9965, 0.9965, 0.9965)},
		{geom.NewPoint(0, 0.7071, -0.7071), colors.NewColor(0.6232, 0.6232, 0.6232)},
	}

	for ti, tt := range tests {
		t.Run(t.Name()+strconv.Itoa(ti), func(t *testing.T) {
			eyeV := eye.Sub(tt.p).Normalize()
			normalV := geom.NewVector(tt.p.X, tt.p.Y, tt.p.Z)
			result := Lighting(m, shape, light, tt.p, eyeV, normalV, 1.0)
			require.Equal(t, tt.expect, result)
		})
	}

}
