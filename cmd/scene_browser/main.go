package main

import (
	"flag"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/pkg/profile"
	"log"
	//_ "net/http/pprof"
	"os"
)

const (
	width = 555
	fov   = 0.45
)

var cpuProfile = flag.String("cpuprofile", "", "write cpu profiling information")
var heapProfile = flag.String("heapprofile", "", "write heap memory profiling information")
var allocProfile = flag.String("allocprofile", "", "write memory alloc profiling information")
var traceProfile = flag.String("traceprofile", "", "write trace profiling information")

func main() {
	//runtime.SetBlockProfileRate(100_000_000) // WARNING: Can cause some CPU overhead
	//file, _ := os.Create("./block.pprof")
	//defer pprof.Lookup("block").WriteTo(file, 0)

	//go func() {
	//	log.Println(http.ListenAndServe("localhost:6060", nil))
	//}()

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
