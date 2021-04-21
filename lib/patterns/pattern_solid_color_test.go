package patterns

import (
	"github.com/stretchr/testify/assert"
	"go-raytrace/lib/colors"
	"go-raytrace/lib/geom"
	"testing"
)

func Test_SolidColorPattern(t *testing.T) {
	p := NewSolidColorPattern(colors.NewColor(0.66, 0.55, 0.44))

	assert.Equal(t, colors.NewColor(0.66, 0.55, 0.44), p.ColorAt(geom.NewPoint(0, 0, 0)))
	assert.Equal(t, colors.NewColor(0.66, 0.55, 0.44), p.ColorAt(geom.NewPoint(1.25, 0, 0)))
	assert.Equal(t, colors.NewColor(0.66, 0.55, 0.44), p.ColorAt(geom.NewPoint(0, 0, 2.51)))
	assert.Equal(t, colors.NewColor(0.66, 0.55, 0.44), p.ColorAt(geom.NewPoint(0.708, -22, 0.708)))
}
