package cmd

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/mwiater/ghost/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// LoginsCmd represents the logins command
var LoginsCmd = &cobra.Command{
	Use:   "logins",
	Short: "Displays recent login attempts and current user sessions.",
	Long: `Retrieves and displays a list of recent login attempts, including successful and failed attempts,
as well as currently logged-in users with their login times and IP addresses (if available).`,
	Run: func(cmd *cobra.Command, args []string) {
		// Retrieve the 'count' flag value
		count := viper.GetInt("count")

		logins, err := RunLogins(count)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		PrintLogins(logins)
	},
}

// LoginEntry holds details about a login attempt or session.
type LoginEntry struct {
	User      string
	Terminal  string
	Host      string
	Time      string
	Status    string
	IPAddress string
}

// RunLogins retrieves the login entries without printing.
// It accepts 'count' to limit the number of entries.
func RunLogins(count int) ([]LoginEntry, error) {
	return GetLogins(count)
}

// PrintLogins displays the login entries in a formatted table.
func PrintLogins(logins []LoginEntry) {
	t := utils.Table("DarkSimple", "loginsCmd")
	t.AppendHeader(table.Row{"User", "Terminal", "Host", "Time", "Status", "IP Address"})

	for _, entry := range logins {
		t.AppendRow(table.Row{
			entry.User,
			entry.Terminal,
			entry.Host,
			entry.Time,
			entry.Status,
			entry.IPAddress,
		})
	}

	fmt.Println()
	t.Render()
	fmt.Println()
}

func init() {
	RootCmd.AddCommand(LoginsCmd)

	// Define flags with default values
	LoginsCmd.PersistentFlags().IntP("count", "c", 10, "Number of login entries to display")

	// Bind flags to viper
	viper.BindPFlag("count", LoginsCmd.PersistentFlags().Lookup("count"))
}

// GetLogins retrieves login information based on the operating system.
// It accepts 'count' to limit the number of entries returned.
func GetLogins(count int) ([]LoginEntry, error) {
	if runtime.GOOS == "windows" {
		return getLoginsWindows(count)
	}
	return getLoginsUnix(count)
}

// getLoginsUnix retrieves login information on Unix-based systems (Linux, macOS).
func getLoginsUnix(count int) ([]LoginEntry, error) {
	var entries []LoginEntry

	// Current logged-in users using 'who' command
	whoCmd := exec.Command("who")
	whoOutput, err := whoCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute 'who' command: %v", err)
	}

	whoLines := strings.Split(string(whoOutput), "\n")
	for _, line := range whoLines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 5 {
			continue
		}
		user := parts[0]
		terminal := parts[1]
		time := strings.Join(parts[2:5], " ")
		host := "-"
		ip := "-"
		if len(parts) >= 6 {
			host = parts[5]
			if strings.HasPrefix(host, "(") && strings.HasSuffix(host, ")") {
				ip = host[1 : len(host)-1]
				host = "-"
			}
		}
		entry := LoginEntry{
			User:      user,
			Terminal:  terminal,
			Host:      host,
			Time:      time,
			Status:    "Active",
			IPAddress: ip,
		}
		entries = append(entries, entry)

		if len(entries) >= count {
			break
		}
	}

	// Check if we need to fetch recent login attempts
	if len(entries) < count {
		remaining := count - len(entries)
		// Recent login attempts using 'last' command with '-n' to limit entries
		lastCmd := exec.Command("last", "-n", fmt.Sprintf("%d", remaining), "-a")
		lastOutput, err := lastCmd.Output()
		if err != nil {
			// 'last' might not be available on all Unix systems
			// Return current sessions only
			return entries, nil
		}

		lastLines := strings.Split(string(lastOutput), "\n")
		for _, line := range lastLines {
			if strings.TrimSpace(line) == "" || strings.HasPrefix(line, "wtmp") {
				continue
			}
			parts := strings.Fields(line)
			if len(parts) < 10 {
				continue
			}
			user := parts[0]
			terminal := parts[1]
			// Date and time can vary, so we'll join a range
			time := strings.Join(parts[2:8], " ")
			status := parts[8]
			host := parts[9]
			ip := "-"
			if strings.HasPrefix(host, "(") && strings.HasSuffix(host, ")") {
				ip = host[1 : len(host)-1]
				host = "-"
			}
			entry := LoginEntry{
				User:      user,
				Terminal:  terminal,
				Host:      host,
				Time:      time,
				Status:    status,
				IPAddress: ip,
			}
			entries = append(entries, entry)

			if len(entries) >= count {
				break
			}
		}
	}

	return entries, nil
}

// getLoginsWindows retrieves login information on Windows systems.
func getLoginsWindows(count int) ([]LoginEntry, error) {
	var entries []LoginEntry

	// Get currently logged-in users using 'query user' command
	queryCmd := exec.Command("query", "user")
	queryOutput, err := queryCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute 'query user' command: %v", err)
	}

	queryLines := strings.Split(string(queryOutput), "\n")
	for _, line := range queryLines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "USERNAME") {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 5 {
			continue
		}
		user := parts[0]
		terminal := parts[1]
		time := parts[3] + " " + parts[4]
		status := parts[2]
		// Windows 'query user' doesn't provide IP addresses by default
		entry := LoginEntry{
			User:      user,
			Terminal:  terminal,
			Host:      "-",
			Time:      time,
			Status:    status,
			IPAddress: "-",
		}
		entries = append(entries, entry)

		if len(entries) >= count {
			break
		}
	}

	// Check if we need to fetch recent login attempts
	if len(entries) < count {
		remaining := count - len(entries)
		// Retrieve recent login attempts from the Security event log
		// This requires PowerShell commands
		powershellCmd := fmt.Sprintf(`Get-EventLog -LogName Security -InstanceId 4624,4625 -Newest %d | Select-Object TimeGenerated, @{Name="User";Expression={$_.ReplacementStrings[5]}}, @{Name="IP";Expression={$_.ReplacementStrings[18]}}`, remaining)
		cmd := exec.Command("powershell", "-Command", powershellCmd)
		psOutput, err := cmd.Output()
		if err != nil {
			// If PowerShell command fails, skip recent logins
			return entries, nil
		}

		psLines := strings.Split(string(psOutput), "\n")
		for _, line := range psLines {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "TimeGenerated") {
				continue
			}
			parts := strings.Fields(line)
			if len(parts) < 3 {
				continue
			}
			time := parts[0] + " " + parts[1]
			user := parts[2]
			ip := "-"
			if len(parts) >= 4 {
				ip = parts[3]
			}
			entry := LoginEntry{
				User:      user,
				Terminal:  "-",
				Host:      "-",
				Time:      time,
				Status:    "Recent",
				IPAddress: ip,
			}
			entries = append(entries, entry)

			if len(entries) >= count {
				break
			}
		}
	}

	return entries, nil
}
