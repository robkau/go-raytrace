package patterns

import (
	"github.com/robkau/go-raytrace/lib/colors"
	"github.com/robkau/go-raytrace/lib/geom"
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
		return p.a.ColorAtShape(p.worldPointToObjectPoint, t)
	}
	return p.b.ColorAtShape(p.worldPointToObjectPoint, t)
}

func (p *RingPattern) ColorAtShape(wtof WorldToObjectF, t geom.Tuple) colors.Color {
	return ColorAtShape(p, wtof, t)
}
