package patterns

import (
	"go-raytrace/lib/colors"
	"go-raytrace/lib/geom"
	"math"
)

type CheckerPattern struct {
	basePattern
	a Pattern
	b Pattern
}

func NewCheckerPattern(a, b Pattern) *CheckerPattern {
	return &CheckerPattern{
		basePattern: newBasePattern(),
		a:           a,
		b:           b,
	}
}

func (p *CheckerPattern) ColorAt(t geom.Tuple) colors.Color {
	if int(math.Floor(t.X)+math.Floor(t.Y)+math.Floor(t.Z))%2 == 0 {
		return p.a.ColorAtShape(p.t, t)
	}
	return p.b.ColorAtShape(p.t, t)
}

func (p *CheckerPattern) ColorAtShape(st geom.X4Matrix, t geom.Tuple) colors.Color {
	return ColorAtShape(p, st, t)
}
