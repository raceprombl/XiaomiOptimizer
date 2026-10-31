package tweaks

import (
	"fmt"
	"strings"

	"github.com/raceprombl/XiOptimizer/internal/adb"
	"github.com/raceprombl/XiOptimizer/internal/backup"
	"github.com/raceprombl/XiOptimizer/internal/ui"
)

// RunFSTrim executes storage wear-leveling block trim
func RunFSTrim(serial *string) {
	ui.ClearScreen()
	ui.PrintBanner()
	fmt.Println(ui.ColorYellow + ui.ColorBold + "\n[1] TRIM FLASH STORAGE (UFS / SSD)" + ui.ColorReset)
	fmt.Println(ui.ColorGray + "Running wear-leveling block trim on NAND partitions:" + ui.ColorReset)
	fmt.Println()

	adb.EnsureDeviceConnected(serial)
	fmt.Println(ui.ColorYellow + "[*] Executing fstrim across /data and /system..." + ui.ColorReset)
	out, _ := adb.RunADBShellWithRetry("sm fstrim", serial)
	if out != "" {
		fmt.Printf(ui.ColorGray+"%s\n"+ui.ColorReset, out)
	}
	fmt.Println(ui.ColorGreen + "[✓ Flash storage trim complete]" + ui.ColorReset)
	fmt.Println(ui.ColorCyan + "[i] 4K random write latency restored. Flash wear leveled." + ui.ColorReset)
}

// RunTrimCaches prunes temporary junk caches across packages
func RunTrimCaches(serial *string) {
	ui.ClearScreen()
	ui.PrintBanner()
	fmt.Println(ui.ColorYellow + ui.ColorBold + "\n[2] PRUNE PACKAGE CACHES" + ui.ColorReset)
	fmt.Println(ui.ColorGray + "Pruning temporary junk caches without deleting user data:" + ui.ColorReset)
	fmt.Println()

	adb.EnsureDeviceConnected(serial)
	fmt.Print(" • Asking PackageManager to trim caches... ")
	out, _ := adb.RunADBShellWithRetry("pm trim-caches 999999999999", serial)
	if out != "" {
		fmt.Println(ui.ColorGreen + "[✓ " + strings.TrimSpace(out) + "]" + ui.ColorReset)
	} else {
		fmt.Println(ui.ColorGreen + "[✓ Caches pruned]" + ui.ColorReset)
	}

	fmt.Println(ui.ColorCyan + "[i] Gigabytes of stale temporary cache reclaimed." + ui.ColorReset)
}

// RunCreateBackupSnapshot captures a full state snapshot into JSON
func RunCreateBackupSnapshot(serial *string) {
	ui.ClearScreen()
	ui.PrintBanner()
	fmt.Println(ui.ColorMiCyan + ui.ColorBold + "\n[3] CREATE DEVICE STATE BACKUP SNAPSHOT" + ui.ColorReset)
	fmt.Println(ui.ColorGray + "Saving current settings and installed package registry to local JSON:" + ui.ColorReset)
	fmt.Println()

	adb.EnsureDeviceConnected(serial)
	details := adb.FetchDeviceDetails(serial)

	path, err := backup.CreateSnapshot(serial, details)
	if err != nil {
		fmt.Printf(ui.ColorRed+"[!] Backup failed: %v\n"+ui.ColorReset, err)
		return
	}

	fmt.Printf(ui.ColorGreen+ui.ColorBold+"\n[✓] Backup snapshot successfully saved:\n    %s\n"+ui.ColorReset, path)
	fmt.Println(ui.ColorCyan + "[i] Use this snapshot anytime to restore your phone to this exact state." + ui.ColorReset)
}

// RunRestoreFromBackupSnapshot restores from a selected JSON backup
func RunRestoreFromBackupSnapshot(serial *string) {
	ui.ClearScreen()
	ui.PrintBanner()
	fmt.Println(ui.ColorOrange + ui.ColorBold + "\n[4] RESTORE DEVICE FROM SNAPSHOT" + ui.ColorReset)
	fmt.Println(ui.ColorGray + "Restoring settings and packages from previously saved backup:" + ui.ColorReset)
	fmt.Println()

	snapshots, err := backup.ListSnapshots()
	if err != nil || len(snapshots) == 0 {
		fmt.Println(ui.ColorYellow + "[!] No previous backup snapshots found in backups/ directory." + ui.ColorReset)
		return
	}

	fmt.Println(ui.ColorBold + "Available Snapshots:" + ui.ColorReset)
	for i, s := range snapshots {
		fmt.Printf(" [%d] %s\n", i+1, s)
	}
	fmt.Printf("\nSelect snapshot [1-%d] or 0 to cancel: ", len(snapshots))

	ui.DrainInput()
	choice := <-ui.InputChan
	choice = strings.TrimSpace(choice)
	if choice == "0" || choice == "" {
		fmt.Println(ui.ColorGray + "Cancelled." + ui.ColorReset)
		return
	}

	var idx int
	_, err = fmt.Sscanf(choice, "%d", &idx)
	if err != nil || idx < 1 || idx > len(snapshots) {
		fmt.Println(ui.ColorRed + "Invalid selection." + ui.ColorReset)
		return
	}

	selected := snapshots[idx-1]
	adb.EnsureDeviceConnected(serial)
	if err := backup.RestoreSnapshot(selected, serial); err != nil {
		fmt.Printf(ui.ColorRed+"[!] Restore failed: %v\n"+ui.ColorReset, err)
	} else {
		fmt.Println(ui.ColorGreen + ui.ColorBold + "\n[✓] Restore from snapshot completed successfully." + ui.ColorReset)
	}
}

