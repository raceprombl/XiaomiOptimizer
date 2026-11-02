package tweaks

import (
	"fmt"
	"strings"
	"time"

	"github.com/raceprombl/XiOptimizer/internal/adb"
	"github.com/raceprombl/XiOptimizer/internal/ui"
)

// RunSetAdGuardDNS configures system-wide AdGuard DNS over TLS
func RunSetAdGuardDNS(serial *string) {
	ui.ClearScreen()
	ui.PrintBanner()
	fmt.Println(ui.ColorCyan + ui.ColorBold + "\n[1] ENABLE ADGUARD PRIVATE DNS (ADBLOCK)" + ui.ColorReset)
	fmt.Println(ui.ColorGray + "Enabling encrypted DNS over TLS to block in-app and browser ads system-wide:" + ui.ColorReset)
	fmt.Println()

	adb.EnsureDeviceConnected(serial)
	fmt.Print(" • Setting Private DNS provider to dns.adguard-dns.com... ")
	_, _ = adb.RunADBShellWithRetry("settings put global private_dns_mode hostname", serial)
	_, _ = adb.RunADBShellWithRetry("settings put global private_dns_specifier dns.adguard-dns.com", serial)
	fmt.Println(ui.ColorGreen + "[✓ Enabled]" + ui.ColorReset)

	fmt.Println(ui.ColorCyan + "\n[i] AdGuard Private DNS is active. In-app popups and banner ads are blocked." + ui.ColorReset)
}

// RunSetCloudflareDNS configures low-latency 1.1.1.1 DNS over TLS
func RunSetCloudflareDNS(serial *string) {
	ui.ClearScreen()
	ui.PrintBanner()
	fmt.Println(ui.ColorCyan + ui.ColorBold + "\n[2] ENABLE CLOUDFLARE 1.1.1.1 PRIVATE DNS" + ui.ColorReset)
	fmt.Println(ui.ColorGray + "Enabling ultra-low latency DNS over TLS for competitive online gaming:" + ui.ColorReset)
	fmt.Println()

	adb.EnsureDeviceConnected(serial)
	fmt.Print(" • Setting Private DNS provider to one.one.one.one... ")
	_, _ = adb.RunADBShellWithRetry("settings put global private_dns_mode hostname", serial)
	_, _ = adb.RunADBShellWithRetry("settings put global private_dns_specifier one.one.one.one", serial)
	fmt.Println(ui.ColorGreen + "[✓ Enabled]" + ui.ColorReset)

	fmt.Println(ui.ColorCyan + "\n[i] Cloudflare 1.1.1.1 active. Faster domain resolution and lower game latency." + ui.ColorReset)
}

// RunResetPrivateDNS restores automatic opportunistic DNS mode
func RunResetPrivateDNS(serial *string) {
	ui.ClearScreen()
	ui.PrintBanner()
	fmt.Println(ui.ColorCyan + ui.ColorBold + "\n[3] RESET PRIVATE DNS TO AUTOMATIC" + ui.ColorReset)
	fmt.Println(ui.ColorGray + "Restoring default ISP/network DNS provider:" + ui.ColorReset)
	fmt.Println()

	adb.EnsureDeviceConnected(serial)
	fmt.Print(" • Restoring opportunistic DNS mode... ")
	_, _ = adb.RunADBShellWithRetry("settings put global private_dns_mode opportunistic", serial)
	_, _ = adb.RunADBShellWithRetry("settings delete global private_dns_specifier", serial)
	fmt.Println(ui.ColorGreen + "[✓ Restored]" + ui.ColorReset)

	fmt.Println(ui.ColorCyan + "\n[i] Private DNS returned to standard automatic DHCP configuration." + ui.ColorReset)
}

