// cmd/largestdirs.go
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/mwiater/ghost/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Dir represents a directory with its path, depth, and size information.
type Dir struct {
	Path            string
	Depth           int
	BytesSize       int64
	PrettyBytesSize string
}

// DirScanner encapsulates the state and methods for scanning directories.
type DirScanner struct {
	rootPath   string
	maxDepth   int
	minDirSize int64
	dirs       []Dir
	visited    map[string]bool
}

// NewDirScanner initializes and returns a new DirScanner.
func NewDirScanner(rootPath string, maxDepth int, minDirSizeMB int) *DirScanner {
	return &DirScanner{
		rootPath:   rootPath,
		maxDepth:   maxDepth,
		minDirSize: int64(minDirSizeMB) * 1_000_000, // Convert MB to Bytes
		dirs:       []Dir{},
		visited:    make(map[string]bool),
	}
}

// ReadDirDepth recursively scans directories up to the specified depth.
func (ds *DirScanner) ReadDirDepth(dirPath string, currentDepth int) error {
	// Calculate the current depth relative to the root path
	relativePath := strings.TrimPrefix(dirPath, ds.rootPath)
	if relativePath == dirPath {
		// If TrimPrefix didn't remove anything, ensure it doesn't start with a separator
		relativePath = strings.TrimPrefix(dirPath, string(filepath.Separator))
	}
	currentDepth = len(strings.Split(strings.Trim(relativePath, string(filepath.Separator)), string(filepath.Separator)))

	// Stop recursion if current depth exceeds the specified max depth
	if currentDepth > ds.maxDepth {
		return nil
	}

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("error reading directory '%s': %v", dirPath, err)
	}

	var dirSize int64

	for _, entry := range entries {
		entryPath := filepath.Join(dirPath, entry.Name())

		if entry.IsDir() {
			// Skip if already visited
			if ds.visited[entryPath] {
				continue
			}
			ds.visited[entryPath] = true

			// Recursively calculate the size of the subdirectory
			subDirSize, err := ds.DirSizeBytes(entryPath)
			if err != nil {
				return err
			}

			// Include the subdirectory based on minDirSize
			if subDirSize >= ds.minDirSize || ds.minDirSize == 0 {
				ds.dirs = append(ds.dirs, Dir{
					Path:            entryPath,
					Depth:           currentDepth + 1,
					BytesSize:       subDirSize,
					PrettyBytesSize: PrettyBytes(subDirSize),
				})
			}

			// Accumulate subdirectory size to the current directory's size
			dirSize += subDirSize

			// Recurse into the subdirectory if within maxDepth
			if currentDepth < ds.maxDepth {
				if err := ds.ReadDirDepth(entryPath, currentDepth+1); err != nil {
					return err
				}
			}
		} else {
			// Add file size to the current directory's size
			info, err := entry.Info()
			if err != nil {
				return fmt.Errorf("error getting info for file '%s': %v", entryPath, err)
			}
			dirSize += info.Size()
		}
	}

	// Include the current directory based on minDirSize
	if (dirSize >= ds.minDirSize || ds.minDirSize == 0) && currentDepth <= ds.maxDepth {
		ds.dirs = append(ds.dirs, Dir{
			Path:            dirPath,
			Depth:           currentDepth,
			BytesSize:       dirSize,
			PrettyBytesSize: PrettyBytes(dirSize),
		})
	}

	return nil
}

// DirSizeBytes calculates the total size of files in a directory recursively.
func (ds *DirScanner) DirSizeBytes(dirPath string) (int64, error) {
	var size int64
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// If the path is inaccessible, skip it
			if os.IsPermission(err) {
				return nil
			}
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	if err != nil {
		return 0, fmt.Errorf("error walking the path '%s': %v", dirPath, err)
	}
	return size, nil
}

// PrintLargestDirsFound displays the largest directories in a formatted table.
func (ds *DirScanner) PrintLargestDirsFound(maxResults int) {
	if len(ds.dirs) == 0 {
		fmt.Println("No directories found.")
		return
	}

	// Sort directories by size in descending order
	sort.Slice(ds.dirs, func(i, j int) bool {
		return ds.dirs[i].BytesSize > ds.dirs[j].BytesSize
	})

	// Limit the number of results to maxResults
	if len(ds.dirs) > maxResults {
		ds.dirs = ds.dirs[:maxResults]
	}

	// Initialize table with "DarkSimple" style
	t := utils.Table("DarkSimple", "largestDirsCmd")
	t.AppendHeader(table.Row{"Directory Path", "Size (MB)"})
	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, Align: text.AlignLeft, WidthMax: 60},
		{Number: 2, Align: text.AlignRight},
	})

	// Populate the table with directory data
	for _, dir := range ds.dirs {
		sizeMB := float64(dir.BytesSize) / (1024 * 1024)
		t.AppendRow(table.Row{dir.Path, fmt.Sprintf("%.2f", sizeMB)}, table.RowConfig{
			AutoMerge: true,
		})
	}

	// Render the table output
	fmt.Println()
	t.Render()
	fmt.Println()
}

// PrettyBytes formats bytes as a human-readable string.
func PrettyBytes(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit && exp < len("kMGTPE"); n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}

// LargestDirsCmd represents the largestdirs command.
var LargestDirsCmd = &cobra.Command{
	Use:   "largestdirs",
	Short: "Lists the largest directories in a specified directory, sorted by size.",
	Long: `Recursively searches for directories in the specified directory (or current directory by default) 
and lists them in descending order by size, including their absolute paths.`,
	Run: func(cmd *cobra.Command, args []string) {
		path := viper.GetString("path")
		depth := viper.GetInt("depth")
		minDirSize := viper.GetInt("mindirsize")

		// Initialize DirScanner
		scanner := NewDirScanner(path, depth, minDirSize)

		// Start scanning from the root path
		if err := scanner.ReadDirDepth(path, 0); err != nil {
			fmt.Fprintf(os.Stderr, "Error scanning directories: %v\n", err)
			os.Exit(1)
		}

		// Print the results
		scanner.PrintLargestDirsFound(10) // Display top 10 largest directories
	},
}

func init() {
	RootCmd.AddCommand(LargestDirsCmd)

	// Define flags with default values
	LargestDirsCmd.PersistentFlags().IntP("depth", "d", 1, "Depth of directory tree to display")
	LargestDirsCmd.PersistentFlags().IntP("mindirsize", "s", 0, "Only display directories larger than this threshold in MB.")
	LargestDirsCmd.PersistentFlags().StringP("path", "p", getDefaultPath(), "Path of the directory to scan")

	// Bind flags to viper
	viper.BindPFlag("depth", LargestDirsCmd.PersistentFlags().Lookup("depth"))
	viper.BindPFlag("mindirsize", LargestDirsCmd.PersistentFlags().Lookup("mindirsize"))
	viper.BindPFlag("path", LargestDirsCmd.PersistentFlags().Lookup("path"))
}

// getDefaultPath returns the current working directory or exits on error.
func getDefaultPath() string {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting current working directory: %v\n", err)
		os.Exit(1)
	}
	return filepath.Clean(pwd)
}
