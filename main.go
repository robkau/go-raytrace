package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"go-raytrace/lib/colors"
	"go-raytrace/lib/geom"
	"go-raytrace/lib/patterns"
	"go-raytrace/lib/shapes"
	"go-raytrace/lib/view"
	"log"
	"math"
	"math/rand"
	"os"
	"time"
)

const (
	width = 600
)

func main() {
	rand.Seed(time.Now().UnixNano())

	var floor shapes.Shape = shapes.NewPlane()
	floor = floor.SetTransform(geom.Scale(10, 0.01, 10))
	m := floor.GetMaterial()
	m.Color = colors.NewColor(1, 0.9, 0.9)
	m.Pattern = patterns.NewGradientPattern(patterns.NewSolidColorPattern(colors.NewColor(1, 0.2, 0.4)), patterns.NewPositionAsColorPattern())
	m.Specular = 0
	floor = floor.SetMaterial(m)

	var leftWall shapes.Shape = shapes.NewPlane()
	leftWall = leftWall.SetTransform(geom.Translate(0, 0, 5).MulX4Matrix(geom.RotateY(-math.Pi / 4).MulX4Matrix(geom.RotateX(math.Pi / 2))))

	var rightWall shapes.Shape = shapes.NewPlane()
	rightWall = rightWall.SetTransform(geom.Translate(0, 0, 5).MulX4Matrix(geom.RotateY(math.Pi / 4).MulX4Matrix(geom.RotateX(math.Pi / 2))))

	var middle shapes.Shape = shapes.NewSphere()
	middle = middle.SetTransform(geom.Translate(-0.5, 1, 0.5))
	m = middle.GetMaterial()
	m.Pattern = patterns.NewPositionAsColorPattern()
	m.Color = colors.NewColor(0.1, 1, 0.5)
	m.Diffuse = 0.7
	m.Specular = 0.3
	middle = middle.SetMaterial(m)

	var right shapes.Shape = shapes.NewSphere()
	right = right.SetTransform(geom.Translate(1.5, 0.5, -0.5).MulX4Matrix(geom.Scale(0.5, 0.5, 0.5)))
	m = right.GetMaterial()
	m.Color = colors.NewColor(0.5, 1, 0.1)
	m.Diffuse = 0.7
	m.Specular = 0.3
	right = right.SetMaterial(m)

	var left shapes.Shape = shapes.NewSphere()
	left = left.SetTransform(geom.Translate(-1.5, 0.33, -0.75).MulX4Matrix(geom.Scale(0.33, 0.33, 0.33)))
	m = left.GetMaterial()
	m.Color = colors.NewColor(1, 0.8, 0.1)
	m.Diffuse = 0.7
	m.Specular = 0.3
	left = left.SetMaterial(m)

	w := view.NewWorld()
	w.AddObject(floor)
	w.AddObject(rightWall)
	w.AddObject(leftWall)
	w.AddObject(middle)
	w.AddObject(right)
	w.AddObject(left)

	w.AddLight(shapes.NewPointLight(geom.NewPoint(-10, 10, -10), colors.White()))

	c := view.NewCamera(width, width, math.Pi/3)
	c.Transform = geom.ViewTransform(geom.NewPoint(0, 1.5, -3),
		geom.NewPoint(0, 1, 0),
		geom.NewVector(0, 1, 0))

	g := &Game{
		c:      c,
		w:      w,
		canvas: view.NewCanvas(c.HSize, c.VSize),
		ng:     1,
		per:    10,
		last:   time.Now(),
	}

	go func() {
		for {
			pc := g.c.PixelChan(g.w, 8)
			for p := range pc {
				g.canvas.SetPixel(p.X, p.Y, p.C)
			}
			g.c.Transform = g.c.Transform.MulX4Matrix(geom.RotateY(math.Pi / 30))
		}

	}()

	f, err := os.Open("intro.wav")
	if err == nil {
		ac := audio.NewContext(44100)
		d, err := wav.Decode(ac, f)
		if err == nil {
			ap, err := audio.NewPlayer(ac, d)
			if err == nil {
				ap.Play()
			}
		}
	}

	ebiten.SetWindowSize(width, width)
	ebiten.SetWindowTitle("go-raytrace")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

type Game struct {
	count  int
	c      view.Camera
	w      view.World
	canvas view.Canvas

	ng   int
	last time.Time
	per  int
}

func (g *Game) Update() error {
	g.count++
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(ebiten.NewImageFromImage(g.canvas.ToImage()), op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}
