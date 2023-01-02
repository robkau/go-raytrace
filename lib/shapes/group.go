package shapes

import (
	"github.com/robkau/go-raytrace/lib/geom"
	"github.com/robkau/go-raytrace/lib/materials"
)

type Group interface {
	Shape
	GetChildren() []Shape
	AddChild(s Shape)
	PartitionChildren() (left, right Group)
}

type group struct {
	t          geom.X4Matrix
	parent     Group
	children   []Shape
	m          materials.Material
	shadowless bool
	id         string

	bounds *BoundingBox
}

func NewGroup() Group {
	return &group{
		t:        geom.NewIdentityMatrixX4(),
		children: make([]Shape, 0),
		bounds:   NewEmptyBoundingBox(),
	}
}

func (g *group) Intersect(ray geom.Ray) *Intersections {
	return Intersect(ray, g.t, g.LocalIntersect)
}

func (g *group) LocalIntersect(r geom.Ray) *Intersections {
	xs := NewIntersections()

	if g.bounds.Intersect(r) {
		// add the intersection for each child in the group
		// should return sorted by t
		for _, s := range g.children {
			xs.AddFrom(s.Intersect(r))
		}
	}

	return xs
}

func (g *group) BoundsOf() *BoundingBox {
	b := NewEmptyBoundingBox()
	for _, c := range g.children {
		childBoundsTransformed := ParentSpaceBoundsOf(c)
		b.AddBoundingBoxes(childBoundsTransformed)
	}
	return b
}

func (g *group) Divide(threshold int) {
	if threshold <= len(g.children) {
		g.Invalidate()
		left, right := g.PartitionChildren()
		if len(left.GetChildren()) > 0 {
			g.AddChild(left)
		}
		if len(right.GetChildren()) > 0 {
			g.AddChild(right)
		}
	}

	for _, child := range g.children {
		child.Divide(threshold)
	}
}

func (g *group) Invalidate() {
	g.bounds = g.BoundsOf()

	if g.parent != nil {
		g.parent.Invalidate()
	}
}

func (g *group) GetTransform() geom.X4Matrix {
	return g.t
}

func (g *group) SetTransform(matrix geom.X4Matrix) {
	g.t = matrix
}

func (g *group) WorldToObject(p geom.Tuple) geom.Tuple {
	if g.parent != nil {
		p = g.parent.WorldToObject(p)
	}
	return g.t.Invert().MulTuple(p)
}

func (g *group) NormalToWorld(normal geom.Tuple) geom.Tuple {
	normal = g.t.Invert().Transpose().MulTuple(normal)
	normal.C = 0
	normal = normal.Normalize()

	if g.parent != nil {
		normal = g.parent.NormalToWorld(normal)
	}
	return normal
}

func (g *group) GetParent() Group {
	return g.parent
}

func (g *group) SetParent(gr Group) {
	g.parent = gr
}

func (g *group) GetChildren() []Shape {
	return g.children
}

func (g *group) AddChild(s Shape) {
	s.SetParent(g)
	g.children = append(g.children, s)
}

func (g *group) PartitionChildren() (left, right Group) {
	left = NewGroup()
	right = NewGroup()

	leftBounds, rightBounds := g.bounds.SplitBounds()
	// zero-length slice with the same underlying array
	remainingChildren := g.children[:0]

	for _, c := range g.children {
		if leftBounds.ContainsBox(ParentSpaceBoundsOf(c)) {
			left.AddChild(c)
		} else if rightBounds.ContainsBox(ParentSpaceBoundsOf(c)) {
			right.AddChild(c)
		} else {
			remainingChildren = append(remainingChildren, c)
		}
	}

	g.children = remainingChildren
	g.Invalidate()
	return
}

func (g *group) NormalAt(tuple geom.Tuple, _ Intersection) geom.Tuple {
	panic("calling me on a group is a logic error")
}

func (g *group) GetMaterial() materials.Material {
	if !materials.IsZeroMaterial(g.m) {
		return g.m
	}
	if g.parent != nil {
		return g.parent.GetMaterial()
	}
	return materials.ZeroMaterial()
}

func (g *group) SetMaterial(material materials.Material) {
	g.m = material
}

func (g *group) Id() string {
	return g.id
}

func (g *group) GetShadowless() bool {
	if g.shadowless {
		return true
	}
	if g.parent != nil {
		return g.parent.GetShadowless()
	}
	return false
}

func (g *group) SetShadowless(s bool) {
	g.shadowless = s
}

func (g *group) GetShaded() bool {
	panic("calling me on a group is a logic error")
}

func (g *group) SetShaded(s bool) {
	panic("calling me on a group is a logic error")
}
