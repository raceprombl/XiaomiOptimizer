package adb

import (
	"bufio"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const (
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBold   = "\033[1m"
)

// Path to resolved adb binary
var resolvedAdbPath string

// SetAdbPath overrides or sets the resolved adb path
func SetAdbPath(path string) {
	resolvedAdbPath = path
}

// GetAdbPath returns current adb binary path
func GetAdbPath() string {
	return resolvedAdbPath
}

// RunADB executes an adb command with arguments
func RunADB(args ...string) (string, error) {
	cmd := exec.Command(resolvedAdbPath, args...)
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

// RunADBShell executes a shell command on the connected Android device
func RunADBShell(cmdStr string) (string, error) {
	return RunADB("shell", cmdStr)
}

// RunADBShellForDevice executes a shell command targeting a specific device serial
func RunADBShellForDevice(cmdStr string, serial *string) (string, error) {
	if serial != nil && *serial != "" {
		return RunADB("-s", *serial, "shell", cmdStr)
	}
	return RunADB("shell", cmdStr)
}

// ConnectWiFi attempts to connect to an ADB daemon over TCP/IP
func ConnectWiFi(address string) (string, error) {
	address = strings.TrimSpace(address)
	if address == "" {
		return "", fmt.Errorf("empty IP address provided")
	}
	if !strings.Contains(address, ":") {
		address = address + ":5555"
	}
	return RunADB("connect", address)
}

// PairWiFi pairs with an Android 11+ wireless debugging host using a 6-digit code
func PairWiFi(address string, code string) (string, error) {
	address = strings.TrimSpace(address)
	code = strings.TrimSpace(code)
	if address == "" || code == "" {
		return "", fmt.Errorf("address and pairing code are required")
	}
	return RunADB("pair", address, code)
}

// EnableTCPMode restarts the device's adbd daemon listening on specified TCP port
func EnableTCPMode(port int, serial *string) (string, error) {
	if port <= 0 {
		port = 5555
	}
	portStr := strconv.Itoa(port)
	if serial != nil && *serial != "" {
		return RunADB("-s", *serial, "tcpip", portStr)
	}
	return RunADB("tcpip", portStr)
}

// DisconnectWiFi disconnects an existing TCP/IP adb session
func DisconnectWiFi(address string) (string, error) {
	address = strings.TrimSpace(address)
	if address == "" {
		return RunADB("disconnect")
	}
	return RunADB("disconnect", address)
}

// GetDeviceWiFiIP retrieves the device's local WLAN IPv4 address
func GetDeviceWiFiIP(serial *string) string {
	// Try 1: ip -f inet addr show wlan0
	out, _ := RunADBShellForDevice("ip -f inet addr show wlan0", serial)
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "inet ") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				ipMask := fields[1]
				parts := strings.Split(ipMask, "/")
				if len(parts) >= 1 && strings.Count(parts[0], ".") == 3 {
					return parts[0]
				}
			}
		}
	}

	// Try 2: ip route
	outRoute, _ := RunADBShellForDevice("ip route", serial)
	for _, line := range strings.Split(outRoute, "\n") {
		if strings.Contains(line, "wlan0") && strings.Contains(line, "src ") {
			fields := strings.Fields(line)
			for i, f := range fields {
				if f == "src" && i+1 < len(fields) {
					ip := fields[i+1]
					if strings.Count(ip, ".") == 3 {
						return ip
					}
				}
			}
		}
	}

	// Try 3: getprop dhcp.wlan0.ipaddress
	propIP, _ := RunADBShellForDevice("getprop dhcp.wlan0.ipaddress", serial)
	propIP = strings.TrimSpace(propIP)
	if strings.Count(propIP, ".") == 3 {
		return propIP
	}

	return ""
}

// IsDisconnectedError checks if an error output signifies a disconnected or offline device
func IsDisconnectedError(out string, err error) bool {
	if err != nil {
		s := strings.ToLower(out)
		if strings.Contains(s, "device not found") ||
			strings.Contains(s, "no devices/emulators found") ||
			strings.Contains(s, "device offline") ||
			strings.Contains(s, "device unauthorized") ||
			strings.Contains(s, "error: closed") ||
			strings.Contains(s, "connection reset") ||
			strings.Contains(s, "transport error") {
			return true
		}
	}
	return false
}

