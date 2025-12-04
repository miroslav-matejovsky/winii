package runtimes

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVcRedistRuntimesAudit(t *testing.T) {
	vcRedist, err := DoVCRedistAudit()
	require.NoError(t, err)
	require.NotNil(t, vcRedist)
	for _, v := range vcRedist {
		t.Logf("Version: %s, Architecture: %s, Installed: %t\nPath: %s", v.Version, v.Architecture, v.Installed, v.Path)
		require.NotEmpty(t, v.Version)
		require.NotEmpty(t, v.Architecture)
		require.NotEmpty(t, v.Path)
		require.Contains(t, []string{"x64", "x86", "X64", "X86"}, v.Architecture)
	}
}
