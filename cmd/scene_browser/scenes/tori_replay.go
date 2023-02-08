package scenes

import (
	"fmt"
	"github.com/robkau/coordinate_supplier"
	"github.com/robkau/go-raytrace/lib/colors"
	"github.com/robkau/go-raytrace/lib/geom"
	"github.com/robkau/go-raytrace/lib/materials"
	"github.com/robkau/go-raytrace/lib/parse"
	"github.com/robkau/go-raytrace/lib/patterns"
	"github.com/robkau/go-raytrace/lib/shapes"
	"github.com/robkau/go-raytrace/lib/view"
	canvas2 "github.com/robkau/go-raytrace/lib/view/canvas"
	"strings"
)

func NewToriReplayScene() (*view.World, []CameraLocation) {
	w := view.NewWorld()

	sceneSpacing := 6.5
	outerShellRadius := 3.15
	innerShellRadius := 2.65
	cameraDistance := sceneSpacing * 8
	wallDistance := cameraDistance * 1.1

	shellsPerLine := 3

	pr, err := parse.ParseReaderAsTori(strings.NewReader(parse.ReplayFile))
	if err != nil {
		panic("err parse")
	}

	// draw a glass sphere around each tori frame.
	cs, err := coordinate_supplier.NewCoordinateSupplierAtomic(coordinate_supplier.CoordinateSupplierOptions{
		Width:  shellsPerLine,
		Height: 10000,
		Depth:  1,
		Order:  coordinate_supplier.Asc,
		Repeat: false,
	})
	if err != nil {
		panic(err)
	}
	for range pr.P0Positions {
		ball := shapes.NewSphere()
		z, y, _, done := cs.Next()
		if done {
			panic("out of coordinates")
		}
		ball.SetMaterial(materials.NewGlassMaterial())
		m := ball.GetMaterial()
		//m.Specular = 0.2
		//m.Transparency = 0.9875
		ball.SetTransform(geom.Translate(0, 2+float64(y)*sceneSpacing, float64(z)*sceneSpacing).MulX4Matrix(geom.Scale(outerShellRadius, outerShellRadius, outerShellRadius)))
		ball.SetShadowless(true)

		hollowCenter := shapes.NewSphere()
		hollowCenter.SetMaterial(materials.NewGlassMaterial())
		hollowCenter.SetTransform(geom.Translate(0, 2+float64(y)*sceneSpacing, float64(z)*sceneSpacing).MulX4Matrix(geom.Scale(innerShellRadius, innerShellRadius, innerShellRadius)))
		hollowCenter.SetShadowless(true)
		m = hollowCenter.GetMaterial()
		//m.Color = colors.NewColor(1, 1, 1)
		//m.Diffuse = 0
		//m.Ambient = 0
		//m.Specular = 0.4
		//m.Shininess = 300
		//m.Transparency = 0.9875
		//m.Reflective = 0.9
		//m.RefractiveIndex = 1.0000034
		m.Diffuse = 0
		m.Ambient = 0
		m.Specular = 0.9
		m.Shininess = 300
		m.Transparency = 0.9
		m.Reflective = 0.9
		m.RefractiveIndex = 1.0000034
		hollowCenter.SetMaterial(m)

		w.AddObject(ball)
		w.AddObject(hollowCenter)
	}

	g := pr.AllScenes(shellsPerLine, sceneSpacing)
	w.AddObject(g)

	// camera points to center of displayed tori frames
	c := g.BoundsOf().Center()

	// skybox sphere
	var skybox = shapes.NewSphere()
	skybox.SetTransform(geom.Scale(wallDistance, wallDistance, wallDistance))
	m := skybox.GetMaterial()
	canvas, err := canvas2.CanvasFromPPMZipFile("data/ppm/satara_night_hdr.ppm.zip")
	if err != nil {
		panic(fmt.Sprintf("loading tokyo ppm to canvas: %s", err.Error()))
	}
	m.Pattern = patterns.NewTextureMapPattern(patterns.NewUVImage(canvas), patterns.SphericalMap)
	m.Ambient = 1
	m.Specular = 0
	m.Diffuse = 0
	skybox.SetMaterial(m)

	w.AddObject(skybox)

	//lizard, err := parse.ParseObjFile("data/obj/LizardFolkOBJ.obj")
	//if err != nil {
	//	log.Fatalf("failed parsing obj file: %s", err.Error())
	//}
	//// todo do this inside parsing and scale for each dimension and scale by largest required
	//lizard.SetTransform(lizard.GetTransform().MulX4Matrix(geom.Scale(2/lizard.BoundsOf().Max.Y, 2/lizard.BoundsOf().Max.Y, 2/lizard.BoundsOf().Max.Y)).MulX4Matrix(geom.Translate(8, 4.7, 0)).MulX4Matrix(geom.RotateY(math.Pi / 1.25)))
	//m = materials.Material{}
	//m.Pattern = patterns.NewSolidColorPattern(colors.Green())
	//m.Ambient = 0.15
	//m.Diffuse = 0.15
	//m.Specular = 0.1
	//m.Shininess = 50
	//m.Reflective = 0.1
	//lizard.SetMaterial(m)
	//lizard.SetTransform(geom.Translate(0, 9.85, 6.45).MulX4Matrix(geom.Scale(2.4, 2.4, 2.4).MulX4Matrix(lizard.GetTransform())))
	//w.AddObject(lizard)
	//
	//newLizard, err := parse.ParseObjFile("data/obj/LizardFolkOBJ.obj")
	//if err != nil {
	//	log.Fatalf("failed parsing obj file: %s", err.Error())
	//}
	//newLizard.SetTransform(geom.RotateY(math.Pi))
	//newLizard.SetTransform(lizard.GetTransform().Copy().MulX4Matrix(newLizard.GetTransform()))
	//newLizard.SetTransform(geom.Translate(0, 0, sceneSpacing).MulX4Matrix(newLizard.GetTransform()))
	//
	//newLizard.SetMaterial(lizard.GetMaterial())
	//w.AddObject(newLizard)

	//w.AddPointLight(shapes.NewPointLight(geom.NewPoint(sceneSpacing*4, wallDistance*0.8, -wallDistance/2), colors.NewColorFromHex("ffffd7").MulBy(2)))
	//w.AddPointLight(shapes.NewPointLight(geom.NewPoint(sceneSpacing*2, sceneSpacing*3, sceneSpacing*3), colors.NewColorFromHex("af005f").MulBy(3)))
	w.AddPointLight(shapes.NewPointLight(geom.NewPoint(wallDistance*0.4, wallDistance*0.4, wallDistance*0.4), colors.NewColorFromHex("f06553").MulBy(0.15)))
	w.AddPointLight(shapes.NewPointLight(geom.NewPoint(wallDistance*0.4, wallDistance*0.4, -wallDistance*0.4), colors.NewColorFromHex("7f00ff").MulBy(0.15)))
	w.AddAreaLight(shapes.NewAreaLight(geom.NewPoint(0, wallDistance*0.8, 0), geom.NewVector(wallDistance/2, 0, 0), 9, geom.NewVector(0, 0, wallDistance/2), 9, colors.NewColorFromHex("00afaf").MulBy(6), nil))
	w.AddAreaLight(shapes.NewAreaLight(geom.NewPoint(0, 0, wallDistance*0.8), geom.NewVector(wallDistance/2, 0, 0), 9, geom.NewVector(0, wallDistance/2, 0), 9, colors.NewColorFromHex("af005f").MulBy(3), nil))

	//w.AddPointLight(shapes.NewPointLight(geom.NewPoint(cameraDistance/2, sceneSpacing * 6, 0), colors.NewColorFromHex("af005f").MulBy(2)))

	cLookingAt := geom.NewPoint(0, c.Y, c.Z).Add(geom.NewPoint(0, 1.25, 0))

	w.Divide(8)
	return w, basicRotatedCameras(cLookingAt, cameraDistance)
}

func basicRotatedCameras(lookingAt geom.Tuple, distance float64) []CameraLocation {
	views := 6
	cs := make([]CameraLocation, views)

	cs[0].At = geom.NewPoint(lookingAt.X+distance, lookingAt.Y, lookingAt.Z)
	cs[0].LookingAt = lookingAt
	cs[1].At = geom.NewPoint(lookingAt.X-distance, lookingAt.Y, lookingAt.Z)
	cs[1].LookingAt = lookingAt

	cs[2].At = geom.NewPoint(lookingAt.X, lookingAt.Y+distance, lookingAt.Z)
	cs[2].LookingAt = lookingAt
	cs[3].At = geom.NewPoint(lookingAt.X, lookingAt.Y-distance, lookingAt.Z)
	cs[3].LookingAt = lookingAt

	cs[4].At = geom.NewPoint(lookingAt.X, lookingAt.Y, lookingAt.Z+distance)
	cs[4].LookingAt = lookingAt
	cs[5].At = geom.NewPoint(lookingAt.X, lookingAt.Y, lookingAt.Z-distance)
	cs[5].LookingAt = lookingAt

	return cs
}
