package shapes

import "github.com/robkau/go-raytrace/lib/geom"

type WrappedRay struct {
	geom.Ray
	Is *Intersections
}

func (w *WrappedRay) Reset() {
	w.Ray = geom.Ray{}
	w.Is.Reset()
}
