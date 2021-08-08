package obj_parse

import (
	"github.com/robkau/go-raytrace/lib/geom"
	"github.com/robkau/go-raytrace/lib/shapes"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

const (
	testObjFileGroups = `
		 v -1 1 0
         v -1 0 0
         v  1 0 0 
         v  1 1 0
         
         g FirstGroup
         f 1 2 3
         g SecondGroup
         f 1 3 4
`
)

func Test_ParseIgnoresUnknownLines(t *testing.T) {
	lines :=
		`There was...
        ... tr....
        ...set out\
	    in a re,,,,,
        the.....
`
	parsed, err := parseReader(strings.NewReader(lines))

	require.NoError(t, err)
	require.Equal(t, 5, parsed.linesIgnored)
}

func Test_ParseVertexLines(t *testing.T) {
	lines :=
		`v -1 1 0
         v -1.0000 0.5000 0.0000
         v 1 0 0
         v 1 1 0`

	parsed, err := parseReader(strings.NewReader(lines))
	require.NoError(t, err)

	t1, err := parsed.getVertex(1)
	require.NoError(t, err)
	t2, err := parsed.getVertex(2)
	require.NoError(t, err)
	t3, err := parsed.getVertex(3)
	require.NoError(t, err)
	t4, err := parsed.getVertex(4)
	require.NoError(t, err)

	require.Equal(t, geom.NewPoint(-1, 1, 0), t1)
	require.Equal(t, geom.NewPoint(-1, 0.5, 0), t2)
	require.Equal(t, geom.NewPoint(1, 0, 0), t3)
	require.Equal(t, geom.NewPoint(1, 1, 0), t4)
}

func Test_ParseTriangleFaces(t *testing.T) {
	lines :=
		`v -1 1 0
         v -1 0 0
         v 1 0 0
         v 1 1 0

		 f 1 2 3
         f 1 3 4
		`

	parsed, err := parseReader(strings.NewReader(lines))
	require.NoError(t, err)

	p1, err := parsed.getVertex(1)
	require.NoError(t, err)
	p2, err := parsed.getVertex(2)
	require.NoError(t, err)
	p3, err := parsed.getVertex(3)
	require.NoError(t, err)
	p4, err := parsed.getVertex(4)
	require.NoError(t, err)

	g := parsed.defaultGroup
	require.Len(t, g.GetChildren(), 2)
	require.Equal(t, p1, g.GetChildren()[0].(*shapes.Triangle).Vertices()[0])
	require.Equal(t, p2, g.GetChildren()[0].(*shapes.Triangle).Vertices()[1])
	require.Equal(t, p3, g.GetChildren()[0].(*shapes.Triangle).Vertices()[2])
	require.Equal(t, p1, g.GetChildren()[1].(*shapes.Triangle).Vertices()[0])
	require.Equal(t, p3, g.GetChildren()[1].(*shapes.Triangle).Vertices()[1])
	require.Equal(t, p4, g.GetChildren()[1].(*shapes.Triangle).Vertices()[2])
}

func Test_ParsePolygon_Triangulated(t *testing.T) {
	lines :=
		`v -1 1 0
         v -1 0 0
         v 1 0 0
         v 1 1 0
		 v 0 2 0

		 f 1 2 3 4 5
		`

	parsed, err := parseReader(strings.NewReader(lines))
	require.NoError(t, err)

	c := parsed.defaultGroup.GetChildren()
	require.Len(t, c, 3)

	t1 := parsed.defaultGroup.GetChildren()[0].(*shapes.Triangle)
	t2 := parsed.defaultGroup.GetChildren()[1].(*shapes.Triangle)
	t3 := parsed.defaultGroup.GetChildren()[2].(*shapes.Triangle)

	p1, err := parsed.getVertex(1)
	require.NoError(t, err)
	p2, err := parsed.getVertex(2)
	require.NoError(t, err)
	p3, err := parsed.getVertex(3)
	require.NoError(t, err)
	p4, err := parsed.getVertex(4)
	require.NoError(t, err)
	p5, err := parsed.getVertex(5)
	require.NoError(t, err)

	require.Equal(t, p1, t1.Vertices()[0])
	require.Equal(t, p2, t1.Vertices()[1])
	require.Equal(t, p3, t1.Vertices()[2])
	require.Equal(t, p1, t2.Vertices()[0])
	require.Equal(t, p3, t2.Vertices()[1])
	require.Equal(t, p4, t2.Vertices()[2])
	require.Equal(t, p1, t3.Vertices()[0])
	require.Equal(t, p4, t3.Vertices()[1])
	require.Equal(t, p5, t3.Vertices()[2])
}

func Test_ParseTriangles_Grouped(t *testing.T) {
	parsed, err := parseReader(strings.NewReader(testObjFileGroups))
	require.NoError(t, err)

	c := parsed.defaultGroup.GetChildren()
	require.Len(t, c, 2)

	gs := []shapes.Group{}
	for _, gr := range parsed.defaultGroup.GetChildren() {
		if grp, ok := gr.(shapes.Group); ok {
			gs = append(gs, grp)
		}
	}
	require.Len(t, gs, 2)
	require.Equal(t, parsed.namedGroups["FirstGroup"], gs[0])
	require.Equal(t, parsed.namedGroups["SecondGroup"], gs[1])

	t1 := gs[0].GetChildren()[0].(*shapes.Triangle)
	t2 := gs[1].GetChildren()[0].(*shapes.Triangle)

	p1, err := parsed.getVertex(1)
	require.NoError(t, err)
	p2, err := parsed.getVertex(2)
	require.NoError(t, err)
	p3, err := parsed.getVertex(3)
	require.NoError(t, err)
	p4, err := parsed.getVertex(4)
	require.NoError(t, err)

	require.Equal(t, p1, t1.Vertices()[0])
	require.Equal(t, p2, t1.Vertices()[1])
	require.Equal(t, p3, t1.Vertices()[2])
	require.Equal(t, p1, t2.Vertices()[0])
	require.Equal(t, p3, t2.Vertices()[1])
	require.Equal(t, p4, t2.Vertices()[2])
}

func Test_VertexNormals(t *testing.T) {
	f := `
		vn 0 0 1
		vn 0.707 0 -0.707
		vn 1 2 3
		`
	p, err := parseReader(strings.NewReader(f))
	require.NoError(t, err)

	n, err := p.getNormal(1)
	require.NoError(t, err)
	require.Equal(t, geom.NewVector(0, 0, 1), n)

	n, err = p.getNormal(2)
	require.NoError(t, err)
	require.Equal(t, geom.NewVector(0.707, 0, -0.707), n)

	n, err = p.getNormal(3)
	require.NoError(t, err)
	require.Equal(t, geom.NewVector(1, 2, 3), n)
}

func Test_FacesWithNormals(t *testing.T) {
	f := `
		v 0 1 0
		v -1 0 0
		v 1 0 0
		
		vn -1 0 0
		vn 1 0 0
		vn 0 1 0
		
		f 1//3 2//1 3//2
		f 1/0/3 2/102/1 3/14/2
		`
	p, err := parseReader(strings.NewReader(f))
	require.NoError(t, err)

	g := p.defaultGroup
	// parser should create smooth triangles when face normals are present
	t1 := g.GetChildren()[0].(*shapes.SmoothTriangle)
	t2 := g.GetChildren()[1].(*shapes.SmoothTriangle)

	v, err := p.getVertex(1)
	require.NoError(t, err)
	require.Equal(t, v, t1.Vertices()[0])

	v, err = p.getVertex(2)
	require.NoError(t, err)
	require.Equal(t, v, t1.Vertices()[1])

	v, err = p.getVertex(3)
	require.NoError(t, err)
	require.Equal(t, v, t1.Vertices()[2])

	ns := t1.Normals()
	n, err := p.getNormal(3)
	require.NoError(t, err)
	require.Equal(t, n, ns[0])

	n, err = p.getNormal(1)
	require.NoError(t, err)
	require.Equal(t, n, ns[1])

	n, err = p.getNormal(2)
	require.NoError(t, err)
	require.Equal(t, n, ns[2])

	require.Equal(t, t2, t1)
}
