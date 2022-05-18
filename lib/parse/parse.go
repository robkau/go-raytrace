package parse

import (
	"github.com/pkg/errors"
	"github.com/robkau/go-raytrace/lib/shapes"
	"io"
	"os"
	"path/filepath"
)

func ParseObjFile(path string) (g shapes.Group, err error) {

	path, err = filepath.Abs(path)
	if err != nil {
		return nil, errors.Wrap(err, "calculate abs path")
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "open file")
	}
	defer f.Close()

	return ParseObjReader(f)
}

func ParseObjReader(content io.Reader) (g shapes.Group, e error) {
	o, err := parseReaderAsObj(content)
	if err != nil {
		return nil, errors.Wrap(err, "failed parsing")
	}
	return o.defaultGroup, nil
}
