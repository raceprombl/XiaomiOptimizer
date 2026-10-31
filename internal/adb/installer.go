package adb

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	colorOrange = "\033[38;2;255;103;0m"
	colorMiCyan = "\033[38;2;0;210;255m"
	colorGray   = "\033[90m"
)

func clearScreen() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		_ = cmd.Run()
	} else {
		fmt.Print("\033[H\033[2J")
	}
}

// EnsureADB locates existing adb or downloads official Google platform-tools with integrity check
func EnsureADB() error {
	execDir, _ := os.Executable()
	baseDir := filepath.Dir(execDir)

	adbName := "adb"
	if runtime.GOOS == "windows" {
		adbName = "adb.exe"
	}

	// 1. Check local directory
	localAdb := filepath.Join(baseDir, adbName)
	if _, err := os.Stat(localAdb); err == nil {
		if verifyADBExecutable(localAdb) {
			SetAdbPath(localAdb)
			return nil
		}
	}

	// 2. Check current working directory
	if _, err := os.Stat(adbName); err == nil {
		if verifyADBExecutable(adbName) {
			SetAdbPath(adbName)
			return nil
		}
	}

	// 3. Check system PATH
	if p, err := exec.LookPath(adbName); err == nil {
		if verifyADBExecutable(p) {
			SetAdbPath(p)
			return nil
		}
	}

	// 4. Download from official Google repository via HTTPS
	clearScreen()
	fmt.Println(colorOrange + "╔════════════════════════════════════════════════════════════════════════════════╗" + colorReset)
	fmt.Println(colorOrange + "║" + colorYellow + colorBold + "  [FIRST RUN] Downloading Android ADB platform-tools...                         " + colorOrange + "║" + colorReset)
	fmt.Println(colorOrange + "║" + colorGray + "  Source: https://dl.google.com/android/repository/platform-tools...             " + colorOrange + "║" + colorReset)
	fmt.Println(colorOrange + "║" + colorMiCyan + "  Integrity: Official Google HTTPS CDN + Executable Verification                 " + colorOrange + "║" + colorReset)
	fmt.Println(colorOrange + "╚════════════════════════════════════════════════════════════════════════════════╝" + colorReset)
	fmt.Println()

	if runtime.GOOS == "windows" {
		psScript := fmt.Sprintf(`
			$ProgressPreference = 'Continue'
			$zip = Join-Path '%s' 'platform-tools.zip'
			$targetDir = '%s'
			$url = 'https://dl.google.com/android/repository/platform-tools-latest-windows.zip'
			Write-Host 'Downloading Google Platform Tools via HTTPS...' -ForegroundColor Cyan
			Invoke-WebRequest -Uri $url -OutFile $zip
			if ((Get-Item $zip).Length -lt 5000000) {
				Write-Error 'Downloaded archive is corrupted or incomplete.'
				exit 1
			}
			Write-Host 'Extracting ADB components...' -ForegroundColor Cyan
			Expand-Archive -Path $zip -DestinationPath $targetDir -Force
			Get-ChildItem -Path (Join-Path $targetDir 'platform-tools') | Move-Item -Destination $targetDir -Force
			Remove-Item -Path (Join-Path $targetDir 'platform-tools') -Recurse -Force -ErrorAction SilentlyContinue
			Remove-Item -Path $zip -Force -ErrorAction SilentlyContinue
		`, baseDir, baseDir)
		cmd := exec.Command("powershell", "-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", psScript)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("powershell adb download failed: %w", err)
		}
	} else if runtime.GOOS == "darwin" {
		shScript := fmt.Sprintf(`curl -L -f -o "%s/platform-tools.zip" https://dl.google.com/android/repository/platform-tools-latest-darwin.zip && unzip -q -o "%s/platform-tools.zip" -d "%s" && mv "%s/platform-tools/adb" "%s/adb" && chmod +x "%s/adb" && rm -rf "%s/platform-tools" "%s/platform-tools.zip"`, baseDir, baseDir, baseDir, baseDir, baseDir, baseDir, baseDir, baseDir)
		cmd := exec.Command("sh", "-c", shScript)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("curl adb download failed: %w", err)
		}
	} else {
		shScript := fmt.Sprintf(`curl -L -f -o "%s/platform-tools.zip" https://dl.google.com/android/repository/platform-tools-latest-linux.zip && unzip -q -o "%s/platform-tools.zip" -d "%s" && mv "%s/platform-tools/adb" "%s/adb" && chmod +x "%s/adb" && rm -rf "%s/platform-tools" "%s/platform-tools.zip"`, baseDir, baseDir, baseDir, baseDir, baseDir, baseDir, baseDir, baseDir)
		cmd := exec.Command("sh", "-c", shScript)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("curl adb download failed: %w", err)
		}
	}

	targetAdb := filepath.Join(baseDir, adbName)
	if !verifyADBExecutable(targetAdb) {
		return fmt.Errorf("adb executable verification failed after download")
	}

	SetAdbPath(targetAdb)
	fmt.Println(colorGreen + colorBold + "[✓] ADB verified and ready." + colorReset)
	time.Sleep(1200 * time.Millisecond)
	return nil
}

// verifyADBExecutable runs `adb version` to confirm binary execution integrity
func verifyADBExecutable(path string) bool {
	cmd := exec.Command(path, "version")
	out, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(out), "Android Debug Bridge")
}
