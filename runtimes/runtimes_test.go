package runtimes

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProvideInsights(t *testing.T) {
	insights, err := ProvideInsights()
	require.NoError(t, err)
	require.NotNil(t, insights)
	require.NotNil(t, insights.VCRedistRuntimes)
	require.NotNil(t, insights.DotNetRuntimes)
}
