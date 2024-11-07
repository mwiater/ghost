package cmd

import (
	"bufio"
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/mwiater/ghost/utils"
	"github.com/spf13/cobra"
)

// GPUInfoCmd represents the gpuinfo command
var GPUInfoCmd = &cobra.Command{
	Use:   "gpuinfo",
	Short: "Displays GPU information, including model, memory, and driver version.",
	Long:  `Provides detailed information about the system's GPU(s), such as the model, memory capacity, driver version, and current utilization.`,
	Run: func(cmd *cobra.Command, args []string) {
		gpus, err := RunGPUInfo()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		PrintGPUInfo(gpus)
	},
}

// GPU represents details about a GPU.
type GPU struct {
	Model         string
	Memory        string
	DriverVersion string
	Utilization   string
}

// RunGPUInfo retrieves GPU information without printing.
func RunGPUInfo() ([]GPU, error) {
	return GetGPUInfo()
}

// PrintGPUInfo displays the GPU information in a formatted table.
func PrintGPUInfo(gpus []GPU) {
	t := utils.Table("DarkSimple", "gpuinfoCmd")
	t.AppendHeader(table.Row{"Model", "Memory", "Driver Version", "Utilization"})

	for _, gpu := range gpus {
		t.AppendRow(table.Row{
			gpu.Model,
			gpu.Memory,
			gpu.DriverVersion,
			gpu.Utilization,
		})
	}

	fmt.Println()
	t.Render()
	fmt.Println()
}

func init() {
	RootCmd.AddCommand(GPUInfoCmd)

	// Define flags with default values (if any)
	// For this command, no additional flags are necessary
	// However, you can add flags here if future enhancements are needed

	// Example:
	// GPUInfoCmd.PersistentFlags().BoolP("json", "j", false, "Output in JSON format")
	// viper.BindPFlag("json", GPUInfoCmd.PersistentFlags().Lookup("json"))
}

// GetGPUInfo retrieves GPU information based on the operating system.
func GetGPUInfo() ([]GPU, error) {
	if runtime.GOOS == "windows" {
		return getGPUInfoWindows()
	}
	return getGPUInfoUnix()
}

// getGPUInfoUnix retrieves GPU information on Unix-based systems (Linux, macOS).
func getGPUInfoUnix() ([]GPU, error) {
	var gpus []GPU

	// Check if 'nvidia-smi' is available
	_, err := exec.LookPath("nvidia-smi")
	if err == nil {
		// Use 'nvidia-smi' to get detailed GPU info
		cmd := exec.Command("nvidia-smi", "--query-gpu=name,memory.total,driver_version,utilization.gpu", "--format=csv,noheader,nounits")
		output, err := cmd.Output()
		if err != nil {
			return nil, fmt.Errorf("failed to execute 'nvidia-smi': %v", err)
		}

		scanner := bufio.NewScanner(strings.NewReader(string(output)))
		for scanner.Scan() {
			line := scanner.Text()
			parts := strings.Split(line, ",")
			if len(parts) < 4 {
				continue
			}
			gpu := GPU{
				Model:         strings.TrimSpace(parts[0]),
				Memory:        strings.TrimSpace(parts[1]) + " MB",
				DriverVersion: strings.TrimSpace(parts[2]),
				Utilization:   strings.TrimSpace(parts[3]) + " %",
			}
			gpus = append(gpus, gpu)
		}
		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("error reading 'nvidia-smi' output: %v", err)
		}
	} else {
		// Fallback to 'lspci' for non-NVIDIA GPUs
		cmd := exec.Command("lspci", "-mm")
		output, err := cmd.Output()
		if err != nil {
			return nil, fmt.Errorf("failed to execute 'lspci': %v", err)
		}

		scanner := bufio.NewScanner(strings.NewReader(string(output)))
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, "VGA compatible controller") || strings.Contains(line, "3D controller") {
				// Example line format:
				// "01:00.0 \"VGA compatible controller\" \"NVIDIA Corporation\" \"GP104 [GeForce GTX 1070]\" -r06\/00\/04"
				parts := strings.Split(line, "\"")
				if len(parts) >= 6 {
					model := strings.TrimSpace(parts[5])
					gpu := GPU{
						Model:         model,
						Memory:        "N/A",
						DriverVersion: "N/A",
						Utilization:   "N/A",
					}
					gpus = append(gpus, gpu)
				}
			}
		}
		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("error reading 'lspci' output: %v", err)
		}
	}

	// Additional GPU information can be fetched here if needed

	return gpus, nil
}

// getGPUInfoWindows retrieves GPU information on Windows systems.
func getGPUInfoWindows() ([]GPU, error) {
	var gpus []GPU

	// Use WMIC to get GPU details
	cmd := exec.Command("wmic", "path", "win32_VideoController", "get", "name,adapterram,driverversion")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute WMIC command: %v", err)
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	// Skip the header line
	scanner.Scan()

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// WMIC outputs data in fixed-width columns
		// Split based on multiple spaces
		parts := strings.Fields(line)
		if len(parts) < 3 {
			continue
		}

		// Extract adapter RAM and convert to MB
		adapterRAMBytes := parts[len(parts)-2]
		adapterRAM, err := parseAdapterRAM(adapterRAMBytes)
		if err != nil {
			adapterRAM = "Unknown"
		}

		driverVersion := parts[len(parts)-1]
		model := strings.Join(parts[:len(parts)-2], " ")

		gpu := GPU{
			Model:         model,
			Memory:        adapterRAM + " MB",
			DriverVersion: driverVersion,
			Utilization:   "N/A", // Utilization is not readily available via WMIC
		}
		gpus = append(gpus, gpu)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading WMIC output: %v", err)
	}

	// Attempt to get GPU utilization using PowerShell (requires administrative privileges)
	utilGpus, err := getGPUUtilizationWindows()
	if err == nil && len(utilGpus) > 0 {
		// Merge utilization data
		for i := range gpus {
			if i < len(utilGpus) {
				gpus[i].Utilization = utilGpus[i] + " %"
			}
		}
	}

	return gpus, nil
}

// parseAdapterRAM converts adapter RAM from bytes to megabytes.
func parseAdapterRAM(adapterRAM string) (string, error) {
	bytes, err := strconv.ParseInt(adapterRAM, 10, 64)
	if err != nil {
		return "", err
	}
	mb := bytes / (1024 * 1024)
	return fmt.Sprintf("%d", mb), nil
}

// getGPUUtilizationWindows attempts to retrieve GPU utilization using PowerShell.
func getGPUUtilizationWindows() ([]string, error) {
	var utilizations []string

	powershellCmd := `Get-Counter '\GPU Engine(*)\Utilization Percentage' | Select -ExpandProperty CounterSamples | Select -ExpandProperty CookedValue`
	cmd := exec.Command("powershell", "-Command", powershellCmd)
	output, err := cmd.Output()
	if err != nil {
		return utilizations, fmt.Errorf("failed to execute PowerShell command for GPU utilization: %v", err)
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		utilizations = append(utilizations, line)
	}

	if err := scanner.Err(); err != nil {
		return utilizations, fmt.Errorf("error reading PowerShell output: %v", err)
	}

	return utilizations, nil
}
