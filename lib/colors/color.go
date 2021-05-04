package colors

import (
	"go-raytrace/lib/geom"
	"math"
	"math/rand"
)

type Color struct {
	R float64
	G float64
	B float64
}

var (
	presetColors = []Color{Black(), Blue(), Green(), Red(), White()}
)

func NewColor(r, g, b float64) Color {
	return Color{R: r, G: g, B: b}
}

func Black() Color {
	return Color{R: 0, G: 0, B: 0}
}

func Blue() Color {
	return Color{R: 0, G: 0, B: 1}
}

func Brown() Color {
	return Color{R: 0.40, G: 0.26, B: 0.13}
}

func Green() Color {
	return Color{R: 0, G: 1, B: 0}
}

func Red() Color {
	return Color{R: 1, G: 0, B: 0}
}

func White() Color {
	return Color{R: 1, G: 1, B: 1}
}

func RandomColor() Color {
	index := rand.Intn(len(presetColors))
	return presetColors[index]
}

func RandomAnyColor() Color {
	return NewColor(rand.Float64(), rand.Float64(), rand.Float64())
}

func (c Color) Add(other Color) Color {
	return Color{
		c.R + other.R,
		c.G + other.G,
		c.B + other.B,
	}
}

func (c Color) Sub(other Color) Color {
	return Color{
		c.R - other.R,
		c.G - other.G,
		c.B - other.B,
	}
}
func (c Color) Mul(other Color) Color {
	return Color{
		c.R * other.R,
		c.G * other.G,
		c.B * other.B,
	}
}

func (c Color) MulBy(x float64) Color {
	return Color{
		c.R * x,
		c.G * x,
		c.B * x,
	}
}

func (c Color) Equal(other Color) bool {
	return geom.AlmostEqual(c.R, other.R) &&
		geom.AlmostEqual(c.G, other.G) &&
		geom.AlmostEqual(c.B, other.B)
}

func (c Color) RoundTo(places int) Color {
	scale := math.Pow10(places)
	return Color{
		R: math.Round(c.R*scale) / scale,
		G: math.Round(c.G*scale) / scale,
		B: math.Round(c.B*scale) / scale,
	}
}
