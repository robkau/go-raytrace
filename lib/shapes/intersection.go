package shapes

import (
	"github.com/robkau/go-raytrace/lib/geom"
	"math"
	"sort"
	"sync"
)

type Intersection struct {
	T     float64
	O     Shape
	U     float64
	V     float64
	UvSet bool
}

func NewIntersection(t float64, o Shape) Intersection {
	return Intersection{T: t, O: o}
}

func NewIntersectionWithUV(t float64, o Shape, u, v float64) Intersection {
	return Intersection{T: t, O: o, U: u, V: v, UvSet: true}
}

func (i Intersection) Hit() bool {
	return i.T > 0
}

type IntersectionComputed struct {
	t          float64
	Object     Shape
	inside     bool
	point      geom.Tuple
	OverPoint  geom.Tuple
	UnderPoint geom.Tuple
	Eyev       geom.Tuple
	Normalv    geom.Tuple
	Reflectv   geom.Tuple
	N1         float64
	N2         float64
}

func (i Intersection) Compute(r geom.Ray, xs *Intersections) IntersectionComputed {
	c := IntersectionComputed{}
	c.t = i.T
	c.Object = i.O
	c.point = r.Position(c.t)
	c.Eyev = r.Direction.Neg()

	if len(xs.I) == 0 {
		c.Normalv = c.Object.NormalAt(c.point, Intersection{})
	} else {
		for _, x := range xs.I {
			if geom.AlmostEqual(i.T, x.T) && i.O.Id() == x.O.Id() {
				c.Normalv = c.Object.NormalAt(c.point, x)
			}
		}
	}

	if c.Normalv.Dot(c.Eyev) < 0 {
		c.inside = true
		c.Normalv = c.Normalv.Neg()
	}

	// reflection
	c.Reflectv = r.Direction.Reflect(c.Normalv)
	c.OverPoint = c.point.Add(c.Normalv.Mul(geom.FloatComparisonEpsilon))
	c.UnderPoint = c.point.Sub(c.Normalv.Mul(geom.FloatComparisonEpsilon))

	// refraction
	containers := []Shape{}
	for _, x := range xs.I {
		if geom.AlmostEqual(i.T, x.T) && i.O.Id() == x.O.Id() {
			if len(containers) == 0 {
				c.N1 = 1
			} else {
				c.N1 = containers[len(containers)-1].GetMaterial().RefractiveIndex
			}
		}

		found := false
		for j, s := range containers {
			if x.O.Id() == s.Id() {
				found = true
				containers = removeIndex(containers, j)
				break
			}
		}
		if !found {
			containers = append(containers, x.O)
		}

		if geom.AlmostEqual(i.T, x.T) && i.O.Id() == x.O.Id() {
			if len(containers) == 0 {
				c.N2 = 1
			} else {
				c.N2 = containers[len(containers)-1].GetMaterial().RefractiveIndex
			}
			break
		}
	}

	return c
}

func (ic IntersectionComputed) Schlick() float64 {
	cos := ic.Eyev.Dot(ic.Normalv)

	if ic.N1 > ic.N2 {
		n := ic.N1 / ic.N2
		sin2T := n * n * (1.0 - cos*cos)
		if sin2T > 1 {
			return 1
		}

		cosT := math.Sqrt(1.0 - sin2T)
		cos = cosT
	}

	r0 := math.Pow((ic.N1-ic.N2)/(ic.N1+ic.N2), 2)
	return r0 + (1-r0)*math.Pow(1-cos, 5)
}

var intersectionsPool = sync.Pool{
	New: func() interface{} { return SortableIntersections{} },
}

type Intersections struct {
	I SortableIntersections
}

type SortableIntersections []Intersection

func (so SortableIntersections) Len() int      { return len(so) }
func (so SortableIntersections) Swap(i, j int) { so[i], so[j] = so[j], so[i] }
func (so SortableIntersections) Less(i, j int) bool {
	return so[i].T < so[j].T
}

func NewIntersections(s ...Intersection) *Intersections {
	i := &Intersections{
		//I: intersectionsPool.Get().(SortableIntersections), // todo i actually slow things down!
		I: make(SortableIntersections, 0, len(s)),
	}
	i.Add(s...)
	return i
}

func (is *Intersections) Hit() (Intersection, bool) {
	for _, i := range is.I {
		if i.Hit() {
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
	sort.Sort(is.I)
}

func (is *Intersections) AddFrom(s *Intersections) {
	is.Add(s.I...)
}

func (is *Intersections) Release() {
	//is.I = is.I[:0]
	//intersectionsPool.Put(is.I)
}

func removeIndex(s []Shape, index int) []Shape {
	return append(s[:index], s[index+1:]...)
}
