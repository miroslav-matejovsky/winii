// Package winii provides Windows inventory insights.
//
// This package offers a high-level API to gather comprehensive insights about a Windows system,
// including system details, installed runtimes, network configurations, and Windows services.
//
// The main entry point is ProvideInventoryInsight() for customizable insights.
//
// Example:
//
//	options := InsightsOptions{
//		WinServiceOptions: winservices.InventoryInsightOption{},
//	}
//	insights, err := ProvideInventoryInsight(options)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("CPU cores: %d\n", insights.System.CPU.Count)
package winii

import (
	"fmt"

	"github.com/miroslav-matejovsky/winii/network"
	"github.com/miroslav-matejovsky/winii/runtimes"
	"github.com/miroslav-matejovsky/winii/system"
	"github.com/miroslav-matejovsky/winii/winservices"
)

// ProvideInventoryInsight allows gathering inventory insights with custom options.
// The options parameter specifies which insights to include and how to filter them.
// It returns a pointer to an Insights struct or an error if gathering fails.
func ProvideInventoryInsight(options InsightsOptions) (*Insights, error) {
	return ProvideInsights(options)
}

// InsightsOptions configures the gathering of inventory insights.
// It allows customizing which Windows services to include via WinServiceOptions.
type InsightsOptions struct {
	WinServiceOptions winservices.InventoryInsightOption
}

// Insights contains the gathered inventory data for a Windows system.
// It includes system details, runtimes, network info, and a list of Windows services.
type Insights struct {
	System          system.Insights
	Runtimes        runtimes.Insights
	Network         network.Insights
	WindowsServices []winservices.ServiceDetails
}

// ProvideInsights gathers inventory insights based on the provided options.
// It collects system, runtime, network, and Windows service information.
// This is the core function used by the convenience wrappers.
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
