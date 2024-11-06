package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// IsGoRun indicates whether the application is being run via `go run`.
// This variable is set in the main.go file.
var IsGoRun bool

// RootCmd represents the base command when called without any subcommands.
// It serves as the entry point for all utilities and tools available within the application.
var RootCmd = &cobra.Command{
	Use:   "ghost",
	Short: "Network diagnostics and system info toolkit.",
	Long:  `A versatile toolkit for network diagnostics and system information gathering, offering developers a suite of commands to scan networks, retrieve system details, and perform IP and port analyses.`,
}

// Execute adds all child commands to the root command and sets the flags appropriately.
// This function is called by main.main() and only needs to be called once for RootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// init initializes the RootCmd and sets up flags for the base command.
// Persistent flags are global for the application, while local flags apply to specific actions.
func init() {
	// Example of defining a persistent flag:
	// RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.golangutils.yaml)")

	// Example of defining a local flag:
	// RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
