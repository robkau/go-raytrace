package colors

import (
	"fmt"
	"github.com/robkau/go-raytrace/lib/geom"
	"github.com/robkau/go-raytrace/lib/util"
	"image/color"
	"math"
	"math/rand"
	"strings"
)

type Color struct {
	R float64
	G float64
	B float64
}

const maxColorValue = 255
const minColorValue = 255

func (c Color) RGBA() (r, g, b, a uint32) {
	return uint32(util.Clamp(c.R*float64(maxColorValue), minColorValue, maxColorValue)),
		uint32(util.Clamp(c.G*float64(maxColorValue), minColorValue, maxColorValue)),
		uint32(util.Clamp(c.B*float64(maxColorValue), minColorValue, maxColorValue)),
		1
}

var (
	presetColors = []Color{Black(), Blue(), Green(), Red(), White()}
)

func NewColor(r, g, b float64) Color {
	return Color{R: r, G: g, B: b}
}

func NewColorFromStdlibColor(c color.Color) Color {
	r, g, b, a := c.RGBA()
	// stdlib returns in format 0xffff (65535) so convert to 0xff (255)
	return NewColorFromRGBA(uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a>>8))
}

func NewColorFromRGBA(r, g, b, a uint8) Color {
	return NewColor(
		(float64(a)/255*float64(r))/255,
		float64(a)/255*float64(g)/255,
		float64(a)/255*float64(b)/255,
	)
}

func NewColorFromHex(hex string) Color {
	// adapted from https://stackoverflow.com/a/54200713/1723695

	c := color.RGBA{}
	c.A = 0xff

	hex = strings.ToLower(hex)

	hexToByte := func(b byte) byte {
		switch {
		case b >= '0' && b <= '9':
			return b - '0'
		case b >= 'a' && b <= 'f':
			return 10 + b - 'a'
		case b >= 'A' && b <= 'F':
			return 10 + b - 'A'
		}
		return 0
	}

	switch len(hex) {
	case 8:
		c.R = (hexToByte(hex[0]) << 4) + hexToByte(hex[1])
		c.G = (hexToByte(hex[2]) << 4) + hexToByte(hex[3])
		c.B = (hexToByte(hex[4]) << 4) + hexToByte(hex[5])
		c.A = (hexToByte(hex[6]) << 4) + hexToByte(hex[7])
	case 6:
		c.R = (hexToByte(hex[0]) << 4) + hexToByte(hex[1])
		c.G = (hexToByte(hex[2]) << 4) + hexToByte(hex[3])
		c.B = (hexToByte(hex[4]) << 4) + hexToByte(hex[5])
	case 4:
		c.R = hexToByte(hex[0]) * 17
		c.G = hexToByte(hex[1]) * 17
		c.B = hexToByte(hex[2]) * 17
		c.A = hexToByte(hex[3]) * 17
	case 3:
		c.R = hexToByte(hex[0]) * 17
		c.G = hexToByte(hex[1]) * 17
		c.B = hexToByte(hex[2]) * 17
	}

	return NewColorFromRGBA(c.R, c.G, c.B, c.A)
}

func Black() Color {
	return Color{R: 0, G: 0, B: 0}
}

func Blue() Color {
	return Color{R: 0, G: 0, B: 1}
}

func Brown() Color {
	return Color{R: 1, G: 0.5, B: 0}
}

func Cyan() Color {
	return Color{R: 0, G: 1, B: 1}
}

func Purple() Color {
	return Color{R: 1, G: 0, B: 1}
}

func Green() Color {
	return Color{R: 0, G: 1, B: 0}
}

func Red() Color {
	return Color{R: 1, G: 0, B: 0}
}

func Yellow() Color {
	return Color{R: 1, G: 1, B: 0}
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

func (c Color) Mag() float64 {
	return c.R + c.B + c.G
}

func (c Color) Equal(other Color) bool {
	return geom.AlmostEqualWithPrecision(c.R, other.R, 1e-3) &&
		geom.AlmostEqualWithPrecision(c.G, other.G, 1e-3) &&
		geom.AlmostEqualWithPrecision(c.B, other.B, 1e-3)
}

func (c Color) RoundTo(places int) Color {
	scale := math.Pow10(places)
	return Color{
		R: math.Round(c.R*scale) / scale,
		G: math.Round(c.G*scale) / scale,
		B: math.Round(c.B*scale) / scale,
	}
}

func (c Color) String() string {
	return fmt.Sprintf("%f:%f:%f", c.R, c.G, c.B)
}
