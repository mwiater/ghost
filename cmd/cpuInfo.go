package cmd

import (
	"fmt"
	"reflect"
	"runtime"
	"sync"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/mwiater/ghost/utils"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/spf13/cobra"
)

// CpuInfoCmd represents the cpuinfo command
var CpuInfoCmd = &cobra.Command{
	Use:   "cpuinfo",
	Short: "Displays detailed CPU information such as model, cores, and frequency.",
	Long:  `Retrieves and displays detailed information about the CPU, including model name, number of cores, and base frequency.`,
	Run: func(cmd *cobra.Command, args []string) {
		cpuDetails, err := RunCpuInfo()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		PrintCpuInfo(cpuDetails)
	},
}

// CpuInfo holds details about the CPU.
type CpuInfo struct {
	ModelName string
	Cores     int
	Frequency string
}

// GetCpuInfo gathers CPU information using gopsutil/cpu with concurrency.
func GetCpuInfo() ([]CpuInfo, error) {
	// Get CPU info
	infoStats, err := cpu.Info()
	if err != nil {
		return nil, err
	}

	// Limit concurrency to the number of available CPUs
	concurrency := runtime.NumCPU()
	var wg sync.WaitGroup
	wg.Add(len(infoStats))

	cpuDetails := make([]CpuInfo, len(infoStats))
	sem := make(chan struct{}, concurrency)

	for i, cpuStat := range infoStats {
		sem <- struct{}{}
		go func(i int, cpuStat cpu.InfoStat) {
			defer func() {
				wg.Done()
				<-sem
			}()
			cpuDetails[i] = CpuInfo{
				ModelName: cpuStat.ModelName,
				Cores:     int(cpuStat.Cores),
				Frequency: fmt.Sprintf("%.2f GHz", cpuStat.Mhz/1000),
			}
		}(i, cpuStat)
	}

	wg.Wait()
	return cpuDetails, nil
}

// RunCpuInfo retrieves CPU information without printing.
func RunCpuInfo() ([]CpuInfo, error) {
	return GetCpuInfo()
}

// PrintCpuInfo displays the CPU information in a formatted table.
func PrintCpuInfo(cpuDetails []CpuInfo) {
	// Use utils.Table to create a table with "DarkSimple" style for alternate row shading
	t := utils.Table("DarkSimple", "cpuInfoCmd")
	t.AppendHeader(table.Row{"CPU Info", "Value"})

	for i, cpuInfo := range cpuDetails {
		t.AppendRow(table.Row{fmt.Sprintf("CPU %d", i+1), ""})

		// Reflect on cpuInfo struct to iterate through fields for display
		v := reflect.ValueOf(cpuInfo)
		for j := 0; j < v.NumField(); j++ {
			field := v.Type().Field(j)
			value := v.Field(j).Interface()
			t.AppendRow(table.Row{field.Name, value})
		}
	}

	// Render the table
	fmt.Println()
	t.Render()
	fmt.Println()
}

func init() {
	RootCmd.AddCommand(CpuInfoCmd)
}
