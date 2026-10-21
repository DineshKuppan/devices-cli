package devices

import (
	"errors"
	"fmt"
	"os/exec"
	"reflect"
	"strings"
	"testing"
)

// TestExecuteCommandLogic_Success tests a successful command execution
func TestExecuteCommandLogic_Success(t *testing.T) {
	if _, err := exec.LookPath("echo"); err != nil {
		t.Skip("echo command not found, skipping TestExecuteCommandLogic_Success")
	}
	lines, err := executeCommandLogic("echo", "hello", "world")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	expectedLines := []string{"hello world"}
	if !reflect.DeepEqual(lines, expectedLines) {
		t.Errorf("expected lines %v, got %v", expectedLines, lines)
	}
}

// TestExecuteCommandLogic_CommandFailure tests a failing command
func TestExecuteCommandLogic_CommandFailure(t *testing.T) {
	nonExistentCommand := "a_command_that_should_not_exist_anywhere"
	lines, err := executeCommandLogic(nonExistentCommand)
	if err == nil {
		t.Fatalf("expected an error, got nil")
	}
	if lines != nil {
		t.Errorf("expected no lines, got %v", lines)
	}
	expectedErrorPart := fmt.Sprintf("failed to execute [%s ]", nonExistentCommand)
	if !strings.Contains(err.Error(), expectedErrorPart) {
		t.Errorf("expected error message to contain '%s', got '%s'", expectedErrorPart, err.Error())
	}
}

// TestExecuteCommandLogic_NoOutput tests a command that produces no stdout output
func TestExecuteCommandLogic_NoOutput(t *testing.T) {
	if _, err := exec.LookPath("echo"); err != nil {
		t.Skip("echo command not found, skipping TestExecuteCommandLogic_NoOutput")
	}
	lines, err := executeCommandLogic("echo", "-n")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if lines != nil {
		t.Errorf("expected nil lines for no output, got %v", lines)
	}
}

// Mock function for ExecuteCommand
type mockExecuteCmd func(command string, args ...string) ([]string, error)

// Helper to temporarily set a mock for ExecuteCommand and restore it
func withMockExecuteCommand(mockFunc mockExecuteCmd, f func()) {
	originalExecuteCommand := ExecuteCommand
	ExecuteCommand = mockFunc
	defer func() {
		ExecuteCommand = originalExecuteCommand
	}()
	f()
}

// Tests for ListIOSDevices
func TestListIOSDevices_Success(t *testing.T) {
	expectedDevices := []string{"deviceid1", "deviceid2"}
	withMockExecuteCommand(func(command string, args ...string) ([]string, error) {
		if command != "idevice_id" || args[0] != "-l" {
			t.Fatalf("ExecuteCommand called with unexpected command/args: %s %v", command, args)
		}
		return expectedDevices, nil
	}, func() {
		devices, err := ListIOSDevices()
		if err != nil {
			t.Fatalf("ListIOSDevices() returned error: %v", err)
		}
		if !reflect.DeepEqual(devices, expectedDevices) {
			t.Errorf("Expected devices %v, got %v", expectedDevices, devices)
		}
	})
}

func TestListIOSDevices_NoDevices(t *testing.T) {
	withMockExecuteCommand(func(command string, args ...string) ([]string, error) {
		if command != "idevice_id" || args[0] != "-l" {
			t.Fatalf("ExecuteCommand called with unexpected command/args: %s %v", command, args)
		}
		return nil, nil // Simulate no devices found
	}, func() {
		devices, err := ListIOSDevices()
		if err != nil {
			t.Fatalf("ListIOSDevices() returned error: %v", err)
		}
		if devices != nil {
			t.Errorf("Expected no devices (nil), got %v", devices)
		}
	})
}

func TestListIOSDevices_Error(t *testing.T) {
	expectedError := errors.New("idevice_id command failed")
	withMockExecuteCommand(func(command string, args ...string) ([]string, error) {
		if command != "idevice_id" || args[0] != "-l" {
			t.Fatalf("ExecuteCommand called with unexpected command/args: %s %v", command, args)
		}
		return nil, expectedError
	}, func() {
		_, err := ListIOSDevices()
		if err == nil {
			t.Fatal("ListIOSDevices() expected an error, got nil")
		}
		if !strings.Contains(err.Error(), expectedError.Error()) {
			t.Errorf("Expected error containing '%v', got '%v'", expectedError, err)
		}
	})
}

// Tests for ListAndroidDevices
func TestListAndroidDevices_Success(t *testing.T) {
	mockOutput := []string{
		"List of devices attached",
		"emulator-5554\tdevice",
		"R58M22XYZAB\tdevice",
	}
	expectedDevices := []string{"emulator-5554", "R58M22XYZAB"}
	withMockExecuteCommand(func(command string, args ...string) ([]string, error) {
		if command != "adb" || args[0] != "devices" {
			t.Fatalf("ExecuteCommand called with unexpected command/args: %s %v", command, args)
		}
		return mockOutput, nil
	}, func() {
		devices, err := ListAndroidDevices()
		if err != nil {
			t.Fatalf("ListAndroidDevices() returned error: %v", err)
		}
		if !reflect.DeepEqual(devices, expectedDevices) {
			t.Errorf("Expected devices %v, got %v", expectedDevices, devices)
		}
	})
}

