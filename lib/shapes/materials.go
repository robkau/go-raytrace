package shapes

import (
	"go-raytrace/lib/colors"
	"go-raytrace/lib/patterns"
)

type Material struct {
	Color      colors.Color
	Pattern    patterns.Pattern
	Ambient    float64
	Diffuse    float64
	Specular   float64
	Shininess  float64
	Reflective float64
}

func NewMaterial() Material {
	return Material{
		Color:      colors.White(),
		Ambient:    0.1,
		Diffuse:    0.9,
		Specular:   0.9,
		Shininess:  200,
		Reflective: 0,
	}
}
