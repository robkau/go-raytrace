package scenes

import (
	"github.com/robkau/go-raytrace/lib/colors"
	"github.com/robkau/go-raytrace/lib/geom"
	"github.com/robkau/go-raytrace/lib/materials"
	"github.com/robkau/go-raytrace/lib/patterns"
	"github.com/robkau/go-raytrace/lib/shapes"
	"github.com/robkau/go-raytrace/lib/view"
	"math"
)

func NewWavyCarpetSpheres() *Scene {
	var floor shapes.Shape = shapes.NewPlane()
	m := floor.GetMaterial()
	m.Color = colors.NewColor(1, 0.9, 0.9)
	// alternating stripe patterns
	p1 := patterns.NewStripePattern(patterns.NewSolidColorPattern(colors.RandomAnyColor()), patterns.NewSolidColorPattern(colors.RandomAnyColor()))
	p1.SetTransform(geom.RotateY(math.Pi / 2).MulX4Matrix(geom.Scale(0.2, 0.2, 0.2)))
	p2 := patterns.NewStripePattern(patterns.NewSolidColorPattern(colors.RandomAnyColor()), patterns.NewSolidColorPattern(colors.RandomAnyColor()))
	p2.SetTransform(geom.Scale(0.2, 0.2, 0.2))
	// with perlin displacement in a checkerboard
	p3 := patterns.NewCheckerPattern(patterns.NewPerlinPattern(p1, 0.3, 0.2, 3), patterns.NewPerlinPattern(p2, 0.3, 0.7, 7))
	// with spraypaint style
	p4 := patterns.NewSprayPaintPattern(p3, 0.04)
	m.Pattern = p4
	m.Specular = 0.2
	floor.SetMaterial(m)

	var middle shapes.Shape = shapes.NewSphere()
	middle.SetTransform(geom.Translate(-0.5, 1, 0.5))
	m = middle.GetMaterial()
	m.Pattern = patterns.NewPerlinPattern(patterns.NewStripePattern(patterns.NewSolidColorPattern(colors.RandomColor()), patterns.NewSolidColorPattern(colors.RandomAnyColor())), 0.5, 0.6, 4)
	m.Pattern.SetTransform(geom.RotateX(math.Pi / 3).MulX4Matrix(geom.Scale(0.3, 0.3, 0.3)))
	m.Color = colors.NewColor(0.1, 1, 0.5)
	m.Diffuse = 0.7
	m.Specular = 0.3
	middle.SetMaterial(m)

	var right shapes.Shape = shapes.NewSphere()
	right.SetTransform(geom.Translate(1.5, 0.5, -0.5).MulX4Matrix(geom.Scale(0.5, 0.5, 0.5)))
	m = right.GetMaterial()
	m.Color = colors.NewColor(0.5, 1, 0.1)
	m.Diffuse = 0.7
	m.Specular = 0.3
	m.Pattern = patterns.NewPerlinPattern(patterns.NewStripePattern(patterns.NewSolidColorPattern(colors.RandomColor()), patterns.NewSolidColorPattern(colors.RandomAnyColor())), 0.4, 0.7, 7)
	m.Pattern.SetTransform(geom.RotateY(math.Pi / 3).MulX4Matrix(geom.Scale(0.4, 0.4, 0.4)))
	right.SetMaterial(m)

	var left shapes.Shape = shapes.NewSphere()
	left.SetTransform(geom.Translate(-1.3, 2.4, -0.75).MulX4Matrix(geom.Scale(0.23, 0.23, 0.23)))
	m = left.GetMaterial()
	m.Color = colors.NewColor(1, 0.8, 0.1)
	m.Diffuse = 0.7
	m.Specular = 0.3
	m.Pattern = patterns.NewPerlinPattern(patterns.NewRingPattern(patterns.NewSolidColorPattern(colors.RandomAnyColor()), patterns.NewSolidColorPattern(colors.RandomColor())), 0.3, 0.8, 3)
	m.Pattern.SetTransform(geom.RotateZ(math.Pi / 3).MulX4Matrix(geom.Scale(0.5, 0.5, 0.5)))
	left.SetMaterial(m)

	var glass shapes.Shape = shapes.NewSphere()
	glass.SetMaterial(materials.NewGlassMaterial())
	glass.SetTransform(geom.Translate(-1.3, 2.4, 3.75).MulX4Matrix(geom.RotateY(-math.Pi / 6.2).MulX4Matrix(geom.RotateZ(math.Pi / 8))).MulX4Matrix(geom.Scale(1.73, 1.73, 0.13)))
	m = glass.GetMaterial()
	m.Reflective = 1
	m.Transparency = 0
	glass.SetMaterial(m)

	g := shapes.NewGroup()
	g.AddChild(floor)
	g.AddChild(middle)
	g.AddChild(right)
	g.AddChild(left)
	g.AddChild(glass)

	g.SetTransform(geom.RotateY(math.Pi / 15))

	w := view.NewWorld()
	w.AddObject(g)

	w.AddLight(shapes.NewPointLight(geom.NewPoint(-10, 10, -10), colors.White()))

	cameraPos := geom.NewPoint(2, 4, -3)
	cameraLookingAt := geom.NewPoint(0, 1, 0)

	return NewScene(w, CameraLocation{cameraPos, cameraLookingAt})
}
