package runtimes

import "golang.org/x/sys/windows/registry"

// WinReg defines an interface for Windows registry operations.
type WinReg interface {
	// OpenKeyReadOnly opens a registry key for read-only access.
	OpenKeyReadOnly(path string) (registry.Key, error)
}

// NewWinReg creates a new instance of WinReg.
func NewWinReg() WinReg {
	return &winReg{}
}

type winReg struct {
}

// OpenKeyReadOnly opens a registry key for read-only access.
func (w *winReg) OpenKeyReadOnly(path string) (registry.Key, error) {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, path, registry.QUERY_VALUE|registry.ENUMERATE_SUB_KEYS|registry.READ)

	if err != nil {
		return k, err
	}
	return k, nil
}
