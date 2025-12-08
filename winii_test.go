package winii

import (
	"fmt"
	"log"
	"regexp"

	"github.com/miroslav-matejovsky/winii/winservices"
)

// ExampleProvideInsights demonstrates how to gather inventory insights
// with custom options, filtering Windows services by name.
func ExampleProvideInsights() {
	options := InsightsOptions{
		EnableSystem:          true,
		EnableRuntimes:        true,
		EnableNetwork:         true,
		EnableWindowsServices: true,
		WinServiceOptions: winservices.InventoryInsightOption{
			WinServiceNameRegex: regexp.MustCompile(".*service.*"),
		},
	}
	insights, err := ProvideInsights(options)
	if err != nil {
		log.Fatal(err)
	}
	if insights.System != nil {
		fmt.Printf("System CPU cores: %d\n", insights.System.CPUMax)
	}
	if insights.Runtimes != nil {
		fmt.Printf("Runtimes: %+v\n", insights.Runtimes)
	}
	if insights.Network != nil {
		fmt.Printf("Network: %+v\n", insights.Network)
	}
	fmt.Printf("Filtered Windows Services: %d services found\n", len(insights.WindowsServices))
}
