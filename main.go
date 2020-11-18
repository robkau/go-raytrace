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
	width = 300
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
					g.imgRw.Lock()
					g.canvas.setPixel(p.x, p.y, p.c)
					g.imgRw.Unlock()
				}
				g.c.transform = g.c.transform.mulX4Matrix(rotateY(math.Pi / 30))
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