// RunWiFiLowLatency enforces low-latency Wi-Fi HAL and fixes roaming jitter
func RunWiFiLowLatency(serial *string) {
	ui.ClearScreen()
	ui.PrintBanner()
	fmt.Println(ui.ColorCyan + ui.ColorBold + "\n[4] WI-FI LOW-LATENCY HAL & ROAMING FIX" + ui.ColorReset)
	fmt.Println(ui.ColorGray + "Enforcing low latency Wi-Fi mode and eliminating packet jitter:" + ui.ColorReset)
	fmt.Println()

	adb.EnsureDeviceConnected(serial)
	fmt.Print(" • Forcing Wi-Fi HAL low latency mode... ")
	_, _ = adb.RunADBShellWithRetry("cmd wifi force-low-latency-mode enabled", serial)
	_, _ = adb.RunADBShellWithRetry("settings put global wifi_suspend_optimizations_enabled 0", serial)
	fmt.Println(ui.ColorGreen + "[✓ Enabled]" + ui.ColorReset)

	fmt.Print(" • Disabling Xiaomi aggressive network switching... ")
	_, _ = adb.RunADBShellWithRetry("settings put global cloud_weak_network_switch_enabled off", serial)
	fmt.Println(ui.ColorGreen + "[✓ Done]" + ui.ColorReset)

	fmt.Println(ui.ColorCyan + "\n[i] Wi-Fi latency stabilized. Packet drops and jitter during gaming reduced." + ui.ColorReset)
}

// RunEnableWirelessADB switches current USB session to wireless ADB over Wi-Fi
func RunEnableWirelessADB(serial *string) {
	ui.ClearScreen()
	ui.PrintBanner()
	fmt.Println(ui.ColorCyan + ui.ColorBold + "\n[5] ENABLE WIRELESS ADB OVER WI-FI" + ui.ColorReset)
	fmt.Println(ui.ColorGray + "Configure device for wireless ADB and disconnect USB cable:" + ui.ColorReset)
	fmt.Println()

	adb.EnsureDeviceConnected(serial)

	if strings.Contains(*serial, ":") {
		fmt.Printf(ui.ColorGreen+"[✓] Device is ALREADY connected via Wi-Fi: %s\n"+ui.ColorReset, *serial)
		fmt.Println(ui.ColorCyan + "[i] You can already safely use XiOptimizer without any USB cables." + ui.ColorReset)
		return
	}

	fmt.Print(" • Detecting phone's local Wi-Fi IP address... ")
	ip := adb.GetDeviceWiFiIP(serial)
	if ip != "" {
		fmt.Printf(ui.ColorGreen+"[✓ Found: %s]\n"+ui.ColorReset, ip)
	} else {
		fmt.Println(ui.ColorYellow + "[– Could not auto-detect]" + ui.ColorReset)
		fmt.Print(ui.ColorBold + "   Please enter phone's Wi-Fi IP (from Wi-Fi Settings): " + ui.ColorReset)
		line := <-ui.InputChan
		ip = strings.TrimSpace(line)
		if ip == "" {
			fmt.Println(ui.ColorRed + "[!] Canceled. Wi-Fi IP address is required." + ui.ColorReset)
			return
		}
	}

	fmt.Print(" • Restarting adbd in TCP/IP mode on port 5555... ")
	_, err := adb.EnableTCPMode(5555, serial)
	if err != nil {
		fmt.Println(ui.ColorRed + "[✗ Failed]" + ui.ColorReset)
		fmt.Printf(ui.ColorRed+"    Error: %v\n"+ui.ColorReset, err)
		return
	}
	fmt.Println(ui.ColorGreen + "[✓ Ready]" + ui.ColorReset)

	time.Sleep(1 * time.Second)

	target := ip + ":5555"
	fmt.Printf(" • Connecting to %s... ", target)
	out, _ := adb.ConnectWiFi(target)
	if strings.Contains(out, "connected to") {
		fmt.Println(ui.ColorGreen + "[✓ Connected]" + ui.ColorReset)
		*serial = target
		adb.SaveWiFiHistory(target)
		fmt.Println(ui.ColorGreen + ui.ColorBold + "\n[🎉] WIRELESS ADB ACTIVATED!" + ui.ColorReset)
		fmt.Println(ui.ColorCyan + "[i] You can now safely UNPLUG the USB cable." + ui.ColorReset)
		fmt.Println(ui.ColorCyan + "[i] XiOptimizer will continue operating fully wirelessly over Wi-Fi." + ui.ColorReset)
	} else {
		fmt.Println(ui.ColorYellow + "[i Connection pending]" + ui.ColorReset)
		fmt.Printf("   Output: %s\n", out)
		fmt.Println(ui.ColorGray + "   Ensure phone and PC are on the same Wi-Fi router." + ui.ColorReset)
	}
}

// RunAllNetworkBatch configures AdGuard DNS and low-latency Wi-Fi
func RunAllNetworkBatch(serial *string) {
	RunSetAdGuardDNS(serial)
	fmt.Println()
	RunWiFiLowLatency(serial)
}
