package ui

import (
	"fmt"
	"strings"

	"github.com/raceprombl/XiOptimizer/internal/adb"
)

// PrintDeviceCard displays device metadata in an exact 82-column card layout
func PrintDeviceCard(d adb.DeviceDetails) {
	top := ColorCyan + "┌" + strings.Repeat("─", 80) + "┐" + ColorReset
	header := ColorCyan + "│ " + ColorBold + FitWidth("CONNECTED DEVICE INFORMATION", 78) + ColorReset + ColorCyan + " │" + ColorReset
	mid := ColorCyan + "├" + strings.Repeat("─", 25) + "┬" + strings.Repeat("─", 54) + "┤" + ColorReset
	bot := ColorCyan + "└" + strings.Repeat("─", 25) + "┴" + strings.Repeat("─", 54) + "┘" + ColorReset

	items := []struct {
		label string
		val   string
		color string
	}{
		{"Model & Codename:", fmt.Sprintf("%s (%s)", d.Model, d.Codename), ColorBold},
		{"Serial Number (SN):", d.Serial, ColorReset},
		{"Android Version:", d.AndroidVer, ColorGreen},
		{"OS / ROM Version:", d.OSVersion, ColorOrange + ColorBold},
		{"Battery Status:", d.Battery, ColorYellow},
		{"ADB Connection:", "USB Debugging Authorized [OK]", ColorGreen},
	}

	fmt.Println(top)
	fmt.Println(header)
	fmt.Println(mid)
	for _, it := range items {
		l := FitWidth(it.label, 23)
		v := FitWidth(it.val, 52)
		fmt.Println(ColorCyan + "│ " + ColorReset + l + ColorCyan + " │ " + it.color + v + ColorReset + ColorCyan + " │" + ColorReset)
	}
	fmt.Println(bot)
}
