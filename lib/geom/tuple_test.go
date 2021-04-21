package geom

import (
	"github.com/stretchr/testify/assert"
	"math"
)
import "testing"

func Test_IsPoint(t *testing.T) {
	u := Tuple{4.3, -4.2, 3.1, point}

	assert.True(t, u.IsPoint())
	assert.False(t, u.IsVector())
}

func Test_IsVector(t *testing.T) {
	u := Tuple{4.3, -4.2, 3.1, Vector}

	assert.False(t, u.IsPoint())
	assert.True(t, u.IsVector())
}

func Test_NewPoint(t *testing.T) {
	u := NewPoint(4, -4, 3)

	assert.Equal(t, Tuple{4, -4, 3, point}, u)
	assert.True(t, u.IsPoint())
	assert.False(t, u.IsVector())
}

func Test_NewVector(t *testing.T) {
	u := NewVector(4, -4, 3)

	assert.Equal(t, Tuple{4, -4, 3, Vector}, u)
	assert.True(t, u.IsVector())
	assert.False(t, u.IsPoint())
}

func Test_Equals(t *testing.T) {
	a := NewTuple(1, 2, 3, 4)
	b := NewTuple(1, 2, 3, 4)
	c := NewTuple(1.000001, 2, 3, 4)

	assert.True(t, a.Equals(a))
	assert.True(t, b.Equals(b))
	assert.True(t, c.Equals(c))

	assert.True(t, a.Equals(b))
	assert.True(t, b.Equals(a))

	assert.False(t, a.Equals(c))
	assert.False(t, b.Equals(c))
	assert.False(t, c.Equals(a))
	assert.False(t, c.Equals(b))
}

func Test_AddTuple_VectorAndVector(t *testing.T) {
	a := NewVector(1, 2, 3)
	b := NewVector(-1, -1, -1)

	c := a.Add(b)

	assert.Equal(t, Tuple{0, 1, 2, Vector}, c)
	assert.False(t, c.IsPoint())
	assert.True(t, c.IsVector())
}

func Test_AddTuple_PointAndVector(t *testing.T) {
	a := NewPoint(1, 2, 3)
	b := NewVector(-1, -1, -1)

	c := a.Add(b)

	assert.Equal(t, Tuple{0, 1, 2, point}, c)
	assert.True(t, c.IsPoint())
	assert.False(t, c.IsVector())
}

func Test_AddTuple_VectorAndPoint(t *testing.T) {
	a := NewVector(1, 2, 3)
	b := NewPoint(-1, -1, -1)

	c := a.Add(b)

	assert.Equal(t, Tuple{0, 1, 2, point}, c)
	assert.True(t, c.IsPoint())
	assert.False(t, c.IsVector())
}

func Test_AddTuple_PointAndPoint(t *testing.T) {
	a := NewPoint(1, 2, 3)
	b := NewPoint(-1, -1, -1)

	c := a.Add(b)

	// it's now nonsense
	// ... but the book specifically says to implement it this way...
	assert.False(t, c.IsPoint())
	assert.False(t, c.IsVector())
	assert.Equal(t, Tuple{0, 1, 2, 2}, c)
}

func Test_SubTuple_VectorFromVector(t *testing.T) {
	a := NewVector(3, 2, 1)
	b := NewVector(5, 6, 7)

	c := a.Sub(b)

	assert.Equal(t, Tuple{-2, -4, -6, Vector}, c)
	assert.False(t, c.IsPoint())
	assert.True(t, c.IsVector())
}

func Test_SubTuple_VectorFromPoint(t *testing.T) {
	a := NewPoint(3, 2, 1)
	b := NewVector(5, 6, 7)

	c := a.Sub(b)

	assert.Equal(t, Tuple{-2, -4, -6, point}, c)
	assert.True(t, c.IsPoint())
	assert.False(t, c.IsVector())
}

func Test_SubTuple_PointFromVector(t *testing.T) {
	a := NewVector(3, 2, 1)
	b := NewPoint(5, 6, 7)

	c := a.Sub(b)

	// it's now nonsense
	// ... but the book specifically says to implement it this way...
	assert.Equal(t, Tuple{-2, -4, -6, -1}, c)
	assert.False(t, c.IsPoint())
	assert.False(t, c.IsVector())
}

