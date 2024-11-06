package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/mwiater/ghost/utils"
	"github.com/spf13/cobra"
)

// FindCmd represents the find command
var FindCmd = &cobra.Command{
	Use:   "find",
	Short: "Finds files with names containing a specified substring.",
	Long:  `Searches for files with names that contain a specified substring, optionally within a specified directory. Returns the list of matching files with their absolute paths.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Error: A search substring is required.")
			return
		}
		searchTerm := args[0]
		dir := "."

		if len(args) > 1 {
			dir = args[1]
		}

		matches, err := RunFind(searchTerm, dir)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		PrintFindResults(matches)
	},
}

// FindFile holds details about a found file.
type FindFile struct {
	Path string
}

// GetMatchingFiles searches for files that contain the searchTerm in their name.
func GetMatchingFiles(searchTerm, startDir string) ([]FindFile, error) {
	var matches []FindFile

	// Walk through the directory recursively
	err := filepath.WalkDir(startDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// If the file name contains the search term, add it to matches
		if !d.IsDir() && strings.Contains(d.Name(), searchTerm) {
			absPath, err := filepath.Abs(path)
			if err != nil {
				return err
			}
			matches = append(matches, FindFile{Path: absPath})
		}
		return nil
	})

	return matches, err
}

// RunFind retrieves matching files without printing.
func RunFind(searchTerm, startDir string) ([]FindFile, error) {
	return GetMatchingFiles(searchTerm, startDir)
}

// PrintFindResults displays the list of matching files in a formatted table.
func PrintFindResults(matches []FindFile) {
	if len(matches) == 0 {
		fmt.Println("No matching files found.")
		return
	}

	// Use utils.Table to create a table with "DarkSimple" style for alternate row shading
	t := utils.Table("DarkSimple", "findCmd")
	t.AppendHeader(table.Row{"Matching Files"})

	// Add each matching file's path to the table
	for _, file := range matches {
		t.AppendRow(table.Row{file.Path})
	}

	// Render the table
	fmt.Println()
	t.Render()
	fmt.Println()
}

func init() {
	RootCmd.AddCommand(FindCmd)
}
