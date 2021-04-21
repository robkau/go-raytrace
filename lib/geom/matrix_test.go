package geom

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_New2x2(t *testing.T) {
	m := NewX2Matrix()

	assert.Len(t, m.b, 4)
}

func Test_GetSet2x2(t *testing.T) {
	m := NewX2Matrix()
	m.Set(0, 0, 1.5)
	m.Set(0, 1, 2.7)
	m.Set(1, 0, 3.9)
	m.Set(1, 1, -4.1)

	assert.Equal(t, 1.5, m.Get(0, 0))
	assert.Equal(t, 2.7, m.Get(0, 1))
	assert.Equal(t, 3.9, m.Get(1, 0))
	assert.Equal(t, -4.1, m.Get(1, 1))
}

func Test_Equals2x2(t *testing.T) {
	m := NewX2MatrixWith(0, 1, 2, 3)
	n := NewX2MatrixWith(0, 1, 2, 3)
	o := NewX2MatrixWith(0, 1.00000001, 2, 3)

	assert.True(t, m.Equals(m))
	assert.True(t, n.Equals(n))
	assert.True(t, m.Equals(n))
	assert.True(t, n.Equals(m))
	assert.True(t, !m.Equals(o))
	assert.True(t, !n.Equals(o))
	assert.True(t, !o.Equals(m))
	assert.True(t, !o.Equals(n))
}

func Test_DeterminantX2(t *testing.T) {
	m := NewX2MatrixWith(
		1, 5,
		-3, 2)

	assert.Equal(t, 17.0, m.Determinant())
}

func Test_New3x3(t *testing.T) {
	m := NewX3Matrix()

	assert.Len(t, m.b, 9)
}

func Test_GetSet3x3(t *testing.T) {
	m := NewX3Matrix()
	m.Set(0, 0, 1.5)
	m.Set(0, 1, 2.7)
	m.Set(0, 2, 2.9)
	m.Set(1, 0, 3.9)
	m.Set(1, 1, -4.1)
	m.Set(1, 2, -4.3)
	m.Set(2, 0, 3.12)
	m.Set(2, 1, -4.122)
	m.Set(2, 2, -4.3333)

	assert.Equal(t, 1.5, m.Get(0, 0))
	assert.Equal(t, 2.7, m.Get(0, 1))
	assert.Equal(t, 2.9, m.Get(0, 2))
	assert.Equal(t, 3.9, m.Get(1, 0))
	assert.Equal(t, -4.1, m.Get(1, 1))
	assert.Equal(t, -4.3, m.Get(1, 2))
	assert.Equal(t, 3.12, m.Get(2, 0))
	assert.Equal(t, -4.122, m.Get(2, 1))
	assert.Equal(t, -4.3333, m.Get(2, 2))
}

func Test_Equals3x3(t *testing.T) {
	m := NewX3MatrixWith(0, 1, 2, 3, 4, 5, 6, 7, 8)
	n := NewX3MatrixWith(0, 1, 2, 3, 4, 5, 6, 7, 8)
	o := NewX3MatrixWith(0, 1.00000001, 2, 3, 4, 5, 6, 7, 8)

	assert.True(t, m.Equals(m))
	assert.True(t, n.Equals(n))
	assert.True(t, m.Equals(n))
	assert.True(t, n.Equals(m))
	assert.True(t, !m.Equals(o))
	assert.True(t, !n.Equals(o))
	assert.True(t, !o.Equals(m))
	assert.True(t, !o.Equals(n))
}

func Test_X3_Minor(t *testing.T) {
	m := NewX3MatrixWith(3, 5, 0, 2, -1, -7, 6, -1, 5)

	b := m.Submatrix(1, 0)

	assert.Equal(t, 25.0, b.Determinant())
	assert.Equal(t, 25.0, m.Minor(1, 0))
}

func Test_X3_Cofactor(t *testing.T) {
	m := NewX3MatrixWith(3, 5, 0, 2, -1, -7, 6, -1, 5)

	assert.Equal(t, -12.0, m.Minor(0, 0))
	assert.Equal(t, -12.0, m.Cofactor(0, 0))
	assert.Equal(t, 25.0, m.Minor(1, 0))
	assert.Equal(t, -25.0, m.Cofactor(1, 0))
}

func Test_X3_Submatrix(t *testing.T) {
	m := NewX3MatrixWith(1, 5, 0, -3, 2, 7, 0, 6, -3)

	n := m.Submatrix(0, 2)

	expected := NewX2MatrixWith(-3, 2, 0, 6)
	assert.Equal(t, expected, n)
}

