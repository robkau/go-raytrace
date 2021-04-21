package view

import (
	"go-raytrace/lib/colors"
	"go-raytrace/lib/geom"
	"math"
	"sync"
)

type Camera struct {
	HSize      int
	VSize      int
	fov        float64
	halfWidth  float64
	halfHeight float64
	pixelSize  float64
	Transform  geom.X4Matrix
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

func (c Camera) Render(w World, numGoRoutines int) Canvas {
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
				c := w.ColorAt(r)
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

func (c Camera) PixelChan(w World, numGoRoutines int) <-chan PixelInfo {
	pi := make(chan PixelInfo, numGoRoutines*8)

	cs := &coordinateSupplier{
		width:  c.HSize,
		height: c.VSize,
		rw:     sync.RWMutex{},
	}

	go func() {
		wg := sync.WaitGroup{}
		for i := 0; i < numGoRoutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for x, y := cs.next(); x != nextDone && y != nextDone; x, y = cs.next() {
					r := c.rayForPixel(x, y)
					c := w.ColorAt(r)

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
	x      int
	y      int
	width  int
	height int
	done   bool
	rw     sync.RWMutex
}

const nextDone = -1

func (c *coordinateSupplier) next() (x, y int) {
	c.rw.Lock()
	defer c.rw.Unlock()

	if c.done {
		return nextDone, nextDone
	}

	x = c.x
	y = c.y

	c.x++
	if c.x >= c.width {
		c.x = 0
		c.y++
	}
	if c.y >= c.height {
		c.done = true
	}
	return
}
