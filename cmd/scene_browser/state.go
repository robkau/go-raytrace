package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/robkau/coordinate_supplier"
	"github.com/robkau/go-raytrace/cmd/scene_browser/scenes"
	"github.com/robkau/go-raytrace/lib/geom"
	"github.com/robkau/go-raytrace/lib/view"
	"math"
	"runtime"
)

// state struct implements ebiten.Game interface
type state struct {
	frameCount   int
	currentScene int
	loc          scenes.CameraLocation
	scenes       []*scenes.Scene
	canvas       *view.Canvas
}

func start() *state {
	s := &state{
		scenes: scenes.LoadScenes(
			scenes.NewWavyCarpetSpheres,
			//scenes.NewGroupTransformsScene,

		),
		canvas: view.NewCanvas(width, width),
		loc: scenes.CameraLocation{
			At:        geom.NewPoint(2, 2, 2),
			LookingAt: geom.NewPoint(0, 0, 0),
		},
	}

	// render middle
	go func() {
		for {
			// todo: render option for # of rays picking random pixels & size to multiply color by
			// (list of pixels hit for each frame and last N frames, average by distance then time)

			pc, err := view.Render(s.scenes[s.currentScene].W, view.NewCameraAt(width, width, fov, s.loc.At, s.loc.LookingAt), 2, int(float64(runtime.NumCPU())), coordinate_supplier.Asc)
			if err != nil {
				fmt.Println("failed create render")
				panic(err)
			}
			for p := range pc {
				s.canvas.SetPixel(p.X, p.Y, p.C)
			}
			//g.c.Transform = g.c.Transform.MulX4Matrix(geom.RotateY(math.Pi / 19))
		}
	}()

	return s
}

func (s *state) Update() error {
	s.frameCount++

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		s.loc.RotateAroundY(math.Pi / 12)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		s.loc.RotateAroundY(-math.Pi / 12)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		s.loc.At = s.loc.At.Mul(0.9)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		s.loc.At = s.loc.At.Mul(1.1)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyC) {
		s.loc.At.Y++
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyZ) {
		s.loc.At.Y--
	}

	// canvas is updated in background goroutine.
	return nil
}

func (s *state) Draw(screen *ebiten.Image) {
	// render current frame progress
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(ebiten.NewImageFromImage(s.canvas.ToImage()), op)
}

func (s *state) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}
