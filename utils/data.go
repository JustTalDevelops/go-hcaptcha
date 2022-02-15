package utils

import (
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"strings"
)

var (
	// FrameSize is the size of the hCaptcha frame.
	FrameSize = [2]int{400, 600}
	// TileImageSize is the size of the tile image.
	TileImageSize = [2]int{123, 123}
	// TileImageStartPosition is the start position of the tile image.
	TileImageStartPosition = [2]int{11, 130}
	// TileImagePadding is the padding between the tile images.
	TileImagePadding = [2]int{5, 6}
	// VerifyButtonPosition is the position of the verify button.
	VerifyButtonPosition = [2]int{314, 559}

	// TilesPerPage is the number of tiles per page.
	TilesPerPage = 9
	// TilesPerRow is the number of tiles per row.
	TilesPerRow = 3

	// VersionRegex to parse our region.
	VersionRegex = regexp.MustCompile("com\\/captcha\\/v1\\/(.*?)\\/")

	// Version is the lastest supported version.
	Version = "44fc726"
	// AssetVersion is the latest supported version of the assets.
	AssetVersion = "4acef65c"
)

// UpdateVersion updates the cached version.
func UpdateVersion() error {
	resp, err := http.Get("https://hcaptcha.com/1/api.js")
	if err != nil {
		return err
	}
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	match := VersionRegex.FindStringSubmatch(string(body))

	Version = match[1]
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)

	return nil
}

// UpdateAssetVersion updates the cached hCaptcha asset version.
func UpdateAssetVersion(siteKey string) error {
	req, err := http.NewRequest(
		"GET",
		"https://hcaptcha.com/checksiteconfig?v="+Version+"&host=dashboard.hcaptcha.com"+
			"&sitekey="+siteKey+"&sc=1&swa=1",
		nil,
	)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.102 Safari/537.36")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return err
	}

	encodedJWT := data["c"].(map[string]interface{})["req"].(string)
	AssetVersion = strings.TrimPrefix(ParseJWT(encodedJWT)["l"].(string), "https://newassets.hcaptcha.com/c/")
	return nil
}
