package network

import (
	"net"
	"os"

	"github.com/yusufpapurcu/wmi"
)

type NetworkInsights struct {
	DomainName  string
	HostName    string
	IPAddresses []string
}

type Win32_ComputerSystem struct {
	Domain string
}

func ProvideNetworkInsights() (*NetworkInsights, error) {
	insights := &NetworkInsights{}

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

	// Get IPAddresses (IPv4 only, excluding loopback)
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && ipnet.IP.To4() != nil {
				insights.IPAddresses = append(insights.IPAddresses, ipnet.IP.String())
			}
		}
	}

	return insights, nil
}
