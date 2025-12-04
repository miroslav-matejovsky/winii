package runtimes

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAudit(t *testing.T) {
	auditResult, err := DoAudit()
	require.NoError(t, err)
	require.NotNil(t, auditResult)
	require.NotNil(t, auditResult.VCRedistRuntimes)
	require.NotNil(t, auditResult.DotNetRuntimes)
}
