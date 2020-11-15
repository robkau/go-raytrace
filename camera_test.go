package main

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func Test_NewCamera(t *testing.T) {
	c := newCamera(160, 120, math.Pi/2)

	assert.Equal(t, 160, c.hSize)
	assert.Equal(t, 120, c.vSize)
	assert.Equal(t, math.Pi/2, c.fov)
	assert.Equal(t, newIdentityMatrixX4(), c.transform)
}

func Test_PixelSize_Horizontal(t *testing.T) {
	c := newCamera(200, 125, math.Pi/2)

	assert.Equal(t, 0.01, c.pixelSize)
}

func Test_PixelSize_Vertical(t *testing.T) {
	c := newCamera(125, 200, math.Pi/2)

	assert.Equal(t, 0.01, c.pixelSize)
}

func Test_RayThrough_Center(t *testing.T) {
	c := newCamera(201, 101, math.Pi/2)

	r := c.rayForPixel(100, 50)

	assert.Equal(t, newPoint(0, 0, 0), r.origin)
	assert.Equal(t, newVector(0, 0, -1), r.direction)

}

func Test_RayThrough_Corner(t *testing.T) {
	c := newCamera(201, 101, math.Pi/2)

	r := c.rayForPixel(0, 0)

	assert.Equal(t, newPoint(0, 0, 0), r.origin)
	assert.Equal(t, newVector(0.66519, 0.33259, -0.66851).roundTo(5), r.direction.roundTo(5))

}

func Test_RayThrough_CameraTransformed(t *testing.T) {
	c := newCamera(201, 101, math.Pi/2)
	c.transform = rotateY(math.Pi / 4).mulX4Matrix(translate(0, -2, 5))

	r := c.rayForPixel(100, 50)

	assert.Equal(t, newPoint(0, 2, -5), r.origin)
	assert.Equal(t, newVector(math.Sqrt(2)/2, 0, -math.Sqrt(2)/2).roundTo(5), r.direction.roundTo(5))
}

func Test_RenderWorld(t *testing.T) {
	w := defaultWorld()
	c := newCamera(11, 11, math.Pi/2)
	from := newPoint(0, 0, -5)
	to := newPoint(0, 0, 0)
	up := newVector(0, 1, 0)
	c.transform = viewTransform(from, to, up)

	image := c.render(w)

	assert.Equal(t, color{0.38066, 0.47583, 0.2855}, image.getPixel(5, 5).roundTo(5))
}
