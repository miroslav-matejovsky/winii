package winservices

import "fmt"

func ProvideInsights() ([]ServiceDetails, error) {
	svcManager := NewWinSvcManager()
	defer func() { _ = svcManager.Disconnect() }()

	services, err := svcManager.AllServices()
	if err != nil {
		return nil, fmt.Errorf("failed to get all services: %w", err)
	}
	return services, nil
}

func (s *WinSvcManager) AllServices() ([]ServiceDetails, error) {
	services, err := s.ListServices()
	if err != nil {
		return nil, fmt.Errorf("failed to list services: %w", err)
	}
	var serviceDetails []ServiceDetails
	for _, serviceName := range services {
		details, err := s.GetServiceDetails(serviceName)
		if err != nil {
			return nil, fmt.Errorf("failed to get details for service %s: %w", serviceName, err)
		}
		serviceDetails = append(serviceDetails, *details)
	}
	return serviceDetails, nil
}
