// package runtimes provides functionality to audit installer runtimes on Windows systems.
// It includes tools for checking installed Visual C++ Redistributables and .NET runtime components.
package runtimes

import "fmt"

type AuditResult struct {
	VCRedistRuntimes []VCRedistRuntime
	DotNetRuntimes   []DotNetRuntime
}

func DoAudit() (*AuditResult, error) {
	vcRedist, err := DoVCRedistAudit()
	if err != nil {
		return nil, fmt.Errorf("failed to audit Visual C++ Redistributable runtimes: %w", err)
	}
	dotNet, err := DotNetRuntimesAuditResult()
	if err != nil {
		return nil, fmt.Errorf("failed to audit .NET runtimes: %w", err)
	}
	return &AuditResult{
		VCRedistRuntimes: vcRedist,
		DotNetRuntimes:   dotNet,
	}, nil
}
