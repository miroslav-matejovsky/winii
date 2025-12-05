package winii

import (
	"fmt"
	"log"
	"regexp"

	"github.com/miroslav-matejovsky/winii/winservices"
)

// ExampleProvideInventoryInsight demonstrates how to gather inventory insights
// with custom options, filtering Windows services by name.
func ExampleProvideInventoryInsight() {
	options := InsightsOptions{
		WinServiceOptions: winservices.InventoryInsightOption{
			WinServiceNameRegex: regexp.MustCompile(".*service.*"),
		},
	}
	insights, err := ProvideInventoryInsight(options)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Filtered Windows Services: %d services found\n", len(insights.WindowsServices))
}
