package obj_parse

import "github.com/robkau/go-raytrace/lib/geom"

type ObjParsed struct {
	TriangleCoordinates []geom.Tuple
	LinesIgnored        int
}

func newObjParsed() *ObjParsed {
	return &ObjParsed{
		TriangleCoordinates: []geom.Tuple{geom.NewPoint(0, 0, 0)}, // 0th element is always ignored
	}
}

// todo only access through methods to hide the 0 - call getvertex and error if 0
func (o *ObjParsed) AddVertex(v geom.Tuple) {
	o.TriangleCoordinates = append(o.TriangleCoordinates, v)
}
