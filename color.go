package main

type color struct {
	r float64
	g float64
	b float64
}

func (c color) add(other color) color {
	return color{
		c.r + other.r,
		c.g + other.g,
		c.b + other.b,
	}
}

func (c color) sub(other color) color {
	return color{
		c.r - other.r,
		c.g - other.g,
		c.b - other.b,
	}
}
func (c color) mul(other color) color {
	return color{
		c.r * other.r,
		c.g * other.g,
		c.b * other.b,
	}
}

func (c color) mulBy(x float64) color {
	return color{
		c.r * x,
		c.g * x,
		c.b * x,
	}
}

func (c color) equal(other color) bool {
	return almostEqual(c.r, other.r) &&
		almostEqual(c.g, other.g) &&
		almostEqual(c.b, other.b)
}
