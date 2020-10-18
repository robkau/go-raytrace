package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_New2x2(t *testing.T) {
	m := newX2Matrix()

	assert.Len(t, m.b, 4)
}

func Test_GetSet2x2(t *testing.T) {
	m := newX2Matrix()
	m.set(0, 0, 1.5)
	m.set(0, 1, 2.7)
	m.set(1, 0, 3.9)
	m.set(1, 1, -4.1)

	assert.Equal(t, 1.5, m.get(0, 0))
	assert.Equal(t, 2.7, m.get(0, 1))
	assert.Equal(t, 3.9, m.get(1, 0))
	assert.Equal(t, -4.1, m.get(1, 1))
}

func Test_Equals2x2(t *testing.T) {
	m := newX2MatrixWith(0, 1, 2, 3)
	n := newX2MatrixWith(0, 1, 2, 3)
	o := newX2MatrixWith(0, 1.00000001, 2, 3)

	assert.True(t, m.equals(m))
	assert.True(t, n.equals(n))
	assert.True(t, m.equals(n))
	assert.True(t, n.equals(m))
	assert.True(t, !m.equals(o))
	assert.True(t, !n.equals(o))
	assert.True(t, !o.equals(m))
	assert.True(t, !o.equals(n))
}

func Test_DeterminantX2(t *testing.T) {
	m := newX2MatrixWith(
		1, 5,
		-3, 2)

	assert.Equal(t, 17.0, m.determinant())
}

func Test_New3x3(t *testing.T) {
	m := newX3Matrix()

	assert.Len(t, m.b, 9)
}

func Test_GetSet3x3(t *testing.T) {
	m := newX3Matrix()
	m.set(0, 0, 1.5)
	m.set(0, 1, 2.7)
	m.set(0, 2, 2.9)
	m.set(1, 0, 3.9)
	m.set(1, 1, -4.1)
	m.set(1, 2, -4.3)
	m.set(2, 0, 3.12)
	m.set(2, 1, -4.122)
	m.set(2, 2, -4.3333)

	assert.Equal(t, 1.5, m.get(0, 0))
	assert.Equal(t, 2.7, m.get(0, 1))
	assert.Equal(t, 2.9, m.get(0, 2))
	assert.Equal(t, 3.9, m.get(1, 0))
	assert.Equal(t, -4.1, m.get(1, 1))
	assert.Equal(t, -4.3, m.get(1, 2))
	assert.Equal(t, 3.12, m.get(2, 0))
	assert.Equal(t, -4.122, m.get(2, 1))
	assert.Equal(t, -4.3333, m.get(2, 2))
}

func Test_Equals3x3(t *testing.T) {
	m := newX3MatrixWith(0, 1, 2, 3, 4, 5, 6, 7, 8)
	n := newX3MatrixWith(0, 1, 2, 3, 4, 5, 6, 7, 8)
	o := newX3MatrixWith(0, 1.00000001, 2, 3, 4, 5, 6, 7, 8)

	assert.True(t, m.equals(m))
	assert.True(t, n.equals(n))
	assert.True(t, m.equals(n))
	assert.True(t, n.equals(m))
	assert.True(t, !m.equals(o))
	assert.True(t, !n.equals(o))
	assert.True(t, !o.equals(m))
	assert.True(t, !o.equals(n))
}

func Test_X3_Minor(t *testing.T) {
	m := newX3MatrixWith(3, 5, 0, 2, -1, -7, 6, -1, 5)

	b := m.submatrix(1, 0)

	assert.Equal(t, 25.0, b.determinant())
	assert.Equal(t, 25.0, m.minor(1, 0))
}

func Test_X3_Cofactor(t *testing.T) {
	m := newX3MatrixWith(3, 5, 0, 2, -1, -7, 6, -1, 5)

	assert.Equal(t, -12.0, m.minor(0, 0))
	assert.Equal(t, -12.0, m.cofactor(0, 0))
	assert.Equal(t, 25.0, m.minor(1, 0))
	assert.Equal(t, -25.0, m.cofactor(1, 0))
}

func Test_X3_Submatrix(t *testing.T) {
	m := newX3MatrixWith(1, 5, 0, -3, 2, 7, 0, 6, -3)

	n := m.submatrix(0, 2)

	expected := newX2MatrixWith(-3, 2, 0, 6)
	assert.Equal(t, expected, n)
}

