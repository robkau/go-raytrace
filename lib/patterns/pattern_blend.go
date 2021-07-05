package patterns

import (
	"github.com/robkau/go-raytrace/lib/colors"
	"github.com/robkau/go-raytrace/lib/geom"
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
	added := p.a.ColorAtShape(p.worldPointToObjectPoint, t).Add(p.b.ColorAtShape(p.worldPointToObjectPoint, t))
	return added.MulBy(p.intensity)
}

func (p *BlendPattern) ColorAtShape(wtof WorldToObjectF, t geom.Tuple) colors.Color {
	return ColorAtShape(p, wtof, t)
}
