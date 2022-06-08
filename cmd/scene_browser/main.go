package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"log"
	//_ "net/http/pprof"
	"os"
)

const (
	width = 333
	fov   = 0.45
)

func main() {
	//runtime.SetBlockProfileRate(100_000_000) // WARNING: Can cause some CPU overhead
	//file, _ := os.Create("./block.pprof")
	//defer pprof.Lookup("block").WriteTo(file, 0)

	//go func() {
	//	log.Println(http.ListenAndServe("localhost:6060", nil))
	//}()

	sb := start()

	f, err := os.Open("intro.wav")
	if err == nil {
		defer f.Close()
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

	ebiten.SetWindowSize(width, width)
	ebiten.SetWindowTitle("go-raytrace")
	if err := ebiten.RunGame(sb); err != nil {
		log.Fatal(err)
	}
}