func Test_X3_Determinant(t *testing.T) {
	m := newX3MatrixWith(1, 2, 6, -5, 8, -4, 2, 6, 4)

	assert.Equal(t, 56.0, m.cofactor(0, 0))
	assert.Equal(t, 12.0, m.cofactor(0, 1))
	assert.Equal(t, -46.0, m.cofactor(0, 2))
	assert.Equal(t, -196.0, m.determinant())
}

func Test_New4x4(t *testing.T) {
	m := newX4Matrix()

	assert.Len(t, m.b, 16)
}

func Test_GetSet4x4(t *testing.T) {
	m := newX4Matrix()
	m.set(0, 0, 1.5)
	m.set(0, 1, 2.7)
	m.set(0, 2, 2.9)
	m.set(0, 3, 2.915)
	m.set(1, 0, 3.9)
	m.set(1, 1, -4.1)
	m.set(1, 2, -4.3)
	m.set(1, 3, -8.3)
	m.set(2, 0, 3.12)
	m.set(2, 1, -4.122)
	m.set(2, 2, -4.3333)
	m.set(2, 3, 4.3333)
	m.set(3, 0, 1.1)
	m.set(3, 1, 1.33)
	m.set(3, 2, 1.66)
	m.set(3, 3, 0.02)

	assert.Equal(t, 1.5, m.get(0, 0))
	assert.Equal(t, 2.7, m.get(0, 1))
	assert.Equal(t, 2.9, m.get(0, 2))
	assert.Equal(t, 2.915, m.get(0, 3))
	assert.Equal(t, 3.9, m.get(1, 0))
	assert.Equal(t, -4.1, m.get(1, 1))
	assert.Equal(t, -4.3, m.get(1, 2))
	assert.Equal(t, -8.3, m.get(1, 3))
	assert.Equal(t, 3.12, m.get(2, 0))
	assert.Equal(t, -4.122, m.get(2, 1))
	assert.Equal(t, -4.3333, m.get(2, 2))
	assert.Equal(t, 4.3333, m.get(2, 3))
	assert.Equal(t, 1.1, m.get(3, 0))
	assert.Equal(t, 1.33, m.get(3, 1))
	assert.Equal(t, 1.66, m.get(3, 2))
	assert.Equal(t, 0.02, m.get(3, 3))
}

func Test_Equals4x4(t *testing.T) {
	m := newX4MatrixWith(0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15)
	n := newX4MatrixWith(0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15)
	o := newX4MatrixWith(0, 1.00000001, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15)

	assert.True(t, m.equals(m))
	assert.True(t, n.equals(n))
	assert.True(t, m.equals(n))
	assert.True(t, n.equals(m))
	assert.True(t, !m.equals(o))
	assert.True(t, !n.equals(o))
	assert.True(t, !o.equals(m))
	assert.True(t, !o.equals(n))
}

func Test_MulX4Matrix(t *testing.T) {
	m := newX4MatrixWith(1, 2, 3, 4, 5, 6, 7, 8, 9, 8, 7, 6, 5, 4, 3, 2)
	n := newX4MatrixWith(-2, 1, 2, 3, 3, 2, 1, -1, 4, 3, 6, 5, 1, 2, 7, 8)

	o := m.mulX4Matrix(n)

	expected := newX4MatrixWith(20, 22, 50, 48, 44, 54, 114, 108, 40, 58, 110, 102, 16, 26, 46, 42)
	assert.True(t, o.equals(expected))
}

func Test_MulX4Tuple(t *testing.T) {
	m := newX4MatrixWith(1, 2, 3, 4, 2, 4, 4, 2, 8, 6, 4, 1, 0, 0, 0, 1)
	tu := newTuple(1, 2, 3, 1)

	expected := newTuple(18, 24, 33, 1)
	assert.Equal(t, expected, m.mulTuple(tu))

}

func Test_IdentityMatrixX4(t *testing.T) {
	i := newIdentityMatrixX4()
	m := newX4MatrixWith(1, 2, 3, 4, 5, 6, 7, 8, -1, -2, -3, -4, -5, -6, -7, -5)

	n := m.mulX4Matrix(i)

	assert.True(t, m.equals(n))
	assert.True(t, n.equals(m))
}

func Test_TransposeX4(t *testing.T) {
	m := newX4MatrixWith(0, 9, 3, 0, 9, 8, 0, 8, 1, 8, 5, 3, 0, 0, 5, 8)

	n := m.transpose()

	expected := newX4MatrixWith(0, 9, 1, 0, 9, 8, 8, 0, 3, 0, 5, 5, 0, 8, 3, 8)
	assert.Equal(t, expected, n)
}

