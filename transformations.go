package main

func viewTransform(from, to, up tuple) x4Matrix {
	forward := to.sub(from).normalize()
	upNormal := up.normalize()
	left := cross(forward, upNormal)
	trueUp := cross(left, forward)

	orientation := newX4MatrixWith(
		left.x, left.y, left.z, 0,
		trueUp.x, trueUp.y, trueUp.z, 0,
		-forward.x, -forward.y, -forward.z, 0,
		0, 0, 0, 1)

	return orientation.mulX4Matrix(translate(-from.x, -from.y, -from.z))
}
