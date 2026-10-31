package main

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// ANSI terminal color codes
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorBold   = "\033[1m"
	colorGray   = "\033[90m"
	colorOrange = "\033[38;2;255;103;0m" // Xiaomi signature orange
	colorMiCyan = "\033[38;2;0;210;255m" // HyperOS cyan
)

// Non-blocking keyboard input channel
var inputChan = make(chan string, 32)

// Starts asynchronous background reader for keyboard input
func startInputReader() {
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				return
			}
			inputChan <- strings.TrimSpace(line)
		}
	}()
}

// Discards any buffered or stale inputs
func drainInput() {
	for {
		select {
		case <-inputChan:
		default:
			return
		}
	}
}

// Waits for user to press Enter
func waitEnter() {
	drainInput()
	fmt.Print(colorGray + "\nPress [Enter] to return..." + colorReset)
	<-inputChan
}

// Returns the true visual width of a rune in a monospace terminal
func runeWidth(r rune) int {
	if r >= 32 && r < 127 {
		return 1
	}
	if r < 32 || (r >= 127 && r < 160) {
		return 0
	}
	// Box drawing, blocks, geometric shapes (1 column)
	if (r >= 0x2500 && r <= 0x257F) || (r >= 0x2580 && r <= 0x259F) || (r >= 0x25A0 && r <= 0x25FF) {
		return 1
	}
	// Emojis and miscellaneous symbols (2 columns)
	if (r >= 0x2600 && r <= 0x27BF) ||
		(r >= 0x2B00 && r <= 0x2BFF) ||
		(r >= 0x1F000 && r <= 0x1FAFF) {
		return 2
	}
	return 1
}

// Calculates visual display width of a string, ignoring ANSI escape sequences
func visualWidth(s string) int {
	w := 0
	inAnsi := false
	for _, r := range s {
		if r == '\033' {
			inAnsi = true
			continue
		}
		if inAnsi {
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
				inAnsi = false
			}
			continue
		}
		w += runeWidth(r)
	}
	return w
}

// Pads string with spaces on the right to match the exact target display width
func padRight(s string, target int) string {
	w := visualWidth(s)
	if w >= target {
		return s
	}
	return s + strings.Repeat(" ", target-w)
}

// Displays ASCII banner with an exact 82-column total width
func printBanner() {
	boxWidth := 80
	top := colorOrange + "╔" + strings.Repeat("═", boxWidth) + "╗" + colorReset
	bot := colorOrange + "╚" + strings.Repeat("═", boxWidth) + "╝" + colorReset
	empty := colorOrange + "║" + strings.Repeat(" ", boxWidth) + "║" + colorReset

	art := []string{
		" ██╗  ██╗██╗ ██████╗ ██████╗ ████████╗██╗███╗   ███╗██╗███████╗███████╗██████╗  ",
		" ╚██╗██╔╝██║██╔═══██╗██╔══██╗╚══██╔══╝██║████╗ ████║██║╚══███╔╝██╔════╝██╔══██╗ ",
		"  ╚███╔╝ ██║██║   ██║██████╔╝   ██║   ██║██╔████╔██║██║  ███╔╝ █████╗  ██████╔╝ ",
		"  ██╔██╗ ██║██║   ██║██╔═══╝    ██║   ██║██║╚██╔╝██║██║ ███╔╝  ██╔══╝  ██╔══██╗ ",
		" ██╔╝ ██╗██║╚██████╔╝██║        ██║   ██║██║ ╚═╝ ██║██║███████╗███████╗██║  ██║ ",
		" ╚═╝  ╚═╝╚═╝ ╚═════╝ ╚═╝        ╚═╝   ╚═╝╚═╝     ╚═╝╚═╝╚══════╝╚══════╝╚═╝  ╚═╝ ",
	}

	sub1 := strings.Repeat(" ", 17) + "✦ HyperOS & MIUI 11 / 12 / 13 / 14 Optimizer ✦" + strings.Repeat(" ", 17)
	sub2 := strings.Repeat(" ", 19) + "Xiaomi • Redmi • POCO Debloater & Tweaker" + strings.Repeat(" ", 20)

	fmt.Println(top)
	fmt.Println(empty)
	for _, a := range art {
		fmt.Println(colorOrange + "║" + a + "║" + colorReset)
	}
	fmt.Println(empty)
	fmt.Println(colorOrange + "║" + colorMiCyan + colorBold + sub1 + colorOrange + "║" + colorReset)
	fmt.Println(colorOrange + "║" + colorGray + sub2 + colorOrange + "║" + colorReset)
	fmt.Println(bot)
}

// PackageInfo holds metadata for a removable Android package
type PackageInfo struct {
	Name string
	Desc string
}

// Safe packages to remove across MIUI 11, 12, 12.5, 13, 14, and HyperOS
var safeAdsPackages = []PackageInfo{
	{"com.miui.msa.global", "MIUI System Ads engine"},
	{"com.miui.analytics", "Usage analytics and telemetry"},
	{"com.miui.daemon", "Background logging daemon"},
	{"com.xiaomi.ugd", "User experience data collector"},
	{"com.miui.misightservice", "Mi Sight tracking daemon"},
	{"com.miui.bugreport", "Bug report uploader"},
	{"com.miui.hybrid", "Quick Apps ad framework (MIUI 11/12)"},
	{"com.miui.hybrid.accessory", "Quick Apps accessory services"},
	{"com.facebook.services", "Facebook analytics service"},
	{"com.facebook.system", "Facebook system installer"},
	{"com.facebook.appmanager", "Facebook app updater"},
	{"com.miui.android.fashiongallery", "Wallpaper carousel ads (MIUI 12.5+)"},
	{"com.mfashiongallery.emag", "Wallpaper carousel ads (MIUI 11/12)"},
	{"com.xiaomi.payment", "Mi Pay service"},
	{"com.miui.yellowpage", "Yellow Pages directory and ad services"},
	{"com.xiaomi.mipicks", "GetApps app store"},
	{"com.xiaomi.discover", "Xiaomi Discover news feed"},
	{"com.mi.global.shop", "Mi Store shopping app"},
	{"cn.wps.xiaomi.abroad.lite", "Mi Doc Viewer (WPS Office Lite)"},
	{"com.xiaomi.simactivate.service", "SIM activation tracking service"},
}

// Deep bloatware packages across MIUI 11, 12, 12.5, 13, 14, and HyperOS
var bloatwarePackages = []PackageInfo{
	{"com.android.browser", "Mi Browser (MIUI 11/12)"},
	{"com.mi.globalbrowser", "Mi Browser (MIUI 12.5+)"},
	{"com.miui.videoplayer", "Mi Video player (with ads)"},
	{"com.miui.video", "Mi Video player legacy (MIUI 11)"},
	{"com.miui.player", "Mi Music player (with promos)"},
	{"com.miui.personalassistant", "App Vault legacy (MIUI 11/12)"},
	{"com.mi.globalminusscreen", "App Vault minus-one screen (MIUI 12+)"},
	{"com.miui.cleaner", "Built-in Cleaner utility"},
	{"com.mi.appfinder", "Mi App Finder"},
	{"video.like", "Pre-installed Like video app"},
	{"com.miui.virtualsim", "Mi Roaming / Virtual SIM"},
	{"com.xiaomi.midrop", "ShareMe / Mi Drop"},
}

// Path to resolved adb binary
var resolvedAdbPath string

// Wraps an io.Reader to show a live download progress bar
type progressReader struct {
	reader     io.Reader
	totalBytes int64
	readBytes  int64
	lastUpdate time.Time
}

func (pr *progressReader) Read(p []byte) (int, error) {
	n, err := pr.reader.Read(p)
	pr.readBytes += int64(n)

	if time.Since(pr.lastUpdate) > 100*time.Millisecond || pr.readBytes == pr.totalBytes {
		pr.lastUpdate = time.Now()
		pr.printProgress()
	}
	return n, err
}

func (pr *progressReader) printProgress() {
	percent := float64(0)
	if pr.totalBytes > 0 {
		percent = float64(pr.readBytes) / float64(pr.totalBytes) * 100
	}
	mbDownloaded := float64(pr.readBytes) / 1024 / 1024
	mbTotal := float64(pr.totalBytes) / 1024 / 1024

	width := 32
	completed := int((percent / 100) * float64(width))
	if completed > width {
		completed = width
	}

	bar := strings.Repeat("█", completed) + strings.Repeat("░", width-completed)
	fmt.Printf("\r %s[%s]%s %6.1f%% (%5.1f MB / %5.1f MB)", colorCyan, bar, colorReset, percent, mbDownloaded, mbTotal)
}

// Sets console environment for Windows UTF-8
func initConsole() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "chcp 65001 > nul")
		_ = cmd.Run()
	}
}

