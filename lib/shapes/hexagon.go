package shapes

import (
	"github.com/robkau/go-raytrace/lib/geom"
	"math"
)

// a hexagon is simply a group of spheres and cylinders
func NewHexagon() Shape {
	hex := NewGroup()
	for i := 0; i < 6; i++ {
		side := hexagonSide()
		side.SetTransform(geom.RotateY(float64(i) * math.Pi / 3))
		hex.AddChild(side)
	}
	return hex
}

func hexagonSide() Shape {
	g := NewGroup()
	g.AddChild(hexagonCorner())
	g.AddChild(hexagonEdge())
	return g
}

func hexagonCorner() Shape {
	c := NewSphere()
	c.SetTransform(geom.Translate(0, 0, -1).MulX4Matrix(geom.Scale(0.25, 0.25, 0.25)))
	return c
}

func hexagonEdge() Shape {
	e := NewCylinder(0, 1, false)
	e.SetTransform(geom.Translate(0, 0, -1).MulX4Matrix(geom.RotateY(-math.Pi / 6)).MulX4Matrix(geom.RotateZ(-math.Pi / 2)).MulX4Matrix(geom.Scale(0.25, 1, 0.25)))
	return e
}
