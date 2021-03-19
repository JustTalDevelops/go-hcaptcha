package main

// Screen contains the screen data for HCaptcha motion data requests.
type Screen struct {
	// AvailableWidth is the available width that the website can use.
	AvailableWidth int `json:"availWidth"`
	// AvailableHeight is the available height that the website can use.
	AvailableHeight int `json:"availHeight"`
	// Width is the width of the user's device.
	Width int `json:"width"`
	// Height is the height of the user's device.
	Height int `json:"height"`
	// ColorDepth is the color depth of the user's screen.
	ColorDepth int `json:"colorDepth"`
	// PixelDepth is the pixel depth of the user's screen.
	PixelDepth int `json:"pixelDepth"`
	// AvailableLeft is all the screen space available on the left side of the screen.
	AvailableLeft int `json:"availLeft"`
	// AvailableTop is all the screen space available on the top of the screen.
	AvailableTop int `json:"availTop"`
}

// NewViewport contains viewport, device, and user specific data.
// It is unclear why the HCaptcha developers decided to name it NewViewport, as it contains other data.
type NewViewport struct {
	// VendorSub ...
	VendorSub string `json:"vendorSub"`
	// ProductSub ...
	ProductSub string `json:"productSub"`
	// Vendor is the vendor of the user's browser.
	Vendor string `json:"vendor"`
	// MaxTouchPoints contains the maximum touch points for a device.
	// If the user is on a PC, this value is always 0.
	MaxTouchPoints int `json:"maxTouchPoints"`
	// UserActivation ...
	UserActivation struct{} `json:"userActivation"`
	// Brave is a value only present if the user is using the Brave browser.
	Brave struct{} `json:"brave"`
	// GlobalPrivacyControl ...
	GlobalPrivacyControl bool `json:"globalPrivacyControl"`
	// DoNotTrack contains tracking options that HCaptcha should avoid.
	DoNotTrack []string `json:"doNotTrack"`
	// Geolocation ...
	Geolocation struct{} `json:"geolocation"`
	// Connection ...
	Connection struct{} `json:"connection"`
	// WebKitTemporaryStorage ...
	WebKitTemporaryStorage struct{} `json:"webKitTemporaryStorage"`
	// WebKitPersistentStorage ...
	WebKitPersistentStorage struct{} `json:"webKitPersistentStorage"`
	// HardwareConcurrency ...
	HardwareConcurrency int `json:"hardwareConcurrency"`
	// CookieEnabled is a boolean only true if cookies are enabled on the browser.
	CookieEnabled bool `json:"cookieEnabled"`
	// AppCodeName is usually the author of the browser. (e.g. Mozilla)
	AppCodeName string `json:"appCodeName"`
	// AppName is the name of the browser. (e.g. Firefox)
	AppName string `json:"appName"`
	// AppVersion is the version of the browser, which is just split from the UserAgent.
	// It is unclear why they send AppVersion when it can just be pulled from the UserAgent.
	AppVersion string `json:"appVersion"`
	// Platform ...
	Platform string `json:"platform"`
	// Product ...
	Product string `json:"product"`
	// UserAgent is the user agent of the user.
	UserAgent string `json:"userAgent"`
	// Language is the language the user's system is configured on.
	Language string `json:"language"`
	// Languages contain all languages on the user's system.
	Languages []string `json:"languages"`
	// Online ...
	Online bool `json:"onLine"`
	// Webdriver ...
	Webdriver bool `json:"webdriver"`
	// MediaCapabilities ...
	MediaCapabilities struct{} `json:"mediaCapabilities"`
	// Permissions ...
	Permissions struct{} `json:"permissions"`
	// Locks ...
	Locks struct{} `json:"locks"`
	// WakeLock ...
	WakeLock struct{} `json:"wakeLock"`
	// USB ...
	USB struct{} `json:"usb"`
	// MediaSession ...
	MediaSession struct{} `json:"mediaSession"`
	// Clipboard ...
	Clipboard struct{} `json:"clipboard"`
	// Credentials ...
	Credentials struct{} `json:"credentials"`
	// Keyboard ...
	Keyboard struct{} `json:"keyboard"`
	// MediaDevices ...
	MediaDevices struct{} `json:"mediaDevices"`
	// Storage ...
	Storage struct{} `json:"storage"`
	// ServiceWorker ...
	ServiceWorker struct{} `json:"serviceWorker"`
	// DeviceMemory is the amount of memory allocated to the browser (in gigabytes).
	// It is unclear why this field is marked as device memory, as it only uses browser memory.
	DeviceMemory int `json:"deviceMemory"`
	// Hid ...
	Hid struct{} `json:"hid"`
	// Presentation ...
	Presentation struct{} `json:"presentation"`
	// Bluetooth ...
	Bluetooth struct{} `json:"bluetooth"`
	// XR ...
	XR struct{} `json:"xr"`
	// Plugins contain all plugins installed on the user's browser.
	Plugins []string `json:"plugins"`
}

// TopLevel is the top level data. It is unclear what the name means.
type TopLevel struct {
	// Version is the version of TopLevel data being sent.
	// Usually, this is just 1.
	Version int `json:"v"`
	// Start is the time that the TopLevel data was initially recorded.
	// This is in Unix seconds.
	Start int64 `json:"st"`
	// Screen contains the screen data of the user's device and browser.
	Screen Screen `json:"sc"`
	// NewViewport contains viewport, device, and user data.
	NewViewport NewViewport `json:"nv"`
	// Direct is the direct link to the page the user is on.
	Direct string `json:"dr"`
	// WN contains 1 piece of movement data. It is unclear what this does.
	WN [][]int64 `json:"wn"`
	// WN contains movements across the X and Y coordinates.
	XY [][]float64 `json:"xy"`
	// Movements contain movement data. The size is usually around 500.
	Movements [][]int64 `json:"mm"`
}

// MotionData contains motion data from the user before starting the captcha.
type MotionData struct {
	// Version is the version of MotionData being sent.
	// Usually, this is just 1.
	Version int `json:"v"`
	// Start is the time that the MotionData was initially recorded.
	// This is in Unix seconds.
	Start int64 `json:"st"`
	// Movements contain movement data. The size is usually around 90.
	Movements [][]int64 `json:"mm"`
	// WN contains 1 piece of movement data. It is unclear what this does.
	Md [][]int64 `json:"md"`
	// WN contains 1 piece of movement data. It is unclear what this does.
	Mu [][]int64 `json:"mu"`
	// TopLevel contains the TopLevel data of the request.
	TopLevel TopLevel `json:"topLevel"`
	// Session ...
	Session []string `json:"session"`
	// WidgetList contains a list of widgets enabled on the target site.
	// It is unclear what this does or what these widgets are.
	WidgetList []string `json:"widgetList"`
	// WidgetId ...
	WidgetId string `json:"widgetId"`
	// Href ...
	Href string `json:"href"`
}
