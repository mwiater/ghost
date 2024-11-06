package cmd

import (
	"fmt"
	"reflect"
	"runtime"
	"sync"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/mwiater/ghost/utils"
	"github.com/shirou/gopsutil/disk"
	"github.com/spf13/cobra"
)

// FsInfoCmd represents the fsinfo command
var FsInfoCmd = &cobra.Command{
	Use:   "fsinfo",
	Short: "Displays filesystem information, including type, total space, and available space.",
	Long:  `Retrieves and displays information about each mounted filesystem, such as the filesystem type, total space, used space, and available space.`,
	Run: func(cmd *cobra.Command, args []string) {
		fsDetails, err := RunFsInfo()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		PrintFsInfo(fsDetails)
	},
}

// FsInfo holds details about each filesystem.
type FsInfo struct {
	Filesystem     string
	Type           string
	TotalSpace     string
	UsedSpace      string
	AvailableSpace string
	UsedPercent    string
}

// GetFsInfo gathers filesystem information for each mounted volume, using concurrency.
func GetFsInfo() ([]FsInfo, error) {
	// Retrieve list of partitions
	partitions, err := disk.Partitions(false)
	if err != nil {
		return nil, err
	}

	concurrency := runtime.NumCPU()
	var wg sync.WaitGroup
	wg.Add(len(partitions))

	fsDetails := make([]FsInfo, len(partitions))
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

			fsDetails[i] = FsInfo{
				Filesystem:     partition.Device,
				Type:           usageStat.Fstype,
				TotalSpace:     fmt.Sprintf("%.2f GB", float64(usageStat.Total)/1e9),
				UsedSpace:      fmt.Sprintf("%.2f GB", float64(usageStat.Used)/1e9),
				AvailableSpace: fmt.Sprintf("%.2f GB", float64(usageStat.Free)/1e9),
				UsedPercent:    fmt.Sprintf("%.2f%%", usageStat.UsedPercent),
			}
		}(i, partition)
	}

	wg.Wait()
	return fsDetails, nil
}

// RunFsInfo retrieves filesystem information without printing.
func RunFsInfo() ([]FsInfo, error) {
	return GetFsInfo()
}

// PrintFsInfo displays the filesystem information in a formatted table.
func PrintFsInfo(fsDetails []FsInfo) {
	// Use utils.Table to create a table with "DarkSimple" style for alternate row shading
	t := utils.Table("DarkSimple", "fsInfoCmd")
	t.AppendHeader(table.Row{"Filesystem", "Type", "Total Space", "Used Space", "Available Space", "Used Percent"})

	// Add each filesystem's details to the table
	for _, fsInfo := range fsDetails {
		v := reflect.ValueOf(fsInfo)
		row := make([]interface{}, v.NumField())
		for i := 0; i < v.NumField(); i++ {
			row[i] = v.Field(i).Interface()
		}
		t.AppendRow(row)
	}

	// Render the table
	fmt.Println()
	t.Render()
	fmt.Println()
}

func init() {
	RootCmd.AddCommand(FsInfoCmd)
}
