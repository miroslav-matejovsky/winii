package inventory

import (
	"github.com/miroslav-matejovsky/winii/runtimes"
	"github.com/miroslav-matejovsky/winii/winservices"
)

type Result struct {
	Runtimes        runtimes.Insights
	WindowsServices []winservices.ServiceDetails
}

func ProvideInsights() (*Result, error) {
	runtimeInsights, err := runtimes.ProvideInsights()
	if err != nil {
		return nil, err
	}
	services, err := winservices.ProvideInsights()
	if err != nil {
		return nil, err
	}
	return &Result{
		Runtimes:        *runtimeInsights,
		WindowsServices: services,
	}, nil
}
