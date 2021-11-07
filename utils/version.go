package utils

import (
	"encoding/json"
	"net/http"
	"strings"
)

var (
	// version is the cached hCaptcha version.
	version string
	// assetVersion is the cached hCaptcha asset version.
	assetVersion string
)

// Version returns the current hCaptcha version.
func Version() string {
	return version
}

// AssetVersion returns the current hCaptcha asset version.
func AssetVersion() string {
	return assetVersion
}

// updateVersion updates the cached version.
func updateVersion() {
	resp, err := http.Get("https://hcaptcha.com/1/api.js")
	if err != nil {
		panic(err)
	}
	version = strings.Split(strings.Split(resp.Request.URL.String(), "v1/")[1], "/")[0]
	defer resp.Body.Close()
}

// updateAssetVersion updates the cached hCaptcha asset version.
func updateAssetVersion() {
	req, err := http.NewRequest(
		"GET",
		"https://hcaptcha.com/checksiteconfig?v="+version+"&host=dashboard.hcaptcha.com"+
			"&sitekey=13257c82-e129-4f09-a733-2a7cb3102832&sc=1&swa=1",
		nil,
	)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.69 Safari/537.36")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		panic(err)
	}

	encodedJWT := data["c"].(map[string]interface{})["req"].(string)
	assetVersion = strings.TrimPrefix(ParseJWT(encodedJWT)["l"].(string), "https://newassets.hcaptcha.com/c/")
}
