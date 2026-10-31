package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/raceprombl/XiOptimizer/internal/adb"
	"github.com/raceprombl/XiOptimizer/internal/tweaks"
	"github.com/raceprombl/XiOptimizer/internal/ui"
)

// initConsole configures terminal code page for UTF-8 on Windows
func initConsole() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "chcp 65001 > nul")
		_ = cmd.Run()
	}
}

// buildMenuTree constructs the complete interactive menu hierarchy
func buildMenuTree() *ui.Menu {
	mainMenu := &ui.Menu{
		ID:       "main",
		Title:    "XIOPTIMIZER — HYPEROS & MIUI MASTER CONTROL",
		Subtitle: "Select a category below to optimize your Xiaomi device",
	}

	// [1] Debloat & Telemetry
	debloatSub := &ui.Menu{
		ID:       "debloat",
		Title:    "DEBLOAT & TELEMETRY HUB",
		Subtitle: "Strip out ad trackers, system telemetry, and unnecessary vendor apps",
	}
	debloatSub.Items = []ui.MenuItem{
		{Key: "1", Icon: "🛡", Title: "Safe Ads & Telemetry Cleanup", Description: "MSA, Analytics, Quick Apps, Facebook", Handler: tweaks.RunSafeAdsCleanup},
		{Key: "2", Icon: "🗑", Title: "Deep Bloatware Removal", Description: "Mi Browser, Video, Music, App Vault", Handler: tweaks.RunBloatwareCleanup},
		{Key: "3", Icon: "📦", Title: "Restore Removed Packages", Description: "Reinstall uninstalled system packages", Handler: tweaks.RunRestorePackagesOnly},
		{Key: "A", Icon: "⚡", Title: "Apply Full Debloat (1 + 2)", Description: "Recommended for clean ROM & privacy", Handler: tweaks.RunFullDebloatBatch},
		{Key: "0", Icon: "↩", Title: "Back to Main Menu", Description: "Return to category selection"},
	}

	// [2] Performance & Gaming
	perfSub := &ui.Menu{
		ID:       "perf",
		Title:    "PERFORMANCE & GAMING HUB",
		Subtitle: "Bypass thermal throttling, lock refresh rates, and optimize ART cache",
	}
	perfSub.Items = []ui.MenuItem{
		{Key: "1", Icon: "⚡", Title: "Remove Joyose Gaming Throttle", Description: "Bypasses FPS drops & screen dimming", Handler: tweaks.RunJoyoseRemoval},
		{Key: "2", Icon: "🚀", Title: "Lock Display Refresh to 120Hz", Description: "Forces 120Hz in apps without 60Hz drops", Handler: tweaks.RunLock120Hz},
		{Key: "3", Icon: "🧠", Title: "Disable Virtual RAM (ExtM Swap)", Description: "Stops UFS storage swap micro-stutters", Handler: tweaks.RunDisableVirtualRAM},
		{Key: "4", Icon: "👻", Title: "Disable Phantom Process Killer", Description: "Prevents Android from killing Termux/bots", Handler: tweaks.RunDisablePhantomProcessKiller},
		{Key: "5", Icon: "⚙", Title: "ART Speed-Profile AOT Compile", Description: "Pre-compiles hot bytecode for apps", Handler: tweaks.RunARTCompilation},
		{Key: "A", Icon: "⚡", Title: "Apply All Performance Tweaks", Description: "Ultimate gaming & responsiveness setup", Handler: tweaks.RunAllPerformanceBatch},
		{Key: "0", Icon: "↩", Title: "Back to Main Menu", Description: "Return to category selection"},
	}

	// [3] Battery & Deep Sleep
	batterySub := &ui.Menu{
		ID:       "battery",
		Title:    "BATTERY & DEEP SLEEP HUB",
		Subtitle: "Reduce idle drain, enable fast Doze mode, and stop radio polling",
	}
	batterySub.Items = []ui.MenuItem{
		{Key: "1", Icon: "🔋", Title: "Aggressive Doze & Deep Sleep", Description: "Fast sleep timeout when screen turns off", Handler: tweaks.RunAggressiveDoze},
		{Key: "2", Icon: "📡", Title: "Disable Wi-Fi/BT Background Scan", Description: "Prevents background location polling", Handler: tweaks.RunDisableScanning},
		{Key: "3", Icon: "🌙", Title: "Standby & Telemetry Tuning", Description: "Stops cellular hold on Wi-Fi & crash logs", Handler: tweaks.RunStandbyOptimizations},
		{Key: "4", Icon: "📝", Title: "Minimize Logd Ring Buffer (64K)", Description: "Halts CPU wakeups from debug log daemons", Handler: tweaks.RunLimitLogdBuffer},
		{Key: "A", Icon: "⚡", Title: "Apply All Battery Optimizations", Description: "Maximizes screen-on-time & standby battery", Handler: tweaks.RunAllBatteryBatch},
		{Key: "0", Icon: "↩", Title: "Back to Main Menu", Description: "Return to category selection"},
	}

	// [4] Network & Private DNS
	netSub := &ui.Menu{
		ID:       "network",
		Title:    "NETWORK & PRIVATE DNS HUB",
		Subtitle: "Configure encrypted DNS-over-TLS adblocking and Wi-Fi low latency",
	}
	netSub.Items = []ui.MenuItem{
		{Key: "1", Icon: "🛡", Title: "Enable AdGuard Private DNS", Description: "Blocks banner & in-app ads system-wide", Handler: tweaks.RunSetAdGuardDNS},
		{Key: "2", Icon: "⚡", Title: "Enable Cloudflare 1.1.1.1 DNS", Description: "Ultra-low latency DNS for gaming ping", Handler: tweaks.RunSetCloudflareDNS},
		{Key: "3", Icon: "📶", Title: "Wi-Fi Low-Latency HAL & Roam", Description: "Fixes Wi-Fi jitter & roaming packet drops", Handler: tweaks.RunWiFiLowLatency},
		{Key: "4", Icon: "🔄", Title: "Reset DNS to Automatic Mode", Description: "Restores default DHCP resolver settings", Handler: tweaks.RunResetPrivateDNS},
		{Key: "A", Icon: "⚡", Title: "Apply AdGuard + Low-Latency Wi-Fi", Description: "Best combination for adblock & connection", Handler: tweaks.RunAllNetworkBatch},
		{Key: "0", Icon: "↩", Title: "Back to Main Menu", Description: "Return to category selection"},
	}

	// [5] UI, Fluidity & Visual Fixes
	uiSub := &ui.Menu{
		ID:       "ui",
		Title:    "UI, FLUIDITY & VISUAL HUB",
		Subtitle: "Tune animation speed, freeze cached background apps, and unlock blur",
	}
	uiSub.Items = []ui.MenuItem{
		{Key: "1", Icon: "🎨", Title: "0.5x Fast System Animations", Description: "Halves window, transition & animator delays", Handler: tweaks.RunFastAnimations},
		{Key: "2", Icon: "🧊", Title: "Native Cached Apps Freezer", Description: "Freezes background apps in RAM, saves CPU", Handler: tweaks.RunCachedAppsFreezer},
		{Key: "3", Icon: "👆", Title: "Touch Latency & Response Tuning", Description: "Optimizes pointer speed & touch sampling", Handler: tweaks.RunTouchTuning},
		{Key: "4", Icon: "✨", Title: "Unlock Control Center Blur", Description: "Fixes ugly grey background on budget phones", Handler: tweaks.RunFixGreyBlur},
		{Key: "5", Icon: "➖", Title: "Hide Gesture Bar Indicator Line", Description: "True immersive fullscreen display", Handler: tweaks.RunHideGestureLine},
		{Key: "A", Icon: "⚡", Title: "Apply All UI Fluidity Tweaks", Description: "Smooth animations, blur & instant feedback", Handler: tweaks.RunAllUIBatch},
		{Key: "0", Icon: "↩", Title: "Back to Main Menu", Description: "Return to category selection"},
	}

	// [6] Special System Flags
	specialSub := &ui.Menu{
		ID:       "special",
		Title:    "SPECIAL SYSTEM FLAGS & AUDIO HUB",
		Subtitle: "Unlock experimental features, camera flags, and bypass warnings",
	}
	specialSub.Items = []ui.MenuItem{
		{Key: "1", Icon: "🔊", Title: "Bypass Safe Headphone Volume", Description: "Disables 60% volume limiter popups", Handler: tweaks.RunBypassSafeVolume},
		{Key: "2", Icon: "📷", Title: "Unlock Hidden Camera Lab", Description: "Enables experimental camera feature toggles", Handler: tweaks.RunUnlockCameraLab},
		{Key: "3", Icon: "⚡", Title: "Bypass APK Verification Lag", Description: "Skips slow security scanner on ADB installs", Handler: tweaks.RunBypassInstallerLag},
		{Key: "A", Icon: "⚡", Title: "Apply All Special Flags", Description: "Instant installs, max volume & camera lab", Handler: tweaks.RunAllSpecialBatch},
		{Key: "0", Icon: "↩", Title: "Back to Main Menu", Description: "Return to category selection"},
	}

	// [7] Maintenance, Storage & State Snapshots
	maintSub := &ui.Menu{
		ID:       "maintenance",
		Title:    "SYSTEM MAINTENANCE, STORAGE & BACKUPS",
		Subtitle: "Flash TRIM, cache pruning, device JSON snapshots, and precise restores",
	}
	maintSub.Items = []ui.MenuItem{
		{Key: "1", Icon: "🧹", Title: "TRIM Flash Storage (UFS / SSD)", Description: "Trims wear-leveled NAND blocks via fstrim", Handler: tweaks.RunFSTrim},
		{Key: "2", Icon: "🧼", Title: "Prune All Package Caches", Description: "Reclaims gigabytes of temporary junk safely", Handler: tweaks.RunTrimCaches},
		{Key: "3", Icon: "📸", Title: "Create Device State Snapshot", Description: "Saves current settings and packages to JSON", Handler: tweaks.RunCreateBackupSnapshot},
		{Key: "4", Icon: "🔄", Title: "Restore From Backup Snapshot", Description: "Restores exact previous settings from file", Handler: tweaks.RunRestoreFromBackupSnapshot},
		{Key: "5", Icon: "⚙", Title: "Reset All Tweaks to Defaults", Description: "Reinstalls bloatware, resets 1.0x scales & Doze", Handler: tweaks.RunResetAllTweaks},
		{Key: "A", Icon: "⚡", Title: "Run Routine Maintenance (1 + 2)", Description: "Recommended weekly care for storage speed", Handler: tweaks.RunMaintenanceBatch},
		{Key: "0", Icon: "↩", Title: "Back to Main Menu", Description: "Return to category selection"},
	}

	// Root Main Menu items
	mainMenu.Items = []ui.MenuItem{
		{Key: "1", Icon: "🛡", Title: "Debloat & Telemetry", Description: "Safe Ads, Bloatware, Restore apps", SubMenu: debloatSub},
		{Key: "2", Icon: "⚡", Title: "Performance & Gaming", Description: "Joyose, 120Hz, Virtual RAM off, ART", SubMenu: perfSub},
		{Key: "3", Icon: "🔋", Title: "Battery & Deep Sleep", Description: "Aggressive Doze, Wi-Fi/BT scan off, Logd", SubMenu: batterySub},
		{Key: "4", Icon: "🌐", Title: "Network & Private DNS", Description: "AdGuard adblock, Cloudflare, Low-Latency", SubMenu: netSub},
		{Key: "5", Icon: "🎨", Title: "UI & Fluidity", Description: "0.5x animations, apps freezer, touch, blur", SubMenu: uiSub},
		{Key: "6", Icon: "🔮", Title: "Special Flags & Audio", Description: "Safe volume bypass, Camera lab, APK lag", SubMenu: specialSub},
		{Key: "7", Icon: "🛠", Title: "Maintenance & Backups", Description: "TRIM storage, cache clean, JSON snapshots", SubMenu: maintSub},
		{Key: "0", Icon: "🚪", Title: "Exit", Description: "Safely quit XiOptimizer", Handler: nil},
	}

	return mainMenu
}

