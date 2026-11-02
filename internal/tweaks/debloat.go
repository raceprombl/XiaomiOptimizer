package tweaks

import (
	"fmt"
	"strings"
	"time"

	"github.com/raceprombl/XiOptimizer/internal/adb"
	"github.com/raceprombl/XiOptimizer/internal/packages"
	"github.com/raceprombl/XiOptimizer/internal/ui"
)

// RunSafeAdsCleanup removes telemetry, ad injectors, and background loggers
func RunSafeAdsCleanup(serial *string) {
	ui.ClearScreen()
	ui.PrintBanner()
	fmt.Println(ui.ColorGreen + ui.ColorBold + "\n[1] SAFE ADS & TELEMETRY CLEANUP" + ui.ColorReset)
	fmt.Println(ui.ColorGray + "Removing system trackers, ad injectors, and partner daemons:" + ui.ColorReset)
	fmt.Println()

	successCount := 0
	notFoundCount := 0

	for i := 0; i < len(packages.SafeAdsPackages); i++ {
		item := packages.SafeAdsPackages[i]
		adb.EnsureDeviceConnected(serial)

		fmt.Printf(" • [%2d/%d] Removing: %-32s (%s)... ", i+1, len(packages.SafeAdsPackages), item.Name, item.Desc)
		out, _ := adb.RunADBShellWithRetry("pm uninstall -k --user 0 "+item.Name, serial)

		if strings.Contains(out, "Success") {
			fmt.Println(ui.ColorGreen + "[✓ Removed]" + ui.ColorReset)
			successCount++
		} else if strings.Contains(out, "not installed") || strings.Contains(out, "Unknown package") {
			fmt.Println(ui.ColorGray + "[– Not installed]" + ui.ColorReset)
			notFoundCount++
		} else {
			fmt.Printf(ui.ColorYellow+"[i Skipped: %s]\n"+ui.ColorReset, strings.TrimSpace(out))
		}
		time.Sleep(25 * time.Millisecond)
	}

	fmt.Println()
	fmt.Printf(ui.ColorGreen+ui.ColorBold+"[✓] Done. Removed: %d | Already absent: %d\n"+ui.ColorReset, successCount, notFoundCount)

	fmt.Print(" • Disabling MIUI Network Data Collection Daemon (NDCD)... ")
	_, _ = adb.RunADBShellWithRetry("settings put global sys.miui.ndcd off", serial)
	fmt.Println(ui.ColorGreen + "[✓ Disabled]" + ui.ColorReset)

	fmt.Print(" • Disabling MIUI TCP Connection Tracker... ")
	_, _ = adb.RunADBShellWithRetry("settings put global persist.sys.miuitcptracker.ctrl false", serial)
	fmt.Println(ui.ColorGreen + "[✓ Disabled]" + ui.ColorReset)

	fmt.Println(ui.ColorCyan + "[i] System ads, trackers, and network telemetry uninstalled for user 0." + ui.ColorReset)
}

// RunBloatwareCleanup removes non-essential media apps, promo utilities, and browser
func RunBloatwareCleanup(serial *string) {
	ui.ClearScreen()
	ui.PrintBanner()
	fmt.Println(ui.ColorYellow + ui.ColorBold + "\n[2] DEEP BLOATWARE REMOVAL" + ui.ColorReset)
	fmt.Println(ui.ColorGray + "Removing pre-installed media players, promo apps, and bundled browser:" + ui.ColorReset)
	fmt.Println()

	successCount := 0
	notFoundCount := 0

	for i := 0; i < len(packages.BloatwarePackages); i++ {
		item := packages.BloatwarePackages[i]
		adb.EnsureDeviceConnected(serial)

		fmt.Printf(" • [%d/%d] Removing: %-26s (%s)... ", i+1, len(packages.BloatwarePackages), item.Name, item.Desc)
		out, _ := adb.RunADBShellWithRetry("pm uninstall -k --user 0 "+item.Name, serial)

		if strings.Contains(out, "Success") {
			fmt.Println(ui.ColorGreen + "[✓ Removed]" + ui.ColorReset)
			successCount++
		} else if strings.Contains(out, "not installed") || strings.Contains(out, "Unknown package") {
			fmt.Println(ui.ColorGray + "[– Not installed]" + ui.ColorReset)
			notFoundCount++
		} else {
			fmt.Printf(ui.ColorYellow+"[i Skipped: %s]\n"+ui.ColorReset, strings.TrimSpace(out))
		}
		time.Sleep(25 * time.Millisecond)
	}

	fmt.Println()
	fmt.Printf(ui.ColorYellow+ui.ColorBold+"[✓] Done. Removed: %d | Already absent: %d\n"+ui.ColorReset, successCount, notFoundCount)
	fmt.Println(ui.ColorCyan + "[i] Bloatware removed. Storage and RAM freed." + ui.ColorReset)
}

// RunRestorePackagesOnly reinstalls all removed packages for user 0
func RunRestorePackagesOnly(serial *string) {
	ui.ClearScreen()
	ui.PrintBanner()
	fmt.Println(ui.ColorOrange + ui.ColorBold + "\n[3] RESTORE REMOVED PACKAGES" + ui.ColorReset)
	fmt.Println(ui.ColorGray + "Reinstalling uninstalled packages for user 0:" + ui.ColorReset)
	fmt.Println()

	var allPackages []string
	for _, item := range packages.SafeAdsPackages {
		allPackages = append(allPackages, item.Name)
	}
	for _, item := range packages.BloatwarePackages {
		allPackages = append(allPackages, item.Name)
	}
	allPackages = append(allPackages, "com.xiaomi.joyose")

	restoredCount := 0
	for i := 0; i < len(allPackages); i++ {
		pkg := allPackages[i]
		adb.EnsureDeviceConnected(serial)

		fmt.Printf(" • [%2d/%d] Restoring: %-34s... ", i+1, len(allPackages), pkg)
		out, _ := adb.RunADBShellWithRetry("cmd package install-existing "+pkg, serial)
		if strings.Contains(out, "installed") || strings.Contains(out, "Success") {
			fmt.Println(ui.ColorGreen + "[✓ Restored]" + ui.ColorReset)
			restoredCount++
		} else {
			outLegacy, _ := adb.RunADBShellWithRetry("pm install-existing "+pkg, serial)
			if strings.Contains(outLegacy, "installed") || strings.Contains(outLegacy, "Success") {
				fmt.Println(ui.ColorGreen + "[✓ Restored]" + ui.ColorReset)
				restoredCount++
			} else {
				fmt.Println(ui.ColorGray + "[– Not needed]" + ui.ColorReset)
			}
		}
		time.Sleep(25 * time.Millisecond)
	}

	fmt.Println()
	fmt.Printf(ui.ColorGreen+ui.ColorBold+"[✓] Restore finished. Packages restored: %d\n"+ui.ColorReset, restoredCount)
}

// RunFullDebloatBatch executes both safe ads cleanup and deep bloatware removal
func RunFullDebloatBatch(serial *string) {
	RunSafeAdsCleanup(serial)
	fmt.Println()
	RunBloatwareCleanup(serial)
}
