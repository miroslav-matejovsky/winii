package winii

import (
	"testing"

	"github.com/miroslav-matejovsky/winii/winservices"
	"github.com/stretchr/testify/require"
)

func TestInventoryInsights(t *testing.T) {
	options := InsightsOptions{
		WinServiceOptions: winservices.InventoryInsightOption{},
	}
	insight, err := ProvideInsights(options)
	require.NoError(t, err)
	require.NotEmpty(t, insight)
}
