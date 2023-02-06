package canvas

import (
	"github.com/robkau/go-raytrace/lib/colors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strconv"
	"strings"
	"testing"
)

func Test_NewCanvas(t *testing.T) {
	c := NewCanvas(10, 20)

	for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {
			assert.Equal(t, colors.NewColor(0, 0, 0), c.GetPixel(x, y))
		}
	}
}

func Test_SetPixel(t *testing.T) {
	c := NewCanvas(10, 20)

	c.SetPixel(3, 2, colors.NewColor(255, 0, 0))

	assert.Equal(t, colors.NewColor(255, 0, 0), c.GetPixel(3, 2))
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

func Test_CanvasFromPPM_InvalidHeader(t *testing.T) {
	ppmFile := `P32
1 1
255
0 0 0`

	_, err := CanvasFromPPMReader(strings.NewReader(ppmFile))
	require.Error(t, err)
}

func Test_CanvasFromPPM_HasCorrectDimensions(t *testing.T) {
	ppmFile := `P3
10 2
255
0 0 0  0 0 0  0 0 0  0 0 0  0 0 0
0 0 0  0 0 0  0 0 0  0 0 0  0 0 0
0 0 0  0 0 0  0 0 0  0 0 0  0 0 0
0 0 0  0 0 0  0 0 0  0 0 0  0 0 0
`

	c, err := CanvasFromPPMReader(strings.NewReader(ppmFile))
	require.NoError(t, err)

	require.Equal(t, 10, c.width)
	require.Equal(t, 2, c.height)
}

func Test_CanvasFromPPM_HasCorrectColorData(t *testing.T) {
	ppmFile := `P3
4 3
255
255 127 0  0 127 255  127 255 0  255 255 255
0 0 0  255 0 0  0 255 0  0 0 255
255 255 0  0 255 255  255 0 255  127 127 127
`

	c, err := CanvasFromPPMReader(strings.NewReader(ppmFile))
	require.NoError(t, err)

	type tc struct {
		x             int
		y             int
		expectedColor colors.Color
	}

	tcs := []tc{
		{0, 0, colors.NewColor(1, 0.498, 0)},
		{1, 0, colors.NewColor(0, 0.498, 1)},
		{2, 0, colors.NewColor(0.498, 1, 0)},
		{3, 0, colors.NewColor(1, 1, 1)},
		{0, 1, colors.NewColor(0, 0, 0)},
		{1, 1, colors.NewColor(1, 0, 0)},
		{2, 1, colors.NewColor(0, 1, 0)},
		{3, 1, colors.NewColor(0, 0, 1)},
		{0, 2, colors.NewColor(1, 1, 0)},
		{1, 2, colors.NewColor(0, 1, 1)},
		{2, 2, colors.NewColor(1, 0, 1)},
		{3, 2, colors.NewColor(0.498, 0.498, 0.498)},
	}

	for i, tc := range tcs {
		t.Run(t.Name()+strconv.Itoa(i), func(t *testing.T) {
			p := c.GetPixel(tc.x, tc.y)
			require.True(t, tc.expectedColor.Equal(p), "wanted", tc.expectedColor, "got", p)
		})
	}
}

func Test_CanvasFromPPM_IgnoresComments(t *testing.T) {
	ppmFile := `P3
    # this is a comment
    2 1
    # this, too
    255
    # another comment
    255 255 255
    # oh, no, comments in the pixel data!
    255 0 255
`

	c, err := CanvasFromPPMReader(strings.NewReader(ppmFile))
	require.NoError(t, err)

	require.True(t, colors.NewColor(1, 1, 1).Equal(c.GetPixel(0, 0)))
	require.True(t, colors.NewColor(1, 0, 1).Equal(c.GetPixel(1, 0)))
}

func Test_CanvasFromPPM_ColorOverMultipleLines(t *testing.T) {
	ppmFile := `P3
    1 1
    255
    51
    153

    204
`

	c, err := CanvasFromPPMReader(strings.NewReader(ppmFile))
	require.NoError(t, err)

	require.True(t, colors.NewColor(0.2, 0.6, 0.8).Equal(c.GetPixel(0, 0)))
}

func Test_CanvasFromPPM_ColorsScaledByValue(t *testing.T) {
	ppmFile := `P3
    2 2
    100
    100 100 100  50 50 50
    75 50 25  0 0 0
`

	c, err := CanvasFromPPMReader(strings.NewReader(ppmFile))
	require.NoError(t, err)

	require.True(t, colors.NewColor(0.75, 0.5, 0.25).Equal(c.GetPixel(0, 1)))
}
