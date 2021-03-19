package main

type Screen struct {
	AvailableWidth int `json:"availWidth"`
	AvailableHeight int `json:"availHeight"`
	Width int `json:"width"`
	Height int `json:"height"`
	ColorDepth int `json:"colorDepth"`
	PixelDepth int `json:"pixelDepth"`
	AvailLeft int `json:"availLeft"`
	AvailTop int `json:"availTop"`
}

type NewViewport struct {
	VendorSub string `json:"vendorSub"`
	ProductSub string `json:"productSub"`
	Vendor string `json:"vendor"`
	MaxTouchPoints int `json:"maxTouchPoints"`
	UserActivation struct{} `json:"userActivation"`
	Brave struct{} `json:"brave"`
	GlobalPrivacyControl bool `json:"globalPrivacyControl"`
	DoNotTrack []string `json:"doNotTrack"`
	Geolocation struct{} `json:"geolocation"`
	Connection struct{} `json:"connection"`
	WebKitTemporaryStorage struct{} `json:"webKitTemporaryStorage"`
	WebKitPersistentStorage struct{} `json:"webKitPersistentStorage"`
	HardwareConcurrency int `json:"hardwareConcurrency"`
	CookieEnabled bool `json:"cookieEnabled"`
	AppCodeName string `json:"appCodeName"`
	AppName string `json:"appName"`
	AppVersion string `json:"appVersion"`
	Platform string `json:"platform"`
	Product string `json:"product"`
	UserAgent string `json:"userAgent"`
	Language string `json:"language"`
	Languages []string `json:"languages"`
	Online bool `json:"onLine"`
	Webdriver bool `json:"webdriver"`
	MediaCapabilities struct{} `json:"mediaCapabilities"`
	Permissions struct{} `json:"permissions"`
	Locks struct{} `json:"locks"`
	WakeLock struct{} `json:"wakeLock"`
	USB struct{} `json:"usb"`
	MediaSession struct{} `json:"mediaSession"`
	Clipboard struct{} `json:"clipboard"`
	Credentials struct{} `json:"credentials"`
	Keyboard struct{} `json:"keyboard"`
	MediaDevices struct{} `json:"mediaDevices"`
	Storage struct{} `json:"storage"`
	ServiceWorker struct{} `json:"serviceWorker"`
	DeviceMemory int `json:"deviceMemory"`
	Hid struct{} `json:"hid"`
	Presentation struct{} `json:"presentation"`
	Bluetooth struct{} `json:"bluetooth"`
	XR struct{} `json:"xr"`
	Plugins []string `json:"plugins"`
}

type TopLevel struct {
	Version int `json:"v"`
	Start int64 `json:"st"`
	Screen Screen `json:"sc"`
	NewViewport NewViewport `json:"nv"`
	Direct string `json:"dr"`
	WN [][]int64 `json:"wn"`
	XY [][]float64 `json:"xy"`
	Movements [][]int64 `json:"mm"`
}

type MotionData struct {
	Version int `json:"v"`
	Start int64 `json:"st"`
	Movements [][]int64 `json:"mm"`
	Md [][]int64 `json:"md"`
	Mu [][]int64 `json:"mu"`
	TopLevel TopLevel `json:"topLevel"`
	Session []string `json:"session"`
	WidgetList []string `json:"widgetList"`
	WidgetId string `json:"widgetId"`
	Href string `json:"href"`
}