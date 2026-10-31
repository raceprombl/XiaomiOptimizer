package packages

// PackageInfo holds metadata for a removable Android package
type PackageInfo struct {
	Name string
	Desc string
}

// SafeAdsPackages includes ad daemons, telemetry, and analytic trackers safe to remove across MIUI 11-14 & HyperOS
var SafeAdsPackages = []PackageInfo{
	{"com.miui.msa.global", "MIUI System Ads engine"},
	{"com.miui.analytics", "Usage analytics and telemetry"},
	{"com.miui.daemon", "Background logging daemon"},
	{"com.xiaomi.ugd", "User experience data collector"},
	{"com.miui.misightservice", "Mi Sight tracking daemon"},
	{"com.miui.bugreport", "Bug report uploader"},
	{"com.miui.hybrid", "Quick Apps ad framework (MIUI 11/12)"},
	{"com.miui.hybrid.accessory", "Quick Apps accessory services"},
	{"com.facebook.services", "Facebook analytics service"},
	{"com.facebook.system", "Facebook system installer"},
	{"com.facebook.appmanager", "Facebook app updater"},
	{"com.miui.android.fashiongallery", "Wallpaper carousel ads (MIUI 12.5+)"},
	{"com.mfashiongallery.emag", "Wallpaper carousel ads (MIUI 11/12)"},
	{"com.xiaomi.payment", "Mi Pay service"},
	{"com.miui.yellowpage", "Yellow Pages directory and ad services"},
	{"com.xiaomi.mipicks", "GetApps app store"},
	{"com.xiaomi.discover", "Xiaomi Discover news feed"},
	{"com.mi.global.shop", "Mi Store shopping app"},
	{"cn.wps.xiaomi.abroad.lite", "Mi Doc Viewer (WPS Office Lite)"},
	{"com.xiaomi.simactivate.service", "SIM activation tracking service"},
}

// BloatwarePackages includes removable vendor applications, promos, and media players
var BloatwarePackages = []PackageInfo{
	{"com.android.browser", "Mi Browser (MIUI 11/12)"},
	{"com.mi.globalbrowser", "Mi Browser (MIUI 12.5+)"},
	{"com.miui.videoplayer", "Mi Video player (with ads)"},
	{"com.miui.video", "Mi Video player legacy (MIUI 11)"},
	{"com.miui.player", "Mi Music player (with promos)"},
	{"com.miui.personalassistant", "App Vault legacy (MIUI 11/12)"},
	{"com.mi.globalminusscreen", "App Vault minus-one screen (MIUI 12+)"},
	{"com.miui.cleaner", "Built-in Cleaner utility"},
	{"com.mi.appfinder", "Mi App Finder"},
	{"video.like", "Pre-installed Like video app"},
	{"com.miui.virtualsim", "Mi Roaming / Virtual SIM"},
	{"com.xiaomi.midrop", "ShareMe / Mi Drop"},
}