func TestListAndroidDevices_NoDevices_HeaderOnly(t *testing.T) {
	mockOutput := []string{"List of devices attached"} 
	withMockExecuteCommand(func(command string, args ...string) ([]string, error) {
		if command != "adb" || args[0] != "devices" {
			t.Fatalf("ExecuteCommand called with unexpected command/args: %s %v", command, args)
		}
		return mockOutput, nil
	}, func() {
		devices, err := ListAndroidDevices()
		if err != nil {
			t.Fatalf("ListAndroidDevices() returned error: %v", err)
		}
		if devices != nil {
			t.Errorf("Expected no devices (nil), got %v", devices)
		}
	})
}

func TestListAndroidDevices_NoDevices_EmptyOutputFromExecute(t *testing.T) {
	withMockExecuteCommand(func(command string, args ...string) ([]string, error) {
		if command != "adb" || args[0] != "devices" {
			t.Fatalf("ExecuteCommand called with unexpected command/args: %s %v", command, args)
		}
		return nil, nil
	}, func() {
		devices, err := ListAndroidDevices()
		if err != nil {
			t.Fatalf("ListAndroidDevices() returned error: %v", err)
		}
		if devices != nil {
			t.Errorf("Expected no devices (nil), got %v", devices)
		}
	})
}

func TestListAndroidDevices_Error(t *testing.T) {
	expectedError := errors.New("adb command failed")
	withMockExecuteCommand(func(command string, args ...string) ([]string, error) {
		if command != "adb" || args[0] != "devices" {
			t.Fatalf("ExecuteCommand called with unexpected command/args: %s %v", command, args)
		}
		return nil, expectedError
	}, func() {
		_, err := ListAndroidDevices()
		if err == nil {
			t.Fatal("ListAndroidDevices() expected an error, got nil")
		}
		if !errors.Is(err, expectedError) && err.Error() != expectedError.Error() {
             // Comparing the error string directly as ListAndroidDevices returns the error from ExecuteCommand as is.
			t.Errorf("Expected error '%v', got '%v'", expectedError, err)
		}
	})
}

func TestListAndroidDevices_UnauthorizedDevice(t *testing.T) {
	mockOutput := []string{
		"List of devices attached",
		"R58M22XYZAB\tunauthorized",
		"emulator-5554\tdevice",
	}
	expectedDevices := []string{"emulator-5554"}
	withMockExecuteCommand(func(command string, args ...string) ([]string, error) {
		if command != "adb" || args[0] != "devices" {
			t.Fatalf("ExecuteCommand called with unexpected command/args: %s %v", command, args)
		}
		return mockOutput, nil
	}, func() {
		devices, err := ListAndroidDevices()
		if err != nil {
			t.Fatalf("ListAndroidDevices() returned error: %v", err)
		}
		if !reflect.DeepEqual(devices, expectedDevices) {
			t.Errorf("Expected devices %v, got %v", expectedDevices, devices)
		}
	})
}

func TestListAndroidDevices_EmulatorDevice(t *testing.T) {
	mockOutput := []string{
		"List of devices attached",
		"emulator-5554\tdevice",
	}
	expectedDevices := []string{"emulator-5554"}
	withMockExecuteCommand(func(command string, args ...string) ([]string, error) {
		if command != "adb" || args[0] != "devices" {
			t.Fatalf("ExecuteCommand called with unexpected command/args: %s %v", command, args)
		}
		return mockOutput, nil
	}, func() {
		devices, err := ListAndroidDevices()
		if err != nil {
			t.Fatalf("ListAndroidDevices() returned error: %v", err)
		}
		if !reflect.DeepEqual(devices, expectedDevices) {
			t.Errorf("Expected devices %v, got %v", expectedDevices, devices)
		}
	})
}

func TestListAndroidDevices_DeviceAndOffline(t *testing.T) {
	mockOutput := []string{
		"List of devices attached",
		"device1\tdevice",
		"device2\toffline",
	}
	expectedDevices := []string{"device1"}
	withMockExecuteCommand(func(command string, args ...string) ([]string, error) {
		if command != "adb" || args[0] != "devices" {
			t.Fatalf("ExecuteCommand called with unexpected command/args: %s %v", command, args)
		}
		return mockOutput, nil
	}, func() {
		devices, err := ListAndroidDevices()
		if err != nil {
			t.Fatalf("ListAndroidDevices() returned error: %v", err)
		}
		if !reflect.DeepEqual(devices, expectedDevices) {
			t.Errorf("Expected devices %v, got %v", expectedDevices, devices)
		}
	})
}

func TestListAndroidDevices_ExtraFields(t *testing.T) {
	mockOutput := []string{
		"List of devices attached",
		"device1\tdevice product:N2G47H model:Nexus_5x device:bullhead",
	}
	expectedDevices := []string{"device1"}
	withMockExecuteCommand(func(command string, args ...string) ([]string, error) {
		if command != "adb" || args[0] != "devices" {
			t.Fatalf("ExecuteCommand called with unexpected command/args: %s %v", command, args)
		}
		return mockOutput, nil
	}, func() {
		devices, err := ListAndroidDevices()
		if err != nil {
			t.Fatalf("ListAndroidDevices() returned error: %v", err)
		}
		if !reflect.DeepEqual(devices, expectedDevices) {
			t.Errorf("Expected devices %v, got %v", expectedDevices, devices)
		}
	})
}

// Test for DisplayDevices (basic panic test)
func TestDisplayDevices_RunsWithoutPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("DisplayDevices panicked: %v", r)
		}
	}()
	DisplayDevices("TestType", []string{"device1", "device2"})
	DisplayDevices("TestTypeEmpty", []string{})
	DisplayDevices("TestTypeNil", nil)
}
