package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_NewLight(t *testing.T) {
	intensity := color{1, 1, 1}
	position := newPoint(0, 0, 0)

	l := newPointLight(position, intensity)

	assert.Equal(t, intensity, l.intensity)
	assert.Equal(t, position, l.position)
}
