package geom

import "math"

const (
	FloatComparisonEpsilon = 1e-9
)

func AlmostEqual(a float64, b float64) bool {
	return math.Abs(a-b) < FloatComparisonEpsilon
}
