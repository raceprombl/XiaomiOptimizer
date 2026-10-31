package ui

import (
	"fmt"
	"strings"
)

// PrintBanner displays the ASCII banner strictly within 82 columns
func PrintBanner() {
	boxWidth := 80
	top := ColorOrange + "╔" + strings.Repeat("═", boxWidth) + "╗" + ColorReset
	bot := ColorOrange + "╚" + strings.Repeat("═", boxWidth) + "╝" + ColorReset
	empty := ColorOrange + "║" + strings.Repeat(" ", boxWidth) + "║" + ColorReset

	art := []string{
		" ██╗  ██╗██╗ ██████╗ ██████╗ ████████╗██╗███╗   ███╗██╗███████╗███████╗██████╗  ",
		" ╚██╗██╔╝██║██╔═══██╗██╔══██╗╚══██╔══╝██║████╗ ████║██║╚══███╔╝██╔════╝██╔══██╗ ",
		"  ╚███╔╝ ██║██║   ██║██████╔╝   ██║   ██║██╔████╔██║██║  ███╔╝ █████╗  ██████╔╝ ",
		"  ██╔██╗ ██║██║   ██║██╔═══╝    ██║   ██║██║╚██╔╝██║██║ ███╔╝  ██╔══╝  ██╔══██╗ ",
		" ██╔╝ ██╗██║╚██████╔╝██║        ██║   ██║██║ ╚═╝ ██║██║███████╗███████╗██║  ██║ ",
		" ╚═╝  ╚═╝╚═╝ ╚═════╝ ╚═╝        ╚═╝   ╚═╝╚═╝     ╚═╝╚═╝╚══════╝╚══════╝╚═╝  ╚═╝ ",
	}

	sub1 := Center("✦ HyperOS & MIUI 11 / 12 / 13 / 14 Optimizer ✦", boxWidth)
	sub2 := Center("Xiaomi • Redmi • POCO Debloater & Tweaker", boxWidth)

	fmt.Println(top)
	fmt.Println(empty)
	for _, a := range art {
		fmt.Println(ColorOrange + "║" + FitWidth(a, boxWidth) + ColorOrange + "║" + ColorReset)
	}
	fmt.Println(empty)
	fmt.Println(ColorOrange + "║" + ColorMiCyan + ColorBold + FitWidth(sub1, boxWidth) + ColorOrange + "║" + ColorReset)
	fmt.Println(ColorOrange + "║" + ColorGray + FitWidth(sub2, boxWidth) + ColorOrange + "║" + ColorReset)
	fmt.Println(bot)
}
