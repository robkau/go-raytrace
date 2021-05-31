package view

import (
	"fmt"
	"go-raytrace/lib/colors"
	"image"
	gocolor "image/color"
	"strings"
	"sync"
)

type Canvas struct {
	pixels []colors.Color
	width  int
	height int
	rw     sync.RWMutex
}

const (
	ppmFileHeader    = "P3"
	ppmMaxWidth      = 70
	ppmMinColorValue = 0
	ppmMaxColorValue = 255
)

func NewCanvas(width, height int) *Canvas {
	return &Canvas{
		pixels: make([]colors.Color, width*height),
		width:  width,
		height: height,
	}
}

func newCanvasWith(width int, height int, p colors.Color) *Canvas {
	c := &Canvas{
		pixels: make([]colors.Color, width*height),
		width:  width,
		height: height,
	}
	for i := range c.pixels {
		c.pixels[i] = p
	}
	return c
}

func (c *Canvas) getPixel(x int, y int) colors.Color {
	c.rw.RLock()
	defer c.rw.RUnlock()
	return c.pixels[y*c.width+x]
}

func (c *Canvas) SetPixel(x int, y int, col colors.Color) {
	c.rw.Lock()
	defer c.rw.Unlock()
	c.pixels[y*c.width+x] = col
}

func (c *Canvas) ToImage() image.Image {
	upLeft := image.Point{0, 0}
	lowRight := image.Point{c.width, c.height}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	for x := 0; x < c.width; x++ {
		for y := 0; y < c.height; y++ {
			p := c.getPixel(x, y)
			img.Set(x, y, gocolor.RGBA{
				R: uint8(clamp(p.R*float64(ppmMaxColorValue), ppmMinColorValue, ppmMaxColorValue)),
				G: uint8(clamp(p.G*float64(ppmMaxColorValue), ppmMinColorValue, ppmMaxColorValue)),
				B: uint8(clamp(p.B*float64(ppmMaxColorValue), ppmMinColorValue, ppmMaxColorValue)),
				A: 0xff})
		}
	}
	return img
}

func (c *Canvas) toPPM() string {

	b := strings.Builder{}
	b.WriteString(fmt.Sprintf("%s\n", ppmFileHeader))
	b.WriteString(fmt.Sprintf("%d %d\n", c.width, c.height))
	b.WriteString(fmt.Sprintf("%d\n", ppmMaxColorValue))

	lineLength := 0
	for i, p := range c.pixels {

		rc := fmt.Sprintf("%d", clamp(p.R*float64(ppmMaxColorValue), ppmMinColorValue, ppmMaxColorValue))
		gc := fmt.Sprintf("%d", clamp(p.G*float64(ppmMaxColorValue), ppmMinColorValue, ppmMaxColorValue))
		bc := fmt.Sprintf("%d", clamp(p.B*float64(ppmMaxColorValue), ppmMinColorValue, ppmMaxColorValue))
		rcl := len(rc)
		gcl := len(gc)
		bcl := len(bc)

		if i%c.width == 0 && i != 0 {
			// finished writing a row
			b.WriteString("\n")
			lineLength = 0
		}

		// red
		if lineLength == 0 {
			// special case: first character of a line
			b.WriteString(rc)
			lineLength += rcl
		} else {
			if lineLength+rcl+1 > ppmMaxWidth {
				// not enough space in line. break it
				b.WriteString("\n")
				lineLength = 0
			} else {
				// continue line
				b.WriteString(" ")
				lineLength += 1
			}
			b.WriteString(rc)
			lineLength += rcl
		}

		// green
		if lineLength+1+gcl > ppmMaxWidth {
			// not enough space in line. break it
			b.WriteString("\n")
			lineLength = 0
		} else {
			// continue line
			b.WriteString(" ")
			lineLength += 1
		}
		b.WriteString(gc)
		lineLength += gcl

		// blue
		if lineLength+1+bcl > ppmMaxWidth {
			// not enough space in line. break it
			b.WriteString("\n")
			lineLength = 0
		} else {
			// continue line
			b.WriteString(" ")
			lineLength += 1
		}
		b.WriteString(bc)
		lineLength += bcl

	}
	b.WriteString("\n")
	return b.String()
}

func clamp(f float64, min int, max int) int {
	r := int(f + 0.5) // round it
	if r <= min {
		r = min
	}

	if r >= max {
		r = max
	}

	return r
}
