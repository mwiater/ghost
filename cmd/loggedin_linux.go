//go:build linux || darwin
// +build linux darwin

package cmd

import (
	"github.com/shirou/gopsutil/v4/host"
)

// GetLoggedInUsers retrieves the currently logged-in users on Unix-based systems.
func GetLoggedInUsers() ([]LoggedInUser, error) {
	users, err := host.Users()
	if err != nil {
		return nil, err
	}

	var loggedInUsers []LoggedInUser
	for _, user := range users {
		loggedInUsers = append(loggedInUsers, LoggedInUser{
			User:     user.User,
			Terminal: user.Terminal,
			Host:     user.Host,
		})
	}
	return loggedInUsers, nil
}