// CheckDeviceState queries adb devices to inspect connected status
func CheckDeviceState() (serial string, status string) {
	out, err := RunADB("devices")
	if err != nil {
		return "", "error"
	}

	lines := strings.Split(out, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "List of devices") || strings.HasPrefix(line, "*") {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			return parts[0], parts[1]
		}
	}
	return "", "none"
}

// EnsureDeviceConnected blocks until device is reconnected during an active operation
func EnsureDeviceConnected(serial *string) {
	curSerial, state := CheckDeviceState()
	if state == "device" {
		if *serial == "" {
			*serial = curSerial
		}
		return
	}

	isWiFi := serial != nil && strings.Contains(*serial, ":")
	if isWiFi {
		fmt.Printf(colorYellow+colorBold+"\n[!] Wi-Fi connection to %s lost. Attempting to reconnect..."+colorReset+"\n", *serial)
	} else {
		fmt.Println(colorYellow + colorBold + "\n[!] Connection lost. Please plug your phone back in via USB to continue..." + colorReset)
	}

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if isWiFi {
			_, _ = ConnectWiFi(*serial)
		}

		s, st := CheckDeviceState()
		if st == "unauthorized" {
			fmt.Printf("\r%s[!] Device %s detected, but not authorized. Confirm debugging on phone...%s", colorYellow, s, colorReset)
			continue
		}
		if st == "device" {
			*serial = s
			fmt.Println(colorGreen + colorBold + "\n[✓] Phone reconnected. Resuming...\n" + colorReset)
			time.Sleep(600 * time.Millisecond)
			return
		}
		dots := strings.Repeat(".", int(time.Now().Unix()%4))
		if isWiFi {
			fmt.Printf("\r%s[⏳] Reconnecting to Wi-Fi (%s)%s%-4s%s", colorYellow, *serial, colorBold, dots, colorReset)
		} else {
			fmt.Printf("\r%s[⏳] Waiting for phone to be reconnected%s%-4s%s", colorYellow, colorBold, dots, colorReset)
		}
	}
}

// RunADBShellWithRetry executes an adb shell command with automatic retry upon reconnection
func RunADBShellWithRetry(cmdStr string, serial *string) (string, error) {
	for {
		EnsureDeviceConnected(serial)
		out, err := RunADBShellForDevice(cmdStr, serial)
		if IsDisconnectedError(out, err) {
			EnsureDeviceConnected(serial)
			continue
		}
		return out, err
	}
}

// GetProp reads a system property via getprop with reconnection retry
func GetProp(prop string, serial *string) string {
	res, err := RunADBShellWithRetry("getprop "+prop, serial)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(res)
}

// BrandInfo holds manufacturer and brand verification details
type BrandInfo struct {
	IsXiaomi   bool
	BrandName  string
	Vendor     string
	VerifiedBy string
}

