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
		Title:    "XIOPTIMIZER - HYPEROS & MIUI MASTER CONTROL",
		Subtitle: "Select a category below to optimize your Xiaomi device",
	}

	// [1] Debloat & Telemetry
	debloatSub := &ui.Menu{
		ID:       "debloat",
		Title:    "DEBLOAT & TELEMETRY HUB",
		Subtitle: "Strip out ad trackers, system telemetry, and unnecessary vendor apps",
	}
	debloatSub.Items = []ui.MenuItem{
		{Key: "1", Icon: "🔒", Title: "Safe Ads & Telemetry Cleanup", Description: "MSA, Analytics, Quick Apps, Facebook", Handler: tweaks.RunSafeAdsCleanup},
		{Key: "2", Icon: "🗑️", Title: "Deep Bloatware Removal", Description: "Mi Browser, Video, Music, App Vault", Handler: tweaks.RunBloatwareCleanup},
		{Key: "3", Icon: "📦", Title: "Restore Removed Packages", Description: "Reinstall uninstalled system packages", Handler: tweaks.RunRestorePackagesOnly},
		{Key: "A", Icon: "⚡", Title: "Apply Full Debloat (1 + 2)", Description: "Recommended for clean ROM & privacy", Handler: tweaks.RunFullDebloatBatch},
		{Key: "0", Icon: "🔙", Title: "Back to Main Menu", Description: "Return to category selection"},
	}

	// [2] Performance & Gaming
	perfSub := &ui.Menu{
		ID:       "perf",
		Title:    "PERFORMANCE & GAMING HUB",
		Subtitle: "Bypass thermal throttling, lock refresh rates, and optimize ART cache",
	}
	perfSub.Items = []ui.MenuItem{
		{Key: "1", Icon: "⚡", Title: "Remove Joyose Gaming Throttle", Description: "Bypasses FPS drops & screen dimming", Handler: tweaks.RunJoyoseRemoval},
		{Key: "2", Icon: "🚀", Title: "Lock Display Refresh to 120Hz", Description: "Forces 120Hz in apps without 60Hz drop", Handler: tweaks.RunLock120Hz},
		{Key: "3", Icon: "🧠", Title: "Disable Virtual RAM (Swap)", Description: "Stops UFS storage swap micro-stutters", Handler: tweaks.RunDisableVirtualRAM},
		{Key: "4", Icon: "👻", Title: "Disable Phantom Process Killer", Description: "Prevents Android from killing Termux", Handler: tweaks.RunDisablePhantomProcessKiller},
		{Key: "5", Icon: "⚙️", Title: "ART Speed-Profile AOT Compile", Description: "Pre-compiles hot bytecode for apps", Handler: tweaks.RunARTCompilation},
		{Key: "6", Icon: "🎮", Title: "Vulkan HWUI UI Renderer", Description: "Enables SkiaVK backend for UI", Handler: tweaks.RunVulkanRenderer},
		{Key: "7", Icon: "⚡", Title: "Game Low-Latency & uclamp", Description: "Reduces touch lag & frame queue", Handler: tweaks.RunGameLowLatency},
		{Key: "A", Icon: "⚡", Title: "Apply All Performance Tweaks", Description: "Ultimate gaming & responsiveness setup", Handler: tweaks.RunAllPerformanceBatch},
		{Key: "0", Icon: "🔙", Title: "Back to Main Menu", Description: "Return to category selection"},
	}

	// [3] Battery & Deep Sleep
	batterySub := &ui.Menu{
		ID:       "battery",
		Title:    "BATTERY & DEEP SLEEP HUB",
		Subtitle: "Reduce idle drain, enable fast Doze mode, and stop radio polling",
	}
	batterySub.Items = []ui.MenuItem{
		{Key: "1", Icon: "🔋", Title: "Aggressive Doze & Deep Sleep", Description: "Fast sleep timeout when screen is off", Handler: tweaks.RunAggressiveDoze},
		{Key: "2", Icon: "📡", Title: "Disable Background Wi-Fi/BT", Description: "Prevents background location polling", Handler: tweaks.RunDisableScanning},
		{Key: "3", Icon: "🌙", Title: "Standby & Telemetry Tuning", Description: "Stops cell hold on Wi-Fi & crash logs", Handler: tweaks.RunStandbyOptimizations},
		{Key: "4", Icon: "📝", Title: "Minimize Logd Buffer (64K)", Description: "Halts CPU wakeups from debug daemons", Handler: tweaks.RunLimitLogdBuffer},
		{Key: "A", Icon: "⚡", Title: "Apply All Battery Tweaks", Description: "Maximizes screen-on & standby battery", Handler: tweaks.RunAllBatteryBatch},
		{Key: "0", Icon: "🔙", Title: "Back to Main Menu", Description: "Return to category selection"},
	}

	// [4] Network & Private DNS
	netSub := &ui.Menu{
		ID:       "network",
		Title:    "NETWORK & PRIVATE DNS HUB",
		Subtitle: "Configure encrypted DNS-over-TLS adblocking and Wi-Fi low latency",
	}
	netSub.Items = []ui.MenuItem{
		{Key: "1", Icon: "🔒", Title: "Enable AdGuard Private DNS", Description: "Blocks banner & in-app ads system-wide", Handler: tweaks.RunSetAdGuardDNS},
		{Key: "2", Icon: "⚡", Title: "Enable Cloudflare 1.1.1.1 DNS", Description: "Ultra-low latency DNS for gaming ping", Handler: tweaks.RunSetCloudflareDNS},
		{Key: "3", Icon: "📶", Title: "Wi-Fi Low-Latency HAL & Roam", Description: "Fixes Wi-Fi jitter & roaming drops", Handler: tweaks.RunWiFiLowLatency},
		{Key: "4", Icon: "🔄", Title: "Reset DNS to Automatic Mode", Description: "Restores default DHCP resolver setup", Handler: tweaks.RunResetPrivateDNS},
		{Key: "5", Icon: "🌐", Title: "Enable Wireless ADB (Wi-Fi)", Description: "Switches to wireless mode on port 5555", Handler: tweaks.RunEnableWirelessADB},
		{Key: "A", Icon: "⚡", Title: "Apply AdGuard + Low-Latency", Description: "Best combination for adblock & speed", Handler: tweaks.RunAllNetworkBatch},
		{Key: "0", Icon: "🔙", Title: "Back to Main Menu", Description: "Return to category selection"},
	}

	// [5] UI, Fluidity & Visual Fixes
	uiSub := &ui.Menu{
		ID:       "ui",
		Title:    "UI, FLUIDITY & VISUAL HUB",
		Subtitle: "Tune animation speed, freeze cached background apps, and unlock blur",
	}
	uiSub.Items = []ui.MenuItem{
		{Key: "1", Icon: "🎨", Title: "0.5x Fast System Animations", Description: "Halves window, transition & delays", Handler: tweaks.RunFastAnimations},
		{Key: "2", Icon: "🧊", Title: "Native Freezer & RAM Compact", Description: "Freezes background apps & compacts RAM", Handler: tweaks.RunCachedAppsFreezer},
		{Key: "3", Icon: "👆", Title: "Touch Latency & Response", Description: "Optimizes pointer speed & touch rate", Handler: tweaks.RunTouchTuning},
		{Key: "4", Icon: "✨", Title: "Unlock Control Center Blur", Description: "Fixes ugly grey background on budget", Handler: tweaks.RunFixGreyBlur},
		{Key: "5", Icon: "➖", Title: "Hide Gesture Bar Indicator", Description: "True immersive fullscreen display", Handler: tweaks.RunHideGestureLine},
		{Key: "A", Icon: "⚡", Title: "Apply All UI Fluidity Tweaks", Description: "Smooth animations, blur & quick tap", Handler: tweaks.RunAllUIBatch},
		{Key: "0", Icon: "🔙", Title: "Back to Main Menu", Description: "Return to category selection"},
	}

	// [6] Special System Flags
	specialSub := &ui.Menu{
		ID:       "special",
		Title:    "SPECIAL SYSTEM FLAGS & AUDIO HUB",
		Subtitle: "Unlock experimental features, camera flags, and bypass warnings",
	}
	specialSub.Items = []ui.MenuItem{
		{Key: "1", Icon: "🔊", Title: "Bypass Safe Headphone Volume", Description: "Disables 60% volume limiter popups", Handler: tweaks.RunBypassSafeVolume},
		{Key: "2", Icon: "📷", Title: "Unlock Hidden Camera Lab", Description: "Enables experimental camera toggles", Handler: tweaks.RunUnlockCameraLab},
		{Key: "3", Icon: "⚡", Title: "Bypass APK Verification Lag", Description: "Skips slow security scanner on install", Handler: tweaks.RunBypassInstallerLag},
		{Key: "4", Icon: "🛡️", Title: "AMOLED Burn-in Protection", Description: "Enables pixel shift for OLED status", Handler: tweaks.RunBurnInProtection},
		{Key: "A", Icon: "⚡", Title: "Apply All Special Flags", Description: "Instant installs, max volume & camera", Handler: tweaks.RunAllSpecialBatch},
		{Key: "0", Icon: "🔙", Title: "Back to Main Menu", Description: "Return to category selection"},
	}

	// [7] Maintenance, Storage & State Snapshots
	maintSub := &ui.Menu{
		ID:       "maintenance",
		Title:    "SYSTEM MAINTENANCE, STORAGE & BACKUPS",
		Subtitle: "Flash TRIM, cache pruning, device JSON snapshots, and precise restores",
	}
	maintSub.Items = []ui.MenuItem{
		{Key: "1", Icon: "🧹", Title: "TRIM Flash Storage (UFS/SSD)", Description: "Trims wear-leveled NAND blocks fstrim", Handler: tweaks.RunFSTrim},
		{Key: "2", Icon: "🧼", Title: "Prune All Package Caches", Description: "Reclaims gigabytes of junk safely", Handler: tweaks.RunTrimCaches},
		{Key: "3", Icon: "📸", Title: "Create Device State Snapshot", Description: "Saves current settings & apps to JSON", Handler: tweaks.RunCreateBackupSnapshot},
		{Key: "4", Icon: "🔄", Title: "Restore From Backup Snapshot", Description: "Restores exact previous settings", Handler: tweaks.RunRestoreFromBackupSnapshot},
		{Key: "5", Icon: "⚙️", Title: "Reset All Tweaks to Defaults", Description: "Reinstalls bloat, resets 1.0x scales", Handler: tweaks.RunResetAllTweaks},
		{Key: "A", Icon: "⚡", Title: "Run Routine Maintenance (1+2)", Description: "Recommended weekly care for storage", Handler: tweaks.RunMaintenanceBatch},
		{Key: "0", Icon: "🔙", Title: "Back to Main Menu", Description: "Return to category selection"},
	}

	// Root Main Menu items
	mainMenu.Items = []ui.MenuItem{
		{Key: "1", Icon: "🔒", Title: "Debloat & Telemetry", Description: "Safe Ads, Bloatware, Restore apps", SubMenu: debloatSub},
		{Key: "2", Icon: "⚡", Title: "Performance & Gaming", Description: "Joyose, 120Hz, Vulkan, Low-Lag, ART", SubMenu: perfSub},
		{Key: "3", Icon: "🔋", Title: "Battery & Deep Sleep", Description: "Aggressive Doze, Wi-Fi/BT scan, Logd", SubMenu: batterySub},
		{Key: "4", Icon: "🌐", Title: "Network & Private DNS", Description: "AdGuard adblock, Cloudflare, Low-Ping", SubMenu: netSub},
		{Key: "5", Icon: "🎨", Title: "UI & Fluidity", Description: "0.5x animations, apps freezer, touch", SubMenu: uiSub},
		{Key: "6", Icon: "🔮", Title: "Special Flags & Audio", Description: "Safe volume, Camera lab, Burn-in", SubMenu: specialSub},
		{Key: "7", Icon: "🔧", Title: "Maintenance & Backups", Description: "TRIM storage, cache clean, snapshots", SubMenu: maintSub},
		{Key: "0", Icon: "🚪", Title: "Exit", Description: "Safely quit XiOptimizer", Handler: nil},
	}

	return mainMenu
}

