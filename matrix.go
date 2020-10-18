package main

import "math"

// This file provides structs to represent 2x2, 3x3, and 4x4 matrices, and their associated operations.
type x2Matrix struct {
	b []float64
}

func newX2Matrix() x2Matrix {
	return x2Matrix{
		b: make([]float64, 2*2),
	}
}

func newX2MatrixWith(aa, ab, ba, bb float64) x2Matrix {
	m := newX2Matrix()
	m.set(0, 0, aa)
	m.set(0, 1, ab)
	m.set(1, 0, ba)
	m.set(1, 1, bb)
	return m
}

func (m x2Matrix) get(r, c int) float64 {
	return m.b[r*2+c]
}

func (m x2Matrix) set(r int, c int, v float64) {
	m.b[r*2+c] = v
}

func (m x2Matrix) equals(o x2Matrix) bool {
	return compareFloatSlice(m.b, o.b)
}

func (m x2Matrix) determinant() float64 {
	return m.get(0, 0)*m.get(1, 1) - m.get(0, 1)*m.get(1, 0)
}

type x3Matrix struct {
	b []float64
}

func newX3Matrix() x3Matrix {
	return x3Matrix{
		b: make([]float64, 3*3),
	}
}

func newX3MatrixWith(aa, ab, ac, ba, bb, bc, ca, cb, cc float64) x3Matrix {
	m := newX3Matrix()
	m.set(0, 0, aa)
	m.set(0, 1, ab)
	m.set(0, 2, ac)
	m.set(1, 0, ba)
	m.set(1, 1, bb)
	m.set(1, 2, bc)
	m.set(2, 0, ca)
	m.set(2, 1, cb)
	m.set(2, 2, cc)
	return m
}

func (m x3Matrix) get(r, c int) float64 {
	return m.b[r*3+c]
}

func (m x3Matrix) set(r int, c int, v float64) {
	m.b[r*3+c] = v
}

func (m x3Matrix) equals(o x3Matrix) bool {
	return compareFloatSlice(m.b, o.b)
}

func (m x3Matrix) submatrix(rr, cr int) x2Matrix {
	n := newX2Matrix()

	rOffset := 0
	for r := 0; r < 3; r++ {
		if r == rr {
			rOffset++
			continue
		}
		cOffset := 0
		for c := 0; c < 3; c++ {
			if c == cr {
				cOffset++
				continue
			}
			n.set(r-rOffset, c-cOffset, m.get(r, c))
		}
	}
	return n
}

func (m x3Matrix) minor(rr, cr int) float64 {
	n := m.submatrix(rr, cr)
	return n.determinant()
}

func (m x3Matrix) cofactor(rr, cr int) float64 {
	n := m.minor(rr, cr)
	if (rr+cr)%2 != 0 {
		n *= -1
	}
	return n
}

func (m x3Matrix) determinant() float64 {
	d := 0.0
	for c := 0; c < 3; c++ {
		d += m.get(0, c) * m.cofactor(0, c)
	}
	return d
}

type x4Matrix struct {
	b []float64
}

func newX4Matrix() x4Matrix {
	return x4Matrix{
		b: make([]float64, 4*4),
	}
}

func newX4MatrixWith(aa, ab, ac, ad, ba, bb, bc, bd, ca, cb, cc, cd, da, db, dc, dd float64) x4Matrix {
	m := newX4Matrix()
	m.set(0, 0, aa)
	m.set(0, 1, ab)
	m.set(0, 2, ac)
	m.set(0, 3, ad)
	m.set(1, 0, ba)
	m.set(1, 1, bb)
	m.set(1, 2, bc)
	m.set(1, 3, bd)
	m.set(2, 0, ca)
	m.set(2, 1, cb)
	m.set(2, 2, cc)
	m.set(2, 3, cd)
	m.set(3, 0, da)
	m.set(3, 1, db)
	m.set(3, 2, dc)
	m.set(3, 3, dd)
	return m
}

