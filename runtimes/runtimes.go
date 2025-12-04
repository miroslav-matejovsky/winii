// package runtimes provides insights about installed runtimes on Windows systems.
// It includes tools for checking installed Visual C++ Redistributables and .NET runtime components.
package runtimes

import "fmt"

type Insights struct {
	VCRedistRuntimes []VCRedistRuntime
	DotNetRuntimes   []DotNetRuntime
}

func ProvideInsights() (*Insights, error) {
	vcRedist, err := ProvideVCRedistInsights()
	if err != nil {
		return nil, fmt.Errorf("failed to provide Visual C++ Redistributable runtimes insights: %w", err)
	}
	dotNet, err := ProvideDotNetInsights()
	if err != nil {
		return nil, fmt.Errorf("failed to provide .NET runtimes insights: %w", err)
	}
	return &Insights{
		VCRedistRuntimes: vcRedist,
		DotNetRuntimes:   dotNet,
	}, nil
}
