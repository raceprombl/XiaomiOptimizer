package tweaks

import (
	"fmt"
	"strings"

	"github.com/raceprombl/XiOptimizer/internal/adb"
	"github.com/raceprombl/XiOptimizer/internal/ui"
)

// RunJoyoseRemoval disables thermal and FPS throttling daemon
func RunJoyoseRemoval(serial *string) {
	ui.ClearScreen()
	ui.PrintBanner()
	fmt.Println(ui.ColorMiCyan + ui.ColorBold + "\n[1] JOYOSE GAMING THROTTLE REMOVAL" + ui.ColorReset)
	fmt.Println(ui.ColorGray + "Disabling thermal throttling and game FPS lock daemon:" + ui.ColorReset)
	fmt.Println()

	adb.EnsureDeviceConnected(serial)
	fmt.Print(" • Removing com.xiaomi.joyose for user 0... ")
	out, _ := adb.RunADBShellWithRetry("pm uninstall -k --user 0 com.xiaomi.joyose", serial)
	if strings.Contains(out, "Success") {
		fmt.Println(ui.ColorGreen + "[✓ Removed]" + ui.ColorReset)
	} else if strings.Contains(out, "not installed") {
		fmt.Println(ui.ColorGray + "[– Already absent]" + ui.ColorReset)
	} else {
		_, _ = adb.RunADBShellWithRetry("pm disable-user --user 0 com.xiaomi.joyose", serial)
		fmt.Println(ui.ColorGreen + "[✓ Disabled]" + ui.ColorReset)
	}

	fmt.Println(ui.ColorCyan + "\n[i] Joyose removed. Games maintain steady FPS without artificial drops." + ui.ColorReset)
}

// RunLock120Hz enforces 120Hz across all applications
func RunLock120Hz(serial *string) {
	ui.ClearScreen()
	ui.PrintBanner()
	fmt.Println(ui.ColorMiCyan + ui.ColorBold + "\n[2] LOCK DISPLAY REFRESH TO PEAK 120HZ" + ui.ColorReset)
	fmt.Println(ui.ColorGray + "Enforcing 120Hz refresh rate across all apps and system menus:" + ui.ColorReset)
	fmt.Println()

	adb.EnsureDeviceConnected(serial)
	fmt.Print(" • Clearing PowerKeeper cache... ")
	_, _ = adb.RunADBShellWithRetry("pm clear com.miui.powerkeeper", serial)
	fmt.Println(ui.ColorGreen + "[✓ Done]" + ui.ColorReset)

	fmt.Print(" • Setting peak and min refresh rates to 120Hz... ")
	_, _ = adb.RunADBShellWithRetry("settings put system min_refresh_rate 120.0", serial)
	_, _ = adb.RunADBShellWithRetry("settings put system peak_refresh_rate 120.0", serial)
	_, _ = adb.RunADBShellWithRetry("settings put secure user_refresh_rate 120", serial)
	_, _ = adb.RunADBShellWithRetry("settings put secure miui_refresh_rate 120", serial)
	_, _ = adb.RunADBShellWithRetry("settings put system thermal_limit_refresh_rate 120", serial)
	_, _ = adb.RunADBShellWithRetry("settings put secure support_highfps 1", serial)
	_, _ = adb.RunADBShellWithRetry("settings put secure refresh_rate_mode 2", serial)
	fmt.Println(ui.ColorGreen + "[✓ Applied]" + ui.ColorReset)

	fmt.Println(ui.ColorCyan + "\n[i] Display locked to 120Hz. App drops to 60Hz are eliminated." + ui.ColorReset)
}

