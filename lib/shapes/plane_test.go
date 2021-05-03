package shapes

import (
	"github.com/stretchr/testify/assert"
	"go-raytrace/lib/geom"
	"testing"
)

func Test_PlaneNormal_Constant(t *testing.T) {
	p := NewPlane()
	n1 := p.LocalNormalAt(geom.NewPoint(0, 0, 0))
	n2 := p.LocalNormalAt(geom.NewPoint(10, 0, -10))
	n3 := p.LocalNormalAt(geom.NewPoint(-5, 10, 150))

	assert.Equal(t, geom.NewVector(0, 1, 0), n1)
	assert.Equal(t, geom.NewVector(0, 1, 0), n2)
	assert.Equal(t, geom.NewVector(0, 1, 0), n3)
}

func Test_Plane_Intersect_Parallel_Ray(t *testing.T) {
	p := NewPlane()
	r := geom.RayWith(geom.NewPoint(0, 10, 0), geom.NewVector(0, 0, 1))

	xs := p.LocalIntersect(r)

	assert.Len(t, xs.I, 0)
}

func Test_Plane_Intersect_Coplanar_Ray(t *testing.T) {
	p := NewPlane()
	r := geom.RayWith(geom.NewPoint(0, 0, 0), geom.NewVector(0, 0, 1))

	xs := p.LocalIntersect(r)

	assert.Len(t, xs.I, 0)
}

func Test_Plane_Intersect_Ray_Above(t *testing.T) {
	p := NewPlane()
	r := geom.RayWith(geom.NewPoint(0, 1, 0), geom.NewVector(0, -1, 0))

	xs := p.LocalIntersect(r)

	assert.Len(t, xs.I, 1)
	assert.Equal(t, 1.0, xs.I[0].T)
	assert.Equal(t, p, xs.I[0].O)
}

func Test_Plane_Intersect_Ray_Below(t *testing.T) {
	p := NewPlane()
	r := geom.RayWith(geom.NewPoint(0, -1, 0), geom.NewVector(0, 1, 0))

	xs := p.LocalIntersect(r)

	assert.Len(t, xs.I, 1)
	assert.Equal(t, 1.0, xs.I[0].T)
	assert.Equal(t, p, xs.I[0].O)
}
