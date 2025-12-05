package runtimes

import (
	"bufio"
	"errors"
	"os/exec"
	"strings"
)

// DotNetRuntime represents a .NET runtime installation on a Windows system.
type DotNetRuntime struct {
	// Full name of the runtime, one of "Microsoft.AspNetCore.App", "Microsoft.NETCore.App", "Microsoft.WindowsDesktop.App"
	Type     string
	Version  string
	Location string
}

// ProvideDotNetInsights holds the insights about installed .NET runtimes.
// It uses dotnet --list-runtimes to gather information about installed runtimes.
func ProvideDotNetInsights() ([]DotNetRuntime, error) {
	// Check if dotnet command exists
	if _, err := exec.LookPath("dotnet"); err != nil {
		// dotnet not found, return empty list
		return []DotNetRuntime{}, nil
	}

	cmd := exec.Command("dotnet", "--list-runtimes")
	output, err := cmd.Output()
	if err != nil {
		return nil, errors.New("failed to execute dotnet --list-runtimes: " + err.Error())
	}

	var runtimes []DotNetRuntime
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, " ", 3)
		if len(parts) != 3 {
			return nil, errors.New("unexpected output format from dotnet --list-runtimes: " + line)
		}
		// Assume format: Name Version [Path]
		fullName := parts[0]
		version := parts[1]
		path := strings.Trim(parts[2], "[]")
		runtime := DotNetRuntime{
			Type:     fullName,
			Version:  version,
			Location: path,
		}
		runtimes = append(runtimes, runtime)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return runtimes, nil
}
