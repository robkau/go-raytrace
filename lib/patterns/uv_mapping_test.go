package patterns

import (
	"github.com/robkau/go-raytrace/lib/colors"
	"github.com/robkau/go-raytrace/lib/geom"
	"github.com/robkau/go-raytrace/lib/view/canvas"
	"github.com/stretchr/testify/require"
	"math"
	"strconv"
	"strings"
	"testing"
)

func Test_CheckerPattern2D(t *testing.T) {
	checkers := NewCheckerPatternUV(2, 2, colors.Black(), colors.White())

	type tc struct {
		u        float64
		v        float64
		expected colors.Color
	}

	tcs := []tc{
		{0.0, 0.0, colors.Black()},
		{0.5, 0.0, colors.White()},
		{0.0, 0.5, colors.White()},
		{0.5, 0.5, colors.Black()},
		{1.0, 1.0, colors.Black()},
	}

	for i, tc := range tcs {
		t.Run(t.Name()+strconv.Itoa(i), func(t *testing.T) {
			require.Equal(t, tc.expected, UvPatternAt(checkers, tc.u, tc.v))
		})
	}
}

func Test_SphericalMapping_3dPoint(t *testing.T) {

	type tc struct {
		p         geom.Tuple
		expectedU float64
		expectedV float64
	}

	tcs := []tc{
		{geom.NewPoint(0, 0, -1), 0.0, 0.5},
		{geom.NewPoint(1, 0, 0), 0.25, 0.5},
		{geom.NewPoint(0, 0, 1), 0.5, 0.5},
		{geom.NewPoint(-1, 0, 0), 0.75, 0.5},
		{geom.NewPoint(0, 1, 0), 0.5, 1.0},
		{geom.NewPoint(0, -1, 0), 0.5, 0.0},
		{geom.NewPoint(math.Sqrt2/2.0, math.Sqrt2/2.0, 0), 0.25, 0.75},
	}

	for i, tc := range tcs {
		t.Run(t.Name()+strconv.Itoa(i), func(t *testing.T) {
			u, v := SphericalMap(tc.p)
			require.Equal(t, tc.expectedU, u)
			require.Equal(t, tc.expectedV, v)
		})
	}
}

func Test_PlanarMapping_3dPoint(t *testing.T) {

	type tc struct {
		p         geom.Tuple
		expectedU float64
		expectedV float64
	}

	tcs := []tc{
		{geom.NewPoint(0.25, 0, 0.5), 0.25, 0.5},
		{geom.NewPoint(0.25, 0, -0.25), 0.25, 0.75},
		{geom.NewPoint(0.25, 0.5, -0.25), 0.25, 0.75},
		{geom.NewPoint(1.25, 0, 0.5), 0.25, 0.5},
		{geom.NewPoint(0.25, 0, -1.75), 0.25, 0.25},
		{geom.NewPoint(1, 0, -1), 0.0, 0.0},
		{geom.NewPoint(0, 0, 0), 0.0, 0.0},
		{geom.NewPoint(0.0001, 0, -0.0001), 0.0001, 0.9999},
		{geom.NewPoint(1.0001, 0, -1.0001), 0.0001, 0.9999},
		{geom.NewPoint(5.0001, 0, -5.0001), 0.0001, 0.9999},
	}

	for i, tc := range tcs {
		t.Run(t.Name()+strconv.Itoa(i), func(t *testing.T) {
			u, v := PlanarMap(tc.p)
			require.True(t, geom.AlmostEqual(tc.expectedU, u))
			require.True(t, geom.AlmostEqual(tc.expectedV, v))
		})
	}
}

func Test_CylindricalMapping_3dPoint(t *testing.T) {

	type tc struct {
		p         geom.Tuple
		expectedU float64
		expectedV float64
	}

	tcs := []tc{
		{geom.NewPoint(0, 0, -1), 0.0, 0.0},
		{geom.NewPoint(0, 0.5, -1), 0.0, 0.5},
		{geom.NewPoint(0, 1, -1), 0.0, 0.0},
		{geom.NewPoint(0.70711, 0.5, -0.70711), 0.125, 0.5},
		{geom.NewPoint(1, 0.5, 0), 0.25, 0.5},
		{geom.NewPoint(0.70711, 0.5, 0.70711), 0.375, 0.5},
		{geom.NewPoint(0, -0.25, 1), 0.5, 0.75},
		{geom.NewPoint(-0.70711, 0.5, 0.70711), 0.625, 0.5},
		{geom.NewPoint(-1, 1.25, 0), 0.75, 0.25},
		{geom.NewPoint(-0.70711, 0.5, -0.70711), 0.875, 0.5},
	}

	for i, tc := range tcs {
		t.Run(t.Name()+strconv.Itoa(i), func(t *testing.T) {
			u, v := CylindricalMap(tc.p)
			require.True(t, geom.AlmostEqual(tc.expectedU, u))
			require.True(t, geom.AlmostEqual(tc.expectedV, v))
		})
	}
}

