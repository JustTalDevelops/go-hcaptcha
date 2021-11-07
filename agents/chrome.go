package agents

import (
	"github.com/justtaldevelops/hcaptcha-solver-go/utils"
	"time"
)

// Chrome is the agent for Google Chrome.
type Chrome struct {
	// screenSize is the screen size of the Chrome browser.
	screenSize screenSize
	// availableScreenSize is the available screen size of the Chrome browser.
	availableScreenSize screenSize
	// cpuCount is the number of CPUs allowed to the Chrome browser.
	cpuCount int
	// memorySize is the memory size in gigabytes allowed to the Chrome browser.
	memorySize int
	// unixOffset is the offset of the unix timestamps.
	unixOffset int64
}

// NewChrome creates a new Chrome agent.
func NewChrome() *Chrome {
	possibleScreenSizes := [][]screenSize{
		{{1920, 1080}, {1920, 1040}},
		{{2560, 1440}, {2560, 1400}},
	}
	possibleCpuCounts := []int{2, 4, 8, 16}
	possibleMemorySizes := []int{2, 4, 8, 16}

	screenSizes := possibleScreenSizes[utils.Between(0, len(possibleScreenSizes)-1)]

	return &Chrome{
		screenSize:          screenSizes[0],
		availableScreenSize: screenSizes[1],
		cpuCount:            possibleCpuCounts[utils.Between(0, len(possibleCpuCounts)-1)],
		memorySize:          possibleMemorySizes[utils.Between(0, len(possibleMemorySizes)-1)],
	}
}

// ScreenProperties ...
func (c *Chrome) ScreenProperties() map[string]interface{} {
	return map[string]interface{}{
		"availWidth":  c.availableScreenSize[0],
		"availHeight": c.availableScreenSize[1],
		"width":       c.screenSize[0],
		"height":      c.screenSize[1],
		"colorDepth":  24,
		"pixelDepth":  24,
		"availLeft":   0,
		"availTop":    0,
	}
}

// NavigatorProperties ...
func (c *Chrome) NavigatorProperties() map[string]interface{} {
	return map[string]interface{}{
		"vendorSub":               "",
		"productSub":              "20030107",
		"vendor":                  "Google Inc.",
		"maxTouchPoints":          0,
		"userActivation":          struct{}{},
		"doNotTrack":              "1",
		"geolocation":             struct{}{},
		"connection":              struct{}{},
		"webkitTemporaryStorage":  struct{}{},
		"webkitPersistentStorage": struct{}{},
		"hardwareConcurrency":     c.cpuCount,
		"cookieEnabled":           true,
		"appCodeName":             "Mozilla",
		"appName":                 "Netscape",
		"appVersion":              "5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/" + chromeVersion + " Safari/537.36",
		"platform":                "Win32",
		"product":                 "Gecko",
		"userAgent":               "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/" + chromeVersion + " Safari/537.36",
		"language":                "en-US",
		"languages":               []string{"en-US"},
		"onLine":                  true,
		"webdriver":               false,
		"pdfViewerEnabled":        true,
		"scheduling":              struct{}{},
		"bluetooth":               struct{}{},
		"clipboard":               struct{}{},
		"credentials":             struct{}{},
		"keyboard":                struct{}{},
		"managed":                 struct{}{},
		"mediaDevices":            struct{}{},
		"storage":                 struct{}{},
		"serviceWorker":           struct{}{},
		"wakeLock":                struct{}{},
		"deviceMemory":            c.memorySize,
		"ink":                     struct{}{},
		"hid":                     struct{}{},
		"locks":                   struct{}{},
		"mediaCapabilities":       struct{}{},
		"mediaSession":            struct{}{},
		"permissions":             struct{}{},
		"presentation":            struct{}{},
		"serial":                  struct{}{},
		"virtualKeyboard":         struct{}{},
		"usb":                     struct{}{},
		"xr":                      struct{}{},
		"userAgentData": map[string]interface{}{
			"brands": []map[string]interface{}{
				{"brand": "Chromium", "version": shortChromeVersion},
				{"brand": "Google Chrome", "version": shortChromeVersion},
				{"brand": ";Not A Brand", "version": "99"},
			},
			"mobile": false,
		},
		"plugins": []string{
			"internal-pdf-viewer",
			"internal-pdf-viewer",
			"internal-pdf-viewer",
			"internal-pdf-viewer",
			"internal-pdf-viewer",
		},
	}
}

// Unix ...
func (c *Chrome) Unix(asMilliseconds bool) int64 {
	t := time.Now().UnixNano() / int64(time.Millisecond)
	t += c.unixOffset
	if !asMilliseconds {
		t /= 1000
	}

	return t
}

// OffsetUnix ...
func (c *Chrome) OffsetUnix(offset int64) {
	c.unixOffset += offset
}

// ResetUnix ...
func (c *Chrome) ResetUnix() {
	if c.unixOffset > 0 {
		time.Sleep(time.Duration(c.unixOffset) * time.Millisecond)
	}
	c.unixOffset = 0
}
