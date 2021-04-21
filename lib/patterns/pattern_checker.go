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

func NewCheckerPatternC(a, b Pattern) Pattern {
	return NewCheckerPattern(a, b)
}

func (p *CheckerPattern) ColorAt(t geom.Tuple) colors.Color {
	if int(math.Abs(t.X)+math.Abs(t.Y)+math.Abs(t.Z))%2 == 0 {
		return p.a.ColorAt(t)
	}
	return p.b.ColorAt(t)
}

func (p *CheckerPattern) ColorAtShape(st geom.X4Matrix, t geom.Tuple) colors.Color {
	return ColorAtShape(p, st, t)
}
