package patterns

import (
	"go-raytrace/lib/colors"
	"go-raytrace/lib/geom"
)

type PositionAsColorPattern struct {
	basePattern
}

func NewPositionAsColorPattern() *PositionAsColorPattern {
	return &PositionAsColorPattern{
		basePattern: newBasePattern(),
	}
}

func (p *PositionAsColorPattern) ColorAt(t geom.Tuple) colors.Color {
	return colors.NewColor(
		t.X,
		t.Y,
		t.Z,
	)
}

func (p *PositionAsColorPattern) ColorAtShape(st geom.X4Matrix, t geom.Tuple) colors.Color {
	return ColorAtShape(p, st, t)
}
