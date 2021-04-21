package patterns

import (
	"go-raytrace/lib/colors"
	"go-raytrace/lib/geom"
	"math"
)

type StripePattern struct {
	basePattern
	a Pattern
	b Pattern
}

func NewStripePattern(a, b Pattern) *StripePattern {
	return &StripePattern{
		basePattern: newBasePattern(),
		a:           a,
		b:           b,
	}
}

func NewStripePatternC(a, b Pattern) Pattern {
	return NewStripePattern(a, b)
}

func (p *StripePattern) ColorAt(t geom.Tuple) colors.Color {
	if int(math.Floor(t.X))%2 == 0 {
		return p.a.ColorAt(t)
	}
	return p.b.ColorAt(t)
}

func (p *StripePattern) ColorAtShape(st geom.X4Matrix, t geom.Tuple) colors.Color {
	return ColorAtShape(p, st, t)
}
