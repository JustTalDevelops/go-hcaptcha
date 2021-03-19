package main

import (
	"github.com/mxschmitt/playwright-go"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
)

func getHsw(host, sitekey, userAgent string, page playwright.Page) (hsw string, original string, err error) {
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
