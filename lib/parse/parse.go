package parse

import (
	"github.com/pkg/errors"
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
		return parseReaderAsTori(content)
	default:
		panic("unknown format")
	}
}
