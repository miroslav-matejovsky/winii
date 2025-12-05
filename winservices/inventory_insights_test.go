package winservices

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProvideInsights(t *testing.T) {
	insight, err := ProvideInsights(
		InventoryInsightOption{},
	)
	require.NoError(t, err)
	require.NotEmpty(t, insight)
}