// runMenuLoop manages hierarchical navigation with live USB connection surveillance
func runMenuLoop(root *ui.Menu, serial *string) bool {
	current := root

	for {
		_, state := adb.CheckDeviceState()
		if state != "device" {
			return false
		}

		ui.ClearScreen()
		ui.PrintBanner()
		fmt.Println()
		details := adb.FetchDeviceDetails(serial)
		ui.PrintDeviceCard(details)
		fmt.Println()

		ui.RenderBoxedMenu(current)

		choice, ok := ui.PromptSelectionWithCableWatch(serial)
		if !ok {
			return false // Cable disconnected -> drop immediately back to USB scan
		}

		if choice == "0" {
			if current.Parent == nil {
				ui.ClearScreen()
				ui.PrintBanner()
				fmt.Println(ui.ColorGreen + "\nThank you for using XiOptimizer. Exiting safely.\n" + ui.ColorReset)
				time.Sleep(1 * time.Second)
				os.Exit(0)
			}
			current = current.Parent
			continue
		}

		var matched *ui.MenuItem
		for i := range current.Items {
			if current.Items[i].Key == choice {
				matched = &current.Items[i]
				break
			}
		}

		if matched == nil {
			fmt.Println(ui.ColorRed + "\n[!] Invalid selection. Please choose an option from the menu." + ui.ColorReset)
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
			ui.WaitEnter()
		}
	}
}

