package runtimes

import (
	"fmt"

	"golang.org/x/sys/windows/registry"
)

// VCRedistRuntime represents information about a Visual C++ Redistributable runtime.
// It includes version, architecture, installation status, and registry path.
type VCRedistRuntime struct {
	Version      string
	Architecture string // either x86 or x64
	Installed    bool
	Path         string
}

// ProvideVCRedistInsights scans the Windows registry for installed Visual C++ Redistributable runtimes.
// It returns a slice of VcRedistRuntime structs or an error if the scan fails.
// This is equivalent to running the PowerShell command:
// Get-ItemProperty -Path "HKLM:\SOFTWARE\WOW6432Node\Microsoft\VisualStudio\*\VC\Runtimes\*" | Format-Table
func ProvideVCRedistInsights() ([]VCRedistRuntime, error) {
	var vcRedistRuntimes []VCRedistRuntime
	winreg := NewWinReg()
	keyPath := `SOFTWARE\WOW6432Node\Microsoft\VisualStudio`
	k, err := winreg.OpenKeyReadOnly(keyPath)
	if err != nil {
		if err == registry.ErrNotExist {
			return vcRedistRuntimes, nil
		}
		return nil, fmt.Errorf("failed to open registry key: %w", err)
	}
	defer func() { _ = k.Close() }()

	subKeys, err := k.ReadSubKeyNames(-1)
	if err != nil {
		return nil, fmt.Errorf("failed to read sub keys: %w", err)
	}

	for _, subKey := range subKeys {
		runtimeKeyPath := keyPath + `\` + subKey + `\VC\Runtimes`
		runtimeKey, err := winreg.OpenKeyReadOnly(runtimeKeyPath)
		if err != nil {
			if err == registry.ErrNotExist {
				continue
			}
			return nil, fmt.Errorf("failed to open runtime registry key %s: %w", runtimeKeyPath, err)
		}
		defer func() { _ = runtimeKey.Close() }()

		runtimeSubKeys, err := runtimeKey.ReadSubKeyNames(-1)
		if err != nil {
			return nil, fmt.Errorf("failed to read runtime sub keys %s : %w", runtimeSubKeys, err)
		}
		for _, runtimeSubKey := range runtimeSubKeys {
			regPath := runtimeKeyPath + `\` + runtimeSubKey
			vcRedist, err := winreg.OpenKeyReadOnly(regPath)
			if err != nil {
				return nil, fmt.Errorf("failed to open full version registry key: %w", err)
			}
			defer func() { _ = vcRedist.Close() }()

			version, _, err := vcRedist.GetStringValue("Version")
			if err != nil {
				if err == registry.ErrNotExist {
					continue
				}
				return nil, fmt.Errorf("failed to read full version %s: %w", regPath, err)
			}
			installed, _, err := vcRedist.GetIntegerValue("Installed")
			if err != nil {
				if err == registry.ErrNotExist {
					continue
				}
				return nil, fmt.Errorf("failed to read installed %s: %w", regPath, err)
			}

			// Collecting redist runtimes
			vcRedistRuntimes = append(vcRedistRuntimes, VCRedistRuntime{
				Version:      version,
				Architecture: runtimeSubKey,
				Installed:    installed == 1,
				Path:         "HKLM\\" + regPath,
			})
		}
	}

	return vcRedistRuntimes, nil

}
