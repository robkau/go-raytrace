package view

import (
	"context"
	"fmt"
	"github.com/robkau/coordinate_supplier"
	"github.com/robkau/go-raytrace/lib/colors"
	"github.com/robkau/go-raytrace/lib/geom"
	"github.com/robkau/go-raytrace/lib/shapes"
	"math"
	"sync"
)

type World struct {
	objects     []shapes.Shape
	pointLights []shapes.PointLight
	areaLights  []shapes.AreaLight
}

func NewWorld() *World {
	return &World{
		objects:     []shapes.Shape{},
		pointLights: []shapes.PointLight{},
		areaLights:  []shapes.AreaLight{},
	}
}

func defaultWorld() *World {
	w := NewWorld()

	var s shapes.Shape = shapes.NewSphere()
	m := s.GetMaterial()
	m.Color = colors.NewColor(0.8, 1.0, 0.6)
	m.Diffuse = 0.7
	m.Specular = 0.2
	s.SetMaterial(m)
	w.AddObject(s)

	s = shapes.NewSphere()
	s.SetTransform(geom.Scale(0.5, 0.5, 0.5))
	w.AddObject(s)

	l := shapes.NewPointLight(geom.NewPoint(-10, 10, -10), colors.White())
	w.AddPointLight(l)

	return w
}

func (w *World) AddObject(s shapes.Shape) {
	w.objects = append(w.objects, s)
}

func (w *World) AddPointLight(l shapes.PointLight) {
	w.pointLights = append(w.pointLights, l)
}

func (w *World) AddAreaLight(l shapes.AreaLight) {
	w.areaLights = append(w.areaLights, l)
}

func (w *World) Intersect(r geom.Ray) *shapes.Intersections {
	is := shapes.NewIntersections()
	for _, s := range w.objects {
		xs := s.Intersect(r)
		is.AddFrom(xs)
	}
	return is
}

func (w *World) ShadeHit(c shapes.IntersectionComputed, remaining int) colors.Color {
	col := colors.NewColor(0, 0, 0)

	for _, l := range w.pointLights {
		col = col.Add(shapes.Lighting(c.Object.GetMaterial(), c.Object, l, c.OverPoint, c.Eyev, c.Normalv, IntensityAt(l, c.OverPoint, w)))
	}

	for _, l := range w.areaLights {
		col = col.Add(shapes.Lighting(c.Object.GetMaterial(), c.Object, l, c.OverPoint, c.Eyev, c.Normalv, IntensityAtAreaLight(l, c.OverPoint, w)))
	}

	reflected := w.ReflectedColor(c, remaining)
	refracted := w.RefractedColor(c, remaining)

	m := c.Object.GetMaterial()
	if m.Reflective > 0 && m.Transparency > 0 {
		// todo scale by transparency?
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

func (w *World) Divide(threshold int) {
	for _, c := range w.objects {
		c.Divide(threshold)
	}
}

func (w *World) BoundsOf() *shapes.BoundingBox {
	b := shapes.NewEmptyBoundingBox()
	for _, c := range w.objects {
		b.AddBoundingBoxes(c.BoundsOf())
	}
	return b
}

func (w *World) ColorAt(r geom.Ray, remaining int) colors.Color {
	is := w.Intersect(r)
	defer is.Release()
	i, ok := is.Hit()
	if !ok {
		return colors.Black()
	}

	cs := i.Compute(r, is)
	return w.ShadeHit(cs, remaining)
}

func (w *World) IsShadowed(lightPosition geom.Tuple, p geom.Tuple) bool {
	v := lightPosition.Sub(p)
	distance := v.Mag()
	direction := v.Normalize()

	r := geom.RayWith(p, direction)
	intersections := w.Intersect(r)
	h, ok := intersections.Hit()
	// shadowless object does not cast shadows onto other objects
	if ok && h.T < distance && !h.O.GetShadowless() {
		return true
	}
	return false
}

func IntensityAt(p shapes.PointLight, pt geom.Tuple, w *World) float64 {
	v := w.IsShadowed(p.Position, pt)
	if v {
		return 0.0
	}
	return 1.0
}

func IntensityAtAreaLight(l shapes.AreaLight, pt geom.Tuple, w *World) float64 {
	total := 0.0

	for u := 0; u < l.USteps; u++ {
		for v := 0; v < l.VSteps; v++ {
			lightPosition := l.PointOnLight(v, u) // flipping fixes test. ?
			if !w.IsShadowed(lightPosition, pt) {
				total += 1.0
			}
		}
	}
	return total / float64(l.Samples)
}

// todo replace c.rencder?
// or wrap this with an option to consume all and then return the image and delete c.render

func Render(ctx context.Context, w *World, c Camera, rayBounces int, numGoRoutines int, renderMode coordinate_supplier.Order) (<-chan PixelInfo, error) {
	pi := make(chan PixelInfo, numGoRoutines*2)

	cs, err := coordinate_supplier.NewCoordinateSupplierAtomic(coordinate_supplier.CoordinateSupplierOptions{
		Width:  c.HSize,
		Height: c.VSize,
		Depth:  1,
		Order:  renderMode,
		Repeat: false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed create coordinate supplier: %w", err)
	}

	go func() {
		wg := sync.WaitGroup{}
		for i := 0; i < numGoRoutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for x, y, _, done := cs.Next(); !done; x, y, _, done = cs.Next() {
					select {
					case <-ctx.Done():
						// cancel render goroutine
						return
					default:
						// noop
					}

					r := c.rayForPixel(x, y)
					c := w.ColorAt(r, rayBounces)

					pi <- PixelInfo{
						X: x,
						Y: y,
						C: c,
					}
				}
			}()
		}
		wg.Wait()
		close(pi)
	}()

	return pi, nil
}
