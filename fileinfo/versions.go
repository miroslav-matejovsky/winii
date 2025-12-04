package fileinfo

import (
	"fmt"

	"golang.org/x/sys/windows"
)

type Versions struct {
	FileVersion    WinFileVersion
	ProductVersion WinFileVersion
}

type WinFileVersion struct {
	Major uint16
	Minor uint16
	Patch uint16
	Build uint16
}

func (f WinFileVersion) String() string {
	return fmt.Sprintf("%d.%d.%d.%d", f.Major, f.Minor, f.Patch, f.Build)
}

// newWinFileInfo creates a new WinFileInfo from the given VS_FIXEDFILEINFO.
func newWinFileInfo(vsFixedInfo *windows.VS_FIXEDFILEINFO) *Versions {
	return &Versions{
		FileVersion: WinFileVersion{
			Major: uint16(vsFixedInfo.FileVersionMS >> 16),
			Minor: uint16(vsFixedInfo.FileVersionMS & 0xffff),
			Patch: uint16(vsFixedInfo.FileVersionLS >> 16),
			Build: uint16(vsFixedInfo.FileVersionLS & 0xffff),
		},
		ProductVersion: WinFileVersion{
			Major: uint16(vsFixedInfo.ProductVersionMS >> 16),
			Minor: uint16(vsFixedInfo.ProductVersionMS & 0xffff),
			Patch: uint16(vsFixedInfo.ProductVersionLS >> 16),
			Build: uint16(vsFixedInfo.ProductVersionLS & 0xffff),
		},
	}
}
