package main

import (
	"errors"
	"github.com/corpix/uarand"
	"github.com/mxschmitt/playwright-go"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"sync"
)

// HSW is an HSW token, containing both the C and N values used for requests to HCaptcha.
type HSW struct {
	// C and N are tokens used by HCaptcha on the getcaptcha endpoint as a form of authorization.
	// These can not be easily spoofed, as they contain browser/device level data.
	C, N string
}

// HSWWorker is a worker for an HSW pool.
// Workers should always be created through the NewWorker function.
type HSWWorker struct {
	playwrightInstance *playwright.Playwright
	playwrightBrowser  playwright.Browser
	playwrightPage     playwright.Page
	pool               *HSWPool
	running            bool
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
	worker := &HSWWorker{pool: pool, running: true}

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
		var hsw HSW
		for {
			hsw, err = generateHsw(pool.host, pool.siteKey, worker.playwrightPage)

			// We run the check after the HSW has been generated,
			// because we are more likely to be already generating HSW
			// when a shutdown is requested.
			if !worker.running {
				break
			}
			if err != nil {
				pool.log.Error(err)
				continue
			}
			pool.hswPoolMutex.Lock()
			pool.hswPool = append(pool.hswPool, hsw)
			pool.hswPoolMutex.Unlock()
		}
	}()

	return worker, nil
}

// HSWPool is a pool of HSW tokens that get regenerated on the fly by HSW workers.
// These are used to generate valid captcha codes.
type HSWPool struct {
	hswPool       []HSW
	hswPoolMutex  *sync.Mutex
	workers       []*HSWWorker
	hswScriptUrl  *string
	host, siteKey string
	log           *logrus.Logger
}

// NewHSWPool creates a new HSW pool with the amount of workers specified.
func NewHSWPool(host, siteKey, hswScriptUrl string, log *logrus.Logger, workers int) (pool *HSWPool, err error) {
	if workers <= 0 {
		return nil, errors.New("invalid worker amount for pool")
	}
	pool = &HSWPool{
		hswPoolMutex: &sync.Mutex{},
		hswScriptUrl: playwright.String(hswScriptUrl),
		host:         host,
		siteKey:      siteKey,
		log:          log,
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
func (h *HSWPool) GetHSW() (HSW, error) {
	for {
		if !h.ActiveWorkers() {
			return HSW{}, errors.New("no active workers found")
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

// generateHsw sends a request to the HCaptcha site config system for a HSW token.
// Then, we use the original token provided to us to generate HSW to send in our captcha requests.
// The HSW is generated through Playwright. On initial startup, Playwright (running Firefox, WebKit seems to fail)
// opens a new empty tab, injects the HSW script by HCaptcha, and then evaluates the HSW using the page evaluate
// function. It then returns the response, plus the HSW token and an error in case anything went wrong.
func generateHsw(host, siteKey string, page playwright.Page) (hsw HSW, err error) {
	req, err := http.NewRequest("GET", "https://hcaptcha.com/checksiteconfig?host="+host+"&siteKey="+siteKey+"&sc=1&swa=1", nil)
	if err != nil {
		return HSW{}, err
	}
	req.Header.Set("User-Agent", uarand.GetRandom())
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return HSW{}, err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return HSW{}, err
	}
	result := gjson.ParseBytes(b)
	pResp, err := page.Evaluate("hsw(\"" + result.Get("c.req").String() + "\");")
	if err != nil {
		return HSW{}, err
	}
	return HSW{
		C: result.Get("c").String(),
		N: pResp.(string),
	}, nil
}
