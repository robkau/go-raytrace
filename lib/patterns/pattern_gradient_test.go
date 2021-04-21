package patterns

import (
	"github.com/stretchr/testify/assert"
	"go-raytrace/lib/colors"
	"go-raytrace/lib/geom"
	"testing"
)

func Test_GradientPattern(t *testing.T) {
	p := NewGradientPattern(NewSolidColorPattern(colors.White()), NewSolidColorPattern(colors.Black()))

	assert.Equal(t, colors.White(), p.ColorAt(geom.NewPoint(0, 0, 0)))
	assert.Equal(t, colors.NewColor(0.75, 0.75, 0.75), p.ColorAt(geom.NewPoint(0.25, 0, 0)))
	assert.Equal(t, colors.NewColor(0.5, 0.5, 0.5), p.ColorAt(geom.NewPoint(0.5, 0, 0)))
	assert.Equal(t, colors.NewColor(0.25, 0.25, 0.25), p.ColorAt(geom.NewPoint(0.75, 0, 0)))
	assert.Equal(t, colors.White(), p.ColorAt(geom.NewPoint(1, 0, 0)))
	assert.Equal(t, colors.NewColor(0.75, 0.75, 0.75), p.ColorAt(geom.NewPoint(1.25, 0, 0)))
}
