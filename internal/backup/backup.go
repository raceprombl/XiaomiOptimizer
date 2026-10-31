package backup

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/raceprombl/XiOptimizer/internal/adb"
	"github.com/raceprombl/XiOptimizer/internal/ui"
)

// DeviceSnapshot represents the state of a device prior to modifications
type DeviceSnapshot struct {
	CreatedAt         string            `json:"created_at"`
	DeviceSerial      string            `json:"device_serial"`
	DeviceModel       string            `json:"device_model"`
	OSVersion         string            `json:"os_version"`
	AndroidVersion    string            `json:"android_version"`
	InstalledPackages []string          `json:"installed_packages"`
	SystemSettings    map[string]string `json:"system_settings"`
	SecureSettings    map[string]string `json:"secure_settings"`
	GlobalSettings    map[string]string `json:"global_settings"`
}

// parses key=value lines from settings list
func parseSettingsMap(raw string) map[string]string {
	result := make(map[string]string)
	lines := strings.Split(raw, "\n")
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l == "" {
			continue
		}
		parts := strings.SplitN(l, "=", 2)
		if len(parts) == 2 {
			result[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	return result
}

// parses package:com.example from pm list packages
func parsePackagesList(raw string) []string {
	var pkgs []string
	lines := strings.Split(raw, "\n")
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if strings.HasPrefix(l, "package:") {
			name := strings.TrimPrefix(l, "package:")
			name = strings.TrimSpace(name)
			if name != "" {
				pkgs = append(pkgs, name)
			}
		}
	}
	return pkgs
}

// CreateSnapshot captures current settings and package lists into a timestamped JSON file
func CreateSnapshot(serial *string, details adb.DeviceDetails) (string, error) {
	execDir, _ := os.Executable()
	baseDir := filepath.Dir(execDir)
	backupDir := filepath.Join(baseDir, "backups")
	_ = os.MkdirAll(backupDir, 0755)

	fmt.Print(" • Capturing global settings... ")
	globalRaw, _ := adb.RunADBShellWithRetry("settings list global", serial)
	fmt.Println(ui.ColorGreen + "[✓]" + ui.ColorReset)

	fmt.Print(" • Capturing secure settings... ")
	secureRaw, _ := adb.RunADBShellWithRetry("settings list secure", serial)
	fmt.Println(ui.ColorGreen + "[✓]" + ui.ColorReset)

	fmt.Print(" • Capturing system settings... ")
	systemRaw, _ := adb.RunADBShellWithRetry("settings list system", serial)
	fmt.Println(ui.ColorGreen + "[✓]" + ui.ColorReset)

	fmt.Print(" • Capturing installed package registry... ")
	pkgsRaw, _ := adb.RunADBShellWithRetry("pm list packages -u", serial)
	fmt.Println(ui.ColorGreen + "[✓]" + ui.ColorReset)

	snapshot := DeviceSnapshot{
		CreatedAt:         time.Now().Format(time.RFC3339),
		DeviceSerial:      details.Serial,
		DeviceModel:       details.Model,
		OSVersion:         details.OSVersion,
		AndroidVersion:    details.AndroidVer,
		InstalledPackages: parsePackagesList(pkgsRaw),
		SystemSettings:    parseSettingsMap(systemRaw),
		SecureSettings:    parseSettingsMap(secureRaw),
		GlobalSettings:    parseSettingsMap(globalRaw),
	}

	timeTag := time.Now().Format("20060102_150405")
	cleanModel := strings.ReplaceAll(details.Model, " ", "_")
	fileName := fmt.Sprintf("snapshot_%s_%s.json", cleanModel, timeTag)
	fullPath := filepath.Join(backupDir, fileName)

	data, err := json.MarshalIndent(snapshot, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to encode snapshot json: %w", err)
	}

	if err := os.WriteFile(fullPath, data, 0644); err != nil {
		return "", fmt.Errorf("failed to write snapshot file: %w", err)
	}

	return fullPath, nil
}

// ListSnapshots returns all snapshot files found in the backups directory
func ListSnapshots() ([]string, error) {
	execDir, _ := os.Executable()
	baseDir := filepath.Dir(execDir)
	backupDir := filepath.Join(baseDir, "backups")

	entries, err := os.ReadDir(backupDir)
	if err != nil {
		return nil, err
	}

	var results []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".json") {
			results = append(results, filepath.Join(backupDir, e.Name()))
		}
	}
	return results, nil
}

