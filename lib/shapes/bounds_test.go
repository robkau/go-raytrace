package shapes

import (
	"github.com/robkau/go-raytrace/lib/geom"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math"
	"testing"
)

func Test_Bounds_MinMax(t *testing.T) {
	type args struct {
		p1  geom.Tuple
		p2  geom.Tuple
		min geom.Tuple
		max geom.Tuple
	}
	tests := []struct {
		name string
		args args
	}{
		{"p1 smaller p2", args{geom.NewPoint(0.1, 0, -5), geom.NewPoint(1, 2, 3), geom.NewPoint(0.1, 0, -5), geom.NewPoint(1, 2, 3)}},
		{"p1 equal p2", args{geom.NewPoint(1.2, 4, -5), geom.NewPoint(1.2, 4, -5), geom.NewPoint(1.2, 4, -5), geom.NewPoint(1.2, 4, -5)}},
		{"p1 some larger than p2", args{geom.NewPoint(0, -999, 5), geom.NewPoint(-1, -2, -3), geom.NewPoint(-1, -999, -3), geom.NewPoint(0, -2, 5)}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := newBounds(
				tt.args.p1,
				tt.args.p2,
			)

			require.Equal(t, tt.args.min, b.Min)
			require.Equal(t, tt.args.max, b.Max)
		})
	}
}

func Test_GroupBounds_TransformedChildren(t *testing.T) {
	// c1 extends up 8 units from the origin, 0.5 in other directions
	c1 := NewCube()
	c1.SetTransform(geom.Translate(0, 4, 0).MulX4Matrix(geom.Scale(0.5, 4, 0.5)))

	// c2 extends 2 units in all directions from origin
	c2 := NewCube()
	c2.SetTransform(geom.Scale(2, 2, 2))

	g := NewGroup()
	g.AddChild(c1)
	g.AddChild(c2)

	gb := g.Bounds()

	require.Equal(t, geom.NewPoint(-2, -2, -2), gb.Min)
	require.Equal(t, geom.NewPoint(2, 8, 2), gb.Max)
}

func Test_GroupBounds_TransformedGroupAndChild(t *testing.T) {
	// c extends 1 unit in all directions from origin of 1,0,0
	c1 := NewCube()
	c1.SetTransform(geom.Translate(1, 0, 0))

	g := NewGroup()
	g.AddChild(c1)

	// rotated 180 deg around Y axis and scaled x3
	g.SetTransform(geom.RotateY(math.Pi).MulX4Matrix(geom.Scale(3, 3, 3)))

	gb := g.Bounds()

	assert.True(t, geom.NewPoint(-6, -3, -3).Equals(gb.Min))
	assert.True(t, geom.NewPoint(0, 3, 3).Equals(gb.Max))
}

func Test_ShapeHasBounds(t *testing.T) {
	ts := newTestShape()
	s := NewSphere()

	tsb := ts.Bounds()
	sb := s.Bounds()

	require.True(t, tsb.Min.X < tsb.Max.X)
	require.True(t, tsb.Min.Y < tsb.Max.Y)
	require.True(t, tsb.Min.Z < tsb.Max.Z)

	require.True(t, sb.Min.X < sb.Max.X)
	require.True(t, sb.Min.Y < sb.Max.Y)
	require.True(t, sb.Min.Z < sb.Max.Z)
}

func Test_TransformedShapeBounds_Scale_Sphere(t *testing.T) {
	s := NewSphere()
	s.SetTransform(geom.Scale(2, 2, 2))

	sb := s.Bounds()

	require.Equal(t, float64(-2), sb.Min.X)
	require.Equal(t, float64(-2), sb.Min.Y)
	require.Equal(t, float64(-2), sb.Min.Z)

	require.Equal(t, float64(2), sb.Max.X)
	require.Equal(t, float64(2), sb.Max.Y)
	require.Equal(t, float64(2), sb.Max.Z)
}

func Test_TransformedShapeBounds_Scale_Cube(t *testing.T) {
	c := NewCube()
	c.SetTransform(geom.Scale(0.5, 4, 2))

	cb := c.Bounds()

	require.Equal(t, geom.NewPoint(-0.5, -4, -2), cb.Min)
	require.Equal(t, geom.NewPoint(0.5, 4, 2), cb.Max)
}

func Test_TransformedShapeBounds_Translate(t *testing.T) {
	s := NewSphere()
	s.SetTransform(geom.Translate(2, 2, 2))

	sb := s.Bounds()

	require.Equal(t, float64(1), sb.Min.X)
	require.Equal(t, float64(1), sb.Min.Y)
	require.Equal(t, float64(1), sb.Min.Z)

	require.Equal(t, float64(3), sb.Max.X)
	require.Equal(t, float64(3), sb.Max.Y)
	require.Equal(t, float64(3), sb.Max.Z)
}

func Test_TransformedShapeBounds_Rotate45(t *testing.T) {
	s := NewCube()
	s.SetTransform(geom.RotateY(math.Pi / 4))

	sb := s.Bounds()

	require.True(t, geom.AlmostEqual(-math.Sqrt2, sb.Min.X))
	require.True(t, geom.AlmostEqual(-1, sb.Min.Y))
	require.True(t, geom.AlmostEqual(-math.Sqrt2, sb.Min.Z))

	require.True(t, geom.AlmostEqual(math.Sqrt2, sb.Max.X))
	require.True(t, geom.AlmostEqual(1, sb.Max.Y))
	require.True(t, geom.AlmostEqual(math.Sqrt2, sb.Max.Z))
}

func Test_TransformedShapeBounds_Rotate90(t *testing.T) {
	s := NewCube()
	s.SetTransform(geom.RotateY(math.Pi / 2))

	sb := s.Bounds()

	require.Equal(t, float64(-1), sb.Min.X)
	require.Equal(t, float64(-1), sb.Min.Y)
	require.Equal(t, float64(-1), sb.Min.Z)

	require.Equal(t, float64(1), sb.Max.X)
	require.Equal(t, float64(1), sb.Max.Y)
	require.Equal(t, float64(1), sb.Max.Z)
}
