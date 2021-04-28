package patterns

import (
	"go-raytrace/lib/colors"
	"go-raytrace/lib/geom"
	"math/rand"
)

type PerlinPattern struct {
	basePattern
	p Pattern
}

func NewPerlinPattern(p Pattern) *PerlinPattern {
	return &PerlinPattern{
		basePattern: newBasePattern(),
		p:           p,
	}
}

func (p *PerlinPattern) ColorAt(t geom.Tuple) colors.Color {
	// todo: get perlin displacement and apply it to point
	scaling := 0.2
	t.X += rand.Float64() * scaling
	t.Y += rand.Float64() * scaling
	t.Z += rand.Float64() * scaling
	return p.p.ColorAtShape(p.t, t)
}

func (p *PerlinPattern) ColorAtShape(st geom.X4Matrix, t geom.Tuple) colors.Color {
	return ColorAtShape(p, st, t)
}