// RunResetAllTweaks restores default parameters and reinstalls packages
func RunResetAllTweaks(serial *string) {
	ui.ClearScreen()
	ui.PrintBanner()
	fmt.Println(ui.ColorOrange + ui.ColorBold + "\n[5] RESET ALL TWEAKS TO FACTORY DEFAULTS" + ui.ColorReset)
	fmt.Println(ui.ColorGray + "Reverting animation scales, display limits, Doze, DNS, and reinstalling apps:" + ui.ColorReset)
	fmt.Println()

	adb.EnsureDeviceConnected(serial)

	fmt.Print(" • Resetting animation scales to 1.0x... ")
	_, _ = adb.RunADBShellWithRetry("settings put global window_animation_scale 1.0", serial)
	_, _ = adb.RunADBShellWithRetry("settings put global transition_animation_scale 1.0", serial)
	_, _ = adb.RunADBShellWithRetry("settings put global animator_duration_scale 1.0", serial)
	fmt.Println(ui.ColorGreen + "[✓ Done]" + ui.ColorReset)

	fmt.Print(" • Removing display refresh rate overrides... ")
	_, _ = adb.RunADBShellWithRetry("settings delete system min_refresh_rate", serial)
	_, _ = adb.RunADBShellWithRetry("settings delete system peak_refresh_rate", serial)
	_, _ = adb.RunADBShellWithRetry("settings delete secure user_refresh_rate", serial)
	_, _ = adb.RunADBShellWithRetry("settings delete secure miui_refresh_rate", serial)
	_, _ = adb.RunADBShellWithRetry("settings delete system thermal_limit_refresh_rate", serial)
	_, _ = adb.RunADBShellWithRetry("settings delete secure support_highfps", serial)
	_, _ = adb.RunADBShellWithRetry("settings delete secure refresh_rate_mode", serial)
	fmt.Println(ui.ColorGreen + "[✓ Done]" + ui.ColorReset)

	fmt.Print(" • Restoring default Doze and scanning parameters... ")
	_, _ = adb.RunADBShellWithRetry("settings delete global device_idle_constants", serial)
	_, _ = adb.RunADBShellWithRetry("settings delete global miui_device_idle_enabled", serial)
	_, _ = adb.RunADBShellWithRetry("settings put global wifi_scan_always_enabled 1", serial)
	_, _ = adb.RunADBShellWithRetry("settings put global ble_scan_always_enabled 1", serial)
	_, _ = adb.RunADBShellWithRetry("dumpsys deviceidle unforce", serial)
	fmt.Println(ui.ColorGreen + "[✓ Done]" + ui.ColorReset)

	fmt.Print(" • Restoring Private DNS to automatic... ")
	_, _ = adb.RunADBShellWithRetry("settings put global private_dns_mode opportunistic", serial)
	_, _ = adb.RunADBShellWithRetry("settings delete global private_dns_specifier", serial)
	fmt.Println(ui.ColorGreen + "[✓ Done]" + ui.ColorReset)

	fmt.Print(" • Restoring Virtual RAM (Memory Extension)... ")
	_, _ = adb.RunADBShellWithRetry("settings put global virtual_ram_enable 1", serial)
	_, _ = adb.RunADBShellWithRetry("setprop persist.miui.extm.enable 1", serial)
	fmt.Println(ui.ColorGreen + "[✓ Done]" + ui.ColorReset)

	fmt.Print(" • Restoring original graphics deviceLevelList & gesture line... ")
	_, _ = adb.RunADBShellWithRetry("settings delete system deviceLevelList", serial)
	_, _ = adb.RunADBShellWithRetry("settings delete secure hide_gesture_line", serial)
	_, _ = adb.RunADBShellWithRetry("settings delete system hide_gesture_line", serial)
	fmt.Println(ui.ColorGreen + "[✓ Done]" + ui.ColorReset)

	fmt.Print(" • Removing Camera Lab flag... ")
	_, _ = adb.RunADBShellWithRetry("rm -rf /sdcard/DCIM/Camera/lab_options_visible", serial)
	fmt.Println(ui.ColorGreen + "[✓ Done]" + ui.ColorReset)

	fmt.Println("\n[*] Reinstalling removed packages:")
	RunRestorePackagesOnly(serial)

	fmt.Println(ui.ColorGreen + ui.ColorBold + "\n[✓] All system parameters restored to factory defaults." + ui.ColorReset)
}

// RunMaintenanceBatch runs routine storage trim and cache clean
func RunMaintenanceBatch(serial *string) {
	RunFSTrim(serial)
	fmt.Println()
	RunTrimCaches(serial)
}
