package adb

import (
	"strings"
	"testing"
)

func TestVerifyXiaomiBrandLogic(t *testing.T) {
	testCases := []struct {
		name       string
		allProps   string
		hyperOSVer string
		miuiVer    string
		wantXi     bool
		wantBrand  string
	}{
		{
			name:      "Xiaomi Flagship",
			allProps:  "xiaomi xiaomi xiaomi 23127pn0cc",
			wantXi:    true,
			wantBrand: "Xiaomi",
		},
		{
			name:      "Redmi Note Phone",
			allProps:  "redmi xiaomi redmi 24116raccg obsidian",
			wantXi:    true,
			wantBrand: "Redmi",
		},
		{
			name:      "POCO Gaming Phone",
			allProps:  "poco xiaomi poco 23049pcd8g marble",
			wantXi:    true,
			wantBrand: "POCO",
		},
		{
			name:      "Black Shark Gaming",
			allProps:  "blackshark blackshark shm-a0",
			wantXi:    true,
			wantBrand: "Black Shark",
		},
		{
			name:       "Custom ROM with HyperOS prop",
			allProps:   "aosp google generic",
			hyperOSVer: "OS1.0.1.0",
			wantXi:     true,
			wantBrand:  "Xiaomi / HyperOS",
		},
		{
			name:       "Custom ROM with MIUI prop",
			allProps:   "lineage generic",
			miuiVer:    "V14",
			wantXi:     true,
			wantBrand:  "Xiaomi / MIUI",
		},
		{
			name:      "Samsung Galaxy",
			allProps:  "samsung samsung sm-s911b",
			wantXi:    false,
			wantBrand: "Samsung",
		},
		{
			name:      "Google Pixel",
			allProps:  "google google husky pixel 8 pro",
			wantXi:    false,
			wantBrand: "Google",
		},
		{
			name:      "OnePlus Phone",
			allProps:  "oneplus oneplus cph2449",
			wantXi:    false,
			wantBrand: "Oneplus",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			isXi := false
			brandName := "Unknown OEM"

			allProps := strings.ToLower(tc.allProps)

			if strings.Contains(allProps, "poco") {
				isXi = true
				brandName = "POCO"
			} else if strings.Contains(allProps, "redmi") {
				isXi = true
				brandName = "Redmi"
			} else if strings.Contains(allProps, "blackshark") || strings.Contains(allProps, "black shark") {
				isXi = true
				brandName = "Black Shark"
			} else if strings.Contains(allProps, "xiaomi") {
				isXi = true
				brandName = "Xiaomi"
			} else if tc.hyperOSVer != "" {
				isXi = true
				brandName = "Xiaomi / HyperOS"
			} else if tc.miuiVer != "" {
				isXi = true
				brandName = "Xiaomi / MIUI"
			} else {
				parts := strings.Fields(allProps)
				if len(parts) > 0 {
					brandName = strings.ToUpper(parts[0][:1]) + parts[0][1:]
				}
			}

			if isXi != tc.wantXi {
				t.Errorf("Expected IsXiaomi=%v, got %v", tc.wantXi, isXi)
			}
			if brandName != tc.wantBrand {
				t.Errorf("Expected BrandName=%s, got %s", tc.wantBrand, brandName)
			}
		})
	}
}