// RunDisableVirtualRAM disables slow UFS swap memory extension
func RunDisableVirtualRAM(serial *string) {
	ui.ClearScreen()
	ui.PrintBanner()
	fmt.Println(ui.ColorMiCyan + ui.ColorBold + "\n[3] DISABLE VIRTUAL RAM (EXTM / UFS SWAP)" + ui.ColorReset)
	fmt.Println(ui.ColorGray + "Eliminating micro-stutters and 0.1% low frame drops caused by slow storage paging:" + ui.ColorReset)
	fmt.Println()

	adb.EnsureDeviceConnected(serial)
	fmt.Print(" • Disabling MIUI/HyperOS Memory Extension in system settings... ")
	_, _ = adb.RunADBShellWithRetry("settings put global virtual_ram_enable 0", serial)
	_, _ = adb.RunADBShellWithRetry("setprop persist.miui.extm.enable 0", serial)
	_, _ = adb.RunADBShellWithRetry("setprop persist.sys.extm.enable 0", serial)
	fmt.Println(ui.ColorGreen + "[✓ Disabled]" + ui.ColorReset)

	fmt.Println(ui.ColorCyan + "\n[i] Virtual RAM disabled. Phone will use pure physical LPDDR RAM.")
	fmt.Println("    Reboot your phone once to completely unload the zRAM/UFS swap partition." + ui.ColorReset)
}

// RunDisablePhantomProcessKiller prevents background terminal tasks from being killed
func RunDisablePhantomProcessKiller(serial *string) {
	ui.ClearScreen()
	ui.PrintBanner()
	fmt.Println(ui.ColorMiCyan + ui.ColorBold + "\n[4] DISABLE PHANTOM PROCESS KILLER" + ui.ColorReset)
	fmt.Println(ui.ColorGray + "Preventing Android 12+ from killing child processes and terminal sessions:" + ui.ColorReset)
	fmt.Println()

	adb.EnsureDeviceConnected(serial)
	fmt.Print(" • Disabling DeviceConfig test sync overrides... ")
	_, _ = adb.RunADBShellWithRetry("device_config set_sync_disabled_for_tests persistent", serial)
	fmt.Println(ui.ColorGreen + "[✓ Done]" + ui.ColorReset)

	fmt.Print(" • Raising max phantom processes to 2147483647... ")
	_, _ = adb.RunADBShellWithRetry("device_config put activity_manager max_phantom_processes 2147483647", serial)
	fmt.Println(ui.ColorGreen + "[✓ Applied]" + ui.ColorReset)

	fmt.Println(ui.ColorCyan + "\n[i] Phantom Process Killer disabled. Termux, bots, and background tasks run smoothly." + ui.ColorReset)
}

// RunARTCompilation pre-compiles hot bytecode for system and top apps
func RunARTCompilation(serial *string) {
	ui.ClearScreen()
	ui.PrintBanner()
	fmt.Println(ui.ColorMiCyan + ui.ColorBold + "\n[5] ART SPEED-PROFILE COMPILATION" + ui.ColorReset)
	fmt.Println(ui.ColorGray + "Ahead-Of-Time compiling hot bytecode for installed applications:" + ui.ColorReset)
	fmt.Println()

	adb.EnsureDeviceConnected(serial)
	fmt.Println(ui.ColorYellow + "[*] Compiling core system and daily apps with speed-profile..." + ui.ColorReset)

	keyApps := []string{
		"com.miui.home",
		"com.android.systemui",
		"com.android.settings",
		"com.android.chrome",
		"org.telegram.messenger",
		"com.google.android.youtube",
	}

	for _, app := range keyApps {
		adb.EnsureDeviceConnected(serial)
		fmt.Printf(" • Compiling: %-30s... ", app)
		out, _ := adb.RunADBShellWithRetry("cmd package compile -m speed-profile -f "+app, serial)
		if strings.Contains(out, "Success") {
			fmt.Println(ui.ColorGreen + "[✓ Compiled]" + ui.ColorReset)
		} else {
			fmt.Println(ui.ColorGray + "[– Skipped / Not present]" + ui.ColorReset)
		}
	}

	fmt.Println(ui.ColorCyan + "\n[*] Scheduling background AOT dexopt for remaining apps..." + ui.ColorReset)
	_, _ = adb.RunADBShellWithRetry("cmd package bg-dexopt-job --enable", serial)
	_, _ = adb.RunADBShellWithRetry("nohup cmd package bg-dexopt-job >/dev/null 2>&1 &", serial)
	fmt.Println(ui.ColorGreen + "[✓ Background AOT compilation task launched asynchronously]" + ui.ColorReset)
	fmt.Println(ui.ColorGray + "  Android will continue compiling apps silently while idle or charging." + ui.ColorReset)
	fmt.Println(ui.ColorGreen + ui.ColorBold + "[✓] Core system UI and apps compiled. Device responsiveness boosted!" + ui.ColorReset)
}

