package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ColorAdd(t *testing.T) {
	a := color{0.9, 0.6, 0.75}
	b := color{0.7, 0.1, 0.25}

	c := a.add(b)

	assert.True(t, color{1.6, 0.7, 1.0}.equal(c))
}

func Test_ColorSub(t *testing.T) {
	a := color{0.9, 0.6, 0.75}
	b := color{0.7, 0.1, 0.25}

	c := a.sub(b)

	assert.True(t, color{0.2, 0.5, 0.5}.equal(c))

}

func Test_ColorMul(t *testing.T) {
	a := color{1, 0.2, 0.4}
	b := color{0.9, 1, 0.1}

	c := a.mul(b)
	d := b.mul(a)

	assert.True(t, color{0.9, 0.2, 0.04}.equal(c))
	assert.Equal(t, c, d)
}

func Test_ColorMulBy(t *testing.T) {
	a := color{0.2, 0.3, 0.4}

	b := a.mulBy(2)

	assert.True(t, color{0.4, 0.6, 0.8}.equal(b))
}
