package scenes

import (
	"github.com/robkau/go-raytrace/lib/colors"
	"github.com/robkau/go-raytrace/lib/geom"
	"github.com/robkau/go-raytrace/lib/shapes"
	"github.com/robkau/go-raytrace/lib/view"
)

func NewGroupGridScene() *Scene {
	w := view.NewWorld()
	cameraPos := geom.NewPoint(0.01, 3, 0.01)
	cameraLookingAt := geom.NewPoint(10, 0, 0)

	g := shapes.NewGroup()
	g.SetTransform(geom.Scale(2, 2, 2))
	s := shapes.NewCube()
	s.SetTransform(geom.Translate(5, 0, 0))
	m := s.GetMaterial()
	m.Color = colors.Blue()
	s.SetMaterial(m)
	g.AddChild(s)

	floor := shapes.NewPlane()
	floor.SetTransform(geom.Translate(0, -2, 0))

	// light above
	w.AddLight(shapes.NewPointLight(geom.NewPoint(5, 10, -3), colors.NewColor(1.9, 1.4, 1.4)))
	w.AddObject(g)
	w.AddObject(floor)

	w.Divide(8)
	return NewScene(w, CameraLocation{cameraPos, cameraLookingAt})
}
