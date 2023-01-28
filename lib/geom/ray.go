package geom

type Ray struct {
	Origin    Tuple
	Direction Tuple
}

func RayWith(origin, direction Tuple) Ray {
	return Ray{origin, direction}
}

func (r Ray) Position(t float64) Tuple {
	return r.Origin.Add(r.Direction.Mul(t))
}

func (r Ray) Transform(m *X4Matrix) Ray {
	return RayWith(
		m.MulTuple(r.Origin),
		m.MulTuple(r.Direction),
	)
}
