package scenes

import (
	"github.com/robkau/go-raytrace/lib/colors"
	"github.com/robkau/go-raytrace/lib/geom"
	"github.com/robkau/go-raytrace/lib/materials"
	"github.com/robkau/go-raytrace/lib/patterns"
	"github.com/robkau/go-raytrace/lib/shapes"
	"github.com/robkau/go-raytrace/lib/view"
)

func NewCappedCylinderScene(width int) (view.World, view.Camera) {
	w := view.NewWorld()
	cameraPos := geom.NewPoint(15, 15, 15)
	cameraLookingAt := geom.NewPoint(0, 5, 0)

	// cylinder
	cyl := shapes.NewCylinder(0, 5.9, true)
	cyl.M.Color = colors.Red()

	// glass sphere partially enveloping cylinder
	gs := shapes.NewSphere()
	gs.SetMaterial(materials.NewGlassMaterial())
	gs.SetTransform(geom.Translate(0, 7, 0).MulX4Matrix(geom.Scale(2.6, 2.6, 2.6)))

	// with a conical hat on top
	cone := shapes.NewUnitCone(true)
	cone.SetTransform(geom.Scale(1.3, 1, 1.3).MulX4Matrix(geom.Translate(0, 10.4, 0)))
	m := cone.GetMaterial()
	m.Pattern = patterns.NewRingPattern(patterns.NewSolidColorPattern(colors.RandomAnyColor()), patterns.NewSolidColorPattern(colors.RandomAnyColor()))
	m.Pattern.SetTransform(geom.Scale(0.2, 0.2, 0.2))
	cone.SetMaterial(m)

	// floor and ceiling as one cube
	// (the ceiling is too short but it looks cool)
	var floorAndCeiling = sizedCubeAt(0, 10, 0, 100, 10, 100)
	m = floorAndCeiling.GetMaterial()
	m.Color = colors.Brown()
	m.Reflective = 0
	m.Transparency = 0
	floorAndCeiling.SetMaterial(m)

	// walls as another cube
	var walls = sizedCubeAt(0, 0, 0, 20, 100, 20)
	m = walls.GetMaterial()
	m.Reflective = 0
	m.Transparency = 0
	m.Color = colors.Blue()
	walls.SetMaterial(m)

	// light above
	w.AddLight(shapes.NewPointLight(geom.NewPoint(5, 10, -3), colors.NewColor(1.9, 1.4, 1.4)))

	w.AddObject(cyl)
	w.AddObject(gs)
	w.AddObject(cone)
	w.AddObject(walls)
	w.AddObject(floorAndCeiling)

	c := view.NewCamera(width, width, 0.45)
	c.Transform = geom.ViewTransform(cameraPos,
		cameraLookingAt,
		geom.UpVector())

	return w, c
}
