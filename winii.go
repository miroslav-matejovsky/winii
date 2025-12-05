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
//		EnableSystem: true,
//		EnableRuntimes: true,
//		EnableNetwork: true,
//		EnableWindowsServices: true,
//	}
//	insights, err := ProvideInventoryInsight(options)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	if insights.System != nil {
//	    fmt.Printf("CPU cores: %d\n", insights.System.CPUMax)
//	}
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
// It allows customizing which components to include and how to filter Windows services.
type InsightsOptions struct {
	// EnableSystem controls whether to gather system insights.
	EnableSystem bool
	// EnableRuntimes controls whether to gather runtime insights.
	EnableRuntimes bool
	// EnableNetwork controls whether to gather network insights.
	EnableNetwork bool
	// EnableWindowsServices controls whether to gather Windows services insights.
	EnableWindowsServices bool
	// WinServiceOptions specifies filtering options for Windows services.
	WinServiceOptions winservices.InventoryInsightOption
}

// Insights contains the gathered inventory data for a Windows system.
// It includes system details, runtimes, network info, and a list of Windows services.
// Fields may be nil if the corresponding component was not enabled in options.
type Insights struct {
	System          *system.Insights
	Runtimes        *runtimes.Insights
	Network         *network.Insights
	WindowsServices []winservices.ServiceDetails
}

// ProvideInsights gathers inventory insights based on the provided options.
// It conditionally collects system, runtime, network, and Windows service information
// based on the Enable* flags. Disabled components will have nil or zero values in the result.
func ProvideInsights(options InsightsOptions) (*Insights, error) {
	insights := &Insights{}

	if options.EnableSystem {
		systemInsights, err := system.GetSystemInsights()
		if err != nil {
			return nil, fmt.Errorf("failed to get system insights: %w", err)
		}
		insights.System = systemInsights
	}

	if options.EnableRuntimes {
		runtimeInsights, err := runtimes.ProvideInsights()
		if err != nil {
			return nil, fmt.Errorf("failed to provide runtimes insights: %w", err)
		}
		insights.Runtimes = runtimeInsights
	}

	if options.EnableWindowsServices {
		servicesInsights, err := winservices.ProvideInsights(options.WinServiceOptions)
		if err != nil {
			return nil, fmt.Errorf("failed to provide windows services insights: %w", err)
		}
		insights.WindowsServices = servicesInsights
	}

	if options.EnableNetwork {
		networkInsights, err := network.ProvideNetworkInsights()
		if err != nil {
			return nil, fmt.Errorf("failed to provide network insights: %w", err)
		}
		insights.Network = networkInsights
	}

	return insights, nil
}
