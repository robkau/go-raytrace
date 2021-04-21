package geom

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func Test_Translate(t *testing.T) {
	tr := Translate(5, -3, 2)
	p := NewPoint(-3, 4, 5)

	pTr := tr.MulTuple(p)

	assert.True(t, NewPoint(2, 1, 7).Equals(pTr))
}

func Test_Translate_Inverse(t *testing.T) {
	tr := Translate(5, -3, 2)
	iTr := tr.Invert()
	p := NewPoint(-3, 4, 5)

	pTr := iTr.MulTuple(p)

	assert.True(t, NewPoint(-8, 7, 3).Equals(pTr))
}

func Test_Translate_Noop_Vector(t *testing.T) {
	tr := Translate(5, -3, 2)
	v := NewVector(-3, 4, 5)

	pTr := tr.MulTuple(v)

	assert.True(t, v.Equals(pTr))
}

func Test_Scale_Point(t *testing.T) {
	p := NewPoint(-4, 6, 8)
	s := Scale(2, 3, 4)

	ps := s.MulTuple(p)

	assert.Equal(t, NewPoint(-8, 18, 32), ps)
}

func Test_Scale_Vector(t *testing.T) {
	p := NewVector(-4, 6, 8)
	s := Scale(2, 3, 4)

	ps := s.MulTuple(p)

	assert.Equal(t, NewVector(-8, 18, 32), ps)
}

func Test_Scale_Inverse(t *testing.T) {
	p := NewVector(-4, 6, 8)
	s := Scale(2, 3, 4)
	si := s.Invert()

	psi := si.MulTuple(p)

	assert.Equal(t, NewVector(-2, 2, 2), psi)
}

func Test_Scale_Reflect(t *testing.T) {
	p := NewPoint(2, 3, 4)
	s := Scale(-1, 1, 1)

	ps := s.MulTuple(p)

	assert.Equal(t, NewPoint(-2, 3, 4), ps)
}

func Test_Rotate_X(t *testing.T) {
	p := NewPoint(0, 1, 0)
	halfQuarter := RotateX(math.Pi / 4)
	fullQuarter := RotateX(math.Pi / 2)

	hp := halfQuarter.MulTuple(p)
	fp := fullQuarter.MulTuple(p)

	assert.True(t, NewPoint(0, math.Sqrt2/2, math.Sqrt2/2).Equals(hp))
	assert.True(t, NewPoint(0, 0, 1).Equals(fp))
}

func Test_Rotate_X_Inverse(t *testing.T) {
	p := NewPoint(0, 1, 0)
	halfQuarter := RotateX(math.Pi / 4)
	hqi := halfQuarter.Invert()

	hpi := hqi.MulTuple(p)

	assert.True(t, NewPoint(0, math.Sqrt2/2, -math.Sqrt2/2).Equals(hpi))
}

func Test_Rotate_Y(t *testing.T) {
	p := NewPoint(0, 0, 1)
	halfQuarter := RotateY(math.Pi / 4)
	fullQuarter := RotateY(math.Pi / 2)

	hp := halfQuarter.MulTuple(p)
	fp := fullQuarter.MulTuple(p)

	assert.True(t, NewPoint(math.Sqrt2/2, 0, math.Sqrt2/2).Equals(hp))
	assert.True(t, NewPoint(1, 0, 0).Equals(fp))
}

func Test_Rotate_Z(t *testing.T) {
	p := NewPoint(0, 1, 0)
	halfQuarter := RotateZ(math.Pi / 4)
	fullQuarter := RotateZ(math.Pi / 2)

	hp := halfQuarter.MulTuple(p)
	fp := fullQuarter.MulTuple(p)

	assert.True(t, NewPoint(-math.Sqrt2/2, math.Sqrt2/2, 0).Equals(hp))
	assert.True(t, NewPoint(-1, 0, 0).Equals(fp))
}

func Test_ShearXY(t *testing.T) {
	p := NewPoint(2, 3, 4)
	s := shear(1, 0, 0, 0, 0, 0)

	assert.True(t, NewPoint(5, 3, 4).Equals(s.MulTuple(p)))
}

func Test_ShearXZ(t *testing.T) {
	p := NewPoint(2, 3, 4)
	s := shear(0, 1, 0, 0, 0, 0)

	assert.True(t, NewPoint(6, 3, 4).Equals(s.MulTuple(p)))
}

func Test_ShearYX(t *testing.T) {
	p := NewPoint(2, 3, 4)
	s := shear(0, 0, 1, 0, 0, 0)

	assert.True(t, NewPoint(2, 5, 4).Equals(s.MulTuple(p)))
}

func Test_ShearYZ(t *testing.T) {
	p := NewPoint(2, 3, 4)
	s := shear(0, 0, 0, 1, 0, 0)

	assert.True(t, NewPoint(2, 7, 4).Equals(s.MulTuple(p)))
}

func Test_ShearZX(t *testing.T) {
	p := NewPoint(2, 3, 4)
	s := shear(0, 0, 0, 0, 1, 0)

	assert.True(t, NewPoint(2, 3, 6).Equals(s.MulTuple(p)))
}

func Test_ShearZY(t *testing.T) {
	p := NewPoint(2, 3, 4)
	s := shear(0, 0, 0, 0, 0, 1)

	assert.True(t, NewPoint(2, 3, 7).Equals(s.MulTuple(p)))
}

func Test_Individual_Transforms(t *testing.T) {
	p := NewPoint(1, 0, 1)
	a := RotateX(math.Pi / 2)
	b := Scale(5, 5, 5)
	c := Translate(10, 5, 7)

	p2 := a.MulTuple(p)
	p3 := b.MulTuple(p2)
	p4 := c.MulTuple(p3)

	assert.True(t, NewPoint(1, -1, 0).Equals(p2))
	assert.True(t, NewPoint(5, -5, 0).Equals(p3))
	assert.True(t, NewPoint(15, 0, 7).Equals(p4))
}

func Test_Chained_Transforms(t *testing.T) {
	p := NewPoint(1, 0, 1)
	a := RotateX(math.Pi / 2)
	b := Scale(5, 5, 5)
	c := Translate(10, 5, 7)

	d := c.MulX4Matrix(b).MulX4Matrix(a)
	p2 := d.MulTuple(p)

	assert.True(t, NewPoint(15, 0, 7).Equals(p2))
}
