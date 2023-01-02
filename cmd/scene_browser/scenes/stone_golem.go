package scenes

import (
	"github.com/robkau/go-raytrace/lib/colors"
	"github.com/robkau/go-raytrace/lib/geom"
	"github.com/robkau/go-raytrace/lib/materials"
	"github.com/robkau/go-raytrace/lib/parse"
	"github.com/robkau/go-raytrace/lib/patterns"
	"github.com/robkau/go-raytrace/lib/shapes"
	"github.com/robkau/go-raytrace/lib/view"
	"log"
	"math"
)

func NewStoneGolemScene() *Scene {
	w := view.NewWorld()
	cameraPos := geom.NewPoint(4, 3, 7)
	cameraLookingAt := geom.NewPoint(0, 1, 0)

	g, err := parse.ParseObjFile("data/obj/stone.obj")
	if err != nil {
		log.Fatalf("failed parsing obj file: %s", err.Error())
	}
	// todo do this inside parsing and scale for each dimension and scale by largest required
	g.SetTransform(g.GetTransform().MulX4Matrix(geom.Scale(2/g.BoundsOf().Max.Y, 2/g.BoundsOf().Max.Y, 2/g.BoundsOf().Max.Y)).MulX4Matrix(geom.Translate(0, 4.7, 0)))
	m := materials.Material{}
	m.Pattern = patterns.NewSolidColorPattern(colors.White())
	m.Ambient = 0.2
	m.Diffuse = 0.2
	m.Specular = 0.1
	g.SetMaterial(m)

	lizard, err := parse.ParseObjFile("data/obj/LizardFolkOBJ.obj")
	if err != nil {
		log.Fatalf("failed parsing obj file: %s", err.Error())
	}
	// todo do this inside parsing and scale for each dimension and scale by largest required
	lizard.SetTransform(lizard.GetTransform().MulX4Matrix(geom.Scale(2/lizard.BoundsOf().Max.Y, 2/lizard.BoundsOf().Max.Y, 2/lizard.BoundsOf().Max.Y)).MulX4Matrix(geom.Translate(8, 4.7, 0)).MulX4Matrix(geom.RotateY(math.Pi)))
	m = materials.Material{}
	m.Pattern = patterns.NewSolidColorPattern(colors.Green())
	m.Ambient = 0.2
	m.Diffuse = 0.2
	m.Specular = 0.1
	lizard.SetMaterial(m)

	sphere := shapes.NewSphere()
	sphere.SetMaterial(materials.NewGlassMaterial())
	sphere.SetTransform(geom.Translate(0, 2, 2.5))

	car, err := parse.ParseObjFile("data/obj/uploads_files_3205191_supra.obj")
	if err != nil {
		log.Fatalf("failed parsing obj file: %s", err.Error())
	}
	car.SetTransform(geom.Translate(1, 0, 0))
	m = materials.Material{}
	m.Pattern = patterns.NewSolidColorPattern(colors.Red())
	m.Ambient = 0.2
	m.Diffuse = 0.2
	m.Specular = 0.1
	m.Reflective = 0
	// todo wrong when transparency != 0 ?
	m.Transparency = 0
	car.SetMaterial(m)

	dragon, err := parse.ParseObjFile("data/obj/dragon.obj")
	if err != nil {
		log.Fatalf("failed parsing obj file: %s", err.Error())
	}
	dragon.SetTransform(dragon.GetTransform().MulX4Matrix(geom.Scale(2/dragon.BoundsOf().Max.Y, 2/dragon.BoundsOf().Max.Y, 2/dragon.BoundsOf().Max.Y)).MulX4Matrix(geom.Translate(-4.5, 0, 0)).MulX4Matrix(geom.RotateY(math.Pi / 1.8)))
	m = materials.Material{}
	m.Pattern = patterns.NewSolidColorPattern(colors.RandomAnyColor())
	m.Ambient = 0.2
	m.Diffuse = 0.2
	m.Specular = 0.1
	m.Reflective = 0
	// todo wrong when transparency != 0 ?
	m.Transparency = 0
	dragon.SetMaterial(m)

	// floor and ceiling as one cube
	var floorAndCeiling = sizedCubeAt(0, 10, 0, 100, 10, 100)
	m = floorAndCeiling.GetMaterial()
	m.Color = colors.Brown()
	m.Reflective = 0
	m.Transparency = 0
	floorAndCeiling.SetMaterial(m)

	// walls as another cube
	var walls = sizedCubeAt(0, 0, 0, 25, 100, 25)
	m = walls.GetMaterial()
	m.Reflective = 0
	m.Transparency = 0
	m.Color = colors.Blue()
	walls.SetMaterial(m)

	w.AddObject(g)
	w.AddObject(lizard)
	//w.AddObject(sphere)
	w.AddObject(car)
	w.AddObject(dragon)
	w.AddObject(floorAndCeiling)
	w.AddObject(walls)

	//w.AddLight(shapes.NewPointLight(geom.NewPoint(0, 13, 0), colors.NewColor(1.9, 1.4, 1.4)))
	w.AddLight(shapes.NewPointLight(cameraPos, colors.NewColor(1.9, 1.4, 1.4)))

	w.Divide(64)

	return NewScene(
		w,
		CameraLocation{
			At:        cameraPos,
			LookingAt: cameraLookingAt,
		},
		CameraLocation{
			At:        cameraPos.Add(geom.NewPoint(0, 6, -4)),
			LookingAt: cameraLookingAt,
		},
		CameraLocation{
			At:        cameraPos.Add(geom.NewPoint(-2, 3, -1)),
			LookingAt: cameraLookingAt,
		},
	)
}
