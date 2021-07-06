package shapes

import (
	"github.com/robkau/go-raytrace/lib/geom"
	"math"
)

type Plane struct {
	baseShape
}

func NewPlane() *Plane {
	return &Plane{
		baseShape: newBaseShape(),
	}
}

func (p *Plane) Bounds() Bounds {
	return newBounds(geom.NewPoint(math.Inf(-1), 0, math.Inf(-1)), geom.NewPoint(math.Inf(1), 0, math.Inf(1))).TransformTo(p.t)
}

func (p *Plane) NormalAt(pt geom.Tuple) geom.Tuple {
	return NormalAt(p, pt, p.LocalNormalAt)
}

func (p *Plane) LocalNormalAt(pt geom.Tuple) geom.Tuple {
	return geom.NewVector(0, 1, 0)
}

func (p *Plane) Intersect(r geom.Ray) Intersections {
	return Intersect(r, p.t, p.LocalIntersect)
}

func (p *Plane) LocalIntersect(r geom.Ray) Intersections {
	if math.Abs(r.Direction.Y) < geom.FloatComparisonEpsilon {
		return NewIntersections()
	}
	t := -r.Origin.Y / r.Direction.Y
	return NewIntersections(NewIntersection(t, p))
}
