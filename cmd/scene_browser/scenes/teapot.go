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

func NewTeapotScene() (*view.World, []CameraLocation) {
	w := view.NewWorld()
	cameraPos := geom.NewPoint(85, 10, -10)
	cameraLookingAt := geom.NewPoint(0, 5, -10)

	g, err := parse.ParseObjFile("data/obj/teapot_lowpoly_no_face_normals.obj")
	if err != nil {
		log.Fatalf("failed parsing obj file: %s", err.Error())
	}
	g.SetTransform(g.GetTransform().MulX4Matrix(geom.Translate(0, 0, -20)).MulX4Matrix(geom.RotateX(-math.Pi / 2)))
	m := materials.NewMaterial()
	m.Pattern = patterns.NewSolidColorPattern(colors.Green())
	m.Ambient = 0.3
	m.Diffuse = 0.3
	m.Specular = 0.3
	m.Shininess = 0.2
	g.SetMaterial(m)

	g2, err := parse.ParseObjFile("data/obj/teapot_lowpoly.obj")
	if err != nil {
		log.Fatalf("failed parsing obj file: %s", err.Error())
	}
	g2.SetTransform(g2.GetTransform().MulX4Matrix(geom.Translate(0, 0, 0)).MulX4Matrix(geom.RotateX(-math.Pi / 2)))
	m = materials.NewMaterial()
	m.Pattern = patterns.NewSolidColorPattern(colors.Green())
	m.Ambient = 0.3
	m.Diffuse = 0.3
	m.Specular = 0.3
	m.Shininess = 0.2
	g2.SetMaterial(m)

	// floor and ceiling as one cube
	var floorAndCeiling = sizedCubeAt(0, 10, 0, 100, 30, 100)
	m = floorAndCeiling.GetMaterial()
	m.Color = colors.Brown()
	m.Reflective = 0
	m.Transparency = 0
	m.Specular = 0
	floorAndCeiling.SetMaterial(m)

	// walls as another cube
	var walls = sizedCubeAt(0, 0, 0, 20, 100, 20)
	m = walls.GetMaterial()
	m.Reflective = 0
	m.Transparency = 0
	m.Color = colors.Blue()
	walls.SetMaterial(m)

	w.AddObject(g)
	w.AddObject(g2)
	w.AddObject(floorAndCeiling)
	//w.AddObject(walls)

	w.Divide(16)

	w.AddPointLight(shapes.NewPointLight(cameraPos, colors.NewColor(1.9, 1.4, 1.4)))

	return w, []CameraLocation{CameraLocation{cameraPos, cameraLookingAt}}
}
