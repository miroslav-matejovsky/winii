package winservices

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProvideInsights(t *testing.T) {

	svcManager := NewWinSvcManager()
	defer func() { _ = svcManager.Disconnect() }()

	allServices, err := svcManager.ListServices()
	require.NoError(t, err)
	require.NotEmpty(t, allServices)

	insight, err := ProvideInsights(
		InventoryInsightOption{},
	)
	require.NoError(t, err)
	require.Len(t, insight, len(allServices))
}

func TestProvideInsightsWithFilter(t *testing.T) {
	insight, err := ProvideInsights(
		InventoryInsightOption{
			WinServiceNameRegex: regexp.MustCompile(`(?i)^wuauserv$`),
		},
	)
	require.NoError(t, err)
	require.Len(t, insight, 1)
	require.Equal(t, "wuauserv", insight[0].Name)
}
