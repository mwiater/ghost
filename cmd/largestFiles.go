package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/mwiater/ghost/utils"
	"github.com/spf13/cobra"
)

// LargestFilesCmd represents the largestfiles command
var LargestFilesCmd = &cobra.Command{
	Use:   "largestfiles",
	Short: "Lists the largest files in a specified directory, sorted by size.",
	Long:  `Recursively searches for files in the specified directory (or current directory by default) and lists them in descending order by size, including their absolute paths.`,
	Run: func(cmd *cobra.Command, args []string) {
		dir, _ := cmd.Flags().GetString("directory")
		results, _ := cmd.Flags().GetInt("results")

		files, err := RunLargestFiles(dir, results)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		PrintLargestFiles(files)
	},
}

// FileSize holds details about a file and its size.
type FileSize struct {
	Path string
	Size int64
}

// GetLargestFiles retrieves files sorted by size in descending order using concurrency.
func GetLargestFiles(startDir string, maxResults int) ([]FileSize, error) {
	var files []FileSize
	var mu sync.Mutex
	var wg sync.WaitGroup
	concurrency := 10 // Set a limit to the number of concurrent goroutines
	sem := make(chan struct{}, concurrency)

	// Walk through the directory recursively
	err := filepath.WalkDir(startDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Only consider files (not directories)
		if !d.IsDir() {
			sem <- struct{}{} // Acquire a slot in the semaphore
			wg.Add(1)
			go func(path string) {
				defer func() {
					<-sem // Release the slot in the semaphore
					wg.Done()
				}()

				// Get file info concurrently
				fileInfo, err := os.Stat(path)
				if err != nil {
					return
				}

				// Lock access to the files slice while appending
				mu.Lock()
				files = append(files, FileSize{
					Path: path,
					Size: fileInfo.Size(),
				})
				mu.Unlock()
			}(path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Wait for all goroutines to complete
	wg.Wait()

	// Sort files by size in descending order
	sort.Slice(files, func(i, j int) bool {
		return files[i].Size > files[j].Size
	})

	// Limit the number of results
	if maxResults > 0 && len(files) > maxResults {
		files = files[:maxResults]
	}

	return files, nil
}

// RunLargestFiles retrieves the largest files without printing.
func RunLargestFiles(startDir string, maxResults int) ([]FileSize, error) {
	return GetLargestFiles(startDir, maxResults)
}

// PrintLargestFiles displays the largest files in a formatted table.
func PrintLargestFiles(files []FileSize) {
	if len(files) == 0 {
		fmt.Println("No files found.")
		return
	}

	// Use utils.Table to create a table with "DarkSimple" style for alternate row shading
	t := utils.Table("DarkSimple", "largestFilesCmd")
	t.AppendHeader(table.Row{"File Path", "Size (MB)"})

	// Add each file's details to the table
	for _, file := range files {
		sizeMB := float64(file.Size) / (1024 * 1024)
		t.AppendRow(table.Row{file.Path, fmt.Sprintf("%.2f", sizeMB)})
	}

	fmt.Println()
	t.Render()
	fmt.Println()
}

func init() {
	RootCmd.AddCommand(LargestFilesCmd)

	// Define flags for directory and results
	LargestFilesCmd.Flags().StringP("directory", "d", ".", "Directory to scan")
	LargestFilesCmd.Flags().IntP("results", "r", 20, "Number of results to display")
}
