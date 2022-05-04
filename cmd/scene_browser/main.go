package main

import (
	"flag"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/pkg/profile"
	"log"
	"os"
)

const (
	width = 444
	fov   = 0.45
)

var cpuProfile = flag.String("cpuprofile", "", "write cpu profiling information")
var heapProfile = flag.String("heapprofile", "", "write heap memory profiling information")
var allocProfile = flag.String("allocprofile", "", "write memory alloc profiling information")
var traceProfile = flag.String("traceprofile", "", "write trace profiling information")

func main() {
	flag.Parse()
	if *cpuProfile != "" {
		defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()
	} else if *heapProfile != "" {
		defer profile.Start(profile.MemProfileHeap, profile.ProfilePath(".")).Stop()
	} else if *allocProfile != "" {
		defer profile.Start(profile.MemProfileAllocs, profile.ProfilePath(".")).Stop()
	} else if *traceProfile != "" {
		defer profile.Start(profile.TraceProfile, profile.ProfilePath(".")).Stop()
	}

	sb := start()

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
	if err := ebiten.RunGame(sb); err != nil {
		log.Fatal(err)
	}
}