func (m x4Matrix) get(r, c int) float64 {
	return m.b[r*4+c]
}

func (m x4Matrix) set(r int, c int, v float64) {
	m.b[r*4+c] = v
}

func (m x4Matrix) equals(o x4Matrix) bool {
	return compareFloatSlice(m.b, o.b)
}

func (m x4Matrix) transpose() x4Matrix {
	n := newX4Matrix()
	for r := 0; r < 4; r++ {
		for c := 0; c < 4; c++ {
			n.set(c, r, m.get(r, c))
		}
	}
	return n
}

func (m x4Matrix) mulX4Matrix(o x4Matrix) x4Matrix {
	n := newX4Matrix()
	for r := 0; r < 4; r++ {
		for c := 0; c < 4; c++ {
			n.set(r, c,
				m.get(r, 0)*o.get(0, c)+
					m.get(r, 1)*o.get(1, c)+
					m.get(r, 2)*o.get(2, c)+
					m.get(r, 3)*o.get(3, c))
		}
	}
	return n
}

func (m x4Matrix) mulTuple(t tuple) tuple {
	n := newTuple(0, 0, 0, 0)

	n.x = m.get(0, 0)*t.x +
		m.get(0, 1)*t.y +
		m.get(0, 2)*t.z +
		m.get(0, 3)*t.c

	n.y = m.get(1, 0)*t.x +
		m.get(1, 1)*t.y +
		m.get(1, 2)*t.z +
		m.get(1, 3)*t.c

	n.z = m.get(2, 0)*t.x +
		m.get(2, 1)*t.y +
		m.get(2, 2)*t.z +
		m.get(2, 3)*t.c

	n.c = m.get(3, 0)*t.x +
		m.get(3, 1)*t.y +
		m.get(3, 2)*t.z +
		m.get(3, 3)*t.c

	return n
}

func (m x4Matrix) roundTo(places int) x4Matrix {
	n := newX4Matrix()
	scale := math.Pow10(places)

	for r := 0; r < 4; r++ {
		for c := 0; c < 4; c++ {
			n.set(r, c, math.Round(m.get(r, c)*scale)/scale)
		}
	}

	return n
}

func (m x4Matrix) submatrix(rr, cr int) x3Matrix {
	n := newX3Matrix()

	rOffset := 0
	for r := 0; r < 4; r++ {
		if r == rr {
			rOffset++
			continue
		}
		cOffset := 0
		for c := 0; c < 4; c++ {
			if c == cr {
				cOffset++
				continue
			}
			n.set(r-rOffset, c-cOffset, m.get(r, c))
		}
	}
	return n
}

func (m x4Matrix) minor(rr, cr int) float64 {
	n := m.submatrix(rr, cr)
	return n.determinant()
}

func (m x4Matrix) cofactor(rr, cr int) float64 {
	n := m.minor(rr, cr)
	if (rr+cr)%2 != 0 {
		n *= -1
	}
	return n
}

func (m x4Matrix) determinant() float64 {
	d := 0.0
	for c := 0; c < 4; c++ {
		d += m.get(0, c) * m.cofactor(0, c)
	}
	return d
}

func (m x4Matrix) invertable() bool {
	return m.determinant() != 0
}

func (m x4Matrix) invert() x4Matrix {
	if !m.invertable() {
		panic("matrix not invertable")
	}

	n := newX4Matrix()
	d := m.determinant()

	for r := 0; r < 4; r++ {
		for c := 0; c < 4; c++ {
			co := m.cofactor(r, c)
			n.set(c, r, co/d)
		}
	}
	return n
}

func newIdentityMatrixX4() x4Matrix {
	return newX4MatrixWith(1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1)
}

func compareFloatSlice(m, o []float64) bool {
	if len(m) != len(o) {
		return false
	}
	for i := range m {
		if !almostEqual(m[i], o[i]) {
			return false
		}
	}
	return true
}
