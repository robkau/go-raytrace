package parse

import (
	"bufio"
	"fmt"
	"github.com/robkau/go-raytrace/lib/geom"
	"github.com/robkau/go-raytrace/lib/shapes"
	"io"
	"math"
	"strconv"
	"strings"
)

func parseReaderAsTori(content io.Reader) (shapes.Group, error) {
	g := shapes.NewGroup()

	scanner := bufio.NewScanner(content)
	for scanner.Scan() {
		pos := parsePositionLine(scanner.Text())
		for _, p := range pos {
			sp := shapes.NewSphere()
			sp.SetTransform(geom.RotateX(-math.Pi / 2).MulX4Matrix(geom.Translate(p.X, p.Y, p.Z)).MulX4Matrix(geom.Scale(0.14, 0.14, 0.14)))
			g.AddChild(sp)
		}
	}

	return g, nil
}

func parsePositionLine(positionLine string) []geom.Tuple {
	// get positions
	s := strings.Split(positionLine, "; ")

	// split positions
	ss := strings.Split(s[1], " ")

	// position strings to float
	// iterate in chunks of 3 (XYZ)
	ts := make([]geom.Tuple, 0)
	for i := 0; i < len(ss)-2; i += 3 {
		x, err := strconv.ParseFloat(ss[i], 64)
		if err != nil {
			panic(fmt.Sprintf("invalid pos float: %s", err.Error()))
		}
		y, err := strconv.ParseFloat(ss[i+1], 64)
		if err != nil {
			panic(fmt.Sprintf("invalid pos float: %s", err.Error()))
		}
		z, err := strconv.ParseFloat(ss[i+2], 64)
		if err != nil {
			panic(fmt.Sprintf("invalid pos float: %s", err.Error()))
		}
		ts = append(ts, geom.NewPoint(x, y, z))
	}

	return ts
}
