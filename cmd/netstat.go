package cmd

import (
	"fmt"
	"log"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/mwiater/ghost/utils"
	"github.com/shirou/gopsutil/v4/net"
	"github.com/spf13/cobra"
)

// NetstatCmd defines the Cobra command for displaying active network connections on the system.
// It gathers information about current TCP/UDP connections and presents them in a table format.
var NetstatCmd = &cobra.Command{
	Use:   "netstat",
	Short: "Displays active network connections on the system",
	Run: func(cmd *cobra.Command, args []string) {
		connections, err := RunNetstat()
		if err != nil {
			log.Fatalf("Error fetching network connections: %v", err)
		}
		PrintConnections(connections)
	},
}

// init registers the NetstatCmd with the root command when this package is imported.
func init() {
	RootCmd.AddCommand(NetstatCmd)
}

// RunNetstat retrieves all active network connections without printing.
func RunNetstat() ([]net.ConnectionStat, error) {
	return net.Connections("all")
}

// PrintConnections formats and displays the network connection information.
// It presents connection details including protocol, local address, remote address, and state using the go-pretty table package.
func PrintConnections(conns []net.ConnectionStat) {
	// Create a new table using utils.Table function for consistent styling
	t := utils.Table("DarkSimple", "Active Network Connections")
	t.AppendHeader(table.Row{"Protocol", "Local Address", "Remote Address", "State"})

	for _, conn := range conns {
		protocol := mapProtocol(conn.Type)
		localAddr := fmt.Sprintf("%s:%d", conn.Laddr.IP, conn.Laddr.Port)
		remoteAddr := fmt.Sprintf("%s:%d", conn.Raddr.IP, conn.Raddr.Port)
		if conn.Raddr.IP == "" {
			remoteAddr = "N/A"
		}

		// Append each connection to the table
		t.AppendRow(table.Row{
			protocol,
			localAddr,
			remoteAddr,
			conn.Status,
		})
	}

	// Render the table
	fmt.Println()
	t.Render()
	fmt.Println("Netstat complete.")
}

// mapProtocol converts a protocol type represented by an integer into a human-readable string.
// It maps the protocol type to either "TCP", "UDP", or "UNKNOWN" based on the provided value.
func mapProtocol(protocolType uint32) string {
	switch protocolType {
	case 1:
		return "TCP"
	case 2:
		return "UDP"
	default:
		return "UNKNOWN"
	}
}
