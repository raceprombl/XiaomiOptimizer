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

// RuneWidth returns the visual display width of a rune in a monospace terminal
func RuneWidth(r rune) int {
	// Control and invisible characters
	if r < 32 || (r >= 127 && r < 160) {
		return 0
	}
	// Variation selectors and zero-width joiners
	if (r >= 0xFE00 && r <= 0xFE0F) || (r >= 0xE0100 && r <= 0xE01EF) || r == 0x200D || r == 0x200B || r == 0xFEFF {
		return 0
	}
	// Combining diacritical marks
	if (r >= 0x0300 && r <= 0x036F) || (r >= 0x1AB0 && r <= 0x1AFF) || (r >= 0x1DC0 && r <= 0x1DFF) || (r >= 0x20D0 && r <= 0x20FF) || (r >= 0xFE20 && r <= 0xFE2F) {
		return 0
	}
	// Standard ASCII printable
	if r >= 32 && r < 127 {
		return 1
	}
	// Box drawing, block elements, geometric shapes (1 column)
	if (r >= 0x2500 && r <= 0x257F) || (r >= 0x2580 && r <= 0x259F) || (r >= 0x25A0 && r <= 0x25FF) {
		return 1
	}
	// Narrow symbols and dingbats (1 column in monospace fonts)
	if r == 0x2726 || r == 0x2727 || r == 0x2605 || r == 0x2606 || (r >= 0x2713 && r <= 0x2718) || r == 0x2022 || r == 0x00B7 {
		return 1
	}
	// Arrows (1 column)
	if r >= 0x2190 && r <= 0x21FF {
		return 1
	}
	// Emojis and miscellaneous graphic symbols (2 columns)
	if (r >= 0x2600 && r <= 0x27BF) ||
		(r >= 0x2B00 && r <= 0x2BFF) ||
		(r >= 0x1F000 && r <= 0x1FAFF) {
		return 2
	}
	// CJK ideographs (2 columns)
	if (r >= 0x2E80 && r <= 0x9FFF) || (r >= 0xAC00 && r <= 0xD7AF) || (r >= 0xF900 && r <= 0xFAFF) {
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

// FitWidth fits a string into an exact visual width target by padding or safely clipping runes
func FitWidth(s string, target int) string {
	curW := 0
	var sb strings.Builder
	inAnsi := false
	var ansiBuf strings.Builder
	hasColor := false

	runes := []rune(s)
	for i := 0; i < len(runes); i++ {
		r := runes[i]
		if r == '\033' {
			inAnsi = true
			ansiBuf.Reset()
			ansiBuf.WriteRune(r)
			continue
		}
		if inAnsi {
			ansiBuf.WriteRune(r)
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
				inAnsi = false
				sb.WriteString(ansiBuf.String())
				hasColor = true
			}
			continue
		}

		rw := RuneWidth(r)
		if curW+rw > target {
			break
		}
		sb.WriteRune(r)
		curW += rw
	}

	if hasColor {
		sb.WriteString(ColorReset)
	}

	if curW < target {
		sb.WriteString(strings.Repeat(" ", target-curW))
	}
	return sb.String()
}

// PadRight guarantees the string matches the exact target display width
func PadRight(s string, target int) string {
	return FitWidth(s, target)
}

// Center centers a string within the target width
func Center(s string, target int) string {
	w := VisualWidth(s)
	if w >= target {
		return FitWidth(s, target)
	}
	left := (target - w) / 2
	right := target - w - left
	return strings.Repeat(" ", left) + s + strings.Repeat(" ", right)
}
