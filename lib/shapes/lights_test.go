package shapes

import (
	"github.com/stretchr/testify/assert"
	"go-raytrace/lib/colors"
	"go-raytrace/lib/geom"
	"testing"
)

func Test_NewLight(t *testing.T) {
	intensity := colors.White()
	position := geom.NewPoint(0, 0, 0)

	l := NewPointLight(position, intensity)

	assert.Equal(t, intensity, l.Intensity)
	assert.Equal(t, position, l.Position)
}
