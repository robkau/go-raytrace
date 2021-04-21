package view

import (
	"go-raytrace/lib/colors"
	"go-raytrace/lib/geom"
	"go-raytrace/lib/shapes"
)

type World struct {
	objects      []shapes.Shape
	lightSources []shapes.PointLight
}

func NewWorld() World {
	return World{
		objects:      []shapes.Shape{},
		lightSources: []shapes.PointLight{},
	}
}

func defaultWorld() World {
	w := NewWorld()

	var s shapes.Shape = shapes.NewSphere()
	m := s.GetMaterial()
	m.Color = colors.NewColor(0.8, 1.0, 0.6)
	m.Diffuse = 0.7
	m.Specular = 0.2
	s = s.SetMaterial(m)
	w.AddObject(s)

	s = shapes.NewSphere()
	s = s.SetTransform(geom.Scale(0.5, 0.5, 0.5))
	w.AddObject(s)

	l := shapes.NewPointLight(geom.NewPoint(-10, 10, -10), colors.White())
	w.AddLight(l)

	return w
}

func (w *World) AddObject(s shapes.Shape) {
	w.objects = append(w.objects, s)
}

func (w *World) AddLight(l shapes.PointLight) {
	w.lightSources = append(w.lightSources, l)
}

func (w *World) Intersect(r geom.Ray) shapes.Intersections {
	is := shapes.NewIntersections()
	for _, s := range w.objects {
		xs := s.Intersect(r)
		is.AddFrom(xs)
	}
	return is
}

func (w *World) ShadeHit(c shapes.IntersectionComputed) colors.Color {
	col := colors.NewColor(0, 0, 0)
	for _, l := range w.lightSources {
		col = col.Add(shapes.Lighting(c.Object.GetMaterial(), c.Object, l, c.OverPoint, c.Eyev, c.Normalv, w.IsShadowed(c.OverPoint)))
	}

	return col
}

func (w *World) ColorAt(r geom.Ray) colors.Color {
	is := w.Intersect(r)
	i, ok := is.Hit()
	if !ok {
		return colors.Black()
	}

	cs := i.Compute(r)
	return w.ShadeHit(cs)
}

func (w *World) IsShadowed(p geom.Tuple) bool {
	for _, l := range w.lightSources {
		v := l.Position.Sub(p)
		distance := v.Mag()
		direction := v.Normalize()

		r := geom.RayWith(p, direction)
		intersections := w.Intersect(r)
		h, ok := intersections.Hit()
		if ok && h.T < distance {
			return true
		} else {
			continue
		}
	}
	return false
}
