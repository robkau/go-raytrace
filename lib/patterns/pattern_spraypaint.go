package patterns

import (
	"go-raytrace/lib/colors"
	"go-raytrace/lib/geom"
	"math/rand"
)

type SprayPaintPattern struct {
	basePattern
	p Pattern
	w float64
}

func NewSprayPaintPattern(p Pattern, width float64) *SprayPaintPattern {
	return &SprayPaintPattern{
		basePattern: newBasePattern(),
		p:           p,
		w:           width,
	}
}

func (p *SprayPaintPattern) ColorAt(t geom.Tuple) colors.Color {
	t.X += rand.Float64() * p.w
	t.Y += rand.Float64() * p.w
	t.Z += rand.Float64() * p.w
	return p.p.ColorAtShape(p.t, t)
}

func (p *SprayPaintPattern) ColorAtShape(st geom.X4Matrix, t geom.Tuple) colors.Color {
	return ColorAtShape(p, st, t)
}