package cmd

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/mwiater/ghost/utils"
	gopsutilNet "github.com/shirou/gopsutil/v4/net"
	"github.com/spf13/cobra"
)

// NetworkInterfacesInterface defines an interface for retrieving network interface information.
type NetworkInterfacesInterface interface {
	Info() ([]gopsutilNet.InterfaceStat, error)
}

// NetworkInterfacesInfo implements NetworkInterfacesInterface using gopsutil/net
// to fetch network interface information.
type NetworkInterfacesInfo struct{}

// Info lists all the network interfaces on the host. It returns a slice of gopsutilNet.InterfaceStat
// and an error if there's an issue during retrieval.
func (NetworkInterfacesInfo) Info() ([]gopsutilNet.InterfaceStat, error) {
	return gopsutilNet.Interfaces()
}

// GetNetworkInterfacesInfo retrieves a list of network interfaces on the host using the provided
// NetworkInterfacesInterface. It returns a slice of gopsutilNet.InterfaceStat and an error in case of failure.
func GetNetworkInterfacesInfo(networkInterfacesInterface NetworkInterfacesInterface) ([]gopsutilNet.InterfaceStat, error) {
	return networkInterfacesInterface.Info()
}

// NetworkInterfacesCmd represents the `networkinterfaces` command, which lists all
// network interfaces available on the host.
var NetworkInterfacesCmd = &cobra.Command{
	Use:   "networkinterfaces",
	Short: "Lists all network interfaces on the host.",
	Long: `Gathers information on each network interface available on the system,
which is important for network configuration and troubleshooting.`,
	Run: func(cmd *cobra.Command, args []string) {
		networkInterfacesInfo, err := RunNetworkInterfacesInfo()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		PrintNetworkInterfacesInfo(networkInterfacesInfo)
	},
}

// RunNetworkInterfacesInfo retrieves network interfaces information without printing.
func RunNetworkInterfacesInfo() ([]gopsutilNet.InterfaceStat, error) {
	networkInterfacesInterface := NetworkInterfacesInfo{}
	return GetNetworkInterfacesInfo(networkInterfacesInterface)
}

// PrintNetworkInterfacesInfo displays the network interfaces information in a formatted table.
func PrintNetworkInterfacesInfo(networkInterfacesInfo []gopsutilNet.InterfaceStat) {
	// Use utils.Table to create a table with "DarkSimple" style for alternate row shading
	t := utils.Table("DarkSimple", "networkInterfacesCmd")
	t.AppendHeader(table.Row{"Network Interfaces", "Value"})

	// Iterate over the network interfaces and populate the table
	for _, iface := range networkInterfacesInfo {
		v := reflect.ValueOf(iface)
		for i := 0; i < v.NumField(); i++ {
			field := v.Type().Field(i)
			value := v.Field(i).Interface()

			// Handle fields that contain slices or nested structs
			switch reflect.TypeOf(value).Kind() {
			case reflect.Slice:
				var strValues []string
				s := reflect.ValueOf(value)
				for j := 0; j < s.Len(); j++ {
					elem := s.Index(j).Interface()
					if reflect.TypeOf(elem).Kind() == reflect.Struct {
						// If the slice contains structs, extract and display their fields
						nestedValue := reflect.ValueOf(elem)
						var nestedStrValues []string
						for k := 0; k < nestedValue.NumField(); k++ {
							nestedField := nestedValue.Type().Field(k)
							nestedFieldValue := nestedValue.Field(k).Interface()
							nestedStrValues = append(nestedStrValues, fmt.Sprintf("%s: %v", nestedField.Name, nestedFieldValue))
						}
						strValues = append(strValues, strings.Join(nestedStrValues, ", "))
					} else {
						strValues = append(strValues, fmt.Sprintf("%v", elem))
					}
				}

				// Add each element of the slice to the table
				for idx, strValue := range strValues {
					if idx == 0 {
						t.AppendRow(table.Row{field.Name, strValue})
					} else {
						t.AppendRow(table.Row{"", strValue})
					}
				}
			default:
				t.AppendRow(table.Row{field.Name, value})
			}
		}

		t.AppendRow(table.Row{"----------", ""})
	}

	// Render the table
	fmt.Println()
	t.Render()
	fmt.Println()
}

// init initializes the `networkInterfaces` command and adds it to the RootCmd.
func init() {
	RootCmd.AddCommand(NetworkInterfacesCmd)
}
