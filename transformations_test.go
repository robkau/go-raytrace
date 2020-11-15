package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_DefaultTransformationMatrix(t *testing.T) {
	from := newPoint(0, 0, 0)
	to := newPoint(0, 0, -1)
	up := newVector(0, 1, 0)

	tr := viewTransform(from, to, up)

	assert.Equal(t, newIdentityMatrixX4(), tr)
}

func Test_PositiveZTransformationMatrix(t *testing.T) {
	from := newPoint(0, 0, 0)
	to := newPoint(0, 0, 1)
	up := newVector(0, 1, 0)

	tr := viewTransform(from, to, up)

	assert.Equal(t, scale(-1, 1, -1), tr)
}

func Test_WorldMovedTransformationMatrix(t *testing.T) {
	from := newPoint(0, 0, 8)
	to := newPoint(0, 0, 0)
	up := newVector(0, 1, 0)

	tr := viewTransform(from, to, up)

	assert.Equal(t, translate(0, 0, -8), tr)
}

func Test_ArbitraryTransformationMatrix(t *testing.T) {
	from := newPoint(1, 3, 2)
	to := newPoint(4, -2, 8)
	up := newVector(1, 1, 0)

	tr := viewTransform(from, to, up).roundTo(5)

	expected := newX4MatrixWith(
		-0.50709, 0.50709, 0.67612, -2.36643,
		0.76772, 0.60609, 0.12122, -2.82843,
		-0.35857, 0.59761, -0.71714, 0.0,
		0.0, 0.0, 0.0, 1.0,
	)

	assert.Equal(t, expected, tr)
}
