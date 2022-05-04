package shapes

import (
	"github.com/robkau/go-raytrace/lib/geom"
	"github.com/robkau/go-raytrace/lib/materials"
)

type Group interface {
	Shape
	GetChildren() []Shape
	AddChild(s Shape)
}

type group struct {
	t          *geom.X4Matrix
	parent     Group
	children   []Shape
	m          *materials.Material
	shadowless bool
	id         string

	//bounds    Bounds
}

func NewGroup() Group {
	return &group{
		t:        geom.NewIdentityMatrixX4(),
		children: make([]Shape, 0),
	}
}

func (g *group) Intersect(ray geom.Ray) Intersections {
	// if ray does not intersect groups bounding box - skip
	if !g.Bounds().Intersects(ray) {
		return NewIntersections()
	}
	return Intersect(ray, g.t, g.LocalIntersect)
}

func (g *group) LocalIntersect(r geom.Ray) Intersections {
	xs := NewIntersections()

	// add the intersection for each child in the group
	// should return sorted by t
	for _, s := range g.children {
		xs.AddFrom(s.Intersect(r))
	}

	return xs
}

func (g *group) Bounds() Bounds {
	b := newBounds()
	for _, c := range g.children {
		childBounds := c.Bounds()
		childBoundsTransformed := childBounds.TransformTo(g.t)
		b = b.Add(childBoundsTransformed.Min, childBoundsTransformed.Max)
	}
	return b
}

func (g *group) GetTransform() *geom.X4Matrix {
	return g.t
}

func (g *group) SetTransform(matrix *geom.X4Matrix) {
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

func (g *group) NormalAt(tuple geom.Tuple, _ Intersection) geom.Tuple {
	panic("calling me on a group is a logic error")
}

func (g *group) GetMaterial() *materials.Material {
	if g.m != nil {
		return g.m
	}
	if g.parent != nil {
		return g.parent.GetMaterial()
	}
	return nil
}

func (g *group) SetMaterial(material *materials.Material) {
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
