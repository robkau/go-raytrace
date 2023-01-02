package patterns

import (
	"github.com/robkau/go-raytrace/lib/colors"
	"github.com/robkau/go-raytrace/lib/geom"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_RingPattern(t *testing.T) {
	p := NewRingPattern(NewSolidColorPattern(colors.White()), NewSolidColorPattern(colors.Black()))

	assert.Equal(t, colors.White(), p.ColorAt(geom.ZeroPoint()))
	assert.Equal(t, colors.Black(), p.ColorAt(geom.NewPoint(1, 0, 0)))
	assert.Equal(t, colors.Black(), p.ColorAt(geom.NewPoint(0, 0, 1)))
	assert.Equal(t, colors.Black(), p.ColorAt(geom.NewPoint(0.708, 0, 0.708)))
}
