package cmd

import (
	"fmt"
	"reflect"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/mwiater/ghost/utils"
	"github.com/shirou/gopsutil/mem"
	"github.com/spf13/cobra"
)

// MemInfoCmd represents the meminfo command
var MemInfoCmd = &cobra.Command{
	Use:   "meminfo",
	Short: "Displays memory usage statistics, including total, used, and free memory.",
	Long:  `Retrieves and displays memory usage information, including total memory, used memory, free memory, and memory usage percentage.`,
	Run: func(cmd *cobra.Command, args []string) {
		memInfo, err := RunMemInfo()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		PrintMemInfo(memInfo)
	},
}

// MemInfo holds details about memory usage.
type MemInfo struct {
	Total       string
	Used        string
	Free        string
	UsedPercent string
}

// GetMemInfo gathers memory information using gopsutil/mem.
func GetMemInfo() (*MemInfo, error) {
	// Retrieve memory stats
	v, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	// Populate MemInfo struct
	memInfo := &MemInfo{
		Total:       fmt.Sprintf("%.2f GB", float64(v.Total)/1e9),
		Used:        fmt.Sprintf("%.2f GB", float64(v.Used)/1e9),
		Free:        fmt.Sprintf("%.2f GB", float64(v.Free)/1e9),
		UsedPercent: fmt.Sprintf("%.2f%%", v.UsedPercent),
	}

	return memInfo, nil
}

// RunMemInfo retrieves memory information without printing.
func RunMemInfo() (*MemInfo, error) {
	return GetMemInfo()
}

// PrintMemInfo displays the memory information in a formatted table.
func PrintMemInfo(memInfo *MemInfo) {
	// Use utils.Table to create a table with "DarkSimple" style for alternate row shading
	t := utils.Table("DarkSimple", "memInfoCmd")
	t.AppendHeader(table.Row{"Memory Info", "Value"})

	// Reflect on memInfo struct to iterate through fields for display
	v := reflect.ValueOf(memInfo).Elem()
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
	RootCmd.AddCommand(MemInfoCmd)
}
