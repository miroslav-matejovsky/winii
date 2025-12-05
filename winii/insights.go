package winii

import (
	"fmt"

	"github.com/miroslav-matejovsky/winii/network"
	"github.com/miroslav-matejovsky/winii/runtimes"
	"github.com/miroslav-matejovsky/winii/system"
	"github.com/miroslav-matejovsky/winii/winservices"
)

type InsightsOptions struct {
	WinServiceOptions winservices.InventoryInsightOption
}

type Insights struct {
	System          system.Insights
	Runtimes        runtimes.Insights
	Network         network.Insights
	WindowsServices []winservices.ServiceDetails
}

func ProvideInsights(options InsightsOptions) (*Insights, error) {
	systemInsights, err := system.GetSystemInsights()
	if err != nil {
		return nil, fmt.Errorf("failed to get system insights: %w", err)
	}
	runtimeInsights, err := runtimes.ProvideInsights()
	if err != nil {
		return nil, fmt.Errorf("failed to provide runtimes insights: %w", err)
	}
	servicesInsights, err := winservices.ProvideInsights(options.WinServiceOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to provide windows services insights: %w", err)
	}
	networkInsights, err := network.ProvideNetworkInsights()
	if err != nil {
		return nil, fmt.Errorf("failed to provide network insights: %w", err)
	}
	return &Insights{
		System:          *systemInsights,
		Runtimes:        *runtimeInsights,
		Network:         *networkInsights,
		WindowsServices: servicesInsights,
	}, nil
}
