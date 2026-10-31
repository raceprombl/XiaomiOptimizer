package ui

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// ClearScreen clears the terminal screen across Windows, macOS, and Linux
func ClearScreen() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		_ = cmd.Run()
	} else {
		fmt.Print("\033[H\033[2J")
	}
}

// RuneWidth returns the true visual display width of a rune in a monospace terminal
func RuneWidth(r rune) int {
	if r >= 32 && r < 127 {
		return 1
	}
	if r < 32 || (r >= 127 && r < 160) {
		return 0
	}
	// Box drawing, blocks, geometric shapes (1 column)
	if (r >= 0x2500 && r <= 0x257F) || (r >= 0x2580 && r <= 0x259F) || (r >= 0x25A0 && r <= 0x25FF) {
		return 1
	}
	// Emojis and miscellaneous symbols (2 columns)
	if (r >= 0x2600 && r <= 0x27BF) ||
		(r >= 0x2B00 && r <= 0x2BFF) ||
		(r >= 0x1F000 && r <= 0x1FAFF) {
		return 2
	}
	return 1
}

// VisualWidth calculates the visual display width of a string, ignoring ANSI escape sequences
func VisualWidth(s string) int {
	w := 0
	inAnsi := false
	for _, r := range s {
		if r == '\033' {
			inAnsi = true
			continue
		}
		if inAnsi {
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
				inAnsi = false
			}
			continue
		}
		w += RuneWidth(r)
	}
	return w
}

// PadRight pads string with spaces on the right to match the exact target display width
func PadRight(s string, target int) string {
	w := VisualWidth(s)
	if w >= target {
		return s
	}
	return s + strings.Repeat(" ", target-w)
}