func Test_TransposeX4Identity(t *testing.T) {
	m := newIdentityMatrixX4()

	n := m.transpose()

	assert.Equal(t, newIdentityMatrixX4(), n)
}

func Test_X4_Submatrix(t *testing.T) {
	m := newX4MatrixWith(-6, 1, 1, 6, -8, 5, 8, 6, -1, 0, 8, 2, -7, 1, -1, 1)

	n := m.submatrix(2, 1)

	expected := newX3MatrixWith(-6, 1, 6, -8, 8, 6, -7, -1, 1)
	assert.Equal(t, expected, n)
}

func Test_X4_Determinant(t *testing.T) {
	m := newX4MatrixWith(-2, -8, 3, 5, -3, 1, 7, 3, 1, 2, -9, 6, -6, 7, 7, -9)

	assert.Equal(t, 690.0, m.cofactor(0, 0))
	assert.Equal(t, 447.0, m.cofactor(0, 1))
	assert.Equal(t, 210.0, m.cofactor(0, 2))
	assert.Equal(t, 51.0, m.cofactor(0, 3))
	assert.Equal(t, -4071.0, m.determinant())
}

func Test_X4_Invertable(t *testing.T) {
	m := newX4MatrixWith(6, 4, 4, 4, 5, 5, 7, 6, 4, -9, 3, -7, 9, 1, 7, -6)

	assert.Equal(t, -2120.0, m.determinant())
	assert.True(t, m.invertable())
}

func Test_X4_NotInvertable(t *testing.T) {
	m := newX4MatrixWith(-4, 2, -2, -3, 9, 6, 2, 6, 0, -5, 1, -5, 0, 0, 0, 0)

	assert.Equal(t, 0.0, m.determinant())
	assert.False(t, m.invertable())
}

func Test_X4_Invert_1(t *testing.T) {
	m := newX4MatrixWith(-5, 2, 6, -8, 1, -5, 1, 8, 7, 7, -6, -7, 1, -3, 7, 4)

	n := m.invert()
	nr := n.roundTo(5)

	assert.Equal(t, 532.0, m.determinant())
	assert.Equal(t, -160.0, m.cofactor(2, 3))
	assert.Equal(t, -160.0/532, n.get(3, 2))
	assert.Equal(t, 105.0, m.cofactor(3, 2))
	assert.Equal(t, 105.0/532, n.get(2, 3))
	assert.True(t, newX4MatrixWith(
		0.21805, 0.45113, 0.24060, -0.04511,
		-0.80827, -1.45677, -0.44361, 0.52068,
		-0.07895, -0.22368, -0.05263, 0.19737,
		-0.52256, -0.81391, -0.30075, 0.30639,
	).equals(nr))
}

func Test_X4_Invert_2(t *testing.T) {
	m := newX4MatrixWith(8, -5, 9, 2, 7, 5, 6, 1, -6, 0, 9, 6, -3, 0, -9, -4)

	n := m.invert()
	nr := n.roundTo(5)

	assert.True(t, newX4MatrixWith(
		-0.15385, -0.15385, -0.28205, -0.53846,
		-0.07692, 0.12308, 0.02564, 0.03077,
		0.35897, 0.35897, 0.43590, 0.92308,
		-0.69231, -0.69231, -0.76923, -1.92308).equals(nr))
}

func Test_X4_Invert_3(t *testing.T) {
	m := newX4MatrixWith(9, 3, 0, 9, -5, -2, -6, -3, -4, 9, 6, 4, -7, 6, 6, 2)

	n := m.invert()
	nr := n.roundTo(5)

	assert.True(t, newX4MatrixWith(
		-0.04074, -0.07778, 0.14444, -0.22222,
		-0.07778, 0.03333, 0.36667, -0.33333,
		-0.02901, -0.14630, -0.10926, 0.12963,
		0.17778, 0.06667, -0.26667, 0.33333).equals(nr))
}

func Test_X4_MulByInverse(t *testing.T) {
	m := newX4MatrixWith(3, 9, 7, 3, 3, -8, 2, -9, -4, 4, 4, 1, -6, 5, -1, 1)
	n := newX4MatrixWith(8, 2, 2, 2, 3, -1, 7, 0, 7, 0, 5, 4, 6, -2, 0, 5)

	mn := m.mulX4Matrix(n)

	assert.True(t, m.equals(mn.mulX4Matrix(n.invert())))
}
