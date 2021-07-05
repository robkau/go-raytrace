package patterns

import (
	"github.com/robkau/go-raytrace/lib/colors"
	"github.com/robkau/go-raytrace/lib/geom"
	"math"
)

type GradientPattern struct {
	basePattern
	a Pattern
	b Pattern
}

func NewGradientPattern(a, b Pattern) *GradientPattern {
	return &GradientPattern{
		basePattern: newBasePattern(),
		a:           a,
		b:           b,
	}
}

func (p *GradientPattern) ColorAt(t geom.Tuple) colors.Color {
	distance := p.b.ColorAt(t).Sub(p.a.ColorAtShape(p.worldPointToObjectPoint, t))
	fraction := t.X - math.Floor(t.X)
	return p.a.ColorAtShape(p.worldPointToObjectPoint, t).Add(distance.MulBy(fraction))
}

func (p *GradientPattern) ColorAtShape(wtof WorldToObjectF, t geom.Tuple) colors.Color {
	return ColorAtShape(p, wtof, t)
}
