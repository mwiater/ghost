package cmd

import (
	"bufio"
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/mwiater/ghost/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// TracerouteCmd represents the traceroute command
var TracerouteCmd = &cobra.Command{
	Use:   "traceroute",
	Short: "Performs a traceroute to a specified IP address.",
	Long:  `Executes a traceroute from the current location to a specified IP address and displays detailed hop information.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Retrieve flags
		destination := viper.GetString("destination")
		maxHops := viper.GetInt("maxHops")

		// Execute traceroute
		hops, err := RunTraceroute(destination, maxHops)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		// Display traceroute results
		PrintTraceroute(hops)
	},
}

// TracerouteHop holds details about a single hop in the traceroute.
type TracerouteHop struct {
	HopNumber int
	Hostname  string
	IP        string
	RTTs      [3]string
}

// RunTraceroute executes the traceroute command and retrieves the hop information.
func RunTraceroute(destination string, maxHops int) ([]TracerouteHop, error) {
	return GetTraceroute(destination, maxHops)
}

// PrintTraceroute displays the traceroute hops in a formatted table.
func PrintTraceroute(hops []TracerouteHop) {
	t := utils.Table("DarkSimple", "tracerouteCmd")
	t.AppendHeader(table.Row{"Hop", "Hostname", "IP Address", "RTT1 (ms)", "RTT2 (ms)", "RTT3 (ms)"})

	for _, hop := range hops {
		t.AppendRow(table.Row{
			hop.HopNumber,
			hop.Hostname,
			hop.IP,
			hop.RTTs[0],
			hop.RTTs[1],
			hop.RTTs[2],
		})
	}

	fmt.Println()
	t.Render()
	fmt.Println()
}

func init() {
	RootCmd.AddCommand(TracerouteCmd)

	// Define flags with default values
	TracerouteCmd.PersistentFlags().StringP("destination", "d", "google.com", "Destination IP address or hostname for traceroute")
	TracerouteCmd.PersistentFlags().IntP("maxHops", "m", 30, "Maximum number of hops to trace")

	// Bind flags to viper
	viper.BindPFlag("destination", TracerouteCmd.PersistentFlags().Lookup("destination"))
	viper.BindPFlag("maxHops", TracerouteCmd.PersistentFlags().Lookup("maxHops"))
}

// GetTraceroute retrieves traceroute information based on the operating system.
func GetTraceroute(destination string, maxHops int) ([]TracerouteHop, error) {
	if runtime.GOOS == "windows" {
		return getTracerouteWindows(destination, maxHops)
	}
	return getTracerouteUnix(destination, maxHops)
}

// getTracerouteUnix retrieves traceroute information on Unix-based systems (Linux, macOS).
func getTracerouteUnix(destination string, maxHops int) ([]TracerouteHop, error) {
	var hops []TracerouteHop

	// Check if 'traceroute' is available
	cmdName := "traceroute"
	if _, err := exec.LookPath(cmdName); err != nil {
		return nil, fmt.Errorf("'traceroute' command not found. Please install it to use this feature.")
	}

	// Prepare the command arguments
	args := []string{"-m", strconv.Itoa(maxHops), destination}

	cmd := exec.Command(cmdName, args...)

	// Capture combined output (stdout and stderr)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to execute 'traceroute' command: %v", err)
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	lineNumber := 0
	currentHop := TracerouteHop{}

	for scanner.Scan() {
		line := scanner.Text()
		lineNumber++

		// Skip the first line which typically contains the destination info
		if strings.HasPrefix(line, "traceroute") {
			continue
		}

		// Handle lines like "1?: [LOCALHOST] pmtu 1500"
		if strings.Contains(line, "pmtu") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		// Parse hop number
		hopNumStr := fields[0]
		hopNumStr = strings.TrimSuffix(hopNumStr, "?")
		hopNumStr = strings.TrimSuffix(hopNumStr, ":")
		hopNum, err := strconv.Atoi(hopNumStr)
		if err != nil {
			continue // Skip lines that don't start with a hop number
		}

		// Initialize or reset TracerouteHop
		if currentHop.HopNumber != hopNum {
			// If we are starting a new hop, append the previous one if it exists
			if currentHop.HopNumber != 0 {
				hops = append(hops, currentHop)
			}
			currentHop = TracerouteHop{
				HopNumber: hopNum,
				Hostname:  "-",
				IP:        "-",
				RTTs:      [3]string{"-", "-", "-"},
			}
		}

		// Determine if the line contains '*' indicating a timeout
		if strings.Contains(line, "*") || strings.Contains(line, "no reply") {
			currentHop.RTTs = [3]string{"*", "*", "*"}
			continue
		}

		// Extract hostname and IP
		// Check if IP is in parentheses
		if strings.Contains(line, "(") && strings.Contains(line, ")") {
			parts := strings.SplitN(line, "(", 2)
			currentHop.Hostname = strings.TrimSpace(parts[0])
			ipPart := strings.SplitN(parts[1], ")", 2)[0]
			currentHop.IP = ipPart
			// Extract RTTs
			rtts := extractRTTs(parts[1])
			currentHop.RTTs = rtts
		} else {
			// No hostname, only IP
			currentHop.Hostname = "-"
			currentHop.IP = fields[1]
			// Extract RTTs
			rtts := extractRTTs(line)
			currentHop.RTTs = rtts
		}
	}

	// Append the last hop if it exists
	if currentHop.HopNumber != 0 {
		hops = append(hops, currentHop)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading traceroute output: %v", err)
	}

	return hops, nil
}

// getTracerouteWindows retrieves traceroute information on Windows systems.
func getTracerouteWindows(destination string, maxHops int) ([]TracerouteHop, error) {
	var hops []TracerouteHop

	// Windows uses 'tracert' command
	// '/h' specifies the maximum number of hops
	cmd := exec.Command("tracert", "-h", strconv.Itoa(maxHops), destination)

	// Capture combined output (stdout and stderr)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to execute 'tracert' command: %v", err)
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))

	for scanner.Scan() {
		line := scanner.Text()

		// Skip header lines
		if strings.HasPrefix(line, "Tracing route to") || strings.HasPrefix(line, "over a maximum of") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 5 {
			continue
		}

		// Parse hop number
		hopNumStr := fields[0]
		hopNum, err := strconv.Atoi(hopNumStr)
		if err != nil {
			continue // Skip lines that don't start with a hop number
		}

		hop := TracerouteHop{
			HopNumber: hopNum,
			Hostname:  "-",
			IP:        "-",
			RTTs:      [3]string{"-", "-", "-"},
		}

		// Check for timeout
		if strings.Contains(line, "Request timed out.") {
			hop.RTTs = [3]string{"*", "*", "*"}
			hops = append(hops, hop)
			continue
		}

		// Extract RTTs and IP/Hostname
		// Example line: "  2     2 ms     2 ms     2 ms  10.0.0.1"
		if len(fields) >= 5 {
			rtt1 := fields[1]
			rtt2 := fields[2]
			rtt3 := fields[3]
			hostIP := fields[4]

			// Attempt to separate hostname and IP if available
			hostname := "-"
			ip := hostIP
			if strings.Contains(hostIP, "(") && strings.Contains(hostIP, ")") {
				// Hostname and IP are present
				parts := strings.SplitN(hostIP, "(", 2)
				hostname = strings.TrimSpace(parts[0])
				ip = strings.TrimSuffix(strings.TrimSpace(parts[1]), ")")
			}

			hop.Hostname = hostname
			hop.IP = ip
			hop.RTTs[0] = strings.TrimSuffix(rtt1, "ms")
			hop.RTTs[1] = strings.TrimSuffix(rtt2, "ms")
			hop.RTTs[2] = strings.TrimSuffix(rtt3, "ms")
		}

		hops = append(hops, hop)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading tracert output: %v", err)
	}

	return hops, nil
}

// extractRTTs parses RTT values from a traceroute line.
// It returns an array of three RTTs as strings.
func extractRTTs(line string) [3]string {
	var rtts [3]string
	// Split the line by "ms" to extract RTT values
	parts := strings.Split(line, "ms")
	count := 0
	for i := 0; i < len(parts)-1 && count < 3; i++ {
		rtt := strings.TrimSpace(parts[i])
		if rtt == "*" || rtt == "?" || strings.ToLower(rtt) == "no" {
			rtts[count] = "*"
		} else {
			rtts[count] = rtt
		}
		count++
	}
	// Fill remaining RTTs with "-"
	for ; count < 3; count++ {
		rtts[count] = "-"
	}
	return rtts
}
