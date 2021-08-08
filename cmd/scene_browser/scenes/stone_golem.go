package scenes

import (
	"github.com/robkau/go-raytrace/lib/colors"
	"github.com/robkau/go-raytrace/lib/geom"
	"github.com/robkau/go-raytrace/lib/obj_parse"
	"github.com/robkau/go-raytrace/lib/patterns"
	"github.com/robkau/go-raytrace/lib/shapes"
	"github.com/robkau/go-raytrace/lib/view"
	"log"
)

func NewStoneGolemScene(width int) (view.World, view.Camera) {
	w := view.NewWorld()
	cameraPos := geom.NewPoint(15, 15, 15)
	cameraLookingAt := geom.NewPoint(0, 5, 0)

	g, err := obj_parse.ParseFile("data/obj/stone.obj")
	if err != nil {
		log.Fatalf("failed parsing obj file: %s", err.Error())
	}
	g = obj_parse.CollapseGroups(5, g)
	m := g.GetMaterial()
	m.Pattern = patterns.NewSolidColorPattern(colors.White())
	m.Ambient = 0.3
	m.Diffuse = 0.3
	g.SetMaterial(m)

	// floor and ceiling as one cube
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

	w.AddObject(g)
	w.AddObject(floorAndCeiling)
	w.AddObject(walls)

	w.AddLight(shapes.NewPointLight(geom.NewPoint(0, 13, 0), colors.NewColor(1.9, 1.4, 1.4)))

	c := view.NewCamera(width, width, 0.45)
	c.Transform = geom.ViewTransform(cameraPos,
		cameraLookingAt,
		geom.UpVector())

	return w, c
}
