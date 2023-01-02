package shapes

import (
	"github.com/robkau/go-raytrace/lib/geom"
	"github.com/stretchr/testify/require"
	"testing"
)

func newTestSmoothTriangle() *SmoothTriangle {
	return NewSmoothTriangle(geom.NewPoint(0, 1, 0), geom.NewPoint(-1, 0, 0), geom.NewPoint(1, 0, 0), geom.NewVector(0, 1, 0), geom.NewVector(-1, 0, 0), geom.NewVector(1, 0, 0))
}

func Test_NewSmoothTriangle(t *testing.T) {
	st := newTestSmoothTriangle()

	require.Equal(t, geom.NewPoint(0, 1, 0), st.p1)
	require.Equal(t, geom.NewPoint(-1, 0, 0), st.p2)
	require.Equal(t, geom.NewPoint(1, 0, 0), st.p3)
	require.Equal(t, geom.NewVector(0, 1, 0), st.n1)
	require.Equal(t, geom.NewVector(-1, 0, 0), st.n2)
	require.Equal(t, geom.NewVector(1, 0, 0), st.n3)
}

func Test_SmoothTriangleIntersectStoresUV(t *testing.T) {
	r := geom.RayWith(geom.NewPoint(-0.2, 0.3, -2), geom.NewVector(0, 0, 1))

	xs := newTestSmoothTriangle().LocalIntersect(r)

	require.Len(t, xs.I, 1)
	require.True(t, xs.I[0].UvSet)
	require.True(t, geom.AlmostEqual(0.45, xs.I[0].U))
	require.True(t, geom.AlmostEqual(0.25, xs.I[0].V))
}

func Test_SmoothTriangle_InterpolatesNormalWithUV(t *testing.T) {
	i := NewIntersectionWithUV(1, newTestSmoothTriangle(), 0.45, 0.25)
	n := newTestSmoothTriangle().NormalAt(geom.ZeroPoint(), i)

	require.Equal(t, geom.NewVector(-0.5547, 0.83205, 0), n.RoundTo(5))
}

func Test_NormalAtSmoothTriangle(t *testing.T) {
	i := NewIntersectionWithUV(1, newTestSmoothTriangle(), 0.45, 0.25)
	r := geom.RayWith(geom.NewPoint(-0.2, 0.3, -2), geom.NewVector(0, 0, -1))
	xs := NewIntersections(i)

	comps := i.Compute(r, xs)

	require.Equal(t, geom.NewVector(-0.5547, 0.83205, 0), comps.Normalv.RoundTo(5))
}
