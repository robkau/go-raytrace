package patterns

import (
	"github.com/stretchr/testify/assert"
	"go-raytrace/lib/colors"
	"go-raytrace/lib/geom"
	"testing"
)

func Test_BlendPattern_RedBlue(t *testing.T) {
	p := NewBlendPattern(NewSolidColorPattern(colors.Red()), NewSolidColorPattern(colors.Blue()), 0.5)

	c := colors.NewColor(0.5, 0, 0.5)
	assert.Equal(t, c, p.ColorAt(geom.NewPoint(0, 0, 0)))
	assert.Equal(t, c, p.ColorAt(geom.NewPoint(0.99, 0, 0)))
	assert.Equal(t, c, p.ColorAt(geom.NewPoint(1.01, 0, 0)))
	assert.Equal(t, c, p.ColorAt(geom.NewPoint(10.01, 300, -249)))
}
