package hcaptcha

import (
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
	"time"
)

// Solver is an HCaptcha solver instance.
type Solver struct {
	site, siteKey string
	hswScriptUrl  string
	proxies       []string
	hswPool       *HSWPool
	sRand         *rand.Rand
	userAgent     string
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
	for {
		if deadline.After(time.Now()) {
			code, err := s.SolveOnce()
			if err != nil {
				continue
			}
			return code, nil
		} else {
			return "", errors.New("failed to solve captcha before deadline")
		}
	}
}

// SolveOnce attempts to solve once. If successful,
// it returns a code and otherwise returns an error.
func (s *Solver) SolveOnce() (code string, err error) {
	hsw, err := s.hswPool.GetHSW()
	if err != nil {
		return "", err
	}

	timestamp := s.makeTimestamp() + s.randomFromRange(30, 120)
	movements, err := s.getMouseMovements(timestamp)
	if err != nil {
		return "", err
	}

	motionData := url.Values{}
	motionData.Add("st", strconv.Itoa(int(timestamp)))
	motionData.Add("dct", strconv.Itoa(int(timestamp)))
	motionData.Add("mm", movements)

	form := url.Values{}
	form.Add("sitekey", s.siteKey)
	form.Add("host", s.site)
	form.Add("n", hsw.N)
	form.Add("c", hsw.C)
	form.Add("motionData", motionData.Encode())

	req, err := http.NewRequest("POST", "https://hcaptcha.com/getcaptcha", strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", s.userAgent)

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
		taskResponses[t.Key] = []string{s.randomTrueFalse()}
	}

	timestamp = s.makeTimestamp() + s.randomFromRange(30, 120)
	movements, err = s.getMouseMovements(timestamp)
	if err != nil {
		return "", err
	}

	motionData = url.Values{}
	motionData.Add("st", strconv.Itoa(int(timestamp)))
	motionData.Add("dct", strconv.Itoa(int(timestamp)))
	motionData.Add("mm", movements)

	form = url.Values{}
	form.Add("answers", url.Values(taskResponses).Encode())
	form.Add("sitekey", s.siteKey)
	form.Add("serverdomain", s.site)
	form.Add("job_mode", job)
	form.Add("motionData", motionData.Encode())

	req, err = http.NewRequest("POST", "https://hcaptcha.com/checkcaptcha/"+key, strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", s.userAgent)

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

// Close closes all workers currently running.
func (s *Solver) Close() {
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
func (s *Solver) randomTrueFalse() string {
	return strconv.FormatBool(s.sRand.Intn(2) == 1)
}

// getMouseMovements returns random mouse movements based on a timestamp.
func (s *Solver) getMouseMovements(timestamp int64) (string, error) {
	motionCount := s.randomFromRange(1000, 10000)
	var mouseMovements [][3]int64
	for i := 0; i < int(motionCount); i++ {
		timestamp += s.randomFromRange(0, 10)
		mouseMovements = append(mouseMovements, [3]int64{s.randomFromRange(0, 600), s.randomFromRange(0, 600), timestamp})
	}
	b, err := json.Marshal(mouseMovements)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// randomFromRange generates a random number from the range provided.
func (s *Solver) randomFromRange(min, max int) int64 {
	return int64(s.sRand.Intn(max-min) + min)
}

// makeTimestamp generates a timestamp in unix milliseconds to be sent to HCaptcha.
func (s *Solver) makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// NewSolver creates a new instance of an HCaptcha solver.
func NewSolver(site string, workers ...int) (*Solver, error) {
	if len(workers) == 0 {
		workers = append(workers, DefaultWorkerAmount)
	}
	siteKey := uuid.New().String()
	pool, err := NewHSWPool(site, siteKey, DefaultScriptUrl, logrus.New(), workers[0])
	if err != nil {
		return nil, err
	}
	return &Solver{site: site, siteKey: siteKey, hswScriptUrl: DefaultScriptUrl, hswPool: pool, sRand: rand.New(rand.NewSource(time.Now().UnixNano())), userAgent: DefaultUserAgent}, nil
}

// NewSolverWithProxies creates a new instance of an HCaptcha solver, along with proxies.
func NewSolverWithProxies(site string, proxies []string, workers ...int) (s *Solver, err error) {
	s, err = NewSolver(site, workers...)
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
