// +build windows

package cmd

import (
	"fmt"
	"net"
	"os/exec"
	"strings"
)

// runARPScan performs ARP scanning on Windows using the 'arp -a' command.
func runARPScan() ([]ARPResult, error) {
	cmd := exec.Command("arp", "-a")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error running arp -a: %w", err)
	}

	var results []ARPResult
	lines := strings.Split(string(output), "\n")

	// Parse the output to extract IP and MAC addresses
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) >= 2 && net.ParseIP(fields[0]) != nil {
			results = append(results, ARPResult{
				IPAddress:  fields[0],
				MACAddress: fields[1],
			})
		}
	}

	return results, nil
}
