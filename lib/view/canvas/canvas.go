package canvas

import (
	"archive/zip"
	"bufio"
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"github.com/robkau/go-raytrace/lib/colors"
	"image"
	gocolor "image/color"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
)

type Canvas struct {
	pixels []colors.Color
	width  int
	height int
	rw     sync.RWMutex
}

const (
	ppmFileHeader    = "P3"
	ppmMaxWidth      = 70
	ppmMinColorValue = 0
	ppmMaxColorValue = 255
)

func NewCanvas(width, height int) *Canvas {
	return &Canvas{
		pixels: make([]colors.Color, width*height),
		width:  width,
		height: height,
	}
}

func newCanvasWith(width int, height int, p colors.Color) *Canvas {
	c := &Canvas{
		pixels: make([]colors.Color, width*height),
		width:  width,
		height: height,
	}
	for i := range c.pixels {
		c.pixels[i] = p
	}
	return c
}

func (c *Canvas) GetSize() (width, height int) {
	return c.width, c.height
}

func (c *Canvas) GetPixel(x int, y int) colors.Color {
	c.rw.RLock()
	defer c.rw.RUnlock()
	return c._getPixel(x, y)
}

func (c *Canvas) _getPixel(x int, y int) colors.Color {
	return c.pixels[y*c.width+x]
}

func (c *Canvas) SetPixel(x int, y int, col colors.Color) {
	c.rw.Lock()
	defer c.rw.Unlock()
	c.pixels[y*c.width+x] = col
}

func (c *Canvas) ToImage() image.Image {
	upLeft := image.Point{0, 0}
	lowRight := image.Point{c.width, c.height}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})
	c.fillImage(img)
	return img
}

func (c *Canvas) fillImage(img *image.RGBA) {
	c.rw.RLock()
	defer c.rw.RUnlock()

	for x := 0; x < c.width; x++ {
		for y := 0; y < c.height; y++ {
			p := c._getPixel(x, y)
			img.Set(x, y, gocolor.RGBA{
				R: uint8(clamp(p.R*float64(ppmMaxColorValue), ppmMinColorValue, ppmMaxColorValue)),
				G: uint8(clamp(p.G*float64(ppmMaxColorValue), ppmMinColorValue, ppmMaxColorValue)),
				B: uint8(clamp(p.B*float64(ppmMaxColorValue), ppmMinColorValue, ppmMaxColorValue)),
				A: 0xff})
		}
	}
}

func (c *Canvas) toPPM() string {

	b := strings.Builder{}
	b.WriteString(fmt.Sprintf("%s\n", ppmFileHeader))
	b.WriteString(fmt.Sprintf("%d %d\n", c.width, c.height))
	b.WriteString(fmt.Sprintf("%d\n", ppmMaxColorValue))

	lineLength := 0
	for i, p := range c.pixels {

		rc := fmt.Sprintf("%d", clamp(p.R*float64(ppmMaxColorValue), ppmMinColorValue, ppmMaxColorValue))
		gc := fmt.Sprintf("%d", clamp(p.G*float64(ppmMaxColorValue), ppmMinColorValue, ppmMaxColorValue))
		bc := fmt.Sprintf("%d", clamp(p.B*float64(ppmMaxColorValue), ppmMinColorValue, ppmMaxColorValue))
		rcl := len(rc)
		gcl := len(gc)
		bcl := len(bc)

		if i%c.width == 0 && i != 0 {
			// finished writing a row
			b.WriteString("\n")
			lineLength = 0
		}

		// red
		if lineLength == 0 {
			// special case: first character of a line
			b.WriteString(rc)
			lineLength += rcl
		} else {
			if lineLength+rcl+1 > ppmMaxWidth {
				// not enough space in line. break it
				b.WriteString("\n")
				lineLength = 0
			} else {
				// continue line
				b.WriteString(" ")
				lineLength += 1
			}
			b.WriteString(rc)
			lineLength += rcl
		}

		// green
		if lineLength+1+gcl > ppmMaxWidth {
			// not enough space in line. break it
			b.WriteString("\n")
			lineLength = 0
		} else {
			// continue line
			b.WriteString(" ")
			lineLength += 1
		}
		b.WriteString(gc)
		lineLength += gcl

		// blue
		if lineLength+1+bcl > ppmMaxWidth {
			// not enough space in line. break it
			b.WriteString("\n")
			lineLength = 0
		} else {
			// continue line
			b.WriteString(" ")
			lineLength += 1
		}
		b.WriteString(bc)
		lineLength += bcl

	}
	b.WriteString("\n")
	return b.String()
}

