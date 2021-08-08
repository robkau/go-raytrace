package shapes

import (
	"github.com/robkau/go-raytrace/lib/geom"
	"github.com/stretchr/testify/require"
	"math"
	"testing"
)

func Test_NewGroup(t *testing.T) {
	g := NewGroup()

	require.Equal(t, geom.NewIdentityMatrixX4(), g.GetTransform())
	require.Len(t, g.GetChildren(), 0)
}

func Test_TestShapeHasParentAttribute(t *testing.T) {
	s := newTestShape()

	require.Nil(t, s.GetParent())
}

func Test_BaseShapeGetSetParentAttribute(t *testing.T) {
	s := NewCube()

	require.Nil(t, s.GetParent())

	g := NewGroup()
	s.SetParent(g)

	require.Equal(t, g, s.GetParent())

	s.SetParent(nil)
	require.Nil(t, s.GetParent())
}

func Test_GroupAddChild(t *testing.T) {
	g := NewGroup()
	s := NewCube()

	g.AddChild(s)

	require.Len(t, g.GetChildren(), 1)
	require.Equal(t, s, g.GetChildren()[0])
	require.Equal(t, g, s.GetParent())
}

func Test_EmptyGroupRayIntersect(t *testing.T) {
	g := NewGroup()
	r := geom.RayWith(geom.NewPoint(0, 0, 0), geom.NewVector(0, 0, 1))

	xs := g.LocalIntersect(r)

	require.Len(t, xs.I, 0)
}

func Test_GroupRayIntersects(t *testing.T) {
	g := NewGroup()
	s1 := NewSphere()
	s2 := NewSphere()
	s3 := NewSphere()
	s2.SetTransform(geom.Translate(0, 0, -3))
	s3.SetTransform(geom.Translate(5, 0, 0))
	g.AddChild(s1)
	g.AddChild(s2)
	g.AddChild(s3)

	r := geom.RayWith(geom.NewPoint(0, 0, -5), geom.NewVector(0, 0, 1))

	xs := g.LocalIntersect(r)

	require.Len(t, xs.I, 4)
	require.Equal(t, s2, xs.I[0].O)
	require.Equal(t, s2, xs.I[1].O)
	require.Equal(t, s1, xs.I[2].O)
	require.Equal(t, s1, xs.I[3].O)
}

func Test_TransformedGroupRayIntersects(t *testing.T) {
	g := NewGroup()
	g.SetTransform(geom.Scale(2, 2, 2))
	s := NewSphere()
	s.SetTransform(geom.Translate(5, 0, 0))
	g.AddChild(s)

	r := geom.RayWith(geom.NewPoint(10, 0, -10), geom.NewVector(0, 0, 1))

	xs := g.Intersect(r)

	require.Len(t, xs.I, 2)
}

func Test_WorldToObjectSpace(t *testing.T) {
	g1 := NewGroup()
	g1.SetTransform(geom.RotateY(math.Pi / 2))
	g2 := NewGroup()
	g2.SetTransform(geom.Scale(2, 2, 2))
	g1.AddChild(g2)
	s := NewSphere()
	s.SetTransform(geom.Translate(5, 0, 0))
	g2.AddChild(s)

	p := s.WorldToObject(geom.NewPoint(-2, 0, -10))

	require.Equal(t, geom.NewPoint(0, 0, -1), p.RoundTo(5))
}

func Test_NormalVectorObjectToWorldSpace(t *testing.T) {
	g1 := NewGroup()
	g1.SetTransform(geom.RotateY(math.Pi / 2))
	g2 := NewGroup()
	g2.SetTransform(geom.Scale(1, 2, 3))
	g1.AddChild(g2)
	s := NewSphere()
	s.SetTransform(geom.Translate(5, 0, 0))
	g2.AddChild(s)

	n := s.NormalToWorld(geom.NewVector(math.Sqrt(3)/3, math.Sqrt(3)/3, math.Sqrt(3)/3))

	require.Equal(t, geom.NewVector(0.2857, 0.4286, -0.8571), n.RoundTo(4))
}

func Test_NormalOnChildObject(t *testing.T) {
	g1 := NewGroup()
	g1.SetTransform(geom.RotateY(math.Pi / 2))
	g2 := NewGroup()
	g2.SetTransform(geom.Scale(1, 2, 3))
	g1.AddChild(g2)
	s := NewSphere()
	s.SetTransform(geom.Translate(5, 0, 0))
	g2.AddChild(s)

	n := s.NormalAt(geom.NewPoint(1.7321, 1.1547, -5.5774), Intersection{})

	require.Equal(t, geom.NewVector(0.286, 0.429, -0.857), n.RoundTo(3))
}
