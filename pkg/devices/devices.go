package devices

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// executeCommandLogic is the actual implementation for ExecuteCommand
func executeCommandLogic(command string, args ...string) ([]string, error) {
	cmd := exec.Command(command, args...)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to execute [%s %s]: %v", command, strings.Join(args, " "), err)
	}

	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) == 1 && lines[0] == "" { // Check if the output is effectively empty
		return nil, nil // No devices or empty output
	}


	return lines, nil
}

// ExecuteCommand is a variable that can be replaced by a mock in tests.
var ExecuteCommand = executeCommandLogic

// listIOSDevicesLogic is the actual implementation for ListIOSDevices
func listIOSDevicesLogic() ([]string, error) {
	return ExecuteCommand("idevice_id", "-l")
}

// ListIOSDevices is a variable that can be replaced by a mock in tests.
var ListIOSDevices = listIOSDevicesLogic

// listAndroidDevicesLogic is the actual implementation for ListAndroidDevices
func listAndroidDevicesLogic() ([]string, error) {
	lines, err := ExecuteCommand("adb", "devices")
	if err != nil {
		// If 'adb devices' itself fails, return the error
		return nil, err
	}
	if len(lines) == 0 || (len(lines) == 1 && lines[0] == "") { 
		// No devices connected or header only with no actual devices
		return nil, nil
	}
	// In case of "List of devices attached" and nothing else
	if len(lines) == 1 && strings.HasPrefix(lines[0], "List of devices attached") {
		return nil, nil
	}


	var devices []string
	// Skip the first line if it's "List of devices attached"
	startLine := 0
	if len(lines) > 0 && strings.HasPrefix(lines[0], "List of devices attached") {
		startLine = 1
	}

	for _, line := range lines[startLine:] {
		fields := strings.Fields(line)
		// Ensure at least 2 fields and the second is "device"
		// The first field is the device ID.
		if len(fields) >= 2 && fields[1] == "device" { 
			devices = append(devices, fields[0])
		}
	}
	if len(devices) == 0 { // if after processing, no devices were found
		return nil, nil
	}
	return devices, nil
}

// ListAndroidDevices is a variable that can be replaced by a mock in tests.
var ListAndroidDevices = listAndroidDevicesLogic

// displayDevicesLogic is the actual implementation for DisplayDevices
func displayDevicesLogic(deviceType string, devices []string) {
	fmt.Printf("Connected Devices: %s\n", deviceType)
	if len(devices) > 0 {
		for _, id := range devices {
			fmt.Printf("- %s Device: %s\n", deviceType, id)
		}
	} else {
		fmt.Printf("No %s devices connected.\n", deviceType)
	}
}

// DisplayDevices is a variable that can be replaced by a mock in tests.
var DisplayDevices = displayDevicesLogic
