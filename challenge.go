package hcaptcha

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/iancoleman/orderedmap"
	"github.com/justtaldevelops/go-hcaptcha/agents"
	"github.com/justtaldevelops/go-hcaptcha/algorithm"
	"github.com/justtaldevelops/go-hcaptcha/screen"
	"github.com/justtaldevelops/go-hcaptcha/utils"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// Challenge is an hCaptcha challenge.
type Challenge struct {
	host, siteKey      string
	url, widgetID      string
	id, token          string
	category, question string
	tasks              []Task
	log                *logrus.Logger
	agent              agents.Agent
	proof              algorithm.Proof
	top, frame         *EventRecorder
}

// ChallengeOptions contains special options that can be applied to new solvers.
type ChallengeOptions struct {
	// Logger is the logger to use for logging.
	Logger *logrus.Logger
	// Proxies is a list of proxies to use for solving.
	Proxies []string
}

// Task is a task assigned by hCaptcha.
type Task struct {
	// Image is the image to represent the task.
	Image []byte
	// Key is the task key, used when referencing answers.
	Key string
	// Index is the index of the task.
	Index int
}

// basicChallengeOptions is a set of default options for a basic solver.
func basicChallengeOptions(options *ChallengeOptions) {
	if options.Logger == nil {
		options.Logger = logrus.New()
		options.Logger.Formatter = &logrus.TextFormatter{ForceColors: true}
		options.Logger.Level = logrus.DebugLevel
	}
}

// NewChallenge creates a new hCaptcha challenge.
func NewChallenge(url, siteKey string, opts ...ChallengeOptions) (*Challenge, error) {
	if len(opts) == 0 {
		opts = append(opts, ChallengeOptions{})
	}

	options := opts[0]
	basicChallengeOptions(&options)

	c := &Challenge{
		host:     strings.ToLower(strings.Split(strings.Split(url, "://")[1], "/")[0]),
		siteKey:  siteKey,
		url:      url,
		log:      options.Logger,
		widgetID: utils.WidgetID(),
		agent:    agents.NewChrome(),
	}
	c.agent.OffsetUnix(-10)
	c.setupFrames()

	c.log.Debug("Verifying site configuration...")
	err := c.siteConfig()
	if err != nil {
		return nil, err
	}
	c.log.Info("Requesting captcha...")
	err = c.requestCaptcha()
	if err != nil {
		return nil, err
	}

	return c, nil
}

// Solve solves the challenge with the provided solver.
func (c *Challenge) Solve(solver Solver) error {
	c.log.Debugf("Solving challenge with %T...", solver)
	if len(c.token) > 0 {
		return nil
	}

	split := strings.Split(c.question, " ")
	object := strings.Replace(strings.Replace(split[len(split)-1], "motorbus", "bus", 1), "airplane", "aeroplane", 1)

	c.log.Debugf(`The type of challenge is "%v"`, c.category)
	c.log.Debugf(`The target object is "%v"`, object)

	answers := solver.Solve(c.category, object, c.tasks)

	c.log.Debugf("Decided on %v/%v of the tasks given!", len(answers), len(c.tasks))
	c.log.Debug("Simulating mouse movements on tiles...")

	c.simulateMouseMovements(answers)
	c.agent.ResetUnix()

	answersAsMap := orderedmap.New()
	for _, answer := range c.tasks {
		answersAsMap.Set(answer.Key, strconv.FormatBool(c.answered(answer, answers)))
	}

	motionData := orderedmap.New()
	frameData := c.frame.Data()
	for _, key := range frameData.Keys() {
		value, _ := frameData.Get(key)
		motionData.Set(key, value)
	}
	motionData.Set("topLevel", c.top.Data())
	motionData.Set("v", 1)

	encodedMotionData, err := json.Marshal(motionData)
	if err != nil {
		return err
	}

	m := orderedmap.New()
	m.Set("v", utils.Version())
	m.Set("job_mode", c.category)
	m.Set("answers", answersAsMap)
	m.Set("serverdomain", c.host)
	m.Set("sitekey", c.siteKey)
	m.Set("motionData", string(encodedMotionData))
	m.Set("n", c.proof.Proof)
	m.Set("c", c.proof.Request)

	b, err := json.Marshal(m)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "https://hcaptcha.com/checkcaptcha/"+c.id+"?s="+c.siteKey, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Authority", "hcaptcha.com")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("User-Agent", c.agent.UserAgent())
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "https://newassets.hcaptcha.com")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	response := gjson.ParseBytes(b)
	if !response.Get("pass").Bool() {
		return fmt.Errorf("submit request was rejected")
	}

	c.log.Info("Successfully completed challenge!")
	c.token = response.Get("generated_pass_UUID").String()
	return nil
}