func Test_TextureMapping_SphericalMap(t *testing.T) {
	type tc struct {
		p             geom.Tuple
		expectedColor colors.Color
	}

	tcs := []tc{
		{geom.NewPoint(0.4315, 0.4670, 0.7719), colors.White()},
		{geom.NewPoint(-0.9654, 0.2552, -0.0534), colors.Black()},
		{geom.NewPoint(0.1039, 0.7090, 0.6975), colors.White()},
		{geom.NewPoint(-0.4986, -0.7856, -0.3663), colors.Black()},
		{geom.NewPoint(-0.0317, -0.9395, 0.3411), colors.Black()},
		{geom.NewPoint(0.4809, -0.7721, 0.4154), colors.Black()},
		{geom.NewPoint(0.0285, -0.9612, -0.2745), colors.Black()},
		{geom.NewPoint(-0.5734, -0.2162, -0.7903), colors.White()},
		{geom.NewPoint(0.7688, -0.1470, 0.6223), colors.Black()},
		{geom.NewPoint(-0.7652, 0.2175, 0.6060), colors.Black()},
	}

	c := NewCheckerPatternUV(16, 8, colors.Black(), colors.White())
	pattern := NewTextureMapPattern(c, SphericalMap)

	for i, tc := range tcs {
		t.Run(t.Name()+strconv.Itoa(i), func(t *testing.T) {
			require.Equal(t, tc.expectedColor, pattern.ColorAt(tc.p))
		})
	}
}

func Test_UvAlignCheck(t *testing.T) {
	m := colors.NewColor(1, 1, 1)
	ul := colors.NewColor(1, 0, 0)
	ur := colors.NewColor(1, 1, 0)
	bl := colors.NewColor(0, 1, 0)
	br := colors.NewColor(0, 1, 1)
	pattern := NewUVAlignCheck(m, ul, ur, bl, br)

	type tc struct {
		u        float64
		v        float64
		expected colors.Color
	}

	tcs := []tc{
		{0.5, 0.5, m},
		{0.1, 0.9, ul},
		{0.9, 0.9, ur},
		{0.1, 0.1, bl},
		{0.9, 0.1, br},
	}

	for i, tc := range tcs {
		t.Run(t.Name()+strconv.Itoa(i), func(t *testing.T) {
			require.Equal(t, tc.expected, UvPatternAt(pattern, tc.u, tc.v))
		})
	}
}

func Test_CubeFaceFromPoint(t *testing.T) {
	type tc struct {
		p        geom.Tuple
		expected CubeFace
	}

	tcs := []tc{
		{geom.NewPoint(-1, 0.5, -0.25), Left},
		{geom.NewPoint(1.1, -0.75, 0.8), Right},
		{geom.NewPoint(0.1, 0.6, 0.9), Front},
		{geom.NewPoint(-0.7, 0, -2), Back},
		{geom.NewPoint(0.5, 1, 0.9), Up},
		{geom.NewPoint(-0.2, -1.3, 1.1), Down},
	}

	for i, tc := range tcs {
		t.Run(t.Name()+strconv.Itoa(i), func(t *testing.T) {
			require.Equal(t, tc.expected, CubeFaceFromPoint(tc.p))
		})
	}
}

func Test_UVMap_CubeFront(t *testing.T) {
	type tc struct {
		p         geom.Tuple
		expectedU float64
		expectedV float64
	}

	tcs := []tc{
		{geom.NewPoint(-0.5, 0.5, 1), 0.25, 0.75},
		{geom.NewPoint(0.5, -0.5, 1), 0.75, 0.25},
	}

	for i, tc := range tcs {
		t.Run(t.Name()+strconv.Itoa(i), func(t *testing.T) {
			u, v := CubeUvFront(tc.p)
			require.Equal(t, tc.expectedU, u)
			require.Equal(t, tc.expectedV, v)
		})
	}
}

