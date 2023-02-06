package main

import (
	"bytes"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"io"
	"log"
	"math/rand"
	"path"
	"time"

	//_ "net/http/pprof"
	"os"
)

const (
	width = 1080
	fov   = 0.45
)

func main() {
	//defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()

	//go func() {
	//	log.Println(http.ListenAndServe("localhost:6060", nil))
	//}()

	sb := start()

	// play a random intro file, sometimes.
	rand.Seed(time.Now().UnixNano())
	audioIndex := rand.Intn(1)
	f, err := os.Open(path.Join("data", "audio", fmt.Sprintf("intro_%d.wav", audioIndex)))
	if err == nil {
		var fBytes = &bytes.Buffer{}
		if _, err := io.Copy(fBytes, f); err == nil {
			ac := audio.NewContext(44100)
			d, err := wav.Decode(ac, fBytes)
			if err == nil {
				ap, err := audio.NewPlayer(ac, d)
				if err == nil {
					ap.Play()
					defer ap.Close()
				}

			}
		}
		f.Close()
	}

	ebiten.SetWindowSize(width, width)
	ebiten.SetWindowTitle("go-raytrace")
	if err := ebiten.RunGame(sb); err != nil {
		log.Fatal(err)
	}
}
