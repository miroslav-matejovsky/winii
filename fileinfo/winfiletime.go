package fileinfo

import (
	"fmt"
	"time"

	"golang.org/x/sys/windows"
)

type FileTime struct {
	CreationTime   time.Time
	LastAccessTime time.Time
	LastWriteTime  time.Time
}

// getFileTime retrieves the creation, last access, and last write times of the file.
func (wf *WinFileInfo) getFileTime() (*FileTime, error) {
	// Convert path to UTF-16
	utf16Path, err := windows.UTF16PtrFromString(wf.path)
	if err != nil {
		return nil, fmt.Errorf("failed to convert path to UTF-16: %w", err)
	}

	// Open file with required access flags
	handle, err := windows.CreateFile(
		utf16Path,
		windows.FILE_READ_EA,
		windows.FILE_SHARE_READ,
		nil,
		windows.OPEN_EXISTING,
		windows.FILE_FLAG_BACKUP_SEMANTICS,
		0,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer func() {
		if err := windows.Close(handle); err != nil {
			fmt.Printf("failed to close file handle: %v\n", err)
		}
	}()

	var ctime, atime, wtime windows.Filetime
	err = windows.GetFileTime(handle, &ctime, &atime, &wtime)
	if err != nil {
		return nil, fmt.Errorf("failed to get file time: %w", err)
	}
	return &FileTime{
		CreationTime:   time.Unix(0, ctime.Nanoseconds()),
		LastAccessTime: time.Unix(0, atime.Nanoseconds()),
		LastWriteTime:  time.Unix(0, wtime.Nanoseconds()),
	}, nil
}
