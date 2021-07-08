package shapes

import (
	"github.com/robkau/go-raytrace/lib/geom"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_(t *testing.T) {

}

func TestCube_RayIntersects(t *testing.T) {
	type args struct {
		origin    geom.Tuple
		direction geom.Tuple
	}
	type want struct {
		t1 float64
		t2 float64
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{"pos x", args{geom.NewPoint(5, 0.5, 0), geom.NewVector(-1, 0, 0)}, want{4, 6}},
		{"neg x", args{geom.NewPoint(-5, 0.5, 0), geom.NewVector(1, 0, 0)}, want{4, 6}},
		{"pos y", args{geom.NewPoint(0.5, 5, 0), geom.NewVector(0, -1, 0)}, want{4, 6}},
		{"neg y", args{geom.NewPoint(0.5, -5, 0), geom.NewVector(0, 1, 0)}, want{4, 6}},
		{"pos z", args{geom.NewPoint(0.5, 0, 5), geom.NewVector(0, 0, -1)}, want{4, 6}},
		{"neg z", args{geom.NewPoint(0.5, 0, -5), geom.NewVector(0, 0, 1)}, want{4, 6}},
		{"inside +x", args{geom.NewPoint(0.5, 0, -0.3), geom.NewVector(1, 0, 0)}, want{-1.5, 0.5}},
		{"inside -y", args{geom.NewPoint(0.5, 0, -0.3), geom.NewVector(0, -1, 0)}, want{-1, 1}},
		{"inside +z", args{geom.NewPoint(0.5, 0, -0.3), geom.NewVector(0, 0, 1)}, want{-0.7, 1.3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCube()
			r := geom.RayWith(tt.args.origin, tt.args.direction)

			xs := c.LocalIntersect(r)

			assert.Len(t, xs.I, 2)
			assert.Equal(t, tt.want.t1, xs.I[0].T)
			assert.Equal(t, tt.want.t2, xs.I[1].T)
		})
	}
}

func TestCube_RayMisses(t *testing.T) {
	type args struct {
		origin    geom.Tuple
		direction geom.Tuple
	}
	tests := []struct {
		name string
		args args
	}{
		{"neg x", args{geom.NewPoint(-2, 0, 0), geom.NewVector(0.2673, 0.5345, 0.8018)}},
		{"neg y", args{geom.NewPoint(0, -2, 0), geom.NewVector(0.8018, 0.2673, 0.5345)}},
		{"neg z", args{geom.NewPoint(0, 0, -2), geom.NewVector(0.5345, 0.8018, 0.2673)}},
		{"pos xz", args{geom.NewPoint(2, 0, 2), geom.NewVector(0, 0, -1)}},
		{"pos yz", args{geom.NewPoint(0, 2, 2), geom.NewVector(0, -1, 0)}},
		{"pos xy", args{geom.NewPoint(2, 2, 0), geom.NewVector(-1, 0, 0)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCube()
			r := geom.RayWith(tt.args.origin, tt.args.direction)

			xs := c.LocalIntersect(r)

			assert.Len(t, xs.I, 0)
		})
	}
}

func TestCube_NormalAt(t *testing.T) {
	type args struct {
		point  geom.Tuple
		normal geom.Tuple
	}
	tests := []struct {
		name string
		args args
	}{
		{"pos x", args{geom.NewPoint(1, 0.5, -0.8), geom.NewVector(1, 0, 0)}},
		{"neg x", args{geom.NewPoint(-1, -0.2, 0.9), geom.NewVector(-1, 0, 0)}},
		{"pos y", args{geom.NewPoint(-0.4, 1, -0.1), geom.NewVector(0, 1, 0)}},
		{"neg y", args{geom.NewPoint(0.3, -1, -0.7), geom.NewVector(0, -1, 0)}},
		{"pos z", args{geom.NewPoint(-0.6, 0.3, 1), geom.NewVector(0, 0, 1)}},
		{"neg z", args{geom.NewPoint(0.4, 0.4, -1), geom.NewVector(0, 0, -1)}},
		{"pos xx", args{geom.NewPoint(1, 1, 1), geom.NewVector(1, 0, 0)}},
		{"neg xx", args{geom.NewPoint(-1, -1, -1), geom.NewVector(-1, 0, 0)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCube()
			p := tt.args.point
			normal := c.LocalNormalAt(p)

			require.Equal(t, tt.args.normal, normal)
		})
	}
}
