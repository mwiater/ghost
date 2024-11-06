package cmd

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/mwiater/ghost/utils"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

// PortDetail holds comprehensive information about an open port.
type PortDetail struct {
	Port     int
	Process  string
	PID      string
	Owner    string
	Protocol string
	State    string
	Local    string
	Foreign  string
}

// PortScannerCmd defines the Cobra command for scanning a range of ports on a specified host.
var PortScannerCmd = &cobra.Command{
	Use:   "portscanner",
	Short: "Scans a range of ports on a specified host",
	Run: func(cmd *cobra.Command, args []string) {
		host, _ := cmd.Flags().GetString("host")
		startPort, _ := cmd.Flags().GetInt("start-port")
		endPort, _ := cmd.Flags().GetInt("end-port")

		numWorkers := runtime.NumCPU() // Limit concurrency to the number of available CPUs
		openPorts := RunPortScanner(host, startPort, endPort, numWorkers)
		PrintPortScanSummary(openPorts, host)
	},
}

// init registers the PortScannerCmd with the root command and defines command-line flags.
func init() {
	RootCmd.AddCommand(PortScannerCmd)
	PortScannerCmd.Flags().StringP("host", "H", "localhost", "Host to scan")
	PortScannerCmd.Flags().IntP("start-port", "s", 1, "Starting port to scan")
	PortScannerCmd.Flags().IntP("end-port", "e", 1024, "Ending port to scan")
}

// RunPortScanner executes the port scanning process for a specified host and port range without printing.
func RunPortScanner(host string, startPort, endPort, numWorkers int) []PortDetail {
	return scanPortsConcurrently(host, startPort, endPort, numWorkers)
}

// scanPortsConcurrently scans ports using multiple workers, displaying progress with a progress bar.
func scanPortsConcurrently(host string, start, end, numWorkers int) []PortDetail {
	var openPorts []PortDetail
	var mu sync.Mutex // Mutex to protect access to the openPorts slice

	totalPorts := end - start + 1
	updateFrequency := 20 // Frequency of progress bar updates
	progressBar := progressbar.NewOptions(totalPorts,
		progressbar.OptionSetDescription("Scanning ports"),
		progressbar.OptionFullWidth(),
	)

	var wg sync.WaitGroup
	portCh := make(chan int, numWorkers) // Channel to distribute ports to workers
	var scannedPorts int                 // Track total number of ports scanned

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for port := range portCh {
				if scanPort(host, port) {
					details := getPortDetails(host, port)
					mu.Lock()
					openPorts = append(openPorts, details)
					mu.Unlock()
				}

				// Increment the number of scanned ports
				mu.Lock()
				scannedPorts++
				// Only update the progress bar every nth port scan
				if scannedPorts%updateFrequency == 0 || scannedPorts == totalPorts {
					progressBar.Describe(fmt.Sprintf("%d/%d ports scanned (%d open ports so far)", scannedPorts, totalPorts, len(openPorts)))
					progressBar.Add(updateFrequency)
				}
				mu.Unlock()
			}
		}()
	}

	// Distribute ports to workers
	for port := start; port <= end; port++ {
		portCh <- port
	}
	close(portCh) // Close the channel to signal workers to stop

	// Wait for all workers to finish
	wg.Wait()

	fmt.Println() // Print a new line after progress bar completes
	return openPorts
}

