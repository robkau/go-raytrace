package colors

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ColorAdd(t *testing.T) {
	a := Color{0.9, 0.6, 0.75}
	b := Color{0.7, 0.1, 0.25}

	c := a.Add(b)

	assert.True(t, Color{1.6, 0.7, 1.0}.Equal(c))
}

func Test_ColorSub(t *testing.T) {
	a := Color{0.9, 0.6, 0.75}
	b := Color{0.7, 0.1, 0.25}

	c := a.Sub(b)

	assert.True(t, Color{0.2, 0.5, 0.5}.Equal(c))

}

func Test_ColorMul(t *testing.T) {
	a := Color{1, 0.2, 0.4}
	b := Color{0.9, 1, 0.1}

	c := a.Mul(b)
	d := b.Mul(a)

	assert.True(t, Color{0.9, 0.2, 0.04}.Equal(c))
	assert.Equal(t, c, d)
}

func Test_ColorMulBy(t *testing.T) {
	a := Color{0.2, 0.3, 0.4}

	b := a.MulBy(2)

	assert.True(t, Color{0.4, 0.6, 0.8}.Equal(b))
}
