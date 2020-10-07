package main

// This file provides structs to represent 2x2, 3x3, and 4x4 matrices, and their associated operations.
type x2Matrix struct {
	b []float64
}

func newX2Matrix() x2Matrix {
	return x2Matrix{
		b: make([]float64, 2*2),
	}
}

func (m x2Matrix) get(r, c int) float64 {
	return m.b[r*2+c]
}

func (m x2Matrix) set(r int, c int, v float64) {
	m.b[r*2+c] = v
}

func (m x2Matrix) equals(o x2Matrix) {
	// todo: almostEqual on each element
}

type x3Matrix struct {
	b []float64
}

func newX3Matrix() x3Matrix {
	return x3Matrix{
		b: make([]float64, 3*3),
	}
}

func (m x3Matrix) get(r, c int) float64 {
	return m.b[r*3+c]
}

func (m x3Matrix) set(r int, c int, v float64) {
	m.b[r*3+c] = v
}

func (m x3Matrix) equals(o x3Matrix) {
	// todo: almostEqual on each element
}

type x4Matrix struct {
	b []float64
}

func newX4Matrix() x4Matrix {
	return x4Matrix{
		b: make([]float64, 4*4),
	}
}

func (m x4Matrix) get(r, c int) float64 {
	return m.b[r*4+c]
}

func (m x4Matrix) set(r int, c int, v float64) {
	m.b[r*4+c] = v
}

func (m x4Matrix) equals(o x4Matrix) {
	// todo: almostEqual on each element
}
