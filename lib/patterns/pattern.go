package patterns

import (
	"go-raytrace/lib/colors"
	"go-raytrace/lib/geom"
	"go-raytrace/lib/shapes"
)

type Pattern interface {
	ColorAt(geom.Tuple) colors.Color
	ColorAtShape(s shapes.Shape, t geom.Tuple) colors.Color
	SetTransform(t geom.X4Matrix)
	GetTransform() geom.X4Matrix
}

// invert ray from object's transformation matrix then call pattern-specific logic
func ColorAtShape(p Pattern, wtof func(geom.Tuple) geom.Tuple, worldPoint geom.Tuple) colors.Color {
	objectPoint := wtof(worldPoint)
	patternPoint := p.GetTransform().Invert().MulTuple(objectPoint)
	return p.ColorAt(patternPoint)
}
