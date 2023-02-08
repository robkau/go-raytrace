package shapes

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_Jitterer(t *testing.T) {
	j := NewJitterSequence(0.1, 0.5, 1.0)

	require.Equal(t, 0.1, j.Next())
	require.Equal(t, 0.5, j.Next())
	require.Equal(t, 1.0, j.Next())
	require.Equal(t, 0.1, j.Next())
}
