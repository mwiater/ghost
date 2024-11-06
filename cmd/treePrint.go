package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// TreePrintCmd represents the treeprint command
var TreePrintCmd = &cobra.Command{
	Use:   "treeprint",
	Short: "Displays a tree-like structure of files and directories.",
	Long: `Recursively displays the directory structure in a tree format. 
Each directory and file is shown with indentation to represent its level in the hierarchy.`,
	Run: func(cmd *cobra.Command, args []string) {
		dir := "."
		if len(args) > 0 {
			dir = args[0]
		}

		tree, err := RunTreePrint(dir, ignoreList)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		PrintTree(tree)
	},
}

var ignoreList []string

func init() {
	// Define the ignore flag as a comma-separated list of directory names
	TreePrintCmd.Flags().StringSliceVar(&ignoreList, "ignore", []string{}, "Comma-separated list of directories to ignore")
	RootCmd.AddCommand(TreePrintCmd)
}

// RunTreePrint recursively generates a tree structure as a formatted string.
func RunTreePrint(rootDir string, ignoreList []string) (string, error) {
	var result strings.Builder
	result.WriteString(rootDir + "\n")
	err := treePrint(rootDir, "", ignoreList, true, &result)
	if err != nil {
		return "", err
	}
	return result.String(), nil
}

// treePrint is a recursive function that appends files and directories in a tree format to the result.
func treePrint(path, prefix string, ignoreList []string, isLast bool, result *strings.Builder) error {
	// Get directory contents
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	// Filter entries to ignore specified directories
	var filteredEntries []os.DirEntry
	for _, entry := range entries {
		if entry.IsDir() && shouldIgnore(entry.Name(), ignoreList) {
			continue
		}
		filteredEntries = append(filteredEntries, entry)
	}

	// Iterate over filtered directory entries
	for i, entry := range filteredEntries {
		// Check if this entry is the last one in the directory
		isLastEntry := i == len(filteredEntries)-1

		// Append the appropriate prefix and entry name to the result
		result.WriteString(prefix)
		if isLastEntry {
			result.WriteString("└── ")
		} else {
			result.WriteString("├── ")
		}
		result.WriteString(entry.Name() + "\n")

		// If entry is a directory, recursively add its contents with updated prefix
		if entry.IsDir() {
			subDir := filepath.Join(path, entry.Name())
			newPrefix := prefix
			if isLastEntry {
				newPrefix += "    "
			} else {
				newPrefix += "│   "
			}
			if err := treePrint(subDir, newPrefix, ignoreList, isLastEntry, result); err != nil {
				return err
			}
		}
	}

	return nil
}

// PrintTree displays the tree structure from a formatted string.
func PrintTree(tree string) {
	fmt.Println(tree)
}

// shouldIgnore checks if a directory should be ignored based on the ignore list.
func shouldIgnore(dirName string, ignoreList []string) bool {
	for _, ignoreDir := range ignoreList {
		if dirName == ignoreDir {
			return true
		}
	}
	return false
}
