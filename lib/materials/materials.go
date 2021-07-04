package materials

import (
	"go-raytrace/lib/colors"
	"go-raytrace/lib/patterns"
)

type Material struct {
	Color           colors.Color
	Pattern         patterns.Pattern
	Ambient         float64
	Diffuse         float64
	Specular        float64
	Shininess       float64
	Reflective      float64
	Transparency    float64
	RefractiveIndex float64
}

func NewMaterial() Material {
	return Material{
		Color:           colors.White(),
		Ambient:         0.1,
		Diffuse:         0.9,
		Specular:        0.9,
		Shininess:       200,
		Reflective:      0,
		Transparency:    0,
		RefractiveIndex: 1,
	}
}

func NewGlassMaterial() Material {
	m := NewMaterial()
	m.Diffuse = 0
	m.Ambient = 0
	m.Specular = 0.9
	m.Shininess = 300
	m.Transparency = 0.95
	m.Reflective = 0.9
	m.RefractiveIndex = 1.5
	return m
}
