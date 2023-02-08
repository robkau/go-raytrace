package shapes

import (
	"github.com/robkau/go-raytrace/lib/geom"
	"github.com/stretchr/testify/require"
	"math"
	"strconv"
	"testing"
)

func Test_EmptyBoundingBox(t *testing.T) {
	box := NewBoundingBox(geom.ZeroPoint(), geom.ZeroPoint())

	require.Equal(t, geom.NewPoint(math.Inf(1), math.Inf(1), math.Inf(1)), box.Min)
	require.Equal(t, geom.NewPoint(math.Inf(-1), math.Inf(-1), math.Inf(-1)), box.Max)
	require.Equal(t, NewEmptyBoundingBox(), box)
}

func Test_BoundingBox(t *testing.T) {
	box := NewBoundingBox(geom.NewPoint(-1, -2, -3), geom.NewPoint(3, 2, 1))

	require.Equal(t, geom.NewPoint(-1, -2, -3), box.Min)
	require.Equal(t, geom.NewPoint(3, 2, 1), box.Max)
}

func Test_BoundingBox_AddPoints(t *testing.T) {
	box := NewEmptyBoundingBox()

	box.Add(geom.NewPoint(-5, 2, 0))
	box.Add(geom.NewPoint(7, 0, -3))

	require.Equal(t, geom.NewPoint(-5, 0, -3), box.Min)
	require.Equal(t, geom.NewPoint(7, 2, 0), box.Max)
}

func Test_BoundingBox_Transform(t *testing.T) {
	box := NewBoundingBox(geom.NewPoint(-1, -1, -1), geom.NewPoint(1, 1, 1))

	box.Transform(geom.RotateX(math.Pi / 4).MulX4Matrix(geom.RotateY(math.Pi / 4)))

	require.True(t, geom.NewPoint(-1.414213562373095, -1.7071067811865475, -1.7071067811865475).Equals(box.Min))
	require.True(t, geom.NewPoint(1.414213562373095, 1.7071067811865475, 1.7071067811865475).Equals(box.Max))
}

func Test_BoundingBox_AddBoxes(t *testing.T) {
	box1 := NewBoundingBox(geom.NewPoint(-5, -2, 0), geom.NewPoint(7, 4, 4))
	box2 := NewBoundingBox(geom.NewPoint(8, -7, -2), geom.NewPoint(14, 2, 8))

	box1.AddBoundingBoxes(box2)

	require.Equal(t, geom.NewPoint(-5, -7, -2), box1.Min)
	require.Equal(t, geom.NewPoint(14, 4, 8), box1.Max)
}

func Test_BoundingBox_Contains(t *testing.T) {
	box := NewBoundingBox(geom.NewPoint(5, -2, 0), geom.NewPoint(11, 4, 7))

	type args struct {
		val    geom.Tuple
		expect bool
	}

	tests := []args{
		{geom.NewPoint(5, -2, 0), true},
		{geom.NewPoint(11, 4, 7), true},
		{geom.NewPoint(8, 1, 3), true},
		{geom.NewPoint(3, 0, 3), false},
		{geom.NewPoint(8, -4, 3), false},
		{geom.NewPoint(8, 1, -1), false},
		{geom.NewPoint(13, 1, 3), false},
		{geom.NewPoint(8, 5, 3), false},
		{geom.NewPoint(8, 1, 8), false},
	}

	for ti, tt := range tests {
		t.Run(t.Name()+strconv.Itoa(ti), func(t *testing.T) {
			c := box.Contains(tt.val)
			require.Equal(t, tt.expect, c)
		})
	}
}

func Test_BoundingBox_ContainsBox(t *testing.T) {
	box := NewBoundingBox(geom.NewPoint(5, -2, 0), geom.NewPoint(11, 4, 7))

	type args struct {
		val    *BoundingBox
		expect bool
	}

	tests := []args{
		{NewBoundingBox(geom.NewPoint(5, -2, 0), geom.NewPoint(11, 4, 7)), true},
		{NewBoundingBox(geom.NewPoint(6, -1, 1), geom.NewPoint(10, 3, 6)), true},
		{NewBoundingBox(geom.NewPoint(4, -3, -1), geom.NewPoint(10, 3, 6)), false},
		{NewBoundingBox(geom.NewPoint(6, -1, 1), geom.NewPoint(12, 5, 8)), false},
	}

	for ti, tt := range tests {
		t.Run(t.Name()+strconv.Itoa(ti), func(t *testing.T) {
			c := box.ContainsBox(tt.val)
			require.Equal(t, tt.expect, c)
		})
	}
}

