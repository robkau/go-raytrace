package parse

import (
	"github.com/robkau/go-raytrace/lib/geom"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func Test_ParsePositionsFromReplay(t *testing.T) {
	pr, err := ParseReaderAsTori(strings.NewReader(ReplayFile))
	require.NoError(t, err)

	require.Len(t, pr.P0Positions, 7)
	require.Len(t, pr.P0Positions[0].parts, 21)
	require.Equal(t, geom.NewPoint(1.00000000, 0.35000002, 2.59000014), pr.P0Positions[0].parts[0])
	require.Equal(t, geom.NewPoint(0.80249145, -0.79330789, 0.75910965), pr.P0Positions[6].parts[0])
	require.Equal(t, geom.NewPoint(1.00000000, -0.54999995, 2.59000014), pr.P1Positions[0].parts[0])
	require.Equal(t, geom.NewPoint(0.84716686, -3.50764995, 0.59010877), pr.P1Positions[7].parts[0])
}

func Test_ParsePositionLine(t *testing.T) {
	positionLine := `POS 0; 1.00000000 0.35000002 2.59000014 1.00000000 0.39999998 2.14000010 1.00000000 0.39999998 1.89000010 1.00000000 0.44999999 1.69000005 1.00000000 0.50000000 1.49000000 0.75000000 0.39999998 2.09000014 0.44999999 0.39999998 2.24000000 0.05000000 0.39999998 2.24000000 1.25000000 0.39999998 2.09000014 1.54999995 0.39999998 2.24000000 1.95000005 0.39999998 2.24000000 -0.34999999 0.34999996 2.24000000 2.34999990 0.34999996 2.24000000 0.80000001 0.50000000 1.39000010 1.20000005 0.50000000 1.39000010 0.80000001 0.50000000 1.04000007 1.20000005 0.50000000 1.04000007 1.20000005 0.50000000 0.43999999 0.80000001 0.50000000 0.43999999 0.80000001 0.39999998 0.04000000 1.20000005 0.39999998 0.04000000`

	player, positions := parsePositionLine(positionLine)

	require.Equal(t, 0, player)
	require.Len(t, positions.parts, 21)
}
