package shapes

import (
	"github.com/google/uuid"
	"go-raytrace/lib/geom"
)

type Shape interface {
	Intersect(geom.Ray) Intersections
	NormalAt(geom.Tuple) geom.Tuple
	GetTransform() geom.X4Matrix
	SetTransform(matrix geom.X4Matrix)
	GetMaterial() Material
	SetMaterial(Material)
	Id() string
	GetShadowless() bool
	SetShadowless(s bool)
	GetShaded() bool
	SetShaded(s bool)
}

type baseShape struct {
	t          geom.X4Matrix
	M          Material
	id         string
	shadowless bool
	unshaded   bool
}

func newBaseShape() baseShape {
	return baseShape{
		t:  geom.NewIdentityMatrixX4(),
		M:  NewMaterial(),
		id: newId(),
	}
}

func (b *baseShape) GetTransform() geom.X4Matrix {
	return b.t
}

func (b *baseShape) SetTransform(matrix geom.X4Matrix) {
	b.t = matrix
}

func (b *baseShape) GetMaterial() Material {
	return b.M
}

func (b *baseShape) SetMaterial(material Material) {
	b.M = material
}

func (b *baseShape) Id() string {
	return b.id
}

func (b *baseShape) GetShadowless() bool {
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

// invert ray from object's transformation matrix then call shape-specific normal logic
func NormalAt(p geom.Tuple, t geom.X4Matrix, lnaf func(geom.Tuple) geom.Tuple) geom.Tuple {
	localPoint := t.Invert().MulTuple(p)
	localNormal := lnaf(localPoint)
	worldNormal := t.Invert().Transpose().MulTuple(localNormal)
	worldNormal.C = geom.Vector
	return worldNormal.Normalize()
}

// invert ray from object's transformation matrix then call shape-specific intersection logic
func Intersect(r geom.Ray, t geom.X4Matrix, lif func(geom.Ray) Intersections) Intersections {
	lr := r.Transform(t.Invert())
	return lif(lr)
}

func newId() string {
	u, err := uuid.NewRandom()
	if err != nil {
		panic("fail create uuid for shape")
	}
	return u.String()
}
