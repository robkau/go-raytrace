package view

import (
	"github.com/robkau/go-raytrace/lib/colors"
	"github.com/robkau/go-raytrace/lib/geom"
	"math"
	"math/rand"
	"sync"
)

type Camera struct {
	HSize      int
	VSize      int
	fov        float64
	halfWidth  float64
	halfHeight float64
	pixelSize  float64
	Transform  *geom.X4Matrix
}

func NewCamera(hs int, vs int, fov float64) Camera {
	c := Camera{
		HSize:     hs,
		VSize:     vs,
		fov:       fov,
		Transform: geom.NewIdentityMatrixX4(),
	}

	halfView := math.Tan(c.fov / 2)
	aspect := float64(c.HSize) / float64(c.VSize)
	if aspect >= 1 {
		c.halfWidth = halfView
		c.halfHeight = halfView / aspect
	} else {
		c.halfWidth = halfView * aspect
		c.halfHeight = halfView
	}
	c.pixelSize = c.halfWidth * 2 / float64(c.HSize)

	return c
}

func (c Camera) rayForPixel(px int, py int) geom.Ray {
	xOffset := (float64(px) + 0.5) * c.pixelSize
	yOffset := (float64(py) + 0.5) * c.pixelSize

	worldX := c.halfWidth - xOffset
	worldY := c.halfHeight - yOffset

	pixel := c.Transform.Invert().MulTuple(geom.NewPoint(worldX, worldY, -1))
	origin := c.Transform.Invert().MulTuple(geom.NewPoint(0, 0, 0))
	direction := pixel.Sub(origin).Normalize()

	return geom.RayWith(origin, direction)
}

func (c Camera) Render(w World, rayBounces int, numGoRoutines int) *Canvas {
	image := NewCanvas(c.HSize, c.VSize)

	wg := sync.WaitGroup{}
	pixelsPerWorker := len(image.pixels) / numGoRoutines
	for i := 0; i < len(image.pixels); i += pixelsPerWorker {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := i; j < i+pixelsPerWorker && j < len(image.pixels); j++ {
				x := j % c.VSize
				y := j / c.VSize

				r := c.rayForPixel(x, y)
				c := w.ColorAt(r, rayBounces)
				image.SetPixel(x, y, c)
			}
		}(i)
	}
	wg.Wait()
	return image
}

type PixelInfo struct {
	X           int
	Y           int
	C           colors.Color
	LastInFrame bool
}

func (c Camera) PixelChan(w World, rayBounces int, numGoRoutines int, renderMode NextMode) <-chan PixelInfo {
	pi := make(chan PixelInfo, numGoRoutines*8)

	cs := newCoordinateSupplier(c.HSize, c.VSize, renderMode, false)

	go func() {
		wg := sync.WaitGroup{}
		for i := 0; i < numGoRoutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for x, y := cs.next(); x != nextDone && y != nextDone; x, y = cs.next() {
					r := c.rayForPixel(x, y)
					c := w.ColorAt(r, rayBounces)

					pi <- PixelInfo{
						X: x,
						Y: y,
						C: c,
					}
				}
			}()
		}
		wg.Wait()
		close(pi)
	}()

	return pi
}

type coordinateSupplier struct {
	// todo option for downscaling (hand out 2x2 squares with 1 random selected inside)
	coordinates []coordinate
	at          int
	repeat      bool
	mode        NextMode

	// todo replace with atomic mod?
	rw sync.RWMutex
}

const nextDone = -1

type coordinate struct {
	x int
	y int
}

func newCoordinateSupplier(width, height int, mode NextMode, repeat bool) *coordinateSupplier {
	cs := &coordinateSupplier{
		repeat: repeat,
		rw:     sync.RWMutex{},
	}

	switch mode {
	case Asc:
		cs.coordinates = makeAscCoordinates(width, height)
	case Random:
		cs.coordinates = makeAscCoordinates(width, height)
		rand.Shuffle(len(cs.coordinates), func(i, j int) { cs.coordinates[i], cs.coordinates[j] = cs.coordinates[j], cs.coordinates[i] })
	case Desc:
		cs.coordinates = makeAscCoordinates(width, height)
		i := 0
		j := len(cs.coordinates) - 1
		for i < j {
			cs.coordinates[i], cs.coordinates[j] = cs.coordinates[j], cs.coordinates[i]
			i++
			j--
		}
	default:
		panic("unknown mode")
	}

	return cs
}

func makeAscCoordinates(width, height int) []coordinate {
	coordinates := make([]coordinate, 0, width*height)
	var atX, atY int
	for {
		coordinates = append(coordinates, coordinate{
			x: atX,
			y: atY,
		})

		atX++
		if atX >= width {
			atX = 0
			atY++
		}
		if atY >= height {
			break
		}
	}
	return coordinates
}

func (c *coordinateSupplier) next() (x, y int) {
	c.rw.Lock()
	defer c.rw.Unlock()

	if c.at >= len(c.coordinates) {
		if c.repeat {
			c.at = 0
		} else {
			return nextDone, nextDone
		}
	}

	defer func() { c.at++ }()
	return c.coordinates[c.at].x, c.coordinates[c.at].y
}

type NextMode uint

const (
	Asc NextMode = iota
	Desc
	Random
)
