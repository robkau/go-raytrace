package main

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func Test_Translate(t *testing.T) {
	tr := translate(5, -3, 2)
	p := newPoint(-3, 4, 5)

	pTr := tr.mulTuple(p)

	assert.True(t, newPoint(2, 1, 7).equals(pTr))
}

func Test_Translate_Inverse(t *testing.T) {
	tr := translate(5, -3, 2)
	iTr := tr.invert()
	p := newPoint(-3, 4, 5)

	pTr := iTr.mulTuple(p)

	assert.True(t, newPoint(-8, 7, 3).equals(pTr))
}

func Test_Translate_Noop_Vector(t *testing.T) {
	tr := translate(5, -3, 2)
	v := newVector(-3, 4, 5)

	pTr := tr.mulTuple(v)

	assert.True(t, v.equals(pTr))
}

func Test_Scale_Point(t *testing.T) {
	p := newPoint(-4, 6, 8)
	s := scale(2, 3, 4)

	ps := s.mulTuple(p)

	assert.Equal(t, newPoint(-8, 18, 32), ps)
}

func Test_Scale_Vector(t *testing.T) {
	p := newVector(-4, 6, 8)
	s := scale(2, 3, 4)

	ps := s.mulTuple(p)

	assert.Equal(t, newVector(-8, 18, 32), ps)
}

func Test_Scale_Inverse(t *testing.T) {
	p := newVector(-4, 6, 8)
	s := scale(2, 3, 4)
	si := s.invert()

	psi := si.mulTuple(p)

	assert.Equal(t, newVector(-2, 2, 2), psi)
}

func Test_Scale_Reflect(t *testing.T) {
	p := newPoint(2, 3, 4)
	s := scale(-1, 1, 1)

	ps := s.mulTuple(p)

	assert.Equal(t, newPoint(-2, 3, 4), ps)
}

func Test_Rotate_X(t *testing.T) {
	p := newPoint(0, 1, 0)
	halfQuarter := rotateX(math.Pi / 4)
	fullQuarter := rotateX(math.Pi / 2)

	hp := halfQuarter.mulTuple(p)
	fp := fullQuarter.mulTuple(p)

	assert.True(t, newPoint(0, math.Sqrt2/2, math.Sqrt2/2).equals(hp))
	assert.True(t, newPoint(0, 0, 1).equals(fp))
}

func Test_Rotate_X_Inverse(t *testing.T) {
	p := newPoint(0, 1, 0)
	halfQuarter := rotateX(math.Pi / 4)
	hqi := halfQuarter.invert()

	hpi := hqi.mulTuple(p)

	assert.True(t, newPoint(0, math.Sqrt2/2, -math.Sqrt2/2).equals(hpi))
}

func Test_Rotate_Y(t *testing.T) {
	p := newPoint(0, 0, 1)
	halfQuarter := rotateY(math.Pi / 4)
	fullQuarter := rotateY(math.Pi / 2)

	hp := halfQuarter.mulTuple(p)
	fp := fullQuarter.mulTuple(p)

	assert.True(t, newPoint(math.Sqrt2/2, 0, math.Sqrt2/2).equals(hp))
	assert.True(t, newPoint(1, 0, 0).equals(fp))
}

func Test_Rotate_Z(t *testing.T) {
	p := newPoint(0, 1, 0)
	halfQuarter := rotateZ(math.Pi / 4)
	fullQuarter := rotateZ(math.Pi / 2)

	hp := halfQuarter.mulTuple(p)
	fp := fullQuarter.mulTuple(p)

	assert.True(t, newPoint(-math.Sqrt2/2, math.Sqrt2/2, 0).equals(hp))
	assert.True(t, newPoint(-1, 0, 0).equals(fp))
}

func Test_ShearXY(t *testing.T) {
	p := newPoint(2, 3, 4)
	s := shear(1, 0, 0, 0, 0, 0)

	assert.True(t, newPoint(5, 3, 4).equals(s.mulTuple(p)))
}

func Test_ShearXZ(t *testing.T) {
	p := newPoint(2, 3, 4)
	s := shear(0, 1, 0, 0, 0, 0)

	assert.True(t, newPoint(6, 3, 4).equals(s.mulTuple(p)))
}

func Test_ShearYX(t *testing.T) {
	p := newPoint(2, 3, 4)
	s := shear(0, 0, 1, 0, 0, 0)

	assert.True(t, newPoint(2, 5, 4).equals(s.mulTuple(p)))
}

func Test_ShearYZ(t *testing.T) {
	p := newPoint(2, 3, 4)
	s := shear(0, 0, 0, 1, 0, 0)

	assert.True(t, newPoint(2, 7, 4).equals(s.mulTuple(p)))
}

func Test_ShearZX(t *testing.T) {
	p := newPoint(2, 3, 4)
	s := shear(0, 0, 0, 0, 1, 0)

	assert.True(t, newPoint(2, 3, 6).equals(s.mulTuple(p)))
}

func Test_ShearZY(t *testing.T) {
	p := newPoint(2, 3, 4)
	s := shear(0, 0, 0, 0, 0, 1)

	assert.True(t, newPoint(2, 3, 7).equals(s.mulTuple(p)))
}

func Test_Individual_Transforms(t *testing.T) {
	p := newPoint(1, 0, 1)
	a := rotateX(math.Pi / 2)
	b := scale(5, 5, 5)
	c := translate(10, 5, 7)

	p2 := a.mulTuple(p)
	p3 := b.mulTuple(p2)
	p4 := c.mulTuple(p3)

	assert.True(t, newPoint(1, -1, 0).equals(p2))
	assert.True(t, newPoint(5, -5, 0).equals(p3))
	assert.True(t, newPoint(15, 0, 7).equals(p4))
}

func Test_Chained_Transforms(t *testing.T) {
	p := newPoint(1, 0, 1)
	a := rotateX(math.Pi / 2)
	b := scale(5, 5, 5)
	c := translate(10, 5, 7)

	d := c.mulX4Matrix(b).mulX4Matrix(a)
	p2 := d.mulTuple(p)

	assert.True(t, newPoint(15, 0, 7).equals(p2))
}
