package util

func Clamp(f float64, min int, max int) int {
	r := int(f + 0.5) // round it
	if r <= min {
		r = min
	}

	if r >= max {
		r = max
	}

	return r
}