func Test_BoundingBox_Intersect(t *testing.T) {
	box := NewBoundingBox(geom.NewPoint(-1, -1, -1), geom.NewPoint(1, 1, 1))

	type args struct {
		origin    geom.Tuple
		direction geom.Tuple
		expect    bool
	}

	tests := []args{
		{geom.NewPoint(5, 0.5, 0), geom.NewVector(-1, 0, 0), true},
		{geom.NewPoint(-5, 0.5, 0), geom.NewVector(1, 0, 0), true},
		{geom.NewPoint(0.5, 5, 0), geom.NewVector(0, -1, 0), true},
		{geom.NewPoint(0.5, -5, 0), geom.NewVector(0, 1, 0), true},
		{geom.NewPoint(0.5, 0, 5), geom.NewVector(0, 0, -1), true},
		{geom.NewPoint(0.5, 0, -5), geom.NewVector(0, 0, 1), true},
		{geom.NewPoint(0, 0.5, 0), geom.NewVector(0, 0, 1), true},
		{geom.NewPoint(-2, 0, 0), geom.NewVector(2, 4, 6), false},
		{geom.NewPoint(0, -2, 0), geom.NewVector(6, 2, 4), false},
		{geom.NewPoint(0, 0, -2), geom.NewVector(4, 6, 2), false},
		{geom.NewPoint(2, 0, 2), geom.NewVector(0, 0, -1), false},
		{geom.NewPoint(0, 2, 2), geom.NewVector(0, -1, 0), false},
		{geom.NewPoint(2, 2, 0), geom.NewVector(-1, 0, 0), false},
	}

	for ti, tt := range tests {
		t.Run(t.Name()+strconv.Itoa(ti), func(t *testing.T) {
			v := box.Intersect(geom.RayWith(tt.origin, tt.direction))
			require.Equal(t, v, tt.expect)
		})
	}
}

func Test_BoundingBoxNonCubic_Intersect(t *testing.T) {
	box := NewBoundingBox(geom.NewPoint(5, -2, 0), geom.NewPoint(11, 4, 7))

	type args struct {
		origin    geom.Tuple
		direction geom.Tuple
		expect    bool
	}

	tests := []args{
		{geom.NewPoint(15, 1, 2), geom.NewVector(-1, 0, 0), true},
		{geom.NewPoint(-5, -1, 4), geom.NewVector(1, 0, 0), true},
		{geom.NewPoint(7, 6, 5), geom.NewVector(0, -1, 0), true},
		{geom.NewPoint(9, -5, 6), geom.NewVector(0, 1, 0), true},
		{geom.NewPoint(8, 2, 12), geom.NewVector(0, 0, -1), true},
		{geom.NewPoint(6, 0, -5), geom.NewVector(0, 0, 1), true},
		{geom.NewPoint(8, 1, 3.5), geom.NewVector(0, 0, 1), true},
		{geom.NewPoint(9, -1, -8), geom.NewVector(2, 4, 6), false},
		{geom.NewPoint(8, 3, -4), geom.NewVector(6, 2, 4), false},
		{geom.NewPoint(9, -1, -2), geom.NewVector(4, 6, 2), false},
		{geom.NewPoint(4, 0, 9), geom.NewVector(0, 0, -1), false},
		{geom.NewPoint(8, 6, -1), geom.NewVector(0, -1, 0), false},
		{geom.NewPoint(12, 5, 4), geom.NewVector(-1, 0, 0), false},
	}

	for ti, tt := range tests {
		t.Run(t.Name()+strconv.Itoa(ti), func(t *testing.T) {
			v := box.Intersect(geom.RayWith(tt.origin, tt.direction))
			require.Equal(t, v, tt.expect)
		})
	}
}

func Test_SplitPerfectCube(t *testing.T) {
	box := NewBoundingBox(geom.NewPoint(-1, -4, -5), geom.NewPoint(9, 6, 5))

	left, right := box.SplitBounds()

	require.Equal(t, geom.NewPoint(-1, -4, -5), left.Min)
	require.Equal(t, geom.NewPoint(4, 6, 5), left.Max)
	require.Equal(t, geom.NewPoint(4, -4, -5), right.Min)
	require.Equal(t, geom.NewPoint(9, 6, 5), right.Max)
}

func Test_SplitXWideCube(t *testing.T) {
	box := NewBoundingBox(geom.NewPoint(-1, -2, -3), geom.NewPoint(9, 5.5, 3))

	left, right := box.SplitBounds()

	require.Equal(t, geom.NewPoint(-1, -2, -3), left.Min)
	require.Equal(t, geom.NewPoint(4, 5.5, 3), left.Max)
	require.Equal(t, geom.NewPoint(4, -2, -3), right.Min)
	require.Equal(t, geom.NewPoint(9, 5.5, 3), right.Max)
}

func Test_SplitYWideCube(t *testing.T) {
	box := NewBoundingBox(geom.NewPoint(-1, -2, -3), geom.NewPoint(5, 8, 3))

	left, right := box.SplitBounds()

	require.Equal(t, geom.NewPoint(-1, -2, -3), left.Min)
	require.Equal(t, geom.NewPoint(5, 3, 3), left.Max)
	require.Equal(t, geom.NewPoint(-1, 3, -3), right.Min)
	require.Equal(t, geom.NewPoint(5, 8, 3), right.Max)
}

func Test_SplitZWideCube(t *testing.T) {
	box := NewBoundingBox(geom.NewPoint(-1, -2, -3), geom.NewPoint(5, 3, 7))

	left, right := box.SplitBounds()

	require.Equal(t, geom.NewPoint(-1, -2, -3), left.Min)
	require.Equal(t, geom.NewPoint(5, 3, 2), left.Max)
	require.Equal(t, geom.NewPoint(-1, -2, 2), right.Min)
	require.Equal(t, geom.NewPoint(5, 3, 7), right.Max)
}
