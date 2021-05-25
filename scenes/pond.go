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
	// todo NewWater(xyzwh)
	var waterSurface shapes.Shape = shapes.NewPlane()
	m := waterSurface.GetMaterial()
	m.Color = colors.NewColor(0, 0.5, 1)
	m.Diffuse = 0.1
	m.Ambient = 0.1
	m.Specular = 0.5
	m.Shininess = 300
	m.Transparency = 0.5
	m.Reflective = 0.3
	m.RefractiveIndex = 1.13333
	waterSurface = waterSurface.SetMaterial(m)
	waterSurface = waterSurface.SetShadowless(true)
	waterSurface = waterSurface.SetShaded(false)

	// floor under water
	var dirtSurface shapes.Shape = shapes.NewPlane()
	m = dirtSurface.GetMaterial()
	m.Color = colors.NewColor(5, 5, 5)
	m.Specular = 0.1
	m.Diffuse = 0.7
	m.Pattern = patterns.NewPerlinPattern(patterns.NewCheckerPattern(patterns.NewSolidColorPattern(colors.Brown()), patterns.NewSolidColorPattern(colors.Black())), 1, 0.5, 3)
	dirtSurface = dirtSurface.SetMaterial(m)
	dirtSurface = dirtSurface.SetTransform(geom.Translate(0, -10, 0))

	// sphere underwater on floor
	var middle shapes.Shape = shapes.NewSphere()
	middle = middle.SetTransform(geom.Translate(-1.5, -5, -1.5))
	m = middle.GetMaterial()
	m.Pattern = patterns.NewCheckerPattern(patterns.NewSolidColorPattern(colors.Red()), patterns.NewSolidColorPattern(colors.White()))
	m.Pattern.SetTransform(geom.RotateX(math.Pi / 3).MulX4Matrix(geom.Scale(0.3, 0.3, 0.3)))
	m.Diffuse = 0.7
	m.Specular = 0.3
	middle = middle.SetMaterial(m)

	// halfway submerged sphere
	var floater shapes.Shape = shapes.NewSphere()
	floater = floater.SetTransform(geom.Translate(-2.5, -0.25, 1.5).MulX4Matrix(geom.Scale(2.3, 2.3, 2.3)))
	m = floater.GetMaterial()
	m.Pattern = patterns.NewCheckerPattern(patterns.NewSolidColorPattern(colors.Red()), patterns.NewSolidColorPattern(colors.White()))
	m.Diffuse = 0.7
	m.Specular = 0.3
	m.Pattern.SetTransform(geom.Scale(0.4, 0.4, 0.4))
	floater = floater.SetMaterial(m)
	floater = floater.SetShadowless(true)

	// light above plane
	w.AddLight(shapes.NewPointLight(geom.NewPoint(2, 12, -5), colors.NewColor(1.9, 1.4, 1.4)))
	w.AddObject(waterSurface)
	w.AddObject(dirtSurface)
	w.AddObject(middle)
	w.AddObject(floater)

	c := view.NewCamera(width, width, 0.45)
	c.Transform = geom.ViewTransform(geom.NewPoint(18, 5, -10),
		geom.NewPoint(0, 0, 0),
		geom.UpVector())

	return w, c
}
