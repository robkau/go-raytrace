package patterns

import (
	"github.com/robkau/go-raytrace/lib/colors"
	"github.com/robkau/go-raytrace/lib/geom"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_PositionAtColorPattern_ObjectTransformed(t *testing.T) {
	p := NewPositionAsColorPattern()

	c := ColorAtShape(p, geom.DoubleScale, geom.NewPoint(2, 3, 4))

	require.Equal(t, colors.NewColor(1, 1.5, 2), c)
}

func Test_PositionAtColorPattern_PatternTransformed(t *testing.T) {
	p := NewPositionAsColorPattern()

	c := ColorAtShape(p, geom.DoubleScale, geom.NewPoint(2, 3, 4))

	require.Equal(t, colors.NewColor(1, 1.5, 2), c)
}

func Test_PositionAtColorPattern_PatternAndObjectTransformed(t *testing.T) {
	p := NewPositionAsColorPattern()
	p.SetTransform(geom.Translate(0.5, 1, 1.5))

	c := ColorAtShape(p, geom.DoubleScale, geom.NewPoint(2.5, 3, 3.5))

	require.Equal(t, colors.NewColor(0.75, 0.5, 0.25), c)
}
