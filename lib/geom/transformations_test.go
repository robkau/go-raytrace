package geom

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_DefaultTransformationMatrix(t *testing.T) {
	from := NewPoint(0, 0, 0)
	to := NewPoint(0, 0, -1)
	up := NewVector(0, 1, 0)

	tr := ViewTransform(from, to, up)

	assert.Equal(t, NewIdentityMatrixX4(), tr)
}

func Test_PositiveZTransformationMatrix(t *testing.T) {
	from := NewPoint(0, 0, 0)
	to := NewPoint(0, 0, 1)
	up := NewVector(0, 1, 0)

	tr := ViewTransform(from, to, up)

	assert.Equal(t, Scale(-1, 1, -1), tr)
}

func Test_WorldMovedTransformationMatrix(t *testing.T) {
	from := NewPoint(0, 0, 8)
	to := NewPoint(0, 0, 0)
	up := NewVector(0, 1, 0)

	tr := ViewTransform(from, to, up)

	assert.Equal(t, Translate(0, 0, -8), tr)
}

func Test_ArbitraryTransformationMatrix(t *testing.T) {
	from := NewPoint(1, 3, 2)
	to := NewPoint(4, -2, 8)
	up := NewVector(1, 1, 0)

	tr := ViewTransform(from, to, up).RoundTo(5)

	expected := NewX4MatrixWith(
		-0.50709, 0.50709, 0.67612, -2.36643,
		0.76772, 0.60609, 0.12122, -2.82843,
		-0.35857, 0.59761, -0.71714, 0.0,
		0.0, 0.0, 0.0, 1.0,
	)

	assert.Equal(t, expected, tr)
}
