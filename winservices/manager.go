package winservices

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

var (
	ErrServiceNotFound = errors.New("windows service not found")
)

// WinSvcManager implements Windows service management using the golang.org/x/sys/windows/svc/mgr package.
// It manages a connection to the Windows service manager with lazy initialization.
type WinSvcManager struct {
	mgr     *mgr.Mgr
	mgrInit sync.Once
	mgrErr  error
}

// NewWinSvcManager creates a new instance of WinSvcManager for managing Windows services.
func NewWinSvcManager() *WinSvcManager {
	return &WinSvcManager{}
}

// Connect connects to the Windows service manager.
// Does not need to be called explicitly, as it is called
// automatically when needed.
// Can be used to check if the service manager is available.
func (s *WinSvcManager) Connect() error {
	// If already connected, do nothing
	if s.mgr != nil && s.mgrErr == nil {
		return nil
	}
	// Connect to the service manager
	s.mgrInit.Do(func() {
		s.mgr, s.mgrErr = mgr.Connect()
	})
	return s.mgrErr
}

// Disconnect disconnects from the service manager.
// It should be called when the service manager is no longer needed.
func (s *WinSvcManager) Disconnect() error {
	if s.mgr == nil {
		return nil
	}
	if err := s.mgr.Disconnect(); err != nil {
		return fmt.Errorf("failed to disconnect from service manager: %w", err)
	}
	s.mgr = nil
	return nil
}

// ListServices retrieves a list of all services on the system.
// It returns a slice of strings containing the names of the services.
func (s *WinSvcManager) ListServices() ([]string, error) {
	// Check if the service manager is connected
	if err := s.Connect(); err != nil {
		return nil, fmt.Errorf("failed to connect to service manager: %w", err)
	}
	// List all services
	return s.mgr.ListServices()
}

// ServiceExists checks if a service exists by its name.
// It returns true if the service exists, false otherwise.
func (s *WinSvcManager) ServiceExists(name string) (bool, error) {
	// Check if the service manager is connected
	if err := s.Connect(); err != nil {
		return false, fmt.Errorf("failed to connect to service manager: %w", err)
	}
	// Check if the service exists
	return s.serviceExists(name)
}

// GetServiceState retrieves the state of a specific service.
// It returns the service state as svc.State and a string representation.
func (s *WinSvcManager) GetServiceState(name string) (state svc.State, stateStr string, err error) {
	// Check if the service manager is connected
	if err := s.Connect(); err != nil {
		return 0, "", fmt.Errorf("failed to connect to service manager: %w", err)
	}

	// Check if the service exists
	exists, err := s.serviceExists(name)
	if err != nil {
		return
	}
	if !exists {
		return 0, "", ErrServiceNotFound
	}
	service, err := s.mgr.OpenService(name)
	if err != nil {
		println("Received error: ", err.Error())

		if strings.Contains(err.Error(), "ERROR_SERVICE_DOES_NOT_EXIST") {
			return 0, "", ErrServiceNotFound
		}
		return
	}
	defer func() { _ = service.Close() }()

	status, err := service.Query()
	if err != nil {
		return
	}

	return status.State, stateToString(status.State), nil
}

// It would be better to use error from mgr.Mgr methods, but there is nothing like this in the API
// Using just error string string representation is not reliable
func (s *WinSvcManager) serviceExists(name string) (bool, error) {
	services, err := s.mgr.ListServices()
	if err != nil {
		return false, fmt.Errorf("failed to list services: %w", err)
	}
	for _, service := range services {
		if strings.EqualFold(service, name) {
			return true, nil
		}
	}
	return false, nil
}
