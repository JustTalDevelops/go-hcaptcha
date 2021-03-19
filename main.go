package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/corpix/uarand"
	"github.com/google/go-querystring/query"
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

func main() {
	rand.Seed(time.Now().UnixNano())
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
	page.AddScriptTag(playwright.PageAddScriptTagOptions{URL: playwright.String("https://assets.hcaptcha.com/c/6043b6da/hsw.js")})
	fmt.Println(tryToSolve("minecraftpocket-servers.com", "e6b7bb01-42ff-4114-9245-3d2b7842ed92", page))
}

func tryToSolve(host, siteKey string, page playwright.Page) (code string, err error) {
	userAgent := uarand.GetRandom()

	hsw, original, err := getHsw(host, siteKey, userAgent, page)
	if err != nil {
		return "", err
	}

	timestamp := makeTimestamp() + randomFromRange(30, 120)
	movements, err := getMouseMovements(timestamp, 90)
	if err != nil {
		return "", err
	}

	var newViewport NewViewport
	newViewport.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Safari/537.36"
	newViewport.ProductSub = "20030107"
	newViewport.Vendor = "Google Inc."
	newViewport.GlobalPrivacyControl = true
	newViewport.HardwareConcurrency = 10
	newViewport.CookieEnabled = true
	newViewport.AppCodeName = "Mozilla"
	newViewport.AppName = "Netscape"
	newViewport.Platform = "Win32"
	newViewport.Product = "Gecko"
	newViewport.AppVersion = "5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Safari/537.36"
	newViewport.Language = "en-US"
	newViewport.Languages = append(newViewport.Languages, "en-US")
	newViewport.Languages = append(newViewport.Languages, "en")
	newViewport.Online = true
	newViewport.DeviceMemory = 4
	newViewport.Plugins = []string{}

	var screen Screen
	screen.AvailableWidth = 1920
	screen.AvailableHeight = 1040
	screen.Width = 1920
	screen.Height = 1080
	screen.ColorDepth = 24
	screen.PixelDepth = 24

	timestamp += randomFromRange(30, 120)

	var topLevel TopLevel
	topLevel.NewViewport = newViewport
	topLevel.Version = 1
	topLevel.Start = timestamp
	topLevel.Screen = screen
	topLevel.Direct = "https://minecraftpocket-servers.com/server/80103/"
	topLevel.XY, err = getXYMovements(timestamp)
	if err != nil {
		return "", err
	}

	timestamp += randomFromRange(30, 120)

	topLevel.WN = generateWN(timestamp)

	timestamp += randomFromRange(30, 120)

	topLevel.Movements, err = getMouseMovements(timestamp, 472)
	if err != nil {
		return "", err
	}

	var motionData MotionData
	motionData.Version = 1
	motionData.Start = timestamp
	motionData.Movements = movements
	motionData.Md = append(motionData.Md, makeMovement(makeTimestamp()+randomFromRange(30, 120)))
	motionData.Mu = append(motionData.Mu, makeMovement(makeTimestamp()+randomFromRange(30, 120)))
	motionData.TopLevel = topLevel
	motionData.Session = []string{}
	motionData.WidgetList = append(motionData.WidgetList, "07r95qrrtxj")
	motionData.WidgetId = "07r95qrrtxj"
	motionData.Href = "https://minecraftpocket-servers.com/server/80103/vote/"

	values, err := query.Values(motionData)
	if err != nil {
		return "", err
	}

	form := url.Values{}
	form.Add("siteKey", siteKey)
	form.Add("host", host)
	form.Add("motionData", values.Encode())
	form.Add("n", hsw)
	form.Add("c", original)
	form.Add("v", "6043b6da")

	req, err := http.NewRequest("POST", "https://hcaptcha.com/getcaptcha", strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Origin", "https://assets.hcaptcha.com")
	req.Header.Set("Referer", "https://assets.hcaptcha.com/")
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
	response := gjson.ParseBytes(b)
	fmt.Println(string(b))
	if response.Get("generated_pass_UUID").Exists() {
		return response.Get("generated_pass_UUID").String(), nil
	}

	var tasks []Task
	err = json.Unmarshal([]byte(response.Get("tasklist").String()), &tasks)
	if err != nil {
		return "", err
	}

	//key := response.Get("key").String()
	//job := response.Get("request_type").String()

	taskResponses := url.Values{}
	for _, t := range tasks {
		taskResponses.Add(t.Key, randomTrueFalse())
	}

	return "", errors.New(string(b))
}

func randomTrueFalse() string {
	return strconv.FormatBool(rand.Intn(2) == 1)
}

func getXYMovements(timestamp int64) (mouseMovements [][]float64, err error) {
	lastMovement := timestamp
	var current float64
	for i := 0; i < 67; i++ {
		lastMovement += randomFromRange(0, 10)
		mouseMovements = append(mouseMovements, []float64{0, current, 0.9745889387144993, float64(timestamp)})
		current += randomFloatFromRange(0, 5)
	}

	return
}

func getMouseMovements(timestamp int64, movementAmount int) (mouseMovements [][]int64, err error) {
	lastMovement := timestamp
	for i := 0; i < movementAmount; i++ {
		lastMovement += randomFromRange(0, 10)
		mouseMovements = append(mouseMovements, makeMovement(lastMovement))
	}

	return
}

func generateWN(timestamp int64) [][]int64 {
	// TODO: Figure out WTF these values are and stop hardcoding them
	return [][]int64{{652, 973, 1, timestamp}}
}

func makeMovement(timestamp int64) []int64 {
	return []int64{randomFromRange(0, 500), randomFromRange(0, 500), timestamp}
}

func randomFloatFromRange(min, max int) float64 {
	return float64(rand.Intn(max-min) + min)
}

func randomFromRange(min, max int) int64 {
	return int64(rand.Intn(max-min) + min)
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
