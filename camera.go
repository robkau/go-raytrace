package main

import "math"

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

func (c camera) render(w world) canvas {
	image := newCanvas(c.hSize, c.vSize)

	for y := 0; y < c.vSize-1; y++ {
		for x := 0; x < c.hSize-1; x++ {
			r := c.rayForPixel(x, y)
			c := w.colorAt(r)
			image.setPixel(x, y, c)
		}
	}
	return image
}
