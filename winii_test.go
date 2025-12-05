package winii

import (
	"fmt"
	"log"
	"regexp"

	"github.com/miroslav-matejovsky/winii/winii"
	"github.com/miroslav-matejovsky/winii/winservices"
)

// ExampleProvideAllInventoryInsight demonstrates how to gather all inventory insights
// for the current Windows system using default options.
func ExampleProvideAllInventoryInsight() {
	insights, err := ProvideAllInventoryInsight()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("System: %+v\n", insights.System)
	fmt.Printf("Runtimes: %+v\n", insights.Runtimes)
	fmt.Printf("Network: %+v\n", insights.Network)
	fmt.Printf("Windows Services: %d services found\n", len(insights.WindowsServices))
}

// ExampleProvideMinimalInventoryInsight demonstrates how to gather inventory insights
// with custom options, filtering Windows services by name.
func ExampleProvideMinimalInventoryInsight() {
	options := winii.InsightsOptions{
		WinServiceOptions: winservices.InventoryInsightOption{
			WinServiceNameRegex: regexp.MustCompile(".*service.*"),
		},
	}
	insights, err := ProvideMinimalInventoryInsight(options)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Filtered Windows Services: %d services found\n", len(insights.WindowsServices))
}
