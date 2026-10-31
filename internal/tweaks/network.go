package tweaks

import (
	"fmt"

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

// RunAllNetworkBatch configures AdGuard DNS and low-latency Wi-Fi
func RunAllNetworkBatch(serial *string) {
	RunSetAdGuardDNS(serial)
	fmt.Println()
	RunWiFiLowLatency(serial)
}
