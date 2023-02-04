package patterns

import (
	"github.com/robkau/go-raytrace/lib/colors"
	"github.com/robkau/go-raytrace/lib/geom"
	"github.com/stretchr/testify/require"
	"math"
	"strconv"
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
