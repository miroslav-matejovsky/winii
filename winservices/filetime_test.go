package winservices

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGetFileTime_ValidFile(t *testing.T) {
	// Create a temporary file
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test.txt")
	file, err := os.Create(tempFile)
	require.NoError(t, err)
	err = file.Close()
	require.NoError(t, err)

	// Get file times
	ft, err := getFileTime(tempFile)
	require.NoError(t, err)
	require.NotNil(t, ft)

	// Check that times are set and reasonable (not zero)
	now := time.Now()

	require.False(t, ft.CreationTime.IsZero(), "creation time should not be zero")
	require.False(t, ft.LastAccessTime.IsZero(), "last access time should not be zero")
	require.False(t, ft.LastWriteTime.IsZero(), "last write time should not be zero")

	require.True(t, ft.CreationTime.After(now.Add(-time.Minute)), "creation time should be recent")
	require.True(t, ft.CreationTime.Before(now.Add(time.Minute)), "creation time should be recent")
}

func TestGetFileTime_InvalidPath(t *testing.T) {
	_, err := getFileTime("nonexistent_file.txt")
	require.Error(t, err)
}
