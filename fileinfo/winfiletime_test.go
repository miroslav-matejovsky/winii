package fileinfo

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGetFileTime(t *testing.T) {
	wf, err := NewWinFileInfo("C:\\Windows\\System32\\notepad.exe")
	require.NoError(t, err)

	ft, err := wf.GetFileTime()
	require.NoError(t, err)

	tenYearsAgo := time.Now().AddDate(-10, 0, 0)
	now := time.Now()

	t.Logf("Creation Time: %s", ft.CreationTime.Format(time.RFC1123))
	require.True(t, ft.CreationTime.After(tenYearsAgo))
	require.True(t, ft.CreationTime.Before(now))

	t.Logf("Last Access Time: %s", ft.LastAccessTime.Format(time.RFC1123))
	require.True(t, ft.LastAccessTime.After(tenYearsAgo))
	require.True(t, ft.LastAccessTime.Before(now))

	t.Logf("Last Write Time: %s", ft.LastWriteTime.Format(time.RFC1123))
	require.True(t, ft.LastWriteTime.After(tenYearsAgo))
	require.True(t, ft.LastWriteTime.Before(now))
}
