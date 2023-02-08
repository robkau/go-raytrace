package patterns

import (
	"github.com/robkau/go-raytrace/lib/colors"
	"github.com/robkau/go-raytrace/lib/geom"
)

type Pattern interface {
	ColorAt(geom.Tuple) colors.Color
	ColorAtShape(wtof WorldToObjectF, t geom.Tuple) colors.Color
	SetTransform(t *geom.X4Matrix)
	GetTransform() *geom.X4Matrix
}

type WorldToObjectF func(geom.Tuple) geom.Tuple

// invert ray from object's transformation matrix then call pattern-specific logic
func ColorAtShape(p Pattern, wtof WorldToObjectF, worldPoint geom.Tuple) colors.Color {
	objectPoint := wtof(worldPoint)
	patternPoint := p.GetTransform().Invert().MulTuple(objectPoint)
	return p.ColorAt(patternPoint)
}