// Tasks returns the tasks for the challenge.
func (c *Challenge) Tasks() []Task {
	return c.tasks
}

// Logger returns the logger for the challenge.
func (c *Challenge) Logger() *logrus.Logger {
	return c.log
}

// Category returns the category of the challenge.
func (c *Challenge) Category() string {
	return c.category
}

// Question returns the question of the challenge.
func (c *Challenge) Question() string {
	return c.question
}

// Token returns the response token. This is only valid once the challenge has been solved.
func (c *Challenge) Token() string {
	return c.token
}

// setupFrames sets up the frames for the challenge.
func (c *Challenge) setupFrames() {
	c.top = NewEventRecorder(c.agent)
	c.top.Record()
	c.top.SetData("dr", "")
	c.top.SetData("inv", false)
	c.top.SetData("sc", c.agent.ScreenProperties())
	c.top.SetData("nv", c.agent.NavigatorProperties())
	c.top.SetData("exec", false)
	c.agent.OffsetUnix(int64(utils.Between(200, 400)))
	c.frame = NewEventRecorder(c.agent)
	c.frame.Record()
}

// simulateMouseMovements simulates mouse movements for the hCaptcha API.
func (c *Challenge) simulateMouseMovements(answers []Task) {
	totalPages := int(math.Max(1, float64(len(c.tasks)/utils.TilesPerPage)))
	cursorPos := screen.Point{X: float64(utils.Between(1, 5)), Y: float64(utils.Between(300, 350))}

	rightBoundary := utils.FrameSize[0]
	upBoundary := utils.FrameSize[1]
	opts := &screen.CurveOpts{
		RightBoundary: &rightBoundary,
		UpBoundary:    &upBoundary,
	}

	for page := 0; page < totalPages; page++ {
		pageTiles := c.tasks[page*utils.TilesPerPage : (page+1)*utils.TilesPerPage]
		for _, tile := range pageTiles {
			if !c.answered(tile, answers) {
				continue
			}

			tilePos := screen.Point{
				X: float64((utils.TileImageSize[0] * tile.Index % utils.TilesPerRow) +
					utils.TileImagePadding[0]*tile.Index%utils.TilesPerRow +
					utils.Between(10, utils.TileImageSize[0]) +
					utils.TileImageStartPosition[0]),
				Y: float64((utils.TileImageSize[1] * tile.Index % utils.TilesPerRow) +
					utils.TileImagePadding[1]*tile.Index%utils.TilesPerRow +
					utils.Between(10, utils.TileImageSize[1]) +
					utils.TileImageStartPosition[1]),
			}

			movements := c.generateMouseMovements(cursorPos, tilePos, opts)
			lastMovement := movements[len(movements)-1]
			for _, move := range movements {
				c.frame.RecordEvent(Event{Type: "mm", Point: move.point, Timestamp: move.timestamp})
			}
			// TODO: Add a delay for movement up and down.
			c.frame.RecordEvent(Event{Type: "md", Point: lastMovement.point, Timestamp: lastMovement.timestamp})
			c.frame.RecordEvent(Event{Type: "mu", Point: lastMovement.point, Timestamp: lastMovement.timestamp})
			cursorPos = tilePos
		}

		buttonPos := screen.Point{
			X: float64(utils.VerifyButtonPosition[0] + utils.Between(5, 50)),
			Y: float64(utils.VerifyButtonPosition[1] + utils.Between(5, 15)),
		}

		movements := c.generateMouseMovements(cursorPos, buttonPos, opts)
		lastMovement := movements[len(movements)-1]
		for _, move := range movements {
			c.frame.RecordEvent(Event{Type: "mm", Point: move.point, Timestamp: move.timestamp})
		}
		c.frame.RecordEvent(Event{Type: "md", Point: lastMovement.point, Timestamp: lastMovement.timestamp})
		c.frame.RecordEvent(Event{Type: "mu", Point: lastMovement.point, Timestamp: lastMovement.timestamp})
		cursorPos = buttonPos
	}
}

