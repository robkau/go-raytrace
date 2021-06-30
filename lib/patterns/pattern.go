package patterns

import (
	"go-raytrace/lib/colors"
	"go-raytrace/lib/geom"
)

type Pattern interface {
	ColorAt(geom.Tuple) colors.Color
	ColorAtShape(st geom.X4Matrix, t geom.Tuple) colors.Color
	SetTransform(t geom.X4Matrix)
	GetTransform() geom.X4Matrix
}

// invert ray from object's transformation matrix then call pattern-specific logic
func ColorAtShape(p Pattern, st geom.X4Matrix, worldPoint geom.Tuple) colors.Color {
	objectPoint := st.Invert().MulTuple(worldPoint)
	patternPoint := p.GetTransform().Invert().MulTuple(objectPoint)
	return p.ColorAt(patternPoint)
}
