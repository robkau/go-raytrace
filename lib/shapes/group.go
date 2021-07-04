package shapes

import (
	"go-raytrace/lib/geom"
	"go-raytrace/lib/materials"
)

type Group interface {
	Shape
	GetChildren() []Shape
	AddChild(s Shape)
}

type group struct {
	t        geom.X4Matrix
	parent   Group
	children []Shape
}

func NewGroup() Group {
	return &group{
		t:        geom.NewIdentityMatrixX4(),
		children: make([]Shape, 0),
	}
}

func (g *group) Intersect(ray geom.Ray) Intersections {
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

func (g *group) NormalAt(tuple geom.Tuple) geom.Tuple {
	panic("implement me")
}

func (g *group) GetTransform() geom.X4Matrix {
	return g.t
}

func (g *group) SetTransform(matrix geom.X4Matrix) {
	g.t = matrix
}

func (g *group) GetMaterial() materials.Material {
	panic("implement me")
}

func (g *group) SetMaterial(material materials.Material) {
	panic("implement me")
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

func (g *group) Id() string {
	panic("implement me")
}

func (g *group) GetShadowless() bool {
	panic("implement me")
}

func (g *group) SetShadowless(s bool) {
	panic("implement me")
}

func (g *group) GetShaded() bool {
	panic("implement me")
}

func (g *group) SetShaded(s bool) {
	panic("implement me")
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
