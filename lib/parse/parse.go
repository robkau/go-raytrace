package parse

import (
	"github.com/pkg/errors"
	"github.com/robkau/go-raytrace/lib/colors"
	"github.com/robkau/go-raytrace/lib/materials"
	"github.com/robkau/go-raytrace/lib/patterns"
	"github.com/robkau/go-raytrace/lib/shapes"
	"io"
	"os"
	"path/filepath"
)

type SupportedFormat uint

const (
	Obj SupportedFormat = iota
	Tori
)

func ParseFile(path string, as SupportedFormat) (g shapes.Group, err error) {

	path, err = filepath.Abs(path)
	if err != nil {
		return nil, errors.Wrap(err, "calculate abs path")
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "open file")
	}
	defer f.Close()

	return ParseReader(f, as)
}

func ParseReader(content io.Reader, as SupportedFormat) (g shapes.Group, e error) {
	switch as {
	case Obj:
		o, err := parseReaderAsObj(content)
		if err != nil {
			return nil, errors.Wrap(err, "failed parsing")
		}
		return o.defaultGroup, nil
	case Tori:
		pr, err := parseReaderAsTori(content)

		p0 := pr.p0AllPositions()
		m := materials.NewMaterial()
		m.Pattern = patterns.NewSolidColorPattern(colors.Blue())
		for _, c := range p0.GetChildren() {
			c.SetMaterial(m)
		}

		p1 := pr.p1AllPositions()
		m = materials.NewMaterial()
		m.Pattern = patterns.NewSolidColorPattern(colors.Red())
		for _, c := range p1.GetChildren() {
			c.SetMaterial(m)

			p0.AddChild(c)
		}

		return p0, err
	default:
		panic("unknown format")
	}
}