// runWiFiWizard guides user through connecting to Xiaomi phone over Wi-Fi
func runWiFiWizard() {
	for {
		ui.ClearScreen()
		ui.PrintBanner()
		fmt.Println()

		savedList := adb.LoadWiFiHistory()
		lastSaved := ""
		if len(savedList) > 0 {
			lastSaved = savedList[0]
		}

		const innerWidth = 78
		top := ui.ColorMiCyan + "╔" + strings.Repeat("═", 80) + "╗" + ui.ColorReset
		mid := ui.ColorMiCyan + "╠" + strings.Repeat("═", 80) + "╣" + ui.ColorReset
		bot := ui.ColorMiCyan + "╚" + strings.Repeat("═", 80) + "╝" + ui.ColorReset

		fmt.Println(top)
		fmt.Println(ui.ColorMiCyan + "║ " + ui.FitWidth(" 📶  WIRELESS ADB & WI-FI CONNECTION WIZARD", innerWidth) + ui.ColorMiCyan + " ║" + ui.ColorReset)
		fmt.Println(ui.ColorMiCyan + "║ " + ui.FitWidth(" Connect and optimize your Xiaomi / Redmi / POCO device wirelessly", innerWidth) + ui.ColorMiCyan + " ║" + ui.ColorReset)
		fmt.Println(mid)

		type wizardItem struct {
			key   string
			title string
			desc  string
		}

		var items []wizardItem
		if lastSaved != "" {
			items = []wizardItem{
				{"1", "⚡ Reconnect to Saved Device", lastSaved + " (Remembered)"},
				{"2", "🌐 Connect to New IP Address", "e.g. 192.168.1.50 or :5555"},
				{"3", "🔑 Wireless Debugging (Pairing Code)", "Android 11+ / HyperOS / MIUI 12.5+"},
				{"4", "🔄 One-Click Switch from USB to Wi-Fi", "Configures port 5555 on USB phone"},
				{"5", "🗑️ Clear Saved Wi-Fi History", "Forget remembered device IP"},
				{"0", "🔙 Return to Auto-Scan", "Go back to main device detection loop"},
			}
		} else {
			items = []wizardItem{
				{"1", "⚡ Quick Connect by IP", "e.g. 192.168.1.50 or 192.168.1.50:5555"},
				{"2", "🔑 Wireless Debugging (Pairing Code)", "Android 11+ / HyperOS / MIUI 12.5+"},
				{"3", "🔄 One-Click Switch from USB to Wi-Fi", "Configures port 5555 on USB phone"},
				{"0", "🔙 Return to Auto-Scan", "Go back to main device detection loop"},
			}
		}

		for _, it := range items {
			line := fmt.Sprintf(" [%s] %-36s %s", it.key, it.title, it.desc)
			fmt.Println(ui.ColorMiCyan + "║ " + ui.FitWidth(line, innerWidth) + ui.ColorMiCyan + " ║" + ui.ColorReset)
		}
		fmt.Println(bot)
		fmt.Println()
		fmt.Print(ui.ColorBold + "Select an option and press Enter: " + ui.ColorReset)

		ui.DrainInput()
		choice := <-ui.InputChan
		choice = strings.ToUpper(strings.TrimSpace(choice))

		if choice == "0" || choice == "BACK" || choice == "EXIT" {
			return
		}

		// Normalize choice when saved history is present
		action := choice
		if lastSaved != "" {
			switch choice {
			case "1":
				action = "RECONNECT_SAVED"
			case "2":
				action = "CONNECT_IP"
			case "3":
				action = "PAIR_CODE"
			case "4":
				action = "SWITCH_USB"
			case "5":
				action = "CLEAR_HISTORY"
			}
		} else {
			switch choice {
			case "1":
				action = "CONNECT_IP"
			case "2":
				action = "PAIR_CODE"
			case "3":
				action = "SWITCH_USB"
			}
		}

		switch action {
		case "RECONNECT_SAVED":
			ui.ClearScreen()
			ui.PrintBanner()
			fmt.Println(ui.ColorMiCyan + ui.ColorBold + "\n[1] RECONNECT TO REMEMBERED WI-FI TARGET" + ui.ColorReset)
			fmt.Printf(ui.ColorGray+"Connecting to previously saved address: %s%s%s...\n"+ui.ColorReset, ui.ColorBold, lastSaved, ui.ColorReset)
			fmt.Println()
			fmt.Printf("[⏳] Connecting to %s... ", lastSaved)
			out, _ := adb.ConnectWiFi(lastSaved)
			if strings.Contains(out, "connected to") {
				adb.SaveWiFiHistory(lastSaved)
				fmt.Println(ui.ColorGreen + "[✓ Success!]" + ui.ColorReset)
				fmt.Printf(ui.ColorGreen+"Connected to %s over Wi-Fi.\n"+ui.ColorReset, lastSaved)
				time.Sleep(1 * time.Second)
				return
			}
			fmt.Println(ui.ColorRed + "[✗ Failed]" + ui.ColorReset)
			fmt.Printf(ui.ColorYellow+"Output: %s\n"+ui.ColorReset, out)
			fmt.Println(ui.ColorGray + "\nTroubleshooting tips:")
			fmt.Println("  1. Make sure phone and PC are on the SAME Wi-Fi network.")
			fmt.Println("  2. If the phone rebooted, port 5555 might have closed (re-enable via USB or Pairing Code).")
			ui.WaitEnter()

		case "CLEAR_HISTORY":
			adb.ClearWiFiHistory()
			fmt.Println(ui.ColorGreen + "\n[✓] Saved Wi-Fi history cleared." + ui.ColorReset)
			time.Sleep(1 * time.Second)

		case "CONNECT_IP":
			ui.ClearScreen()
			ui.PrintBanner()
			fmt.Println(ui.ColorMiCyan + ui.ColorBold + "\nCONNECT BY IP ADDRESS" + ui.ColorReset)
			fmt.Println(ui.ColorGray + "Enter the IP address of your phone (must be on the same Wi-Fi network):" + ui.ColorReset)
			fmt.Println()
			fmt.Print(ui.ColorBold + "Enter IP (e.g. 192.168.1.50 or 192.168.1.50:5555): " + ui.ColorReset)
			ui.DrainInput()
			target := strings.TrimSpace(<-ui.InputChan)
			if target == "" {
				continue
			}
			if !strings.Contains(target, ":") {
				target = target + ":5555"
			}
			fmt.Printf("\n[⏳] Connecting to %s... ", target)
			out, _ := adb.ConnectWiFi(target)
			if strings.Contains(out, "connected to") {
				adb.SaveWiFiHistory(target)
				fmt.Println(ui.ColorGreen + "[✓ Success!]" + ui.ColorReset)
				fmt.Printf(ui.ColorGreen+"Connected to %s over Wi-Fi.\n"+ui.ColorReset, target)
				time.Sleep(1 * time.Second)
				return
			}
			fmt.Println(ui.ColorRed + "[✗ Failed]" + ui.ColorReset)
			fmt.Printf(ui.ColorYellow+"Output: %s\n"+ui.ColorReset, out)
			fmt.Println(ui.ColorGray + "\nTroubleshooting tips:")
			fmt.Println("  1. Make sure phone and PC are on the SAME Wi-Fi network.")
			fmt.Println("  2. If using Android 11+ Wireless Debugging, use Option (Pairing Code).")
			fmt.Println("  3. If using port 5555, enable it first via USB.")
			ui.WaitEnter()

		case "PAIR_CODE":
			ui.ClearScreen()
			ui.PrintBanner()
			fmt.Println(ui.ColorMiCyan + ui.ColorBold + "\nWIRELESS DEBUGGING WITH PAIRING CODE (ANDROID 11+ / HYPEROS)" + ui.ColorReset)
			fmt.Println(ui.ColorGray + "Instructions on your Xiaomi / Redmi / POCO phone:" + ui.ColorReset)
			fmt.Println("  1. Open Settings -> Developer Options (Для разработчиков).")
			fmt.Println("  2. Turn ON 'Wireless Debugging' (Беспроводная отладка).")
			fmt.Println("  3. Tap 'Pair device with pairing code' (Подключить с помощью кода).")
			fmt.Println("  4. You will see an IP:PairingPort and a 6-digit Wi-Fi Pairing Code.")
			fmt.Println()

			fmt.Print(ui.ColorBold + "Step 1. Enter IP & Pairing Port from phone (e.g. 192.168.1.50:37841): " + ui.ColorReset)
			ui.DrainInput()
			pairTarget := strings.TrimSpace(<-ui.InputChan)
			if pairTarget == "" {
				continue
			}

			fmt.Print(ui.ColorBold + "Step 2. Enter 6-digit Wi-Fi Pairing Code (e.g. 849201): " + ui.ColorReset)
			ui.DrainInput()
			code := strings.TrimSpace(<-ui.InputChan)
			if code == "" {
				continue
			}

			fmt.Printf("\n[⏳] Pairing with %s... ", pairTarget)
			out, _ := adb.PairWiFi(pairTarget, code)
			if !strings.Contains(out, "Successfully paired") && !strings.Contains(out, "already paired") {
				fmt.Println(ui.ColorRed + "[✗ Pairing Failed]" + ui.ColorReset)
				fmt.Printf(ui.ColorYellow+"Output: %s\n"+ui.ColorReset, out)
				ui.WaitEnter()
				continue
			}
			fmt.Println(ui.ColorGreen + "[✓ Successfully Paired!]" + ui.ColorReset)

			ipBase := pairTarget
			if idx := strings.LastIndex(pairTarget, ":"); idx != -1 {
				ipBase = pairTarget[:idx]
			}

			fmt.Println(ui.ColorGray + "\nNow look at the main 'Wireless Debugging' screen on your phone.")
			fmt.Println("Check the 'IP address & Port' displayed at the top." + ui.ColorReset)
			fmt.Print(ui.ColorBold + "Step 3. Enter connection port (or full IP:Port, e.g. 41235): " + ui.ColorReset)
			ui.DrainInput()
			connInput := strings.TrimSpace(<-ui.InputChan)
			if connInput == "" {
				continue
			}

			connectTarget := connInput
			if !strings.Contains(connInput, ":") {
				connectTarget = ipBase + ":" + connInput
			}

			fmt.Printf("\n[⏳] Connecting to %s... ", connectTarget)
			connOut, _ := adb.ConnectWiFi(connectTarget)
			if strings.Contains(connOut, "connected to") {
				adb.SaveWiFiHistory(connectTarget)
				fmt.Println(ui.ColorGreen + "[✓ Connected!]" + ui.ColorReset)
				fmt.Printf(ui.ColorGreen+"Wireless ADB session active on %s\n"+ui.ColorReset, connectTarget)
				time.Sleep(1 * time.Second)
				return
			}
			fmt.Println(ui.ColorRed + "[✗ Connection Failed]" + ui.ColorReset)
			fmt.Printf(ui.ColorYellow+"Output: %s\n"+ui.ColorReset, connOut)
			ui.WaitEnter()

		case "SWITCH_USB":
			ui.ClearScreen()
			ui.PrintBanner()
			fmt.Println(ui.ColorMiCyan + ui.ColorBold + "\nONE-CLICK SWITCH FROM USB TO WI-FI" + ui.ColorReset)
			fmt.Println(ui.ColorGray + "Automatically configures wireless mode while USB cable is plugged in:" + ui.ColorReset)
			fmt.Println()

			ser, st := adb.CheckDeviceState()
			if st != "device" {
				fmt.Println(ui.ColorRed + "[!] No device currently connected via USB cable." + ui.ColorReset)
				fmt.Println(ui.ColorGray + "    Plug in the USB cable once to enable port 5555, then you can unplug it." + ui.ColorReset)
				ui.WaitEnter()
				continue
			}

			fmt.Printf(" • Detected device: %s%s%s\n", ui.ColorBold, ser, ui.ColorReset)
			fmt.Print(" • Detecting phone's local Wi-Fi IP address... ")
			ip := adb.GetDeviceWiFiIP(&ser)
			if ip != "" {
				fmt.Printf(ui.ColorGreen+"[✓ Found: %s]\n"+ui.ColorReset, ip)
			} else {
				fmt.Println(ui.ColorYellow + "[– Could not auto-detect]" + ui.ColorReset)
				fmt.Print(ui.ColorBold + "   Please enter phone's Wi-Fi IP (from Wi-Fi Settings): " + ui.ColorReset)
				ui.DrainInput()
				ip = strings.TrimSpace(<-ui.InputChan)
				if ip == "" {
					continue
				}
			}

			fmt.Print(" • Setting adbd to port 5555... ")
			_, err := adb.EnableTCPMode(5555, &ser)
			if err != nil {
				fmt.Println(ui.ColorRed + "[✗ Failed]" + ui.ColorReset)
				ui.WaitEnter()
				continue
			}
			fmt.Println(ui.ColorGreen + "[✓ Ready]" + ui.ColorReset)

			time.Sleep(1 * time.Second)
			target := ip + ":5555"
			fmt.Printf(" • Connecting to %s... ", target)
			connOut, _ := adb.ConnectWiFi(target)
			if strings.Contains(connOut, "connected to") {
				adb.SaveWiFiHistory(target)
				fmt.Println(ui.ColorGreen + "[✓ Connected!]" + ui.ColorReset)
				fmt.Println(ui.ColorGreen + ui.ColorBold + "\n[🎉] WIRELESS ADB ACTIVATED!" + ui.ColorReset)
				fmt.Println(ui.ColorCyan + "[i] You can now safely UNPLUG the USB cable." + ui.ColorReset)
				time.Sleep(1500 * time.Millisecond)
				return
			}
			fmt.Println(ui.ColorRed + "[✗ Connection Failed]" + ui.ColorReset)
			fmt.Printf(ui.ColorYellow+"Output: %s\n"+ui.ColorReset, connOut)
			ui.WaitEnter()
		}
	}
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
			return false // Cable / Wi-Fi disconnected -> drop immediately back to scan
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
		fmt.Println(ui.ColorGray + "  • USB Mode:   Connect phone via USB cable (USB Debugging enabled).")
		fmt.Println("  • Wi-Fi Mode: Type 'W' and press Enter to open Wireless Setup Wizard." + ui.ColorReset)
		savedWiFiList := adb.LoadWiFiHistory()
		lastWiFi := adb.GetLastWiFiAddress()
		if lastWiFi != "" {
			fmt.Printf(ui.ColorCyan+"  • Auto-Detect: Searching for remembered Wi-Fi target (%s)...\n"+ui.ColorReset, lastWiFi)
		}
		fmt.Println()

		var connectedSerial string

		ticker := time.NewTicker(1500 * time.Millisecond)
		loopCount := 0
		for {
			serial, state := adb.CheckDeviceState()

			if state == "unauthorized" {
				fmt.Printf("\r%s[!] Device %s detected, but NOT AUTHORIZED. Confirm debugging on phone...%s", ui.ColorYellow, serial, ui.ColorReset)
			} else if state == "device" {
				connectedSerial = serial
				ticker.Stop()
				if strings.Contains(connectedSerial, ":") {
					adb.SaveWiFiHistory(connectedSerial)
				}
				break
			} else {
				dots := strings.Repeat(".", int(time.Now().Unix()%4))
				fmt.Printf("\r%s[🔍] Scanning for phone (USB / Wi-Fi)%s%-4s%s [Press 'W' + Enter for Wi-Fi]: ", ui.ColorCyan, ui.ColorBold, dots, ui.ColorReset)
			}

			// Automatically attempt non-blocking connection to remembered Wi-Fi target
			if loopCount%2 == 0 && len(savedWiFiList) > 0 {
				go func(targets []string) {
					for _, addr := range targets {
						out, _ := adb.ConnectWiFi(addr)
						if strings.Contains(out, "connected to") {
							break
						}
					}
				}(savedWiFiList)
			}
			loopCount++

			select {
			case inp := <-ui.InputChan:
				cmd := strings.ToUpper(strings.TrimSpace(inp))
				if cmd == "W" || cmd == "WIFI" {
					ticker.Stop()
					runWiFiWizard()
					ui.ClearScreen()
					ui.PrintBanner()
					fmt.Println()
					fmt.Println(ui.ColorYellow + ui.ColorBold + "[⏳] Waiting for a Xiaomi / Redmi / POCO phone..." + ui.ColorReset)
					fmt.Println(ui.ColorGray + "  • USB Mode:   Connect phone via USB cable (USB Debugging enabled).")
					fmt.Println("  • Wi-Fi Mode: Type 'W' and press Enter to open Wireless Setup Wizard." + ui.ColorReset)
					savedWiFiList = adb.LoadWiFiHistory()
					lastWiFi = adb.GetLastWiFiAddress()
					if lastWiFi != "" {
						fmt.Printf(ui.ColorCyan+"  • Auto-Detect: Searching for remembered Wi-Fi target (%s)...\n"+ui.ColorReset, lastWiFi)
					}
					fmt.Println()
					ticker = time.NewTicker(1500 * time.Millisecond)
				}
			case <-ticker.C:
			}
		}

		brandInfo := adb.VerifyXiaomiBrand(&connectedSerial)
		if !brandInfo.IsXiaomi {
			ui.ClearScreen()
			ui.PrintBanner()
			fmt.Println()

			const innerWidth = 78
			top := ui.ColorRed + "╔" + strings.Repeat("═", 80) + "╗" + ui.ColorReset
			mid := ui.ColorRed + "╠" + strings.Repeat("═", 80) + "╣" + ui.ColorReset
			bot := ui.ColorRed + "╚" + strings.Repeat("═", 80) + "╝" + ui.ColorReset

			fmt.Println(top)
			fmt.Println(ui.ColorRed + "║ " + ui.FitWidth(" ⚠️  NON-XIAOMI / INCOMPATIBLE DEVICE DETECTED", innerWidth) + ui.ColorRed + " ║" + ui.ColorReset)
			fmt.Println(mid)
			l1 := fmt.Sprintf("  Detected Brand:       %s%s%s", ui.ColorBold, brandInfo.BrandName, ui.ColorReset)
			l2 := fmt.Sprintf("  Manufacturer:         %s%s%s", ui.ColorBold, brandInfo.Vendor, ui.ColorReset)
			l3 := " "
			l4 := "  XiOptimizer is engineered specifically for Xiaomi, Redmi, POCO &"
			l5 := "  Black Shark devices running MIUI 11-14 or HyperOS 1.0-3.0."
			l6 := " "
			l7 := "  Applying Xiaomi system debloat packages and tweaks to a different brand"
			l8 := "  (Samsung, Pixel, OnePlus, etc.) will fail or cause system instability!"
			l9 := " "
			l10 := fmt.Sprintf("  %s[1] Disconnect and wait for a Xiaomi phone (Recommended)%s", ui.ColorGreen, ui.ColorReset)
			l11 := fmt.Sprintf("  %s[2] Force override and proceed anyway at your own risk%s", ui.ColorYellow, ui.ColorReset)

			for _, l := range []string{l1, l2, l3, l4, l5, l6, l7, l8, l9, l10, l11} {
				fmt.Println(ui.ColorRed + "║ " + ui.FitWidth(l, innerWidth) + ui.ColorRed + " ║" + ui.ColorReset)
			}
			fmt.Println(bot)
			fmt.Println()
			fmt.Print(ui.ColorBold + "Select an option [1-2] and press Enter: " + ui.ColorReset)

			choice, ok := ui.PromptSelectionWithCableWatch(&connectedSerial)
			if !ok || choice != "2" {
				continue
			}
		}

		_ = runMenuLoop(menuTree, &connectedSerial)
	}
}
