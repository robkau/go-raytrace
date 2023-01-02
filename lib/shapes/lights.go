package shapes

import (
	"github.com/robkau/go-raytrace/lib/colors"
	"github.com/robkau/go-raytrace/lib/geom"
	"github.com/robkau/go-raytrace/lib/materials"
	"math"
)

type Light interface {
	GetIntensity() colors.Color
	GetPosition() geom.Tuple
	// todo unify Intensity functions
}

type PointLight struct {
	Position  geom.Tuple
	Intensity colors.Color
}

func (p PointLight) GetIntensity() colors.Color {
	return p.Intensity
}

func (p PointLight) GetPosition() geom.Tuple {
	return p.Position
}

func NewPointLight(p geom.Tuple, i colors.Color) PointLight {
	return PointLight{
		Position:  p,
		Intensity: i,
	}
}

type AreaLight struct {
	Corner    geom.Tuple
	UVec      geom.Tuple
	USteps    int
	VVec      geom.Tuple
	VSteps    int
	Samples   int
	Intensity colors.Color
	Position  geom.Tuple
}

func (a AreaLight) GetIntensity() colors.Color {
	return a.Intensity
}

func (a AreaLight) GetPosition() geom.Tuple {
	return a.Position
}

func NewAreaLight(corner geom.Tuple, uVec geom.Tuple, uSteps int, vVec geom.Tuple, vSteps int, intensity colors.Color) AreaLight {
	return AreaLight{
		Corner:    corner,
		UVec:      uVec.Div(float64(uSteps)),
		USteps:    uSteps,
		VVec:      vVec.Div(float64(vSteps)),
		VSteps:    vSteps,
		Samples:   uSteps * vSteps,
		Intensity: intensity,
		Position:  corner.Add(uVec.Div(2)).Add(vVec.Div(2)),
	}
}

func (a AreaLight) PointOnLight(u, v int) geom.Tuple {
	return a.Corner.Add(a.UVec.Mul(float64(u) + 0.5)).Add(a.VVec.Mul(float64(v) + 0.5))
}

// lighting calculates Phong lighting
func Lighting(m materials.Material, s Shape, l Light, p geom.Tuple, eyev geom.Tuple, nv geom.Tuple, intensity float64) colors.Color {
	if !s.GetShaded() {
		// override intensity, cannot be shaded in any way
		intensity = 1.0
	}

	var materialColor colors.Color
	if m.Pattern != nil {
		materialColor = m.Pattern.ColorAtShape(s.WorldToObject, p)
	} else {
		materialColor = m.Color
	}

	effectiveColor := materialColor.Mul(l.GetIntensity())
	lightv := l.GetPosition().Sub(p).Normalize()
	ambient := effectiveColor.MulBy(m.Ambient)

	if intensity == 0 {
		// object is in the shade, no other lighting to calculate
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
			specular = l.GetIntensity().MulBy(m.Specular).MulBy(factor)
		}
	}

	return ambient.Add(diffuse.MulBy(intensity)).Add(specular.MulBy(intensity))
}
