package main

import "math"

func translate(x, y, z float64) x4Matrix {
	m := newIdentityMatrixX4()
	m.set(0, 3, x)
	m.set(1, 3, y)
	m.set(2, 3, z)
	return m
}

func scale(x, y, z float64) x4Matrix {
	m := newIdentityMatrixX4()
	m.set(0, 0, x)
	m.set(1, 1, y)
	m.set(2, 2, z)
	return m
}

func rotateX(rad float64) x4Matrix {
	m := newIdentityMatrixX4()
	m.set(1, 1, math.Cos(rad))
	m.set(1, 2, -math.Sin(rad))
	m.set(2, 1, math.Sin(rad))
	m.set(2, 2, math.Cos(rad))
	return m
}

func rotateY(rad float64) x4Matrix {
	m := newIdentityMatrixX4()
	m.set(0, 0, math.Cos(rad))
	m.set(0, 2, math.Sin(rad))
	m.set(2, 0, -math.Sin(rad))
	m.set(2, 2, math.Cos(rad))
	return m
}

func rotateZ(rad float64) x4Matrix {
	m := newIdentityMatrixX4()
	m.set(0, 0, math.Cos(rad))
	m.set(0, 1, -math.Sin(rad))
	m.set(1, 0, math.Sin(rad))
	m.set(1, 1, math.Cos(rad))
	return m
}

func shear(xy, xz, yx, yz, zx, zy float64) x4Matrix {
	m := newIdentityMatrixX4()
	m.set(0, 1, xy)
	m.set(0, 2, xz)
	m.set(1, 0, yx)
	m.set(1, 2, yz)
	m.set(2, 0, zx)
	m.set(2, 1, zy)
	return m
}
