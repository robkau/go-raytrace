package shapes

import (
	"github.com/robkau/go-raytrace/lib/geom"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_NewTriangle(t *testing.T) {
	p1 := geom.NewPoint(0, 1, 0)
	p2 := geom.NewPoint(-1, 0, 0)
	p3 := geom.NewPoint(1, 0, 0)

	tr := NewTriangle(p1, p2, p3)

	require.Equal(t, p1, tr.p1)
	require.Equal(t, p2, tr.p2)
	require.Equal(t, p3, tr.p3)
	require.Equal(t, geom.NewVector(-1, -1, 0), tr.e1)
	require.Equal(t, geom.NewVector(1, -1, 0), tr.e2)
	require.Equal(t, geom.NewVector(0, 0, -1), tr.normal)
}

func Test_TriangleNormal(t *testing.T) {
	p1 := geom.NewPoint(0, 1, 0)
	p2 := geom.NewPoint(-1, 0, 0)
	p3 := geom.NewPoint(1, 0, 0)

	tr := NewTriangle(p1, p2, p3)

	require.Equal(t, tr.normal, tr.LocalNormalAt(geom.NewPoint(0, 0.5, 0), Intersection{}))
	require.Equal(t, tr.normal, tr.LocalNormalAt(geom.NewPoint(-0.5, 0.75, 0), Intersection{}))
	require.Equal(t, tr.normal, tr.LocalNormalAt(geom.NewPoint(0.5, 0.25, 0), Intersection{}))
}

func Test_IntersectRayParallel(t *testing.T) {
	tr := NewTriangle(geom.NewPoint(0, 1, 0), geom.NewPoint(-1, 0, 0), geom.NewPoint(1, 0, 0))
	r := geom.RayWith(geom.NewPoint(0, -1, -2), geom.NewVector(0, 1, 0))

	xs := tr.LocalIntersect(r)

	require.Len(t, xs.I, 0)
}

func Test_IntersectRayMissesP1P3Edge(t *testing.T) {
	tr := NewTriangle(geom.NewPoint(0, 1, 0), geom.NewPoint(-1, 0, 0), geom.NewPoint(1, 0, 0))
	r := geom.RayWith(geom.NewPoint(1, 1, -2), geom.NewVector(0, 0, 1))

	xs := tr.LocalIntersect(r)

	require.Len(t, xs.I, 0)
}

func Test_IntersectRayMissesP1P2Edge(t *testing.T) {
	tr := NewTriangle(geom.NewPoint(0, 1, 0), geom.NewPoint(-1, 0, 0), geom.NewPoint(1, 0, 0))
	r := geom.RayWith(geom.NewPoint(-1, 1, -2), geom.NewVector(0, 0, 1))

	xs := tr.LocalIntersect(r)

	require.Len(t, xs.I, 0)
}

func Test_IntersectRayMissesP2P3Edge(t *testing.T) {
	tr := NewTriangle(geom.NewPoint(0, 1, 0), geom.NewPoint(-1, 0, 0), geom.NewPoint(1, 0, 0))
	r := geom.RayWith(geom.NewPoint(0, -1, -2), geom.NewVector(0, 0, 1))

	xs := tr.LocalIntersect(r)

	require.Len(t, xs.I, 0)
}

func Test_IntersectRayStrikes(t *testing.T) {
	tr := NewTriangle(geom.NewPoint(0, 1, 0), geom.NewPoint(-1, 0, 0), geom.NewPoint(1, 0, 0))
	r := geom.RayWith(geom.NewPoint(0, 0.5, -2), geom.NewVector(0, 0, 1))

	xs := tr.LocalIntersect(r)

	require.Len(t, xs.I, 1)
	require.Equal(t, 2.0, xs.I[0].T)
}

func Test_TriangleIntersectDoesNotStoreUV(t *testing.T) {
	r := geom.RayWith(geom.NewPoint(-0.2, 0.3, -2), geom.NewVector(0, 0, 1))

	xs := NewTriangle(geom.NewPoint(0, 1, 0), geom.NewPoint(-1, 0, 0), geom.NewPoint(1, 0, 0)).LocalIntersect(r)

	require.Len(t, xs.I, 1)
	require.False(t, xs.I[0].UvSet)
}
