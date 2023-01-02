package view

import (
	"github.com/robkau/go-raytrace/lib/colors"
	"github.com/robkau/go-raytrace/lib/geom"
	"math"
	"sync"
)

type Camera struct {
	HSize      int
	VSize      int
	Fov        float64
	halfWidth  float64
	halfHeight float64
	pixelSize  float64
	Transform  geom.X4Matrix
}

func NewCamera(hs int, vs int, fov float64) Camera {
	c := Camera{
		HSize:     hs,
		VSize:     vs,
		Fov:       fov,
		Transform: geom.NewIdentityMatrixX4(),
	}

	halfView := math.Tan(c.Fov / 2)
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

func NewCameraAt(hs, vs int, fov float64, at geom.Tuple, lookingAt geom.Tuple) Camera {
	c := NewCamera(hs, vs, fov)
	c.Transform = geom.ViewTransform(at,
		lookingAt,
		geom.UpVector())
	return c
}

func (c Camera) rayForPixel(px int, py int) geom.Ray {
	xOffset := (float64(px) + 0.5) * c.pixelSize
	yOffset := (float64(py) + 0.5) * c.pixelSize

	worldX := c.halfWidth - xOffset
	worldY := c.halfHeight - yOffset

	pixel := c.Transform.Invert().MulTuple(geom.NewPoint(worldX, worldY, -1))
	origin := c.Transform.Invert().MulTuple(geom.ZeroPoint())
	direction := pixel.Sub(origin).Normalize()

	return geom.RayWith(origin, direction)
}

func (c Camera) Render(w *World, rayBounces int, numGoRoutines int) *Canvas {
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
	X int
	Y int
	C colors.Color
}
