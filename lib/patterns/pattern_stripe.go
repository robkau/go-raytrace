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

func (p *StripePattern) ColorAt(t geom.Tuple) colors.Color {
	if int(math.Floor(t.X))%2 == 0 {
		return p.a.ColorAtShape(p.t, t)
	}
	return p.b.ColorAtShape(p.t, t)
}

func (p *StripePattern) ColorAtShape(st geom.X4Matrix, t geom.Tuple) colors.Color {
	return ColorAtShape(p, st, t)
}
