package cmd

import (
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/mwiater/ghost/utils" // using the specified module name
	"github.com/spf13/cobra"
)

// Service represents a single service with name, status, and memory usage.
type Service struct {
	Name        string
	Status      string
	MemoryUsage string
}

// ServicesCmd represents the services command
var ServicesCmd = &cobra.Command{
	Use:   "services",
	Short: "Lists running services with their status and memory usage.",
	Long:  `Retrieves and displays a list of all running services, showing the service name, current status, and memory usage.`,
	Run: func(cmd *cobra.Command, args []string) {
		services, err := RunServices()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		PrintServices(services)
	},
}

// RunServices retrieves the list of running services without printing.
func RunServices() ([]Service, error) {
	return GetServices()
}

// PrintServices displays the list of running services in a formatted table.
func PrintServices(services []Service) {
	// Use utils.Table to create a table with "DarkSimple" style for alternate row shading
	t := utils.Table("DarkSimple", "servicesCmd")
	t.AppendHeader(table.Row{"Service Name", "Status", "Memory Usage"})

	// Add each service's details to the table
	for _, svc := range services {
		t.AppendRow(table.Row{svc.Name, svc.Status, svc.MemoryUsage})
	}

	// Render the table
	fmt.Println()
	t.Render()
	fmt.Println()
}

func init() {
	RootCmd.AddCommand(ServicesCmd)
}
