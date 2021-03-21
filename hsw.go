package hcaptcha

import (
	"errors"
	"github.com/mxschmitt/playwright-go"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"sync"
)

// HSWWorker is a worker for an HSW pool.
// Workers should always be created through the NewWorker function.
type HSWWorker struct {
	playwrightInstance *playwright.Playwright
	playwrightBrowser  playwright.Browser
	playwrightPage     playwright.Page
	pool               *HSWPool
	running            bool
	client             *http.Client
}

// Close closes all Playwright instances, and set's the running status to false to close HSW generation.
func (h *HSWWorker) Close() {
	h.playwrightPage.Close()
	h.playwrightBrowser.Close()
	h.playwrightInstance.Stop()

	h.running = false
}

// NewWorker returns a new HSW worker to be used in HSW pools.
func NewWorker(pool *HSWPool) (*HSWWorker, error) {
	var err error
	worker := &HSWWorker{pool: pool, running: true, client: http.DefaultClient}

	worker.playwrightInstance, err = playwright.Run()
	if err != nil {
		return nil, err
	}
	worker.playwrightBrowser, err = worker.playwrightInstance.Firefox.Launch()
	if err != nil {
		return nil, err
	}
	worker.playwrightPage, err = worker.playwrightBrowser.NewPage()
	if err != nil {
		return nil, err
	}
	_, err = worker.playwrightPage.AddScriptTag(playwright.PageAddScriptTagOptions{Url: pool.hswScriptUrl})
	if err != nil {
		return nil, err
	}

	go func() {
		var c string
		for {
			if len(pool.hswPool) == pool.poolLimit {
				// We don't want to overfill the pool with entries that it is never going to touch, so we continue until the pool shrinks.
				continue
			}
			c, err = generateHsw(worker)

			// We run the check after the HSW has been generated,
			// because we are more likely to be already generating HSW
			// when a shutdown is requested.
			if !worker.running {
				break
			}
			if err != nil {
				continue
			}
			pool.hswPoolMutex.Lock()
			pool.hswPool = append(pool.hswPool, c)
			pool.hswPoolMutex.Unlock()
		}
	}()

	return worker, nil
}

// HSWPool is a pool of HSW tokens that get regenerated on the fly by HSW workers.
// These are used to generate valid captcha codes.
type HSWPool struct {
	hswPool       []string
	hswPoolMutex  *sync.Mutex
	workers       []*HSWWorker
	hswScriptUrl  *string
	host, siteKey string
	log           *logrus.Logger
	userAgent     string
	poolLimit     int
}

// NewHSWPool creates a new HSW pool with the amount of workers specified.
func NewHSWPool(host, siteKey, hswScriptUrl string, log *logrus.Logger, poolLimit, workers int) (pool *HSWPool, err error) {
	if workers <= 0 {
		return nil, errors.New("invalid worker amount for pool")
	}
	pool = &HSWPool{
		hswPoolMutex: &sync.Mutex{},
		hswScriptUrl: playwright.String(hswScriptUrl),
		host:         host,
		siteKey:      siteKey,
		log:          log,
		userAgent:    DefaultUserAgent,
		poolLimit:    poolLimit,
	}

	var worker *HSWWorker
	for i := 0; i < workers; i++ {
		worker, err = NewWorker(pool)
		if err != nil {
			return nil, err
		}
		pool.workers = append(pool.workers, worker)
	}

	return
}

// ActiveWorkers returns true if there are any active workers.
func (h *HSWPool) ActiveWorkers() (active bool) {
	for _, w := range h.workers {
		if w.running {
			active = true
		}
	}

	return
}

// GetHSW gets a new HSW token from the HSW pool.
func (h *HSWPool) GetHSW() (string, error) {
	for {
		if !h.ActiveWorkers() {
			return "", errors.New("no active workers found")
		}
		if len(h.hswPool) != 0 {
			h.hswPoolMutex.Lock()
			defer func() {
				h.hswPool[len(h.hswPool)-1], h.hswPool[0] = h.hswPool[0], h.hswPool[len(h.hswPool)-1]
				h.hswPool = h.hswPool[:len(h.hswPool)-1]
				h.hswPoolMutex.Unlock()
			}()
			return h.hswPool[0], nil
		}
	}
}

// Close closes all workers currently running.
func (h *HSWPool) Close() {
	for _, w := range h.workers {
		if w.running {
			w.Close()
		}
	}
}

// evaluateHsw gets the C token for an hCaptcha n token, and returns it.
func evaluateHsw(s *Solver, c string) (string, error) {
	worker := s.hswPool.workers[s.sRand.Intn(len(s.hswPool.workers))]
	pResp, err := worker.playwrightPage.Evaluate("hsw(\"" + gjson.Get(c, "req").String() + "\");")
	if err != nil {
		return "", err
	}
	return pResp.(string), nil
}

// generateHsw sends a request to the HCaptcha site config system for a HSW token and returns that token.
func generateHsw(worker *HSWWorker) (c string, err error) {
	req, err := http.NewRequest("GET", "https://hcaptcha.com/checksiteconfig?host="+worker.pool.host+"&sitekey="+worker.pool.siteKey+"&sc=1&swa=1", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", worker.pool.userAgent)
	resp, err := worker.client.Do(req)
	if err != nil {
		return "", err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	result := gjson.ParseBytes(b)
	return result.Get("c").String(), nil
}
