package main

import (
	"context"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/robkau/coordinate_supplier"
	"github.com/robkau/go-raytrace/cmd/scene_browser/scenes"
	"github.com/robkau/go-raytrace/lib/geom"
	"github.com/robkau/go-raytrace/lib/view"
	"log"
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
	loc           *scenes.CameraLocation
	scenes        []*scenes.Scene
	canvas        *view.Canvas

	cancel context.CancelFunc
}

func start() *state {
	s := &state{
		scenes: scenes.LoadScenes(
			scenes.NewGroupTransformsScene,
			scenes.NewToriReplayScene,
			scenes.NewStoneGolemScene,
			scenes.NewWavyCarpetSpheres,
			scenes.NewTeapotScene,
			scenes.NewPondScene,
			scenes.NewCappedCylinderScene,
			scenes.NewGroupGridScene,
			scenes.NewHollowGlassSphereScene,
			scenes.NewRoomScene,
		),
		canvas: view.NewCanvas(width, width),
		loc: &scenes.CameraLocation{
			At:        geom.NewPoint(2, 2, 2),
			LookingAt: geom.ZeroPoint(),
		},
	}

	var rendered uint32 = 0
	var pixelsPerRenderStat uint32 = 3000

	// render middle
	go func() {
		// todo race when scene changed.

		for {
			// todo: render option for # of rays picking random pixels & size to multiply color by
			// (list of pixels hit for each frame and last N frames, average by distance then time)

			var ctx context.Context
			ctx, s.cancel = context.WithCancel(context.Background())

			s.loc = &s.scenes[s.currentScene].Cs[s.currentCamera]
			pc, err := view.Render(ctx, s.scenes[s.currentScene].W, view.NewCameraAt(width, width, fov, s.loc.At, s.loc.LookingAt), 3, int(float64(runtime.NumCPU())/4), coordinate_supplier.Random)
			if err != nil {
				fmt.Println("failed create render")
				log.Fatalf(err.Error())
			}
			tLastRenderStat := time.Now()
			for p := range pc {
				if n := atomic.AddUint32(&rendered, 1); n%pixelsPerRenderStat == 0 {
					fmt.Printf("Writing %f pixels/sec\n", float64(pixelsPerRenderStat)/time.Since(tLastRenderStat).Seconds())
					tLastRenderStat = time.Now()
				}
				s.canvas.SetPixel(p.X, p.Y, p.C)
			}

			// todo wait if fully drawn
		}
	}()

	return s
}

func (s *state) Update() error {
	s.frameCount++

	if inpututil.IsKeyJustPressed(ebiten.KeyW) {
		s.loc.At = s.loc.At.Sub(s.loc.At.Sub(s.loc.LookingAt).Normalize())
		s.cancel()
		s.canvas = view.NewCanvas(width, width)

	}
	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		s.loc.At = s.loc.At.Add(s.loc.At.Sub(s.loc.LookingAt).Normalize())
		s.cancel()
		s.canvas = view.NewCanvas(width, width)
	}

	// todo cross product for 90 degree motiuon
	// todo A D Q E

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		s.loc.RotateAroundY(math.Pi / 12)
		s.cancel()
		s.canvas = view.NewCanvas(width, width)

	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		s.loc.RotateAroundY(-math.Pi / 12)
		s.cancel()
		s.canvas = view.NewCanvas(width, width)
	}

	// move toward origin
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		s.loc.At = s.loc.At.Mul(0.9)
		s.cancel()
		s.canvas = view.NewCanvas(width, width)
	}
	// move away from origin
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		s.loc.At = s.loc.At.Mul(1.1)
		s.cancel()
		s.canvas = view.NewCanvas(width, width)
	}
	// translate up
	if inpututil.IsKeyJustPressed(ebiten.KeyC) {
		s.loc.At.Y++
		s.cancel()
		s.canvas = view.NewCanvas(width, width)
	}
	// translate down
	if inpututil.IsKeyJustPressed(ebiten.KeyZ) {
		s.loc.At.Y--
		s.cancel()
		s.canvas = view.NewCanvas(width, width)
	}
	// next scene
	if inpututil.IsKeyJustPressed(ebiten.KeyN) {
		s.currentScene++
		if s.currentScene >= len(s.scenes) {
			s.currentScene = 0
		}
		s.currentCamera = 0
		s.cancel()
		s.canvas = view.NewCanvas(width, width)
	}
	// next camera
	if inpututil.IsKeyJustPressed(ebiten.KeyM) {
		s.currentCamera++
		if s.currentCamera >= len(s.scenes[s.currentScene].Cs) {
			s.currentCamera = 0
		}
		s.cancel()
		s.canvas = view.NewCanvas(width, width)
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
