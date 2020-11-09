package main

import (
	"io/ioutil"
)

func main() {

	wallSize := 7.0
	half := wallSize / 2
	wallZ := 10.0
	canvasPixels := 750
	pixelSize := wallSize / float64(canvasPixels)

	origin := newPoint(0, 0, -5)
	canvas := newCanvas(canvasPixels, canvasPixels)
	s := newSphere()
	s.m.color = color{1, 0.2, 1}

	lightPosition := newPoint(-10, 10, -10)
	lightColor := color{1.0, 1.0, 1.0}
	light := newPointLight(lightPosition, lightColor)

	for i := 0; i < canvasPixels; i++ {
		worldY := half - pixelSize*float64(i)

		for j := 0; j < canvasPixels; j++ {
			worldX := -half + pixelSize*float64(j)
			position := newPoint(worldX, worldY, wallZ)
			r := rayWith(origin, position.sub(origin).normalize())

			xs := s.intersect(r)
			t, ok := xs.hit()
			if ok {
				pos := r.position(t.t)
				normal := t.o.normalAt(pos)
				eye := r.direction.neg()

				c := lighting(t.o.m, light, pos, eye, normal)

				canvas.setPixel(j, i, c)
			}
		}

	}

	ioutil.WriteFile("output.ppm", []byte(canvas.toPPM()), 0755)
}
