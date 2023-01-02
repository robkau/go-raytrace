package parse

import (
	"bufio"
	"fmt"
	"github.com/pkg/errors"
	"github.com/robkau/go-raytrace/lib/geom"
	"github.com/robkau/go-raytrace/lib/shapes"
	"io"
	"strconv"
	"strings"
)

func parseReaderAsObj(content io.Reader) (parsed *objParsed, err error) {
	parsed = newObjParsed()
	scanner := bufio.NewScanner(content)
	for scanner.Scan() {
		err = parseLine(scanner.Text(), parsed)
		if err != nil {
			return nil, errors.Wrap(err, "failed parsing line")
		}
	}
	if err = scanner.Err(); err != nil {
		return nil, errors.Wrap(err, "failed scanning input")
	}

	return
}

func parseLine(line string, parsed *objParsed) error {
	line = strings.TrimSpace(line)
	if strings.HasPrefix(line, "v ") {
		return parseVertexLine(line, parsed)
	}
	if strings.HasPrefix(line, "vn ") {
		return parseVertexNormalLine(line, parsed)
	}
	if strings.HasPrefix(line, "f ") {
		return parseFaceLine(line, parsed)
	}
	if strings.HasPrefix(line, "g ") {
		return parseGroupLine(line, parsed)
	}

	parsed.linesIgnored++
	return nil
}

func parseVertexLine(line string, parsed *objParsed) error {
	line = strings.TrimPrefix(line, "v ")
	points := strings.FieldsFunc(line,
		func(c rune) bool {
			return c == ' '
		})
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

	parsed.addVertex(geom.NewPoint(p1, p2, p3))
	return nil
}

func parseVertexNormalLine(line string, parsed *objParsed) error {
	line = strings.TrimPrefix(line, "vn ")
	points := strings.FieldsFunc(line,
		func(c rune) bool {
			return c == ' '
		})
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

	parsed.addNormal(geom.NewVector(p1, p2, p3))
	return nil
}

