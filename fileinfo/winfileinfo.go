// package fileinfo provides utilities for retrieving file information on Windows systems.
package fileinfo

import (
	"fmt"
	"os"

	"golang.org/x/sys/windows"
)

// WinFileInfo represents a file on the Windows filesystem.
// This file must exist in the OS
// afero in memory cannot be used because native Windows APIs are used to retrieve file information.
type WinFileInfo struct {
	path string
}

// NewWinFileInfo creates a new WinFile for the given path.
// It returns an error if the file does not exist or if there is an error checking the file.
func NewWinFileInfo(path string) (*WinFileInfo, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("file does not exist: %s", path)
	}
	if err != nil {
		return nil, fmt.Errorf("error checking file: %s", err)
	}
	return &WinFileInfo{path: path}, nil
}

// GetFileTime retrieves the file time information for the file.
// It returns a WinFileTime struct containing the file time information.
func (wf *WinFileInfo) GetFileTime() (*FileTime, error) {
	return wf.getFileTime()
}

// GetVersions retrieves the file version information for the file.
// It returns a WinFileInfo struct containing the file version information.
func (wf *WinFileInfo) GetVersions() (*Versions, error) {
	ffi, err := wf.GetFixedFileInfo()
	if err != nil {
		return nil, err
	}
	return newWinFileInfo(ffi), nil
}

// GetFixedFileInfo retrieves the fixed file information for the file.
// It returns a windows.VS_FIXEDFILEINFO struct containing the fixed file information.
func (wf *WinFileInfo) GetFixedFileInfo() (*windows.VS_FIXEDFILEINFO, error) {
	winver, err := initWinVer(wf.path)
	if err != nil {
		return nil, err
	}
	return winver.queryFixedFileInfo()
}

// GetCertificates retrieves the embedded certificates from the file.
// It returns a slice of x509.Certificate pointers or an error if the operation fails.
func (wf *WinFileInfo) GetCertificates() (*Certificates, error) {
	certs, err := getCertificates(wf.path)
	if err != nil {
		return nil, fmt.Errorf("failed to extract certificates: %v", err)
	}
	return certs, nil
}
