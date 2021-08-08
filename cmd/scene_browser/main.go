package main

import (
	"flag"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/robkau/go-raytrace/cmd/scene_browser/scenes"
	"github.com/robkau/go-raytrace/lib/geom"
	"github.com/robkau/go-raytrace/lib/view"
	"log"
	"math"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

const (
	width = 444
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	w, c := scenes.NewTeapotScene(width)

	g := &Game{
		c:      c,
		w:      w,
		canvas: view.NewCanvas(c.HSize, c.VSize),
	}

	go func() {
		// start rendering in background. draw frames to canvas
		for {
			start := time.Now()
			pc := g.c.PixelChan(g.w, 5, int(float64(runtime.NumCPU())/1.5), view.Asc)
			for p := range pc {
				g.canvas.SetPixel(p.X, p.Y, p.C)
			}
			elapsed := time.Now().Sub(start)
			fmt.Printf("Took %d to render\n", elapsed)

			return

			// update camera between frames
			g.c.Transform = g.c.Transform.MulX4Matrix(geom.RotateY(math.Pi / 19))

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
				defer ap.Close()
			}

		}
	}
	defer f.Close()

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
	canvas *view.Canvas
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
