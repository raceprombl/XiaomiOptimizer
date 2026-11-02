package adb

import (
	"os"
	"testing"
)

func TestWiFiHistoryManagement(t *testing.T) {
	tempFile := "test_wifi_history.json"
	defer os.Remove(tempFile)

	list := filterUniqueAddresses([]string{
		"192.168.1.100:5555",
		"192.168.1.100:5555",
		"192.168.1.101:37841",
		"",
		"   ",
	})

	if len(list) != 2 {
		t.Fatalf("Expected 2 unique items, got %d", len(list))
	}
	if list[0] != "192.168.1.100:5555" || list[1] != "192.168.1.101:37841" {
		t.Errorf("Unexpected filter order: %v", list)
	}
}
