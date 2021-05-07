package scenes

import (
	"go-raytrace/lib/colors"
	"go-raytrace/lib/geom"
	"go-raytrace/lib/shapes"
	"go-raytrace/lib/view"
)

func sizedCubeAt(x, y, z, w, h, d float64) shapes.Shape {
	var c shapes.Shape = shapes.NewCube()
	c = c.SetTransform(geom.Translate(x, y, z).MulX4Matrix(geom.Scale(w, h, d)))
	m := c.GetMaterial()
	m.Reflective = 0.05
	m.Diffuse = 0.4
	m.Specular = 0.7
	m.Ambient = 0.1
	m.Shininess = 90
	c = c.SetMaterial(m)
	return c
}

func NewRoomScene(width int) (view.World, view.Camera) {
	w := view.NewWorld()
	cameraPos := geom.NewPoint(15, 15, 15)
	cameraLookingAt := geom.NewPoint(0, 5, 0)

	// table made of cubes
	tlul := sizedCubeAt(-3, 2, -3, 0.2, 2, 0.2)
	m := tlul.GetMaterial()
	m.Color = colors.NewColor(0.2, 0.9, 0.2)
	tlul = tlul.SetMaterial(m)
	w.AddObject(tlul)
	tlur := sizedCubeAt(-3, 2, 3, 0.2, 2, 0.2)
	tlur = tlur.SetMaterial(m)
	w.AddObject(tlur)
	tldr := sizedCubeAt(3, 2, 3, 0.2, 2, 0.2)
	tldr = tldr.SetMaterial(m)
	w.AddObject(tldr)
	tldl := sizedCubeAt(3, 2, -3, 0.2, 2, 0.2)
	tldl = tldl.SetMaterial(m)
	w.AddObject(tldl)
	tt := sizedCubeAt(0, 4.1, 0, 3.19, 0.15, 3.19)
	tt = tt.SetMaterial(m)
	w.AddObject(tt)

	// a few cubes on table
	cl := sizedCubeAt(-0.5, 4.44, -1.5, 0.2, 0.2, 0.2)
	m.Color = colors.Blue()
	cl = cl.SetMaterial(m)
	w.AddObject(cl)
	cr := sizedCubeAt(-1, 4.53, 0.75, 0.3, 0.3, 0.3)
	m.Color = colors.Red()
	cr = cr.SetMaterial(m)
	w.AddObject(cr)

	// floor and ceiling as one cube
	var floorAndCeiling = sizedCubeAt(0, 10, 0, 100, 10, 100)
	m = floorAndCeiling.GetMaterial()
	m.Color = colors.Brown()
	m.Reflective = 0
	m.Transparency = 0
	floorAndCeiling = floorAndCeiling.SetMaterial(m)

	// walls as another cube
	var walls = sizedCubeAt(0, 0, 0, 20, 100, 20)
	m = walls.GetMaterial()
	m.Reflective = 0
	m.Transparency = 0
	m.Color = colors.Blue()
	walls = walls.SetMaterial(m)

	// light above
	w.AddLight(shapes.NewPointLight(geom.NewPoint(5, 10, -3), colors.NewColor(1.9, 1.4, 1.4)))
	w.AddObject(floorAndCeiling)
	w.AddObject(walls)

	c := view.NewCamera(width, width, 0.45)
	c.Transform = geom.ViewTransform(cameraPos,
		cameraLookingAt,
		geom.NewVector(0, 1, 0))

	return w, c
}
