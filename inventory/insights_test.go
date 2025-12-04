package inventory

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInventoryInsights(t *testing.T) {
	// TODO
	t.Skip("not yet")
	insight, err := ProvideInsights()
	require.NoError(t, err)
	require.NotEmpty(t, insight)
}