// scanPort checks if a specific port on the host is open by attempting to establish a TCP connection.
func scanPort(host string, port int) bool {
	address := net.JoinHostPort(host, strconv.Itoa(port))
	conn, err := net.DialTimeout("tcp", address, 1*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// PrintPortScanSummary displays the final summary of the port scan results in a table.
func PrintPortScanSummary(openPorts []PortDetail, host string) {
	fmt.Println("\n--- Port Scan Summary ---")

	// Prepare the table using utils.Table for consistent formatting
	t := utils.Table("DarkSimple", "Port Scan Results")
	t.AppendHeader(table.Row{
		"Port",
		"Protocol",
		"Local Address",
		"Foreign Address",
		"State",
		"Process",
		"PID",
		"Owner",
	})

	if len(openPorts) == 0 {
		// If no open ports are found, show a message
		t.AppendRow(table.Row{"-", "-", "-", "-", "-", "-", "-", "-"})
	} else {
		// Add each open port to the table
		for _, port := range openPorts {
			t.AppendRow(table.Row{
				strconv.Itoa(port.Port),
				port.Protocol,
				port.Local,
				port.Foreign,
				port.State,
				port.Process,
				port.PID,
				port.Owner,
			})
		}
	}

	// Render the table
	fmt.Println()
	t.Render()
	fmt.Println()
}

// getPortDetails retrieves detailed port information depending on the operating system.
func getPortDetails(host string, port int) PortDetail {
	var detail PortDetail
	detail.Port = port

	switch runtime.GOOS {
	case "linux", "darwin":
		// On Linux/macOS, use lsof to find the process using the open port
		cmd := exec.Command("lsof", "-i", fmt.Sprintf("TCP:%d", port), "-sTCP:LISTEN")
		output, err := cmd.Output()
		if err != nil {
			detail.Process = "N/A"
			detail.PID = "N/A"
			detail.Owner = "N/A"
			detail.Protocol = "TCP"
			detail.State = "LISTEN"
			detail.Local = fmt.Sprintf("%s:%d", host, port)
			detail.Foreign = "N/A"
			return detail
		}

		// Parse lsof output
		scanner := bufio.NewScanner(bytes.NewReader(output))
		firstLine := true
		for scanner.Scan() {
			line := scanner.Text()
			if firstLine {
				// Skip header line
				firstLine = false
				continue
			}
			fields := strings.Fields(line)
			if len(fields) >= 9 {
				detail.Process = fields[0]
				detail.PID = fields[1]
				detail.Owner = fields[2]
				detail.Protocol = fields[4]
				detail.State = fields[7]
				detail.Local = fields[8]
				detail.Foreign = "N/A"
				break
			}
		}

	case "windows":
		// On Windows, use netstat to find the process using the open port
		cmd := exec.Command("netstat", "-ano")
		output, err := cmd.Output()
		if err != nil {
			detail.Process = "N/A"
			detail.PID = "N/A"
			detail.Owner = "N/A"
			detail.Protocol = "N/A"
			detail.State = "N/A"
			detail.Local = fmt.Sprintf("%s:%d", host, port)
			detail.Foreign = "N/A"
			return detail
		}

		// Parse netstat output
		scanner := bufio.NewScanner(bytes.NewReader(output))
		found := false
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, fmt.Sprintf(":%d ", port)) || strings.Contains(line, fmt.Sprintf(":%d\r", port)) {
				fields := strings.Fields(line)
				if len(fields) >= 5 {
					detail.Protocol = fields[0]
					detail.Local = fields[1]
					detail.Foreign = fields[2]
					detail.State = fields[3]
					detail.PID = fields[4]
					found = true
					break
				}
			}
		}

		if found {
			// Get process name from PID
			processName := getProcessNameWindows(detail.PID)
			detail.Process = processName

			// Get owner (username) from PID
			owner := getProcessOwnerWindows(detail.PID)
			detail.Owner = owner
		} else {
			detail.Process = "N/A"
			detail.PID = "N/A"
			detail.Owner = "N/A"
			detail.Protocol = "TCP"
			detail.State = "LISTEN"
			detail.Local = fmt.Sprintf("%s:%d", host, port)
			detail.Foreign = "N/A"
		}
	default:
		detail.Process = "Unsupported OS"
		detail.PID = "N/A"
		detail.Owner = "N/A"
		detail.Protocol = "N/A"
		detail.State = "N/A"
		detail.Local = fmt.Sprintf("%s:%d", host, port)
		detail.Foreign = "N/A"
	}

	return detail
}

// getProcessNameWindows retrieves the process name given a PID on Windows.
func getProcessNameWindows(pid string) string {
	cmd := exec.Command("tasklist", "/FI", fmt.Sprintf("PID eq %s", pid), "/FO", "CSV", "/NH")
	output, err := cmd.Output()
	if err != nil {
		return "N/A"
	}

	// Parse CSV output
	fields := parseCSVLine(string(output))
	if len(fields) >= 1 {
		return fields[0]
	}
	return "N/A"
}

// getProcessOwnerWindows retrieves the owner (username) given a PID on Windows.
func getProcessOwnerWindows(pid string) string {
	// Use 'wmic' to get the owner
	cmd := exec.Command("wmic", "process", "where", fmt.Sprintf("ProcessId=%s", pid), "get", "Owner", "/FORMAT:LIST")
	output, err := cmd.Output()
	if err != nil {
		return "N/A"
	}

	// Parse the output
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Owner=") {
			owner := strings.TrimPrefix(line, "Owner=")
			return owner
		}
	}
	return "N/A"
}

// parseCSVLine parses a single CSV line and returns the fields.
func parseCSVLine(line string) []string {
	reader := csv.NewReader(strings.NewReader(line))
	reader.FieldsPerRecord = -1 // Variable number of fields
	records, err := reader.Read()
	if err != nil {
		return []string{}
	}
	return records
}
