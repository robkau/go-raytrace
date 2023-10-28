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
	"os"
	"path"
	"time"

	"github.com/grafana/pyroscope-go"
)

const (
	width = 1200
	fov   = 0.45
)

func main() {
	p, err := pyroscope.Start(pyroscope.Config{
		ApplicationName: "raytracer",

		// replace this with the address of pyroscope server
		ServerAddress: "http://localhost:4040",

		// you can disable logging by setting this to nil
		Logger: pyroscope.StandardLogger,

		// Optional HTTP Basic authentication (Grafana Cloud)
		///BasicAuthUser:     "<User>",
		//BasicAuthPassword: "<Password>",
		// Optional Pyroscope tenant ID (only needed if using multi-tenancy). Not needed for Grafana Cloud.
		// TenantID:          "<TenantID>",

		// by default all profilers are enabled,
		// but you can select the ones you want to use:
		//ProfileTypes: []pyroscope.ProfileType{
		//	pyroscope.ProfileCPU,
		//	pyroscope.ProfileAllocObjects,
		//	pyroscope.ProfileAllocSpace,
		//	pyroscope.ProfileInuseObjects,
		//	pyroscope.ProfileInuseSpace,
		//},
	})
	if err != nil {
		log.Fatal(err)
	}

	defer p.Stop()

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
