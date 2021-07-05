package patterns

import (
	"github.com/robkau/go-raytrace/lib/colors"
	"github.com/robkau/go-raytrace/lib/geom"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_CheckerPattern_RepeatsX(t *testing.T) {
	p := NewCheckerPattern(NewSolidColorPattern(colors.White()), NewSolidColorPattern(colors.Black()))

	assert.Equal(t, colors.White(), p.ColorAt(geom.NewPoint(0, 0, 0)))
	assert.Equal(t, colors.White(), p.ColorAt(geom.NewPoint(0.99, 0, 0)))
	assert.Equal(t, colors.Black(), p.ColorAt(geom.NewPoint(1.01, 0, 0)))

	assert.Equal(t, colors.White(), p.ColorAt(geom.NewPoint(0, 0, 4)))
	assert.Equal(t, colors.White(), p.ColorAt(geom.NewPoint(0.99, 0, 4)))
	assert.Equal(t, colors.Black(), p.ColorAt(geom.NewPoint(1.01, 0, 4)))
}

func Test_CheckerPattern_RepeatsY(t *testing.T) {
	p := NewCheckerPattern(NewSolidColorPattern(colors.White()), NewSolidColorPattern(colors.Black()))

	assert.Equal(t, colors.White(), p.ColorAt(geom.NewPoint(0, 0, 0)))
	assert.Equal(t, colors.White(), p.ColorAt(geom.NewPoint(0, 0.99, 0)))
	assert.Equal(t, colors.Black(), p.ColorAt(geom.NewPoint(0, 1.01, 0)))
}

func Test_CheckerPattern_RepeatsZ(t *testing.T) {
	p := NewCheckerPattern(NewSolidColorPattern(colors.White()), NewSolidColorPattern(colors.Black()))

	assert.Equal(t, colors.White(), p.ColorAt(geom.NewPoint(0, 0, 0)))
	assert.Equal(t, colors.White(), p.ColorAt(geom.NewPoint(0, 0, 0.99)))
	assert.Equal(t, colors.Black(), p.ColorAt(geom.NewPoint(0, 0, 1.01)))
}
