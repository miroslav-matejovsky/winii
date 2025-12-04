package fileinfo

import (
	"testing"

	"github.com/bi-zone/go-fileversion"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// bi-zone go-fileversion is used to retrieve the file version of a file for testing purposes.
// With the test we ensure that the file version is correctly retrieved.
func TestBiZoneFileVersion(t *testing.T) {
	file := `C:\Windows\System32\notepad.exe`
	gfv, err := fileversion.New(file)
	require.NoError(t, err)
	fixedInfo := gfv.FixedInfo()
	fv := fixedInfo.FileVersion
	// the actual version has Path and Build numbers swapped, there is an opened bug
	// https://github.com/bi-zone/go-fileversion/issues/3
	assert.Equal(t, "6.2.3624.26100", fv.String())
	assert.Equal(t, uint16(6), fv.Major)
	assert.Equal(t, uint16(2), fv.Minor)
	assert.Equal(t, uint16(3624), fv.Patch)
	assert.Equal(t, uint16(26100), fv.Build)
}
