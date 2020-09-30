package main

import "fmt"

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
		newVector(1, 1, 0).normalize(),
	}

	e := environment{
		newVector(0, -0.1, 0),
		newVector(-0.01, 0, 0),
	}

	for {
		fmt.Printf(fmt.Sprintf("Current position: %v\n", p.x))
		if p.x.y <= 0 {
			break
		}
		p = tick(e, p)
	}
}
