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

	positionLine := `POS 0; 1.00000000 0.35000002 2.59000014 1.00000000 0.39999998 2.14000010 1.00000000 0.39999998 1.89000010 1.00000000 0.44999999 1.69000005 1.00000000 0.50000000 1.49000000 0.75000000 0.39999998 2.09000014 0.44999999 0.39999998 2.24000000 0.05000000 0.39999998 2.24000000 1.25000000 0.39999998 2.09000014 1.54999995 0.39999998 2.24000000 1.95000005 0.39999998 2.24000000 -0.34999999 0.34999996 2.24000000 2.34999990 0.34999996 2.24000000 0.80000001 0.50000000 1.39000010 1.20000005 0.50000000 1.39000010 0.80000001 0.50000000 1.04000007 1.20000005 0.50000000 1.04000007 1.20000005 0.50000000 0.43999999 0.80000001 0.50000000 0.43999999 0.80000001 0.39999998 0.04000000 1.20000005 0.39999998 0.04000000`

	g, err := parse.ParseReader(strings.NewReader(positionLine), parse.Tori)
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

	cs[0].At = geom.NewPoint(lookingAt.X+distance, lookingAt.Y, lookingAt.Z)
	cs[0].LookingAt = lookingAt

	cs[1].At = geom.NewPoint(lookingAt.X+0.001, lookingAt.Y+distance, lookingAt.Z) // todo not seeing anything.
	cs[1].LookingAt = lookingAt

	cs[2].At = geom.NewPoint(lookingAt.X, lookingAt.Y, lookingAt.Z+distance)
	cs[2].LookingAt = lookingAt

	return cs
}