func Test_SubTuple_PointFromPoint(t *testing.T) {
	a := NewPoint(3, 2, 1)
	b := NewPoint(5, 6, 7)

	c := a.Sub(b)

	assert.False(t, c.IsPoint())
	assert.True(t, c.IsVector())
	assert.Equal(t, Tuple{-2, -4, -6, Vector}, c)
}

func Test_Sub_VectorFromZV(t *testing.T) {
	zv := Tuple{}
	a := NewVector(1, -2, 3)

	b := zv.Sub(a)

	assert.Equal(t, Tuple{-1, 2, -3, Vector}, b)
}

func Test_Neg_Tuple(t *testing.T) {
	a := NewVector(1, -2, 3)

	b := a.Neg()

	assert.Equal(t, Tuple{-1, 2, -3, Vector}, b)
}

func Test_Mul_Tuple_ByScalar(t *testing.T) {
	a := Tuple{1, -2, 3, -4}

	b := a.Mul(3.5)

	assert.Equal(t, Tuple{3.5, -7, 10.5, -14}, b)
}

func Test_Mul_Tuple_ByFraction(t *testing.T) {
	a := Tuple{1, -2, 3, -4}

	b := a.Mul(0.5)

	assert.Equal(t, Tuple{0.5, -1, 1.5, -2}, b)
}

func Test_Div_Tuple(t *testing.T) {
	a := Tuple{1, -2, 3, -4}

	b := a.Div(2)

	assert.Equal(t, Tuple{0.5, -1, 1.5, -2}, b)
}

func Test_Mag_UnitVectorX(t *testing.T) {
	a := NewVector(1, 0, 0)

	mag := a.Mag()

	assert.Equal(t, 1.0, mag)
}

func Test_Mag_UnitVectorY(t *testing.T) {
	a := NewVector(0, 1, 0)

	mag := a.Mag()

	assert.Equal(t, 1.0, mag)
}

func Test_Mag_UnitVectorZ(t *testing.T) {
	a := NewVector(0, 0, 1)

	mag := a.Mag()

	assert.Equal(t, 1.0, mag)
}

func Test_Mag_Vector(t *testing.T) {
	a := NewVector(1, 2, 3)

	mag := a.Mag()

	assert.Equal(t, math.Sqrt(14), mag)
}

func Test_Mag_NegVector(t *testing.T) {
	a := NewVector(-1, -2, -3)

	mag := a.Mag()

	assert.Equal(t, math.Sqrt(14), mag)
}

func Test_Normalize_1D(t *testing.T) {
	a := NewVector(4, 0, 0)

	b := a.Normalize()

	assert.Equal(t, Tuple{1, 0, 0, Vector}, b)
}

func Test_Normalize_3D(t *testing.T) {
	a := NewVector(1, 2, 3)

	b := a.Normalize()

	assert.Equal(t, Tuple{1 / math.Sqrt(14), 2 / math.Sqrt(14), 3 / math.Sqrt(14), Vector}, b)
}

func Test_Normalized_Magnitude(t *testing.T) {
	a := NewVector(1, 2, 3)

	b := a.Normalize()

	assert.Equal(t, 1.0, b.Mag())
}

func Test_DotProduct(t *testing.T) {
	a := NewVector(1, 2, 3)
	b := NewVector(2, 3, 4)

	c := a.Dot(b)
	d := b.Dot(a)

	assert.Equal(t, 20.0, c)
	assert.Equal(t, c, d)
}

func Test_CrossProductA(t *testing.T) {
	a := NewVector(1, 2, 3)
	b := NewVector(2, 3, 4)

	c := Cross(a, b)

	assert.Equal(t, NewVector(-1, 2, -1), c)
}

func Test_CrossProductB(t *testing.T) {
	a := NewVector(1, 2, 3)
	b := NewVector(2, 3, 4)

	c := Cross(b, a)

	assert.Equal(t, NewVector(1, -2, 1), c)
}

func Test_Reflect45(t *testing.T) {
	v := NewVector(1, -1, 0)
	n := NewVector(0, 1, 0)

	r := v.Reflect(n)

	assert.Equal(t, NewVector(1, 1, 0), r)
}

func Test_ReflectSlanted(t *testing.T) {
	v := NewVector(0, -1, 0)
	n := NewVector(math.Sqrt2/2, math.Sqrt2/2, 0)

	r := v.Reflect(n)

	assert.True(t, NewVector(1, 0, 0).Equals(r))
}
