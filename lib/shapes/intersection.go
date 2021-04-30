package shapes

import (
	"go-raytrace/lib/geom"
	"sort"
)

type Intersection struct {
	T float64
	o Shape
}

func NewIntersection(t float64, o Shape) Intersection {
	return Intersection{t, o}
}

type IntersectionComputed struct {
	t         float64
	Object    Shape
	inside    bool
	point     geom.Tuple
	OverPoint geom.Tuple
	Eyev      geom.Tuple
	Normalv   geom.Tuple
	Reflectv  geom.Tuple
}

func (i Intersection) Compute(r geom.Ray) IntersectionComputed {
	c := IntersectionComputed{}
	c.t = i.T
	c.Object = i.o
	c.point = r.Position(c.t)
	c.Eyev = r.Direction.Neg()
	c.Normalv = c.Object.NormalAt(c.point)

	if c.Normalv.Dot(c.Eyev) < 0 {
		c.inside = true
		c.Normalv = c.Normalv.Neg()
	}

	c.Reflectv = r.Direction.Reflect(c.Normalv)
	c.OverPoint = c.point.Add(c.Normalv.Mul(geom.FloatComparisonEpsilon))

	return c
}

type Intersections struct {
	I []Intersection
}

func NewIntersections(s ...Intersection) Intersections {
	is := Intersections{}
	is.Add(s...)
	return is
}

func (is Intersections) Hit() (Intersection, bool) {
	for _, i := range is.I {
		if i.T >= 0 {
			return i, true
		}
	}
	return Intersection{}, false
}

func (is *Intersections) Add(s ...Intersection) {
	if len(s) < 1 {
		return
	}
	for _, i := range s {
		is.I = append(is.I, i)
	}
	sort.Slice(is.I, func(i, j int) bool {
		return is.I[i].T < is.I[j].T
	})
}

func (is *Intersections) AddFrom(s Intersections) {
	is.Add(s.I...)
}
