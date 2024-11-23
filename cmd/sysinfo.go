package cmd

import (
	"fmt"
	"reflect"
	"runtime"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/mwiater/ghost/utils"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/spf13/cobra"
)

// SysInfoCmd represents the sysinfo command
var SysInfoCmd = &cobra.Command{
	Use:   "sysinfo",
	Short: "Displays system information such as OS, architecture, and uptime.",
	Long:  `Retrieves and displays details about the operating system, architecture, kernel version, and system uptime.`,
	Run: func(cmd *cobra.Command, args []string) {
		info, err := RunSysInfo()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		PrintSysInfo(info)
	},
}

// SystemInfo holds details about the system.
type SystemInfo struct {
	OS           string
	Architecture string
	Kernel       string
	Uptime       string
}

// GetSysInfo gathers system information based on the current platform.
func GetSysInfo() (*SystemInfo, error) {
	hostInfo, err := host.Info()
	if err != nil {
		return nil, err
	}

	uptime := time.Duration(hostInfo.Uptime) * time.Second
	info := &SystemInfo{
		OS:           fmt.Sprintf("%s %s", hostInfo.Platform, hostInfo.PlatformVersion),
		Architecture: runtime.GOARCH,
		Kernel:       hostInfo.KernelVersion,
		Uptime:       uptime.String(),
	}

	return info, nil
}

// RunSysInfo retrieves system information without printing.
func RunSysInfo() (*SystemInfo, error) {
	return GetSysInfo()
}

// PrintSysInfo displays the system information in a formatted table.
func PrintSysInfo(info *SystemInfo) {
	// Use utils.Table to create a table with "DarkSimple" style for alternate row shading
	t := utils.Table("DarkSimple", "sysInfoCmd")
	t.AppendHeader(table.Row{"System Info", "Value"})

	// Reflect on info struct to iterate through fields for display
	v := reflect.ValueOf(info).Elem()
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
	RootCmd.AddCommand(SysInfoCmd)
}
