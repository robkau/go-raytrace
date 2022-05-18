package shapes

import (
	"github.com/robkau/go-raytrace/lib/geom"
)

type SmoothTriangle struct {
	*Triangle
	n1 geom.Tuple
	n2 geom.Tuple
	n3 geom.Tuple
}

func NewSmoothTriangle(p1, p2, p3, n1, n2, n3 geom.Tuple) *SmoothTriangle {
	t := &SmoothTriangle{
		Triangle: NewTriangle(p1, p2, p3),
		n1:       n1,
		n2:       n2,
		n3:       n3,
	}
	return t
}

func (t *SmoothTriangle) Intersect(ray geom.Ray) *Intersections {
	return Intersect(ray, t.t, t.LocalIntersect)
}

func (t *SmoothTriangle) LocalIntersect(r geom.Ray) *Intersections {
	tHit, u, v, ok := t.localIntersectHits(r)
	if !ok {
		return NewIntersections()
	}

	return NewIntersections(NewIntersectionWithUV(tHit, t, u, v))
}

func (t *SmoothTriangle) NormalAt(p geom.Tuple, i Intersection) geom.Tuple {
	return NormalAt(t, p, t.LocalNormalAt, i)
}

func (t *SmoothTriangle) LocalNormalAt(p geom.Tuple, i Intersection) geom.Tuple {
	if !i.UvSet {
		return t.Triangle.LocalNormalAt(p, i)
	}
	return t.n2.Mul(i.U).Add(t.n3.Mul(i.V)).Add(t.n1.Mul(1 - i.U - i.V))
}

func (t *SmoothTriangle) Normals() []geom.Tuple {
	return []geom.Tuple{t.n1, t.n2, t.n3}
}