func Test_UVMap_CubeBack(t *testing.T) {
	type tc struct {
		p         geom.Tuple
		expectedU float64
		expectedV float64
	}

	tcs := []tc{
		{geom.NewPoint(0.5, 0.5, -1), 0.25, 0.75},
		{geom.NewPoint(-0.5, -0.5, -1), 0.75, 0.25},
	}

	for i, tc := range tcs {
		t.Run(t.Name()+strconv.Itoa(i), func(t *testing.T) {
			u, v := CubeUvBack(tc.p)
			require.Equal(t, tc.expectedU, u)
			require.Equal(t, tc.expectedV, v)
		})
	}
}

func Test_UVMap_CubeLeft(t *testing.T) {
	type tc struct {
		p         geom.Tuple
		expectedU float64
		expectedV float64
	}

	tcs := []tc{
		{geom.NewPoint(-1, 0.5, -0.5), 0.25, 0.75},
		{geom.NewPoint(-1, -0.5, 0.5), 0.75, 0.25},
	}

	for i, tc := range tcs {
		t.Run(t.Name()+strconv.Itoa(i), func(t *testing.T) {
			u, v := CubeUvLeft(tc.p)
			require.Equal(t, tc.expectedU, u)
			require.Equal(t, tc.expectedV, v)
		})
	}
}

func Test_UVMap_CubeRight(t *testing.T) {
	type tc struct {
		p         geom.Tuple
		expectedU float64
		expectedV float64
	}

	tcs := []tc{
		{geom.NewPoint(1, 0.5, 0.5), 0.25, 0.75},
		{geom.NewPoint(1, -0.5, -0.5), 0.75, 0.25},
	}

	for i, tc := range tcs {
		t.Run(t.Name()+strconv.Itoa(i), func(t *testing.T) {
			u, v := CubeUvRight(tc.p)
			require.Equal(t, tc.expectedU, u)
			require.Equal(t, tc.expectedV, v)
		})
	}
}

func Test_UVMap_CubeUp(t *testing.T) {
	type tc struct {
		p         geom.Tuple
		expectedU float64
		expectedV float64
	}

	tcs := []tc{
		{geom.NewPoint(-0.5, 1, -0.5), 0.25, 0.75},
		{geom.NewPoint(0.5, 1, 0.5), 0.75, 0.25},
	}

	for i, tc := range tcs {
		t.Run(t.Name()+strconv.Itoa(i), func(t *testing.T) {
			u, v := CubeUvUp(tc.p)
			require.Equal(t, tc.expectedU, u)
			require.Equal(t, tc.expectedV, v)
		})
	}
}

func Test_UVMap_CubeDown(t *testing.T) {
	type tc struct {
		p         geom.Tuple
		expectedU float64
		expectedV float64
	}

	tcs := []tc{
		{geom.NewPoint(-0.5, -1, 0.5), 0.25, 0.75},
		{geom.NewPoint(0.5, -1, -0.5), 0.75, 0.25},
	}

	for i, tc := range tcs {
		t.Run(t.Name()+strconv.Itoa(i), func(t *testing.T) {
			u, v := CubeUvDown(tc.p)
			require.Equal(t, tc.expectedU, u)
			require.Equal(t, tc.expectedV, v)
		})
	}
}

