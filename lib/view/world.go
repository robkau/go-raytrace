package view

import (
	"go-raytrace/lib/colors"
	"go-raytrace/lib/geom"
	"go-raytrace/lib/shapes"
	"math"
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

func (w *World) ShadeHit(c shapes.IntersectionComputed, remaining int) colors.Color {
	col := colors.NewColor(0, 0, 0)
	for _, l := range w.lightSources {
		col = col.Add(shapes.Lighting(c.Object.GetMaterial(), c.Object, l, c.OverPoint, c.Eyev, c.Normalv, w.IsShadowed(c.OverPoint)))
	}

	reflected := w.ReflectedColor(c, remaining)
	refracted := w.RefractedColor(c, remaining)

	m := c.Object.GetMaterial()
	if m.Reflective > 0 && m.Transparency > 0 {
		reflectance := c.Schlick()
		return col.Add(reflected.MulBy(reflectance)).Add(refracted.MulBy(1 - reflectance))
	}

	return col.Add(reflected).Add(refracted)
}

func (w *World) ReflectedColor(c shapes.IntersectionComputed, remaining int) colors.Color {
	col := colors.NewColor(0, 0, 0)

	if remaining <= 0 {
		return col
	}

	if c.Object.GetMaterial().Reflective == 0 {
		return col
	}

	reflectRay := geom.RayWith(c.OverPoint, c.Reflectv)
	col = w.ColorAt(reflectRay, remaining-1)

	return col.MulBy(c.Object.GetMaterial().Reflective)
}

func (w *World) RefractedColor(c shapes.IntersectionComputed, remaining int) colors.Color {
	col := colors.NewColor(0, 0, 0)
	if remaining == 0 || // limited recursion
		c.Object.GetMaterial().Transparency == 0 { // opaque material
		return col
	}

	nRatio := c.N1 / c.N2
	cosI := c.Eyev.Dot(c.Normalv)
	sin2T := nRatio * nRatio * (1 - cosI*cosI)
	if sin2T > 1 {
		// total internal refraction
		return col
	}

	cosT := math.Sqrt(1.0 - sin2T)

	direction := c.Normalv.Mul(nRatio*cosI - cosT).Sub(c.Eyev.Mul(nRatio))
	refractedRay := geom.RayWith(c.UnderPoint, direction)

	col = w.ColorAt(refractedRay, remaining-1).MulBy(c.Object.GetMaterial().Transparency)
	return col
}

func (w *World) ColorAt(r geom.Ray, remaining int) colors.Color {
	is := w.Intersect(r)
	i, ok := is.Hit()
	if !ok {
		return colors.Black()
	}

	cs := i.Compute(r, is)
	return w.ShadeHit(cs, remaining)
}

func (w *World) IsShadowed(p geom.Tuple) bool {
	for _, l := range w.lightSources {
		v := l.Position.Sub(p)
		distance := v.Mag()
		direction := v.Normalize()

		r := geom.RayWith(p, direction)
		intersections := w.Intersect(r)
		h, ok := intersections.Hit()
		if ok && h.T < distance && !h.O.GetShadowless() {
			return true
		} else {
			continue
		}
	}
	return false
}
