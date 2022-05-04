package shapes

import (
	"github.com/google/uuid"
	"github.com/robkau/go-raytrace/lib/geom"
	"github.com/robkau/go-raytrace/lib/materials"
)

type Shape interface {
	Intersect(geom.Ray) Intersections
	LocalIntersect(r geom.Ray) Intersections
	NormalAt(at geom.Tuple, i Intersection) geom.Tuple
	WorldToObject(p geom.Tuple) geom.Tuple
	NormalToWorld(normal geom.Tuple) geom.Tuple
	GetTransform() *geom.X4Matrix
	SetTransform(matrix *geom.X4Matrix)
	GetMaterial() *materials.Material
	SetMaterial(*materials.Material)
	Bounds() Bounds
	GetShadowless() bool
	SetShadowless(s bool)
	GetShaded() bool
	SetShaded(s bool)
	GetParent() Group
	SetParent(g Group)
	Id() string
}

type baseShape struct {
	parent     Group
	t          *geom.X4Matrix
	m          *materials.Material
	id         string
	shadowless bool
	unshaded   bool
}

// todo https://forum.raytracerchallenge.com/thread/144/marching-cubes-isosurfaces

func newBaseShape() baseShape {
	return baseShape{
		t:  geom.NewIdentityMatrixX4(),
		m:  materials.NewMaterial(),
		id: uuid.NewString(),
	}
}

func (b *baseShape) GetTransform() *geom.X4Matrix {
	return b.t
}

func (b *baseShape) SetTransform(matrix *geom.X4Matrix) {
	b.t = matrix
}

func (b *baseShape) GetMaterial() *materials.Material {
	if b.parent != nil {
		m := b.parent.GetMaterial()
		if m != nil {
			return m
		}
	}
	if b.m == nil {
		panic("getmaterial returned nil for shape")
	}
	return b.m
}

func (b *baseShape) SetMaterial(material *materials.Material) {
	b.m = material
}

func (b *baseShape) WorldToObject(p geom.Tuple) geom.Tuple {
	if b.parent != nil {
		p = b.parent.WorldToObject(p)
	}
	return b.t.Invert().MulTuple(p)
}

func (b *baseShape) NormalToWorld(normal geom.Tuple) geom.Tuple {
	normal = b.t.Invert().Transpose().MulTuple(normal)
	normal.C = 0
	normal = normal.Normalize()

	if b.parent != nil {
		normal = b.parent.NormalToWorld(normal)
	}
	return normal
}

func (b *baseShape) GetShadowless() bool {
	if b.parent != nil {
		return b.parent.GetShadowless()
	}
	return b.shadowless
}

func (b *baseShape) SetShadowless(s bool) {
	b.shadowless = s
}

func (b *baseShape) GetShaded() bool {
	return !b.unshaded
}

func (b *baseShape) SetShaded(s bool) {
	b.unshaded = !s
}

func (b *baseShape) GetParent() Group {
	return b.parent
}

func (b *baseShape) SetParent(g Group) {
	b.parent = g
}

func (b *baseShape) Id() string {
	return b.id
}

// invert ray from object's transformation matrix then call shape-specific normal logic
func NormalAt(s Shape, p geom.Tuple, lnaf func(p geom.Tuple, in Intersection) geom.Tuple, i Intersection) geom.Tuple {
	lp := s.WorldToObject(p)
	ln := lnaf(lp, i)
	return s.NormalToWorld(ln)
}

// invert ray from object's transformation matrix then call shape-specific intersection logic
func Intersect(r geom.Ray, t *geom.X4Matrix, lif func(geom.Ray) Intersections) Intersections {
	lr := r.Transform(t.Invert())
	return lif(lr)
}
