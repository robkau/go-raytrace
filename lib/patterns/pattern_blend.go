package patterns

import (
	"go-raytrace/lib/colors"
	"go-raytrace/lib/geom"
)

type BlendPattern struct {
	basePattern
	a         Pattern
	b         Pattern
	intensity float64
}

func NewBlendPattern(a, b Pattern, i float64) *BlendPattern {
	return &BlendPattern{
		basePattern: newBasePattern(),
		a:           a,
		b:           b,
		intensity:   i,
	}
}

func (p *BlendPattern) ColorAt(t geom.Tuple) colors.Color {
	added := p.a.ColorAtShape(p.t, t).Add(p.b.ColorAtShape(p.t, t))
	return added.MulBy(p.intensity)
}

func (p *BlendPattern) ColorAtShape(st geom.X4Matrix, t geom.Tuple) colors.Color {
	return ColorAtShape(p, st, t)
}
