package shapes

import (
	"github.com/robkau/go-raytrace/lib/colors"
	"github.com/robkau/go-raytrace/lib/geom"
	"github.com/robkau/go-raytrace/lib/materials"
	"math"
)

type Light interface {
	GetIntensity() colors.Color
	GetSamples() []geom.Tuple
	// todo unify Intensity functions
}

type PointLight struct {
	Position  geom.Tuple
	Intensity colors.Color
}

func (p PointLight) GetIntensity() colors.Color {
	return p.Intensity
}

func (p PointLight) GetSamples() []geom.Tuple {
	return []geom.Tuple{p.Position}
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
	Seq       Sequence
	Center    geom.Tuple
}

func (a AreaLight) GetIntensity() colors.Color {
	return a.Intensity
}

func (a AreaLight) GetSamples() []geom.Tuple {
	s := make([]geom.Tuple, 0, a.Samples)
	for i := 0; i < a.USteps; i++ {
		for j := 0; j < a.VSteps; j++ {
			s = append(s, a.PointOnLight(i, j))
		}
	}

	return s
}

func NewAreaLight(corner geom.Tuple, uVec geom.Tuple, uSteps int, vVec geom.Tuple, vSteps int, intensity colors.Color, seq Sequence) AreaLight {
	if seq == nil {
		seq = NewRandomSequence()
	}
	return AreaLight{
		Corner:    corner,
		UVec:      uVec.Div(float64(uSteps)),
		USteps:    uSteps,
		VVec:      vVec.Div(float64(vSteps)),
		VSteps:    vSteps,
		Samples:   uSteps * vSteps,
		Intensity: intensity,
		Seq:       seq,
		Center:    corner.Add(uVec.Div(2)).Add(vVec.Div(2)),
	}
}

func (a AreaLight) PointOnLight(u, v int) geom.Tuple {
	uJit := a.Seq.Next()
	vJit := a.Seq.Next()
	return a.Corner.Add(a.UVec.Mul(float64(u) + uJit)).Add(a.VVec.Mul(float64(v) + vJit))
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
	ambient := effectiveColor.MulBy(m.Ambient)
	if intensity == 0 {
		// object is in the shade, no other lighting to calculate
		return ambient
	}

	sum := colors.Black()

	numSamples := 0
	for _, sample := range l.GetSamples() {
		numSamples++
		lightv := sample.Sub(p).Normalize()
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
		sum = sum.Add(diffuse)
		sum = sum.Add(specular)
	}

	return ambient.Add(sum.MulBy(intensity / float64(numSamples))) // todo or intensity multiply all?
}
