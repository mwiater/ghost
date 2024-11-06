//go:build linux
// +build linux

package cmd

import (
	"fmt"
	"os/exec"
	"strings"
)

// GetServices retrieves the list of running services for Linux.
func GetServices() ([]Service, error) {
	var services []Service

	// Use `ps` command to list processes with memory usage
	cmd := exec.Command("ps", "-eo", "comm,state,rss")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// Process each line of output
	lines := strings.Split(string(output), "\n")
	for _, line := range lines[1:] {
		fields := strings.Fields(line)
		if len(fields) >= 3 {
			name := fields[0]
			status := fields[1]
			memUsage := fmt.Sprintf("%.2f MB", parseMemory(fields[2])/1024.0)

			services = append(services, Service{
				Name:        name,
				Status:      status,
				MemoryUsage: memUsage,
			})
		}
	}

	return services, nil
}

// parseMemory converts memory usage from KB to float64 for formatting
func parseMemory(mem string) float64 {
	memKb := 0.0
	fmt.Sscanf(mem, "%f", &memKb)
	return memKb
}
