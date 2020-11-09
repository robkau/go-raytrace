package main

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func Test_DefaultMaterial(t *testing.T) {
	m := newMaterial()

	assert.Equal(t, color{1, 1, 1}, m.color)
	assert.Equal(t, 0.1, m.ambient)
	assert.Equal(t, 0.9, m.diffuse)
	assert.Equal(t, 0.9, m.specular)
	assert.Equal(t, 200.0, m.shininess)
}

func Test_Light_Eye_Inline(t *testing.T) {
	m := newMaterial()
	pos := newPoint(0, 0, 0)
	eyev := newVector(0, 0, -1)
	nv := newVector(0, 0, -1)
	light := newPointLight(newPoint(0, 0, -10), color{1, 1, 1})

	r := lighting(m, light, pos, eyev, nv)

	assert.Equal(t, color{1.9, 1.9, 1.9}, r)
}

func Test_Light_Eye_Offset45(t *testing.T) {
	m := newMaterial()
	pos := newPoint(0, 0, 0)
	eyev := newVector(0, math.Sqrt2/2, -math.Sqrt2/2)
	nv := newVector(0, 0, -1)
	light := newPointLight(newPoint(0, 0, -10), color{1, 1, 1})

	r := lighting(m, light, pos, eyev, nv)

	assert.Equal(t, color{1.0, 1.0, 1.0}, r)
}

func Test_Light_Offset45_Eye(t *testing.T) {
	m := newMaterial()
	pos := newPoint(0, 0, 0)
	eyev := newVector(0, 0, -1)
	nv := newVector(0, 0, -1)
	light := newPointLight(newPoint(0, 10, -10), color{1, 1, 1})

	r := lighting(m, light, pos, eyev, nv).roundTo(4)

	assert.Equal(t, color{0.7364, 0.7364, 0.7364}, r)
}

func Test_Light_Eye_Reflected(t *testing.T) {
	m := newMaterial()
	pos := newPoint(0, 0, 0)
	eyev := newVector(0, -math.Sqrt2/2, -math.Sqrt2/2)
	nv := newVector(0, 0, -1)
	light := newPointLight(newPoint(0, 10, -10), color{1, 1, 1})

	r := lighting(m, light, pos, eyev, nv).roundTo(4)

	assert.Equal(t, color{1.6364, 1.6364, 1.6364}, r)
}

func Test_Eye_LightBehindSurface(t *testing.T) {
	m := newMaterial()
	pos := newPoint(0, 0, 0)
	eyev := newVector(0, 0, -1)
	nv := newVector(0, 0, -1)
	light := newPointLight(newPoint(0, 0, 10), color{1, 1, 1})

	r := lighting(m, light, pos, eyev, nv)

	assert.Equal(t, color{0.1, 0.1, 0.1}, r)
}
