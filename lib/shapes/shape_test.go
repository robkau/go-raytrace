package shapes

import (
	"github.com/robkau/go-raytrace/lib/geom"
	"github.com/robkau/go-raytrace/lib/materials"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math"
	"testing"
)

func Test_GetDefaultTransformation(t *testing.T) {
	s := newTestShape()
	require.Equal(t, geom.NewIdentityMatrixX4(), s.GetTransform())
}

func Test_SetTransform(t *testing.T) {
	var s Shape = newTestShape()
	s.SetTransform(geom.Translate(2, 3, 4))
	require.Equal(t, geom.Translate(2, 3, 4), s.GetTransform())
}

func Test_GetDefaultMaterial(t *testing.T) {
	s := newTestShape()
	require.Equal(t, materials.NewMaterial(), s.GetMaterial())
}

func Test_SetMaterial(t *testing.T) {
	var s Shape = newTestShape()
	m := s.GetMaterial()
	m.Ambient = 1
	s.SetMaterial(m)
	require.Equal(t, m, s.GetMaterial())
}

func Test_IntersectScaledShape(t *testing.T) {
	r := geom.RayWith(geom.NewPoint(0, 0, -5), geom.NewVector(0, 0, 1))
	var ts = newTestShape()
	var s Shape = ts
	ts.SetTransform(geom.Scale(2, 2, 2))

	_ = s.Intersect(r)

	assert.Equal(t, geom.NewPoint(0, 0, -2.5), ts.savedRay.Origin)
	assert.Equal(t, geom.NewVector(0, 0, 0.5), ts.savedRay.Direction)
}

func Test_IntersectTranslatedShape(t *testing.T) {
	r := geom.RayWith(geom.NewPoint(0, 0, -5), geom.NewVector(0, 0, 1))
	var ts = newTestShape()
	ts.SetTransform(geom.Translate(5, 0, 0))

	_ = ts.Intersect(r)

	assert.Equal(t, geom.NewPoint(-5, 0, -5), ts.savedRay.Origin)
	assert.Equal(t, geom.NewVector(0, 0, 1), ts.savedRay.Direction)
}

func Test_NormalTranslatedShape(t *testing.T) {
	var ts = newTestShape()
	var s Shape = ts
	ts.SetTransform(geom.Translate(0, 1, 0))

	n := s.NormalAt(geom.NewPoint(0, 1.70711, -0.70711)).RoundTo(5)

	assert.Equal(t, geom.NewVector(0, 0.70711, -0.70711), n)
}

func Test_NormalTransformedShape(t *testing.T) {
	var ts = newTestShape()
	var s Shape = ts
	ts.SetTransform(geom.Scale(1, 0.5, 1).MulX4Matrix(geom.RotateZ(math.Pi / 5)))

	n := s.NormalAt(geom.NewPoint(0, math.Sqrt2/2, -math.Sqrt2/2)).RoundTo(5)

	assert.Equal(t, geom.NewVector(0, 0.97014, -0.24254), n)
}

type testShape struct {
	parent     Group
	t          geom.X4Matrix
	m          materials.Material
	savedRay   geom.Ray
	shadowless bool
	unshaded   bool
}

func newTestShape() *testShape {
	return &testShape{
		t: geom.NewIdentityMatrixX4(),
		m: materials.NewMaterial(),
	}
}

func (t *testShape) Bounds() Bounds {
	return newBounds(geom.NewPoint(-1, -1, -1), geom.NewPoint(1, 1, 1)).TransformTo(t.t)
}

func (t *testShape) Intersect(r geom.Ray) Intersections {
	lr := r.Transform(t.t.Invert())
	t.savedRay = lr
	return Intersections{}
}

func (t *testShape) LocalIntersect(r geom.Ray) Intersections {
	return NewIntersections()
}

func (t *testShape) WorldToObject(p geom.Tuple) geom.Tuple {
	return p
}

func (t *testShape) NormalToWorld(normal geom.Tuple) geom.Tuple {
	return normal
}

func (t *testShape) Id() string {
	return ""
}

func (t *testShape) NormalAt(p geom.Tuple) geom.Tuple {
	localPoint := t.t.Invert().MulTuple(p)
	localNormal := geom.NewVector(localPoint.X, localPoint.Y, localPoint.Z)
	worldNormal := t.t.Invert().Transpose().MulTuple(localNormal)
	worldNormal.C = geom.Vector
	return worldNormal.Normalize()
}

func (t *testShape) GetTransform() geom.X4Matrix {
	return t.t
}

func (t *testShape) SetTransform(m geom.X4Matrix) {
	t.t = m
}

func (t *testShape) GetMaterial() materials.Material {
	return t.m
}

func (t *testShape) SetMaterial(m materials.Material) {
	t.m = m
}

func (t *testShape) GetShadowless() bool {
	return t.shadowless
}

func (t *testShape) SetShadowless(s bool) {
	t.shadowless = s
}

func (t *testShape) GetShaded() bool {
	return !t.unshaded
}

func (t *testShape) SetShaded(s bool) {
	t.unshaded = !s
}

func (t *testShape) GetParent() Group {
	return t.parent
}

func (t *testShape) SetParent(g Group) {
	t.parent = g
}
