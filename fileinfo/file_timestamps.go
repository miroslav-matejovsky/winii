package fileinfo

import (
	"fmt"
	"time"

	"golang.org/x/sys/windows"
)

type FileTimestamps struct {
	CreationTime   time.Time
	LastAccessTime time.Time
	LastWriteTime  time.Time
}

// GetFileTime retrieves the creation, last access, and last write times of the file.
func GetFileTimestamps(path string) (*FileTimestamps, error) {
	// Convert path to UTF-16
	utf16Path, err := windows.UTF16PtrFromString(path)
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

	var createTime, accessTime, writeTime windows.Filetime
	err = windows.GetFileTime(handle, &createTime, &accessTime, &writeTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get file time: %w", err)
	}
	return &FileTimestamps{
		CreationTime:   time.Unix(0, createTime.Nanoseconds()),
		LastAccessTime: time.Unix(0, accessTime.Nanoseconds()),
		LastWriteTime:  time.Unix(0, writeTime.Nanoseconds()),
	}, nil
}