func parseFaceLine(line string, parsed *objParsed) error {
	line = strings.TrimPrefix(line, "f ")
	indexes := strings.FieldsFunc(line,
		func(c rune) bool {
			return c == ' '
		})

	if len(indexes) < 3 {
		return errors.New("3 or more indexes to make a face")
	}

	smooth := false
	for i, _ := range indexes {
		iParts := strings.Split(indexes[i], "/")
		if len(iParts) == 3 {
			smooth = true
		}
	}

	// todo refactorm e
	var triangles []shapes.Shape
	if len(indexes) < 3 {
		return errors.New("must provide 3 or more vertexes for a face")
	}
	if len(indexes) == 3 {
		// triangles
		if smooth {
			i1, err := strconv.Atoi(strings.Split(indexes[0], "/")[0])
			if err != nil {
				return errors.Wrap(err, "parsing vertex index 1")
			}
			i2, err := strconv.Atoi(strings.Split(indexes[1], "/")[0])
			if err != nil {
				return errors.Wrap(err, "parsing vertex index 2")
			}
			i3, err := strconv.Atoi(strings.Split(indexes[2], "/")[0])
			if err != nil {
				return errors.Wrap(err, "parsing vertex index 3")
			}
			p1, err := parsed.getVertex(i1)
			if err != nil {
				return errors.Wrap(err, "getting vertex index 1")
			}
			p2, err := parsed.getVertex(i2)
			if err != nil {
				return errors.Wrap(err, "getting vertex index 2")
			}
			p3, err := parsed.getVertex(i3)
			if err != nil {
				return errors.Wrap(err, "getting vertex index 3")
			}
			in1, err := strconv.Atoi(strings.Split(indexes[0], "/")[2])
			if err != nil {
				return errors.Wrap(err, "parsing normal index 1")
			}
			in2, err := strconv.Atoi(strings.Split(indexes[1], "/")[2])
			if err != nil {
				return errors.Wrap(err, "parsing normal index 2")
			}
			in3, err := strconv.Atoi(strings.Split(indexes[2], "/")[2])
			if err != nil {
				return errors.Wrap(err, "parsing normal index 3")
			}
			n1, err := parsed.getNormal(in1)
			if err != nil {
				return errors.Wrap(err, "getting normal index 1")
			}
			n2, err := parsed.getNormal(in2)
			if err != nil {
				return errors.Wrap(err, "getting normal index 2")
			}
			n3, err := parsed.getNormal(in3)
			if err != nil {
				return errors.Wrap(err, "getting normal index 3")
			}
			triangles = append(triangles, shapes.NewSmoothTriangle(p1, p2, p3, n1, n2, n3))
		} else {
			i1, err := strconv.Atoi(indexes[0])
			if err != nil {
				return errors.Wrap(err, "parsing vertex index 1")
			}
			i2, err := strconv.Atoi(indexes[1])
			if err != nil {
				return errors.Wrap(err, "parsing vertex index 2")
			}
			i3, err := strconv.Atoi(indexes[2])
			if err != nil {
				return errors.Wrap(err, "parsing vertex index 3")
			}
			p1, err := parsed.getVertex(i1)
			if err != nil {
				return errors.Wrap(err, "getting vertex index 1")
			}
			p2, err := parsed.getVertex(i2)
			if err != nil {
				return errors.Wrap(err, "getting vertex index 2")
			}
			p3, err := parsed.getVertex(i3)
			if err != nil {
				return errors.Wrap(err, "getting vertex index 3")
			}
			triangles = append(triangles, shapes.NewTriangle(p1, p2, p3))
		}

	} else {
		// polygons - decompose into multiple triangles
		vs := []geom.Tuple{}
		var ns []geom.Tuple
		for i := 0; i < len(indexes); i++ {
			if smooth {
				// split line
				parts := strings.Split(indexes[i], "/")
				if len(parts) != 3 {
					return errors.New("unknown parts syntax")
				}

				iv, err := strconv.Atoi(parts[0])
				if err != nil {
					return errors.Wrap(err, "parsing vertex index")
				}
				in, err := strconv.Atoi(parts[2])
				if err != nil {
					return errors.Wrap(err, "parsing vertex normal index")
				}
				v, err := parsed.getVertex(iv)
				if err != nil {
					return errors.Wrap(err, "getting vertex by index")
				}
				n, err := parsed.getNormal(in)
				if err != nil {
					return errors.Wrap(err, "getting normal by index")
				}
				vs = append(vs, v)
				ns = append(ns, n)
			} else {
				iv, err := strconv.Atoi(indexes[i])
				if err != nil {
					return errors.Wrap(err, "parsing vertex index")
				}
				v, err := parsed.getVertex(iv)
				if err != nil {
					return errors.Wrap(err, "getting vertex by index")
				}
				vs = append(vs, v)
			}

		}
		// triangulate face
		triangles = fanTriangulation(vs, ns)
	}

	for i := 0; i < len(triangles); i++ {
		parsed.addShape(triangles[i])
	}
	return nil
}

// triangulates a polygon defined by vertexes into component triangle
// if normals is nil, regular triangles are returned
// if normals is not nil, smooth triangles are returned. normals should be same length as vertexes, with a matching normal for each vertex
func fanTriangulation(vertexes, normals []geom.Tuple) (triangles []shapes.Shape) {
	for i := 1; i <= len(vertexes)-2; i++ {
		if len(normals) > 0 {
			// todo this doesnt look any different when rendered?
			t := shapes.NewSmoothTriangle(vertexes[0], vertexes[i], vertexes[i+1], normals[0], normals[i], normals[i+1])
			triangles = append(triangles, t)
		} else {
			t := shapes.NewTriangle(vertexes[0], vertexes[i], vertexes[i+1])
			triangles = append(triangles, t)
		}
	}
	return
}

func parseGroupLine(line string, parsed *objParsed) error {
	line = strings.TrimPrefix(line, "g ")
	indexes := strings.FieldsFunc(line,
		func(c rune) bool {
			return c == ' '
		})

	if len(indexes) < 1 {
		return errors.New("should have a group name")
	}
	parsed.currentNamedGroup = indexes[0]
	return nil
}
