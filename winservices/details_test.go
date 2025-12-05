package winservices

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNotExistingService(t *testing.T) {
	manager := NewWinSvcManager()

	// Test with a non-existing service name
	d, err := manager.GetServiceDetails("NonExistingService1234")
	require.ErrorIs(t, err, ErrServiceNotFound)
	require.Nil(t, d)
}

func TestServiceDetails(t *testing.T) {
	manager := NewWinSvcManager()

	// wuauserv is the Windows Update service, and it should be present on all Windows systems.
	d, err := manager.GetServiceDetails("wuauserv")
	require.NoError(t, err)
	assert.Equal(t, "wuauserv", d.ServiceName)
	assert.Equal(t, "Windows Update", d.DisplayName)
	assert.Equal(t, "Enables the detection, download, and installation of updates for Windows and other programs. If this service is disabled, users of this computer will not be able to use Windows Update or its automatic updating feature, and programs will not be able to use the Windows Update Agent (WUA) API.", d.Description)
	assert.Equal(t, "c:\\windows\\system32\\svchost.exe -k netsvcs -p", strings.ToLower(d.PathToExecutable))
	// Sometimes Manual, sometimes Automatic, sometimes Disabled?
	assert.Contains(t, []string{"Automatic", "Manual", "Disabled"}, d.StartupType)
	// Service status can be "Stopped" on the CI environment
	assert.Contains(t, []string{"Running", "Stopped"}, d.ServiceStatus)
	assert.Equal(t, "Win32 Share Process", d.ServiceType)
	assert.Equal(t, "Normal", d.ErrorControl)
	assert.Equal(t, []string{"rpcss"}, d.Dependencies)
	assert.Equal(t, "LocalSystem", d.ServiceStartName)
	assert.Equal(t, false, d.DelayedAutoStart)

	executableFile := d.Executable.ExecutableFile
	assert.Equal(t, `c:\windows\system32\svchost.exe`, strings.ToLower(executableFile.Path))
	assert.Regexp(t, `^\d+\.\d+\.\d+\.\d+$`, executableFile.Version)
	assert.Regexp(t, `^\d+\.\d+\.\d+\.\d+$`, executableFile.ProductVersion)
	assert.NotEqual(t, time.Time{}, executableFile.CreationTime)
	assert.NotEqual(t, time.Time{}, executableFile.LastAccessTime)
	assert.NotEqual(t, time.Time{}, executableFile.LastWriteTime)
	assert.Empty(t, d.Executable.ConfigFiles, "C:\\Windows	is skipped for config files collection")

	assert.Equal(t, "", d.Recovery.Command)
	assert.Equal(t, "Restart", d.Recovery.FirstFailure)
	assert.Equal(t, 1*time.Minute, d.Recovery.FirstFailureAfter)
	assert.Equal(t, "None", d.Recovery.SecondFailure)
	assert.Equal(t, 0*time.Second, d.Recovery.SecondFailureAfter)
	assert.Equal(t, "None", d.Recovery.SubsequentFailures)
	assert.Equal(t, 0*time.Second, d.Recovery.SubsequentFailuresAfter)
	assert.Equal(t, false, d.Recovery.MoreThan3Actions)
	assert.Equal(t, 24*time.Hour, d.Recovery.ResetFailCountAfter)
}
