package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// executeCommand runs a command and returns the output lines
func executeCommand(command string, args ...string) ([]string, error) {
	cmd := exec.Command(command, args...)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to execute [%s %s]: %v", command, strings.Join(args, " "), err)
	}

	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) == 0 {
		return nil, nil // No devices connected
	}

	return lines, nil
}

// listIOSDevices connects to usbmuxd and lists connected iOS devices
func listIOSDevices() ([]string, error) {
	return executeCommand("idevice_id", "-l")
}

// listAndroidDevices lists connected Android devices using adb
func listAndroidDevices() ([]string, error) {
	lines, err := executeCommand("adb", "devices")
	if err != nil || len(lines) <= 1 {
		return nil, err // No devices connected
	}

	var devices []string
	for _, line := range lines[1:] {
		fields := strings.Fields(line)
		if len(fields) == 2 && fields[1] == "device" {
			devices = append(devices, fields[0])
		}
	}
	return devices, nil
}

// displayDevices prints the list of connected devices
func displayDevices(deviceType string, devices []string) {
	fmt.Printf("Connected Devices: %s\n", deviceType)
	if len(devices) > 0 {
		for _, id := range devices {
			fmt.Printf("- %s Device: %s\n", deviceType, id)
		}
	} else {
		fmt.Printf("No %s devices connected.\n", deviceType)
	}
}

func main() {
	fmt.Println("List of connected Android & IOS devices...")
	iosDevices, iosErr := listIOSDevices()
	androidDevices, androidErr := listAndroidDevices()
	if iosErr == nil {
		displayDevices("iOS", iosDevices)
	} else {
		fmt.Printf("Error listing iOS devices: %v\n", iosErr)
	}

	if androidErr == nil {
		displayDevices("Android", androidDevices)
	} else {
		fmt.Printf("Error listing Android devices: %v\n", androidErr)
	}
}
