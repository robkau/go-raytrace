package shapes

import (
	"github.com/robkau/go-raytrace/lib/colors"
	"github.com/robkau/go-raytrace/lib/geom"
	"github.com/robkau/go-raytrace/lib/materials"
	"math"
)

type PointLight struct {
	Position  geom.Tuple
	Intensity colors.Color
}

func NewPointLight(p geom.Tuple, i colors.Color) PointLight {
	return PointLight{
		Position:  p,
		Intensity: i,
	}
}

// lighting calculates Phong lighting
func Lighting(m materials.Material, s Shape, l PointLight, p geom.Tuple, eyev geom.Tuple, nv geom.Tuple, shaded bool) colors.Color {
	var materialColor colors.Color
	if m.Pattern != nil {
		materialColor = m.Pattern.ColorAtShape(s.GetTransform(), p)
	} else {
		materialColor = m.Color
	}

	effectiveColor := materialColor.Mul(l.Intensity)
	lightv := l.Position.Sub(p).Normalize()
	ambient := effectiveColor.MulBy(m.Ambient)

	if shaded && s.GetShaded() {
		// object is in the shade, and able to be shaded. no lighting to calculate
		return ambient
	}

	lightDotNormal := lightv.Dot(nv)

	diffuse := colors.NewColor(0, 0, 0)
	specular := colors.NewColor(0, 0, 0)
	if lightDotNormal < 0 {
		diffuse = colors.Black()
		specular = colors.Black()
	} else {
		diffuse = effectiveColor.MulBy(m.Diffuse).MulBy(lightDotNormal)
		reflectv := lightv.Neg().Reflect(nv)
		reflectDotEye := reflectv.Dot(eyev)

		if reflectDotEye <= 0 {
			specular = colors.Black()
		} else {
			factor := math.Pow(reflectDotEye, m.Shininess)
			specular = l.Intensity.MulBy(m.Specular).MulBy(factor)
		}
	}

	return ambient.Add(diffuse).Add(specular)
}
