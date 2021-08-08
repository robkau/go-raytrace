package scenes

import "github.com/robkau/go-raytrace/lib/view"

type NewSceneFunc func(width int) (view.World, view.Camera)
