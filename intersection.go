package main

import "sort"

type intersection struct {
	t float64
	o shape
}

func newIntersection(t float64, o shape) intersection {
	return intersection{t, o}
}

type intersectionComputed struct {
	t         float64
	object    shape
	inside    bool
	point     tuple
	overPoint tuple
	eyev      tuple
	normalv   tuple
}

func (i intersection) compute(r ray) intersectionComputed {
	c := intersectionComputed{}
	c.t = i.t
	c.object = i.o
	c.point = r.position(c.t)
	c.eyev = r.direction.neg()
	c.normalv = c.object.normalAt(c.point)

	if c.normalv.dot(c.eyev) < 0 {
		c.inside = true
		c.normalv = c.normalv.neg()
	}

	c.overPoint = c.point.add(c.normalv.mul(floatComparisonEpsilon))

	return c
}

type intersections struct {
	i []intersection
}

func newIntersections(s ...intersection) intersections {
	is := intersections{}
	is.add(s...)
	return is
}

func (is intersections) hit() (intersection, bool) {
	for _, i := range is.i {
		if i.t >= 0 {
			return i, true
		}
	}
	return intersection{}, false
}

func (is *intersections) add(s ...intersection) {
	if len(s) < 1 {
		return
	}
	for _, i := range s {
		is.i = append(is.i, i)
	}
	sort.Slice(is.i, func(i, j int) bool {
		return is.i[i].t < is.i[j].t
	})
}

func (is *intersections) addFrom(s intersections) {
	is.add(s.i...)
}
