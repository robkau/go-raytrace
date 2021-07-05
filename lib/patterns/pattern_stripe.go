package patterns

import (
	"github.com/robkau/go-raytrace/lib/colors"
	"github.com/robkau/go-raytrace/lib/geom"
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
		return p.a.ColorAtShape(p.worldPointToObjectPoint, t)
	}
	return p.b.ColorAtShape(p.worldPointToObjectPoint, t)
}

func (p *StripePattern) ColorAtShape(wtof WorldToObjectF, t geom.Tuple) colors.Color {
	return ColorAtShape(p, wtof, t)
}
