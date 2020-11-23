package main

type world struct {
	objects      []sphere
	lightSources []pointLight
}

func newWorld() world {
	return world{
		objects:      []sphere{},
		lightSources: []pointLight{},
	}
}

func defaultWorld() world {
	w := newWorld()

	s := newSphere()
	s.m.color = color{0.8, 1.0, 0.6}
	s.m.diffuse = 0.7
	s.m.specular = 0.2
	w.addObject(s)

	s = newSphere()
	s = s.setTransform(scale(0.5, 0.5, 0.5))
	w.addObject(s)

	l := newPointLight(newPoint(-10, 10, -10), color{1, 1, 1})
	w.addLight(l)

	return w
}

func (w *world) addObject(s sphere) {
	w.objects = append(w.objects, s)
}

func (w *world) addLight(l pointLight) {
	w.lightSources = append(w.lightSources, l)
}

func (w *world) intersect(r ray) intersections {
	is := newIntersections()
	for _, s := range w.objects {
		xs := s.intersect(r)
		is.addFrom(xs)
	}
	return is
}

func (w *world) shadeHit(c intersectionComputed) color {
	col := color{}
	for _, l := range w.lightSources {
		col = col.add(lighting(c.object.m, l, c.overPoint, c.eyev, c.normalv, w.isShadowed(c.overPoint)))
	}

	return col
}

func (w *world) colorAt(r ray) color {
	is := w.intersect(r)
	i, ok := is.hit()
	if !ok {
		return color{0, 0, 0}
	}

	cs := i.compute(r)
	return w.shadeHit(cs)
}

func (w *world) isShadowed(p tuple) bool {
	for _, l := range w.lightSources {
		v := l.position.sub(p)
		distance := v.mag()
		direction := v.normalize()

		r := rayWith(p, direction)
		intersections := w.intersect(r)
		h, ok := intersections.hit()
		if ok && h.t < distance {
			return true
		} else {
			continue
		}
	}
	return false
}
