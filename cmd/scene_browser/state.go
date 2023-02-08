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
	"github.com/robkau/go-raytrace/lib/view/canvas"
	"log"
	"runtime"
	"sync/atomic"
	"time"
)

// state struct implements ebiten.Game interface
type state struct {
	frameCount       int
	currentScene     int
	currentCamera    int
	rayBounces       int32
	renderGoroutines int32
	loc              *scenes.CameraLocation
	scenes           []*scenes.Scene
	canvas           *canvas.Canvas

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
		canvas: canvas.NewCanvas(width, width),
		loc: &scenes.CameraLocation{
			At:        geom.NewPoint(2, 2, 2),
			LookingAt: geom.ZeroPoint(),
		},
		rayBounces:       3,
		renderGoroutines: int32(runtime.NumCPU() / 3),
	}

	var rendered uint32 = 0
	var pixelsPerRenderStat uint32 = 15000

	// render
	go func() {
		// todo race when scene changed.

		for {
			var ctx context.Context
			ctx, s.cancel = context.WithCancel(context.Background())

			// todo show loading progress
			s.scenes[s.currentScene].Load()

			s.loc = &s.scenes[s.currentScene].Cs[s.currentCamera]
			bounces := atomic.LoadInt32(&s.rayBounces)
			if bounces < 0 {
				bounces = 0
				atomic.StoreInt32(&s.rayBounces, 0)
			}
			renderGoroutines := atomic.LoadInt32(&s.renderGoroutines)
			if renderGoroutines < 0 {
				renderGoroutines = 0
				atomic.StoreInt32(&s.renderGoroutines, 0)

			}
			log.Println("camera at", s.loc.At, "pointed to", s.loc.LookingAt)
			pc, err := view.Render(ctx, s.scenes[s.currentScene].W, view.NewCameraAt(width, width, fov, s.loc.At, s.loc.LookingAt), int(bounces), int(renderGoroutines), coordinate_supplier.Random)
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

			select {
			case <-ctx.Done():
				// wait until this render is done. (don't repeat drawing something already finished)
			}
		}
	}()

	return s
}

func (s *state) Update() error {
	s.frameCount++

	//if inpututil.IsKeyJustPressed(ebiten.KeyW) {
	//	s.loc.At = s.loc.At.Sub(s.loc.At.Sub(s.loc.LookingAt).Normalize())
	//	s.cancel()
	//	s.canvas = canvas.NewCanvas(width, width)
	//}
	//if inpututil.IsKeyJustPressed(ebiten.KeyS) {
	//	s.loc.At = s.loc.At.Add(s.loc.At.Sub(s.loc.LookingAt).Normalize())
	//	s.cancel()
	//	s.canvas = canvas.NewCanvas(width, width)
	//}

	// increase/decrease ray bounces
	if inpututil.IsKeyJustPressed(ebiten.KeyNumpadAdd) {
		b := atomic.AddInt32(&s.rayBounces, 1)
		log.Println(b, "ray bounces")
		s.cancel()
		s.canvas = canvas.NewCanvas(width, width)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyNumpadSubtract) {
		b := atomic.AddInt32(&s.rayBounces, -1)
		log.Println(b, "ray bounces")
		s.cancel()
		s.canvas = canvas.NewCanvas(width, width)
	}

	// increase/decrease render goroutines
	if inpututil.IsKeyJustPressed(ebiten.KeyNumpadMultiply) {
		b := atomic.AddInt32(&s.renderGoroutines, 1)
		log.Println(b, "render goroutines")
		s.cancel()
		s.canvas = canvas.NewCanvas(width, width)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyNumpadDivide) {
		b := atomic.AddInt32(&s.renderGoroutines, -1)
		log.Println(b, "render goroutines")
		s.cancel()
		s.canvas = canvas.NewCanvas(width, width)
	}

	// move toward origin
	if inpututil.IsKeyJustPressed(ebiten.KeyW) {
		s.loc.At = s.loc.At.Mul(0.9)
		s.cancel()
		s.canvas = canvas.NewCanvas(width, width)
	}
	// move away from origin
	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		s.loc.At = s.loc.At.Mul(1.1)
		s.cancel()
		s.canvas = canvas.NewCanvas(width, width)
	}
	// translate left
	if inpututil.IsKeyJustPressed(ebiten.KeyA) {
		s.loc.At.X++
		s.cancel()
		s.canvas = canvas.NewCanvas(width, width)
	}
	// translate right
	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		s.loc.At.X--
		s.cancel()
		s.canvas = canvas.NewCanvas(width, width)
	}
	// translate z left
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		s.loc.At.Z++
		s.cancel()
		s.canvas = canvas.NewCanvas(width, width)
	}
	// translate z right
	if inpututil.IsKeyJustPressed(ebiten.KeyE) {
		s.loc.At.Z--
		s.cancel()
		s.canvas = canvas.NewCanvas(width, width)
	}
	// translate up
	if inpututil.IsKeyJustPressed(ebiten.KeyC) {
		s.loc.At.Y++
		s.cancel()
		s.canvas = canvas.NewCanvas(width, width)
	}
	// translate down
	if inpututil.IsKeyJustPressed(ebiten.KeyZ) {
		s.loc.At.Y--
		s.cancel()
		s.canvas = canvas.NewCanvas(width, width)
	}

	// look left
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		s.loc.LookingAt.X += 0.25
		s.cancel()
		s.canvas = canvas.NewCanvas(width, width)
	}
	// look right
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		s.loc.LookingAt.X -= 0.25
		s.cancel()
		s.canvas = canvas.NewCanvas(width, width)
	}
	// look up
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		s.loc.LookingAt.Y += 0.25
		s.cancel()
		s.canvas = canvas.NewCanvas(width, width)
	}
	// look down
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		s.loc.LookingAt.Y -= 0.25
		s.cancel()
		s.canvas = canvas.NewCanvas(width, width)
	}
	// next scene
	if inpututil.IsKeyJustPressed(ebiten.KeyN) {
		s.currentScene++
		if s.currentScene >= len(s.scenes) {
			s.currentScene = 0
		}
		s.currentCamera = 0
		s.cancel()
		s.canvas = canvas.NewCanvas(width, width)
	}
	// next camera
	if inpututil.IsKeyJustPressed(ebiten.KeyM) {
		s.currentCamera++
		if s.currentCamera >= len(s.scenes[s.currentScene].Cs) {
			s.currentCamera = 0
		}
		s.cancel()
		s.canvas = canvas.NewCanvas(width, width)
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
