package scenes

import (
	"github.com/robkau/go-raytrace/lib/colors"
	"github.com/robkau/go-raytrace/lib/geom"
	"github.com/robkau/go-raytrace/lib/materials"
	"github.com/robkau/go-raytrace/lib/patterns"
	"github.com/robkau/go-raytrace/lib/shapes"
	"github.com/robkau/go-raytrace/lib/view"
)

func NewHollowGlassSphereScene() (*view.World, []CameraLocation) {
	w := view.NewWorld()

	wall := shapes.NewPlane()
	m := wall.GetMaterial()
	m.Pattern = patterns.NewSprayPaintPattern(patterns.NewCheckerPattern(patterns.NewSolidColorPattern(colors.NewColor(0.15, 0.15, 0.15)), patterns.NewSolidColorPattern(colors.NewColor(0.85, 0.85, 0.85))), 0.025)
	m.Ambient = 0.8
	m.Diffuse = 0.2
	m.Specular = 0
	wall.SetMaterial(m)
	wall.SetTransform(geom.Translate(0, 0, 10).MulX4Matrix(geom.RotateX(1.5708)))

	ball := shapes.NewSphere()
	ball.SetMaterial(materials.NewGlassMaterial())

	hollowCenter := shapes.NewSphere()
	hollowCenter.SetMaterial(materials.NewGlassMaterial())
	hollowCenter.SetTransform(geom.Scale(0.5, 0.5, 0.5))
	m = hollowCenter.GetMaterial()
	m.Color = colors.NewColor(1, 1, 1)
	m.Diffuse = 0
	m.Ambient = 0
	m.Specular = 0.9
	m.Shininess = 300
	m.Transparency = 0.9
	m.Reflective = 0.9
	m.RefractiveIndex = 1.0000034
	hollowCenter.SetMaterial(m)

	w.AddObject(wall)
	w.AddObject(ball)
	w.AddObject(hollowCenter)
	w.AddPointLight(shapes.NewPointLight(geom.NewPoint(2, 10, -5), colors.NewColor(0.9, 0.9, 0.9)))

	cameraPos := geom.NewPoint(0, 0, -5)
	cameraLookingAt := geom.ZeroPoint()

	w.Divide(8)
	return w, []CameraLocation{CameraLocation{cameraPos, cameraLookingAt}}
}
