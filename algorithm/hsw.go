package algorithm

import (
	"fmt"
	"github.com/mxschmitt/playwright-go"
)

// page is used for generating the HSW.
var page playwright.Page

// init initializes Playwright for the HSW algorithm.
func init() {
	pw, err := playwright.Run()
	if err != nil {
		panic(err)
	}
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
	})
	if err != nil {
		panic(err)
	}
	page, err = browser.NewPage()
	if err != nil {
		panic(err)
	}
	_, err = page.AddScriptTag(playwright.PageAddScriptTagOptions{Content: playwright.String(script("hsw.js"))})
	if err != nil {
		panic(err)
	}
}

// HSW is one of a few proof algorithms for hCaptcha services.
type HSW struct{}

// Encode ...
func (h *HSW) Encode() string {
	return "hsw"
}

// Prove ...
func (h HSW) Prove(request string) (string, error) {
	resp, err := page.Evaluate(fmt.Sprintf("hsw(\"%v\")", request))
	if err != nil {
		return "", err
	}
	return resp.(string), nil
}