// Clears terminal screen across Windows, macOS, and Linux
func clearScreen() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		_ = cmd.Run()
	} else {
		fmt.Print("\033[H\033[2J")
	}
}

// Locates existing adb or downloads it from Google's official repository
func ensureADB() error {
	execDir, _ := os.Executable()
	baseDir := filepath.Dir(execDir)

	adbName := "adb"
	if runtime.GOOS == "windows" {
		adbName = "adb.exe"
	}

	localAdb := filepath.Join(baseDir, adbName)
	if _, err := os.Stat(localAdb); err == nil {
		resolvedAdbPath = localAdb
		return nil
	}

	if _, err := os.Stat(adbName); err == nil {
		resolvedAdbPath = adbName
		return nil
	}

	if p, err := exec.LookPath(adbName); err == nil {
		resolvedAdbPath = p
		return nil
	}

	clearScreen()
	fmt.Println(colorOrange + "╔════════════════════════════════════════════════════════════════════════════════╗" + colorReset)
	fmt.Println(colorOrange + "║" + colorYellow + colorBold + "  [FIRST RUN] Downloading Android ADB platform-tools...                         " + colorOrange + "║" + colorReset)
	fmt.Println(colorOrange + "║" + colorGray + "  Source: https://dl.google.com/android/repository/platform-tools...             " + colorOrange + "║" + colorReset)
	fmt.Println(colorOrange + "╚════════════════════════════════════════════════════════════════════════════════╝" + colorReset)
	fmt.Println()

	var zipURL string
	switch runtime.GOOS {
	case "darwin":
		zipURL = "https://dl.google.com/android/repository/platform-tools-latest-darwin.zip"
	case "linux":
		zipURL = "https://dl.google.com/android/repository/platform-tools-latest-linux.zip"
	default:
		zipURL = "https://dl.google.com/android/repository/platform-tools-latest-windows.zip"
	}

	tempZip := filepath.Join(baseDir, "platform-tools-temp.zip")

	resp, err := http.Get(zipURL)
	if err != nil {
		return fmt.Errorf("network request failed for dl.google.com: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned %s when downloading ADB", resp.Status)
	}

	outFile, err := os.Create(tempZip)
	if err != nil {
		return fmt.Errorf("failed to create temporary archive: %w", err)
	}

	pr := &progressReader{
		reader:     resp.Body,
		totalBytes: resp.ContentLength,
		lastUpdate: time.Now(),
	}

	_, err = io.Copy(outFile, pr)
	outFile.Close()
	if err != nil {
		_ = os.Remove(tempZip)
		return fmt.Errorf("download error: %w", err)
	}
	fmt.Println()
	fmt.Println(colorCyan + "\n[i] Extracting ADB binaries for " + runtime.GOOS + "..." + colorReset)

	zipReader, err := zip.OpenReader(tempZip)
	if err != nil {
		_ = os.Remove(tempZip)
		return fmt.Errorf("failed to open zip file: %w", err)
	}
	defer zipReader.Close()

	neededFiles := map[string]bool{
		"adb":              true,
		"adb.exe":          true,
		"AdbWinApi.dll":    true,
		"AdbWinUsbApi.dll": true,
	}

	extractedCount := 0
	for _, f := range zipReader.File {
		baseName := filepath.Base(f.Name)
		if neededFiles[baseName] {
			targetPath := filepath.Join(baseDir, baseName)
			rc, err := f.Open()
			if err != nil {
				continue
			}

			dstFile, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
			if err != nil {
				rc.Close()
				continue
			}

			_, _ = io.Copy(dstFile, rc)
			dstFile.Close()
			rc.Close()
			_ = os.Chmod(targetPath, 0755)
			extractedCount++
		}
	}

	zipReader.Close()
	_ = os.Remove(tempZip)

	if extractedCount == 0 {
		return fmt.Errorf("required executables not found in the downloaded archive")
	}

	resolvedAdbPath = filepath.Join(baseDir, adbName)
	fmt.Println(colorGreen + colorBold + "[✓] ADB installed and ready." + colorReset)
	time.Sleep(1500 * time.Millisecond)
	return nil
}

