package materials

import (
	"github.com/robkau/go-raytrace/lib/colors"
	"github.com/robkau/go-raytrace/lib/patterns"
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
	m.Color = colors.NewColor(0, 0, 0.2)
	m.Diffuse = 0.2
	m.Ambient = 0
	m.Specular = 0.9
	m.Shininess = 300
	m.Transparency = 0.95
	m.Reflective = 0.9
	m.RefractiveIndex = 1.5
	return m
}

func IsZeroMaterial(m Material) bool {
	return m.Color == colors.Color{} &&
		m.Ambient == 0 &&
		m.Diffuse == 0 &&
		m.Specular == 0 &&
		m.Shininess == 0 &&
		m.Reflective == 0 &&
		m.Transparency == 0 &&
		m.RefractiveIndex == 0
}

func ZeroMaterial() Material {
	return Material{}
}
