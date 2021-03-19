package main

// Screen contains the screen data for HCaptcha motion data requests.
type Screen struct {
	// AvailableWidth is the available width that the website can use.
	AvailableWidth int `url:"availWidth"`
	// AvailableHeight is the available height that the website can use.
	AvailableHeight int `url:"availHeight"`
	// Width is the width of the user's device.
	Width int `url:"width"`
	// Height is the height of the user's device.
	Height int `url:"height"`
	// ColorDepth is the color depth of the user's screen.
	ColorDepth int `url:"colorDepth"`
	// PixelDepth is the pixel depth of the user's screen.
	PixelDepth int `url:"pixelDepth"`
	// AvailableLeft is all the screen space available on the left side of the screen.
	AvailableLeft int `url:"availLeft"`
	// AvailableTop is all the screen space available on the top of the screen.
	AvailableTop int `url:"availTop"`
}

// NewViewport contains viewport, device, and user specific data.
// It is unclear why the HCaptcha developers decided to name it NewViewport, as it contains other data.
type NewViewport struct {
	// VendorSub ...
	VendorSub string `url:"vendorSub"`
	// ProductSub ...
	ProductSub string `url:"productSub"`
	// Vendor is the vendor of the user's browser.
	Vendor string `url:"vendor"`
	// MaxTouchPoints contains the maximum touch points for a device.
	// If the user is on a PC, this value is always 0.
	MaxTouchPoints int `url:"maxTouchPoints"`
	// UserActivation ...
	UserActivation struct{} `url:"userActivation"`
	// Brave is a value only present if the user is using the Brave browser.
	Brave struct{} `url:"brave"`
	// GlobalPrivacyControl ...
	GlobalPrivacyControl bool `url:"globalPrivacyControl"`
	// DoNotTrack contains tracking options that HCaptcha should avoid.
	DoNotTrack []string `url:"doNotTrack"`
	// Geolocation ...
	Geolocation struct{} `url:"geolocation"`
	// Connection ...
	Connection struct{} `url:"connection"`
	// WebKitTemporaryStorage ...
	WebKitTemporaryStorage struct{} `url:"webKitTemporaryStorage"`
	// WebKitPersistentStorage ...
	WebKitPersistentStorage struct{} `url:"webKitPersistentStorage"`
	// HardwareConcurrency ...
	HardwareConcurrency int `url:"hardwareConcurrency"`
	// CookieEnabled is a boolean only true if cookies are enabled on the browser.
	CookieEnabled bool `url:"cookieEnabled"`
	// AppCodeName is usually the author of the browser. (e.g. Mozilla)
	AppCodeName string `url:"appCodeName"`
	// AppName is the name of the browser. (e.g. Firefox)
	AppName string `url:"appName"`
	// AppVersion is the version of the browser, which is just split from the UserAgent.
	// It is unclear why they send AppVersion when it can just be pulled from the UserAgent.
	AppVersion string `url:"appVersion"`
	// Platform ...
	Platform string `url:"platform"`
	// Product ...
	Product string `url:"product"`
	// UserAgent is the user agent of the user.
	UserAgent string `url:"userAgent"`
	// Language is the language the user's system is configured on.
	Language string `url:"language"`
	// Languages contain all languages on the user's system.
	Languages []string `url:"languages"`
	// Online ...
	Online bool `url:"onLine"`
	// Webdriver ...
	Webdriver bool `url:"webdriver"`
	// MediaCapabilities ...
	MediaCapabilities struct{} `url:"mediaCapabilities"`
	// Permissions ...
	Permissions struct{} `url:"permissions"`
	// Locks ...
	Locks struct{} `url:"locks"`
	// WakeLock ...
	WakeLock struct{} `url:"wakeLock"`
	// USB ...
	USB struct{} `url:"usb"`
	// MediaSession ...
	MediaSession struct{} `url:"mediaSession"`
	// Clipboard ...
	Clipboard struct{} `url:"clipboard"`
	// Credentials ...
	Credentials struct{} `url:"credentials"`
	// Keyboard ...
	Keyboard struct{} `url:"keyboard"`
	// MediaDevices ...
	MediaDevices struct{} `url:"mediaDevices"`
	// Storage ...
	Storage struct{} `url:"storage"`
	// ServiceWorker ...
	ServiceWorker struct{} `url:"serviceWorker"`
	// DeviceMemory is the amount of memory allocated to the browser (in gigabytes).
	// It is unclear why this field is marked as device memory, as it only uses browser memory.
	DeviceMemory int `url:"deviceMemory"`
	// Hid ...
	Hid struct{} `url:"hid"`
	// Presentation ...
	Presentation struct{} `url:"presentation"`
	// Bluetooth ...
	Bluetooth struct{} `url:"bluetooth"`
	// XR ...
	XR struct{} `url:"xr"`
	// Plugins contain all plugins installed on the user's browser.
	Plugins []string `url:"plugins"`
}

// TopLevel is the top level data. It is unclear what the name means.
type TopLevel struct {
	// Version is the version of TopLevel data being sent.
	// Usually, this is just 1.
	Version int `url:"v"`
	// Start is the time that the TopLevel data was initially recorded.
	// This is in Unix seconds.
	Start int64 `url:"st"`
	// Screen contains the screen data of the user's device and browser.
	Screen Screen `url:"sc"`
	// NewViewport contains viewport, device, and user data.
	NewViewport NewViewport `url:"nv"`
	// Direct is the direct link to the page the user is on.
	Direct string `url:"dr"`
	// WN contains 1 piece of movement data. It is unclear what this does.
	WN [][]int64 `url:"wn"`
	// WN contains movements across the X and Y coordinates.
	XY [][]float64 `url:"xy"`
	// Movements contain movement data. The size is usually around 500.
	Movements [][]int64 `url:"mm"`
}

// MotionData contains motion data from the user before starting the captcha.
type MotionData struct {
	// Version is the version of MotionData being sent.
	// Usually, this is just 1.
	Version int `url:"v"`
	// Start is the time that the MotionData was initially recorded.
	// This is in Unix seconds.
	Start int64 `url:"st"`
	// Movements contain movement data. The size is usually around 90.
	Movements [][]int64 `url:"mm"`
	// WN contains 1 piece of movement data. It is unclear what this does.
	Md [][]int64 `url:"md"`
	// WN contains 1 piece of movement data. It is unclear what this does.
	Mu [][]int64 `url:"mu"`
	// TopLevel contains the TopLevel data of the request.
	TopLevel TopLevel `url:"topLevel"`
	// Session ...
	Session []string `url:"session"`
	// WidgetList contains a list of widgets enabled on the target site.
	// It is unclear what this does or what these widgets are.
	WidgetList []string `url:"widgetList"`
	// WidgetId ...
	WidgetId string `url:"widgetId"`
	// Href ...
	Href string `url:"href"`
}