func Test_X3_Determinant(t *testing.T) {
	m := NewX3MatrixWith(1, 2, 6, -5, 8, -4, 2, 6, 4)

	assert.Equal(t, 56.0, m.Cofactor(0, 0))
	assert.Equal(t, 12.0, m.Cofactor(0, 1))
	assert.Equal(t, -46.0, m.Cofactor(0, 2))
	assert.Equal(t, -196.0, m.Determinant())
}

func Test_New4x4(t *testing.T) {
	m := NewX4Matrix()

	assert.Len(t, m.b, 16)
}

func Test_GetSet4x4(t *testing.T) {
	m := NewX4Matrix()
	m.Set(0, 0, 1.5)
	m.Set(0, 1, 2.7)
	m.Set(0, 2, 2.9)
	m.Set(0, 3, 2.915)
	m.Set(1, 0, 3.9)
	m.Set(1, 1, -4.1)
	m.Set(1, 2, -4.3)
	m.Set(1, 3, -8.3)
	m.Set(2, 0, 3.12)
	m.Set(2, 1, -4.122)
	m.Set(2, 2, -4.3333)
	m.Set(2, 3, 4.3333)
	m.Set(3, 0, 1.1)
	m.Set(3, 1, 1.33)
	m.Set(3, 2, 1.66)
	m.Set(3, 3, 0.02)

	assert.Equal(t, 1.5, m.Get(0, 0))
	assert.Equal(t, 2.7, m.Get(0, 1))
	assert.Equal(t, 2.9, m.Get(0, 2))
	assert.Equal(t, 2.915, m.Get(0, 3))
	assert.Equal(t, 3.9, m.Get(1, 0))
	assert.Equal(t, -4.1, m.Get(1, 1))
	assert.Equal(t, -4.3, m.Get(1, 2))
	assert.Equal(t, -8.3, m.Get(1, 3))
	assert.Equal(t, 3.12, m.Get(2, 0))
	assert.Equal(t, -4.122, m.Get(2, 1))
	assert.Equal(t, -4.3333, m.Get(2, 2))
	assert.Equal(t, 4.3333, m.Get(2, 3))
	assert.Equal(t, 1.1, m.Get(3, 0))
	assert.Equal(t, 1.33, m.Get(3, 1))
	assert.Equal(t, 1.66, m.Get(3, 2))
	assert.Equal(t, 0.02, m.Get(3, 3))
}

func Test_Equals4x4(t *testing.T) {
	m := NewX4MatrixWith(0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15)
	n := NewX4MatrixWith(0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15)
	o := NewX4MatrixWith(0, 1.00000001, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15)

	assert.True(t, m.Equals(m))
	assert.True(t, n.Equals(n))
	assert.True(t, m.Equals(n))
	assert.True(t, n.Equals(m))
	assert.True(t, !m.Equals(o))
	assert.True(t, !n.Equals(o))
	assert.True(t, !o.Equals(m))
	assert.True(t, !o.Equals(n))
}

func Test_MulX4Matrix(t *testing.T) {
	m := NewX4MatrixWith(1, 2, 3, 4, 5, 6, 7, 8, 9, 8, 7, 6, 5, 4, 3, 2)
	n := NewX4MatrixWith(-2, 1, 2, 3, 3, 2, 1, -1, 4, 3, 6, 5, 1, 2, 7, 8)

	o := m.MulX4Matrix(n)

	expected := NewX4MatrixWith(20, 22, 50, 48, 44, 54, 114, 108, 40, 58, 110, 102, 16, 26, 46, 42)
	assert.True(t, o.Equals(expected))
}

func Test_MulX4Tuple(t *testing.T) {
	m := NewX4MatrixWith(1, 2, 3, 4, 2, 4, 4, 2, 8, 6, 4, 1, 0, 0, 0, 1)
	tu := NewTuple(1, 2, 3, 1)

	expected := NewTuple(18, 24, 33, 1)
	assert.Equal(t, expected, m.MulTuple(tu))

}

func Test_IdentityMatrixX4(t *testing.T) {
	i := NewIdentityMatrixX4()
	m := NewX4MatrixWith(1, 2, 3, 4, 5, 6, 7, 8, -1, -2, -3, -4, -5, -6, -7, -5)

	n := m.MulX4Matrix(i)

	assert.True(t, m.Equals(n))
	assert.True(t, n.Equals(m))
}

func Test_TransposeX4(t *testing.T) {
	m := NewX4MatrixWith(0, 9, 3, 0, 9, 8, 0, 8, 1, 8, 5, 3, 0, 0, 5, 8)

	n := m.Transpose()

	expected := NewX4MatrixWith(0, 9, 1, 0, 9, 8, 8, 0, 3, 0, 5, 5, 0, 8, 3, 8)
	assert.Equal(t, expected, n)
}

