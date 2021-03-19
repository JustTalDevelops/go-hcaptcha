package main

import (
	"github.com/mxschmitt/playwright-go"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
)

// getHsw sends a request to the HCaptcha site config system for a HSW token.
// Then, we use the original token provided to us to generate HSW to send in our captcha requests.
// The HSW is generated through Playwright. On initial startup, Playwright (running Firefox, WebKit seems to fail)
// opens a new empty tab, injects the HSW script by HCaptcha, and then evaluates the HSW using the page evaluate
// function. It then returns the response, plus the HSW token and an error in case anything went wrong.
func getHsw(host, sitekey, userAgent string, page playwright.Page) (hsw string, token string, err error) {
	req, err := http.NewRequest("GET", "https://hcaptcha.com/checksiteconfig?host="+host+"&sitekey="+sitekey+"&sc=1&swa=1", nil)
	if err != nil {
		return "", "", err
	}
	req.Header.Set("User-Agent", userAgent)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}
	result := gjson.ParseBytes(b)
	pResp, err := page.Evaluate("hsw(\"" + result.Get("c.req").String() + "\");")
	if err != nil {
		return "", "", err
	}
	return pResp.(string), result.Get("c").String(), nil
}
