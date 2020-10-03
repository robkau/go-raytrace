package main

import (
	"io/ioutil"
)

type projectile struct {
	x tuple //point
	v tuple //vector // todo: refactor and take advantage of interface/types to enforce this
}

type environment struct {
	gravity tuple
	wind    tuple
}

func tick(e environment, p projectile) projectile {
	return projectile{
		x: p.x.add(p.v),
		v: p.v.add(e.gravity).add(e.wind),
	}
}

func main() {

	p := projectile{
		newPoint(0, 1, 0),
		newVector(4, 7, 0),
	}

	e := environment{
		newVector(0, -0.1, 0),
		newVector(-0.01, 0, 0),
	}

	canvasHeight := 500
	c := newCanvas(500, 500)

	for {
		if p.x.y <= 0 {
			break
		}

		c.setPixel(int(p.x.x), canvasHeight-int(p.x.y)-1, color{128, 128, 128})

		p = tick(e, p)
	}

	ioutil.WriteFile("output.ppm", []byte(c.toPPM()), 0755)
}
