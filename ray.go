package main

type ray struct {
	origin    tuple
	direction tuple
}

func rayWith(origin, direction tuple) ray {
	return ray{origin, direction}
}

func (r ray) position(t float64) tuple {
	return r.origin.add(r.direction.mul(t))
}

func (r ray) transform(m x4Matrix) ray {
	return rayWith(
		m.mulTuple(r.origin),
		m.mulTuple(r.direction),
	)
}
