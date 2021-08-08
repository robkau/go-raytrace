package geom

func ViewTransform(from, to, up Tuple) *X4Matrix {
	forward := to.Sub(from).Normalize()
	upNormal := up.Normalize()
	left := Cross(forward, upNormal)
	trueUp := Cross(left, forward)

	orientation := NewX4MatrixWith(
		left.X, left.Y, left.Z, 0,
		trueUp.X, trueUp.Y, trueUp.Z, 0,
		-forward.X, -forward.Y, -forward.Z, 0,
		0, 0, 0, 1)

	return orientation.MulX4Matrix(Translate(-from.X, -from.Y, -from.Z))
}
