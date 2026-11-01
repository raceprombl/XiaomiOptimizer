package ui

import (
	"fmt"
	"strings"
	"testing"
)

func TestMenuRowWidths(t *testing.T) {
	const innerWidth = 78

	menus := []struct {
		title string
		sub   string
		items []MenuItem
	}{
		{
			title: "XIOPTIMIZER - HYPEROS & MIUI MASTER CONTROL",
			sub:   "Select a category below to optimize your Xiaomi device",
			items: []MenuItem{
				{Key: "1", Icon: "🔒", Title: "Debloat & Telemetry", Description: "Safe Ads, Bloatware, Restore apps"},
				{Key: "2", Icon: "⚡", Title: "Performance & Gaming", Description: "Joyose, 120Hz, Virtual RAM off, ART"},
				{Key: "3", Icon: "🔋", Title: "Battery & Deep Sleep", Description: "Aggressive Doze, Wi-Fi/BT scan, Logd"},
				{Key: "4", Icon: "🌐", Title: "Network & Private DNS", Description: "AdGuard adblock, Cloudflare, Low-Ping"},
				{Key: "5", Icon: "🎨", Title: "UI & Fluidity", Description: "0.5x animations, apps freezer, touch"},
				{Key: "6", Icon: "🔮", Title: "Special Flags & Audio", Description: "Safe volume bypass, Camera lab, APK"},
				{Key: "7", Icon: "🔧", Title: "Maintenance & Backups", Description: "TRIM storage, cache clean, snapshots"},
				{Key: "0", Icon: "🚪", Title: "Exit", Description: "Safely quit XiOptimizer"},
			},
		},
		{
			title: "DEBLOAT & TELEMETRY HUB",
			sub:   "Strip out ad trackers, system telemetry, and unnecessary vendor apps",
			items: []MenuItem{
				{Key: "1", Icon: "🔒", Title: "Safe Ads & Telemetry Cleanup", Description: "MSA, Analytics, Quick Apps, Facebook"},
				{Key: "2", Icon: "🗑️", Title: "Deep Bloatware Removal", Description: "Mi Browser, Video, Music, App Vault"},
				{Key: "3", Icon: "📦", Title: "Restore Removed Packages", Description: "Reinstall uninstalled system packages"},
				{Key: "A", Icon: "⚡", Title: "Apply Full Debloat (1 + 2)", Description: "Recommended for clean ROM & privacy"},
				{Key: "0", Icon: "🔙", Title: "Back to Main Menu", Description: "Return to category selection"},
			},
		},
		{
			title: "PERFORMANCE & GAMING HUB",
			sub:   "Bypass thermal throttling, lock refresh rates, and optimize ART cache",
			items: []MenuItem{
				{Key: "1", Icon: "⚡", Title: "Remove Joyose Gaming Throttle", Description: "Bypasses FPS drops & screen dimming"},
				{Key: "2", Icon: "🚀", Title: "Lock Display Refresh to 120Hz", Description: "Forces 120Hz in apps without 60Hz drop"},
				{Key: "3", Icon: "🧠", Title: "Disable Virtual RAM (Swap)", Description: "Stops UFS storage swap micro-stutters"},
				{Key: "4", Icon: "👻", Title: "Disable Phantom Process Killer", Description: "Prevents Android from killing Termux"},
				{Key: "5", Icon: "⚙️", Title: "ART Speed-Profile AOT Compile", Description: "Pre-compiles hot bytecode for apps"},
				{Key: "A", Icon: "⚡", Title: "Apply All Performance Tweaks", Description: "Ultimate gaming & responsiveness setup"},
				{Key: "0", Icon: "🔙", Title: "Back to Main Menu", Description: "Return to category selection"},
			},
		},
		{
			title: "BATTERY & DEEP SLEEP HUB",
			sub:   "Reduce idle drain, enable fast Doze mode, and stop radio polling",
			items: []MenuItem{
				{Key: "1", Icon: "🔋", Title: "Aggressive Doze & Deep Sleep", Description: "Fast sleep timeout when screen is off"},
				{Key: "2", Icon: "📡", Title: "Disable Background Wi-Fi/BT", Description: "Prevents background location polling"},
				{Key: "3", Icon: "🌙", Title: "Standby & Telemetry Tuning", Description: "Stops cell hold on Wi-Fi & crash logs"},
				{Key: "4", Icon: "📝", Title: "Minimize Logd Buffer (64K)", Description: "Halts CPU wakeups from debug daemons"},
				{Key: "A", Icon: "⚡", Title: "Apply All Battery Tweaks", Description: "Maximizes screen-on & standby battery"},
				{Key: "0", Icon: "🔙", Title: "Back to Main Menu", Description: "Return to category selection"},
			},
		},
		{
			title: "NETWORK & PRIVATE DNS HUB",
			sub:   "Configure encrypted DNS-over-TLS adblocking and Wi-Fi low latency",
			items: []MenuItem{
				{Key: "1", Icon: "🔒", Title: "Enable AdGuard Private DNS", Description: "Blocks banner & in-app ads system-wide"},
				{Key: "2", Icon: "⚡", Title: "Enable Cloudflare 1.1.1.1 DNS", Description: "Ultra-low latency DNS for gaming ping"},
				{Key: "3", Icon: "📶", Title: "Wi-Fi Low-Latency HAL & Roam", Description: "Fixes Wi-Fi jitter & roaming drops"},
				{Key: "4", Icon: "🔄", Title: "Reset DNS to Automatic Mode", Description: "Restores default DHCP resolver setup"},
				{Key: "5", Icon: "🌐", Title: "Enable Wireless ADB (Wi-Fi)", Description: "Switches to wireless mode on port 5555"},
				{Key: "A", Icon: "⚡", Title: "Apply AdGuard + Low-Latency", Description: "Best combination for adblock & speed"},
				{Key: "0", Icon: "🔙", Title: "Back to Main Menu", Description: "Return to category selection"},
			},
		},
		{
			title: "UI, FLUIDITY & VISUAL HUB",
			sub:   "Tune animation speed, freeze cached background apps, and unlock blur",
			items: []MenuItem{
				{Key: "1", Icon: "🎨", Title: "0.5x Fast System Animations", Description: "Halves window, transition & delays"},
				{Key: "2", Icon: "🧊", Title: "Native Cached Apps Freezer", Description: "Freezes background apps in RAM, saves"},
				{Key: "3", Icon: "👆", Title: "Touch Latency & Response", Description: "Optimizes pointer speed & touch rate"},
				{Key: "4", Icon: "✨", Title: "Unlock Control Center Blur", Description: "Fixes ugly grey background on budget"},
				{Key: "5", Icon: "➖", Title: "Hide Gesture Bar Indicator", Description: "True immersive fullscreen display"},
				{Key: "A", Icon: "⚡", Title: "Apply All UI Fluidity Tweaks", Description: "Smooth animations, blur & quick tap"},
				{Key: "0", Icon: "🔙", Title: "Back to Main Menu", Description: "Return to category selection"},
			},
		},
		{
			title: "SPECIAL SYSTEM FLAGS & AUDIO HUB",
			sub:   "Unlock experimental features, camera flags, and bypass warnings",
			items: []MenuItem{
				{Key: "1", Icon: "🔊", Title: "Bypass Safe Headphone Volume", Description: "Disables 60% volume limiter popups"},
				{Key: "2", Icon: "📷", Title: "Unlock Hidden Camera Lab", Description: "Enables experimental camera toggles"},
				{Key: "3", Icon: "⚡", Title: "Bypass APK Verification Lag", Description: "Skips slow security scanner on install"},
				{Key: "A", Icon: "⚡", Title: "Apply All Special Flags", Description: "Instant installs, max volume & camera"},
				{Key: "0", Icon: "🔙", Title: "Back to Main Menu", Description: "Return to category selection"},
			},
		},
		{
			title: "SYSTEM MAINTENANCE, STORAGE & BACKUPS",
			sub:   "Flash TRIM, cache pruning, device JSON snapshots, and precise restores",
			items: []MenuItem{
				{Key: "1", Icon: "🧹", Title: "TRIM Flash Storage (UFS/SSD)", Description: "Trims wear-leveled NAND blocks fstrim"},
				{Key: "2", Icon: "🧼", Title: "Prune All Package Caches", Description: "Reclaims gigabytes of junk safely"},
				{Key: "3", Icon: "📸", Title: "Create Device State Snapshot", Description: "Saves current settings & apps to JSON"},
				{Key: "4", Icon: "🔄", Title: "Restore From Backup Snapshot", Description: "Restores exact previous settings"},
				{Key: "5", Icon: "⚙️", Title: "Reset All Tweaks to Defaults", Description: "Reinstalls bloat, resets 1.0x scales"},
				{Key: "A", Icon: "⚡", Title: "Run Routine Maintenance (1+2)", Description: "Recommended weekly care for storage"},
				{Key: "0", Icon: "🔙", Title: "Back to Main Menu", Description: "Return to category selection"},
			},
		},
	}

	for _, m := range menus {
		t.Run(m.title, func(t *testing.T) {
			titleLine := fmt.Sprintf(" %s%s%s", ColorBold, m.title, ColorReset)
			titleRow := "║ " + FitWidth(titleLine, innerWidth) + " ║"
			if w := VisualWidth(titleRow); w != 82 {
				t.Errorf("Title row visual width = %d, expected 82. Row: %q", w, titleRow)
			}

			if m.sub != "" {
				subLine := fmt.Sprintf(" %s%s%s", ColorGray, m.sub, ColorReset)
				subRow := "║ " + FitWidth(subLine, innerWidth) + " ║"
				if w := VisualWidth(subRow); w != 82 {
					t.Errorf("Subtitle row visual width = %d, expected 82. Row: %q", w, subRow)
				}
			}

			for _, it := range m.items {
				keyColor := ColorGreen
				if it.Key == "0" {
					keyColor = ColorRed
				} else if it.Key == "A" {
					keyColor = ColorYellow + ColorBold
				}
				keyPart := fmt.Sprintf("%s[%s]%s", keyColor, it.Key, ColorReset)
				iconPart := FitWidth(it.Icon, 2)
				titlePart := FitWidth(it.Title, 30)
				descPart := FitWidth(it.Description, 38)

				row := " " + keyPart + " " + iconPart + " " + titlePart + "  " + ColorGray + descPart + ColorReset
				line := "║ " + FitWidth(row, innerWidth) + " ║"
				if w := VisualWidth(line); w != 82 {
					t.Errorf("Item %s [%s] visual width = %d, expected 82. Line: %q", it.Key, it.Title, w, line)
				}
			}
		})
	}
}

func TestBannerAndCardWidths(t *testing.T) {
	boxWidth := 80
	top := "╔" + strings.Repeat("═", boxWidth) + "╗"
	if w := VisualWidth(top); w != 82 {
		t.Errorf("Banner top visual width = %d, expected 82", w)
	}

	sub1 := Center("✦ HyperOS & MIUI 11 / 12 / 13 / 14 Optimizer ✦", boxWidth)
	line1 := "║" + FitWidth(sub1, boxWidth) + "║"
	if w := VisualWidth(line1); w != 82 {
		t.Errorf("Banner sub1 visual width = %d, expected 82", w)
	}

	sub2 := Center("Xiaomi • Redmi • POCO Debloater & Tweaker", boxWidth)
	line2 := "║" + FitWidth(sub2, boxWidth) + "║"
	if w := VisualWidth(line2); w != 82 {
		t.Errorf("Banner sub2 visual width = %d, expected 82", w)
	}
}