func main() {
	initConsole()

	if err := adb.EnsureADB(); err != nil {
		fmt.Printf(ui.ColorRed+"[ADB ERROR]: %v\n"+ui.ColorReset, err)
		fmt.Println("Press Enter to exit...")
		ui.WaitEnter()
		os.Exit(1)
	}

	ui.StartInputReader()

	menuTree := buildMenuTree()

	for {
		ui.ClearScreen()
		ui.PrintBanner()
		fmt.Println()
		fmt.Println(ui.ColorYellow + ui.ColorBold + "[⏳] Waiting for a Xiaomi / Redmi / POCO phone..." + ui.ColorReset)
		fmt.Println(ui.ColorGray + "  1. Enable 'USB Debugging' in Developer Options on your phone.")
		fmt.Println("  2. Connect the phone to this computer via USB cable.")
		fmt.Println("  3. Allow USB debugging if prompted on your phone's screen." + ui.ColorReset)
		fmt.Println()

		var connectedSerial string

		ticker := time.NewTicker(2 * time.Second)
		for range ticker.C {
			serial, state := adb.CheckDeviceState()

			if state == "unauthorized" {
				fmt.Printf("\r%s[!] Device %s detected, but NOT AUTHORIZED. Confirm debugging on phone...%s", ui.ColorYellow, serial, ui.ColorReset)
				continue
			}

			if state == "device" {
				connectedSerial = serial
				ticker.Stop()
				break
			}

			dots := strings.Repeat(".", int(time.Now().Unix()%4))
			fmt.Printf("\r%s[🔍] Scanning USB ports%s%-4s%s", ui.ColorCyan, ui.ColorBold, dots, ui.ColorReset)
		}

		_ = runMenuLoop(menuTree, &connectedSerial)
	}
}