func CanvasFromPPMZipFile(filepath string) (*Canvas, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("open filepath: %w", err)
	}
	defer f.Close() // todo me

	stat, err := f.Stat()
	zr, err := zip.NewReader(f, stat.Size())
	if err != nil {
		return nil, fmt.Errorf("open zip reader: %w", err)
	}

	if len(zr.File) != 1 {
		return nil, fmt.Errorf("expect one file inside zip but had %d", len(zr.File))
	}

	zf, err := zr.File[0].Open()
	if err != nil {
		return nil, fmt.Errorf("open zip file data: %w", err)
	}
	defer zf.Close() // todo me

	return CanvasFromPPMReader(zf)
}

func CanvasFromPPMFile(filepath string) (*Canvas, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("open filepath: %w", err)
	}
	defer f.Close() // todo me
	return CanvasFromPPMReader(f)
}

func CanvasFromPPMReader(r io.Reader) (*Canvas, error) {
	ppmBytes, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("read bytes: %w", err)
	}

	var c *Canvas
	var colorScale int
	var ps = &ppmCoordinateScanner{}
	var atX = 0
	var atY = 0
	lineNo := 0
	scanner := bufio.NewScanner(bytes.NewReader(ppmBytes))
	for scanner.Scan() {
		line := strings.TrimLeft(scanner.Text(), " ")
		if strings.HasPrefix(line, "#") {
			// skip comment lines entirely.
			continue
		}
		lineNo++
		if lineNo == 1 {
			// file header should be PPM magic bytes
			if line != fmt.Sprintf("%s", ppmFileHeader) {
				return nil, fmt.Errorf("invalid PPM header, should start with %s\\n", ppmFileHeader)
			}
			continue
		}
		if lineNo == 2 {
			// second line dictates width / height
			w, h, err := whFromPPMLine(line)
			if err != nil {
				return nil, fmt.Errorf("get ppm width and height: %w", err)
			}
			c = NewCanvas(w, h)
			continue
		}
		if lineNo == 3 {
			// third line dictates color scale
			colorScale, err = strconv.Atoi(line)
			if err != nil {
				return nil, fmt.Errorf("convert color scale to int: %w", err)
			}
			continue
		}

		// the rest of the lines dictate colors per-pixel
		lineVals := strings.Split(line, " ")
		for _, valS := range lineVals {
			if valS == "" {
				continue
			}
			val, err := strconv.Atoi(valS)
			if err != nil {
				return nil, fmt.Errorf("could not convert color val %s to int on line %d", valS, lineNo)
			}
			if built, col := ps.handleNextValue(val); built {
				c.SetPixel(atX, atY, col.MulBy(1./float64(colorScale)))
				atX++
				if atX >= c.width {
					atX = 0
					atY++
				}
				if atY > c.height {
					panic(fmt.Sprintf("ppm incremented to y value (%d) past height (%d)", atY, c.height))
				}
			}
		}

	}
	if err = scanner.Err(); err != nil {
		return nil, errors.Wrap(err, "scanning ppm data")
	}

	return c, nil
}

type ppmScannerState uint8

const (
	wantR ppmScannerState = iota
	wantG
	wantB
)

type ppmCoordinateScanner struct {
	r     int
	g     int
	b     int
	state ppmScannerState
}

func (p *ppmCoordinateScanner) handleNextValue(val int) (built bool, c colors.Color) {
	switch p.state {
	case wantR:
		p.r = val
		p.state = wantG
		return
	case wantG:
		p.g = val
		p.state = wantB
		return
	case wantB:
		defer p.toNextCoordinate()
		p.b = val
		return true, colors.NewColor(float64(p.r), float64(p.g), float64(p.b))
	default:
		panic(fmt.Sprintf("unmatched ppmCoordinateScanner state of %v", p.state))
	}
}

func (p *ppmCoordinateScanner) toNextCoordinate() {
	p.state = wantR
	p.r = 0
	p.g = 0
	p.b = 0
}

func whFromPPMLine(line string) (int, int, error) {
	parts := strings.Split(line, " ")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("wh line should contain two items but has %d items", len(parts))
	}

	width, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("convert width to int: %w", err)
	}

	height, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("convert height to int: %w", err)
	}

	return width, height, nil
}

func clamp(f float64, min int, max int) int {
	r := int(f + 0.5) // round it
	if r <= min {
		r = min
	}

	if r >= max {
		r = max
	}

	return r
}