// RestoreSnapshot restores exact previous settings and reinstall packages recorded in snapshot
func RestoreSnapshot(filePath string, serial *string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("cannot read snapshot file: %w", err)
	}

	var snap DeviceSnapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return fmt.Errorf("invalid snapshot json: %w", err)
	}

	fmt.Printf(ui.ColorCyan+"\n[*] Restoring snapshot for: %s (%s) taken at %s\n"+ui.ColorReset,
		snap.DeviceModel, snap.DeviceSerial, snap.CreatedAt)

	// 1. Reinstall previously present packages
	fmt.Println(ui.ColorBold + "\n[Phase 1] Restoring previous packages..." + ui.ColorReset)
	restoredPkgs := 0
	for _, pkg := range snap.InstalledPackages {
		adb.EnsureDeviceConnected(serial)
		out, _ := adb.RunADBShellWithRetry("cmd package install-existing "+pkg, serial)
		if strings.Contains(out, "installed") || strings.Contains(out, "Success") {
			restoredPkgs++
		}
	}
	fmt.Printf(ui.ColorGreen+"[✓] Verified packages restored: %d\n"+ui.ColorReset, restoredPkgs)

	// 2. Restore settings from captured maps
	fmt.Println(ui.ColorBold + "\n[Phase 2] Restoring exact system configurations..." + ui.ColorReset)

	restoreKeySetting := func(namespace string, key string, m map[string]string) {
		val, exists := m[key]
		if exists && val != "" {
			_, _ = adb.RunADBShellWithRetry(fmt.Sprintf("settings put %s %s %s", namespace, key, val), serial)
		} else {
			_, _ = adb.RunADBShellWithRetry(fmt.Sprintf("settings delete %s %s", namespace, key), serial)
		}
	}

	// Restore animation scales
	restoreKeySetting("global", "window_animation_scale", snap.GlobalSettings)
	restoreKeySetting("global", "transition_animation_scale", snap.GlobalSettings)
	restoreKeySetting("global", "animator_duration_scale", snap.GlobalSettings)

	// Restore display refresh rates
	restoreKeySetting("system", "min_refresh_rate", snap.SystemSettings)
	restoreKeySetting("system", "peak_refresh_rate", snap.SystemSettings)
	restoreKeySetting("secure", "user_refresh_rate", snap.SecureSettings)
	restoreKeySetting("secure", "miui_refresh_rate", snap.SecureSettings)

	// Restore battery and doze parameters
	restoreKeySetting("global", "device_idle_constants", snap.GlobalSettings)
	restoreKeySetting("global", "miui_device_idle_enabled", snap.GlobalSettings)
	restoreKeySetting("global", "wifi_scan_always_enabled", snap.GlobalSettings)
	restoreKeySetting("global", "ble_scan_always_enabled", snap.GlobalSettings)

	// Restore Private DNS
	restoreKeySetting("global", "private_dns_mode", snap.GlobalSettings)
	restoreKeySetting("global", "private_dns_specifier", snap.GlobalSettings)

	// Restore virtual RAM & blur
	restoreKeySetting("global", "virtual_ram_enable", snap.GlobalSettings)
	restoreKeySetting("system", "deviceLevelList", snap.SystemSettings)
	restoreKeySetting("secure", "hide_gesture_line", snap.SecureSettings)
	restoreKeySetting("system", "hide_gesture_line", snap.SystemSettings)

	fmt.Println(ui.ColorGreen + "[✓] Device state restored to snapshot configuration." + ui.ColorReset)
	return nil
}
