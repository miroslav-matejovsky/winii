// Package network provides utilities for retrieving network insights on Windows systems.
package network

import (
	"net"
	"os"

	"github.com/yusufpapurcu/wmi"
)

// Insights holds network-related information for the system.
type Insights struct {
	DomainName   string   // DomainName is the domain the computer is joined to.
	HostName     string   // HostName is the name of the computer.
	IPAddresses  []string // IPAddresses contains IPv4 and IPv6 addresses.
	MACAddresses []string // MACAddresses contains MAC addresses of network adapters.
	Gateway      string   // Gateway is the default gateway IP address.
	DNSServers   []string // DNSServers contains DNS server IP addresses.
}

// Win32_ComputerSystem represents the WMI class for computer system information.
type Win32_ComputerSystem struct {
	Domain string // Domain is the domain name.
}

// Win32_NetworkAdapterConfiguration represents the WMI class for network adapter configurations.
type Win32_NetworkAdapterConfiguration struct {
	DefaultIPGateway     []string // DefaultIPGateway contains default gateway IPs.
	DNSServerSearchOrder []string // DNSServerSearchOrder contains DNS server IPs.
	IPAddress            []string // IPAddress contains IP addresses.
	MACAddress           string   // MACAddress is the MAC address.
}

// ProvideNetworkInsights retrieves and returns network insights using WMI queries.
// It populates an Insights struct with hostname, domain, IP addresses, MAC addresses, gateway, and DNS servers.
// Returns an error if hostname retrieval fails.
func ProvideNetworkInsights() (*Insights, error) {
	insights := &Insights{}

	// Get HostName
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	insights.HostName = hostname

	// Get DomainName using WMI (native Windows)
	var systems []Win32_ComputerSystem
	err = wmi.Query("SELECT Domain FROM Win32_ComputerSystem", &systems)
	if err == nil && len(systems) > 0 {
		insights.DomainName = systems[0].Domain
	} else {
		insights.DomainName = "" // Fallback if not in domain
	}

	// Get network adapter configurations using WMI
	var adapters []Win32_NetworkAdapterConfiguration
	err = wmi.Query("SELECT DefaultIPGateway, DNSServerSearchOrder, IPAddress, MACAddress FROM Win32_NetworkAdapterConfiguration WHERE IPEnabled = TRUE", &adapters)
	if err == nil {
		for _, adapter := range adapters {
			// Gateway (take first if multiple)
			if len(adapter.DefaultIPGateway) > 0 && insights.Gateway == "" {
				insights.Gateway = adapter.DefaultIPGateway[0]
			}
			// DNS Servers
			insights.DNSServers = append(insights.DNSServers, adapter.DNSServerSearchOrder...)
			// MAC Address
			if adapter.MACAddress != "" {
				insights.MACAddresses = append(insights.MACAddresses, adapter.MACAddress)
			}
			// IP Addresses (IPv4 and IPv6)
			for _, ip := range adapter.IPAddress {
				if net.ParseIP(ip) != nil {
					insights.IPAddresses = append(insights.IPAddresses, ip)
				}
			}
		}
		// Remove duplicates if any
		insights.IPAddresses = removeDuplicates(insights.IPAddresses)
		insights.MACAddresses = removeDuplicates(insights.MACAddresses)
		insights.DNSServers = removeDuplicates(insights.DNSServers)
	}

	return insights, nil
}

// removeDuplicates removes duplicate strings from a slice while preserving order.
func removeDuplicates(slice []string) []string {
	keys := make(map[string]bool)
	var result []string
	for _, item := range slice {
		if !keys[item] {
			keys[item] = true
			result = append(result, item)
		}
	}
	return result
}
