package cmd

import (
	"bufio"
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/mwiater/ghost/utils"
	"github.com/spf13/cobra"
)

// RouteCmd represents the route command
var RouteCmd = &cobra.Command{
	Use:   "routeinfo",
	Short: "Displays the IP routing table and network routes.",
	Long:  `Retrieves and displays the system's IP routing table and network routes in a formatted table.`,
	Run: func(cmd *cobra.Command, args []string) {
		routes, err := RunRoute()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		PrintRoutes(routes)
	},
}

// RouteEntry holds details about a single route.
type RouteEntry struct {
	Destination string
	Genmask     string
	Gateway     string
	Flags       string
	Metric      string
	Ref         string
	Use         string
	Iface       string
}

// RunRoute retrieves the routing table without printing.
func RunRoute() ([]RouteEntry, error) {
	return GetRoute()
}

// PrintRoutes displays the routing entries in a formatted table.
func PrintRoutes(routes []RouteEntry) {
	t := utils.Table("DarkSimple", "routeCmd")
	// Define table headers based on the operating system
	if runtime.GOOS == "windows" {
		t.AppendHeader(table.Row{"Network Destination", "Netmask", "Gateway", "Interface", "Metric"})
		for _, route := range routes {
			t.AppendRow(table.Row{
				route.Destination,
				route.Genmask,
				route.Gateway,
				route.Iface,
				route.Metric,
			})
		}
	} else {
		t.AppendHeader(table.Row{"Destination", "Genmask", "Gateway", "Flags", "Metric", "Ref", "Use", "Iface"})
		for _, route := range routes {
			t.AppendRow(table.Row{
				route.Destination,
				route.Genmask,
				route.Gateway,
				route.Flags,
				route.Metric,
				route.Ref,
				route.Use,
				route.Iface,
			})
		}
	}

	fmt.Println()
	t.Render()
	fmt.Println()
}

func init() {
	RootCmd.AddCommand(RouteCmd)

	// Define flags with default values (if any)
	// Currently, no additional flags are necessary for the 'route' command
	// However, you can add flags here if future enhancements are needed

	// Example:
	// RouteCmd.PersistentFlags().BoolP("json", "j", false, "Output in JSON format")
	// viper.BindPFlag("json", RouteCmd.PersistentFlags().Lookup("json"))
}

// GetRoute retrieves routing information based on the operating system.
func GetRoute() ([]RouteEntry, error) {
	if runtime.GOOS == "windows" {
		return getRouteWindows()
	}
	return getRouteUnix()
}

// getRouteUnix retrieves routing information on Unix-based systems (Linux, macOS).
func getRouteUnix() ([]RouteEntry, error) {
	var routes []RouteEntry

	// Use 'route -n' for Linux and 'netstat -rn' for macOS
	var cmd *exec.Cmd
	if runtime.GOOS == "darwin" {
		// macOS
		cmd = exec.Command("netstat", "-rn")
	} else {
		// Assume Linux
		cmd = exec.Command("route", "-n")
	}

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute routing command: %v", err)
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	lineNumber := 0
	for scanner.Scan() {
		line := scanner.Text()
		lineNumber++

		// Skip header lines
		if runtime.GOOS == "darwin" {
			if lineNumber < 3 {
				continue
			}
		} else {
			if lineNumber < 3 {
				continue
			}
		}

		// Split the line into fields
		fields := strings.Fields(line)
		if runtime.GOOS == "darwin" {
			// macOS netstat -rn output has columns:
			// Destination, Gateway, Flags, Refs, Use, Netif, Expire
			if len(fields) < 7 {
				continue
			}
			route := RouteEntry{
				Destination: fields[0],
				Genmask:     "N/A", // Not provided directly
				Gateway:     fields[1],
				Flags:       fields[2],
				Metric:      "N/A", // Not provided directly
				Ref:         fields[3],
				Use:         fields[4],
				Iface:       fields[5],
			}
			routes = append(routes, route)
		} else {
			// Linux route -n output has columns:
			// Destination, Gateway, Genmask, Flags, Metric, Ref, Use, Iface
			if len(fields) < 8 {
				continue
			}
			route := RouteEntry{
				Destination: fields[0],
				Genmask:     fields[2],
				Gateway:     fields[1],
				Flags:       fields[3],
				Metric:      fields[4],
				Ref:         fields[5],
				Use:         fields[6],
				Iface:       fields[7],
			}
			routes = append(routes, route)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading routing command output: %v", err)
	}

	return routes, nil
}

// getRouteWindows retrieves routing information on Windows systems.
func getRouteWindows() ([]RouteEntry, error) {
	var routes []RouteEntry

	// Use 'route print' command
	cmd := exec.Command("route", "print")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute 'route print' command: %v", err)
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	inIPv4Section := false
	for scanner.Scan() {
		line := scanner.Text()

		// Detect the IPv4 Route Table section
		if strings.Contains(line, "IPv4 Route Table") {
			inIPv4Section = true
			continue
		}

		if inIPv4Section {
			// Skip until headers are found
			if strings.HasPrefix(line, "===") || strings.HasPrefix(line, "Network Destination") {
				continue
			}

			// An empty line signifies the end of the IPv4 section
			if strings.TrimSpace(line) == "" {
				break
			}

			// Split the line into fields based on whitespace
			fields := strings.Fields(line)
			if len(fields) < 5 {
				continue
			}

			route := RouteEntry{
				Destination: fields[0],
				Genmask:     fields[1],
				Gateway:     fields[2],
				Iface:       fields[3],
				Metric:      fields[4],
				Flags:       "N/A", // Flags are not directly available
				Ref:         "N/A", // Not available
				Use:         "N/A", // Not available
			}
			routes = append(routes, route)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading 'route print' output: %v", err)
	}

	return routes, nil
}
