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
		triangleCoordinates: []geom.Tuple{geom.ZeroPoint()},        // 0th element is always ignored
		normals:             []geom.Tuple{geom.NewVector(0, 0, 0)}, // 0th element is always ignored
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

// collapses group, and its child groups, into one group broken into subgroups in cubes of size width
func CollapseGroups(width float64, group shapes.Group) shapes.Group {
	// todo https://forum.raytracerchallenge.com/thread/179/rendering-large-3d-models
	// todo http://raytracerchallenge.com/bonus/bounding-boxes.html
	type coordinate struct {
		x int
		y int
		z int
	}
	combined := shapes.NewGroup()
	combined.SetTransform(group.GetTransform())
	combined.SetMaterial(group.GetMaterial())
	cubedGroups := map[coordinate]shapes.Group{}

	// break shapes down into groups of equally spaced cubes
	for _, shape := range group.GetChildren() {

		sizeX := shape.Bounds().Max.X - shape.Bounds().Min.X
		sizeY := shape.Bounds().Max.Y - shape.Bounds().Min.Y
		sizeZ := shape.Bounds().Max.Z - shape.Bounds().Min.Z

		// if object too big, put on first level of main group,
		// .. to avoid increasing bounds too much for another group
		if sizeX >= width/2 || sizeY >= width/2 || sizeZ >= width/2 {
			combined.AddChild(shape)
			continue
		}

		// object will fit into about one cube
		// find containing square
		center := shape.Bounds().Center()
		coord := coordinate{
			x: int(center.X),
			y: int(center.Y),
			z: int(center.Z),
		}
		// get group for this cube
		cubedGroup, ok := cubedGroups[coord]
		if !ok {
			// create group if not exist
			cubedGroup = shapes.NewGroup()
			cubedGroups[coord] = cubedGroup
		}
		// add shape to it
		cubedGroup.AddChild(shape)
	}

	// add each cube to combined group
	for _, group := range cubedGroups {
		combined.AddChild(group)
	}
	return combined
}
