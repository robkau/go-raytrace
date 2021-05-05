package shapes

import (
	"github.com/google/uuid"
	"go-raytrace/lib/geom"
)

type Shape interface {
	Intersect(geom.Ray) Intersections
	NormalAt(geom.Tuple) geom.Tuple
	GetTransform() geom.X4Matrix
	SetTransform(matrix geom.X4Matrix) Shape
	GetMaterial() Material
	SetMaterial(Material) Shape
	Id() string
	GetShadowless() bool
	SetShadowless(s bool) Shape
	GetShaded() bool
	SetShaded(s bool) Shape
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
