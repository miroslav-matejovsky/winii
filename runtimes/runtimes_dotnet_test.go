package runtimes

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDotNetRuntimesAudit(t *testing.T) {
	dotnet, err := DotNetRuntimesAuditResult()
	require.NoError(t, err)
	require.NotNil(t, dotnet)
	for _, d := range dotnet {
		t.Log("Type: ", d.Type)
		t.Log("Version: ", d.Version)
		t.Log("Location: ", d.Location)
		require.NotEmpty(t, d.Type)
		require.NotEmpty(t, d.Version)
		require.NotEmpty(t, d.Location)
		require.Contains(t, []string{"Microsoft.NETCore.App", "Microsoft.AspNetCore.App", "Microsoft.WindowsDesktop.App"}, d.Type)
	}

}
