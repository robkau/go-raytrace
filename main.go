package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"log"
	"math"
	"runtime"
	"sync"
	"time"
)

const (
	width = 100
)

func main() {

	var floor shape = newPlane()
	floor = floor.setTransform(scale(10, 0.01, 10))
	m := floor.getMaterial()
	m.color = color{1, 0.9, 0.9}
	m.specular = 0
	floor = floor.setMaterial(m)

	var leftWall shape = newPlane()
	leftWall = leftWall.setTransform(translate(0, 0, 5).mulX4Matrix(rotateY(-math.Pi / 4).mulX4Matrix(rotateX(math.Pi / 2))))
	leftWall = leftWall.setMaterial(floor.getMaterial())

	var rightWall shape = newPlane()
	rightWall = rightWall.setTransform(translate(0, 0, 5).mulX4Matrix(rotateY(math.Pi / 4).mulX4Matrix(rotateX(math.Pi / 2))))
	rightWall.setMaterial(floor.getMaterial())

	var middle shape = newSphere()
	middle = middle.setTransform(translate(-0.5, 1, 0.5))
	m = middle.getMaterial()
	m.color = color{0.1, 1, 0.5}
	m.diffuse = 0.7
	m.specular = 0.3
	middle = middle.setMaterial(m)

	var right shape = newSphere()
	right = right.setTransform(translate(1.5, 0.5, -0.5).mulX4Matrix(scale(0.5, 0.5, 0.5)))
	m = right.getMaterial()
	m.color = color{0.5, 1, 0.1}
	m.diffuse = 0.7
	m.specular = 0.3
	right = right.setMaterial(m)

	var left shape = newSphere()
	left = left.setTransform(translate(-1.5, 0.33, -0.75).mulX4Matrix(scale(0.33, 0.33, 0.33)))
	m = left.getMaterial()
	m.color = color{1, 0.8, 0.1}
	m.diffuse = 0.7
	m.specular = 0.3
	left = left.setMaterial(m)

	w := newWorld()
	w.addObject(floor)
	w.addObject(rightWall)
	w.addObject(leftWall)
	w.addObject(middle)
	w.addObject(right)
	w.addObject(left)

	w.addLight(newPointLight(newPoint(-10, 10, -10), color{1, 1, 1}))

	c := newCamera(width, width, math.Pi/3)
	c.transform = viewTransform(newPoint(0, 1.5, -3),
		newPoint(0, 1, 0),
		newVector(0, 1, 0))

	g := &Game{
		c:      c,
		w:      w,
		canvas: newCanvas(c.hSize, c.vSize),
		ng:     1,
		per:    10,
		last:   time.Now(),
	}

	ebiten.SetWindowSize(width, width)
	ebiten.SetWindowTitle("go-raytrace")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

type Game struct {
	count  int
	c      camera
	w      world
	canvas canvas

	ng   int
	last time.Time
	per  int

	imgRw sync.RWMutex
}

func (g *Game) Update() error {
	g.count++

	if g.count == 1 {
		go func() {
			for {
				pi := g.c.pixelChan(g.w, runtime.NumCPU())
				for p := range pi {
					// todo buffer locally and lock less
					// maybe maintain two canvases. update a different one each frame. swap which canvas to display each frame.
					g.imgRw.Lock()
					g.canvas.setPixel(p.x, p.y, p.c)
					g.imgRw.Unlock()
				}
				g.imgRw.Lock()
				g.c.transform = g.c.transform.mulX4Matrix(rotateY(math.Pi / 30))
				g.imgRw.Unlock()
			}
		}()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.imgRw.Lock()
	defer g.imgRw.Unlock()
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(ebiten.NewImageFromImage(g.canvas.toImage()), op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}
