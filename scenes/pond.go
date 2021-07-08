package scenes

import (
	"github.com/robkau/go-raytrace/lib/colors"
	"github.com/robkau/go-raytrace/lib/geom"
	"github.com/robkau/go-raytrace/lib/patterns"
	"github.com/robkau/go-raytrace/lib/shapes"
	"github.com/robkau/go-raytrace/lib/view"
	"math"
)

func NewPondScene(width int) (view.World, view.Camera) {
	w := view.NewWorld()

	// transparent plane
	waterSurface := shapes.NewPlane()
	m := waterSurface.GetMaterial()
	m.Color = colors.NewColor(0, 0.5, 1)
	m.Diffuse = 0.1
	m.Ambient = 0.1
	m.Specular = 0.5
	m.Shininess = 300
	m.Transparency = 0.5
	m.Reflective = 0.3
	m.RefractiveIndex = 1.13333
	waterSurface.SetMaterial(m)
	waterSurface.SetShadowless(true)
	waterSurface.SetShaded(false)

	// floor under water
	dirtSurface := shapes.NewPlane()
	m = dirtSurface.GetMaterial()
	m.Color = colors.NewColor(5, 5, 5)
	m.Specular = 0.1
	m.Diffuse = 0.7
	m.Pattern = patterns.NewPerlinPattern(patterns.NewCheckerPattern(patterns.NewSolidColorPattern(colors.Brown()), patterns.NewSolidColorPattern(colors.Black())), 1, 0.5, 3)
	dirtSurface.SetMaterial(m)
	dirtSurface.SetTransform(geom.Translate(0, -10, 0))

	// sphere underwater on floor
	middle := shapes.NewSphere()
	middle.SetTransform(geom.Translate(-1.5, -5, -1.5))
	m = middle.GetMaterial()
	m.Pattern = patterns.NewCheckerPattern(patterns.NewSolidColorPattern(colors.Red()), patterns.NewSolidColorPattern(colors.White()))
	m.Pattern.SetTransform(geom.RotateX(math.Pi / 3).MulX4Matrix(geom.Scale(0.3, 0.3, 0.3)))
	m.Diffuse = 0.7
	m.Specular = 0.3
	middle.SetMaterial(m)

	// halfway submerged sphere
	floater := shapes.NewSphere()
	floater.SetTransform(geom.Translate(-2.5, -0.25, 1.5).MulX4Matrix(geom.Scale(2.3, 2.3, 2.3)))
	m = floater.GetMaterial()
	m.Pattern = patterns.NewCheckerPattern(patterns.NewSolidColorPattern(colors.Red()), patterns.NewSolidColorPattern(colors.White()))
	m.Diffuse = 0.7
	m.Specular = 0.3
	m.Pattern.SetTransform(geom.Scale(0.4, 0.4, 0.4))
	floater.SetMaterial(m)
	floater.SetShadowless(true)

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
