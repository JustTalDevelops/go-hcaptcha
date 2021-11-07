package utils

import (
	"net/http"
	"strings"
)

// version is the cached hCaptcha version.
var version string

// Version returns the current hCaptcha version.
func Version() string {
	return version
}

// init initializes the version.
func init() {
	resp, err := http.Get("https://hcaptcha.com/1/api.js")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	version = strings.Split(strings.Split(resp.Request.URL.String(), "v1/")[1], "/")[0]
}
