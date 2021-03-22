package hcaptcha

import (
	"bytes"
	vision "cloud.google.com/go/vision/apiv1"
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

// MotionData contains motion data, just used for JSON requests.
type MotionData struct {
	Start       int64      `json:"st"`
	Destination int64      `json:"dct"`
	Movements   [][3]int64 `json:"mm"`
}

// Solver is an HCaptcha solver instance.
type Solver struct {
	site, siteKey string
	proxies       []string
	hswPool       *HSWPool
	sRand         *rand.Rand
	userAgent     string
	log           *logrus.Logger
	vision        *vision.ImageAnnotatorClient
	client        *http.Client
}

// SolverOptions contains special options that can be applied to new solvers.
type SolverOptions struct {
	// WorkerSize is the amount of workers that should be used. The default is 1.
	WorkerSize int
	// HwsLimit is the limit of HSW in the pool. The default is 3.
	HwsLimit int
	// ScriptUrl is the HSW script URL being used.
	ScriptUrl string
	// SiteKey is the site key of the domain.
	SiteKey string
	// UserAgent is the user agent of the solver.
	UserAgent string
	// Log is the Logrus logger.
	Log *logrus.Logger
}

// Task is a task assigned by HCaptcha.
type Task struct {
	// Image is the image to represent the task.
	Image string `json:"datapoint_uri"`
	// Key is the task key, used when referencing answers.
	Key string `json:"task_key"`
}

// ProxiesEnabled returns true if there are any proxies in the solver.
func (s *Solver) ProxiesEnabled() bool {
	return len(s.proxies) != 0
}

// Solve attempts to solve once until a successful code appears.
// It returns an error if it fails to solve the code before the deadline.
func (s *Solver) Solve(deadline time.Time) (string, error) {
	start := time.Now()
	for {
		var code string
		var err error

		if deadline.After(time.Now()) {
			code, err = s.SolveOnce()
			if err != nil {
				s.log.Error(err)
				continue
			}
			s.log.Info("Solved in less than ", time.Now().Sub(start).Seconds(), " seconds!")
			return code, nil
		} else {
			return "", errors.New("failed to solve captcha before deadline")
		}
	}
}

// SolveOnce attempts to solve once. If successful,
// it returns a code and otherwise returns an error.
func (s *Solver) SolveOnce() (code string, err error) {
	c, err := s.hswPool.GetHSW()
	if err != nil {
		return "", err
	}

	n, err := evaluateHsw(s, c)
	if err != nil {
		return "", err
	}

	timestamp := s.makeTimestamp() + s.randomFromRange(30, 120)
	movements, err := s.getMouseMovements(timestamp)

	motionData := url.Values{}
	motionData.Add("st", strconv.Itoa(int(timestamp)))
	motionData.Add("dct", strconv.Itoa(int(timestamp)))
	motionData.Add("mm", movements)

	form := url.Values{}
	form.Add("sitekey", s.siteKey)
	form.Add("host", s.site)
	form.Add("hl", "en")
	form.Add("motionData", "{}")
	form.Add("n", n)
	form.Add("c", c)

	req, err := http.NewRequest("POST", "https://hcaptcha.com/getcaptcha", strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authority", "hcaptcha.com")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", s.userAgent)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Origin", "https://assets.hcaptcha.com")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	response := gjson.Parse(string(b))
	resp.Body.Close()
	if response.Get("generated_pass_UUID").Exists() {
		return response.Get("generated_pass_UUID").String(), nil
	}

	var tasks []Task
	err = json.Unmarshal([]byte(response.Get("tasklist").String()), &tasks)
	if err != nil {
		return "", errors.New(string(b))
	}

	prompt := strings.Split(response.Get("requester_question").Get("en").String(), " ")

	key := response.Get("key").String()
	job := response.Get("request_type").String()

	timestamp = s.makeTimestamp()

	var motionDataJson MotionData
	motionDataJson.Start = timestamp
	motionDataJson.Destination = timestamp
	motionDataJson.Movements = s.getRawMouseMovements(timestamp)

	b, err = json.Marshal(motionDataJson)
	if err != nil {
		return "", err
	}

	var formJson struct {
		Job        string            `json:"job_mode"`
		Answers    map[string]string `json:"answers"`
		Domain     string            `json:"serverdomain"`
		SiteKey    string            `json:"sitekey"`
		MotionData string            `json:"motionData"`
		N          string            `json:"n"`
		C          string            `json:"c"`
	}

	formJson.Answers = make(map[string]string)
	object := prompt[len(prompt)-1]

	var wg sync.WaitGroup
	for _, t := range tasks {
		if s.vision == nil {
			formJson.Answers[t.Key] = strconv.FormatBool(s.randomTrueFalse())
		} else {
			img := t.Image
			key := t.Key
			wg.Add(1)
			go func() {
				ok, err := s.ImageContainsObject(img, object)
				if err != nil {
					s.log.Error(err)
				}
				s.log.Info(img, " ", object, " ", ok)
				formJson.Answers[key] = strconv.FormatBool(ok)
				wg.Done()
			}()
		}
	}

	wg.Wait()

	n, err = evaluateHsw(s, c)
	if err != nil {
		return "", err
	}

	formJson.Job = job
	formJson.Domain = s.site
	formJson.SiteKey = s.siteKey
	formJson.MotionData = string(b)
	formJson.N = n
	formJson.C = c

	b, err = json.Marshal(formJson)
	if err != nil {
		return "", err
	}

	req, err = http.NewRequest("POST", "https://hcaptcha.com/checkcaptcha/"+key, bytes.NewBuffer(b))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authority", "hcaptcha.com")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("User-Agent", s.userAgent)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "https://assets.hcaptcha.com")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	resp, err = s.client.Do(req)
	if err != nil {
		return "", err
	}

	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	response = gjson.Parse(string(b))
	resp.Body.Close()
	if response.Get("generated_pass_UUID").Exists() {
		return response.Get("generated_pass_UUID").String(), nil
	}

	return "", errors.New(string(b))
}

// Close closes all workers currently running.
func (s *Solver) Close() {
	s.vision.Close()
	s.hswPool.Close()
}

// UpdatePoolUserAgent updates both the pool and the solver's user agents.
func (s *Solver) UpdateAllUserAgents(userAgent string) {
	s.UpdatePoolUserAgent(userAgent)
	s.UpdateUserAgent(userAgent)
}

// UpdatePoolUserAgent updates the pool's user agent.
func (s *Solver) UpdatePoolUserAgent(userAgent string) {
	s.hswPool.userAgent = userAgent
}

// UpdateUserAgent updates the solver's user agent.
func (s *Solver) UpdateUserAgent(userAgent string) {
	s.userAgent = userAgent
}

// randomTrueFalse returns a random boolean to be used in task responses.
func (s *Solver) randomTrueFalse() bool {
	return s.sRand.Intn(2) == 1
}

// randomFromRange generates a random number from the range provided.
func (s *Solver) randomFromRange(min, max int) int64 {
	return int64(s.sRand.Intn(max-min) + min)
}

// makeTimestamp generates a timestamp in unix milliseconds to be sent to HCaptcha.
func (s *Solver) makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// ImageContainsObject checks if an image contains an hCaptcha object.
func (s *Solver) ImageContainsObject(image, object string) (bool, error) {
	if object == "motorbus" { // why hCaptcha... why
		object = "bus"
	}

	resp, err := http.Get(image)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	img, err := vision.NewImageFromReader(resp.Body)
	if err != nil {
		return false, err
	}

	annotations, err := s.vision.LocalizeObjects(context.Background(), img, nil)
	if err != nil {
		return false, err
	}

	for _, annotation := range annotations {
		if strings.Contains(strings.ToLower(annotation.Name), strings.ToLower(object)) && annotation.Score > 0.50 {
			return true, nil
		}
	}

	return false, nil
}

// NewSolver creates a new instance of an HCaptcha solver.
func NewSolver(site string, opts ...SolverOptions) (*Solver, error) {
	if len(opts) == 0 {
		opts = append(opts, SolverOptions{})
	}

	// Default options checks
	if opts[0].WorkerSize == 0 {
		opts[0].WorkerSize = DefaultWorkerAmount
	}

	if opts[0].HwsLimit == 0 {
		opts[0].HwsLimit = DefaultHWSLimit
	}

	if opts[0].ScriptUrl == "" {
		opts[0].ScriptUrl = DefaultScriptUrl
	}

	if opts[0].UserAgent == "" {
		opts[0].UserAgent = DefaultUserAgent
	}

	if opts[0].SiteKey == "" {
		opts[0].SiteKey = uuid.New().String()
	}

	if opts[0].Log == nil {
		opts[0].Log = logrus.New()
		opts[0].Log.Formatter = &logrus.TextFormatter{ForceColors: true}
		opts[0].Log.Level = logrus.DebugLevel
	}

	// Initialize the pool
	pool, err := NewHSWPool(site, opts[0].SiteKey, opts[0].ScriptUrl, opts[0].Log, opts[0].HwsLimit, opts[0].WorkerSize)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		opts[0].Log.Error(err)
		opts[0].Log.Error("You can ignore the above error if you aren't using Vision API.")
	}

	return &Solver{client: &http.Client{}, vision: client, log: opts[0].Log, site: site, siteKey: opts[0].SiteKey, hswPool: pool, sRand: rand.New(rand.NewSource(time.Now().UnixNano())), userAgent: opts[0].UserAgent}, nil
}

