package agents

import (
	"github.com/iancoleman/orderedmap"
	"github.com/justtaldevelops/go-hcaptcha/utils"
	"time"
)

// Chrome is the agent for Google Chrome.
type Chrome struct {
	// screenSize is the screen size of the Chrome browser.
	screenSize [2]int
	// availableScreenSize is the available screen size of the Chrome browser.
	availableScreenSize [2]int
	// cpuCount is the number of CPUs allowed to the Chrome browser.
	cpuCount int
	// memorySize is the memory size in gigabytes allowed to the Chrome browser.
	memorySize int
	// unixOffset is the offset of the unix timestamps.
	unixOffset int64
}

// NewChrome creates a new Chrome agent.
func NewChrome() *Chrome {
	possibleScreenSizes := [][][2]int{
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

// UserAgent ...
func (c *Chrome) UserAgent() string {
	return latestChromeAgent
}

// ScreenProperties ...
func (c *Chrome) ScreenProperties() *orderedmap.OrderedMap {
	m := orderedmap.New()
	m.Set("availWidth", c.availableScreenSize[0])
	m.Set("availHeight", c.availableScreenSize[1])
	m.Set("width", c.screenSize[0])
	m.Set("height", c.screenSize[1])
	m.Set("colorDepth", 24)
	m.Set("pixelDepth", 24)
	m.Set("availLeft", 0)
	m.Set("availTop", 0)
	return m
}

// NavigatorProperties ...
func (c *Chrome) NavigatorProperties() *orderedmap.OrderedMap {
	chromium := orderedmap.New()
	chromium.Set("brand", "Chromium")
	chromium.Set("version", shortChromeVersion)

	chrome := orderedmap.New()
	chrome.Set("brand", "Google Chrome")
	chrome.Set("version", shortChromeVersion)

	notAnyBrand := orderedmap.New()
	notAnyBrand.Set("brand", ";Not A Brand")
	notAnyBrand.Set("version", "99")

	userAgentData := orderedmap.New()
	userAgentData.Set("brands", []*orderedmap.OrderedMap{chromium, chrome, notAnyBrand})
	userAgentData.Set("mobile", false)

	m := orderedmap.New()
	m.Set("vendorSub", "")
	m.Set("productSub", "20030107")
	m.Set("vendor", "Google Inc.")
	m.Set("maxTouchPoints", 0)
	m.Set("userActivation", struct{}{})
	m.Set("doNotTrack", "1")
	m.Set("geolocation", struct{}{})
	m.Set("connection", struct{}{})
	m.Set("webkitTemporaryStorage", struct{}{})
	m.Set("webkitPersistentStorage", struct{}{})
	m.Set("hardwareConcurrency", c.cpuCount)
	m.Set("cookieEnabled", true)
	m.Set("appCodeName", "Mozilla")
	m.Set("appName", "Netscape")
	m.Set("appVersion", "5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/"+chromeVersion+" Safari/537.36")
	m.Set("platform", "Win32")
	m.Set("product", "Gecko")
	m.Set("userAgent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/"+chromeVersion+" Safari/537.36")
	m.Set("language", "en-US")
	m.Set("languages", []string{"en-US"})
	m.Set("onLine", true)
	m.Set("webdriver", false)
	m.Set("pdfViewerEnabled", true)
	m.Set("scheduling", struct{}{})
	m.Set("bluetooth", struct{}{})
	m.Set("clipboard", struct{}{})
	m.Set("credentials", struct{}{})
	m.Set("keyboard", struct{}{})
	m.Set("managed", struct{}{})
	m.Set("mediaDevices", struct{}{})
	m.Set("storage", struct{}{})
	m.Set("serviceWorker", struct{}{})
	m.Set("wakeLock", struct{}{})
	m.Set("deviceMemory", c.memorySize)
	m.Set("ink", struct{}{})
	m.Set("hid", struct{}{})
	m.Set("locks", struct{}{})
	m.Set("mediaCapabilities", struct{}{})
	m.Set("mediaSession", struct{}{})
	m.Set("permissions", struct{}{})
	m.Set("presentation", struct{}{})
	m.Set("serial", struct{}{})
	m.Set("virtualKeyboard", struct{}{})
	m.Set("usb", struct{}{})
	m.Set("xr", struct{}{})
	m.Set("userAgentData", userAgentData)
	m.Set("plugins", []string{
		"internal-pdf-viewer",
		"internal-pdf-viewer",
		"internal-pdf-viewer",
		"internal-pdf-viewer",
		"internal-pdf-viewer",
	})

	return m
}

// Unix ...
func (c *Chrome) Unix() int64 {
	t := time.Now().UnixNano() / int64(time.Millisecond)
	t += c.unixOffset

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
