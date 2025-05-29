package main

import (
	"fmt"
	"devices-cli/pkg/devices"
)

func main() {
	fmt.Println("List of connected Android & IOS devices...")
	iosDevices, iosErr := devices.ListIOSDevices()
	androidDevices, androidErr := devices.ListAndroidDevices()
	if iosErr == nil {
		devices.DisplayDevices("iOS", iosDevices)
	} else {
		fmt.Printf("Error listing iOS devices: %v\n", iosErr)
	}

	if androidErr == nil {
		devices.DisplayDevices("Android", androidDevices)
	} else {
		fmt.Printf("Error listing Android devices: %v\n", androidErr)
	}
}
