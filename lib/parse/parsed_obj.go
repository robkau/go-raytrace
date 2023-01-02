package parse

import (
	"errors"
	"fmt"
	"github.com/robkau/go-raytrace/lib/geom"
	"github.com/robkau/go-raytrace/lib/shapes"
)

type objParsed struct {
	triangleCoordinates []geom.Tuple
	normals             []geom.Tuple
	linesIgnored        int
	defaultGroup        shapes.Group
	currentNamedGroup   string
	namedGroups         map[string]shapes.Group
}

func newObjParsed() *objParsed {
	return &objParsed{
		triangleCoordinates: []geom.Tuple{geom.ZeroPoint()},  // 0th element is always ignored
		normals:             []geom.Tuple{geom.ZeroVector()}, // 0th element is always ignored
		defaultGroup:        shapes.NewGroup(),
		currentNamedGroup:   defaultGroup,
		namedGroups:         map[string]shapes.Group{},
	}
}

const defaultGroup = "default"

func (o *objParsed) addVertex(v geom.Tuple) {
	o.triangleCoordinates = append(o.triangleCoordinates, v)
}

func (o *objParsed) getVertex(i int) (geom.Tuple, error) {
	if i <= 0 {
		return geom.Tuple{}, errors.New("index must be >= 1")
	}

	if i >= len(o.triangleCoordinates) {
		return geom.Tuple{}, fmt.Errorf("index %d out of bounds with array length %d", i, len(o.triangleCoordinates))
	}

	return o.triangleCoordinates[i], nil
}

func (o *objParsed) addNormal(v geom.Tuple) {
	o.normals = append(o.normals, v)
}

func (o *objParsed) getNormal(i int) (geom.Tuple, error) {
	if i <= 0 {
		return geom.Tuple{}, errors.New("index must be >= 1")
	}

	if i >= len(o.normals) {
		return geom.Tuple{}, fmt.Errorf("index %d out of bounds with array length %d", i, len(o.normals))
	}

	return o.normals[i], nil
}

func (o *objParsed) addShape(s shapes.Shape) {
	if o.currentNamedGroup == defaultGroup {
		o.defaultGroup.AddChild(s)
		return
	}
	g, ok := o.namedGroups[o.currentNamedGroup]
	if !ok {
		o.namedGroups[o.currentNamedGroup] = shapes.NewGroup()
		g = o.namedGroups[o.currentNamedGroup]
		o.defaultGroup.AddChild(g)
	}
	g.AddChild(s)
}