// VerifyXiaomiBrand checks whether the connected device belongs to Xiaomi or its sub-brands (Redmi, POCO, Black Shark)
func VerifyXiaomiBrand(serial *string) BrandInfo {
	brand := strings.ToLower(GetProp("ro.product.brand", serial))
	mfg := strings.ToLower(GetProp("ro.product.manufacturer", serial))
	vendorBrand := strings.ToLower(GetProp("ro.product.vendor.brand", serial))
	vendorMfg := strings.ToLower(GetProp("ro.product.vendor.manufacturer", serial))
	systemBrand := strings.ToLower(GetProp("ro.product.system.brand", serial))
	systemMfg := strings.ToLower(GetProp("ro.product.system.manufacturer", serial))
	productName := strings.ToLower(GetProp("ro.product.name", serial))

	miuiVer := GetProp("ro.miui.ui.version.name", serial)
	hyperOSVer := GetProp("ro.mi.os.version.name", serial)
	buildIncremental := GetProp("ro.build.version.incremental", serial)

	allProps := strings.Join([]string{brand, mfg, vendorBrand, vendorMfg, systemBrand, systemMfg, productName}, " ")

	// 1. Check for POCO
	if strings.Contains(allProps, "poco") {
		return BrandInfo{
			IsXiaomi:   true,
			BrandName:  "POCO",
			Vendor:     "Xiaomi Corporation",
			VerifiedBy: "ro.product.brand / manufacturer",
		}
	}

	// 2. Check for Redmi
	if strings.Contains(allProps, "redmi") {
		return BrandInfo{
			IsXiaomi:   true,
			BrandName:  "Redmi",
			Vendor:     "Xiaomi Corporation",
			VerifiedBy: "ro.product.brand / manufacturer",
		}
	}

	// 3. Check for Black Shark
	if strings.Contains(allProps, "blackshark") || strings.Contains(allProps, "black shark") {
		return BrandInfo{
			IsXiaomi:   true,
			BrandName:  "Black Shark",
			Vendor:     "Blackshark Technologies",
			VerifiedBy: "ro.product.brand",
		}
	}

	// 4. Check for Xiaomi (Mi)
	if strings.Contains(allProps, "xiaomi") {
		return BrandInfo{
			IsXiaomi:   true,
			BrandName:  "Xiaomi",
			Vendor:     "Xiaomi Corporation",
			VerifiedBy: "ro.product.brand / manufacturer",
		}
	}

	// 5. Check for MIUI or HyperOS props (even on custom or modified ROMs)
	if hyperOSVer != "" {
		return BrandInfo{
			IsXiaomi:   true,
			BrandName:  "Xiaomi / HyperOS",
			Vendor:     "Xiaomi Corporation",
			VerifiedBy: "ro.mi.os.version.name (" + hyperOSVer + ")",
		}
	}
	if miuiVer != "" || strings.HasPrefix(buildIncremental, "V") {
		return BrandInfo{
			IsXiaomi:   true,
			BrandName:  "Xiaomi / MIUI",
			Vendor:     "Xiaomi Corporation",
			VerifiedBy: "ro.miui.ui.version.name (" + miuiVer + ")",
		}
	}

	// Non-Xiaomi device detected
	displayBrand := brand
	if displayBrand == "" {
		displayBrand = mfg
	}
	if displayBrand == "" {
		displayBrand = "Unknown OEM"
	} else {
		displayBrand = strings.ToUpper(displayBrand[:1]) + displayBrand[1:]
	}

	return BrandInfo{
		IsXiaomi:   false,
		BrandName:  displayBrand,
		Vendor:     mfg,
		VerifiedBy: "Non-Xiaomi hardware",
	}
}

// DeviceDetails holds hardware, system, and battery information
type DeviceDetails struct {
	Model       string
	Codename    string
	Serial      string
	AndroidVer  string
	OSVersion   string
	Battery     string
	Brand       string
	IsXiaomi    bool
	BrandStatus string
	IsWireless  bool
}

