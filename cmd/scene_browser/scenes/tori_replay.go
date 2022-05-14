package scenes

import (
	"github.com/robkau/go-raytrace/lib/colors"
	"github.com/robkau/go-raytrace/lib/geom"
	"github.com/robkau/go-raytrace/lib/parse"
	"github.com/robkau/go-raytrace/lib/shapes"
	"github.com/robkau/go-raytrace/lib/view"
	"strings"
)

func NewToriReplayScene() *Scene {
	w := view.NewWorld()

	g, err := parse.ParseReader(strings.NewReader(parse.ReplayFile), parse.Tori)
	if err != nil {
		panic("err parse")
	}

	for _, c := range g.GetChildren() {
		w.AddObject(c)
	}

	// transparent plane
	waterSurface := shapes.NewPlane()
	m := waterSurface.GetMaterial()
	m.Color = colors.NewColor(0, 0.5, 1)
	m.Diffuse = 0.1
	m.Ambient = 0.1
	m.Specular = 0.5
	m.Shininess = 300
	m.Transparency = 0.5
	m.Reflective = 0.3
	m.RefractiveIndex = 1.13333
	waterSurface.SetTransform(geom.Translate(0, -1, 0))
	waterSurface.SetMaterial(m)
	waterSurface.SetShadowless(true)
	waterSurface.SetShaded(false)

	// light above plane
	w.AddLight(shapes.NewPointLight(geom.NewPoint(2, 12, -5), colors.NewColor(1.9, 1.4, 1.4)))
	w.AddObject(waterSurface)

	return NewScene(w, basicRotatedCameras(geom.NewPoint(0, 0, 0), 10)...)
}

func basicRotatedCameras(lookingAt geom.Tuple, distance float64) []CameraLocation {
	views := 3
	cs := make([]CameraLocation, views)

	cs[0].At = geom.NewPoint(lookingAt.X+distance, lookingAt.Y, lookingAt.Z+distance)
	cs[0].LookingAt = lookingAt

	cs[1].At = geom.NewPoint(lookingAt.X-distance, lookingAt.Y+distance, lookingAt.Z-distance) // todo not seeing anything.
	cs[1].LookingAt = lookingAt

	cs[2].At = geom.NewPoint(lookingAt.X-distance, lookingAt.Y+distance/2, lookingAt.Z+distance)
	cs[2].LookingAt = lookingAt

	return cs
}
