package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"go-raytrace/lib/view"
	"go-raytrace/scenes"
	"log"
	"os"
	"runtime"
)

const (
	width = 555
)

func main() {
	w, c := scenes.NewRoomScene(width)

	g := &Game{
		c:      c,
		w:      w,
		canvas: view.NewCanvas(c.HSize, c.VSize),
	}

	go func() {
		// start rendering in background. draw one frame to canvas
		pc := g.c.PixelChan(g.w, 5, runtime.NumCPU()/2)
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

// Game struct implements ebiten.Game
type Game struct {
	count  int
	c      view.Camera
	w      view.World
	canvas view.Canvas
}

func (g *Game) Update() error {
	g.count++
	// not updating anything here. canvas is updated in background goroutine.
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// render current frame progress
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(ebiten.NewImageFromImage(g.canvas.ToImage()), op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}
