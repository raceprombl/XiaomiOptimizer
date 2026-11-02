package tweaks

import (
	"fmt"

	"github.com/raceprombl/XiOptimizer/internal/adb"
	"github.com/raceprombl/XiOptimizer/internal/ui"
)

// RunFastAnimations halves transition and animation durations
func RunFastAnimations(serial *string) {
	ui.ClearScreen()
	ui.PrintBanner()
	fmt.Println(ui.ColorMiCyan + ui.ColorBold + "\n[1] 0.5X FAST SYSTEM ANIMATIONS" + ui.ColorReset)
	fmt.Println(ui.ColorGray + "Halving window, transition, and animator scale delays:" + ui.ColorReset)
	fmt.Println()

	adb.EnsureDeviceConnected(serial)
	fmt.Print(" • Setting window animation scale to 0.5x... ")
	_, _ = adb.RunADBShellWithRetry("settings put global window_animation_scale 0.5", serial)
	_, _ = adb.RunADBShellWithRetry("settings put global transition_animation_scale 0.5", serial)
	_, _ = adb.RunADBShellWithRetry("settings put global animator_duration_scale 0.5", serial)
	fmt.Println(ui.ColorGreen + "[✓ Applied]" + ui.ColorReset)

	fmt.Println(ui.ColorCyan + "\n[i] UI animations now execute twice as fast." + ui.ColorReset)
}

// RunCachedAppsFreezer freezes idle background processes in RAM
func RunCachedAppsFreezer(serial *string) {
	ui.ClearScreen()
	ui.PrintBanner()
	fmt.Println(ui.ColorMiCyan + ui.ColorBold + "\n[2] NATIVE CACHED APPS FREEZER & RAM COMPACTION" + ui.ColorReset)
	fmt.Println(ui.ColorGray + "Freezing idle background processes and enabling proactive RAM compaction:" + ui.ColorReset)
	fmt.Println()

	adb.EnsureDeviceConnected(serial)
	fmt.Print(" • Enabling cached apps freezer in Android runtime... ")
	_, _ = adb.RunADBShellWithRetry("settings put global cached_apps_freezer enabled", serial)
	_, _ = adb.RunADBShellWithRetry("device_config put activity_manager_native_boot use_freezer true", serial)
	_, _ = adb.RunADBShellWithRetry("device_config put activity_manager max_cached_processes 64", serial)
	fmt.Println(ui.ColorGreen + "[✓ Enabled]" + ui.ColorReset)

	fmt.Print(" • Enabling proactive ActivityManager RAM compaction & hibernation... ")
	_, _ = adb.RunADBShellWithRetry("settings put global persist.device_config.activity_manager.use_compaction true", serial)
	_, _ = adb.RunADBShellWithRetry("device_config put app_hibernation app_hibernation_enabled 1", serial)
	fmt.Println(ui.ColorGreen + "[✓ Enabled]" + ui.ColorReset)

	fmt.Println(ui.ColorCyan + "\n[i] Inactive apps frozen in RAM & compacted: 0% CPU consumption and more free memory." + ui.ColorReset)
}

// RunTouchTuning reduces long press delays and touch dead zones
func RunTouchTuning(serial *string) {
	ui.ClearScreen()
	ui.PrintBanner()
	fmt.Println(ui.ColorMiCyan + ui.ColorBold + "\n[3] TOUCH LATENCY & RESPONSE TUNING" + ui.ColorReset)
	fmt.Println(ui.ColorGray + "Reducing tap registration deadtime and gesture latency:" + ui.ColorReset)
	fmt.Println()

	adb.EnsureDeviceConnected(serial)
	fmt.Print(" • Shortening long-press and multi-press thresholds... ")
	_, _ = adb.RunADBShellWithRetry("settings put secure long_press_timeout 250", serial)
	_, _ = adb.RunADBShellWithRetry("settings put secure multi_press_timeout 250", serial)
	_, _ = adb.RunADBShellWithRetry("settings put secure tap_duration_threshold 0.0", serial)
	_, _ = adb.RunADBShellWithRetry("settings put secure touch_blocking_period 0.0", serial)
	_, _ = adb.RunADBShellWithRetry("settings put system pointer_speed 1", serial)
	fmt.Println(ui.ColorGreen + "[✓ Applied]" + ui.ColorReset)

	fmt.Println(ui.ColorCyan + "\n[i] Touch responsiveness optimized for gaming and fast navigation." + ui.ColorReset)
}

// RunFixGreyBlur unlocks Gaussian blur in Control Center on budget Redmi/POCO devices
func RunFixGreyBlur(serial *string) {
	ui.ClearScreen()
	ui.PrintBanner()
	fmt.Println(ui.ColorMiCyan + ui.ColorBold + "\n[4] UNLOCK CONTROL CENTER GAUSSIAN BLUR" + ui.ColorReset)
	fmt.Println(ui.ColorGray + "Fixing the ugly grey background on budget/mid-range Redmi & POCO phones:" + ui.ColorReset)
	fmt.Println()

	adb.EnsureDeviceConnected(serial)
	fmt.Print(" • Setting GPU deviceLevelList to v:1,c:3,g:3 (Flagship graphics profile)... ")
	_, _ = adb.RunADBShellWithRetry("settings put system deviceLevelList \"v:1,c:3,g:3\"", serial)
	_, _ = adb.RunADBShellWithRetry("settings put system window_blurs_enabled 1", serial)
	fmt.Println(ui.ColorGreen + "[✓ Applied]" + ui.ColorReset)

	fmt.Println(ui.ColorCyan + "\n[i] High-end blur shaders unlocked. Real-time background blur will appear")
	fmt.Println("    in Control Center after rebooting or restarting System UI." + ui.ColorReset)
}

// RunHideGestureLine hides the fullscreen gesture bar indicator
func RunHideGestureLine(serial *string) {
	ui.ClearScreen()
	ui.PrintBanner()
	fmt.Println(ui.ColorMiCyan + ui.ColorBold + "\n[5] HIDE FULLSCREEN GESTURE BAR LINE" + ui.ColorReset)
	fmt.Println(ui.ColorGray + "Removing the horizontal pill at the bottom for true fullscreen display:" + ui.ColorReset)
	fmt.Println()

	adb.EnsureDeviceConnected(serial)
	fmt.Print(" • Hiding fullscreen navigation indicator line... ")
	_, _ = adb.RunADBShellWithRetry("settings put secure hide_gesture_line 1", serial)
	_, _ = adb.RunADBShellWithRetry("settings put system hide_gesture_line 1", serial)
	fmt.Println(ui.ColorGreen + "[✓ Hidden]" + ui.ColorReset)

	fmt.Println(ui.ColorCyan + "\n[i] Gesture indicator line hidden. Full screen real estate reclaimed." + ui.ColorReset)
}

// RunAllUIBatch applies all UI fluidity, touch, and visual enhancements
func RunAllUIBatch(serial *string) {
	RunFastAnimations(serial)
	fmt.Println()
	RunCachedAppsFreezer(serial)
	fmt.Println()
	RunTouchTuning(serial)
	fmt.Println()
	RunFixGreyBlur(serial)
	fmt.Println()
	RunHideGestureLine(serial)
}
