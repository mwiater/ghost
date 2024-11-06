/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"reflect"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/mwiater/ghost/utils"
	"github.com/shirou/gopsutil/host"
	"github.com/spf13/cobra"
)

// HostInfoCmd represents the hostInfo command
var HostInfoCmd = &cobra.Command{
	Use:   "hostinfo",
	Short: "Delivers comprehensive host system information.",
	Long: ` Fetches detailed information about the host system, such as uptime,
boot time, and OS specifics, using the gopsutil package. It's a vital function
for system diagnostics and inventory.`,
	Run: func(cmd *cobra.Command, args []string) {
		hostInfo, err := RunHostInfo()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		PrintHostInfo(hostInfo)
	},
}

// HostInterface defines an interface for retrieving host information.
type HostInterface interface {
	Info() (*host.InfoStat, error)
}

// HostInfo implements HostInterface using gopsutil/host to fetch detailed host information.
type HostInfo struct{}

// Info retrieves detailed information about the host, including uptime, processes count,
// operating system, and platform. It returns a pointer to a host.InfoStat and an error if there's
// an issue during retrieval.
func (HostInfo) Info() (*host.InfoStat, error) {
	return host.Info()
}

// GetHostInfo retrieves detailed information about the host using the provided HostInterface.
// It returns a pointer to a host.InfoStat and an error in case of failure.
func GetHostInfo(hostInterface HostInterface) (*host.InfoStat, error) {
	return hostInterface.Info()
}

// RunHostInfo retrieves host information without printing.
func RunHostInfo() (*host.InfoStat, error) {
	hostInterface := HostInfo{}
	return GetHostInfo(hostInterface)
}

// PrintHostInfo displays the host information in a formatted table.
func PrintHostInfo(hostInfo *host.InfoStat) {
	// Use utils.Table to create a table with "DarkSimple" style for alternate row shading
	t := utils.Table("DarkSimple", "hostInfoCmd")
	t.AppendHeader(table.Row{"Host Info", "Value"})

	// Reflect on hostInfo struct to iterate through fields for display
	v := reflect.ValueOf(hostInfo).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		value := v.Field(i).Interface()
		t.AppendRow(table.Row{field.Name, value})
	}

	// Render the table
	fmt.Println()
	t.Render()
	fmt.Println()
}

func init() {
	RootCmd.AddCommand(HostInfoCmd)
}
