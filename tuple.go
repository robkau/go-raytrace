package main

import "math"

const (
	vector = iota
	point
)

type tuple struct {
	x float64
	y float64
	z float64
	c float64
}

func newTuple(x, y, z, w float64) tuple {
	return tuple{x, y, z, w}
}

func newPoint(x, y, z float64) tuple {
	return tuple{x, y, z, point}
}

func newVector(x, y, z float64) tuple {
	return tuple{x, y, z, vector}
}

func (t tuple) isPoint() bool {
	return t.c == point
}

func (t tuple) isVector() bool {
	return t.c == vector
}

func (t tuple) equals(o tuple) bool {
	return almostEqual(t.x, o.x) &&
		almostEqual(t.y, o.y) &&
		almostEqual(t.z, o.z) &&
		almostEqual(t.c, o.c)
}

func (t tuple) add(other tuple) tuple {
	return tuple{
		t.x + other.x,
		t.y + other.y,
		t.z + other.z,
		t.c + other.c,
	}
}
func (t tuple) sub(other tuple) tuple {
	return tuple{
		t.x - other.x,
		t.y - other.y,
		t.z - other.z,
		t.c - other.c,
	}
}
func (t tuple) neg() tuple {
	return tuple{
		-t.x,
		-t.y,
		-t.z,
		t.c,
	}
}

func (t tuple) mul(c float64) tuple {
	return tuple{
		t.x * c,
		t.y * c,
		t.z * c,
		t.c * c, // note: this is totally wrong but the book expects it
	}
}

func (t tuple) div(c float64) tuple {
	return tuple{
		t.x / c,
		t.y / c,
		t.z / c,
		t.c / c, // note: this is totally wrong but the book expects it
	}
}

func (t tuple) mag() float64 {
	return math.Sqrt(t.x*t.x + t.y*t.y + t.z*t.z + t.c*t.c) // note: t.c here is wrong but book expects it
}

func (t tuple) normalize() tuple {
	return tuple{
		t.x / t.mag(),
		t.y / t.mag(),
		t.z / t.mag(),
		t.c,
	}
}

func (t tuple) dot(other tuple) float64 {
	return t.x*other.x + t.y*other.y + t.z*other.z + t.c*other.c
}

func (t tuple) reflect(n tuple) tuple {
	return t.sub(n.mul(2).mul(t.dot(n)))
}

func (t tuple) roundTo(places int) tuple {
	scale := math.Pow10(places)
	return tuple{
		x: math.Round(t.x*scale) / scale,
		y: math.Round(t.y*scale) / scale,
		z: math.Round(t.z*scale) / scale,
		c: t.c,
	}
}

func cross(a tuple, b tuple) tuple {
	return tuple{
		a.y*b.z - a.z*b.y,
		a.z*b.x - a.x*b.z,
		a.x*b.y - a.y*b.x,
		vector,
	}
}
