package geom

import (
	"log"
	"math"
	"sync"
)

// This file provides structs to represent 2x2, 3x3, and 4x4 matrices, and their associated operations.
type X2Matrix struct {
	b [4]float64
}

func NewX2Matrix() X2Matrix {
	return X2Matrix{
		b: [4]float64{},
	}
}

func NewX2MatrixWith(aa, ab, ba, bb float64) X2Matrix {
	m := NewX2Matrix()
	m.b[0] = aa
	m.b[1] = ab
	m.b[2] = ba
	m.b[3] = bb
	return m
}

func (m X2Matrix) Get(r, c int) float64 {
	return m.b[r*2+c]
}

func (m X2Matrix) Set(r int, c int, v float64) X2Matrix {
	ret := NewX2MatrixWith(m.b[0], m.b[1], m.b[2], m.b[3])
	ret.b[r*2+c] = v
	return ret
}

func (m X2Matrix) Equals(o X2Matrix) bool {
	return compareFloat4(m.b, o.b)
}

func (m X2Matrix) Determinant() float64 {
	return m.Get(0, 0)*m.Get(1, 1) - m.Get(0, 1)*m.Get(1, 0)
}

type X3Matrix struct {
	b [9]float64
}

func NewX3Matrix() X3Matrix {
	return X3Matrix{
		b: [9]float64{},
	}
}

func NewX3MatrixWith(aa, ab, ac, ba, bb, bc, ca, cb, cc float64) X3Matrix {
	m := NewX3Matrix()
	m.b[0] = aa
	m.b[1] = ab
	m.b[2] = ac
	m.b[3] = ba
	m.b[4] = bb
	m.b[5] = bc
	m.b[6] = ca
	m.b[7] = cb
	m.b[8] = cc
	return m
}

func (m X3Matrix) Get(r, c int) float64 {
	return m.b[r*3+c]
}

func (m X3Matrix) Set(r int, c int, v float64) X3Matrix {
	ret := NewX3MatrixWith(m.b[0], m.b[1], m.b[2], m.b[3], m.b[4], m.b[5], m.b[6], m.b[7], m.b[8])
	ret.b[r*3+c] = v
	return ret
}

func (m X3Matrix) Equals(o X3Matrix) bool {
	return compareFloat9(m.b, o.b)
}

func (m X3Matrix) Submatrix(rr, cr int) X2Matrix {
	n := NewX2Matrix()

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
			n = n.Set(r-rOffset, c-cOffset, m.Get(r, c))
		}
	}
	return n
}

func (m X3Matrix) Minor(rr, cr int) float64 {
	n := m.Submatrix(rr, cr)
	return n.Determinant()
}

func (m X3Matrix) Cofactor(rr, cr int) float64 {
	n := m.Minor(rr, cr)
	if (rr+cr)%2 != 0 {
		n *= -1
	}
	return n
}

func (m X3Matrix) Determinant() float64 {
	d := 0.0
	for c := 0; c < 3; c++ {
		d += m.Get(0, c) * m.Cofactor(0, c)
	}
	return d
}

type X4Matrix struct {
	b [16]float64
}

func NewX4Matrix() X4Matrix {
	return X4Matrix{
		b: [16]float64{},
	}
}

func NewX4MatrixFrom(t [16]float64) X4Matrix {
	m := NewX4Matrix()
	m.b[0] = t[0]
	m.b[1] = t[1]
	m.b[2] = t[2]
	m.b[3] = t[3]
	m.b[4] = t[4]
	m.b[5] = t[5]
	m.b[6] = t[6]
	m.b[7] = t[7]
	m.b[8] = t[8]
	m.b[9] = t[9]
	m.b[10] = t[10]
	m.b[11] = t[11]
	m.b[12] = t[12]
	m.b[13] = t[13]
	m.b[14] = t[14]
	m.b[15] = t[15]

	return m
}

func NewX4MatrixWith(aa, ab, ac, ad, ba, bb, bc, bd, ca, cb, cc, cd, da, db, dc, dd float64) X4Matrix {
	m := NewX4Matrix()
	m.b[0] = aa
	m.b[1] = ab
	m.b[2] = ac
	m.b[3] = ad
	m.b[4] = ba
	m.b[5] = bb
	m.b[6] = bc
	m.b[7] = bd
	m.b[8] = ca
	m.b[9] = cb
	m.b[10] = cc
	m.b[11] = cd
	m.b[12] = da
	m.b[13] = db
	m.b[14] = dc
	m.b[15] = dd
	return m
}

func (m X4Matrix) Get(r, c int) float64 {
	return m.b[r*4+c]
}

func (m X4Matrix) Set(r int, c int, v float64) X4Matrix {
	ret := NewX4Matrix()
	ret.b = m.b
	ret.b[r*4+c] = v
	return ret
}

