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

	sub1 := strings.Repeat(" ", 17) + "✦ HyperOS & MIUI 11 / 12 / 13 / 14 Optimizer ✦" + strings.Repeat(" ", 17)
	sub2 := strings.Repeat(" ", 19) + "Xiaomi • Redmi • POCO Debloater & Tweaker" + strings.Repeat(" ", 20)

	fmt.Println(top)
	fmt.Println(empty)
	for _, a := range art {
		fmt.Println(ColorOrange + "║" + a + "║" + ColorReset)
	}
	fmt.Println(empty)
	fmt.Println(ColorOrange + "║" + ColorMiCyan + ColorBold + sub1 + ColorOrange + "║" + ColorReset)
	fmt.Println(ColorOrange + "║" + ColorGray + sub2 + ColorOrange + "║" + ColorReset)
	fmt.Println(bot)
}
