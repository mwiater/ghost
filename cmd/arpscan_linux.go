//go:build linux
// +build linux

package cmd

import (
	"fmt"
	"net/netip"

	"github.com/mdlayher/arp"
)

// runARPScan performs ARP scanning on Unix-based systems (Linux/macOS).
func runARPScan() ([]ARPResult, error) {
	iface, err := getInterface()
	if err != nil {
		return nil, fmt.Errorf("error getting interface: %w", err)
	}

	conn, err := arp.Dial(iface)
	if err != nil {
		return nil, fmt.Errorf("error creating ARP connection: %w", err)
	}
	defer conn.Close()

	var results []ARPResult

	// Scan the local network and collect each result
	for ip := 1; ip < 255; ip++ {
		address := fmt.Sprintf("192.168.1.%d", ip)
		hwAddr, err := conn.Resolve(netip.MustParseAddr(address))
		if err != nil {
			continue
		}
		results = append(results, ARPResult{
			IPAddress:  address,
			MACAddress: hwAddr.String(),
		})
	}

	return results, nil
}
