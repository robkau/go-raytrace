package geom

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_NewRay(t *testing.T) {
	r := RayWith(NewPoint(1, 2, 3), NewVector(4, 5, 6))

	assert.True(t, NewPoint(1, 2, 3).Equals(r.Origin))
	assert.True(t, NewVector(4, 5, 6).Equals(r.Direction))
}

func Test_Ray_Position(t *testing.T) {
	r := RayWith(NewPoint(2, 3, 4), NewVector(1, 0, 0))

	assert.True(t, NewPoint(2, 3, 4).Equals(r.Position(0)))
	assert.True(t, NewPoint(3, 3, 4).Equals(r.Position(1)))
	assert.True(t, NewPoint(1, 3, 4).Equals(r.Position(-1)))
	assert.True(t, NewPoint(4.5, 3, 4).Equals(r.Position(2.5)))
}

func Test_Ray_Translate(t *testing.T) {
	r := RayWith(NewPoint(1, 2, 3), NewVector(0, 1, 0))
	m := Translate(3, 4, 5)

	r2 := r.Transform(m)

	assert.Equal(t, NewPoint(4, 6, 8), r2.Origin)
	assert.Equal(t, NewVector(0, 1, 0), r2.Direction)
}

func Test_Ray_Scale(t *testing.T) {
	r := RayWith(NewPoint(1, 2, 3), NewVector(0, 1, 0))
	m := Scale(2, 3, 4)

	r2 := r.Transform(m)

	assert.Equal(t, NewPoint(2, 6, 12), r2.Origin)
	assert.Equal(t, NewVector(0, 3, 0), r2.Direction)
}
