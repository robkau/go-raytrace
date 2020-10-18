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

func Test_Compare_Smaller(t *testing.T) {
	a := 0.00000000001
	b := 0.00000000002
	c := 0.00000000003

	assert.True(t, almostEqual(a+b, c))
}

func Test_Compare_Smallest_Match(t *testing.T) {
	a := 0.00000000000000001
	b := 0.00000000000000002
	c := 0.00000000000000003

	assert.True(t, almostEqual(a+b, c))
}

func Test_Compare_Smallest_NotMatch(t *testing.T) {
	a := 0.000000001
	b := 0.000000002
	c := 0.000000005

	assert.False(t, almostEqual(a+b, c))
}

func Test_Compare_Bigger_Match(t *testing.T) {
	a := 100000000000000.00000000001
	b := 100000000000000.00000000002
	c := 200000000000000.00000000003

	assert.True(t, almostEqual(a+b, c))
}

func Test_Compare_Bigger_NotMatch(t *testing.T) {
	a := 100000000000000.00000000001
	b := 100000000000000.00000000002
	c := 200000000000000.5

	assert.False(t, almostEqual(a+b, c))
}
