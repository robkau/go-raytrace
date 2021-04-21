package geom

import (
	"math"
)

const (
	Vector = iota
	point
)

type Tuple struct {
	X float64
	Y float64
	Z float64
	C float64
}

func NewTuple(x, y, z, w float64) Tuple {
	return Tuple{x, y, z, w}
}

func NewPoint(x, y, z float64) Tuple {
	return Tuple{x, y, z, point}
}

func NewVector(x, y, z float64) Tuple {
	return Tuple{x, y, z, Vector}
}

func (t Tuple) IsPoint() bool {
	return t.C == point
}

func (t Tuple) IsVector() bool {
	return t.C == Vector
}

func (t Tuple) Equals(o Tuple) bool {
	return AlmostEqual(t.X, o.X) &&
		AlmostEqual(t.Y, o.Y) &&
		AlmostEqual(t.Z, o.Z) &&
		AlmostEqual(t.C, o.C)
}

func (t Tuple) Add(other Tuple) Tuple {
	return Tuple{
		t.X + other.X,
		t.Y + other.Y,
		t.Z + other.Z,
		t.C + other.C,
	}
}
func (t Tuple) Sub(other Tuple) Tuple {
	return Tuple{
		t.X - other.X,
		t.Y - other.Y,
		t.Z - other.Z,
		t.C - other.C,
	}
}
func (t Tuple) Neg() Tuple {
	return Tuple{
		-t.X,
		-t.Y,
		-t.Z,
		t.C,
	}
}

func (t Tuple) Mul(c float64) Tuple {
	return Tuple{
		t.X * c,
		t.Y * c,
		t.Z * c,
		t.C * c, // note: this is totally wrong but the book expects it
	}
}

func (t Tuple) Div(c float64) Tuple {
	return Tuple{
		t.X / c,
		t.Y / c,
		t.Z / c,
		t.C / c, // note: this is totally wrong but the book expects it
	}
}

func (t Tuple) Mag() float64 {
	return math.Sqrt(t.X*t.X + t.Y*t.Y + t.Z*t.Z + t.C*t.C) // note: t.C here is wrong but book expects it
}

func (t Tuple) Normalize() Tuple {
	return Tuple{
		t.X / t.Mag(),
		t.Y / t.Mag(),
		t.Z / t.Mag(),
		t.C,
	}
}

func (t Tuple) Dot(other Tuple) float64 {
	return t.X*other.X + t.Y*other.Y + t.Z*other.Z + t.C*other.C
}

func (t Tuple) Reflect(n Tuple) Tuple {
	return t.Sub(n.Mul(2).Mul(t.Dot(n)))
}

func (t Tuple) RoundTo(places int) Tuple {
	scale := math.Pow10(places)
	return Tuple{
		X: math.Round(t.X*scale) / scale,
		Y: math.Round(t.Y*scale) / scale,
		Z: math.Round(t.Z*scale) / scale,
		C: t.C,
	}
}

func Cross(a Tuple, b Tuple) Tuple {
	return Tuple{
		a.Y*b.Z - a.Z*b.Y,
		a.Z*b.X - a.X*b.Z,
		a.X*b.Y - a.Y*b.X,
		Vector,
	}
}
