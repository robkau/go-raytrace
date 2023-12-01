package canvas

import (
	"archive/zip"
	"bufio"
	"bytes"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/pkg/errors"
	"github.com/robkau/go-raytrace/lib/colors"
	"github.com/robkau/go-raytrace/lib/util"
	"image"
	gocolor "image/color"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
)

type Canvas struct {
	pixels *image.RGBA
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
		pixels: image.NewRGBA(image.Rect(0, 0, width, height)),
		width:  width,
		height: height,
	}
}

func newCanvasWith(width int, height int, p colors.Color) *Canvas {
	c := &Canvas{
		pixels: image.NewRGBA(image.Rect(0, 0, width, height)),
		width:  width,
		height: height,
	}
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			c.pixels.Set(x, y, p)
		}
	}
	return c
}

func (c *Canvas) GetSize() (width, height int) {
	return c.width, c.height
}

func (c *Canvas) GetPixel(x int, y int) colors.Color {
	c.rw.RLock()
	defer c.rw.RUnlock()
	return colors.NewColorFromStdlibColor(c.pixels.At(x, y))
}

func (c *Canvas) SetPixel(x int, y int, col colors.Color) {
	c.rw.Lock()
	defer c.rw.Unlock()
	// todo atomic pixels instead. ???

	c.pixels.Set(x, y, col)
}

func (c *Canvas) ToImage() image.Image {
	upLeft := image.Point{0, 0}
	lowRight := image.Point{c.width, c.height}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	for x := 0; x < c.width; x++ {
		for y := 0; y < c.height; y++ {
			p := c.GetPixel(x, y)
			r, g, b, a := p.RGBA()
			img.Set(x, y, gocolor.RGBA{
				R: uint8(r),
				G: uint8(g),
				B: uint8(b),
				A: uint8(a)})
		}
	}
	return img
}

func (c *Canvas) ToEbitenImage(previous *ebiten.Image) *ebiten.Image {
	if previous == nil {
		return ebiten.NewImageFromImage(c.ToImage())
	}

	previous.Clear()

	// protecting c.pixels
	c.rw.RLock()
	defer c.rw.RUnlock()

	previous.WritePixels(c.pixels.Pix)
	return previous
}

func (c *Canvas) Reset() {
	c.rw.Lock()
	defer c.rw.Unlock()
	//c.pixels.Pix = c.pixels.Pix[:0]
	for i, _ := range c.pixels.Pix {
		c.pixels.Pix[i] = 0
	}

}

func (c *Canvas) ToImagePaletted() *image.Paletted {
	palette := []gocolor.Color{
		gocolor.RGBA{0x00, 0x00, 0x00, 0xff},
		gocolor.RGBA{0x00, 0x00, 0xff, 0xff},
		gocolor.RGBA{0x00, 0xff, 0x00, 0xff},
		gocolor.RGBA{0x00, 0xff, 0xff, 0xff},
		gocolor.RGBA{0xff, 0x00, 0x00, 0xff},
		gocolor.RGBA{0xff, 0x00, 0xff, 0xff},
		gocolor.RGBA{0xff, 0xff, 0x00, 0xff},
		gocolor.RGBA{0xff, 0xff, 0xff, 0xff},
	}

	img := image.NewPaletted(image.Rect(0, 0, c.width, c.height), palette)

	for x := 0; x < c.width; x++ {
		for y := 0; y < c.height; y++ {
			p := c.GetPixel(x, y)
			r, g, b, a := p.RGBA()
			img.Set(x, y, gocolor.RGBA{
				R: uint8(r),
				G: uint8(g),
				B: uint8(b),
				A: uint8(a)})
		}
	}

	return img
}

func (c *Canvas) toPPM() string {

	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("%s\n", ppmFileHeader))
	sb.WriteString(fmt.Sprintf("%d %d\n", c.width, c.height))
	sb.WriteString(fmt.Sprintf("%d\n", ppmMaxColorValue))

	lineLength := 0
	for x := 0; x < c.width; x++ {
		for y := 0; y < c.height; y++ {
			col := c.pixels.At(x, y)
			r, g, b, _ := col.RGBA()

			// todo is this right after changing color repr
			rc := fmt.Sprintf("%d", util.Clamp(float64(r)*float64(ppmMaxColorValue), ppmMinColorValue, ppmMaxColorValue))
			gc := fmt.Sprintf("%d", util.Clamp(float64(g)*float64(ppmMaxColorValue), ppmMinColorValue, ppmMaxColorValue))
			bc := fmt.Sprintf("%d", util.Clamp(float64(b)*float64(ppmMaxColorValue), ppmMinColorValue, ppmMaxColorValue))
			rcl := len(rc)
			gcl := len(gc)
			bcl := len(bc)

			if y == c.height-1 {
				// finished writing a row
				sb.WriteString("\n")
				lineLength = 0
			}

			// red
			if lineLength == 0 {
				// special case: first character of a line
				sb.WriteString(rc)
				lineLength += rcl
			} else {
				if lineLength+rcl+1 > ppmMaxWidth {
					// not enough space in line. break it
					sb.WriteString("\n")
					lineLength = 0
				} else {
					// continue line
					sb.WriteString(" ")
					lineLength += 1
				}
				sb.WriteString(rc)
				lineLength += rcl
			}

			// green
			if lineLength+1+gcl > ppmMaxWidth {
				// not enough space in line. break it
				sb.WriteString("\n")
				lineLength = 0
			} else {
				// continue line
				sb.WriteString(" ")
				lineLength += 1
			}
			sb.WriteString(gc)
			lineLength += gcl

			// blue
			if lineLength+1+bcl > ppmMaxWidth {
				// not enough space in line. break it
				sb.WriteString("\n")
				lineLength = 0
			} else {
				// continue line
				sb.WriteString(" ")
				lineLength += 1
			}
			sb.WriteString(bc)
			lineLength += bcl
		}
	}

	sb.WriteString("\n")
	return sb.String()
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
