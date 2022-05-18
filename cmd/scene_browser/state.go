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
	"sync/atomic"
	"time"
)

// state struct implements ebiten.Game interface
type state struct {
	frameCount    int
	currentScene  int
	currentCamera int
	loc           scenes.CameraLocation
	scenes        []*scenes.Scene
	canvas        *view.Canvas
}

func start() *state {
	s := &state{
		scenes: scenes.LoadScenes(
			scenes.NewToriReplayScene,
			scenes.NewGroupTransformsScene,
			scenes.NewStoneGolemScene,
			scenes.NewWavyCarpetSpheres,
		),
		canvas: view.NewCanvas(width, width),
		loc: scenes.CameraLocation{
			At:        geom.NewPoint(2, 2, 2),
			LookingAt: geom.NewPoint(0, 0, 0),
		},
	}

	var rendered uint32 = 0
	var pixelsPerRenderStat uint32 = 100

	// render middle
	go func() {
		for {
			// todo: render option for # of rays picking random pixels & size to multiply color by
			// (list of pixels hit for each frame and last N frames, average by distance then time)

			s.loc = s.scenes[s.currentScene].Cs[s.currentCamera]

			pc, err := view.Render(s.scenes[s.currentScene].W, view.NewCameraAt(width, width, fov, s.loc.At, s.loc.LookingAt), 3, int(float64(runtime.NumCPU())/2), coordinate_supplier.Random)
			if err != nil {
				fmt.Println("failed create render")
				panic(err)
			}
			tLastRenderStat := time.Now()
			for p := range pc {
				if n := atomic.AddUint32(&rendered, 1); n%pixelsPerRenderStat == 0 {
					fmt.Printf("Writing %f pixels/sec\n", float64(pixelsPerRenderStat)/time.Since(tLastRenderStat).Seconds())
					tLastRenderStat = time.Now()
				}
				s.canvas.SetPixel(p.X, p.Y, p.C)
			}
			time.Sleep(24 * time.Hour)
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

	if inpututil.IsKeyJustPressed(ebiten.KeyN) {
		s.currentScene++
		if s.currentScene >= len(s.scenes) {
			s.currentScene = 0
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyM) {
		s.currentCamera++
		if s.currentCamera >= len(s.scenes[s.currentScene].Cs) {
			s.currentCamera = 0
		}
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
