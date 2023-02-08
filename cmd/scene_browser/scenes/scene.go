package scenes

import (
	"github.com/robkau/go-raytrace/lib/geom"
	"github.com/robkau/go-raytrace/lib/view"
	"math"
	"sync/atomic"
)

type Scene struct {
	W      *view.World
	Cs     []CameraLocation
	loadF  NewSceneFunc
	loaded *atomic.Bool
}

func NewScene(loadF NewSceneFunc) *Scene {
	return &Scene{
		W:      nil,
		Cs:     nil,
		loadF:  loadF,
		loaded: &atomic.Bool{}, // todo zerovalue ok? not reference ok?
	}
}

func (s *Scene) Load() {
	if s.loaded.CompareAndSwap(false, true) {
		s.W, s.Cs = s.loadF()
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

type NewSceneFunc func() (*view.World, []CameraLocation)

func LoadScenes(fs ...NewSceneFunc) []*Scene {
	scenes := []*Scene{}
	for _, initF := range fs {
		scenes = append(scenes, NewScene(initF))
	}
	return scenes
}
