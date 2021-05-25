package scenes

import "go-raytrace/lib/view"

type NewSceneFunc func(width int) (view.World, view.Camera)
