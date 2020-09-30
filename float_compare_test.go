package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Compare(t *testing.T) {
	a := 0.1
	b := 0.2
	c := 0.3

	assert.True(t, almostEqual(a+b, c))
}
