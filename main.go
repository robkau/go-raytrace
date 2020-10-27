package main

import "io/ioutil"

func main() {
	wallSize := 7.0
	half := wallSize / 2
	wallZ := 10.0
	canvasPixels := 100
	pixelSize := wallSize / float64(canvasPixels)

	origin := newPoint(0, 0, -5)
	canvas := newCanvas(canvasPixels, canvasPixels)
	s := newSphereWith(scale(1.2, 1, 1))

	for i := 0; i < canvasPixels; i++ {
		worldY := half - pixelSize*float64(i)

		for j := 0; j < canvasPixels; j++ {
			worldX := -half + pixelSize*float64(j)

			position := newPoint(worldX, worldY, wallZ)
			r := rayWith(origin, position.sub(origin).normalize())

			xs := s.intersect(r)
			_, ok := xs.hit()
			if ok {
				canvas.setPixel(j, i, color{255, 0, 0})
			}
		}

	}

	ioutil.WriteFile("output.ppm", []byte(canvas.toPPM()), 0755)
}
