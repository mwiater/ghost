//go:build windows
// +build windows

package cmd

import (
	"fmt"
	"os/user"
)

// GetLoggedInUsers retrieves the current user on Windows.
func GetLoggedInUsers() ([]LoggedInUser, error) {
	currentUser, err := user.Current()
	if err != nil {
		return nil, fmt.Errorf("retrieving logged-in users is not fully supported on Windows")
	}

	loggedInUsers := []LoggedInUser{
		{
			User:     currentUser.Username,
			Terminal: "N/A", // Terminal information is typically not available on Windows.
			Host:     "localhost",
		},
	}
	return loggedInUsers, nil
}