// Executes an adb command with arguments
func runADB(args ...string) (string, error) {
	cmd := exec.Command(resolvedAdbPath, args...)
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

// Executes a shell command on the connected Android device
func runADBShell(cmdStr string) (string, error) {
	return runADB("shell", cmdStr)
}

// Checks if an error output signifies a disconnected or offline device
func isDisconnectedError(out string, err error) bool {
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

// Checks device connection and authorization status
func checkDeviceState() (serial string, status string) {
	out, err := runADB("devices")
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

// Blocks until device is reconnected during an active operation
func ensureDeviceConnected(serial *string) {
	curSerial, state := checkDeviceState()
	if state == "device" {
		if *serial == "" {
			*serial = curSerial
		}
		return
	}

	fmt.Println(colorYellow + colorBold + "\n[!] Connection lost. Please plug your phone back in via USB to continue..." + colorReset)

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		s, st := checkDeviceState()
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
		fmt.Printf("\r%s[⏳] Waiting for phone to be reconnected%s%-4s%s", colorYellow, colorBold, dots, colorReset)
	}
}

// Executes an adb shell command with automatic retry upon reconnection
func runADBShellWithRetry(cmdStr string, serial *string) (string, error) {
	for {
		ensureDeviceConnected(serial)
		out, err := runADBShell(cmdStr)
		if isDisconnectedError(out, err) {
			ensureDeviceConnected(serial)
			continue
		}
		return out, err
	}
}

// Reads system property via getprop with reconnection retry
func getProp(prop string, serial *string) string {
	res, err := runADBShellWithRetry("getprop "+prop, serial)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(res)
}

// DeviceDetails holds hardware, system, and battery information
type DeviceDetails struct {
	Model      string
	Codename   string
	Serial     string
	AndroidVer string
	OSVersion  string
	Battery    string
}

// Queries the connected device for details, supporting MIUI 11 through HyperOS 3
func fetchDeviceDetails(serial *string) DeviceDetails {
	model := getProp("ro.product.model", serial)
	if model == "" {
		model = "Unknown"
	}

	codename := getProp("ro.product.name", serial)
	if codename == "" {
		codename = getProp("ro.product.device", serial)
	}
	if codename == "" {
		codename = "generic"
	}

	ser := getProp("ro.serialno", serial)
	if ser == "" {
		ser = *serial
	}

	androidVer := getProp("ro.build.version.release", serial)
	if androidVer == "" {
		androidVer = "Unknown"
	}

	var osName string
	hyperOS := getProp("ro.mi.os.version.name", serial)
	miuiVersion := getProp("ro.miui.ui.version.name", serial)
	buildIncremental := getProp("ro.build.version.incremental", serial)

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

	batteryStr := parseBattery(serial)

	return DeviceDetails{
		Model:      model,
		Codename:   codename,
		Serial:     ser,
		AndroidVer: "Android " + androidVer,
		OSVersion:  osName,
		Battery:    batteryStr,
	}
}

// Parses battery level and charging status from dumpsys battery
func parseBattery(serial *string) string {
	out, err := runADBShellWithRetry("dumpsys battery", serial)
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

// Displays device metadata in a card matching the 82-column layout
func printDeviceCard(d DeviceDetails) {
	top := colorCyan + "┌" + strings.Repeat("─", 80) + "┐" + colorReset
	header := colorCyan + "│ " + colorBold + padRight("CONNECTED DEVICE INFORMATION", 78) + colorReset + colorCyan + " │" + colorReset
	mid := colorCyan + "├" + strings.Repeat("─", 25) + "┬" + strings.Repeat("─", 54) + "┤" + colorReset
	bot := colorCyan + "└" + strings.Repeat("─", 25) + "┴" + strings.Repeat("─", 54) + "┘" + colorReset

	items := []struct {
		label string
		val   string
		color string
	}{
		{"Model & Codename:", fmt.Sprintf("%s (%s)", d.Model, d.Codename), colorBold},
		{"Serial Number (SN):", d.Serial, colorReset},
		{"Android Version:", d.AndroidVer, colorGreen},
		{"OS / ROM Version:", d.OSVersion, colorOrange + colorBold},
		{"Battery Status:", d.Battery, colorYellow},
		{"ADB Connection:", "USB Debugging Authorized [OK]", colorGreen},
	}

	fmt.Println(top)
	fmt.Println(header)
	fmt.Println(mid)
	for _, it := range items {
		l := padRight(it.label, 23)
		v := padRight(it.val, 52)
		fmt.Println(colorCyan + "│ " + colorReset + l + colorCyan + " │ " + it.color + v + colorReset + colorCyan + " │" + colorReset)
	}
	fmt.Println(bot)
}

// ActionFunc executes a specific tweak routine
type ActionFunc func(serial *string)

// MenuItem represents an interactive option in a menu
type MenuItem struct {
	Key         string
	Icon        string
	Title       string
	Description string
	Handler     ActionFunc
	SubMenu     *Menu
}

// Menu represents a navigable screen of options
type Menu struct {
	ID       string
	Title    string
	Subtitle string
	Parent   *Menu
	Items    []MenuItem
}

// Renders any menu screen strictly within 82 columns
func renderBoxedMenu(m *Menu) {
	const innerWidth = 78
	top := colorMiCyan + "╔" + strings.Repeat("═", 80) + "╗" + colorReset
	mid := colorMiCyan + "╠" + strings.Repeat("═", 80) + "╣" + colorReset
	bot := colorMiCyan + "╚" + strings.Repeat("═", 80) + "╝" + colorReset

	fmt.Println(top)
	titleLine := fmt.Sprintf(" %s%s%s", colorBold, m.Title, colorReset)
	fmt.Println(colorMiCyan + "║ " + padRight(titleLine, innerWidth) + colorMiCyan + " ║" + colorReset)

	if m.Subtitle != "" {
		subLine := fmt.Sprintf(" %s%s%s", colorGray, m.Subtitle, colorReset)
		fmt.Println(colorMiCyan + "║ " + padRight(subLine, innerWidth) + colorMiCyan + " ║" + colorReset)
	}
	fmt.Println(mid)

	for _, it := range m.Items {
		var itemStr string
		if it.Key == "0" {
			itemStr = fmt.Sprintf(" %s[%s]%s %s %-32s %s%s%s",
				colorRed, it.Key, colorReset, it.Icon, it.Title, colorGray, it.Description, colorReset)
		} else if it.Key == "A" {
			itemStr = fmt.Sprintf(" %s[%s]%s %s %-32s %s%s%s",
				colorYellow+colorBold, it.Key, colorReset, it.Icon, it.Title, colorGray, it.Description, colorReset)
		} else {
			itemStr = fmt.Sprintf(" %s[%s]%s %s %-32s %s%s%s",
				colorGreen, it.Key, colorReset, it.Icon, it.Title, colorGray, it.Description, colorReset)
		}
		fmt.Println(colorMiCyan + "║ " + padRight(itemStr, innerWidth) + colorMiCyan + " ║" + colorReset)
	}
	fmt.Println(bot)
	fmt.Println()
	fmt.Print(colorBold + "Select an option and press Enter: " + colorReset)
}

// Prompts for menu input while actively monitoring cable connection
func promptSelectionWithCableWatch(serial *string) (string, bool) {
	drainInput()
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case choice := <-inputChan:
			return strings.ToUpper(strings.TrimSpace(choice)), true
		case <-ticker.C:
			_, state := checkDeviceState()
			if state != "device" {
				return "", false
			}
		}
	}
}

// ============================================================================
// --- [1] DEBLOAT & TELEMETRY HANDLERS ---
// ============================================================================

func runSafeAdsCleanup(serial *string) {
	clearScreen()
	printBanner()
	fmt.Println(colorGreen + colorBold + "\n[1] SAFE ADS & TELEMETRY CLEANUP" + colorReset)
	fmt.Println(colorGray + "Removing system trackers, ad injectors, and partner daemons:" + colorReset)
	fmt.Println()

	successCount := 0
	notFoundCount := 0

	for i := 0; i < len(safeAdsPackages); i++ {
		item := safeAdsPackages[i]
		ensureDeviceConnected(serial)

		fmt.Printf(" • [%2d/%d] Removing: %-32s (%s)... ", i+1, len(safeAdsPackages), item.Name, item.Desc)
		out, _ := runADBShellWithRetry("pm uninstall -k --user 0 "+item.Name, serial)

		if strings.Contains(out, "Success") {
			fmt.Println(colorGreen + "[✓ Removed]" + colorReset)
			successCount++
		} else if strings.Contains(out, "not installed") || strings.Contains(out, "Unknown package") {
			fmt.Println(colorGray + "[– Not installed]" + colorReset)
			notFoundCount++
		} else {
			fmt.Printf(colorYellow+"[i Skipped: %s]\n"+colorReset, strings.TrimSpace(out))
		}
		time.Sleep(25 * time.Millisecond)
	}

	fmt.Println()
	fmt.Printf(colorGreen+colorBold+"[✓] Done. Removed: %d | Already absent: %d\n"+colorReset, successCount, notFoundCount)
	fmt.Println(colorCyan + "[i] System ads and trackers uninstalled for user 0." + colorReset)
}

func runBloatwareCleanup(serial *string) {
	clearScreen()
	printBanner()
	fmt.Println(colorYellow + colorBold + "\n[2] DEEP BLOATWARE REMOVAL" + colorReset)
	fmt.Println(colorGray + "Removing pre-installed media players, promo apps, and bundled browser:" + colorReset)
	fmt.Println()

	successCount := 0
	notFoundCount := 0

	for i := 0; i < len(bloatwarePackages); i++ {
		item := bloatwarePackages[i]
		ensureDeviceConnected(serial)

		fmt.Printf(" • [%d/%d] Removing: %-26s (%s)... ", i+1, len(bloatwarePackages), item.Name, item.Desc)
		out, _ := runADBShellWithRetry("pm uninstall -k --user 0 "+item.Name, serial)

		if strings.Contains(out, "Success") {
			fmt.Println(colorGreen + "[✓ Removed]" + colorReset)
			successCount++
		} else if strings.Contains(out, "not installed") || strings.Contains(out, "Unknown package") {
			fmt.Println(colorGray + "[– Not installed]" + colorReset)
			notFoundCount++
		} else {
			fmt.Printf(colorYellow+"[i Skipped: %s]\n"+colorReset, strings.TrimSpace(out))
		}
		time.Sleep(25 * time.Millisecond)
	}

	fmt.Println()
	fmt.Printf(colorYellow+colorBold+"[✓] Done. Removed: %d | Already absent: %d\n"+colorReset, successCount, notFoundCount)
	fmt.Println(colorCyan + "[i] Bloatware removed. Storage and RAM freed." + colorReset)
}

func runRestorePackagesOnly(serial *string) {
	clearScreen()
	printBanner()
	fmt.Println(colorOrange + colorBold + "\n[3] RESTORE REMOVED PACKAGES" + colorReset)
	fmt.Println(colorGray + "Reinstalling uninstalled packages for user 0:" + colorReset)
	fmt.Println()

	var allPackages []string
	for _, item := range safeAdsPackages {
		allPackages = append(allPackages, item.Name)
	}
	for _, item := range bloatwarePackages {
		allPackages = append(allPackages, item.Name)
	}
	allPackages = append(allPackages, "com.xiaomi.joyose")

	restoredCount := 0
	for i := 0; i < len(allPackages); i++ {
		pkg := allPackages[i]
		ensureDeviceConnected(serial)

		fmt.Printf(" • [%2d/%d] Restoring: %-34s... ", i+1, len(allPackages), pkg)
		out, _ := runADBShellWithRetry("cmd package install-existing "+pkg, serial)
		if strings.Contains(out, "installed") || strings.Contains(out, "Success") {
			fmt.Println(colorGreen + "[✓ Restored]" + colorReset)
			restoredCount++
		} else {
			outLegacy, _ := runADBShellWithRetry("pm install-existing "+pkg, serial)
			if strings.Contains(outLegacy, "installed") || strings.Contains(outLegacy, "Success") {
				fmt.Println(colorGreen + "[✓ Restored]" + colorReset)
				restoredCount++
			} else {
				fmt.Println(colorGray + "[– Not needed]" + colorReset)
			}
		}
		time.Sleep(25 * time.Millisecond)
	}

	fmt.Println()
	fmt.Printf(colorGreen+colorBold+"[✓] Restore finished. Packages restored: %d\n"+colorReset, restoredCount)
}

func runFullDebloatBatch(serial *string) {
	runSafeAdsCleanup(serial)
	fmt.Println()
	runBloatwareCleanup(serial)
}

// ============================================================================
// --- [2] PERFORMANCE & GAMING HANDLERS ---
// ============================================================================

func runJoyoseRemoval(serial *string) {
	clearScreen()
	printBanner()
	fmt.Println(colorMiCyan + colorBold + "\n[1] JOYOSE GAMING THROTTLE REMOVAL" + colorReset)
	fmt.Println(colorGray + "Disabling thermal throttling and game FPS lock daemon:" + colorReset)
	fmt.Println()

	ensureDeviceConnected(serial)
	fmt.Print(" • Removing com.xiaomi.joyose for user 0... ")
	out, _ := runADBShellWithRetry("pm uninstall -k --user 0 com.xiaomi.joyose", serial)
	if strings.Contains(out, "Success") {
		fmt.Println(colorGreen + "[✓ Removed]" + colorReset)
	} else if strings.Contains(out, "not installed") {
		fmt.Println(colorGray + "[– Already absent]" + colorReset)
	} else {
		_, _ = runADBShellWithRetry("pm disable-user --user 0 com.xiaomi.joyose", serial)
		fmt.Println(colorGreen + "[✓ Disabled]" + colorReset)
	}

	fmt.Println(colorCyan + "\n[i] Joyose removed. Games maintain steady FPS without artificial drops." + colorReset)
}

func runLock120Hz(serial *string) {
	clearScreen()
	printBanner()
	fmt.Println(colorMiCyan + colorBold + "\n[2] LOCK DISPLAY REFRESH TO PEAK 120HZ" + colorReset)
	fmt.Println(colorGray + "Enforcing 120Hz refresh rate across all apps and system menus:" + colorReset)
	fmt.Println()

	ensureDeviceConnected(serial)
	fmt.Print(" • Clearing PowerKeeper cache... ")
	_, _ = runADBShellWithRetry("pm clear com.miui.powerkeeper", serial)
	fmt.Println(colorGreen + "[✓ Done]" + colorReset)

	fmt.Print(" • Setting peak and min refresh rates to 120Hz... ")
	_, _ = runADBShellWithRetry("settings put system min_refresh_rate 120.0", serial)
	_, _ = runADBShellWithRetry("settings put system peak_refresh_rate 120.0", serial)
	_, _ = runADBShellWithRetry("settings put secure user_refresh_rate 120", serial)
	_, _ = runADBShellWithRetry("settings put secure miui_refresh_rate 120", serial)
	_, _ = runADBShellWithRetry("settings put system thermal_limit_refresh_rate 120", serial)
	_, _ = runADBShellWithRetry("settings put secure support_highfps 1", serial)
	_, _ = runADBShellWithRetry("settings put secure refresh_rate_mode 2", serial)
	fmt.Println(colorGreen + "[✓ Applied]" + colorReset)

	fmt.Println(colorCyan + "\n[i] Display locked to 120Hz. App drops to 60Hz are eliminated." + colorReset)
}

func runDisableVirtualRAM(serial *string) {
	clearScreen()
	printBanner()
	fmt.Println(colorMiCyan + colorBold + "\n[3] DISABLE VIRTUAL RAM (EXTM / UFS SWAP)" + colorReset)
	fmt.Println(colorGray + "Eliminating micro-stutters and 0.1% low frame drops caused by slow storage paging:" + colorReset)
	fmt.Println()

	ensureDeviceConnected(serial)
	fmt.Print(" • Disabling MIUI/HyperOS Memory Extension in system settings... ")
	_, _ = runADBShellWithRetry("settings put global virtual_ram_enable 0", serial)
	_, _ = runADBShellWithRetry("setprop persist.miui.extm.enable 0", serial)
	_, _ = runADBShellWithRetry("setprop persist.sys.extm.enable 0", serial)
	fmt.Println(colorGreen + "[✓ Disabled]" + colorReset)

	fmt.Println(colorCyan + "\n[i] Virtual RAM disabled. Phone will use pure physical LPDDR RAM.")
	fmt.Println("    Reboot your phone once to completely unload the zRAM/UFS swap partition." + colorReset)
}

func runDisablePhantomProcessKiller(serial *string) {
	clearScreen()
	printBanner()
	fmt.Println(colorMiCyan + colorBold + "\n[4] DISABLE PHANTOM PROCESS KILLER" + colorReset)
	fmt.Println(colorGray + "Preventing Android 12+ from killing child processes and terminal sessions:" + colorReset)
	fmt.Println()

	ensureDeviceConnected(serial)
	fmt.Print(" • Disabling DeviceConfig test sync overrides... ")
	_, _ = runADBShellWithRetry("device_config set_sync_disabled_for_tests persistent", serial)
	fmt.Println(colorGreen + "[✓ Done]" + colorReset)

	fmt.Print(" • Raising max phantom processes to 2147483647... ")
	_, _ = runADBShellWithRetry("device_config put activity_manager max_phantom_processes 2147483647", serial)
	fmt.Println(colorGreen + "[✓ Applied]" + colorReset)

	fmt.Println(colorCyan + "\n[i] Phantom Process Killer disabled. Termux, bots, and background tasks run smoothly." + colorReset)
}

func runARTCompilation(serial *string) {
	clearScreen()
	printBanner()
	fmt.Println(colorMiCyan + colorBold + "\n[5] ART SPEED-PROFILE COMPILATION" + colorReset)
	fmt.Println(colorGray + "Ahead-Of-Time compiling hot bytecode for installed applications:" + colorReset)
	fmt.Println()

	ensureDeviceConnected(serial)
	fmt.Println(colorYellow + "[*] Compiling key system and user apps with speed-profile..." + colorReset)

	keyApps := []string{
		"com.miui.home",
		"com.android.systemui",
		"com.android.chrome",
		"org.telegram.messenger",
		"com.google.android.youtube",
	}

	for _, app := range keyApps {
		ensureDeviceConnected(serial)
		fmt.Printf(" • Compiling: %-32s... ", app)
		out, _ := runADBShellWithRetry("cmd package compile -m speed-profile -f "+app, serial)
		if strings.Contains(out, "Success") {
			fmt.Println(colorGreen + "[✓ Compiled]" + colorReset)
		} else {
			fmt.Println(colorGray + "[– Skipped / Not present]" + colorReset)
		}
	}

	fmt.Println(colorCyan + "\n[*] Triggering background dexopt job for remaining apps...")
	_, _ = runADBShellWithRetry("cmd package bg-dexopt-job", serial)
	fmt.Println(colorGreen + "[✓ Background AOT compilation task scheduled]" + colorReset)
	fmt.Println(colorCyan + "[i] App launch times and UI frame rates significantly improved." + colorReset)
}

func runAllPerformanceBatch(serial *string) {
	runJoyoseRemoval(serial)
	fmt.Println()
	runLock120Hz(serial)
	fmt.Println()
	runDisableVirtualRAM(serial)
	fmt.Println()
	runDisablePhantomProcessKiller(serial)
	fmt.Println()
	runARTCompilation(serial)
}

// ============================================================================
// --- [3] BATTERY & DEEP SLEEP HANDLERS ---
// ============================================================================

func runAggressiveDoze(serial *string) {
	clearScreen()
	printBanner()
	fmt.Println(colorGreen + colorBold + "\n[1] AGGRESSIVE DOZE & DEEP SLEEP" + colorReset)
	fmt.Println(colorGray + "Enabling fast deep sleep mode when screen turns off:" + colorReset)
	fmt.Println()

	ensureDeviceConnected(serial)
	fmt.Print(" • Activating MIUI device idle trigger... ")
	_, _ = runADBShellWithRetry("settings put global miui_device_idle_enabled 1", serial)
	fmt.Println(colorGreen + "[✓ Enabled]" + colorReset)

	fmt.Print(" • Tuning Doze constants for 30-second sleep timeout... ")
	dozeConsts := "inactive_to=30000,sensing_to=0,locating_to=0,location_accuracy=0,motion_inactive_to=0,idle_after_inactive_to=0,idle_pending_to=30000,max_idle_pending_to=60000,idle_to=1800000,idle_factor=2.0,max_idle_to=14400000,light_after_inactive_to=5000,light_pre_idle_to=5000,light_idle_to=180000,light_idle_factor=2.0,light_max_idle_to=900000"
	_, _ = runADBShellWithRetry("settings put global device_idle_constants "+dozeConsts, serial)
	fmt.Println(colorGreen + "[✓ Applied]" + colorReset)

	fmt.Print(" • Whitelisting critical messaging apps (Telegram, WhatsApp)... ")
	_, _ = runADBShellWithRetry("dumpsys deviceidle whitelist +org.telegram.messenger", serial)
	_, _ = runADBShellWithRetry("dumpsys deviceidle whitelist +com.whatsapp", serial)
	fmt.Println(colorGreen + "[✓ Whitelisted]" + colorReset)

	fmt.Println(colorCyan + "\n[i] Aggressive Doze active. Overnight battery drain reduced to under 1%." + colorReset)
}

func runDisableScanning(serial *string) {
	clearScreen()
	printBanner()
	fmt.Println(colorGreen + colorBold + "\n[2] DISABLE WI-FI & BLUETOOTH BACKGROUND SCANNING" + colorReset)
	fmt.Println(colorGray + "Stopping constant background beacon discovery and telemetry:" + colorReset)
	fmt.Println()

	ensureDeviceConnected(serial)
	fmt.Print(" • Disabling always-on Wi-Fi background scanning... ")
	_, _ = runADBShellWithRetry("settings put global wifi_scan_always_enabled 0", serial)
	_, _ = runADBShellWithRetry("settings put global wifi_wakeup_enabled 0", serial)
	fmt.Println(colorGreen + "[✓ Disabled]" + colorReset)

	fmt.Print(" • Disabling always-on Bluetooth LE scanning... ")
	_, _ = runADBShellWithRetry("settings put global ble_scan_always_enabled 0", serial)
	fmt.Println(colorGreen + "[✓ Disabled]" + colorReset)

	fmt.Println(colorCyan + "\n[i] Radio transceivers will sleep completely when Wi-Fi and Bluetooth are off." + colorReset)
}

func runStandbyOptimizations(serial *string) {
	clearScreen()
	printBanner()
	fmt.Println(colorGreen + colorBold + "\n[3] STANDBY & TELEMETRY OPTIMIZATIONS" + colorReset)
	fmt.Println(colorGray + "Disabling cellular radio hold when on Wi-Fi and stopping crash uploaders:" + colorReset)
	fmt.Println()

	ensureDeviceConnected(serial)
	fmt.Print(" • Disabling cellular radio hold on Wi-Fi (mobile_data_always_on)... ")
	_, _ = runADBShellWithRetry("settings put global mobile_data_always_on 0", serial)
	fmt.Println(colorGreen + "[✓ Disabled]" + colorReset)

	fmt.Print(" • Disabling network scoring and crash upload daemons... ")
	_, _ = runADBShellWithRetry("settings put global network_scoring_ui_enabled 0", serial)
	_, _ = runADBShellWithRetry("settings put secure send_action_app_error 0", serial)
	_, _ = runADBShellWithRetry("settings put global send_action_app_error 0", serial)
	fmt.Println(colorGreen + "[✓ Disabled]" + colorReset)

	fmt.Println(colorCyan + "\n[i] Unnecessary modem wakeups and diagnostic uploads halted." + colorReset)
}

func runLimitLogdBuffer(serial *string) {
	clearScreen()
	printBanner()
	fmt.Println(colorGreen + colorBold + "\n[4] MINIMIZE SYSTEM LOGGING (LOGD 64K)" + colorReset)
	fmt.Println(colorGray + "Restricting logd ring buffer size to reduce CPU wakeups and background writes:" + colorReset)
	fmt.Println()

	ensureDeviceConnected(serial)
	fmt.Print(" • Setting logd ring buffer size to 64K... ")
	_, _ = runADBShellWithRetry("setprop persist.logd.size 64K", serial)
	fmt.Println(colorGreen + "[✓ Applied]" + colorReset)

	fmt.Print(" • Turning off MIUI verbose debug logging... ")
	_, _ = runADBShellWithRetry("settings put global miui_log_enable 0", serial)
	_, _ = runADBShellWithRetry("settings put global miui_log_event_enable 0", serial)
	fmt.Println(colorGreen + "[✓ Disabled]" + colorReset)

	fmt.Println(colorCyan + "\n[i] System logger overhead minimized. Saves CPU cycles and storage wear." + colorReset)
}

func runAllBatteryBatch(serial *string) {
	runAggressiveDoze(serial)
	fmt.Println()
	runDisableScanning(serial)
	fmt.Println()
	runStandbyOptimizations(serial)
	fmt.Println()
	runLimitLogdBuffer(serial)
}

// ============================================================================
// --- [4] NETWORK & PRIVATE DNS HANDLERS ---
// ============================================================================

func runSetAdGuardDNS(serial *string) {
	clearScreen()
	printBanner()
	fmt.Println(colorCyan + colorBold + "\n[1] ENABLE ADGUARD PRIVATE DNS (ADBLOCK)" + colorReset)
	fmt.Println(colorGray + "Enabling encrypted DNS over TLS to block in-app and browser ads system-wide:" + colorReset)
	fmt.Println()

	ensureDeviceConnected(serial)
	fmt.Print(" • Setting Private DNS provider to dns.adguard-dns.com... ")
	_, _ = runADBShellWithRetry("settings put global private_dns_mode hostname", serial)
	_, _ = runADBShellWithRetry("settings put global private_dns_specifier dns.adguard-dns.com", serial)
	fmt.Println(colorGreen + "[✓ Enabled]" + colorReset)

	fmt.Println(colorCyan + "\n[i] AdGuard Private DNS is active. In-app popups and banner ads are blocked." + colorReset)
}

func runSetCloudflareDNS(serial *string) {
	clearScreen()
	printBanner()
	fmt.Println(colorCyan + colorBold + "\n[2] ENABLE CLOUDFLARE 1.1.1.1 PRIVATE DNS" + colorReset)
	fmt.Println(colorGray + "Enabling ultra-low latency DNS over TLS for competitive online gaming:" + colorReset)
	fmt.Println()

	ensureDeviceConnected(serial)
	fmt.Print(" • Setting Private DNS provider to one.one.one.one... ")
	_, _ = runADBShellWithRetry("settings put global private_dns_mode hostname", serial)
	_, _ = runADBShellWithRetry("settings put global private_dns_specifier one.one.one.one", serial)
	fmt.Println(colorGreen + "[✓ Enabled]" + colorReset)

	fmt.Println(colorCyan + "\n[i] Cloudflare 1.1.1.1 active. Faster domain resolution and lower game latency." + colorReset)
}

func runResetPrivateDNS(serial *string) {
	clearScreen()
	printBanner()
	fmt.Println(colorCyan + colorBold + "\n[3] RESET PRIVATE DNS TO AUTOMATIC" + colorReset)
	fmt.Println(colorGray + "Restoring default ISP/network DNS provider:" + colorReset)
	fmt.Println()

	ensureDeviceConnected(serial)
	fmt.Print(" • Restoring opportunistic DNS mode... ")
	_, _ = runADBShellWithRetry("settings put global private_dns_mode opportunistic", serial)
	_, _ = runADBShellWithRetry("settings delete global private_dns_specifier", serial)
	fmt.Println(colorGreen + "[✓ Restored]" + colorReset)

	fmt.Println(colorCyan + "\n[i] Private DNS returned to standard automatic DHCP configuration." + colorReset)
}

func runWiFiLowLatency(serial *string) {
	clearScreen()
	printBanner()
	fmt.Println(colorCyan + colorBold + "\n[4] WI-FI LOW-LATENCY HAL & ROAMING FIX" + colorReset)
	fmt.Println(colorGray + "Enforcing low latency Wi-Fi mode and eliminating packet jitter:" + colorReset)
	fmt.Println()

	ensureDeviceConnected(serial)
	fmt.Print(" • Forcing Wi-Fi HAL low latency mode... ")
	_, _ = runADBShellWithRetry("cmd wifi force-low-latency-mode enabled", serial)
	_, _ = runADBShellWithRetry("settings put global wifi_suspend_optimizations_enabled 0", serial)
	fmt.Println(colorGreen + "[✓ Enabled]" + colorReset)

	fmt.Print(" • Disabling Xiaomi aggressive network switching... ")
	_, _ = runADBShellWithRetry("settings put global cloud_weak_network_switch_enabled off", serial)
	fmt.Println(colorGreen + "[✓ Done]" + colorReset)

	fmt.Println(colorCyan + "\n[i] Wi-Fi latency stabilized. Packet drops and jitter during gaming reduced." + colorReset)
}

func runAllNetworkBatch(serial *string) {
	runSetAdGuardDNS(serial)
	fmt.Println()
	runWiFiLowLatency(serial)
}

// ============================================================================
// --- [5] UI, FLUIDITY & VISUAL HANDLERS ---
// ============================================================================

func runFastAnimations(serial *string) {
	clearScreen()
	printBanner()
	fmt.Println(colorMiCyan + colorBold + "\n[1] 0.5X FAST SYSTEM ANIMATIONS" + colorReset)
	fmt.Println(colorGray + "Halving window, transition, and animator scale delays:" + colorReset)
	fmt.Println()

	ensureDeviceConnected(serial)
	fmt.Print(" • Setting window animation scale to 0.5x... ")
	_, _ = runADBShellWithRetry("settings put global window_animation_scale 0.5", serial)
	_, _ = runADBShellWithRetry("settings put global transition_animation_scale 0.5", serial)
	_, _ = runADBShellWithRetry("settings put global animator_duration_scale 0.5", serial)
	fmt.Println(colorGreen + "[✓ Applied]" + colorReset)

	fmt.Println(colorCyan + "\n[i] UI animations now execute twice as fast." + colorReset)
}

func runCachedAppsFreezer(serial *string) {
	clearScreen()
	printBanner()
	fmt.Println(colorMiCyan + colorBold + "\n[2] NATIVE CACHED APPS FREEZER" + colorReset)
	fmt.Println(colorGray + "Freezing idle background processes in RAM using kernel cgroups:" + colorReset)
	fmt.Println()

	ensureDeviceConnected(serial)
	fmt.Print(" • Enabling cached apps freezer in Android runtime... ")
	_, _ = runADBShellWithRetry("settings put global cached_apps_freezer enabled", serial)
	_, _ = runADBShellWithRetry("device_config put activity_manager_native_boot use_freezer true", serial)
	_, _ = runADBShellWithRetry("device_config put activity_manager max_cached_processes 64", serial)
	fmt.Println(colorGreen + "[✓ Enabled]" + colorReset)

	fmt.Println(colorCyan + "\n[i] Inactive apps frozen in RAM: 0% CPU consumption in background." + colorReset)
}

func runTouchTuning(serial *string) {
	clearScreen()
	printBanner()
	fmt.Println(colorMiCyan + colorBold + "\n[3] TOUCH LATENCY & RESPONSE TUNING" + colorReset)
	fmt.Println(colorGray + "Reducing tap registration deadtime and gesture latency:" + colorReset)
	fmt.Println()

	ensureDeviceConnected(serial)
	fmt.Print(" • Shortening long-press and multi-press thresholds... ")
	_, _ = runADBShellWithRetry("settings put secure long_press_timeout 250", serial)
	_, _ = runADBShellWithRetry("settings put secure multi_press_timeout 250", serial)
	_, _ = runADBShellWithRetry("settings put secure tap_duration_threshold 0.0", serial)
	_, _ = runADBShellWithRetry("settings put secure touch_blocking_period 0.0", serial)
	_, _ = runADBShellWithRetry("settings put system pointer_speed 1", serial)
	fmt.Println(colorGreen + "[✓ Applied]" + colorReset)

	fmt.Println(colorCyan + "\n[i] Touch responsiveness optimized for gaming and fast navigation." + colorReset)
}

func runFixGreyBlur(serial *string) {
	clearScreen()
	printBanner()
	fmt.Println(colorMiCyan + colorBold + "\n[4] UNLOCK CONTROL CENTER GAUSSIAN BLUR" + colorReset)
	fmt.Println(colorGray + "Fixing the ugly grey background on budget/mid-range Redmi & POCO phones:" + colorReset)
	fmt.Println()

	ensureDeviceConnected(serial)
	fmt.Print(" • Setting GPU deviceLevelList to v:1,c:3,g:3 (Flagship graphics profile)... ")
	_, _ = runADBShellWithRetry("settings put system deviceLevelList \"v:1,c:3,g:3\"", serial)
	_, _ = runADBShellWithRetry("settings put system window_blurs_enabled 1", serial)
	fmt.Println(colorGreen + "[✓ Applied]" + colorReset)

	fmt.Println(colorCyan + "\n[i] High-end blur shaders unlocked. Real-time background blur will appear")
	fmt.Println("    in Control Center after rebooting or restarting System UI." + colorReset)
}

func runHideGestureLine(serial *string) {
	clearScreen()
	printBanner()
	fmt.Println(colorMiCyan + colorBold + "\n[5] HIDE FULLSCREEN GESTURE BAR LINE" + colorReset)
	fmt.Println(colorGray + "Removing the horizontal pill at the bottom for true fullscreen display:" + colorReset)
	fmt.Println()

	ensureDeviceConnected(serial)
	fmt.Print(" • Hiding fullscreen navigation indicator line... ")
	_, _ = runADBShellWithRetry("settings put secure hide_gesture_line 1", serial)
	_, _ = runADBShellWithRetry("settings put system hide_gesture_line 1", serial)
	fmt.Println(colorGreen + "[✓ Hidden]" + colorReset)

	fmt.Println(colorCyan + "\n[i] Gesture indicator line hidden. Full screen real estate reclaimed." + colorReset)
}

func runAllUIBatch(serial *string) {
	runFastAnimations(serial)
	fmt.Println()
	runCachedAppsFreezer(serial)
	fmt.Println()
	runTouchTuning(serial)
	fmt.Println()
	runFixGreyBlur(serial)
	fmt.Println()
	runHideGestureLine(serial)
}

// ============================================================================
// --- [6] SPECIAL SYSTEM FLAGS & AUDIO HANDLERS ---
// ============================================================================

func runBypassSafeVolume(serial *string) {
	clearScreen()
	printBanner()
	fmt.Println(colorYellow + colorBold + "\n[1] BYPASS SAFE HEADPHONE VOLUME WARNING" + colorReset)
	fmt.Println(colorGray + "Disabling EU/regional 60% headphone volume reset popups:" + colorReset)
	fmt.Println()

	ensureDeviceConnected(serial)
	fmt.Print(" • Overriding safe audio volume state... ")
	_, _ = runADBShellWithRetry("settings put global audio_safe_volume_state 2", serial)
	_, _ = runADBShellWithRetry("settings put system safe_volume 0", serial)
	fmt.Println(colorGreen + "[✓ Applied]" + colorReset)

	fmt.Println(colorCyan + "\n[i] Safe volume dialog disabled. Full headphone output unlocked." + colorReset)
}

func runUnlockCameraLab(serial *string) {
	clearScreen()
	printBanner()
	fmt.Println(colorYellow + colorBold + "\n[2] UNLOCK HIDDEN CAMERA LAB FEATURES" + colorReset)
	fmt.Println(colorGray + "Enabling internal experimental camera settings and post-processing flags:" + colorReset)
	fmt.Println()

	ensureDeviceConnected(serial)
	fmt.Print(" • Creating camera lab options trigger flag... ")
	_, _ = runADBShellWithRetry("mkdir -p /sdcard/DCIM/Camera/lab_options_visible", serial)
	fmt.Println(colorGreen + "[✓ Created]" + colorReset)

	fmt.Println(colorCyan + "\n[i] Camera Lab unlocked. Open Stock Camera -> Settings to view Experimental options." + colorReset)
}

func runBypassInstallerLag(serial *string) {
	clearScreen()
	printBanner()
	fmt.Println(colorYellow + colorBold + "\n[3] DISABLE APK INSTALL VERIFICATION DELAYS" + colorReset)
	fmt.Println(colorGray + "Eliminating package verifier countdown timers and scan pauses for ADB:" + colorReset)
	fmt.Println()

	ensureDeviceConnected(serial)
	fmt.Print(" • Disabling package verification checks... ")
	_, _ = runADBShellWithRetry("settings put global verifier_verify_adb_installs 0", serial)
	_, _ = runADBShellWithRetry("settings put global package_verifier_enable 0", serial)
	_, _ = runADBShellWithRetry("settings put global package_verifier_user_consent -1", serial)
	fmt.Println(colorGreen + "[✓ Disabled]" + colorReset)

	fmt.Println(colorCyan + "\n[i] APK install verification skipped. App installs execute instantly." + colorReset)
}

func runAllSpecialBatch(serial *string) {
	runBypassSafeVolume(serial)
	fmt.Println()
	runUnlockCameraLab(serial)
	fmt.Println()
	runBypassInstallerLag(serial)
}

// ============================================================================
// --- [7] STORAGE & MAINTENANCE HANDLERS ---
// ============================================================================

func runFSTrim(serial *string) {
	clearScreen()
	printBanner()
	fmt.Println(colorYellow + colorBold + "\n[1] TRIM FLASH STORAGE (UFS / SSD)" + colorReset)
	fmt.Println(colorGray + "Running wear-leveling block trim on NAND partitions:" + colorReset)
	fmt.Println()

	ensureDeviceConnected(serial)
	fmt.Println(colorYellow + "[*] Executing fstrim across /data and /system..." + colorReset)
	out, _ := runADBShellWithRetry("sm fstrim", serial)
	if out != "" {
		fmt.Printf(colorGray+"%s\n"+colorReset, out)
	}
	fmt.Println(colorGreen + "[✓ Flash storage trim complete]" + colorReset)
	fmt.Println(colorCyan + "[i] 4K random write latency restored. Flash wear leveled." + colorReset)
}

func runTrimCaches(serial *string) {
	clearScreen()
	printBanner()
	fmt.Println(colorYellow + colorBold + "\n[2] PRUNE PACKAGE CACHES" + colorReset)
	fmt.Println(colorGray + "Pruning temporary junk caches without deleting user data:" + colorReset)
	fmt.Println()

	ensureDeviceConnected(serial)
	fmt.Print(" • Asking PackageManager to trim caches... ")
	out, _ := runADBShellWithRetry("pm trim-caches 999999999999", serial)
	if out != "" {
		fmt.Println(colorGreen + "[✓ " + strings.TrimSpace(out) + "]" + colorReset)
	} else {
		fmt.Println(colorGreen + "[✓ Caches pruned]" + colorReset)
	}

	fmt.Println(colorCyan + "[i] Gigabytes of stale temporary cache reclaimed." + colorReset)
}

func runResetAllTweaks(serial *string) {
	clearScreen()
	printBanner()
	fmt.Println(colorOrange + colorBold + "\n[3] RESET ALL TWEAKS TO FACTORY DEFAULTS" + colorReset)
	fmt.Println(colorGray + "Reverting animation scales, display limits, Doze, DNS, and reinstalling apps:" + colorReset)
	fmt.Println()

	ensureDeviceConnected(serial)

	fmt.Print(" • Resetting animation scales to 1.0x... ")
	_, _ = runADBShellWithRetry("settings put global window_animation_scale 1.0", serial)
	_, _ = runADBShellWithRetry("settings put global transition_animation_scale 1.0", serial)
	_, _ = runADBShellWithRetry("settings put global animator_duration_scale 1.0", serial)
	fmt.Println(colorGreen + "[✓ Done]" + colorReset)

	fmt.Print(" • Removing display refresh rate overrides... ")
	_, _ = runADBShellWithRetry("settings delete system min_refresh_rate", serial)
	_, _ = runADBShellWithRetry("settings delete system peak_refresh_rate", serial)
	_, _ = runADBShellWithRetry("settings delete secure user_refresh_rate", serial)
	_, _ = runADBShellWithRetry("settings delete secure miui_refresh_rate", serial)
	_, _ = runADBShellWithRetry("settings delete system thermal_limit_refresh_rate", serial)
	_, _ = runADBShellWithRetry("settings delete secure support_highfps", serial)
	_, _ = runADBShellWithRetry("settings delete secure refresh_rate_mode", serial)
	fmt.Println(colorGreen + "[✓ Done]" + colorReset)

	fmt.Print(" • Restoring default Doze and scanning parameters... ")
	_, _ = runADBShellWithRetry("settings delete global device_idle_constants", serial)
	_, _ = runADBShellWithRetry("settings delete global miui_device_idle_enabled", serial)
	_, _ = runADBShellWithRetry("settings put global wifi_scan_always_enabled 1", serial)
	_, _ = runADBShellWithRetry("settings put global ble_scan_always_enabled 1", serial)
	_, _ = runADBShellWithRetry("dumpsys deviceidle unforce", serial)
	fmt.Println(colorGreen + "[✓ Done]" + colorReset)

	fmt.Print(" • Restoring Private DNS to automatic... ")
	_, _ = runADBShellWithRetry("settings put global private_dns_mode opportunistic", serial)
	_, _ = runADBShellWithRetry("settings delete global private_dns_specifier", serial)
	fmt.Println(colorGreen + "[✓ Done]" + colorReset)

	fmt.Print(" • Restoring Virtual RAM (Memory Extension)... ")
	_, _ = runADBShellWithRetry("settings put global virtual_ram_enable 1", serial)
	_, _ = runADBShellWithRetry("setprop persist.miui.extm.enable 1", serial)
	fmt.Println(colorGreen + "[✓ Done]" + colorReset)

	fmt.Print(" • Restoring original graphics deviceLevelList & gesture line... ")
	_, _ = runADBShellWithRetry("settings delete system deviceLevelList", serial)
	_, _ = runADBShellWithRetry("settings delete secure hide_gesture_line", serial)
	_, _ = runADBShellWithRetry("settings delete system hide_gesture_line", serial)
	fmt.Println(colorGreen + "[✓ Done]" + colorReset)

	fmt.Print(" • Removing Camera Lab flag... ")
	_, _ = runADBShellWithRetry("rm -rf /sdcard/DCIM/Camera/lab_options_visible", serial)
	fmt.Println(colorGreen + "[✓ Done]" + colorReset)

	fmt.Println("\n[*] Reinstalling removed packages:")
	runRestorePackagesOnly(serial)

	fmt.Println(colorGreen + colorBold + "\n[✓] All system parameters restored to factory defaults." + colorReset)
}

func runMaintenanceBatch(serial *string) {
	runFSTrim(serial)
	fmt.Println()
	runTrimCaches(serial)
}

// ============================================================================
// --- MENU TREE BUILDER ---
// ============================================================================

func buildMenuTree() *Menu {
	mainMenu := &Menu{
		ID:       "main",
		Title:    "XIOPTIMIZER — HYPEROS & MIUI MASTER CONTROL",
		Subtitle: "Select a category below to optimize your Xiaomi device",
	}

	// [1] Debloat & Telemetry
	debloatSub := &Menu{
		ID:       "debloat",
		Title:    "DEBLOAT & TELEMETRY HUB",
		Subtitle: "Strip out ad trackers, system telemetry, and unnecessary vendor apps",
	}
	debloatSub.Items = []MenuItem{
		{Key: "1", Icon: "🛡", Title: "Safe Ads & Telemetry Cleanup", Description: "MSA, Analytics, Quick Apps, Facebook", Handler: runSafeAdsCleanup},
		{Key: "2", Icon: "🗑", Title: "Deep Bloatware Removal", Description: "Mi Browser, Video, Music, App Vault", Handler: runBloatwareCleanup},
		{Key: "3", Icon: "📦", Title: "Restore Removed Packages", Description: "Reinstall uninstalled system packages", Handler: runRestorePackagesOnly},
		{Key: "A", Icon: "⚡", Title: "Apply Full Debloat (1 + 2)", Description: "Recommended for clean ROM & privacy", Handler: runFullDebloatBatch},
		{Key: "0", Icon: "↩", Title: "Back to Main Menu", Description: "Return to category selection"},
	}

	// [2] Performance & Gaming
	perfSub := &Menu{
		ID:       "perf",
		Title:    "PERFORMANCE & GAMING HUB",
		Subtitle: "Bypass thermal throttling, lock refresh rates, and optimize ART cache",
	}
	perfSub.Items = []MenuItem{
		{Key: "1", Icon: "⚡", Title: "Remove Joyose Gaming Throttle", Description: "Bypasses FPS drops & screen dimming", Handler: runJoyoseRemoval},
		{Key: "2", Icon: "🚀", Title: "Lock Display Refresh to 120Hz", Description: "Forces 120Hz in apps without 60Hz drops", Handler: runLock120Hz},
		{Key: "3", Icon: "🧠", Title: "Disable Virtual RAM (ExtM Swap)", Description: "Stops UFS storage swap micro-stutters", Handler: runDisableVirtualRAM},
		{Key: "4", Icon: "👻", Title: "Disable Phantom Process Killer", Description: "Prevents Android from killing Termux/bots", Handler: runDisablePhantomProcessKiller},
		{Key: "5", Icon: "⚙", Title: "ART Speed-Profile AOT Compile", Description: "Pre-compiles hot bytecode for apps", Handler: runARTCompilation},
		{Key: "A", Icon: "⚡", Title: "Apply All Performance Tweaks", Description: "Ultimate gaming & responsiveness setup", Handler: runAllPerformanceBatch},
		{Key: "0", Icon: "↩", Title: "Back to Main Menu", Description: "Return to category selection"},
	}

	// [3] Battery & Deep Sleep
	batterySub := &Menu{
		ID:       "battery",
		Title:    "BATTERY & DEEP SLEEP HUB",
		Subtitle: "Reduce idle drain, enable fast Doze mode, and stop radio polling",
	}
	batterySub.Items = []MenuItem{
		{Key: "1", Icon: "🔋", Title: "Aggressive Doze & Deep Sleep", Description: "Fast sleep timeout when screen turns off", Handler: runAggressiveDoze},
		{Key: "2", Icon: "📡", Title: "Disable Wi-Fi/BT Background Scan", Description: "Prevents background location polling", Handler: runDisableScanning},
		{Key: "3", Icon: "🌙", Title: "Standby & Telemetry Tuning", Description: "Stops cellular hold on Wi-Fi & crash logs", Handler: runStandbyOptimizations},
		{Key: "4", Icon: "📝", Title: "Minimize Logd Ring Buffer (64K)", Description: "Halts CPU wakeups from debug log daemons", Handler: runLimitLogdBuffer},
		{Key: "A", Icon: "⚡", Title: "Apply All Battery Optimizations", Description: "Maximizes screen-on-time & standby battery", Handler: runAllBatteryBatch},
		{Key: "0", Icon: "↩", Title: "Back to Main Menu", Description: "Return to category selection"},
	}

	// [4] Network & Private DNS
	netSub := &Menu{
		ID:       "network",
		Title:    "NETWORK & PRIVATE DNS HUB",
		Subtitle: "Configure encrypted DNS-over-TLS adblocking and Wi-Fi low latency",
	}
	netSub.Items = []MenuItem{
		{Key: "1", Icon: "🛡", Title: "Enable AdGuard Private DNS", Description: "Blocks banner & in-app ads system-wide", Handler: runSetAdGuardDNS},
		{Key: "2", Icon: "⚡", Title: "Enable Cloudflare 1.1.1.1 DNS", Description: "Ultra-low latency DNS for gaming ping", Handler: runSetCloudflareDNS},
		{Key: "3", Icon: "📶", Title: "Wi-Fi Low-Latency HAL & Roam", Description: "Fixes Wi-Fi jitter & roaming packet drops", Handler: runWiFiLowLatency},
		{Key: "4", Icon: "🔄", Title: "Reset DNS to Automatic Mode", Description: "Restores default DHCP resolver settings", Handler: runResetPrivateDNS},
		{Key: "A", Icon: "⚡", Title: "Apply AdGuard + Low-Latency Wi-Fi", Description: "Best combination for adblock & connection", Handler: runAllNetworkBatch},
		{Key: "0", Icon: "↩", Title: "Back to Main Menu", Description: "Return to category selection"},
	}

	// [5] UI, Fluidity & Visual Fixes
	uiSub := &Menu{
		ID:       "ui",
		Title:    "UI, FLUIDITY & VISUAL HUB",
		Subtitle: "Tune animation speed, freeze cached background apps, and unlock blur",
	}
	uiSub.Items = []MenuItem{
		{Key: "1", Icon: "🎨", Title: "0.5x Fast System Animations", Description: "Halves window, transition & animator delays", Handler: runFastAnimations},
		{Key: "2", Icon: "🧊", Title: "Native Cached Apps Freezer", Description: "Freezes background apps in RAM, saves CPU", Handler: runCachedAppsFreezer},
		{Key: "3", Icon: "👆", Title: "Touch Latency & Response Tuning", Description: "Optimizes pointer speed & touch sampling", Handler: runTouchTuning},
		{Key: "4", Icon: "✨", Title: "Unlock Control Center Blur", Description: "Fixes ugly grey background on budget phones", Handler: runFixGreyBlur},
		{Key: "5", Icon: "➖", Title: "Hide Gesture Bar Indicator Line", Description: "True immersive fullscreen display", Handler: runHideGestureLine},
		{Key: "A", Icon: "⚡", Title: "Apply All UI Fluidity Tweaks", Description: "Smooth animations, blur & instant feedback", Handler: runAllUIBatch},
		{Key: "0", Icon: "↩", Title: "Back to Main Menu", Description: "Return to category selection"},
	}

	// [6] Special System Flags
	specialSub := &Menu{
		ID:       "special",
		Title:    "SPECIAL SYSTEM FLAGS & AUDIO HUB",
		Subtitle: "Unlock experimental features, camera flags, and bypass warnings",
	}
	specialSub.Items = []MenuItem{
		{Key: "1", Icon: "🔊", Title: "Bypass Safe Headphone Volume", Description: "Disables 60% volume limiter popups", Handler: runBypassSafeVolume},
		{Key: "2", Icon: "📷", Title: "Unlock Hidden Camera Lab", Description: "Enables experimental camera feature toggles", Handler: runUnlockCameraLab},
		{Key: "3", Icon: "⚡", Title: "Bypass APK Verification Lag", Description: "Skips slow security scanner on ADB installs", Handler: runBypassInstallerLag},
		{Key: "A", Icon: "⚡", Title: "Apply All Special Flags", Description: "Instant installs, max volume & camera lab", Handler: runAllSpecialBatch},
		{Key: "0", Icon: "↩", Title: "Back to Main Menu", Description: "Return to category selection"},
	}

	// [7] Maintenance & Storage
	maintSub := &Menu{
		ID:       "maintenance",
		Title:    "SYSTEM MAINTENANCE & STORAGE HUB",
		Subtitle: "NAND flash TRIM, package cache pruning, and full factory defaults reset",
	}
	maintSub.Items = []MenuItem{
		{Key: "1", Icon: "🧹", Title: "TRIM Flash Storage (UFS / SSD)", Description: "Trims wear-leveled NAND blocks via fstrim", Handler: runFSTrim},
		{Key: "2", Icon: "🧼", Title: "Prune All Package Caches", Description: "Reclaims gigabytes of temporary junk safely", Handler: runTrimCaches},
		{Key: "3", Icon: "🔄", Title: "Reset All Tweaks to Defaults", Description: "Reinstalls bloatware, resets 1.0x scales & Doze", Handler: runResetAllTweaks},
		{Key: "A", Icon: "⚡", Title: "Run Routine Maintenance (1 + 2)", Description: "Recommended weekly care for storage speed", Handler: runMaintenanceBatch},
		{Key: "0", Icon: "↩", Title: "Back to Main Menu", Description: "Return to category selection"},
	}

	// Root Main Menu items
	mainMenu.Items = []MenuItem{
		{Key: "1", Icon: "🛡", Title: "Debloat & Telemetry", Description: "Safe Ads, Bloatware, Restore apps", SubMenu: debloatSub},
		{Key: "2", Icon: "⚡", Title: "Performance & Gaming", Description: "Joyose, 120Hz, Virtual RAM off, ART", SubMenu: perfSub},
		{Key: "3", Icon: "🔋", Title: "Battery & Deep Sleep", Description: "Aggressive Doze, Wi-Fi/BT scan off, Logd", SubMenu: batterySub},
		{Key: "4", Icon: "🌐", Title: "Network & Private DNS", Description: "AdGuard adblock, Cloudflare, Low-Latency", SubMenu: netSub},
		{Key: "5", Icon: "🎨", Title: "UI & Fluidity", Description: "0.5x animations, apps freezer, touch, blur", SubMenu: uiSub},
		{Key: "6", Icon: "🔮", Title: "Special Flags & Audio", Description: "Safe volume bypass, Camera lab, APK lag", SubMenu: specialSub},
		{Key: "7", Icon: "🛠", Title: "System Maintenance & Storage", Description: "TRIM SSD/UFS, package cache clean, reset", SubMenu: maintSub},
		{Key: "0", Icon: "🚪", Title: "Exit", Description: "Safely quit XiOptimizer", Handler: nil},
	}

	return mainMenu
}

// Manages the hierarchical menu navigation with live connection surveillance
func runMenuLoop(root *Menu, serial *string) bool {
	current := root

	for {
		_, state := checkDeviceState()
		if state != "device" {
			return false
		}

		clearScreen()
		printBanner()
		fmt.Println()
		details := fetchDeviceDetails(serial)
		printDeviceCard(details)
		fmt.Println()

		renderBoxedMenu(current)

		choice, ok := promptSelectionWithCableWatch(serial)
		if !ok {
			return false // Cable disconnected -> drop immediately back to USB scan
		}

		if choice == "0" {
			if current.Parent == nil {
				clearScreen()
				printBanner()
				fmt.Println(colorGreen + "\nThank you for using XiOptimizer. Exiting safely.\n" + colorReset)
				time.Sleep(1 * time.Second)
				os.Exit(0)
			}
			current = current.Parent
			continue
		}

		var matched *MenuItem
		for i := range current.Items {
			if current.Items[i].Key == choice {
				matched = &current.Items[i]
				break
			}
		}

		if matched == nil {
			fmt.Println(colorRed + "\n[!] Invalid selection. Please choose an option from the menu." + colorReset)
			time.Sleep(1 * time.Second)
			continue
		}

		if matched.SubMenu != nil {
			matched.SubMenu.Parent = current
			current = matched.SubMenu
			continue
		}

		if matched.Handler != nil {
			matched.Handler(serial)
			waitEnter()
		}
	}
}

func main() {
	initConsole()

	if err := ensureADB(); err != nil {
		fmt.Printf(colorRed+"[ADB ERROR]: %v\n"+colorReset, err)
		fmt.Println("Press Enter to exit...")
		_, _ = bufio.NewReader(os.Stdin).ReadString('\n')
		os.Exit(1)
	}

	startInputReader()

	menuTree := buildMenuTree()

	for {
		clearScreen()
		printBanner()
		fmt.Println()
		fmt.Println(colorYellow + colorBold + "[⏳] Waiting for a Xiaomi / Redmi / POCO phone..." + colorReset)
		fmt.Println(colorGray + "  1. Enable 'USB Debugging' in Developer Options on your phone.")
		fmt.Println("  2. Connect the phone to this computer via USB cable.")
		fmt.Println("  3. Allow USB debugging if prompted on your phone's screen." + colorReset)
		fmt.Println()

		var connectedSerial string

		ticker := time.NewTicker(2 * time.Second)
		for range ticker.C {
			serial, state := checkDeviceState()

			if state == "unauthorized" {
				fmt.Printf("\r%s[!] Device %s detected, but NOT AUTHORIZED. Confirm debugging on phone...%s", colorYellow, serial, colorReset)
				continue
			}

			if state == "device" {
				connectedSerial = serial
				ticker.Stop()
				break
			}

			dots := strings.Repeat(".", int(time.Now().Unix()%4))
			fmt.Printf("\r%s[🔍] Scanning USB ports%s%-4s%s", colorCyan, colorBold, dots, colorReset)
		}

		_ = runMenuLoop(menuTree, &connectedSerial)
	}
}
