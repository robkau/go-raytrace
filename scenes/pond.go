package scenes

import (
	"go-raytrace/lib/colors"
	"go-raytrace/lib/geom"
	"go-raytrace/lib/patterns"
	"go-raytrace/lib/shapes"
	"go-raytrace/lib/view"
	"math"
)

func NewPondScene(width int) (view.World, view.Camera) {
	w := view.NewWorld()

	// transparent plane
	var waterSurface shapes.Shape = shapes.NewPlane()
	m := waterSurface.GetMaterial()
	m.Color = colors.NewColor(0, 0.5, 1)
	m.Diffuse = 0.1
	m.Ambient = 0.1
	m.Specular = 0.5
	m.Shininess = 300
	m.Transparency = 0.3
	m.Reflective = 0.9
	m.RefractiveIndex = 1.233
	waterSurface = waterSurface.SetMaterial(m)
	waterSurface = waterSurface.SetShadowless(true)

	// floor under water
	var dirtSurface shapes.Shape = shapes.NewPlane()
	m = dirtSurface.GetMaterial()
	m.Color = colors.NewColor(5, 5, 5)
	m.Specular = 0.1
	m.Diffuse = 0.7
	m.Pattern = patterns.NewRingPattern(patterns.NewSolidColorPattern(colors.RandomColor()), patterns.NewSolidColorPattern(colors.RandomColor()))
	dirtSurface = dirtSurface.SetMaterial(m)
	dirtSurface = dirtSurface.SetTransform(geom.Translate(0, -10, 0))

	// sphere underwater on floor
	var middle shapes.Shape = shapes.NewSphere()
	middle = middle.SetTransform(geom.Translate(-1.5, -7, -1.5))
	m = middle.GetMaterial()
	m.Pattern = patterns.NewCheckerPattern(patterns.NewSolidColorPattern(colors.Red()), patterns.NewSolidColorPattern(colors.White()))
	m.Pattern.SetTransform(geom.RotateX(math.Pi / 3).MulX4Matrix(geom.Scale(0.3, 0.3, 0.3)))
	m.Diffuse = 0.7
	m.Specular = 0.3
	middle = middle.SetMaterial(m)

	// halfway submerged sphere
	var floater shapes.Shape = shapes.NewSphere()
	floater = floater.SetTransform(geom.Translate(-5.5, -0.25, 0.5).MulX4Matrix(geom.Scale(2.3, 2.3, 2.3)))
	m = floater.GetMaterial()
	m.Pattern = patterns.NewCheckerPattern(patterns.NewSolidColorPattern(colors.Red()), patterns.NewSolidColorPattern(colors.White()))
	m.Diffuse = 0.7
	m.Specular = 0.3
	floater = floater.SetMaterial(m)

	// light above plane
	w.AddLight(shapes.NewPointLight(geom.NewPoint(2, 12, -5), colors.NewColor(0.9, 0.9, 0.9)))
	w.AddObject(waterSurface)
	w.AddObject(dirtSurface)
	w.AddObject(middle)
	w.AddObject(floater)

	c := view.NewCamera(width, width, 0.45)
	c.Transform = geom.ViewTransform(geom.NewPoint(15, 7, -7),
		geom.NewPoint(0, 0, 0),
		geom.NewVector(0, 1, 0))

	return w, c
}
