package winii

import (
	"fmt"

	"github.com/miroslav-matejovsky/winii/network"
	"github.com/miroslav-matejovsky/winii/runtimes"
	"github.com/miroslav-matejovsky/winii/winservices"
)

type InsightsOptions struct {
	WinServiceOptions winservices.InventoryInsightOption
}

type Result struct {
	Runtimes        runtimes.Insights
	Network         network.Insights
	WindowsServices []winservices.ServiceDetails
}

func ProvideInsights(options InsightsOptions) (*Result, error) {
	runtimeInsights, err := runtimes.ProvideInsights()
	if err != nil {
		return nil, fmt.Errorf("failed to provide runtimes insights: %w", err)
	}
	services, err := winservices.ProvideInsights(options.WinServiceOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to provide windows services insights: %w", err)
	}
	network, err := network.ProvideNetworkInsights()
	if err != nil {
		return nil, fmt.Errorf("failed to provide network insights: %w", err)
	}
	return &Result{
		Runtimes:        *runtimeInsights,
		Network:         *network,
		WindowsServices: services,
	}, nil
}
