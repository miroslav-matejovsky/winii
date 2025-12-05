package network

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNetworkInsights(t *testing.T) {
	insights, err := ProvideNetworkInsights()
	require.NoError(t, err)
	require.NotNil(t, insights)
	t.Log(insights)
}
