package main

import (
	"fmt"
	"image"
	gocolor "image/color"
	"strings"
)

type canvas struct {
	pixels []color
	width  int
	height int
}

const (
	ppmFileHeader    = "P3"
	ppmMaxWidth      = 70
	ppmMinColorValue = 0
	ppmMaxColorValue = 255
)

func newCanvas(width, height int) canvas {
	return canvas{
		pixels: make([]color, width*height),
		width:  width,
		height: height,
	}
}

func newCanvasWith(width int, height int, p color) canvas {
	c := canvas{
		pixels: make([]color, width*height),
		width:  width,
		height: height,
	}
	for i := range c.pixels {
		c.pixels[i] = p
	}
	return c
}

func (c canvas) getPixel(x int, y int) color {
	return c.pixels[y*c.width+x]
}

func (c canvas) setPixel(x int, y int, col color) {
	c.pixels[y*c.width+x] = col
}

func (c canvas) toImage() image.Image {
	upLeft := image.Point{0, 0}
	lowRight := image.Point{c.width, c.height}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	for x := 0; x < c.width; x++ {
		for y := 0; y < c.height; y++ {
			p := c.getPixel(x, y)
			img.Set(x, y, gocolor.RGBA{
				uint8(clamp(p.r*float64(ppmMaxColorValue), ppmMinColorValue, ppmMaxColorValue)),
				uint8(clamp(p.g*float64(ppmMaxColorValue), ppmMinColorValue, ppmMaxColorValue)),
				uint8(clamp(p.b*float64(ppmMaxColorValue), ppmMinColorValue, ppmMaxColorValue)),
				0xff})
		}
	}
	return img
}

func (c canvas) toPPM() string {

	b := strings.Builder{}
	b.WriteString(fmt.Sprintf("%s\n", ppmFileHeader))
	b.WriteString(fmt.Sprintf("%d %d\n", c.width, c.height))
	b.WriteString(fmt.Sprintf("%d\n", ppmMaxColorValue))

	lineLength := 0
	for i, p := range c.pixels {

		rc := fmt.Sprintf("%d", clamp(p.r*float64(ppmMaxColorValue), ppmMinColorValue, ppmMaxColorValue))
		gc := fmt.Sprintf("%d", clamp(p.g*float64(ppmMaxColorValue), ppmMinColorValue, ppmMaxColorValue))
		bc := fmt.Sprintf("%d", clamp(p.b*float64(ppmMaxColorValue), ppmMinColorValue, ppmMaxColorValue))
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
