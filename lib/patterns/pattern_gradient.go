package patterns

import (
	"go-raytrace/lib/colors"
	"go-raytrace/lib/geom"
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
	distance := p.b.ColorAt(t).Sub(p.a.ColorAtShape(p.t, t))
	fraction := t.X - math.Floor(t.X)
	return p.a.ColorAtShape(p.t, t).Add(distance.MulBy(fraction))
}

func (p *GradientPattern) ColorAtShape(st geom.X4Matrix, t geom.Tuple) colors.Color {
	return ColorAtShape(p, st, t)
}