// movement is a mouse movement.
type movement struct {
	// point is the mouse's position at the timestamp.
	point screen.Point
	// timestamp is the timestamp of the movement.
	timestamp int64
}

// generateMouseMovements generates mouse movements for simulateMouseMovements.
func (c *Challenge) generateMouseMovements(fromPoint, toPoint screen.Point, opts *screen.CurveOpts) []movement {
	curve := screen.NewHumanCurve(fromPoint, toPoint, opts)
	points := curve.Points()

	resultMovements := make([]movement, len(points))
	for _, point := range points {
		c.agent.OffsetUnix(int64(utils.Between(2, 5)))
		resultMovements = append(resultMovements, movement{point: point, timestamp: c.agent.Unix()})
	}
	return resultMovements
}

// answered returns true if the task provided is in the answers slice.
func (c *Challenge) answered(task Task, answers []Task) bool {
	for _, answer := range answers {
		if answer.Key == task.Key {
			return true
		}
	}
	return false
}

// requestCaptcha gets the captcha from the site.
func (c *Challenge) requestCaptcha() error {
	prev := orderedmap.New()
	prev.Set("escaped", false)
	prev.Set("passed", false)
	prev.Set("expiredChallenge", false)
	prev.Set("expiredResponse", false)

	motionData := orderedmap.New()
	motionData.Set("v", 1)

	frameData := c.frame.Data()
	for _, key := range frameData.Keys() {
		value, _ := frameData.Get(key)
		motionData.Set(key, value)
	}

	motionData.Set("topLevel", c.top.Data())
	motionData.Set("session", struct{}{})
	motionData.Set("widgetList", []string{c.widgetID})
	motionData.Set("widgetId", c.widgetID)
	motionData.Set("href", c.url)
	motionData.Set("prev", prev)

	encodedMotionData, err := json.Marshal(motionData)
	if err != nil {
		return err
	}

	form := url.Values{}
	form.Add("v", utils.Version())
	form.Add("sitekey", c.siteKey)
	form.Add("host", c.host)
	form.Add("hl", "en")
	form.Add("motionData", string(encodedMotionData))
	form.Add("n", c.proof.Proof)
	form.Add("c", c.proof.Request)

	req, err := http.NewRequest("POST", "https://hcaptcha.com/getcaptcha?s="+c.siteKey, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Authority", "hcaptcha.com")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.agent.UserAgent())
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Origin", "https://newassets.hcaptcha.com")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	response := gjson.ParseBytes(b)
	if response.Get("pass").Exists() {
		c.token = response.Get("generated_pass_UUID").String()
		return nil
	}

	success := response.Get("success")
	if success.Exists() && !success.Bool() {
		return fmt.Errorf("challenge creation request was rejected")
	}

	c.id = response.Get("key").String()
	c.category = response.Get("request_type").String()
	c.question = response.Get("requester_question").Get("en").String()

	for index, task := range response.Get("tasklist").Array() {
		resp, err = http.Get(task.Get("datapoint_uri").String())
		if err != nil {
			return err
		}
		b, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		_ = resp.Body.Close()

		c.tasks = append(c.tasks, Task{
			Image: b,
			Key:   task.Get("task_key").String(),
			Index: index,
		})
	}

	request := response.Get("c")
	c.proof, err = algorithm.Solve(request.Get("type").String(), request.Get("req").String())
	if err != nil {
		return err
	}
	return nil
}

// siteConfig verifies a site configuration and returns proof of work for the given challenge.
func (c *Challenge) siteConfig() error {
	req, err := http.NewRequest("GET", "https://hcaptcha.com/checksiteconfig?v="+utils.Version()+"&host="+c.host+"&sitekey="+c.siteKey+"&sc=1&swa=1", nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", c.agent.UserAgent())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	response := gjson.ParseBytes(b)
	if !response.Get("pass").Bool() {
		return fmt.Errorf("site key is invalid")
	}

	request := response.Get("c")
	c.proof, err = algorithm.Solve(request.Get("type").String(), request.Get("req").String())
	if err != nil {
		return err
	}
	return nil
}
