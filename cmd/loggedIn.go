package cmd

import (
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/mwiater/ghost/utils"
	"github.com/spf13/cobra"
)

// LoggedInCmd represents the loggedin command
var LoggedInCmd = &cobra.Command{
	Use:   "loggedin",
	Short: "Displays currently logged-in users.",
	Long:  `Retrieves and displays a list of currently logged-in users, including login times and IP addresses (if available).`,
	Run: func(cmd *cobra.Command, args []string) {
		loggedInUsers, err := RunLoggedIn()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		PrintLoggedInUsers(loggedInUsers)
	},
}

// LoggedInUser holds details about a logged-in user.
type LoggedInUser struct {
	User     string
	Terminal string
	Host     string
}

// RunLoggedIn retrieves the logged-in users without printing.
func RunLoggedIn() ([]LoggedInUser, error) {
	return GetLoggedInUsers()
}

// PrintLoggedInUsers displays the logged-in users in a formatted table.
func PrintLoggedInUsers(loggedInUsers []LoggedInUser) {
	t := utils.Table("DarkSimple", "loggedInCmd")
	t.AppendHeader(table.Row{"User", "Terminal", "Host"})

	for _, user := range loggedInUsers {
		t.AppendRow(table.Row{user.User, user.Terminal, user.Host})
	}

	fmt.Println()
	t.Render()
	fmt.Println()
}

func init() {
	RootCmd.AddCommand(LoggedInCmd)
}