func (m X4Matrix) Equals(o X4Matrix) bool {
	return compareFloat16(m.b, o.b)
}

func (m X4Matrix) Transpose() X4Matrix {
	n := NewX4Matrix()
	for r := 0; r < 4; r++ {
		for c := 0; c < 4; c++ {
			n = n.Set(c, r, m.Get(r, c))
		}
	}
	return n
}

func (m X4Matrix) MulX4Matrix(o X4Matrix) X4Matrix {
	n := NewX4Matrix()
	for r := 0; r < 4; r++ {
		for c := 0; c < 4; c++ {
			n = n.Set(r, c,
				m.Get(r, 0)*o.Get(0, c)+
					m.Get(r, 1)*o.Get(1, c)+
					m.Get(r, 2)*o.Get(2, c)+
					m.Get(r, 3)*o.Get(3, c))
		}
	}
	return n
}

func (m X4Matrix) MulTuple(t Tuple) Tuple {
	n := NewTuple(0, 0, 0, 0)

	n.X = m.Get(0, 0)*t.X +
		m.Get(0, 1)*t.Y +
		m.Get(0, 2)*t.Z +
		m.Get(0, 3)*t.C

	n.Y = m.Get(1, 0)*t.X +
		m.Get(1, 1)*t.Y +
		m.Get(1, 2)*t.Z +
		m.Get(1, 3)*t.C

	n.Z = m.Get(2, 0)*t.X +
		m.Get(2, 1)*t.Y +
		m.Get(2, 2)*t.Z +
		m.Get(2, 3)*t.C

	n.C = m.Get(3, 0)*t.X +
		m.Get(3, 1)*t.Y +
		m.Get(3, 2)*t.Z +
		m.Get(3, 3)*t.C

	return n
}

func (m X4Matrix) RoundTo(places int) X4Matrix {
	n := NewX4Matrix()
	scale := math.Pow10(places)

	for r := 0; r < 4; r++ {
		for c := 0; c < 4; c++ {
			n = n.Set(r, c, math.Round(m.Get(r, c)*scale)/scale)
		}
	}

	return n
}

func (m X4Matrix) Submatrix(rr, cr int) X3Matrix {
	n := NewX3Matrix()

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
			n = n.Set(r-rOffset, c-cOffset, m.Get(r, c))
		}
	}
	return n
}

func (m X4Matrix) Minor(rr, cr int) float64 {
	n := m.Submatrix(rr, cr)
	return n.Determinant()
}

func (m X4Matrix) Cofactor(rr, cr int) float64 {
	n := m.Minor(rr, cr)
	if (rr+cr)%2 != 0 {
		n *= -1
	}
	return n
}

func (m X4Matrix) Determinant() float64 {
	d := 0.0
	for c := 0; c < 4; c++ {
		d += m.Get(0, c) * m.Cofactor(0, c)
	}
	return d
}

// todo make me faster
// time spent on mutex
// time spent on hashing
var rwX4MatrixCache = sync.RWMutex{}
var x4MatrixInvertCache = make(map[X4Matrix]X4Matrix)

func cachedX4MatrixInverter(m X4Matrix) X4Matrix {
	rwX4MatrixCache.RLock()
	if ret, ok := x4MatrixInvertCache[m]; ok {
		rwX4MatrixCache.RUnlock()
		// was cached
		return ret
	}
	rwX4MatrixCache.RUnlock()

	// compute it
	rwX4MatrixCache.Lock()
	defer rwX4MatrixCache.Unlock()
	cm := m.doInvert()
	x4MatrixInvertCache[m] = cm
	return cm
}

func (m X4Matrix) Invertable() bool {
	return m.Determinant() != 0
}

func (m X4Matrix) Invert() X4Matrix {
	return cachedX4MatrixInverter(m)
}

func (m X4Matrix) doInvert() X4Matrix {
	if !m.Invertable() {
		log.Fatalf("not invertable")
	}
	invert := NewX4Matrix()
	d := m.Determinant()
	for r := 0; r < 4; r++ {
		for c := 0; c < 4; c++ {
			co := m.Cofactor(r, c)
			invert = invert.Set(c, r, co/d)
		}
	}
	return invert
}

func (m X4Matrix) Copy() X4Matrix {
	return NewX4MatrixFrom(m.b)
}

func NewIdentityMatrixX4() X4Matrix {
	return NewX4MatrixWith(1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1)
}

func compareFloat4(m, o [4]float64) bool {
	for i := range m {
		if !AlmostEqual(m[i], o[i]) {
			return false
		}
	}
	return true
}

func compareFloat9(m, o [9]float64) bool {
	for i := range m {
		if !AlmostEqual(m[i], o[i]) {
			return false
		}
	}
	return true
}

func compareFloat16(m, o [16]float64) bool {
	for i := range m {
		if !AlmostEqual(m[i], o[i]) {
			return false
		}
	}
	return true
}
