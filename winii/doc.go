// Package winii provides Windows inventory insights.
//
// This package offers a high-level API to gather comprehensive insights about a Windows system,
// including system details, installed runtimes, network configurations, and Windows services.
//
// The main entry points are ProvideAllInventoryInsight() for a full set of insights with default options,
// and ProvideMinimalInventoryInsight() for customizable insights.
//
// Example:
//
//	insights, err := winii.ProvideAllInventoryInsight()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("CPU cores: %d\n", insights.System.CPU.Count)
package winii
