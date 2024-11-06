package cmd

import (
	"os"
	"sort"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/mwiater/ghost/utils"
	"github.com/spf13/cobra"
)

// EnvVarsCmd represents the envvars command
var EnvVarsCmd = &cobra.Command{
	Use:   "envvars",
	Short: "Displays all environment variables.",
	Long:  `Retrieves and displays all environment variables in a consistent, readable format, providing variable names and their values.`,
	Run: func(cmd *cobra.Command, args []string) {
		envVars := RunEnvVars()
		PrintEnvVars(envVars)
	},
}

// GetEnvVars retrieves all environment variables as a sorted list of key-value pairs.
func GetEnvVars() map[string]string {
	envVars := make(map[string]string)
	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		key := strings.TrimSpace(parts[0])
		value := ""
		if len(parts) > 1 {
			value = strings.TrimSpace(parts[1])
		}
		if key != "" && value != "" { // Only include non-empty variables
			envVars[key] = value
		}
	}
	return envVars
}

// RunEnvVars retrieves the environment variables without printing.
func RunEnvVars() map[string]string {
	return GetEnvVars()
}

// PrintEnvVars displays the environment variables in a formatted table.
func PrintEnvVars(envVars map[string]string) {
	// Sort environment variables by name for consistent ordering
	keys := make([]string, 0, len(envVars))
	for k := range envVars {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Use utils.Table to create a table with "DarkSimple" style for alternate row shading
	t := utils.Table("DarkSimple", "envVarsCmd")
	t.AppendHeader(table.Row{"Variable", "Value"})

	// Add each environment variable to the table
	for _, k := range keys {
		t.AppendRow(table.Row{k, envVars[k]})
	}

	// Set column width for Value to wrap long entries
	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 2, WidthMax: 40}, // Set max width for Value column
	})

	// Render the table without extra padding
	t.Style().Box.PaddingLeft = ""
	t.Style().Box.PaddingRight = ""
	t.Render()
}

func init() {
	RootCmd.AddCommand(EnvVarsCmd)
}
