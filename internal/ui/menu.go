package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/raceprombl/XiOptimizer/internal/adb"
)

// Non-blocking keyboard input channel
var InputChan = make(chan string, 32)

// StartInputReader starts asynchronous background reader for keyboard input
func StartInputReader() {
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				return
			}
			InputChan <- strings.TrimSpace(line)
		}
	}()
}

// DrainInput discards any buffered or stale inputs
func DrainInput() {
	for {
		select {
		case <-InputChan:
		default:
			return
		}
	}
}

// WaitEnter waits for user to press Enter
func WaitEnter() {
	DrainInput()
	fmt.Print(ColorGray + "\nPress [Enter] to return..." + ColorReset)
	<-InputChan
}

// ActionFunc executes a specific tweak routine
type ActionFunc func(serial *string)

// MenuItem represents an interactive option in a menu
type MenuItem struct {
	Key         string
	Icon        string
	Title       string
	Description string
	Handler     ActionFunc
	SubMenu     *Menu
}

// Menu represents a navigable screen of options
type Menu struct {
	ID       string
	Title    string
	Subtitle string
	Parent   *Menu
	Items    []MenuItem
}

// RenderBoxedMenu renders any menu screen strictly within 82 columns
func RenderBoxedMenu(m *Menu) {
	const innerWidth = 78
	top := ColorMiCyan + "╔" + strings.Repeat("═", 80) + "╗" + ColorReset
	mid := ColorMiCyan + "╠" + strings.Repeat("═", 80) + "╣" + ColorReset
	bot := ColorMiCyan + "╚" + strings.Repeat("═", 80) + "╝" + ColorReset

	fmt.Println(top)
	titleLine := fmt.Sprintf(" %s%s%s", ColorBold, m.Title, ColorReset)
	fmt.Println(ColorMiCyan + "║ " + PadRight(titleLine, innerWidth) + ColorMiCyan + " ║" + ColorReset)

	if m.Subtitle != "" {
		subLine := fmt.Sprintf(" %s%s%s", ColorGray, m.Subtitle, ColorReset)
		fmt.Println(ColorMiCyan + "║ " + PadRight(subLine, innerWidth) + ColorMiCyan + " ║" + ColorReset)
	}
	fmt.Println(mid)

	for _, it := range m.Items {
		var itemStr string
		if it.Key == "0" {
			itemStr = fmt.Sprintf(" %s[%s]%s %s %-32s %s%s%s",
				ColorRed, it.Key, ColorReset, it.Icon, it.Title, ColorGray, it.Description, ColorReset)
		} else if it.Key == "A" {
			itemStr = fmt.Sprintf(" %s[%s]%s %s %-32s %s%s%s",
				ColorYellow+ColorBold, it.Key, ColorReset, it.Icon, it.Title, ColorGray, it.Description, ColorReset)
		} else {
			itemStr = fmt.Sprintf(" %s[%s]%s %s %-32s %s%s%s",
				ColorGreen, it.Key, ColorReset, it.Icon, it.Title, ColorGray, it.Description, ColorReset)
		}
		fmt.Println(ColorMiCyan + "║ " + PadRight(itemStr, innerWidth) + ColorMiCyan + " ║" + ColorReset)
	}
	fmt.Println(bot)
	fmt.Println()
	fmt.Print(ColorBold + "Select an option and press Enter: " + ColorReset)
}

// PromptSelectionWithCableWatch prompts for menu input while actively monitoring cable connection
func PromptSelectionWithCableWatch(serial *string) (string, bool) {
	DrainInput()
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case choice := <-InputChan:
			return strings.ToUpper(strings.TrimSpace(choice)), true
		case <-ticker.C:
			_, state := adb.CheckDeviceState()
			if state != "device" {
				return "", false
			}
		}
	}
}
