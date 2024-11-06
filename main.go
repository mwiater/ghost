package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mwiater/ghost/cmd"
	"github.com/mwiater/ghost/utils"
)

// checkIsGoRun determines whether the executable is being run via `go run`.
// It checks the executable's path against the system's temporary directory.
// Returns true if the executable is in the temp directory (likely indicating a `go run` execution),
// otherwise, returns false for a compiled binary.
func checkIsGoRun() bool {
	executablePath, err := os.Executable()
	if err != nil {
		fmt.Println("Error checking executable path:", err)
		return false
	}

	tempDir := os.TempDir()

	executablePath, err = filepath.Abs(executablePath)
	if err != nil {
		fmt.Println("Error converting executable path to absolute:", err)
		return false
	}
	tempDir, err = filepath.Abs(tempDir)
	if err != nil {
		fmt.Println("Error converting temp dir to absolute:", err)
		return false
	}

	relPath, err := filepath.Rel(tempDir, executablePath)
	if err != nil {
		fmt.Println("Error checking relative path:", err)
		return false
	}

	return !strings.HasPrefix(relPath, "..") && relPath != "."
}

// main is the entry point of the application. It clears the terminal, checks
// if the executable is being run via `go run` or as a compiled binary, and then
// conditionally registers commands. It also executes the root command defined
// in the cmd package.
func main() {
	utils.ClearTerminal()
	isGoRun := checkIsGoRun()
	if isGoRun {
		cmd.RootCmd.AddCommand(cmd.ListCmd)
	}

	cmd.Execute()
}