func Test_MappedCube_Colors(t *testing.T) {
	pattern := NewPrismaticCube()

	type tc struct {
		p             geom.Tuple
		expectedColor colors.Color
	}

	tcs := []tc{
		// left
		{geom.NewPoint(-1, 0, 0), colors.Yellow()},
		{geom.NewPoint(-1, 0.9, -0.9), colors.Cyan()},
		{geom.NewPoint(-1, 0.9, 0.9), colors.Red()},
		{geom.NewPoint(-1, -0.9, -0.9), colors.Blue()},
		{geom.NewPoint(-1, -0.9, 0.9), colors.Brown()},
		// front
		{geom.NewPoint(0, 0, 1), colors.Cyan()},
		{geom.NewPoint(-0.9, 0.9, 1), colors.Red()},
		{geom.NewPoint(0.9, 0.9, 1), colors.Yellow()},
		{geom.NewPoint(-0.9, -0.9, 1), colors.Brown()},
		{geom.NewPoint(0.9, -0.9, 1), colors.Green()},
		// right
		{geom.NewPoint(1, 0, 0), colors.Red()},
		{geom.NewPoint(1, 0.9, 0.9), colors.Yellow()},
		{geom.NewPoint(1, 0.9, -0.9), colors.Purple()},
		{geom.NewPoint(1, -0.9, 0.9), colors.Green()},
		{geom.NewPoint(1, -0.9, -0.9), colors.White()},
		// back
		{geom.NewPoint(0, 0, -1), colors.Green()},
		{geom.NewPoint(0.9, 0.9, -1), colors.Purple()},
		{geom.NewPoint(-0.9, 0.9, -1), colors.Cyan()},
		{geom.NewPoint(0.9, -0.9, -1), colors.White()},
		{geom.NewPoint(-0.9, -0.9, -1), colors.Blue()},
		// up
		{geom.NewPoint(0, 1, 0), colors.Brown()},
		{geom.NewPoint(-0.9, 1, -0.9), colors.Cyan()},
		{geom.NewPoint(0.9, 1, -0.9), colors.Purple()},
		{geom.NewPoint(-0.9, 1, 0.9), colors.Red()},
		{geom.NewPoint(0.9, 1, 0.9), colors.Yellow()},
		// down
		{geom.NewPoint(0, -1, 0), colors.Purple()},
		{geom.NewPoint(-0.9, -1, 0.9), colors.Brown()},
		{geom.NewPoint(0.9, -1, 0.9), colors.Green()},
		{geom.NewPoint(-0.9, -1, -0.9), colors.Blue()},
		{geom.NewPoint(0.9, -1, -0.9), colors.White()},
	}

	for i, tc := range tcs {
		t.Run(t.Name()+strconv.Itoa(i), func(t *testing.T) {
			require.Equal(t, tc.expectedColor, pattern.ColorAt(tc.p))
		})
	}
}

func Test_UVImage_FromPPMCanvas(t *testing.T) {
	ppmFile := `P3
    10 10
    10
    0 0 0  1 1 1  2 2 2  3 3 3  4 4 4  5 5 5  6 6 6  7 7 7  8 8 8  9 9 9
    1 1 1  2 2 2  3 3 3  4 4 4  5 5 5  6 6 6  7 7 7  8 8 8  9 9 9  0 0 0
    2 2 2  3 3 3  4 4 4  5 5 5  6 6 6  7 7 7  8 8 8  9 9 9  0 0 0  1 1 1
    3 3 3  4 4 4  5 5 5  6 6 6  7 7 7  8 8 8  9 9 9  0 0 0  1 1 1  2 2 2
    4 4 4  5 5 5  6 6 6  7 7 7  8 8 8  9 9 9  0 0 0  1 1 1  2 2 2  3 3 3
    5 5 5  6 6 6  7 7 7  8 8 8  9 9 9  0 0 0  1 1 1  2 2 2  3 3 3  4 4 4
    6 6 6  7 7 7  8 8 8  9 9 9  0 0 0  1 1 1  2 2 2  3 3 3  4 4 4  5 5 5
    7 7 7  8 8 8  9 9 9  0 0 0  1 1 1  2 2 2  3 3 3  4 4 4  5 5 5  6 6 6
    8 8 8  9 9 9  0 0 0  1 1 1  2 2 2  3 3 3  4 4 4  5 5 5  6 6 6  7 7 7
    9 9 9  0 0 0  1 1 1  2 2 2  3 3 3  4 4 4  5 5 5  6 6 6  7 7 7  8 8 8`

	canvas, err := canvas.CanvasFromPPMReader(strings.NewReader(ppmFile))
	require.NoError(t, err)

	pattern := NewUVImage(canvas)

	type tc struct {
		u             float64
		v             float64
		expectedColor colors.Color
	}

	tcs := []tc{
		{0, 0, colors.NewColor(0.9, 0.9, 0.9)},
		{0.3, 0, colors.NewColor(0.2, 0.2, 0.2)},
		{0.6, 0.3, colors.NewColor(0.1, 0.1, 0.1)},
		{1, 1, colors.NewColor(0.9, 0.9, 0.9)},
	}

	for i, tc := range tcs {
		t.Run(t.Name()+strconv.Itoa(i), func(t *testing.T) {
			require.Equal(t, tc.expectedColor, UvPatternAt(pattern, tc.u, tc.v))
		})
	}

}
