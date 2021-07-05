package scenes

import (
	"github.com/robkau/go-raytrace/lib/colors"
	"github.com/robkau/go-raytrace/lib/geom"
	"github.com/robkau/go-raytrace/lib/patterns"
	"github.com/robkau/go-raytrace/lib/shapes"
	"github.com/robkau/go-raytrace/lib/view"
	"math"
)

func makeObjectGroup() shapes.Shape {
	g := shapes.NewGroup()

	// table made of cubes
	tlul := sizedCubeAt(-3, 2, -3, 0.2, 2, 0.2)
	m := tlul.GetMaterial()
	m.Color = colors.NewColor(0.2, 0.9, 0.2)
	tlul.SetMaterial(m)
	g.AddChild(tlul)
	tlur := sizedCubeAt(-3, 2, 3, 0.2, 2, 0.2)
	tlur.SetMaterial(m)
	g.AddChild(tlur)
	tldr := sizedCubeAt(3, 2, 3, 0.2, 2, 0.2)
	tldr.SetMaterial(m)
	g.AddChild(tldr)
	tldl := sizedCubeAt(3, 2, -3, 0.2, 2, 0.2)
	tldl.SetMaterial(m)
	g.AddChild(tldl)
	tt := sizedCubeAt(0, 4.1, 0, 3.19, 0.15, 3.19)
	ttm := tlul.GetMaterial()
	ttm.Pattern = patterns.NewStripePattern(patterns.NewSolidColorPattern(colors.Black()), patterns.NewSolidColorPattern(colors.NewColor(0.2237, 0.328, 0.44235)))
	ttm.Pattern.SetTransform(geom.Scale(0.02, 0.02, 0.02))
	ttm.Reflective = 0.2
	ttm.Shininess = 0.4
	ttm.Specular = 0.1
	tt.SetMaterial(ttm)
	g.AddChild(tt)

	// a few cubes on table
	cl := sizedCubeAt(1.75, 4.64, 1.5, 0.4, 0.4, 0.4)
	m.Color = colors.Blue()
	cl.SetMaterial(m)
	g.AddChild(cl)
	cr := sizedCubeAt(1, 4.83, -0.75, 0.6, 0.6, 0.6)
	m.Color = colors.Red()
	cr.SetMaterial(m)
	g.AddChild(cr)

	return g
}

func stackTable(t geom.X4Matrix) geom.X4Matrix {
	return t.MulX4Matrix(geom.Translate(-1, 4, 1)).MulX4Matrix(geom.RotateY(math.Pi / 9).MulX4Matrix(geom.Scale(0.5, 0.5, 0.5)))
}

func NewGroupTransformsScene(width int) (view.World, view.Camera) {
	w := view.NewWorld()
	cameraPos := geom.NewPoint(15, 15, 15)
	cameraLookingAt := geom.NewPoint(0, 5, 0)

	g1 := makeObjectGroup()

	g2 := makeObjectGroup()
	g2.SetTransform(stackTable(geom.NewIdentityMatrixX4()))

	g3 := makeObjectGroup()
	g3.SetTransform(stackTable(g2.GetTransform()))

	g4 := makeObjectGroup()
	g4.SetTransform(stackTable(g3.GetTransform()))

	g5 := makeObjectGroup()
	g5.SetTransform(stackTable(g4.GetTransform()))

	g6 := makeObjectGroup()
	g6.SetTransform(stackTable(g5.GetTransform()))

	g7 := makeObjectGroup()
	g7.SetTransform(stackTable(g6.GetTransform()))

	// floor and ceiling as one cube
	var floorAndCeiling = sizedCubeAt(0, 10, 0, 100, 10, 100)
	m := floorAndCeiling.GetMaterial()
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
	w.AddObject(g1)
	w.AddObject(g2)
	w.AddObject(g3)
	w.AddObject(g4)
	w.AddObject(g5)
	w.AddObject(g6)
	w.AddObject(g7)
	w.AddObject(floorAndCeiling)
	w.AddObject(walls)

	c := view.NewCamera(width, width, 0.45)
	c.Transform = geom.ViewTransform(cameraPos,
		cameraLookingAt,
		geom.UpVector())

	return w, c
}
