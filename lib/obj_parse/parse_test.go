package obj_parse

import (
	"github.com/robkau/go-raytrace/lib/geom"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func Test_ParseIgnoresUnknownLines(t *testing.T) {
	lines :=
		`There was...
        ... tr....
        ...set out\
	    in a re,,,,,
        the.....
`
	parsed, err := ParseReader(strings.NewReader(lines))

	require.NoError(t, err)
	require.Equal(t, 5, parsed.LinesIgnored)
}

func Test_ParseVertexLines(t *testing.T) {
	lines :=
		`v -1 1 0
         v -1.0000 0.5000 0.0000
         v 1 0 0
         v 1 1 0`

	parsed, err := ParseReader(strings.NewReader(lines))

	require.NoError(t, err)
	require.Equal(t, geom.NewPoint(-1, 1, 0), parsed.TriangleCoordinates(1))
	require.Equal(t, geom.NewPoint(-1, 0.5, 0), parsed.TriangleCoordinates(2))
	require.Equal(t, geom.NewPoint(1, 0, 0), parsed.TriangleCoordinates(3))
	require.Equal(t, geom.NewPoint(1, 1, 0), parsed.TriangleCoordinates(4))
}
