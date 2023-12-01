package shapes

import "github.com/robkau/go-raytrace/lib/geom"

type WrappedRay struct {
	geom.Ray
	Is *Intersections
}

func (w *WrappedRay) Reset() {
	w.Ray.Origin.X = 0
	w.Ray.Origin.Y = 0
	w.Ray.Origin.Z = 0
	w.Ray.Direction.X = 0
	w.Ray.Direction.Y = 0
	w.Ray.Direction.Z = 0
	w.Is.Reset()
}
