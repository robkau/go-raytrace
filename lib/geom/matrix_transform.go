package geom

import "math"

func Translate(x, y, z float64) *X4Matrix {
	m := NewIdentityMatrixX4()
	m = m.Set(0, 3, x)
	m = m.Set(1, 3, y)
	m = m.Set(2, 3, z)
	return m
}

func Scale(x, y, z float64) *X4Matrix {
	m := NewIdentityMatrixX4()
	m = m.Set(0, 0, x)
	m = m.Set(1, 1, y)
	m = m.Set(2, 2, z)
	return m
}

func RotateX(rad float64) *X4Matrix {
	m := NewIdentityMatrixX4()
	m = m.Set(1, 1, math.Cos(rad))
	m = m.Set(1, 2, -math.Sin(rad))
	m = m.Set(2, 1, math.Sin(rad))
	m = m.Set(2, 2, math.Cos(rad))
	return m
}

func RotateY(rad float64) *X4Matrix {
	m := NewIdentityMatrixX4()
	m = m.Set(0, 0, math.Cos(rad))
	m = m.Set(0, 2, math.Sin(rad))
	m = m.Set(2, 0, -math.Sin(rad))
	m = m.Set(2, 2, math.Cos(rad))
	return m
}

func RotateZ(rad float64) *X4Matrix {
	m := NewIdentityMatrixX4()
	m = m.Set(0, 0, math.Cos(rad))
	m = m.Set(0, 1, -math.Sin(rad))
	m = m.Set(1, 0, math.Sin(rad))
	m = m.Set(1, 1, math.Cos(rad))
	return m
}

func shear(xy, xz, yx, yz, zx, zy float64) *X4Matrix {
	m := NewIdentityMatrixX4()
	m = m.Set(0, 1, xy)
	m = m.Set(0, 2, xz)
	m = m.Set(1, 0, yx)
	m = m.Set(1, 2, yz)
	m = m.Set(2, 0, zx)
	m = m.Set(2, 1, zy)
	return m
}
