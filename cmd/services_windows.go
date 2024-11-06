//go:build windows
// +build windows

package cmd

import (
	"fmt"
	"os/exec"
	"strings"
)

// GetServices retrieves the list of running services for Windows.
func GetServices() ([]Service, error) {
	var services []Service

	// Use PowerShell command to get services and memory usage
	cmd := exec.Command("powershell", "-Command", "Get-Process | Select-Object Name, Status, WorkingSet")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// Process each line of output
	lines := strings.Split(string(output), "\n")
	for _, line := range lines[3:] { // Skipping the header lines
		fields := strings.Fields(line)
		if len(fields) >= 3 {
			name := fields[0]
			status := fields[1]
			memUsage := fmt.Sprintf("%.2f MB", parseMemory(fields[2])/1024.0/1024.0)

			services = append(services, Service{
				Name:        name,
				Status:      status,
				MemoryUsage: memUsage,
			})
		}
	}

	return services, nil
}

// parseMemory converts memory usage from bytes to float64 for formatting
func parseMemory(mem string) float64 {
	memBytes := 0.0
	fmt.Sscanf(mem, "%f", &memBytes)
	return memBytes
}
