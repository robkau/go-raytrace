package patterns

import (
	"github.com/robkau/go-raytrace/lib/colors"
	"github.com/robkau/go-raytrace/lib/geom"
	"github.com/robkau/go-raytrace/lib/util"
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
	t.X += util.RandomSource.Float64() * p.w
	t.Y += util.RandomSource.Float64() * p.w
	t.Z += util.RandomSource.Float64() * p.w
	return p.p.ColorAtShape(p.worldPointToObjectPoint, t)
}

func (p *SprayPaintPattern) ColorAtShape(wtof WorldToObjectF, t geom.Tuple) colors.Color {
	return ColorAtShape(p, wtof, t)
}
