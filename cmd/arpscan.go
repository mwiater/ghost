package cmd

import (
	"fmt"
	"net"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/mwiater/ghost/utils"
	"github.com/spf13/cobra"
)

// ARPScannerCmd is the command used to scan the local network using ARP.
var ARPScannerCmd = &cobra.Command{
	Use:   "arpscan",
	Short: "Scans the local network using ARP to find devices",
	Run: func(cmd *cobra.Command, args []string) {
		results, err := RunARPScanner()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		PrintArpScanResults(results)
	},
}

// ARPResult holds the IP and MAC address for each discovered device.
type ARPResult struct {
	IPAddress  string
	MACAddress string
}

// init registers the ARPScannerCmd with the root command when this package is imported.
func init() {
	RootCmd.AddCommand(ARPScannerCmd)
}

// RunARPScanner determines the operating system and calls the appropriate ARP scanning function.
func RunARPScanner() ([]ARPResult, error) {
	return runARPScan()
}

// PrintArpScanResults displays the ARP scan results in a formatted table.
func PrintArpScanResults(results []ARPResult) {
	// Use utils.Table to create a table with "DarkSimple" style for alternate row shading
	t := utils.Table("DarkSimple", "ARP Scan Results")
	t.AppendHeader(table.Row{"IP Address", "MAC Address"})

	// Add each ARP result's details to the table
	for _, result := range results {
		t.AppendRow(table.Row{result.IPAddress, result.MACAddress})
	}

	// Render the table
	fmt.Println()
	t.Render()
	fmt.Println()
}

// getInterface returns the first active network interface that is not a loopback interface.
func getInterface() (*net.Interface, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp != 0 && iface.Flags&net.FlagLoopback == 0 {
			return &iface, nil
		}
	}
	return nil, fmt.Errorf("no valid network interface found")
}
