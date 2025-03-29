## devices-cli: A Golang CLI tool to list all connected devices
### Prerequisites for Cloning the Repository

Before cloning this repository, ensure you have the following prerequisites installed:

1. **Git**: Install Git by following the instructions on the [official Git website](https://git-scm.com/).
    - Verify the installation:
      ```bash
      git --version
      ```

2. **Golang**: Ensure Golang is installed and properly set up. Refer to the [Golang Installation Instructions](#golang-installation-instructions) section above.

3. **Homebrew** (for macOS users): Install Homebrew by following the instructions on the [Homebrew website](https://brew.sh/).
    - Verify the installation:
      ```bash
      brew --version
      ```

Once the prerequisites are installed, you can proceed to clone the repository:
```bash
git clone https://github.com/DineshKuppan/devices-cli.git
cd devices-cli
```

### Golang Installation Instructions

To install Golang, follow these steps:

1. Download the latest version of Go from the [official website](https://golang.org/dl/).
2. Follow the installation instructions for your operating system.
3. Verify the installation by running:
    ```bash
    go version
    ```

### Setting Up macOS Tools for iOS Devices

1. Install `libimobiledevice` tools using Homebrew:
    ```bash
    brew install libimobiledevice
    ```
2. Verify the installation:
    ```bash
    idevice_id -l
    ideviceinfo -u <device_udid>
    ```

### Setting Up ADB for Android Devices

1. Install Android Platform Tools:
    ```bash
    brew install --cask android-platform-tools
    ```
2. Add the Android Platform Tools to your PATH:
    ```bash
    export PATH=$PATH:/path/to/android/platform-tools
    ```
    Replace `/path/to/android/platform-tools` with the actual installation path.
3. Verify the installation:
    ```bash
    adb devices
    ```

### Enabling Developer Settings

#### For iOS Devices
1. Open the **Settings** app on your iOS device.
2. Navigate to **Privacy & Security** > **Developer Mode**.
3. Enable **Developer Mode** and restart your device if prompted.

#### For Android Devices
1. Open the **Settings** app on your Android device.
2. Navigate to **About Phone** (or **About Device**, depending on your device).
3. Locate the **Build Number** and tap it 7 times to enable Developer Options. You may need to enter your device's PIN or password.
4. For detailed instructions specific to your device, refer to [this guide on XDA Developers](https://www.xda-developers.com/android-developer-options/).
5. Go back to **Settings** and open **Developer Options**.
6. Enable **USB Debugging**.


#### For IOS Devices 
For example, the connected device id will be displayed with `UUID` and use `ideviceinfo` command to fetch info about the devices

To learn further information about the commands available for library `libimobiledevice`, please visit this link [Link](https://github.com/libimobiledevice/libimobiledevice)

```bash
idevice_id -l
ideviceinfo -u 00000000-0000000XXXXXXXXX
```

#### For Android

To fetch connected device model details and sdk version

```bash
adb -s <device_serial_id> shell getprop ro.product.model
adb -s <device_serial_id> shell getprop ro.build.version.release
```

##### Output

```bash
Pixel 7a
15
```

### Building and Running the Binary

To build and run the `devices-cli` binary, follow these steps:

1. Build the binary using the `go build` command:
    ```bash
    go build -o devices
    ```
2. Run the binary:
    ```bash
    ./devices
    ```

This will execute the `devices-cli` tool and display the list of connected devices.

```bash
List of connected Android & IOS devices...
Connected Devices: iOS
- iOS Device: 00000000-0000000XXXXXXXXX
Connected Devices: Android
- Android Device: <actual_device_serial_id>
```

### License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

### Reporting Issues

If you encounter any issues or have suggestions for improvements, please open an issue on the [GitHub Issues](https://github.com/DineshKuppan/devices-cli/issues) page. Provide as much detail as possible, including steps to reproduce the issue and any relevant logs or screenshots.