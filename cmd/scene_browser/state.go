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
	"image/gif"
	"log"
	"math"
	"os"
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

	recordingGif uint32
	gifFrame     int
	gifData      gif.GIF

	lastImage *ebiten.Image

	cancel context.CancelFunc
}

func start() *state {
	rp := view.NewRayPool()

	s := &state{
		scenes: scenes.LoadScenes(
			scenes.NewGroupTransformsScene,
			scenes.NewBeadTapestry,
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
		rayBounces:       2,
		renderGoroutines: int32(runtime.NumCPU() / 2),
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
			pc, err := view.Render(ctx, s.scenes[s.currentScene].W, view.NewCameraAt(width, width, fov, s.loc.At, s.loc.LookingAt), int(bounces), int(renderGoroutines), coordinate_supplier.Random, rp)
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

			if s.recordingGif == 1 {
				s.gifFrame++
				// save gif frame from canvas
				s.gifData.Image = append(s.gifData.Image, s.canvas.ToImagePaletted())
				s.gifData.Delay = append(s.gifData.Delay, 0)

				// go next frame
				s.canvas = canvas.NewCanvas(width, width)
				s.scenes[s.currentScene].Cs[s.currentCamera].RotateAroundY(1 * math.Pi / 180)

				if s.gifFrame == 360 { // todo off by one?
					// all done!
					// todo handle errs
					f, _ := os.OpenFile("rgb.gif", os.O_WRONLY|os.O_CREATE, 0600)
					gif.EncodeAll(f, &s.gifData)
					f.Close()
					atomic.StoreUint32(&s.recordingGif, 0)
				}

			} else {
				select {
				case <-ctx.Done():
					// wait after this render is done. (don't repeat drawing something already finished)
				}
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

	if atomic.LoadUint32(&s.recordingGif) == 1 {
		// disable input during gif save.
		return nil
	}

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
		if ebiten.IsKeyPressed(ebiten.KeyShift) {
			// big move
			s.loc.At.X += 25
		} else {
			// normal move
			s.loc.At.X++
		}
		s.cancel()
		s.canvas = canvas.NewCanvas(width, width)
	}
	// translate right
	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		if ebiten.IsKeyPressed(ebiten.KeyShift) {
			// big move
			s.loc.At.X -= 25
		} else {
			// normal move
			s.loc.At.X--
		}
		s.cancel()
		s.canvas = canvas.NewCanvas(width, width)
	}
	// translate z left
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		if ebiten.IsKeyPressed(ebiten.KeyShift) {
			// big move
			s.loc.At.Z += 25
		} else {
			// normal move
			s.loc.At.Z++
		}
		s.cancel()
		s.canvas = canvas.NewCanvas(width, width)
	}
	// translate z right
	if inpututil.IsKeyJustPressed(ebiten.KeyE) {
		if ebiten.IsKeyPressed(ebiten.KeyShift) {
			// big move
			s.loc.At.Z -= 25
		} else {
			// normal move
			s.loc.At.Z--
		}
		s.cancel()
		s.canvas = canvas.NewCanvas(width, width)
	}
	// translate up
	if inpututil.IsKeyJustPressed(ebiten.KeyC) {
		if ebiten.IsKeyPressed(ebiten.KeyShift) {
			// big move
			s.loc.At.Y += 25
		} else {
			// normal move
			s.loc.At.Y++
		}
		s.cancel()
		s.canvas = canvas.NewCanvas(width, width)
	}
	// translate down
	if inpututil.IsKeyJustPressed(ebiten.KeyZ) {
		if ebiten.IsKeyPressed(ebiten.KeyShift) {
			// big move
			s.loc.At.Y -= 25
		} else {
			// normal move
			s.loc.At.X--
		}
		s.cancel()
		s.canvas = canvas.NewCanvas(width, width)
	}

	// look left
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		if ebiten.IsKeyPressed(ebiten.KeyShift) {
			// big move
			s.loc.LookingAt.X += 25
		} else {
			// normal move
			s.loc.LookingAt.X += 0.25
		}

		s.cancel()
		s.canvas = canvas.NewCanvas(width, width)
	}
	// look right
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		if ebiten.IsKeyPressed(ebiten.KeyShift) {
			// big move
			s.loc.LookingAt.X -= 25
		} else {
			// normal move
			s.loc.LookingAt.X -= 0.25
		}
		s.cancel()
		s.canvas.Reset()
	}
	// look up
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		if ebiten.IsKeyPressed(ebiten.KeyShift) {
			// big move
			s.loc.LookingAt.Y += 25
		} else {
			// normal move
			s.loc.LookingAt.Y += 0.25
		}
		s.cancel()
		s.canvas = canvas.NewCanvas(width, width)
	}
	// look down
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		if ebiten.IsKeyPressed(ebiten.KeyShift) {
			// big move
			s.loc.LookingAt.Y -= 25
		} else {
			// normal move
			s.loc.LookingAt.Y -= 0.25
		}
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

	// spin around and record a gif!
	if inpututil.IsKeyJustPressed(ebiten.KeyG) {
		atomic.StoreUint32(&s.recordingGif, 1)
		s.cancel()
		s.canvas = canvas.NewCanvas(width, width)
	}

	// canvas is updated in background goroutine.
	return nil
}

func (s *state) Draw(screen *ebiten.Image) {
	// render current frame progress
	op := &ebiten.DrawImageOptions{}

	s.lastImage = s.canvas.ToEbitenImage(s.lastImage)
	screen.DrawImage(ebiten.NewImageFromImage(s.lastImage), op)
}

func (s *state) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}
