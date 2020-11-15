package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"log"
	"math"
)

const (
	width  = 400
	height = 400
)

func main() {

	floor := newSphere()
	floor = floor.setTransform(scale(10, 0.01, 10))
	floor.m.color = color{1, 0.9, 0.9}
	floor.m.specular = 0

	leftWall := newSphere()
	leftWall = leftWall.setTransform(translate(0, 0, 5).mulX4Matrix(rotateY(-math.Pi / 4).mulX4Matrix(rotateX(math.Pi / 2).mulX4Matrix(scale(10, 0.01, 10)))))
	leftWall.m = floor.m

	rightWall := newSphere()
	rightWall = rightWall.setTransform(translate(0, 0, 5).mulX4Matrix(rotateY(math.Pi / 4).mulX4Matrix(rotateX(math.Pi / 2).mulX4Matrix(scale(10, 0.01, 10)))))
	rightWall.m = floor.m

	middle := newSphere()
	middle = middle.setTransform(translate(-0.5, 1, 0.5))
	middle.m.color = color{0.1, 1, 0.5}
	middle.m.diffuse = 0.7
	middle.m.specular = 0.3

	right := newSphere()
	right = right.setTransform(translate(1.5, 0.5, -0.5).mulX4Matrix(scale(0.5, 0.5, 0.5)))
	right.m.color = color{0.5, 1, 0.1}
	right.m.diffuse = 0.7
	right.m.specular = 0.3

	left := newSphere()
	left = left.setTransform(translate(-1.5, 0.33, -0.75).mulX4Matrix(scale(0.33, 0.33, 0.33)))
	left.m.color = color{1, 0.8, 0.1}
	left.m.diffuse = 0.7
	left.m.specular = 0.3

	w := newWorld()
	w.addObject(floor)
	w.addObject(rightWall)
	w.addObject(leftWall)
	w.addObject(middle)
	w.addObject(right)
	w.addObject(left)

	w.addLight(newPointLight(newPoint(-10, 10, -10), color{1, 1, 1}))

	c := newCamera(width, height, math.Pi/3)
	c.transform = viewTransform(newPoint(0, 1.5, -5),
		newPoint(0, 1, 0),
		newVector(0, 1, 0))

	cnv := c.render(w)
	runnerImage = ebiten.NewImageFromImage(cnv.toImage())

	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("go-raytrace")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}

// ebiten stuff
var (
	runnerImage *ebiten.Image
)

type Game struct {
	count int
}

func (g *Game) Update() error {
	g.count++
	//runnerImage = ebiten.NewImageFromImage(canvas.toImage())
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(runnerImage, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}
