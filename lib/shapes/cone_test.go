package shapes

import (
	"github.com/robkau/go-raytrace/lib/geom"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math"
	"testing"
)

func Test_NewCone_BoundsOf(t *testing.T) {
	s := NewCone(-5, 3, true)
	s.SetTransform(geom.Translate(0, 1, 0)) // no effect

	assert.Equal(t, geom.NewPoint(-5, -5, -5), s.BoundsOf().Min)
	assert.Equal(t, geom.NewPoint(5, 3, 5), s.BoundsOf().Max)
}

func Test_NewInfiniteCone_BoundsOf(t *testing.T) {
	s := NewInfiniteCone()
	s.SetTransform(geom.Translate(0, 1, 0)) // no effect

	assert.Equal(t, geom.NegInfPoint(), s.BoundsOf().Min)
	assert.Equal(t, geom.PosInfPoint(), s.BoundsOf().Max)
}

func Test_ConeIntersections(t *testing.T) {
	type args struct {
		origin    geom.Tuple
		direction geom.Tuple
		t0        float64
		t1        float64
	}
	tests := []struct {
		name string
		args args
	}{
		{"1", args{geom.NewPoint(0, 0, -5), geom.NewVector(0, 0, 1), 5, 5}},
		{"2", args{geom.NewPoint(0, 0, -5), geom.NewVector(1, 1, 1), 8.66025, 8.66025}},
		{"3", args{geom.NewPoint(1, 1, -5), geom.NewVector(-0.5, -1, 1), 4.55006, 49.44994}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewInfiniteCone()
			d := tt.args.direction.Normalize()
			r := geom.RayWith(tt.args.origin, d)

			xs := c.Intersect(r)

			assert.Equal(t, tt.args.t0, geom.RoundTo(xs.I[0].T, 5))
			assert.Equal(t, tt.args.t1, geom.RoundTo(xs.I[1].T, 5))
		})
	}
}

func Test_ConeIntersections_Caps(t *testing.T) {
	type args struct {
		origin    geom.Tuple
		direction geom.Tuple
		count     int
	}
	tests := []struct {
		name string
		args args
	}{
		{"1", args{geom.NewPoint(0, 0, -5), geom.NewVector(0, 1, 0), 0}},
		{"2", args{geom.NewPoint(0, 0, -0.25), geom.NewVector(0, 1, 1), 2}},
		{"3", args{geom.NewPoint(0, 0, -0.25), geom.NewVector(0, 1, 0), 4}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCone(-0.5, 0.5, true)

			d := tt.args.direction.Normalize()
			r := geom.RayWith(tt.args.origin, d)

			xs := c.Intersect(r)

			assert.Len(t, xs.I, tt.args.count)
		})
	}
}

func Test_ConeIntersections_RayParallelToSides(t *testing.T) {
	c := NewInfiniteCone()
	d := geom.NewVector(0, 1, 1).Normalize()
	r := geom.RayWith(geom.NewPoint(0, 0, -1), d)

	xs := c.Intersect(r)

	assert.Equal(t, 0.35355, geom.RoundTo(xs.I[0].T, 5))
}

func Test_ConeNormal_EndCaps(t *testing.T) {
	type args struct {
		point  geom.Tuple
		normal geom.Tuple
	}
	tests := []struct {
		name string
		args args
	}{
		{"1", args{geom.ZeroPoint(), geom.ZeroVector()}},
		{"2", args{geom.NewPoint(1, 1, 1), geom.NewVector(1, -math.Sqrt(2), 1)}},
		{"3", args{geom.NewPoint(-1, -1, 0), geom.NewVector(-1, 1, 0)}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewInfiniteCone()
			n := c.LocalNormalAt(tt.args.point, Intersection{})

			require.Equal(t, tt.args.normal, n)
		})
	}
}
