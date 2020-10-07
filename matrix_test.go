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
	panic("implement me")
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
	panic("implement me")
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
	panic("implement me")
}
