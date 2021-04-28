package patterns

import (
	"go-raytrace/lib/colors"
	"go-raytrace/lib/geom"
	"math"
)

type RingPattern struct {
	basePattern
	a Pattern
	b Pattern
}

func NewRingPattern(a, b Pattern) *RingPattern {
	return &RingPattern{
		basePattern: newBasePattern(),
		a:           a,
		b:           b,
	}
}

func (p *RingPattern) ColorAt(t geom.Tuple) colors.Color {
	if int(math.Floor(math.Sqrt(t.X*t.X+t.Z*t.Z)))%2 == 0 {
		return p.a.ColorAtShape(p.t, t)
	}
	return p.b.ColorAtShape(p.t, t)
}

func (p *RingPattern) ColorAtShape(st geom.X4Matrix, t geom.Tuple) colors.Color {
	return ColorAtShape(p, st, t)
}
