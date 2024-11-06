// Package utils provides utilities for interacting with the terminal and formatting output.
package utils

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

type Terminal struct {
	Height                  int
	Width                   int
	OutputType              string
	NumberOfSupportedColors int
	TERM                    string
	SHELL                   string
	COLORTERM               string
}

// ErrorLevel type for defining constants for the error levels
type ErrorLevel int

// Define constants for the error levels
const (
	Alert ErrorLevel = iota
	Critical
	Error
	Warn
	Notice
	Info
	Debug
)

// colorMap maps error levels to their respective ANSI color codes
var colorMap = map[ErrorLevel]string{
	Alert:    "\033[38;5;201m", // Magenta
	Critical: "\033[38;5;214m", // Orange
	Error:    "\033[38;5;196m", // Light Red
	Warn:     "\033[38;5;226m", // Yellow
	Notice:   "\033[38;5;117m", // Light Blue
	Info:     "\033[38;5;250m", // Gray
	Debug:    "\033[38;5;120m", // Light Green
}

// ClearTerminal clears the terminal screen based on the operating system.
func ClearTerminal() error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	default:
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

// TerminalColor prints the given string to the terminal in the color corresponding to the error level
func TerminalColor(message string, level ErrorLevel) {
	colorCode, ok := colorMap[level]
	if !ok {
		fmt.Println(message)
		return
	}
	fmt.Printf("%s%s\033[0m\n", colorCode, message)
}
