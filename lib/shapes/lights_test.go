package shapes

import (
	"github.com/robkau/go-raytrace/lib/colors"
	"github.com/robkau/go-raytrace/lib/geom"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
)

func Test_NewPointLight(t *testing.T) {
	intensity := colors.White()
	position := geom.ZeroPoint()

	l := NewPointLight(position, intensity)

	assert.Equal(t, intensity, l.Intensity)
	assert.Equal(t, position, l.Position)
}

func Test_NewAreaLight(t *testing.T) {
	corner := geom.ZeroPoint()
	v1 := geom.NewVector(2, 0, 0)
	v2 := geom.NewVector(0, 0, 1)

	light := NewAreaLight(corner, v1, 4, v2, 2, colors.White(), NewJitterSequence(0.5))

	require.Equal(t, corner, light.Corner)
	require.Equal(t, geom.NewVector(0.5, 0, 0), light.UVec)
	require.Equal(t, 4, light.USteps)
	require.Equal(t, geom.NewVector(0, 0, 0.5), light.VVec)
	require.Equal(t, 2, light.VSteps)
	require.Equal(t, 8, light.Samples)
	require.Equal(t, geom.NewPoint(1, 0, 0.5), light.Center)
}

func Test_AreaLight_Points(t *testing.T) {
	corner := geom.ZeroPoint()
	v1 := geom.NewVector(2, 0, 0)
	v2 := geom.NewVector(0, 0, 1)

	light := NewAreaLight(corner, v1, 4, v2, 2, colors.White(), NewJitterSequence(0.5))

	type args struct {
		u      int
		v      int
		expect geom.Tuple
	}

	tests := []args{
		{0, 0, geom.NewPoint(0.25, 0, 0.25)},
		{1, 0, geom.NewPoint(0.75, 0, 0.25)},
		{0, 1, geom.NewPoint(0.25, 0, 0.75)},
		{2, 0, geom.NewPoint(1.25, 0, 0.25)},
		{3, 1, geom.NewPoint(1.75, 0, 0.75)},
	}

	for ti, tt := range tests {
		t.Run(t.Name()+strconv.Itoa(ti), func(t *testing.T) {
			result := light.PointOnLight(tt.u, tt.v)
			require.Equal(t, tt.expect, result)
		})
	}
}
