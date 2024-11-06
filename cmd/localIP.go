package cmd

import (
	"fmt"
	"net"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/mwiater/ghost/utils"
	"github.com/spf13/cobra"
)

// LocalIPCmd represents the `localip` command, which finds and returns the first
// internal IPv4 address, typically within the "192.168" subnet. This is useful for
// services that need to bind to an internal network interface.
var LocalIPCmd = &cobra.Command{
	Use:   "localip",
	Short: "Finds an internal IPv4 address.",
	Long: `Searches for and returns the first internal IPv4 address, typically 
within the "192.168" subnet. If none is found, it returns an error. This is 
useful for services that need to bind to an internal network interface.`,
	Run: func(cmd *cobra.Command, args []string) {
		localIP, err := RunLocalIP()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		PrintLocalIP(localIP)
	},
}

// RunLocalIP retrieves the internal IPv4 address of the local machine.
// It returns the IP address as a string and an error if no address is found.
func RunLocalIP() (string, error) {
	localIP, err := GetInternalIPv4()
	if err != nil {
		return "", err
	}
	return localIP, nil
}

// PrintLocalIP displays the local IP address in a formatted table.
func PrintLocalIP(localIP string) {
	// Create and configure a table to display the local IP address
	t := utils.Table("DarkSimple", "localIPCmd")
	t.AppendHeader(table.Row{"Local IP", "Value"})
	t.AppendRow(table.Row{"IP Address", localIP})

	fmt.Println()
	t.Render()
	fmt.Println()
}

// GetInternalIPv4 searches for and returns the first internal IPv4 address it finds,
// typically one that starts with "192.168". If no such address is found, it returns an error.
func GetInternalIPv4() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		ipNet, ok := addr.(*net.IPNet)
		if ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			ip := ipNet.IP.String()
			if ip[:7] == "192.168" {
				return ip, nil
			}
		}
	}

	return "", fmt.Errorf("no internal IPv4 address found")
}

// init initializes the `localIP` command and adds it to the RootCmd.
// This command allows users to find and display the first internal IPv4 address.
func init() {
	RootCmd.AddCommand(LocalIPCmd)
}
