package main

import (
	"github.com/stretchr/testify/assert"
	"math"
)
import "testing"

func Test_IsPoint(t *testing.T) {
	u := tuple{4.3, -4.2, 3.1, point}

	assert.True(t, u.isPoint())
	assert.False(t, u.isVector())
}

func Test_IsVector(t *testing.T) {
	u := tuple{4.3, -4.2, 3.1, vector}

	assert.False(t, u.isPoint())
	assert.True(t, u.isVector())
}

func Test_NewPoint(t *testing.T) {
	u := newPoint(4, -4, 3)

	assert.Equal(t, tuple{4, -4, 3, point}, u)
	assert.True(t, u.isPoint())
	assert.False(t, u.isVector())
}

func Test_NewVector(t *testing.T) {
	u := newVector(4, -4, 3)

	assert.Equal(t, tuple{4, -4, 3, vector}, u)
	assert.True(t, u.isVector())
	assert.False(t, u.isPoint())
}

func Test_AddTuple_VectorAndVector(t *testing.T) {
	a := newVector(1, 2, 3)
	b := newVector(-1, -1, -1)

	c := a.add(b)

	assert.Equal(t, tuple{0, 1, 2, vector}, c)
	assert.False(t, c.isPoint())
	assert.True(t, c.isVector())
}

func Test_AddTuple_PointAndVector(t *testing.T) {
	a := newPoint(1, 2, 3)
	b := newVector(-1, -1, -1)

	c := a.add(b)

	assert.Equal(t, tuple{0, 1, 2, point}, c)
	assert.True(t, c.isPoint())
	assert.False(t, c.isVector())
}

func Test_AddTuple_VectorAndPoint(t *testing.T) {
	a := newVector(1, 2, 3)
	b := newPoint(-1, -1, -1)

	c := a.add(b)

	assert.Equal(t, tuple{0, 1, 2, point}, c)
	assert.True(t, c.isPoint())
	assert.False(t, c.isVector())
}

func Test_AddTuple_PointAndPoint(t *testing.T) {
	a := newPoint(1, 2, 3)
	b := newPoint(-1, -1, -1)

	c := a.add(b)

	// it's now nonsense
	// ... but the book specifically says to implement it this way...
	assert.False(t, c.isPoint())
	assert.False(t, c.isVector())
	assert.Equal(t, tuple{0, 1, 2, 2}, c)
}

func Test_SubTuple_VectorFromVector(t *testing.T) {
	a := newVector(3, 2, 1)
	b := newVector(5, 6, 7)

	c := a.sub(b)

	assert.Equal(t, tuple{-2, -4, -6, vector}, c)
	assert.False(t, c.isPoint())
	assert.True(t, c.isVector())
}

func Test_SubTuple_VectorFromPoint(t *testing.T) {
	a := newPoint(3, 2, 1)
	b := newVector(5, 6, 7)

	c := a.sub(b)

	assert.Equal(t, tuple{-2, -4, -6, point}, c)
	assert.True(t, c.isPoint())
	assert.False(t, c.isVector())
}

func Test_SubTuple_PointFromVector(t *testing.T) {
	a := newVector(3, 2, 1)
	b := newPoint(5, 6, 7)

	c := a.sub(b)

	// it's now nonsense
	// ... but the book specifically says to implement it this way...
	assert.Equal(t, tuple{-2, -4, -6, -1}, c)
	assert.False(t, c.isPoint())
	assert.False(t, c.isVector())
}

func Test_SubTuple_PointFromPoint(t *testing.T) {
	a := newPoint(3, 2, 1)
	b := newPoint(5, 6, 7)

	c := a.sub(b)

	assert.False(t, c.isPoint())
	assert.True(t, c.isVector())
	assert.Equal(t, tuple{-2, -4, -6, vector}, c)
}

func Test_Sub_VectorFromZV(t *testing.T) {
	zv := tuple{}
	a := newVector(1, -2, 3)

	b := zv.sub(a)

	assert.Equal(t, tuple{-1, 2, -3, vector}, b)
}

func Test_Neg_Tuple(t *testing.T) {
	a := newVector(1, -2, 3)

	b := a.neg()

	assert.Equal(t, tuple{-1, 2, -3, vector}, b)
}

func Test_Mul_Tuple_ByScalar(t *testing.T) {
	a := tuple{1, -2, 3, -4}

	b := a.mul(3.5)

	assert.Equal(t, tuple{3.5, -7, 10.5, -14}, b)
}

func Test_Mul_Tuple_ByFraction(t *testing.T) {
	a := tuple{1, -2, 3, -4}

	b := a.mul(0.5)

	assert.Equal(t, tuple{0.5, -1, 1.5, -2}, b)
}

func Test_Div_Tuple(t *testing.T) {
	a := tuple{1, -2, 3, -4}

	b := a.div(2)

	assert.Equal(t, tuple{0.5, -1, 1.5, -2}, b)
}

func Test_Mag_UnitVectorX(t *testing.T) {
	a := newVector(1, 0, 0)

	mag := a.mag()

	assert.Equal(t, 1.0, mag)
}

func Test_Mag_UnitVectorY(t *testing.T) {
	a := newVector(0, 1, 0)

	mag := a.mag()

	assert.Equal(t, 1.0, mag)
}

func Test_Mag_UnitVectorZ(t *testing.T) {
	a := newVector(0, 0, 1)

	mag := a.mag()

	assert.Equal(t, 1.0, mag)
}

func Test_Mag_Vector(t *testing.T) {
	a := newVector(1, 2, 3)

	mag := a.mag()

	assert.Equal(t, math.Sqrt(14), mag)
}

func Test_Mag_NegVector(t *testing.T) {
	a := newVector(-1, -2, -3)

	mag := a.mag()

	assert.Equal(t, math.Sqrt(14), mag)
}

func Test_Normalize_1D(t *testing.T) {
	a := newVector(4, 0, 0)

	b := a.normalize()

	assert.Equal(t, tuple{1, 0, 0, vector}, b)
}

func Test_Normalize_3D(t *testing.T) {
	a := newVector(1, 2, 3)

	b := a.normalize()

	assert.Equal(t, tuple{1 / math.Sqrt(14), 2 / math.Sqrt(14), 3 / math.Sqrt(14), vector}, b)
}

func Test_Normalized_Magnitude(t *testing.T) {
	a := newVector(1, 2, 3)

	b := a.normalize()

	assert.Equal(t, 1.0, b.mag())
}

func Test_DotProduct(t *testing.T) {
	a := newVector(1, 2, 3)
	b := newVector(2, 3, 4)

	c := a.dot(b)
	d := b.dot(a)

	assert.Equal(t, 20.0, c)
	assert.Equal(t, c, d)
}

func Test_CrossProductA(t *testing.T) {
	a := newVector(1, 2, 3)
	b := newVector(2, 3, 4)

	c := cross(a, b)

	assert.Equal(t, newVector(-1, 2, -1), c)
}

func Test_CrossProductB(t *testing.T) {
	a := newVector(1, 2, 3)
	b := newVector(2, 3, 4)

	c := cross(b, a)

	assert.Equal(t, newVector(1, -2, 1), c)
}
