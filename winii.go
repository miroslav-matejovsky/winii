// Package winii provides a Go library for gathering comprehensive inventory insights
// about Windows systems. It collects information on system details, installed runtimes,
// network configurations, and Windows services.
package winii

import (
	"github.com/miroslav-matejovsky/winii/winii"
	"github.com/miroslav-matejovsky/winii/winservices"
)

// ProvideAllInventoryInsight gathers and returns a full set of inventory insights
// for the current Windows system using default options for all components.
// It returns a pointer to an Insights struct containing detailed system information,
// or an error if the gathering process fails.
func ProvideAllInventoryInsight() (*winii.Insights, error) {
	options := winii.InsightsOptions{
		WinServiceOptions: winservices.InventoryInsightOption{},
	}
	return winii.ProvideInsights(options)
}

// ProvideMinimalInventoryInsight allows gathering inventory insights with custom options.
// The options parameter specifies which insights to include and how to filter them.
// It returns a pointer to an Insights struct or an error if gathering fails.
func ProvideMinimalInventoryInsight(options winii.InsightsOptions) (*winii.Insights, error) {
	return winii.ProvideInsights(options)
}