func Test_TransposeX4Identity(t *testing.T) {
	m := NewIdentityMatrixX4()

	n := m.Transpose()

	assert.Equal(t, NewIdentityMatrixX4(), n)
}

func Test_X4_Submatrix(t *testing.T) {
	m := NewX4MatrixWith(-6, 1, 1, 6, -8, 5, 8, 6, -1, 0, 8, 2, -7, 1, -1, 1)

	n := m.Submatrix(2, 1)

	expected := NewX3MatrixWith(-6, 1, 6, -8, 8, 6, -7, -1, 1)
	assert.Equal(t, expected, n)
}

func Test_X4_Determinant(t *testing.T) {
	m := NewX4MatrixWith(-2, -8, 3, 5, -3, 1, 7, 3, 1, 2, -9, 6, -6, 7, 7, -9)

	assert.Equal(t, 690.0, m.Cofactor(0, 0))
	assert.Equal(t, 447.0, m.Cofactor(0, 1))
	assert.Equal(t, 210.0, m.Cofactor(0, 2))
	assert.Equal(t, 51.0, m.Cofactor(0, 3))
	assert.Equal(t, -4071.0, m.Determinant())
}

func Test_X4_Invertable(t *testing.T) {
	m := NewX4MatrixWith(6, 4, 4, 4, 5, 5, 7, 6, 4, -9, 3, -7, 9, 1, 7, -6)

	assert.Equal(t, -2120.0, m.Determinant())
	assert.True(t, m.Invertable())
}

func Test_X4_NotInvertable(t *testing.T) {
	m := NewX4MatrixWith(-4, 2, -2, -3, 9, 6, 2, 6, 0, -5, 1, -5, 0, 0, 0, 0)

	assert.Equal(t, 0.0, m.Determinant())
	assert.False(t, m.Invertable())
}

func Test_X4_Invert_1(t *testing.T) {
	m := NewX4MatrixWith(-5, 2, 6, -8, 1, -5, 1, 8, 7, 7, -6, -7, 1, -3, 7, 4)

	n := m.Invert()
	nr := n.RoundTo(5)

	assert.Equal(t, 532.0, m.Determinant())
	assert.Equal(t, -160.0, m.Cofactor(2, 3))
	assert.Equal(t, -160.0/532, n.Get(3, 2))
	assert.Equal(t, 105.0, m.Cofactor(3, 2))
	assert.Equal(t, 105.0/532, n.Get(2, 3))
	assert.True(t, NewX4MatrixWith(
		0.21805, 0.45113, 0.24060, -0.04511,
		-0.80827, -1.45677, -0.44361, 0.52068,
		-0.07895, -0.22368, -0.05263, 0.19737,
		-0.52256, -0.81391, -0.30075, 0.30639,
	).Equals(nr))
}

func Test_X4_Invert_2(t *testing.T) {
	m := NewX4MatrixWith(8, -5, 9, 2, 7, 5, 6, 1, -6, 0, 9, 6, -3, 0, -9, -4)

	n := m.Invert()
	nr := n.RoundTo(5)

	assert.True(t, NewX4MatrixWith(
		-0.15385, -0.15385, -0.28205, -0.53846,
		-0.07692, 0.12308, 0.02564, 0.03077,
		0.35897, 0.35897, 0.43590, 0.92308,
		-0.69231, -0.69231, -0.76923, -1.92308).Equals(nr))
}

func Test_X4_Invert_3(t *testing.T) {
	m := NewX4MatrixWith(9, 3, 0, 9, -5, -2, -6, -3, -4, 9, 6, 4, -7, 6, 6, 2)

	n := m.Invert()
	nr := n.RoundTo(5)

	assert.True(t, NewX4MatrixWith(
		-0.04074, -0.07778, 0.14444, -0.22222,
		-0.07778, 0.03333, 0.36667, -0.33333,
		-0.02901, -0.14630, -0.10926, 0.12963,
		0.17778, 0.06667, -0.26667, 0.33333).Equals(nr))
}

func Test_X4_MulByInverse(t *testing.T) {
	m := NewX4MatrixWith(3, 9, 7, 3, 3, -8, 2, -9, -4, 4, 4, 1, -6, 5, -1, 1)
	n := NewX4MatrixWith(8, 2, 2, 2, 3, -1, 7, 0, 7, 0, 5, 4, 6, -2, 0, 5)

	mn := m.MulX4Matrix(n)

	assert.True(t, m.Equals(mn.MulX4Matrix(n.Invert())))
}
