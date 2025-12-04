package fileinfo

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

type winver struct {
	dataPointer unsafe.Pointer
	data        []byte
}

func initWinVer(filePath string) (*winver, error) {

	var zHandle windows.Handle
	// https://learn.microsoft.com/en-us/windows/win32/api/winver/nf-winver-getfileversioninfosizea
	size, err := windows.GetFileVersionInfoSize(filePath, &zHandle)
	if err != nil {
		return nil, fmt.Errorf("failed to get file version info size: %w", err)
	}

	// https://learn.microsoft.com/en-us/windows/win32/api/winver/nf-winver-getfileversioninfoa
	var ignoredHandle uint32 // described as Ignored in the documentation
	buffer := make([]byte, size)
	var lpData = unsafe.Pointer(&buffer[0])
	err = windows.GetFileVersionInfo(filePath, ignoredHandle, size, lpData)
	if err != nil {
		return nil, fmt.Errorf("failed to get file version info: %w", err)
	}
	return &winver{
		dataPointer: lpData,
		data:        buffer,
	}, nil
}

// https://learn.microsoft.com/en-us/windows/win32/api/verrsrc/ns-verrsrc-vs_fixedfileinfo
func (f *winver) queryFixedFileInfo() (*windows.VS_FIXEDFILEINFO, error) {
	var offset uintptr = 0
	var offsetPointer = unsafe.Pointer(&offset)
	var length uint32
	// https://learn.microsoft.com/en-us/windows/win32/api/winver/nf-winver-verqueryvaluea
	err := windows.VerQueryValue(f.dataPointer, `\`, offsetPointer, &length)
	if err != nil {
		return nil, fmt.Errorf("failed to query file version info: %w", err)
	}
	if length == 0 {
		return nil, fmt.Errorf("no version info found")
	}
	start := offset - uintptr(f.dataPointer)
	end := start + uintptr(length)
	data := f.data[start:end]
	fixedFileInfo := *(*windows.VS_FIXEDFILEINFO)(unsafe.Pointer(&data[0]))
	return &fixedFileInfo, nil
}
