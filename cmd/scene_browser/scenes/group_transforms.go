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
	"math/rand"
	"strings"
)

func makeObjectGroup() shapes.Group {
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
	ttm.Pattern.SetTransform(geom.Scale(0.06, 0.06, 0.06))
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

	// hexagon floating above table
	h := shapes.NewHexagon()
	h.SetShadowless(true)
	h.SetTransform(geom.Translate(5, 8.6, 1.5).MulX4Matrix(geom.Scale(1.1, 0.8, 1.1)).MulX4Matrix(geom.RotateX(-math.Pi / 3)).MulX4Matrix(geom.RotateZ(math.Pi / 6)).MulX4Matrix(geom.RotateY(-math.Pi / 6)))
	m = materials.NewMaterial()
	m.Pattern = patterns.NewSolidColorPattern(colors.NewColor(0.7, 0.7, 0))
	m.Specular = 0.05
	m.Ambient = 0.05
	m.Diffuse = 0.05
	m.Transparency = 0.97
	m.Shininess = 0
	h.SetMaterial(m)
	g.AddChild(h)

	// a random frame from toribash replay on each table
	pr, err := parse.ParseReaderAsTori(strings.NewReader(parse.ReplayFile))
	if err != nil {
		log.Fatalf("err parse")
	}
	n := rand.Intn(int(math.Min(float64(len(pr.P0Positions)), float64(len(pr.P1Positions)))))
	pg0 := pr.P0Positions[n].AsGroup()
	pg0.SetTransform(geom.Translate(0, 4.1+parse.ToriSphereWidth, 0).MulX4Matrix(geom.Scale(1.2, 1.2, 1.2)))
	m = materials.NewMaterial()
	m.Pattern = patterns.NewSolidColorPattern(colors.Green())
	for _, c := range pg0.GetChildren() {
		c.SetMaterial(m)
	}
	g.AddChild(pg0)

	pg1 := pr.P1Positions[n].AsGroup()
	pg1.SetTransform(geom.Translate(0, 4.1+parse.ToriSphereWidth, 0).MulX4Matrix(geom.Scale(1.2, 1.2, 1.2)))
	m = materials.NewMaterial()
	m.Pattern = patterns.NewSolidColorPattern(colors.Red())
	for _, c := range pg1.GetChildren() {
		c.SetMaterial(m)
	}
	g.AddChild(pg1)

	g.Divide(8)
	return g
}

func makeGroupOfGroups() shapes.Group {
	g1 := makeObjectGroup()

	g2 := makeObjectGroup()
	g2.SetTransform(stackTable(geom.NewIdentityMatrixX4()))

	g3 := makeObjectGroup()
	g3.SetTransform(stackTable(geom.NewIdentityMatrixX4()))

	g4 := makeObjectGroup()
	g4.SetTransform(stackTable(geom.NewIdentityMatrixX4()))

	g5 := makeObjectGroup()
	g5.SetTransform(stackTable(geom.NewIdentityMatrixX4()))

	g6 := makeObjectGroup()
	g6.SetTransform(stackTable(geom.NewIdentityMatrixX4()))

	g7 := makeObjectGroup()
	g7.SetTransform(stackTable(geom.NewIdentityMatrixX4()))

	g1.AddChild(g2)
	g2.AddChild(g3)
	g3.AddChild(g4)
	g4.AddChild(g5)
	g5.AddChild(g6)
	g6.AddChild(g7)

	return g1
}

func stackTable(t geom.X4Matrix) geom.X4Matrix {
	return t.MulX4Matrix(geom.Translate(-1, 4, 1)).MulX4Matrix(geom.RotateY(math.Pi / 9).MulX4Matrix(geom.Scale(0.5, 0.5, 0.5)))
}

func NewGroupTransformsScene() *Scene {
	//seed := time.Now().UnixNano()
	//fmt.Println("rand seed", seed)
	//rand.Seed(seed)

	w := view.NewWorld()
	cameraPos := geom.NewPoint(15, 15, 15)
	cameraLookingAt := geom.NewPoint(0, 5, 0)

	g1 := makeGroupOfGroups()
	g1.SetTransform(geom.Translate(-7, 0, -3))

	g2 := makeGroupOfGroups()
	g2.SetTransform(geom.Translate(0.5, 0, -3).MulX4Matrix(geom.RotateY(math.Pi)))

	// floor and ceiling as one cube
	var floorAndCeiling = sizedCubeAt(0, 10, 0, 100, 10, 100)
	m := floorAndCeiling.GetMaterial()
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

	// light above
	w.AddAreaLight(shapes.NewAreaLight(geom.NewPoint(0, 10, 0), geom.NewVector(15, 0, 0), 4, geom.NewVector(0, 0, 15), 4, colors.NewColorFromHex("ffffd7").MulBy(2), nil))
	//w.AddPointLight(shapes.NewPointLight(geom.NewPoint(5, 10, -3), colors.NewColor(1.9, 1.4, 1.4)))
	w.AddObject(g1)
	w.AddObject(g2)
	w.AddObject(floorAndCeiling)
	w.AddObject(walls)

	return NewScene(w, CameraLocation{cameraPos, cameraLookingAt})
}
