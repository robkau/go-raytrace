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
	"math/rand"
	"os"
	"time"
)

const (
	width = 555
)

func main() {
	rand.Seed(time.Now().UnixNano())

	var wall shapes.Shape = shapes.NewPlane()
	m := wall.GetMaterial()
	m.Pattern = patterns.NewSprayPaintPattern(patterns.NewCheckerPattern(patterns.NewSolidColorPattern(colors.NewColor(0.15, 0.15, 0.15)), patterns.NewSolidColorPattern(colors.NewColor(0.85, 0.85, 0.85))), 0.05)
	m.Ambient = 0.8
	m.Diffuse = 0.2
	m.Specular = 0
	wall = wall.SetMaterial(m)
	wall = wall.SetTransform(geom.Translate(0, 0, 10).MulX4Matrix(geom.RotateX(1.5708)))

	var ball shapes.Shape = shapes.NewGlassSphere()
	m = ball.GetMaterial()
	m.Color = colors.NewColor(1, 1, 1)
	m.Diffuse = 0
	m.Ambient = 0
	m.Specular = 0.9
	m.Shininess = 300
	m.Transparency = 0.9
	m.Reflective = 0.9
	m.RefractiveIndex = 1.5
	ball = ball.SetMaterial(m)

	var center shapes.Shape = shapes.NewGlassSphere()
	center = center.SetTransform(geom.Scale(0.5, 0.5, 0.5))
	m = center.GetMaterial()
	m.Color = colors.NewColor(1, 1, 1)
	m.Diffuse = 0
	m.Ambient = 0
	m.Specular = 0.9
	m.Shininess = 300
	m.Transparency = 0.9
	m.Reflective = 0.9
	m.RefractiveIndex = 1.0000034
	center = center.SetMaterial(m)

	w := view.NewWorld()
	w.AddObject(wall)
	w.AddObject(ball)
	w.AddObject(center)
	w.AddLight(shapes.NewPointLight(geom.NewPoint(2, 10, -5), colors.NewColor(0.9, 0.9, 0.9)))

	c := view.NewCamera(width, width, 0.45)
	c.Transform = geom.ViewTransform(geom.NewPoint(0, 0, -5),
		geom.NewPoint(0, 0, 0),
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
		// start rendering in background. draw one frame to canvas
		pc := g.c.PixelChan(g.w, 4, 4)
		for p := range pc {
			g.canvas.SetPixel(p.X, p.Y, p.C)
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
	// render frame in progress
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(ebiten.NewImageFromImage(g.canvas.ToImage()), op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}
