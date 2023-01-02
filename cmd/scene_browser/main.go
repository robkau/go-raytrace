package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"log"
	"path"

	//_ "net/http/pprof"
	"os"
)

const (
	width = 700
	fov   = 0.45
)

func main() {
	//defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()

	//go func() {
	//	log.Println(http.ListenAndServe("localhost:6060", nil))
	//}()

	sb := start()

	f, err := os.Open(path.Join("data", "intro.wav"))
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
