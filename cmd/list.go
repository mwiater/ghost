package cmd

import (
	"github.com/mwiater/ghost/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var packages bool
var packageName string

// ListCmd represents the `list` command, which lists all the functions available
// within the 'utils' package. It can list functions for all packages, or functions
// specific to a particular package, along with their signatures and comments.
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all the functions available within the 'utils' package.",
	Long: `Lists functions from the 'utils' package, along with comments,
function argument types, and return types. Optionally, it can list only
the package names or functions for a specific package.`,
	Run: func(cmd *cobra.Command, args []string) {
		utils.ClearTerminal()
		packages := viper.GetBool("packages")
		packageName := viper.GetString("packageName")

		utils.ListPackageFunctions("./utils/", packages, packageName)
	},
}

// init initializes the `list` command and adds it to the RootCmd.
// The command has two flags: 'packages', which lists only package names, and
// 'packageName', which lists functions for a specific package.
func init() {
	// Define flags for the list command
	ListCmd.Flags().BoolVarP(&packages, "packages", "", false, "List package names only")
	ListCmd.Flags().StringVarP(&packageName, "packageName", "p", "", "List functions for a specific package")

	// Bind the flags to Viper for configuration management
	viper.BindPFlag("packages", ListCmd.Flags().Lookup("packages"))
	viper.BindPFlag("packageName", ListCmd.Flags().Lookup("packageName"))
}
