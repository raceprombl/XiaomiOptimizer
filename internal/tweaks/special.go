package tweaks

import (
	"fmt"

	"github.com/raceprombl/XiOptimizer/internal/adb"
	"github.com/raceprombl/XiOptimizer/internal/ui"
)

// RunBypassSafeVolume removes the regional headphone volume limiter
func RunBypassSafeVolume(serial *string) {
	ui.ClearScreen()
	ui.PrintBanner()
	fmt.Println(ui.ColorYellow + ui.ColorBold + "\n[1] BYPASS SAFE HEADPHONE VOLUME WARNING" + ui.ColorReset)
	fmt.Println(ui.ColorGray + "Disabling EU/regional 60% headphone volume reset popups:" + ui.ColorReset)
	fmt.Println()

	adb.EnsureDeviceConnected(serial)
	fmt.Print(" • Overriding safe audio volume state... ")
	_, _ = adb.RunADBShellWithRetry("settings put global audio_safe_volume_state 2", serial)
	_, _ = adb.RunADBShellWithRetry("settings put global audio.safemedia.bypass true", serial)
	_, _ = adb.RunADBShellWithRetry("settings put system safe_volume 0", serial)
	fmt.Println(ui.ColorGreen + "[✓ Applied]" + ui.ColorReset)

	fmt.Println(ui.ColorCyan + "\n[i] Safe volume dialog disabled. Full headphone output unlocked." + ui.ColorReset)
}

// RunUnlockCameraLab creates trigger flag to enable experimental camera settings
func RunUnlockCameraLab(serial *string) {
	ui.ClearScreen()
	ui.PrintBanner()
	fmt.Println(ui.ColorYellow + ui.ColorBold + "\n[2] UNLOCK HIDDEN CAMERA LAB FEATURES" + ui.ColorReset)
	fmt.Println(ui.ColorGray + "Enabling internal experimental camera settings and post-processing flags:" + ui.ColorReset)
	fmt.Println()

	adb.EnsureDeviceConnected(serial)
	fmt.Print(" • Creating camera lab options trigger flag... ")
	_, _ = adb.RunADBShellWithRetry("mkdir -p /sdcard/DCIM/Camera/lab_options_visible", serial)
	fmt.Println(ui.ColorGreen + "[✓ Created]" + ui.ColorReset)

	fmt.Println(ui.ColorCyan + "\n[i] Camera Lab unlocked. Open Stock Camera -> Settings to view Experimental options." + ui.ColorReset)
}

// RunBypassInstallerLag skips package verifier countdown on sideloads
func RunBypassInstallerLag(serial *string) {
	ui.ClearScreen()
	ui.PrintBanner()
	fmt.Println(ui.ColorYellow + ui.ColorBold + "\n[3] DISABLE APK INSTALL VERIFICATION DELAYS" + ui.ColorReset)
	fmt.Println(ui.ColorGray + "Eliminating package verifier countdown timers and scan pauses for ADB:" + ui.ColorReset)
	fmt.Println()

	adb.EnsureDeviceConnected(serial)
	fmt.Print(" • Disabling package verification checks... ")
	_, _ = adb.RunADBShellWithRetry("settings put global verifier_verify_adb_installs 0", serial)
	_, _ = adb.RunADBShellWithRetry("settings put global package_verifier_enable 0", serial)
	_, _ = adb.RunADBShellWithRetry("settings put global package_verifier_user_consent -1", serial)
	fmt.Println(ui.ColorGreen + "[✓ Disabled]" + ui.ColorReset)

	fmt.Println(ui.ColorCyan + "\n[i] APK install verification skipped. App installs execute instantly." + ui.ColorReset)
}

// RunBurnInProtection enables AMOLED pixel-shifting protection for status bar icons
func RunBurnInProtection(serial *string) {
	ui.ClearScreen()
	ui.PrintBanner()
	fmt.Println(ui.ColorYellow + ui.ColorBold + "\n[4] AMOLED DISPLAY BURN-IN PROTECTION" + ui.ColorReset)
	fmt.Println(ui.ColorGray + "Enabling periodic micro-pixel shifting to prevent status bar burn-in:" + ui.ColorReset)
	fmt.Println()

	adb.EnsureDeviceConnected(serial)
	fmt.Print(" • Enabling AMOLED burn-in protection in system framework... ")
	_, _ = adb.RunADBShellWithRetry("settings put global burn_in_protection 1", serial)
	fmt.Println(ui.ColorGreen + "[✓ Enabled]" + ui.ColorReset)

	fmt.Println(ui.ColorCyan + "\n[i] Burn-in protection active. Status bar icons will periodically micro-shift")
	fmt.Println("    to preserve AMOLED panel longevity on Xiaomi/Redmi/POCO devices." + ui.ColorReset)
}

// RunAllSpecialBatch applies all special system flags
func RunAllSpecialBatch(serial *string) {
	RunBypassSafeVolume(serial)
	fmt.Println()
	RunUnlockCameraLab(serial)
	fmt.Println()
	RunBypassInstallerLag(serial)
	fmt.Println()
	RunBurnInProtection(serial)
}
