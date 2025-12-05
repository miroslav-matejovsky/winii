package winii

import (
	"github.com/miroslav-matejovsky/winii/runtimes"
	"github.com/miroslav-matejovsky/winii/winservices"
)

type InsightsOptions struct {
	WinServiceOptions winservices.InventoryInsightOption
}

type Result struct {
	Runtimes        runtimes.Insights
	WindowsServices []winservices.ServiceDetails
}

func ProvideInsights(options InsightsOptions) (*Result, error) {
	runtimeInsights, err := runtimes.ProvideInsights()
	if err != nil {
		return nil, err
	}
	services, err := winservices.ProvideInsights(options.WinServiceOptions)
	if err != nil {
		return nil, err
	}
	return &Result{
		Runtimes:        *runtimeInsights,
		WindowsServices: services,
	}, nil
}
