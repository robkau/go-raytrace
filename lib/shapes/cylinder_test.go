package shapes

import (
	"github.com/robkau/go-raytrace/lib/geom"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math"
	"testing"
)

func Test_NewCylinder_BoundsOf(t *testing.T) {
	s := NewCylinder(-5, 3, true)
	s.SetTransform(geom.Translate(0, 1, 0)) // no effect

	assert.Equal(t, geom.NewPoint(-1, -5, -1), s.BoundsOf().Min)
	assert.Equal(t, geom.NewPoint(1, 3, 1), s.BoundsOf().Max)
}

func Test_NewInfiniteCylinder_BoundsOf(t *testing.T) {
	s := NewInfiniteCylinder()
	s.SetTransform(geom.Translate(0, 1, 0)) // no effect

	assert.Equal(t, geom.NewPoint(-1, math.Inf(-1), -1), s.BoundsOf().Min)
	assert.Equal(t, geom.NewPoint(1, math.Inf(1), 1), s.BoundsOf().Max)
}

func Test_RayMissesCylinder(t *testing.T) {
	type args struct {
		origin    geom.Tuple
		direction geom.Tuple
	}
	tests := []struct {
		name string
		args args
	}{
		{"+x up", args{geom.NewPoint(1, 0, 0), geom.NewVector(0, 1, 0)}},
		{"origin up", args{geom.ZeroPoint(), geom.NewVector(0, 1, 0)}},
		{"-x out", args{geom.NewPoint(0, 0, -5), geom.NewVector(1, 1, 1)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewInfiniteCylinder()
			dir := tt.args.direction.Normalize()
			r := geom.RayWith(tt.args.origin, dir)
			xs := c.LocalIntersect(r)

			require.Len(t, xs.I, 0)
		})
	}
}

func Test_RayHitsCylinder(t *testing.T) {
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
		{"same time", args{geom.NewPoint(1, 0, -5), geom.NewVector(0, 0, 1), 5, 5}},
		{"hit 1", args{geom.NewPoint(0, 0, -5), geom.NewVector(0, 0, 1), 4, 6}},
		{"hit 2", args{geom.NewPoint(0.5, 0, -5), geom.NewVector(0.1, 1, 1), 6.80798, 7.08872}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewInfiniteCylinder()
			dir := tt.args.direction.Normalize()
			r := geom.RayWith(tt.args.origin, dir)
			xs := c.LocalIntersect(r)

			require.Len(t, xs.I, 2)
			require.Equal(t, tt.args.t0, geom.RoundTo(xs.I[0].T, 5))
			require.Equal(t, tt.args.t1, geom.RoundTo(xs.I[1].T, 5))
		})
	}
}

func Test_CylinderNormalVector(t *testing.T) {
	type args struct {
		point  geom.Tuple
		normal geom.Tuple
	}
	tests := []struct {
		name string
		args args
	}{
		{"+x", args{geom.NewPoint(1, 0, 0), geom.NewVector(1, 0, 0)}},
		{"y up", args{geom.NewPoint(0, 5, -1), geom.NewVector(0, 0, -1)}},
		{"y down", args{geom.NewPoint(0, -2, 1), geom.NewVector(0, 0, 1)}},
		{"-xy", args{geom.NewPoint(-1, 1, 0), geom.NewVector(-1, 0, 0)}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewInfiniteCylinder()
			n := c.LocalNormalAt(tt.args.point, Intersection{})
			require.Equal(t, tt.args.normal, n)
		})
	}
}

func Test_InfiniteCylinderMinMax(t *testing.T) {
	c := NewInfiniteCylinder()
	assert.Equal(t, math.Inf(-1), c.minimum)
	assert.Equal(t, math.Inf(1), c.maximum)
}

func Test_ConstrainedCylinderIntersections(t *testing.T) {
	type args struct {
		point     geom.Tuple
		direction geom.Tuple
		count     int
	}
	tests := []struct {
		name string
		args args
	}{
		{"1", args{geom.NewPoint(0, 1.5, 0), geom.NewVector(0.1, 1, 0), 0}},
		{"2", args{geom.NewPoint(0, 3, -5), geom.NewVector(0, 0, 1), 0}},
		{"3", args{geom.NewPoint(0, 0, -5), geom.NewVector(0, 0, 1), 0}},
		{"4", args{geom.NewPoint(0, 2, -5), geom.NewVector(0, 0, 1), 0}},
		{"5", args{geom.NewPoint(0, 1, -5), geom.NewVector(0, 0, 1), 0}},
		{"6", args{geom.NewPoint(0, 1.5, -2), geom.NewVector(0, 0, 1), 2}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCylinder(1, 2, false)
			d := tt.args.direction.Normalize()
			r := geom.RayWith(tt.args.point, d)

			xs := c.Intersect(r)

			assert.Equal(t, tt.args.count, len(xs.I))
		})
	}
}

func Test_ConstrainedCappedCylinderIntersections(t *testing.T) {
	type args struct {
		point     geom.Tuple
		direction geom.Tuple
		count     int
	}
	tests := []struct {
		name string
		args args
	}{
		{"1", args{geom.NewPoint(0, 3, 0), geom.NewVector(0, -1, 0), 2}},
		{"2", args{geom.NewPoint(0, 3, -2), geom.NewVector(0, -1, 2), 2}},
		{"3 - corner case", args{geom.NewPoint(0, 4, -2), geom.NewVector(0, -1, 1), 2}},
		{"4", args{geom.NewPoint(0, 0, -2), geom.NewVector(0, 1, 2), 2}},
		{"5 - corner case", args{geom.NewPoint(0, -1, -2), geom.NewVector(0, 1, 1), 2}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCylinder(1, 2, true)
			d := tt.args.direction.Normalize()
			r := geom.RayWith(tt.args.point, d)

			xs := c.Intersect(r)

			assert.Equal(t, tt.args.count, len(xs.I))
		})
	}
}

func Test_ConstrainedCappedCylinderNormalAt(t *testing.T) {
	type args struct {
		point  geom.Tuple
		normal geom.Tuple
	}
	tests := []struct {
		name string
		args args
	}{
		{"1", args{geom.NewPoint(0, 1, 0), geom.NewVector(0, -1, 0)}},
		{"2", args{geom.NewPoint(0.5, 1, 0), geom.NewVector(0, -1, 0)}},
		{"3", args{geom.NewPoint(0, 1, 0.5), geom.NewVector(0, -1, 0)}},
		{"4", args{geom.NewPoint(0, 2, 0), geom.NewVector(0, 1, 0)}},
		{"5", args{geom.NewPoint(0.5, 2, 0), geom.NewVector(0, 1, 0)}},
		{"6", args{geom.NewPoint(0, 2, 0.5), geom.NewVector(0, 1, 0)}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCylinder(1, 2, true)

			n := c.NormalAt(tt.args.point, Intersection{})

			require.Equal(t, tt.args.normal, n)
		})
	}
}

func Test_NewCylinderUncapped(t *testing.T) {
	c := NewCylinder(1, 2, false)
	require.Equal(t, false, c.capped)
}

func Test_NewCylinderCapped(t *testing.T) {
	c := NewCylinder(1, 2, true)
	require.Equal(t, true, c.capped)
}
