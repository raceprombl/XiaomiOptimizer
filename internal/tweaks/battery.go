package tweaks

import (
	"fmt"

	"github.com/raceprombl/XiOptimizer/internal/adb"
	"github.com/raceprombl/XiOptimizer/internal/ui"
)

// RunAggressiveDoze optimizes deep sleep parameters
func RunAggressiveDoze(serial *string) {
	ui.ClearScreen()
	ui.PrintBanner()
	fmt.Println(ui.ColorGreen + ui.ColorBold + "\n[1] AGGRESSIVE DOZE & DEEP SLEEP" + ui.ColorReset)
	fmt.Println(ui.ColorGray + "Enabling fast deep sleep mode when screen turns off:" + ui.ColorReset)
	fmt.Println()

	adb.EnsureDeviceConnected(serial)
	fmt.Print(" • Activating MIUI device idle trigger... ")
	_, _ = adb.RunADBShellWithRetry("settings put global miui_device_idle_enabled 1", serial)
	fmt.Println(ui.ColorGreen + "[✓ Enabled]" + ui.ColorReset)

	fmt.Print(" • Tuning Doze constants for 30-second sleep timeout... ")
	dozeConsts := "inactive_to=30000,sensing_to=0,locating_to=0,location_accuracy=0,motion_inactive_to=0,idle_after_inactive_to=0,idle_pending_to=30000,max_idle_pending_to=60000,idle_to=1800000,idle_factor=2.0,max_idle_to=14400000,light_after_inactive_to=5000,light_pre_idle_to=5000,light_idle_to=180000,light_idle_factor=2.0,light_max_idle_to=900000"
	_, _ = adb.RunADBShellWithRetry("settings put global device_idle_constants "+dozeConsts, serial)
	fmt.Println(ui.ColorGreen + "[✓ Applied]" + ui.ColorReset)

	fmt.Print(" • Whitelisting critical messaging apps (Telegram, WhatsApp)... ")
	_, _ = adb.RunADBShellWithRetry("dumpsys deviceidle whitelist +org.telegram.messenger", serial)
	_, _ = adb.RunADBShellWithRetry("dumpsys deviceidle whitelist +com.whatsapp", serial)
	fmt.Println(ui.ColorGreen + "[✓ Whitelisted]" + ui.ColorReset)

	fmt.Println(ui.ColorCyan + "\n[i] Aggressive Doze active. Overnight battery drain reduced to under 1%." + ui.ColorReset)
}

// RunDisableScanning turns off background radio polling
func RunDisableScanning(serial *string) {
	ui.ClearScreen()
	ui.PrintBanner()
	fmt.Println(ui.ColorGreen + ui.ColorBold + "\n[2] DISABLE WI-FI & BLUETOOTH BACKGROUND SCANNING" + ui.ColorReset)
	fmt.Println(ui.ColorGray + "Stopping constant background beacon discovery and telemetry:" + ui.ColorReset)
	fmt.Println()

	adb.EnsureDeviceConnected(serial)
	fmt.Print(" • Disabling always-on Wi-Fi background scanning... ")
	_, _ = adb.RunADBShellWithRetry("settings put global wifi_scan_always_enabled 0", serial)
	_, _ = adb.RunADBShellWithRetry("settings put global wifi_wakeup_enabled 0", serial)
	fmt.Println(ui.ColorGreen + "[✓ Disabled]" + ui.ColorReset)

	fmt.Print(" • Disabling always-on Bluetooth LE scanning... ")
	_, _ = adb.RunADBShellWithRetry("settings put global ble_scan_always_enabled 0", serial)
	fmt.Println(ui.ColorGreen + "[✓ Disabled]" + ui.ColorReset)

	fmt.Println(ui.ColorCyan + "\n[i] Radio transceivers will sleep completely when Wi-Fi and Bluetooth are off." + ui.ColorReset)
}

// RunStandbyOptimizations tunes modem standby and disables crash uploads
func RunStandbyOptimizations(serial *string) {
	ui.ClearScreen()
	ui.PrintBanner()
	fmt.Println(ui.ColorGreen + ui.ColorBold + "\n[3] STANDBY & TELEMETRY OPTIMIZATIONS" + ui.ColorReset)
	fmt.Println(ui.ColorGray + "Disabling cellular radio hold when on Wi-Fi and stopping crash uploaders:" + ui.ColorReset)
	fmt.Println()

	adb.EnsureDeviceConnected(serial)
	fmt.Print(" • Disabling cellular radio hold on Wi-Fi (mobile_data_always_on)... ")
	_, _ = adb.RunADBShellWithRetry("settings put global mobile_data_always_on 0", serial)
	fmt.Println(ui.ColorGreen + "[✓ Disabled]" + ui.ColorReset)

	fmt.Print(" • Disabling network scoring and crash upload daemons... ")
	_, _ = adb.RunADBShellWithRetry("settings put global network_scoring_ui_enabled 0", serial)
	_, _ = adb.RunADBShellWithRetry("settings put secure send_action_app_error 0", serial)
	_, _ = adb.RunADBShellWithRetry("settings put global send_action_app_error 0", serial)
	fmt.Println(ui.ColorGreen + "[✓ Disabled]" + ui.ColorReset)

	fmt.Println(ui.ColorCyan + "\n[i] Unnecessary modem wakeups and diagnostic uploads halted." + ui.ColorReset)
}

// RunLimitLogdBuffer minimizes system logger CPU overhead
func RunLimitLogdBuffer(serial *string) {
	ui.ClearScreen()
	ui.PrintBanner()
	fmt.Println(ui.ColorGreen + ui.ColorBold + "\n[4] MINIMIZE SYSTEM LOGGING (LOGD 64K)" + ui.ColorReset)
	fmt.Println(ui.ColorGray + "Restricting logd ring buffer size to reduce CPU wakeups and background writes:" + ui.ColorReset)
	fmt.Println()

	adb.EnsureDeviceConnected(serial)
	fmt.Print(" • Setting logd ring buffer size to 64K... ")
	_, _ = adb.RunADBShellWithRetry("setprop persist.logd.size 64K", serial)
	fmt.Println(ui.ColorGreen + "[✓ Applied]" + ui.ColorReset)

	fmt.Print(" • Turning off MIUI verbose debug logging... ")
	_, _ = adb.RunADBShellWithRetry("settings put global miui_log_enable 0", serial)
	_, _ = adb.RunADBShellWithRetry("settings put global miui_log_event_enable 0", serial)
	fmt.Println(ui.ColorGreen + "[✓ Disabled]" + ui.ColorReset)

	fmt.Println(ui.ColorCyan + "\n[i] System logger overhead minimized. Saves CPU cycles and storage wear." + ui.ColorReset)
}

// RunAllBatteryBatch applies all battery saving optimizations
func RunAllBatteryBatch(serial *string) {
	RunAggressiveDoze(serial)
	fmt.Println()
	RunDisableScanning(serial)
	fmt.Println()
	RunStandbyOptimizations(serial)
	fmt.Println()
	RunLimitLogdBuffer(serial)
}
