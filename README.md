# XiOptimizer

[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go)](https://golang.org)
[![Platform Support](https://img.shields.io/badge/Platforms-Windows%20%7C%20macOS%20%7C%20Linux-orange.svg)](#cross-platform-binaries)
[![GitHub Release](https://img.shields.io/github/v/release/raceprombl/XiOptimizer)](https://github.com/raceprombl/XiOptimizer/releases)

XiOptimizer is a lightweight, dependency-free command-line debloater and system performance tuner for Xiaomi, Redmi, and POCO smartphones running MIUI 11, 12, 12.5, 13, 14, and HyperOS 1.0, 2.0, 3.0.

No root access is required. All operations communicate directly via Android Debug Bridge (ADB) using non-destructive user partition overrides (`pm uninstall -k --user 0` and Android system settings).

```text
╔════════════════════════════════════════════════════════════════════════════════╗
║                                                                                ║
║  ██╗  ██╗██╗ ██████╗ ██████╗ ████████╗██╗███╗   ███╗██╗███████╗███████╗██████╗  ║
║  ╚██╗██╔╝██║██╔═══██╗██╔══██╗╚══██╔══╝██║████╗ ████║██║╚══███╔╝██╔════╝██╔══██╗ ║
║   ╚███╔╝ ██║██║   ██║██████╔╝   ██║   ██║██╔████╔██║██║  ███╔╝ █████╗  ██████╔╝ ║
║   ██╔██╗ ██║██║   ██║██╔═══╝    ██║   ██║██║╚██╔╝██║██║ ███╔╝  ██╔══╝  ██╔══██╗ ║
║  ██╔╝ ██╗██║╚██████╔╝██║        ██║   ██║██║ ╚═╝ ██║██║███████╗███████╗██║  ██║ ║
║  ╚═╝  ╚═╝╚═╝ ╚═════╝ ╚═╝        ╚═╝   ╚═╝╚═╝     ╚═╝╚═╝╚══════╝╚══════╝╚═╝  ╚═╝ ║
║                                                                                ║
║                  ✦ HyperOS & MIUI 11 / 12 / 13 / 14 Optimizer ✦                 ║
║                    Xiaomi • Redmi • POCO Debloater & Tweaker                   ║
╚════════════════════════════════════════════════════════════════════════════════╝

┌────────────────────────────────────────────────────────────────────────────────┐
│ CONNECTED DEVICE INFORMATION                                                   │
├─────────────────────────┬──────────────────────────────────────────────────────┤
│ Model & Codename:       │ 24117RN76G (tanzanite_global)                        │
│ Serial Number (SN):     │ nfirtoemmv7lfaca                                     │
│ Android Version:        │ Android 16                                           │
│ OS / ROM Version:       │ Xiaomi HyperOS 3.0 [OS3.0]                           │
│ Battery Status:         │ 82% [Charging (+), 33.5°C]                           │
│ ADB Connection:         │ USB Debugging Authorized [OK]                        │
└─────────────────────────┴──────────────────────────────────────────────────────┘

╔════════════════════════════════════════════════════════════════════════════════╗
║  XIOPTIMIZER — HYPEROS & MIUI MASTER CONTROL                                   ║
║  Select a category below to optimize your Xiaomi device                        ║
╠════════════════════════════════════════════════════════════════════════════════╣
║  [1] 🛡 Debloat & Telemetry              Safe Ads, Bloatware, Restore apps     ║
║  [2] ⚡ Performance & Gaming             Joyose, 120Hz, Virtual RAM off, ART   ║
║  [3] 🔋 Battery & Deep Sleep             Aggressive Doze, Wi-Fi/BT scan off    ║
║  [4] 🌐 Network & Private DNS            AdGuard adblock, Cloudflare, Low-Ping ║
║  [5] 🎨 UI & Fluidity                    0.5x animations, apps freezer, blur   ║
║  [6] 🔮 Special Flags & Audio            Safe volume bypass, Camera lab, APK   ║
║  [7] 🛠 Maintenance & Backups            TRIM storage, JSON state snapshots    ║
║  [0] 🚪 Exit                             Safely quit XiOptimizer               ║
╚════════════════════════════════════════════════════════════════════════════════╝
```

---

## Architecture & Codebase Layout

The project follows a modular Go package structure designed for maintainability and extensibility:

```text
XiOptimizer/
├── internal/
│   ├── adb/          # ADB client wrappers, device state polling, and HTTPS Google installer
│   ├── backup/       # Device state snapshot manager (JSON capture and exact state restore)
│   ├── packages/     # Curated lists of safe ad trackers and removable bloatware
│   ├── tweaks/       # Modular handlers for Debloat, Performance, Battery, UI, and Network
│   └── ui/           # Strict 82-column TUI layout engine, ASCII banner, and menus
├── backups/          # Storage directory for captured device JSON state snapshots
├── .github/          # GitHub Actions cross-compilation release workflow
├── versioninfo.json  # Windows PE metadata definition
├── main.go           # Application entry point and menu router
└── go.mod            # Pure Go standard library dependencies
```

---

## Key Capabilities

### 1. Debloat & Privacy Hub
* **Safe Ads & Telemetry Cleanup**: Strips out background analytics daemons and ad injection services (`com.miui.msa.global`, `com.miui.analytics`, `com.miui.daemon`, `com.xiaomi.ugd`, `com.miui.misightservice`, `cn.wps.xiaomi.abroad.lite`, Facebook system background services, and wallpaper carousel ads).
* **Deep Bloatware Removal**: Removes bundled vendor applications (`com.android.browser`, `com.mi.globalbrowser`, `com.miui.videoplayer`, `com.miui.player`, `com.mi.globalminusscreen` App Vault, `com.miui.cleaner`, and `video.like`).
* **Instant Reinstallation / Restore**: Restore any or all uninstalled packages back to stock anytime via `cmd package install-existing`.

### 2. Performance & Gaming Hub
* **Joyose Gaming Throttle Removal**: Disables `com.xiaomi.joyose`, bypassing artificial FPS caps and screen dimming in games like Genshin Impact, PUBG Mobile, and Call of Duty.
* **120Hz Display Lock**: Clears PowerKeeper throttling profiles and forces refresh rates to stay locked at 120Hz across apps that usually drop to 60Hz (Chrome, YouTube, Telegram, Maps).
* **Disable Virtual RAM (ExtM / UFS Swap)**: Disables MIUI/HyperOS Memory Extension. Paging RAM into slow flash storage causes severe 0.1% low frame drops and micro-stutters. Turning it off restores smooth frame pacing and preserves flash NAND endurance.
* **Disable Phantom Process Killer**: Increases Android 12+ child process limit to 2147483647, preventing the OS from killing background Termux sessions, compilers, and bots.
* **ART Speed-Profile Ahead-Of-Time Bytecode Compilation**: Compiles hot bytecode for system UI, launcher, and user apps for faster cold-starts and higher UI framerates.

### 3. Battery & Deep Sleep Hub
* **Aggressive Doze & Sleep Timeout**: Shortens idle timeout before device enters deep sleep from 30 minutes to 30 seconds while whitelisting priority messaging apps (Telegram, WhatsApp).
* **Disable Wi-Fi & Bluetooth Background Scanning**: Turns off always-on background beacon discovery and location radio polling.
* **Standby Modem & Telemetry Tuning**: Disables cellular radio active-hold when connected to Wi-Fi (`mobile_data_always_on 0`) and stops background crash uploaders.
* **Minimize System Logd Buffer**: Drops logd ring-buffer size to 64K, cutting background CPU wakeups and diagnostic write wear.

### 4. Network & Private DNS Hub
* **AdGuard Private DNS (DoT)**: System-wide adblocking through encrypted DNS-over-TLS (`dns.adguard-dns.com`). Blocks banner ads and popups across apps and browsers without running battery-heavy VPN apps.
* **Cloudflare 1.1.1.1 Private DNS**: Ultra-low latency DNS resolution for gaming connections.
* **Wi-Fi Low-Latency Mode & Roaming Fix**: Forces Wi-Fi HAL low-latency mode and disables Xiaomi cloud network switching to prevent ping spikes during online gaming.

### 5. UI, Fluidity & Visual Hub
* **0.5x Fast System Animations**: Cuts window, transition, and animator durations in half for instant UI feedback.
* **Native Cached Apps Freezer**: Enables Linux kernel cgroup process freezing (`cached_apps_freezer`), preventing idle background apps from burning CPU cycles.
* **Touch Latency & Response Tuning**: Optimizes tap thresholds, long-press latency, and touch sampling parameters.
* **Unlock Control Center Gaussian Blur**: Overrides GPU `deviceLevelList` to `v:1,c:3,g:3` and enables `window_blurs_enabled`, replacing the ugly grey background on budget Redmi and POCO phones with real-time blur shaders.
* **Hide Gesture Bar Indicator Line**: Reclaims display real estate by hiding the navigation line at the bottom of the screen.

### 6. Special System Flags & Audio Hub
* **Bypass Safe Headphone Volume Warning**: Disables European and regional 60% headphone volume reset popups.
* **Unlock Hidden Camera Lab Features**: Creates the internal trigger directory to enable experimental camera toggles in the stock camera application.
* **Bypass APK Verification Lag**: Skips verification countdowns and security checks on ADB sideloads.

### 7. Storage, Maintenance & State Snapshots
* **Device State Snapshots (JSON)**: Captures a complete snapshot of all Android settings namespaces (`global`, `secure`, `system`) and all installed packages before any changes are made.
* **Exact State Restoration**: Rather than generic defaults, restores your device back to the exact snapshot taken before optimization.
* **Flash Storage TRIM**: Executes `sm fstrim` across `/data` and `/system` to restore 4K random write speeds.
* **Package Cache Pruning**: Frees gigabytes of stale app caches via `pm trim-caches`.

---

## Integrity & Reconnection Safety

* **Strict HTTPS & Executable Verification**: ADB binaries are downloaded exclusively from official Google servers (`https://dl.google.com/android/repository/...`), verified for archive integrity, and validated by testing `adb version` output before execution.
* **Non-Blocking Menu Loop**: Unplugging the USB cable while browsing menus immediately drops back to the USB connection discovery screen without freezing.
* **Interrupted Batch Protection**: If the phone is disconnected mid-batch during uninstalls or parameter updates, XiOptimizer pauses execution, displays a reconnect banner, and resumes from the exact package once the cable is re-inserted.

---

## Cross-Platform Binaries

Standalone binaries compiled for 64-bit systems without external runtime dependencies:

| Operating System | Architecture | Binary File |
| :--- | :--- | :--- |
| **Windows** | x86_64 / amd64 | `XiOptimizer-windows-amd64.exe` (or `XiOptimizer-windows-amd64.zip`) |
| **macOS** (Apple Silicon) | arm64 (M1/M2/M3/M4) | `XiOptimizer-darwin-arm64` |
| **macOS** (Intel) | x86_64 / amd64 | `XiOptimizer-darwin-amd64` |
| **Linux** | x86_64 / amd64 | `XiOptimizer-linux-amd64` |
| **Linux** | arm64 (Raspberry Pi/SBC) | `XiOptimizer-linux-arm64` |

Download pre-built releases directly from the [Releases Page](https://github.com/raceprombl/XiOptimizer/releases).

---

## Quick Start

### 1. Enable USB Debugging on Your Device
1. Open **Settings** -> **About Phone**.
2. Tap **OS Version** or **MIUI Version** 7 times until developer options are unlocked.
3. Go to **Additional Settings** -> **Developer Options**.
4. Enable **USB Debugging**.
5. Enable **Install via USB** and **USB Debugging (Security Settings)** (required on some MIUI builds for package operations).

### 2. Run XiOptimizer

#### On Windows:
Double-click `XiOptimizer.exe` or execute from PowerShell / Command Prompt:
```powershell
.\XiOptimizer.exe
```

#### On macOS:
Open Terminal:
```bash
chmod +x XiOptimizer-darwin-arm64
./XiOptimizer-darwin-arm64
```

#### On Linux:
Open Terminal:
```bash
chmod +x XiOptimizer-linux-amd64
./XiOptimizer-linux-amd64
```

---

## Building from Source

XiOptimizer is written strictly using the Go standard library. Building requires Go 1.22+:

```bash
# Clone the repository
git clone https://github.com/raceprombl/XiOptimizer.git
cd XiOptimizer

# Build for current host operating system
go build -trimpath -o XiOptimizer .

# Cross-compile for Windows
GOOS=windows GOARCH=amd64 go build -trimpath -o XiOptimizer-windows-amd64.exe .

# Cross-compile for macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -trimpath -o XiOptimizer-darwin-arm64 .

# Cross-compile for Linux
GOOS=linux GOARCH=amd64 go build -trimpath -o XiOptimizer-linux-amd64 .
```

---

## License

This project is licensed under the Apache 2.0 License. See the [LICENSE](LICENSE) file for details.