// NewSolverWithProxies creates a new instance of an HCaptcha solver, along with proxies.
func NewSolverWithProxies(site string, proxies []string, opts ...SolverOptions) (s *Solver, err error) {
	s, err = NewSolver(site, opts...)
	if err != nil {
		return
	}
	s.proxies = proxies
	for _, w := range s.hswPool.workers {
		p := proxies[rand.Int()%len(proxies)]
		pSplit := strings.Split(p, ":")
		switch len(pSplit) {
		case 4:
			w.client = &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(&url.URL{
				Scheme: "http",
				User:   url.UserPassword(pSplit[2], pSplit[3]),
				Host:   p,
			})}}
		case 2:
			w.client = &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(&url.URL{
				Scheme: "http",
				Host:   p,
			})}}
		default:
			return nil, errors.New("invalid proxy type: must be ip, port, username, and password or ip and port")
		}
	}

	return
}

// getRawMouseMovements returns random mouse movements based on a timestamp.
func (s *Solver) getRawMouseMovements(timestamp int64) (mouseMovements [][3]int64) {
	motionCount := s.randomFromRange(8000, 10000)
	for i := 0; i < int(motionCount); i++ {
		timestamp += s.randomFromRange(0, 10)
		mouseMovements = append(mouseMovements, [3]int64{s.randomFromRange(0, 500), s.randomFromRange(0, 500), timestamp})
	}

	return
}

// getMouseMovements returns random mouse movements based on a timestamp.
func (s *Solver) getMouseMovements(timestamp int64) (string, error) {
	b, err := json.Marshal(s.getRawMouseMovements(timestamp))
	if err != nil {
		return "", err
	}
	return string(b), nil
}
