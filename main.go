package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/corpix/uarand"
	"github.com/mxschmitt/playwright-go"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Task struct {
	Image string `json:"datapoint_uri"`
	Key   string `json:"task_key"`
}

var sRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func main() {
	pw, err := playwright.Run()
	if err != nil {
		panic(err)
	}
	browser, err := pw.Firefox.Launch()
	if err != nil {
		panic(err)
	}
	page, err := browser.NewPage()
	if err != nil {
		panic(err)
	}
	page.AddScriptTag(playwright.PageAddScriptTagOptions{Url: playwright.String("https://assets.hcaptcha.com/c/6043b6da/hsw.js")})

	for {
		code, err := tryToSolve("cd252234-44d8-4b81-b1d5-d4e14b624834", "minecraftpocket-servers.com", page)
		if err != nil {
			continue
		}
		fmt.Println(code)
		break
	}

	page.Close()
	browser.Close()
	pw.Stop()
}

func tryToSolve(sitekey, host string, page playwright.Page) (code string, err error) {
	userAgent := uarand.GetRandom()

	hsw, err := generateHsw(host, sitekey, page)
	if err != nil {
		return "", err
	}

	timestamp := makeTimestamp() + randomFromRange(30, 120)
	movements, err := getMouseMovements(timestamp)
	if err != nil {
		return "", err
	}

	motionData := url.Values{}
	motionData.Add("st", strconv.Itoa(int(timestamp)))
	motionData.Add("dct", strconv.Itoa(int(timestamp)))
	motionData.Add("mm", movements)

	form := url.Values{}
	form.Add("sitekey", sitekey)
	form.Add("host", host)
	form.Add("n", hsw.N)
	form.Add("c", hsw.C)
	form.Add("motionData", motionData.Encode())

	req, err := http.NewRequest("POST", "https://hcaptcha.com/getcaptcha", strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", userAgent)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	response := gjson.Parse(string(b))
	if response.Get("generated_pass_UUID").Exists() {
		return response.Get("generated_pass_UUID").String(), nil
	}

	var tasks []Task
	err = json.Unmarshal([]byte(response.Get("tasklist").String()), &tasks)
	if err != nil {
		return "", err
	}

	key := response.Get("key").String()
	job := response.Get("request_type").String()

	taskResponses := make(map[string][]string)
	for _, t := range tasks {
		taskResponses[t.Key] = []string{randomTrueFalse()}
	}

	timestamp = makeTimestamp() + randomFromRange(30, 120)
	movements, err = getMouseMovements(timestamp)
	if err != nil {
		return "", err
	}

	motionData = url.Values{}
	motionData.Add("st", strconv.Itoa(int(timestamp)))
	motionData.Add("dct", strconv.Itoa(int(timestamp)))
	motionData.Add("mm", movements)

	form = url.Values{}
	form.Add("answers", url.Values(taskResponses).Encode())
	form.Add("sitekey", sitekey)
	form.Add("serverdomain", host)
	form.Add("job_mode", job)
	form.Add("motionData", motionData.Encode())

	req, err = http.NewRequest("POST", "https://hcaptcha.com/checkcaptcha/"+key, strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", userAgent)

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	response = gjson.Parse(string(b))
	if response.Get("generated_pass_UUID").Exists() {
		return response.Get("generated_pass_UUID").String(), nil
	}

	return "", errors.New(string(b))
}

func randomTrueFalse() string {
	return strconv.FormatBool(sRand.Intn(2) == 1)
}

func getMouseMovements(timestamp int64) (string, error) {
	motionCount := randomFromRange(1000, 10000)
	var mouseMovements [][3]int64
	for i := 0; i < int(motionCount); i++ {
		timestamp += randomFromRange(0, 10)
		mouseMovements = append(mouseMovements, [3]int64{randomFromRange(0, 600), randomFromRange(0, 600), timestamp})
	}
	b, err := json.Marshal(mouseMovements)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func randomFromRange(min, max int) int64 {
	return int64(sRand.Intn(max-min) + min)
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
