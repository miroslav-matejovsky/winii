package winservices

import (
	"fmt"
	"regexp"
)

type InventoryInsightOption struct {
	WinServiceNameRegex *regexp.Regexp
}

// ProvideInsights retrieves and optionally filters Windows services based on the provided options.
// It uses the pre-compiled WinServiceNameRegex from opts to filter services by name (case-insensitive match).
// If opts is nil, all services are returned.
func ProvideInsights(opts InventoryInsightOption) ([]ServiceDetails, error) {
	svcManager := NewWinSvcManager()
	defer func() { _ = svcManager.Disconnect() }()

	services, err := svcManager.AllServices(opts.WinServiceNameRegex)
	if err != nil {
		return nil, fmt.Errorf("failed to get all services: %w", err)
	}

	return services, nil
}

// AllServices retrieves details for all Windows services managed by this WinSvcManager,
// optionally filtering by the provided regexp (matches service names).
// If the regexp is empty or ".*", no filtering is applied.
func (s *WinSvcManager) AllServices(winServiceNameRegexp *regexp.Regexp) ([]ServiceDetails, error) {
	serviceNames, err := s.ListServices()
	if err != nil {
		return nil, fmt.Errorf("failed to list services: %w", err)
	}
	// Filter services if regexp is provided
	if winServiceNameRegexp != nil {
		var filtered []string
		for _, name := range serviceNames {
			if winServiceNameRegexp.MatchString(name) {
				filtered = append(filtered, name)
			}
		}
		serviceNames = filtered
	}
	// Collect details for each service
	var serviceDetails []ServiceDetails
	for _, serviceName := range serviceNames {
		details, err := s.GetServiceDetails(serviceName)
		if err != nil {
			return nil, fmt.Errorf("failed to get details for service %s: %w", serviceName, err)
		}
		serviceDetails = append(serviceDetails, *details)
	}
	return serviceDetails, nil
}
