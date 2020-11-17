package main

import (
	"math"
	"sync"
)

type camera struct {
	hSize      int
	vSize      int
	fov        float64
	halfWidth  float64
	halfHeight float64
	pixelSize  float64
	transform  x4Matrix
}

func newCamera(hs int, vs int, fov float64) camera {
	c := camera{
		hSize:     hs,
		vSize:     vs,
		fov:       fov,
		transform: newIdentityMatrixX4(),
	}

	halfView := math.Tan(c.fov / 2)
	aspect := float64(c.hSize) / float64(c.vSize)
	if aspect >= 1 {
		c.halfWidth = halfView
		c.halfHeight = halfView / aspect
	} else {
		c.halfWidth = halfView * aspect
		c.halfHeight = halfView
	}
	c.pixelSize = c.halfWidth * 2 / float64(c.hSize)

	return c
}

func (c camera) rayForPixel(px int, py int) ray {
	xOffset := (float64(px) + 0.5) * c.pixelSize
	yOffset := (float64(py) + 0.5) * c.pixelSize

	worldX := c.halfWidth - xOffset
	worldY := c.halfHeight - yOffset

	pixel := c.transform.invert().mulTuple(newPoint(worldX, worldY, -1))
	origin := c.transform.invert().mulTuple(newPoint(0, 0, 0))
	direction := pixel.sub(origin).normalize()

	return rayWith(origin, direction)
}

func (c camera) render(w world, numGoRoutines int) canvas {
	image := newCanvas(c.hSize, c.vSize)

	wg := sync.WaitGroup{}
	pixelsPerWorker := len(image.pixels) / numGoRoutines
	for i := 0; i < len(image.pixels); i += pixelsPerWorker {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := i; j < i+pixelsPerWorker && j < len(image.pixels); j++ {
				x := j % c.vSize
				y := j / c.vSize

				r := c.rayForPixel(x, y)
				c := w.colorAt(r)
				image.setPixel(x, y, c)
			}
		}(i)
	}
	wg.Wait()
	return image
}
