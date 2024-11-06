package cmd

import (
	"fmt"
	"net"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/mwiater/ghost/utils"
	"github.com/spf13/cobra"
)

// SubnetCalcCmd defines the Cobra command for calculating network details from a given IP address and subnet (CIDR).
// It calculates and displays the network address, broadcast address, and IP range for the specified subnet.
var SubnetCalcCmd = &cobra.Command{
	Use:   "subnetcalc",
	Short: "Calculates network details for a given IP address and subnet (CIDR)",
	Run: func(cmd *cobra.Command, args []string) {
		cidr, _ := cmd.Flags().GetString("cidr")

		_, ipNet, err := net.ParseCIDR(cidr)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		subnetDetails := RunSubnetCalculator(ipNet)
		PrintSubnetDetails(subnetDetails)
	},
}

// init registers the SubnetCalcCmd with the root command and defines the --cidr flag.
func init() {
	RootCmd.AddCommand(SubnetCalcCmd)
	SubnetCalcCmd.Flags().StringP("cidr", "c", "192.168.1.0/24", "CIDR notation for subnet (e.g., 192.168.1.0/24)")
}

// SubnetDetails holds details about the calculated subnet information.
type SubnetDetails struct {
	NetworkAddress   string
	BroadcastAddress string
	IPRange          string
}

// RunSubnetCalculator calculates the network address, broadcast address, and IP range for the specified subnet.
func RunSubnetCalculator(ipNet *net.IPNet) *SubnetDetails {
	networkAddress := ipNet.IP.String()
	broadcastAddress := calculateBroadcastAddress(ipNet)
	ipRange := fmt.Sprintf("%s - %s", networkAddress, calculateLastIP(ipNet))

	return &SubnetDetails{
		NetworkAddress:   networkAddress,
		BroadcastAddress: broadcastAddress,
		IPRange:          ipRange,
	}
}

// PrintSubnetDetails displays the subnet details in a formatted table.
func PrintSubnetDetails(details *SubnetDetails) {
	// Create a table using utils.Table for consistent formatting
	t := utils.Table("DarkSimple", "Subnet Calculation Results")
	t.AppendHeader(table.Row{"Field", "Value"})

	// Prepare the data for the table
	data := [][]string{
		{"Network Address", details.NetworkAddress},
		{"Broadcast Address", details.BroadcastAddress},
		{"IP Range", details.IPRange},
	}

	// Add the data to the table
	for _, v := range data {
		t.AppendRow(table.Row{v[0], v[1]})
	}

	// Render the table
	fmt.Println()
	t.Render()
	fmt.Println("Subnet calculation complete.")
}

// calculateBroadcastAddress calculates the broadcast address for the given subnet.
// It performs bitwise operations using the IP address and subnet mask to derive the broadcast address.
func calculateBroadcastAddress(ipNet *net.IPNet) string {
	ip := ipNet.IP.To4()
	mask := ipNet.Mask
	broadcast := make(net.IP, len(ip))
	for i := 0; i < len(ip); i++ {
		broadcast[i] = ip[i] | ^mask[i]
	}
	return broadcast.String()
}

// calculateLastIP calculates the last IP address in the given subnet's IP range.
// It uses the network address and subnet mask to determine the highest IP in the range.
func calculateLastIP(ipNet *net.IPNet) string {
	ip := ipNet.IP.To4()
	mask := ipNet.Mask
	lastIP := make(net.IP, len(ip))
	for i := 0; i < len(ip); i++ {
		lastIP[i] = ip[i] | ^mask[i]
	}
	return lastIP.String()
}
