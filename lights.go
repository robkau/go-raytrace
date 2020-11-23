package main

import "math"

type pointLight struct {
	position  tuple
	intensity color
}

func newPointLight(p tuple, i color) pointLight {
	return pointLight{
		position:  p,
		intensity: i,
	}
}

// lighting calculates Phong lighting
func lighting(m material, l pointLight, p tuple, eyev tuple, nv tuple, shaded bool) color {
	effectiveColor := m.color.mul(l.intensity)
	lightv := l.position.sub(p).normalize()
	ambient := effectiveColor.mulBy(m.ambient)

	lightDotNormal := lightv.dot(nv)

	diffuse := color{}
	specular := color{}
	if lightDotNormal < 0 {
		diffuse = color{0, 0, 0}
		specular = color{0, 0, 0}
	} else {
		diffuse = effectiveColor.mulBy(m.diffuse).mulBy(lightDotNormal)
		reflectv := lightv.neg().reflect(nv)
		reflectDotEye := reflectv.dot(eyev)

		if reflectDotEye <= 0 {
			specular = color{0, 0, 0}
		} else {
			factor := math.Pow(reflectDotEye, m.shininess)
			specular = l.intensity.mulBy(m.specular).mulBy(factor)
		}
	}

	if shaded {
		return ambient
	} else {
		return ambient.add(diffuse).add(specular)
	}
}
