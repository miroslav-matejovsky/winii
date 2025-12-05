package system

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
)

// DiskInfo holds information for a single disk partition.
type DiskInfo struct {
	// Path is the mount point of the disk partition.
	Path string
	// Total is the total size of the disk in bytes.
	Total uint64
	// Used is the used space on the disk in bytes.
	Used uint64
	// Free is the free space on the disk in bytes.
	Free uint64
	// UsedPercent is the percentage of disk space used.
	UsedPercent float64
}

// Insights holds system insights including disk, RAM, CPU, and OS information.
type Insights struct {
	// Disks is a slice of disk information for all partitions.
	Disks []DiskInfo
	// RAMUsed is the amount of RAM currently used in bytes.
	RAMUsed uint64
	// RAMTotal is the total amount of RAM in bytes.
	RAMTotal uint64
	// WindowsVer is the version of the Windows operating system.
	WindowsVer string
	// CPUUsage is the current CPU usage percentage.
	CPUUsage float64
	// CPUMax is the number of logical CPU cores.
	CPUMax int
}

// GetSystemInsights retrieves system insights like disk info (all disks), RAM usage and total, Windows version, CPU usage, and max cores.
// It returns an Insights struct or an error if retrieval fails.
func GetSystemInsights() (*Insights, error) {
	insights := &Insights{}

	// Get all disk partitions and their usage
	partitions, err := disk.Partitions(true)
	if err != nil {
		return nil, fmt.Errorf("failed to get partitions: %v", err)
	}
	for _, p := range partitions {
		usage, err := disk.Usage(p.Mountpoint)
		if err != nil {
			continue // Skip errors for individual disks
		}
		insights.Disks = append(insights.Disks, DiskInfo{
			Path:        p.Mountpoint,
			Total:       usage.Total,
			Used:        usage.Used,
			Free:        usage.Free,
			UsedPercent: usage.UsedPercent,
		})
	}

	// Get RAM info
	vm, err := mem.VirtualMemory()
	if err != nil {
		return nil, fmt.Errorf("failed to get RAM info: %v", err)
	}
	insights.RAMUsed = vm.Used
	insights.RAMTotal = vm.Total

	// Get Windows version
	hostInfo, err := host.Info()
	if err != nil {
		return nil, fmt.Errorf("failed to get host info: %v", err)
	}
	insights.WindowsVer = hostInfo.PlatformVersion

	// Get CPU usage (current) and max (logical cores)
	percent, err := cpu.Percent(time.Second, false)
	if err != nil {
		return nil, fmt.Errorf("failed to get CPU percent: %v", err)
	}
	insights.CPUUsage = percent[0] // Overall CPU usage

	cores, err := cpu.Counts(true)
	if err != nil {
		return nil, fmt.Errorf("failed to get CPU counts: %v", err)
	}
	insights.CPUMax = cores

	return insights, nil
}
