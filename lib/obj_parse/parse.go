package obj_parse

import (
	"bufio"
	"fmt"
	"github.com/pkg/errors"
	"github.com/robkau/go-raytrace/lib/geom"
	"io"
	"os"
	"strconv"
	"strings"
)

// Parse reads a .obj format file and returns the triangles

func ParseFile(path string) (*ObjParsed, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "open file")
	}
	defer f.Close()

	return ParseReader(f)
}

func ParseReader(content io.Reader) (parsed *ObjParsed, err error) {
	parsed = newObjParsed()
	scanner := bufio.NewScanner(content)
	for scanner.Scan() {
		err = ParseLine(scanner.Text(), parsed)
		if err != nil {
			return nil, errors.Wrap(err, "failed parsing line")
		}
	}
	if err = scanner.Err(); err != nil {
		return nil, errors.Wrap(err, "failed scanning input")
	}

	return
}

func ParseLine(line string, parsed *ObjParsed) error {
	line = strings.TrimSpace(line)
	if strings.HasPrefix(line, "v") {
		return ParseVertexLine(line, parsed)
	}

	parsed.LinesIgnored++
	return nil
}

func ParseVertexLine(line string, parsed *ObjParsed) error {
	line = strings.TrimPrefix(line, "v")
	points := strings.Split(line, " ")
	if len(points) != 3 {
		return errors.New(fmt.Sprintf("expected 3 points, got %d points", len(points)))
	}

	p1, err := strconv.ParseFloat(points[0], 64)
	if err != nil {
		return errors.Wrap(err, "parsing vertex point 1")
	}
	p2, err := strconv.ParseFloat(points[1], 64)
	if err != nil {
		return errors.Wrap(err, "parsing vertex point 2")
	}
	p3, err := strconv.ParseFloat(points[2], 64)
	if err != nil {
		return errors.Wrap(err, "parsing vertex point 3")
	}

	parsed.AddVertex(geom.NewPoint(p1, p2, p3))
	return nil
}
