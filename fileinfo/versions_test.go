package fileinfo

import (
	"testing"

	"github.com/bi-zone/go-fileversion"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetFileVersionInfo(t *testing.T) {
	file := `C:\Windows\System32\notepad.exe`
	bizone, err := fileversion.New(file)
	require.NoError(t, err)
	expected := bizone.FixedInfo()

	wf, err := NewWinFileInfo(file)
	require.NoError(t, err)

	fi, err := wf.GetVersions()
	require.NoError(t, err)

	assert.Equal(t, expected.FileVersion.Major, fi.FileVersion.Major)
	assert.Equal(t, expected.FileVersion.Minor, fi.FileVersion.Minor)
	// !!! the bug is bi-zone go-fileversion library, the Patch and Build are swapped
	// https://github.com/bi-zone/go-fileversion/issues/3
	assert.Equal(t, expected.FileVersion.Build, fi.FileVersion.Patch)
	assert.Equal(t, expected.FileVersion.Patch, fi.FileVersion.Build)
}
