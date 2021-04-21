package patterns

import (
	"go-raytrace/lib/geom"
)

type basePattern struct {
	t geom.X4Matrix
}

func newBasePattern() basePattern {
	return basePattern{
		t: geom.NewIdentityMatrixX4(),
	}
}

func (b *basePattern) SetTransform(t geom.X4Matrix) {
	b.t = t
}

func (b *basePattern) GetTransform() geom.X4Matrix {
	return b.t
}
