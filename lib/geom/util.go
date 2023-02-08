package geom

func TransformPoint(t *X4Matrix, p Tuple) Tuple {
	return t.Invert().MulTuple(p)
}

func IdentityTransform(p Tuple) Tuple {
	return p
}

func DoubleScale(p Tuple) Tuple {
	return TransformPoint(Scale(2, 2, 2), p)
}
