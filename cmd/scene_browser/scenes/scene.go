package scenes

import (
	"github.com/robkau/go-raytrace/lib/geom"
	"github.com/robkau/go-raytrace/lib/view"
	"math"
)

type NewSceneFunc func() *Scene

type Scene struct {
	W  *view.World
	Cs []CameraLocation
}

func NewScene(w *view.World, cs ...CameraLocation) *Scene {
	return &Scene{
		W:  w,
		Cs: cs,
	}
}

type CameraLocation struct {
	At        geom.Tuple
	LookingAt geom.Tuple
}

func (l *CameraLocation) RotateAroundX(radians float64) {
	diff := l.At.Sub(l.LookingAt)

	yNew := diff.Y*math.Cos(radians) - diff.Z*math.Sin(radians)
	zNew := diff.Y*math.Sin(radians) + diff.Z*math.Cos(radians)

	l.At.Y = yNew
	l.At.Z = zNew
}

func (l *CameraLocation) RotateAroundY(radians float64) {
	diff := l.At.Sub(l.LookingAt)

	xNew := diff.X*math.Cos(radians) - diff.Z*math.Sin(radians)
	zNew := diff.X*math.Sin(radians) + diff.Z*math.Cos(radians)

	l.At.X = xNew
	l.At.Z = zNew
}

func (l *CameraLocation) RotateAroundZ(radians float64) {
	diff := l.At.Sub(l.LookingAt)

	xNew := diff.X*math.Cos(radians) - diff.Y*math.Sin(radians)
	yNew := diff.X*math.Sin(radians) + diff.Y*math.Cos(radians)

	l.At.X = xNew
	l.At.Y = yNew
}

func LoadScenes(fs ...NewSceneFunc) []*Scene {
	scenes := []*Scene{}
	for _, initF := range fs {
		scenes = append(scenes, initF())
	}
	return scenes
}
