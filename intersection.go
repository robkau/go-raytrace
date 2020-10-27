package main

import "sort"

type intersections struct {
	i []intersection
}

type intersection struct {
	t float64
	o sphere
}

func newIntersection(t float64, o sphere) intersection {
	return intersection{t, o}
}

func newIntersections(is ...intersection) intersections {
	in := intersections{}
	for _, i := range is {
		in.i = append(in.i, i)
	}
	sort.Slice(in.i, func(i, j int) bool {
		return in.i[i].t < in.i[j].t
	})
	return in
}

func (is intersections) hit() (intersection, bool) {
	for _, i := range is.i {
		if i.t >= 0 {
			return i, true
		}
	}
	return intersection{}, false
}
