package view

import (
	"github.com/stretchr/testify/assert"
	"go-raytrace/lib/colors"
	"go-raytrace/lib/geom"
	"math"
	"testing"
)

func Test_NewCamera(t *testing.T) {
	c := NewCamera(160, 120, math.Pi/2)

	assert.Equal(t, 160, c.HSize)
	assert.Equal(t, 120, c.VSize)
	assert.Equal(t, math.Pi/2, c.fov)
	assert.Equal(t, geom.NewIdentityMatrixX4(), c.Transform)
}

func Test_PixelSize_Horizontal(t *testing.T) {
	c := NewCamera(200, 125, math.Pi/2)

	assert.Equal(t, 0.01, c.pixelSize)
}

func Test_PixelSize_Vertical(t *testing.T) {
	c := NewCamera(125, 200, math.Pi/2)

	assert.Equal(t, 0.01, c.pixelSize)
}

func Test_RayThrough_Center(t *testing.T) {
	c := NewCamera(201, 101, math.Pi/2)

	r := c.rayForPixel(100, 50)

	assert.Equal(t, geom.NewPoint(0, 0, 0), r.Origin)
	assert.Equal(t, geom.NewVector(0, 0, -1), r.Direction)

}

func Test_RayThrough_Corner(t *testing.T) {
	c := NewCamera(201, 101, math.Pi/2)

	r := c.rayForPixel(0, 0)

	assert.Equal(t, geom.NewPoint(0, 0, 0), r.Origin)
	assert.Equal(t, geom.NewVector(0.66519, 0.33259, -0.66851).RoundTo(5), r.Direction.RoundTo(5))

}

func Test_RayThrough_CameraTransformed(t *testing.T) {
	c := NewCamera(201, 101, math.Pi/2)
	c.Transform = geom.RotateY(math.Pi / 4).MulX4Matrix(geom.Translate(0, -2, 5))

	r := c.rayForPixel(100, 50)

	assert.Equal(t, geom.NewPoint(0, 2, -5), r.Origin)
	assert.Equal(t, geom.NewVector(math.Sqrt(2)/2, 0, -math.Sqrt(2)/2).RoundTo(5), r.Direction.RoundTo(5))
}

func Test_RenderWorld(t *testing.T) {
	w := defaultWorld()
	c := NewCamera(11, 11, math.Pi/2)
	from := geom.NewPoint(0, 0, -5)
	to := geom.NewPoint(0, 0, 0)
	up := geom.NewVector(0, 1, 0)
	c.Transform = geom.ViewTransform(from, to, up)

	image := c.render(w, 1)

	assert.Equal(t, colors.NewColor(0.38066, 0.47583, 0.2855), image.getPixel(5, 5).RoundTo(5))
}
