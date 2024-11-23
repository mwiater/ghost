package cmd

import (
	"fmt"
	"runtime"
	"sync"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/mwiater/ghost/utils"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/spf13/cobra"
)

// DiskUsageCmd represents the diskusage command
var DiskUsageCmd = &cobra.Command{
	Use:   "diskusage",
	Short: "Displays disk usage information, including total, used, and free space.",
	Long:  `Retrieves and displays disk usage statistics for each mounted volume, including total space, used space, free space, and percentage used.`,
	Run: func(cmd *cobra.Command, args []string) {
		diskUsages, err := RunDiskUsage()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		PrintDiskUsage(diskUsages)
	},
}

// DiskUsage holds usage details about each disk.
type DiskUsage struct {
	MountPoint  string
	TotalSpace  string
	UsedSpace   string
	FreeSpace   string
	UsedPercent string
}

// GetDiskUsage gathers disk usage information for each mounted volume, using concurrency.
func GetDiskUsage() ([]DiskUsage, error) {
	// Retrieve list of partitions
	partitions, err := disk.Partitions(false)
	if err != nil {
		return nil, err
	}

	concurrency := runtime.NumCPU()
	var wg sync.WaitGroup
	wg.Add(len(partitions))

	diskUsages := make([]DiskUsage, len(partitions))
	sem := make(chan struct{}, concurrency)

	for i, partition := range partitions {
		sem <- struct{}{}
		go func(i int, partition disk.PartitionStat) {
			defer func() {
				wg.Done()
				<-sem
			}()

			// Get usage stats for the partition
			usageStat, err := disk.Usage(partition.Mountpoint)
			if err != nil {
				return
			}

			diskUsages[i] = DiskUsage{
				MountPoint:  partition.Mountpoint,
				TotalSpace:  fmt.Sprintf("%.2f GB", float64(usageStat.Total)/1e9),
				UsedSpace:   fmt.Sprintf("%.2f GB", float64(usageStat.Used)/1e9),
				FreeSpace:   fmt.Sprintf("%.2f GB", float64(usageStat.Free)/1e9),
				UsedPercent: fmt.Sprintf("%.2f%%", usageStat.UsedPercent),
			}
		}(i, partition)
	}

	wg.Wait()
	return diskUsages, nil
}

// RunDiskUsage retrieves the disk usage data without printing.
func RunDiskUsage() ([]DiskUsage, error) {
	return GetDiskUsage()
}

// PrintDiskUsage displays the disk usage information in a formatted table.
func PrintDiskUsage(diskUsages []DiskUsage) {
	// Use utils.Table to create a table with "DarkSimple" style for alternate row shading
	t := utils.Table("DarkSimple", "diskUsageCmd")
	t.AppendHeader(table.Row{"Mount Point", "Total Space", "Used Space", "Free Space", "Used Percent"})

	// Add each disk's usage details to the table
	for _, diskUsage := range diskUsages {
		t.AppendRow(table.Row{diskUsage.MountPoint, diskUsage.TotalSpace, diskUsage.UsedSpace, diskUsage.FreeSpace, diskUsage.UsedPercent})
	}

	// Render the table
	fmt.Println()
	t.Render()
	fmt.Println()
}

func init() {
	RootCmd.AddCommand(DiskUsageCmd)
}
