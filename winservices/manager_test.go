package winservices

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestListServices(t *testing.T) {
	service := NewWinSvcManager()

	services, err := service.ListServices()
	require.NoError(t, err, "Failed to list services")
	require.NotEmpty(t, services, "Service list should not be empty")
}

func TestServiceExists(t *testing.T) {
	serviceManager := NewWinSvcManager()

	// Test with a known service name
	t.Run("ServiceExists", func(t *testing.T) {
		exists, err := serviceManager.ServiceExists("EventLog")
		require.NoError(t, err, "Failed to check if service exists")
		require.True(t, exists, "Service should exist")
	})

	t.Run("ServiceExistsLowercase", func(t *testing.T) {
		// Test with a known service name in lowercase
		exists, err := serviceManager.ServiceExists("eventlog")
		require.NoError(t, err, "Failed to check if service exists")
		require.True(t, exists, "Service should exist")
	})

	t.Run("ServiceExistsUppercase", func(t *testing.T) {
		// Test with a known service name in uppercase
		exists, err := serviceManager.ServiceExists("EVENTLOG")
		require.NoError(t, err, "Failed to check if service exists")
		require.True(t, exists, "Service should exist")
	})

	// Test with an unknown service name
	t.Run("ServiceDoesNotExist", func(t *testing.T) {
		exists, err := serviceManager.ServiceExists("UnknownServiceABC")
		require.NoError(t, err, "Failed to check if service exists")
		require.False(t, exists, "Service should not exist")
	})
}

func TestGetServiceState(t *testing.T) {
	serviceManager := NewWinSvcManager()

	// Test with a known service name
	t.Run("GetKnownServiceState", func(t *testing.T) {
		state, stateStr, err := serviceManager.GetServiceState("EventLog")
		require.NoError(t, err, "Failed to get service state")
		require.NotEmpty(t, stateStr, "Service state string should not be empty")
		require.NotEmpty(t, state, "Service state should not be empty")
	})

	// Test with a service that is not running
	t.Run("GetKnownServiceStateUppercase", func(t *testing.T) {
		// Test with a known service name in uppercase
		state, stateStr, err := serviceManager.GetServiceState("EVENTLOG")
		require.NoError(t, err, "Failed to get service state")
		require.NotEmpty(t, stateStr, "Service state string should not be empty")
		require.NotEmpty(t, state, "Service state should not be empty")
	})

	t.Run("GetUnknownServiceState", func(t *testing.T) {
		// Test with an unknown service name
		_, _, err := serviceManager.GetServiceState("UnknownServiceCDE")
		require.Error(t, err, "Expected error for unknown service")
		require.Equal(t, ErrServiceNotFound, err)
	})
}
