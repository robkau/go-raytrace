package shapes

import (
	"github.com/robkau/go-raytrace/lib/geom"
	"github.com/robkau/go-raytrace/lib/materials"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func Test_NewSphere_DefaultTransform(t *testing.T) {
	s := NewSphere()

	assert.Equal(t, geom.NewIdentityMatrixX4(), s.t)
}
func Test_NewSphere_HasDefaultMaterial(t *testing.T) {
	s := NewSphere()

	assert.Equal(t, materials.NewMaterial(), s.m)
}

func Test_NewGlassSphere(t *testing.T) {
	s := NewSphere()
	s.SetMaterial(materials.NewGlassMaterial())
	m := s.GetMaterial()

	assert.Equal(t, geom.NewIdentityMatrixX4(), s.GetTransform())
	assert.Equal(t, 0.95, m.Transparency)
	assert.Equal(t, 1.5, m.RefractiveIndex)
}

func Test_Sphere_SetTransform(t *testing.T) {
	var s Shape = NewSphere()
	tr := geom.Translate(2, 3, 4)
	s.SetTransform(tr)

	assert.Equal(t, tr, s.GetTransform())
}

func Test_Sphere_SetMaterial(t *testing.T) {
	s := NewSphere()
	m := materials.NewMaterial()
	m.Ambient = 1.0

	s.m = m

	assert.Equal(t, m, s.m)
}

func Test_RayIntersectSphere(t *testing.T) {
	s := NewSphere()
	r := geom.RayWith(geom.NewPoint(0, 0, -5), geom.NewVector(0, 0, 1))

	xs := s.Intersect(r)

	assert.Len(t, xs.I, 2)
	assert.Equal(t, 4.0, xs.I[0].T)
	assert.Equal(t, 6.0, xs.I[1].T)
}

func Test_RayIntersectSphere_Tangent(t *testing.T) {
	s := NewSphere()
	r := geom.RayWith(geom.NewPoint(0, 1, -5), geom.NewVector(0, 0, 1))

	xs := s.Intersect(r)

	assert.Len(t, xs.I, 2)
	assert.Equal(t, 5.0, xs.I[0].T)
	assert.Equal(t, 5.0, xs.I[1].T)
}

func Test_RayIntersectSphere_Miss(t *testing.T) {
	s := NewSphere()
	r := geom.RayWith(geom.NewPoint(0, 2, -5), geom.NewVector(0, 0, 1))

	xs := s.Intersect(r)

	assert.Len(t, xs.I, 0)
}

func Test_RayIntersectSphere_FromInside(t *testing.T) {
	s := NewSphere()
	r := geom.RayWith(geom.NewPoint(0, 0, 0), geom.NewVector(0, 0, 1))

	xs := s.Intersect(r)

	assert.Len(t, xs.I, 2)
	assert.Equal(t, -1.0, xs.I[0].T)
	assert.Equal(t, 1.0, xs.I[1].T)
}

func Test_RayIntersectSphere_Behind(t *testing.T) {
	s := NewSphere()
	r := geom.RayWith(geom.NewPoint(0, 0, 5), geom.NewVector(0, 0, 1))

	xs := s.Intersect(r)

	assert.Len(t, xs.I, 2)
	assert.Equal(t, -6.0, xs.I[0].T)
	assert.Equal(t, -4.0, xs.I[1].T)
}

func Test_RayIntersectSphere_ObjectSet(t *testing.T) {
	s := NewSphere()
	r := geom.RayWith(geom.NewPoint(0, 0, -5), geom.NewVector(0, 0, 1))

	xs := s.Intersect(r)

	assert.Len(t, xs.I, 2)
	assert.Equal(t, s, xs.I[0].O)
	assert.Equal(t, s, xs.I[1].O)
}

func Test_ScaledSphere_Intersect_Ray(t *testing.T) {
	s := NewSphere()
	s.SetTransform(geom.Scale(2, 2, 2))
	r := geom.RayWith(geom.NewPoint(0, 0, -5), geom.NewVector(0, 0, 1))

	xs := s.Intersect(r)

	assert.Len(t, xs.I, 2)
	assert.Equal(t, 3.0, xs.I[0].T)
	assert.Equal(t, 7.0, xs.I[1].T)
}

func Test_TranslatedSphere_Intersect_Ray(t *testing.T) {
	s := NewSphere()
	s.SetTransform(geom.Translate(5, 0, 0))
	r := geom.RayWith(geom.NewPoint(0, 0, -5), geom.NewVector(0, 0, 1))

	xs := s.Intersect(r)

	assert.Len(t, xs.I, 0)
}

func Test_NormalX(t *testing.T) {
	s := NewSphere()

	n := s.NormalAt(geom.NewPoint(1, 0, 0), Intersection{})

	assert.Equal(t, geom.NewVector(1, 0, 0), n)
}

func Test_NormalY(t *testing.T) {
	s := NewSphere()

	n := s.NormalAt(geom.NewPoint(0, 1, 0), Intersection{})

	assert.Equal(t, geom.NewVector(0, 1, 0), n)
}

func Test_NormalZ(t *testing.T) {
	s := NewSphere()

	n := s.NormalAt(geom.NewPoint(0, 0, 1), Intersection{})

	assert.Equal(t, geom.NewVector(0, 0, 1), n)
}

func Test_NormalXYZ(t *testing.T) {
	s := NewSphere()

	n := s.NormalAt(geom.NewPoint(math.Sqrt(3)/3, math.Sqrt(3)/3, math.Sqrt(3)/3), Intersection{})

	assert.Equal(t, geom.NewVector(math.Sqrt(3)/3, math.Sqrt(3)/3, math.Sqrt(3)/3), n)
}

func Test_Normal_Translated(t *testing.T) {
	s := NewSphere()
	s.SetTransform(geom.Translate(0, 1, 0))

	n := s.NormalAt(geom.NewPoint(0, 1.70711, -0.70711), Intersection{}).RoundTo(5)

	assert.Equal(t, geom.NewVector(0, 0.70711, -0.70711), n)
}

func Test_Normal_ScaledAndRotated(t *testing.T) {
	s := NewSphere()
	s.SetTransform(geom.Scale(1, 0.5, 1).MulX4Matrix(geom.RotateZ(math.Pi / 5)))

	n := s.NormalAt(geom.NewPoint(0, math.Sqrt2/2, -math.Sqrt2/2), Intersection{}).RoundTo(5)

	assert.Equal(t, geom.NewVector(0, 0.97014, -0.24254), n)
}
