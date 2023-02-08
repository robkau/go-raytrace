package geom

import "math"

const (
	FloatComparisonEpsilon = 1e-9
)

func AlmostEqual(a float64, b float64) bool {
	return AlmostEqualWithPrecision(a, b, FloatComparisonEpsilon)
}

func AlmostEqualWithPrecision(a float64, b float64, precision float64) bool {
	return math.Abs(a-b) < precision
}

func RoundTo(n float64, places int) float64 {
	if places < 0 {
		panic("places must be nonnegative")
	}
	if places == 0 {
		return float64(int(n))
	}
	scale := math.Pow10(places)
	return math.Round(n*scale) / scale
}
