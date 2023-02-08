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
	r := geom.RayWith(geom.ZeroPoint(), geom.NewVector(0, 0, 1))

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

func Test_GroupBoundsContainsChildren(t *testing.T) {
	s := NewSphere()
	s.SetTransform(geom.Translate(2, 5, -3).MulX4Matrix(geom.Scale(2, 2, 2)))
	c := NewCylinder(-2, 2, true)
	c.SetTransform(geom.Translate(-4, -1, 4).MulX4Matrix(geom.Scale(0.5, 1, 0.5)))

	g := NewGroup()
	g.AddChild(s)
	g.AddChild(c)

	box := g.BoundsOf()

	require.Equal(t, geom.NewPoint(-4.5, -3, -5), box.Min)
	require.Equal(t, geom.NewPoint(4, 7, 4.5), box.Max)
}

func Test_Group_PartitionChildren(t *testing.T) {
	s1 := NewSphere()
	s1.SetTransform(geom.Translate(-2, 0, 0))
	s2 := NewSphere()
	s2.SetTransform(geom.Translate(2, 0, 0))
	s3 := NewSphere()

	g := NewGroup()
	g.AddChild(s1) // todo change to shape...
	g.AddChild(s2)
	g.AddChild(s3)

	left, right := g.PartitionChildren()

	gChildren := g.GetChildren()
	leftChildren := left.GetChildren()
	rightChildren := right.GetChildren()

	require.Equal(t, []Shape{s3}, gChildren)
	require.Equal(t, []Shape{s1}, leftChildren)
	require.Equal(t, []Shape{s2}, rightChildren)
}

func Test_Group_AddSubgroup(t *testing.T) {
	s1 := NewSphere()
	s2 := NewSphere()

	g := NewGroup()

	sg := NewGroup()
	sg.AddChild(s1)
	sg.AddChild(s2)

	g.AddChild(sg)
	require.Len(t, g.GetChildren(), 1)

	gg := g.GetChildren()[0].(Group)
	require.Equal(t, []Shape{s1, s2}, gg.GetChildren())
}

func Test_Group_Divide(t *testing.T) {
	s1 := NewSphere()
	s1.SetTransform(geom.Translate(-2, -2, 0))
	s2 := NewSphere()
	s2.SetTransform(geom.Translate(-2, 2, 0))
	s3 := NewSphere()
	s3.SetTransform(geom.Scale(4, 4, 4))

	g := NewGroup()
	g.AddChild(s1)
	g.AddChild(s2)
	g.AddChild(s3)

	g.Divide(1)

	require.Equal(t, s3, g.GetChildren()[0])

	gg := g.GetChildren()[1].(Group)
	require.Len(t, gg.GetChildren(), 2)

	ggl := gg.GetChildren()[0].(Group)
	ggr := gg.GetChildren()[1].(Group)
	require.Equal(t, []Shape{s1}, ggl.GetChildren())
	require.Equal(t, []Shape{s2}, ggr.GetChildren())
}

func Test_Group_Divide_TooFewChildren(t *testing.T) {
	s1 := NewSphere()
	s1.SetTransform(geom.Translate(-2, 0, 0))
	s2 := NewSphere()
	s2.SetTransform(geom.Translate(2, 1, 0))
	s3 := NewSphere()
	s3.SetTransform(geom.Translate(2, -1, 0))
	s4 := NewSphere()

	sg := NewGroup()
	sg.AddChild(s1)
	sg.AddChild(s2)
	sg.AddChild(s3)

	g := NewGroup()
	g.AddChild(sg)
	g.AddChild(s4)

	g.Divide(3)

	gc := g.GetChildren()
	require.Len(t, gc, 2)

	require.Equal(t, sg, gc[0])
	require.Equal(t, s4, gc[1])

	sgc := sg.GetChildren()
	require.Len(t, sgc, 2)

	sgcl := sgc[0].(Group)
	sgcr := sgc[1].(Group)

	require.Equal(t, []Shape{s1}, sgcl.GetChildren())
	require.Equal(t, []Shape{s2, s3}, sgcr.GetChildren())
}
