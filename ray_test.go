package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_NewRay(t *testing.T) {
	r := rayWith(newPoint(1, 2, 3), newVector(4, 5, 6))

	assert.True(t, newPoint(1, 2, 3).equals(r.origin))
	assert.True(t, newVector(4, 5, 6).equals(r.direction))
}

func Test_Ray_Position(t *testing.T) {
	r := rayWith(newPoint(2, 3, 4), newVector(1, 0, 0))

	assert.True(t, newPoint(2, 3, 4).equals(r.position(0)))
	assert.True(t, newPoint(3, 3, 4).equals(r.position(1)))
	assert.True(t, newPoint(1, 3, 4).equals(r.position(-1)))
	assert.True(t, newPoint(4.5, 3, 4).equals(r.position(2.5)))
}

func Test_Ray_Translate(t *testing.T) {
	r := rayWith(newPoint(1, 2, 3), newVector(0, 1, 0))
	m := translate(3, 4, 5)

	r2 := r.transform(m)

	assert.Equal(t, newPoint(4, 6, 8), r2.origin)
	assert.Equal(t, newVector(0, 1, 0), r2.direction)
}

func Test_Ray_Scale(t *testing.T) {
	r := rayWith(newPoint(1, 2, 3), newVector(0, 1, 0))
	m := scale(2, 3, 4)

	r2 := r.transform(m)

	assert.Equal(t, newPoint(2, 6, 12), r2.origin)
	assert.Equal(t, newVector(0, 3, 0), r2.direction)
}
