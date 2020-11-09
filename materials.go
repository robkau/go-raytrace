package main

type material struct {
	color     color
	ambient   float64
	diffuse   float64
	specular  float64
	shininess float64
}

func newMaterial() material {
	return material{
		color:     color{1, 1, 1},
		ambient:   0.1,
		diffuse:   0.9,
		specular:  0.9,
		shininess: 200,
	}
}
