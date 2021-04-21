package patterns

import (
	"go-raytrace/lib/colors"
	"go-raytrace/lib/geom"
)

type BlendPattern struct {
	basePattern
	a Pattern
	b Pattern
}

// todo unit test me

func NewBlendPattern(a, b Pattern) *BlendPattern {
	return &BlendPattern{
		basePattern: newBasePattern(),
		a:           a,
		b:           b,
	}
}

func NewBlendPatternC(a, b Pattern) Pattern {
	return NewBlendPattern(a, b)
}

func (p *BlendPattern) ColorAt(t geom.Tuple) colors.Color {
	added := p.a.ColorAt(t).Add(p.b.ColorAt(t))
	return added.MulBy(0.5)
}

func (p *BlendPattern) ColorAtShape(st geom.X4Matrix, t geom.Tuple) colors.Color {
	return ColorAtShape(p, st, t)
}
