package adb

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// WiFiHistoryEntry stores a saved Wi-Fi connection target
type WiFiHistoryEntry struct {
	Address   string `json:"address"`
	Timestamp int64  `json:"timestamp"`
}

// getHistoryFilePath resolves the absolute path to wifi_history.json
func getHistoryFilePath() string {
	execDir, err := os.Executable()
	if err == nil {
		baseDir := filepath.Dir(execDir)
		return filepath.Join(baseDir, "wifi_history.json")
	}
	return "wifi_history.json"
}

// LoadWiFiHistory loads saved Wi-Fi addresses, most recent first
func LoadWiFiHistory() []string {
	filePath := getHistoryFilePath()
	data, err := os.ReadFile(filePath)
	if err != nil {
		return []string{}
	}

	var entries []WiFiHistoryEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		// Fallback for simple string array format
		var simpleList []string
		if err2 := json.Unmarshal(data, &simpleList); err2 == nil {
			return filterUniqueAddresses(simpleList)
		}
		return []string{}
	}

	var result []string
	for _, e := range entries {
		addr := strings.TrimSpace(e.Address)
		if addr != "" {
			result = append(result, addr)
		}
	}
	return filterUniqueAddresses(result)
}

// SaveWiFiHistory persists a successfully connected Wi-Fi address to history
func SaveWiFiHistory(address string) {
	address = strings.TrimSpace(address)
	if address == "" {
		return
	}
	if !strings.Contains(address, ":") {
		address = address + ":5555"
	}

	existing := LoadWiFiHistory()
	var updated []WiFiHistoryEntry
	updated = append(updated, WiFiHistoryEntry{
		Address:   address,
		Timestamp: time.Now().Unix(),
	})

	for _, oldAddr := range existing {
		if !strings.EqualFold(oldAddr, address) {
			updated = append(updated, WiFiHistoryEntry{
				Address:   oldAddr,
				Timestamp: time.Now().Unix(),
			})
		}
		if len(updated) >= 5 {
			break
		}
	}

	data, err := json.MarshalIndent(updated, "", "  ")
	if err == nil {
		_ = os.WriteFile(getHistoryFilePath(), data, 0644)
	}
}

// ClearWiFiHistory wipes all remembered Wi-Fi target entries
func ClearWiFiHistory() {
	_ = os.Remove(getHistoryFilePath())
}

// GetLastWiFiAddress returns the most recently used Wi-Fi address or empty string
func GetLastWiFiAddress() string {
	hist := LoadWiFiHistory()
	if len(hist) > 0 {
		return hist[0]
	}
	return ""
}

// filterUniqueAddresses removes duplicate and invalid Wi-Fi strings
func filterUniqueAddresses(list []string) []string {
	seen := make(map[string]bool)
	var res []string
	for _, item := range list {
		item = strings.TrimSpace(item)
		if item != "" && !seen[item] {
			seen[item] = true
			res = append(res, item)
		}
	}
	return res
}
