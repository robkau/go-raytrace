package main

import "math"

const (
	floatComparisonEpsilon = 1e-9
)

func almostEqual(a float64, b float64) bool {
	return math.Abs(a-b) < floatComparisonEpsilon
}