// FetchDeviceDetails queries the connected device for details, supporting MIUI 11 through HyperOS 3
func FetchDeviceDetails(serial *string) DeviceDetails {
	brandInfo := VerifyXiaomiBrand(serial)

	model := GetProp("ro.product.model", serial)
	if model == "" {
		model = "Unknown"
	}

	codename := GetProp("ro.product.name", serial)
	if codename == "" {
		codename = GetProp("ro.product.device", serial)
	}
	if codename == "" {
		codename = "generic"
	}

	ser := GetProp("ro.serialno", serial)
	if ser == "" {
		ser = *serial
	}

	androidVer := GetProp("ro.build.version.release", serial)
	if androidVer == "" {
		androidVer = "Unknown"
	}

	var osName string
	hyperOS := GetProp("ro.mi.os.version.name", serial)
	miuiVersion := GetProp("ro.miui.ui.version.name", serial)
	buildIncremental := GetProp("ro.build.version.incremental", serial)

	if hyperOS != "" {
		osName = fmt.Sprintf("Xiaomi HyperOS (%s)", hyperOS)
		if strings.HasPrefix(hyperOS, "OS3.") {
			osName = fmt.Sprintf("Xiaomi HyperOS 3.0 [%s]", hyperOS)
		} else if strings.HasPrefix(hyperOS, "OS2.") {
			osName = fmt.Sprintf("Xiaomi HyperOS 2.0 [%s]", hyperOS)
		} else if strings.HasPrefix(hyperOS, "OS1.") {
			osName = fmt.Sprintf("Xiaomi HyperOS 1.0 [%s]", hyperOS)
		}
	} else if miuiVersion != "" {
		verTag := miuiVersion
		switch miuiVersion {
		case "V14":
			verTag = "MIUI 14"
		case "V13":
			verTag = "MIUI 13"
		case "V125":
			verTag = "MIUI 12.5"
		case "V12":
			if strings.Contains(buildIncremental, "12.5") {
				verTag = "MIUI 12.5"
			} else {
				verTag = "MIUI 12"
			}
		case "V11":
			verTag = "MIUI 11"
		default:
			if strings.HasPrefix(miuiVersion, "V") {
				verTag = "MIUI " + strings.TrimPrefix(miuiVersion, "V")
			}
		}

		if buildIncremental != "" {
			osName = fmt.Sprintf("%s (%s)", verTag, buildIncremental)
		} else {
			osName = verTag
		}
	} else {
		if buildIncremental != "" {
			if strings.Contains(buildIncremental, "V11.") {
				osName = fmt.Sprintf("MIUI 11 (%s)", buildIncremental)
			} else if strings.Contains(buildIncremental, "V12.5.") {
				osName = fmt.Sprintf("MIUI 12.5 (%s)", buildIncremental)
			} else if strings.Contains(buildIncremental, "V12.") {
				osName = fmt.Sprintf("MIUI 12 (%s)", buildIncremental)
			} else {
				osName = fmt.Sprintf("MIUI / AOSP (%s)", buildIncremental)
			}
		} else {
			osName = "MIUI / HyperOS Custom"
		}
	}

	batteryStr := ParseBattery(serial)

	brandStatus := fmt.Sprintf("%s (Xiaomi Ecosystem) [Verified ✓]", brandInfo.BrandName)
	if !brandInfo.IsXiaomi {
		brandStatus = fmt.Sprintf("%s [NON-XIAOMI WARNING ✗]", brandInfo.BrandName)
	}

	isWireless := strings.Contains(ser, ":") || (serial != nil && strings.Contains(*serial, ":"))

	return DeviceDetails{
		Model:       model,
		Codename:    codename,
		Serial:      ser,
		AndroidVer:  "Android " + androidVer,
		OSVersion:   osName,
		Battery:     batteryStr,
		Brand:       brandInfo.BrandName,
		IsXiaomi:    brandInfo.IsXiaomi,
		BrandStatus: brandStatus,
		IsWireless:  isWireless,
	}
}

// ParseBattery extracts level, status, and temperature from dumpsys battery
func ParseBattery(serial *string) string {
	out, err := RunADBShellWithRetry("dumpsys battery", serial)
	if err != nil {
		return "N/A"
	}

	level := 0
	statusText := "Discharging"
	tempC := 0.0

	scanner := bufio.NewScanner(strings.NewReader(out))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "level:") {
			parts := strings.Split(line, ":")
			if len(parts) == 2 {
				level, _ = strconv.Atoi(strings.TrimSpace(parts[1]))
			}
		} else if strings.HasPrefix(line, "status:") {
			parts := strings.Split(line, ":")
			if len(parts) == 2 {
				statusCode, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
				switch statusCode {
				case 2:
					statusText = "Charging (+)"
				case 3, 4:
					statusText = "Discharging"
				case 5:
					statusText = "Full (100%)"
				}
			}
		} else if strings.HasPrefix(line, "temperature:") {
			parts := strings.Split(line, ":")
			if len(parts) == 2 {
				rawTemp, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
				if rawTemp > 0 {
					tempC = float64(rawTemp) / 10.0
				}
			}
		}
	}

	if tempC > 0 {
		return fmt.Sprintf("%d%% [%s, %.1f°C]", level, statusText, tempC)
	}
	return fmt.Sprintf("%d%% [%s]", level, statusText)
}