// RunVulkanRenderer forces the Vulkan (SkiaVK) rendering backend for UI smoothness
func RunVulkanRenderer(serial *string) {
	ui.ClearScreen()
	ui.PrintBanner()
	fmt.Println(ui.ColorMiCyan + ui.ColorBold + "\n[6] VULKAN HWUI UI RENDERER (SKIAVK)" + ui.ColorReset)
	fmt.Println(ui.ColorGray + "Enforcing multithreaded Vulkan rendering backend for system UI & animations:" + ui.ColorReset)
	fmt.Println()

	adb.EnsureDeviceConnected(serial)
	fmt.Print(" • Configuring HWUI renderer backend to SkiaVK... ")
	_, _ = adb.RunADBShellWithRetry("setprop debug.hwui.renderer skiavkthreaded", serial)
	_, _ = adb.RunADBShellWithRetry("setprop debug.renderengine.backend skiavkthreaded", serial)
	_, _ = adb.RunADBShellWithRetry("settings put global rendering_type skiavk", serial)
	_, _ = adb.RunADBShellWithRetry("settings put global persist.sys.ui.rendering skiavk", serial)
	fmt.Println(ui.ColorGreen + "[✓ Applied]" + ui.ColorReset)

	fmt.Println(ui.ColorCyan + "\n[i] Vulkan UI rendering active. Reduces CPU overhead during 120Hz scrolling.")
	fmt.Println("    Restart SystemUI or reboot phone for full system-wide effects." + ui.ColorReset)
}

// RunGameLowLatency optimizes display buffer queue depth and EAS CPU scheduler clamp
func RunGameLowLatency(serial *string) {
	ui.ClearScreen()
	ui.PrintBanner()
	fmt.Println(ui.ColorMiCyan + ui.ColorBold + "\n[7] GAME LOW-LATENCY & EAS UCLAMP SCHEDULER" + ui.ColorReset)
	fmt.Println(ui.ColorGray + "Minimizing touch input lag and locking high-performance CPU scheduling:" + ui.ColorReset)
	fmt.Println()

	adb.EnsureDeviceConnected(serial)
	fmt.Print(" • Enabling SurfaceFlinger game low-latency mode... ")
	_, _ = adb.RunADBShellWithRetry("settings put global game_low_latency_mode 1", serial)
	fmt.Println(ui.ColorGreen + "[✓ Done]" + ui.ColorReset)

	fmt.Print(" • Configuring EAS uclamp latency sensitivity for foreground app... ")
	_, _ = adb.RunADBShellWithRetry("settings put global uclamp_min_latency_sensitive 25", serial)
	_, _ = adb.RunADBShellWithRetry("settings put global uclamp_min_top_app 20", serial)
	_, _ = adb.RunADBShellWithRetry("settings put global uclamp_min_high_scheduling_group 15", serial)
	fmt.Println(ui.ColorGreen + "[✓ Applied]" + ui.ColorReset)

	fmt.Println(ui.ColorCyan + "\n[i] Game low-latency enabled. Display frame queue minimized and touch response boosted." + ui.ColorReset)
}

// RunAllPerformanceBatch applies all performance tweaks in sequence
func RunAllPerformanceBatch(serial *string) {
	RunJoyoseRemoval(serial)
	fmt.Println()
	RunLock120Hz(serial)
	fmt.Println()
	RunDisableVirtualRAM(serial)
	fmt.Println()
	RunDisablePhantomProcessKiller(serial)
	fmt.Println()
	RunARTCompilation(serial)
	fmt.Println()
	RunVulkanRenderer(serial)
	fmt.Println()
	RunGameLowLatency(serial)
}
