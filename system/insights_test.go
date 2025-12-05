package system

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSystemInsights(t *testing.T) {
	insights, err := GetSystemInsights()
	require.NoError(t, err)
	require.NotNil(t, insights)
	t.Logf("System Insights: %+v", insights)
}
