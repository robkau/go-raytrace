package patterns

import (
	"go-raytrace/lib/colors"
	"go-raytrace/lib/geom"
)

type SolidColorPattern struct {
	basePattern
	c colors.Color
}

func NewSolidColorPattern(c colors.Color) *SolidColorPattern {
	return &SolidColorPattern{
		basePattern: newBasePattern(),
		c:           c,
	}
}

func (p *SolidColorPattern) ColorAt(t geom.Tuple) colors.Color {
	return p.c
}

func (p *SolidColorPattern) ColorAtShape(st geom.X4Matrix, t geom.Tuple) colors.Color {
	return ColorAtShape(p, st, t)
}
