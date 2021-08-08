package shapes

import (
	"github.com/robkau/go-raytrace/lib/geom"
	"math"
)

type Triangle struct {
	baseShape

	p1     geom.Tuple
	p2     geom.Tuple
	p3     geom.Tuple
	e1     geom.Tuple
	e2     geom.Tuple
	normal geom.Tuple
}

func NewTriangle(p1, p2, p3 geom.Tuple) *Triangle {
	t := &Triangle{
		baseShape: newBaseShape(),
		p1:        p1,
		p2:        p2,
		p3:        p3,
		e1:        p2.Sub(p1),
		e2:        p3.Sub(p1),
	}
	t.normal = geom.Cross(t.e2, t.e1).Normalize()
	return t
}

func (t *Triangle) Intersect(ray geom.Ray) Intersections {
	return Intersect(ray, t.t, t.LocalIntersect)
}

func (t *Triangle) LocalIntersect(r geom.Ray) Intersections {
	tHit, _, _, ok := t.localIntersectHits(r)
	if !ok {
		return NewIntersections()
	}

	return NewIntersections(NewIntersection(tHit, t))
}

func (t *Triangle) localIntersectHits(r geom.Ray) (tHit, u, v float64, createIntersections bool) {
	// Implements Möller–Trumbore triangle intersection algorithm
	dirCrossE2 := geom.Cross(r.Direction, t.e2)
	det := t.e1.Dot(dirCrossE2)
	if math.Abs(det) < geom.FloatComparisonEpsilon {
		return 0, 0, 0, false
	}

	f := 1.0 / det
	p1ToOrigin := r.Origin.Sub(t.p1)
	u = f * p1ToOrigin.Dot(dirCrossE2)
	if u < 0 || u > 1 {
		return 0, 0, 0, false
	}

	originCrossE1 := geom.Cross(p1ToOrigin, t.e1)
	v = f * r.Direction.Dot(originCrossE1)
	if v < 0 || (u+v) > 1 {
		return 0, 0, 0, false
	}

	tHit = f * t.e2.Dot(originCrossE1)
	return tHit, u, v, true
}

func (t *Triangle) NormalAt(p geom.Tuple, i Intersection) geom.Tuple {
	return NormalAt(t, p, t.LocalNormalAt, i)
}

func (t *Triangle) LocalNormalAt(p geom.Tuple, _ Intersection) geom.Tuple {
	return t.normal
}

func (t *Triangle) Bounds() Bounds {
	return newBounds(t.p1, t.p2, t.p3)
}

func (t *Triangle) Vertices() []geom.Tuple {
	return []geom.Tuple{t.p1, t.p2, t.p3}
}
