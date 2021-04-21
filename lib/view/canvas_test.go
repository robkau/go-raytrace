package view

import (
	"github.com/stretchr/testify/assert"
	"go-raytrace/lib/colors"
	"strings"
	"testing"
)

func Test_NewCanvas(t *testing.T) {
	c := NewCanvas(10, 20)

	for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {
			assert.Equal(t, colors.NewColor(0, 0, 0), c.getPixel(x, y))
		}
	}
}

func Test_SetPixel(t *testing.T) {
	c := NewCanvas(10, 20)

	c.SetPixel(3, 2, colors.NewColor(255, 0, 0))

	assert.Equal(t, colors.NewColor(255, 0, 0), c.getPixel(3, 2))
}

func Test_ToPPM_Header(t *testing.T) {
	c := NewCanvas(5, 3)

	s := c.toPPM()
	lines := strings.Split(s, "\n")

	assert.Equal(t, ppmFileHeader, lines[0])
	assert.Equal(t, "5 3", lines[1])
	assert.Equal(t, "255", lines[2])
}

func Test_ToPPM_Content(t *testing.T) {
	c := NewCanvas(5, 3)
	c.SetPixel(0, 0, colors.NewColor(1.5, 0, 0))
	c.SetPixel(2, 1, colors.NewColor(0, 0.5, 0))
	c.SetPixel(4, 2, colors.NewColor(-0.5, 0, 1))

	s := c.toPPM()
	lines := strings.Split(s, "\n")

	assert.Equal(t, "255 0 0 0 0 0 0 0 0 0 0 0 0 0 0", lines[3])
	assert.Equal(t, "0 0 0 0 0 0 0 128 0 0 0 0 0 0 0", lines[4])
	assert.Equal(t, "0 0 0 0 0 0 0 0 0 0 0 0 0 0 255", lines[5])
}

func Test_ToPPM_Content_LongLine(t *testing.T) {
	c := newCanvasWith(10, 2, colors.NewColor(1.0, 0.8, 0.6))

	s := c.toPPM()
	lines := strings.Split(s, "\n")

	assert.Equal(t, "255 204 153 255 204 153 255 204 153 255 204 153 255 204 153 255 204", lines[3])
	assert.Equal(t, "153 255 204 153 255 204 153 255 204 153 255 204 153", lines[4])
	assert.Equal(t, "255 204 153 255 204 153 255 204 153 255 204 153 255 204 153 255 204", lines[5])
	assert.Equal(t, "153 255 204 153 255 204 153 255 204 153 255 204 153", lines[6])
}

func Test_ToPPM_EndsWithNewLine(t *testing.T) {
	c := NewCanvas(5, 3)

	s := c.toPPM()
	lastChar := string(s[len(s)-1])
	assert.Equal(t, "\n", lastChar)
}
