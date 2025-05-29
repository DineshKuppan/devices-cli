package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
	"errors" // For creating mock errors
	"reflect" // For DeepEqual in tests

	"devices-cli/pkg/devices"
)

// captureOutput helper function
func captureOutput(f func()) string {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	if err := w.Close(); err != nil {
		// Consider logging this error if your test framework supports it
		// For now, we'll proceed, but this could mask issues.
	}
	os.Stdout = oldStdout
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r)
	if err != nil {
        // Similar to w.Close(), log if possible.
	}
	return buf.String()
}

// TestMain_SuccessPaths tests the main function for successful device listing
func TestMain_SuccessPaths(t *testing.T) {
	// Save original functions
	originalListIOSDevices := devices.ListIOSDevices
	originalListAndroidDevices := devices.ListAndroidDevices
	originalDisplayDevices := devices.DisplayDevices

	// Defer restoration
	defer func() {
		devices.ListIOSDevices = originalListIOSDevices
		devices.ListAndroidDevices = originalListAndroidDevices
		devices.DisplayDevices = originalDisplayDevices
	}()

	// Mock implementations
	devices.ListIOSDevices = func() ([]string, error) {
		return []string{"iosDevice1"}, nil
	}
	devices.ListAndroidDevices = func() ([]string, error) {
		return []string{"androidDevice1"}, nil
	}
	var displayedIOSType, displayedAndroidType string
	var displayedIOSDevices, displayedAndroidDevices []string

	devices.DisplayDevices = func(deviceType string, devList []string) {
		if deviceType == "iOS" {
			displayedIOSType = deviceType
			displayedIOSDevices = devList
		} else if deviceType == "Android" {
			displayedAndroidType = deviceType
			displayedAndroidDevices = devList
		}
	}

	output := captureOutput(main)

	if !strings.Contains(output, "List of connected Android & IOS devices...") {
		t.Errorf("Expected introductory text not found in output: %s", output)
	}

	if displayedIOSType != "iOS" {
		t.Errorf("Expected DisplayDevices to be called for iOS, but it wasn't")
	}
	if !reflect.DeepEqual(displayedIOSDevices, []string{"iosDevice1"}) {
		t.Errorf("Expected iOS devices %v, got %v", []string{"iosDevice1"}, displayedIOSDevices)
	}

	if displayedAndroidType != "Android" {
		t.Errorf("Expected DisplayDevices to be called for Android, but it wasn't")
	}
	if !reflect.DeepEqual(displayedAndroidDevices, []string{"androidDevice1"}) {
		t.Errorf("Expected Android devices %v, got %v", []string{"androidDevice1"}, displayedAndroidDevices)
	}
	
	// Check for positive output messages (optional, as DisplayDevices is mocked)
	// We are checking if DisplayDevices was called correctly above.
	// If DisplayDevices was not mocked to capture, we'd check os.Stdout here.
}

// TestMain_IOSError tests the main function when ListIOSDevices returns an error
func TestMain_IOSError(t *testing.T) {
	originalListIOSDevices := devices.ListIOSDevices
	originalListAndroidDevices := devices.ListAndroidDevices
	originalDisplayDevices := devices.DisplayDevices
	defer func() {
		devices.ListIOSDevices = originalListIOSDevices
		devices.ListAndroidDevices = originalListAndroidDevices
		devices.DisplayDevices = originalDisplayDevices
	}()

	mockIOSError := errors.New("mock iOS error")
	devices.ListIOSDevices = func() ([]string, error) {
		return nil, mockIOSError
	}
	devices.ListAndroidDevices = func() ([]string, error) {
		return []string{"androidDevice1"}, nil
	}
	
	var displayedAndroid bool
	devices.DisplayDevices = func(deviceType string, devList []string) {
		if deviceType == "Android" {
			displayedAndroid = true
			if !reflect.DeepEqual(devList, []string{"androidDevice1"}) {
				t.Errorf("Expected Android devices %v, got %v", []string{"androidDevice1"}, devList)
			}
		}
		if deviceType == "iOS" {
			t.Errorf("DisplayDevices should not be called for iOS when there's an error")
		}
	}

	output := captureOutput(main)

	if !strings.Contains(output, "Error listing iOS devices: "+mockIOSError.Error()) {
		t.Errorf("Expected iOS error message not found in output: %s", output)
	}
	if !displayedAndroid {
		t.Errorf("Expected Android devices to be displayed even if iOS fails")
	}
}

// TestMain_AndroidError tests the main function when ListAndroidDevices returns an error
func TestMain_AndroidError(t *testing.T) {
	originalListIOSDevices := devices.ListIOSDevices
	originalListAndroidDevices := devices.ListAndroidDevices
	originalDisplayDevices := devices.DisplayDevices
	defer func() {
		devices.ListIOSDevices = originalListIOSDevices
		devices.ListAndroidDevices = originalListAndroidDevices
		devices.DisplayDevices = originalDisplayDevices
	}()

	mockAndroidError := errors.New("mock Android error")
	devices.ListIOSDevices = func() ([]string, error) {
		return []string{"iosDevice1"}, nil
	}
	devices.ListAndroidDevices = func() ([]string, error) {
		return nil, mockAndroidError
	}

	var displayedIOS bool
	devices.DisplayDevices = func(deviceType string, devList []string) {
		if deviceType == "iOS" {
			displayedIOS = true
			if !reflect.DeepEqual(devList, []string{"iosDevice1"}) {
				t.Errorf("Expected iOS devices %v, got %v", []string{"iosDevice1"}, devList)
			}
		}
		if deviceType == "Android" {
			t.Errorf("DisplayDevices should not be called for Android when there's an error")
		}
	}

	output := captureOutput(main)

	if !strings.Contains(output, "Error listing Android devices: "+mockAndroidError.Error()) {
		t.Errorf("Expected Android error message not found in output: %s", output)
	}
	if !displayedIOS {
		t.Errorf("Expected iOS devices to be displayed even if Android fails")
	}
}

// TestMain_BothError tests the main function when both listing functions return errors
func TestMain_BothError(t *testing.T) {
	originalListIOSDevices := devices.ListIOSDevices
	originalListAndroidDevices := devices.ListAndroidDevices
	originalDisplayDevices := devices.DisplayDevices
	defer func() {
		devices.ListIOSDevices = originalListIOSDevices
		devices.ListAndroidDevices = originalListAndroidDevices
		devices.DisplayDevices = originalDisplayDevices
	}()

	mockIOSError := errors.New("mock iOS error")
	mockAndroidError := errors.New("mock Android error")

	devices.ListIOSDevices = func() ([]string, error) {
		return nil, mockIOSError
	}
	devices.ListAndroidDevices = func() ([]string, error) {
		return nil, mockAndroidError
	}
	devices.DisplayDevices = func(deviceType string, devList []string) {
		t.Errorf("DisplayDevices should not be called when both list functions error out")
	}

	output := captureOutput(main)

	if !strings.Contains(output, "Error listing iOS devices: "+mockIOSError.Error()) {
		t.Errorf("Expected iOS error message not found in output: %s", output)
	}
	if !strings.Contains(output, "Error listing Android devices: "+mockAndroidError.Error()) {
		t.Errorf("Expected Android error message not found in output: %s", output)
	}
}
