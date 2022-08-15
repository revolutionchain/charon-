package analytics

import "sync"

type Analytics struct {
	success       int
	failures      int
	lastRequest   int
	lastRequests  []bool
	totalRequests int

	mutex sync.RWMutex
}

func NewAnalytics(requests int) *Analytics {
	analytics := &Analytics{
		success:       0,
		failures:      0,
		lastRequest:   0,
		lastRequests:  make([]bool, requests),
		totalRequests: requests,
	}

	return analytics
}

func (a *Analytics) GetSuccessRate() float32 {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	total := a.success + a.failures

	if total != a.totalRequests {
		// if there are not enough requests, don't let the first request get reported as an alert
		return 1
	}

	return float32(a.success) / float32(total)
}

func (a *Analytics) Success() {
	a.bump(true)
}

func (a *Analytics) Failure() {
	a.bump(false)
}

func (a *Analytics) bump(success bool) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	total := a.success + a.failures

	if success {
		a.success++
	} else {
		a.failures++
	}

	if total >= a.totalRequests {
		// push the last request off
		if a.lastRequests[a.lastRequest] {
			a.success--
		} else {
			a.failures--
		}

		if a.success < 0 {
			panic("internal analytics success count is < 0")
		}

		if a.failures < 0 {
			panic("internal analytics failure count is < 0")
		}
	}

	a.lastRequests[a.lastRequest] = success
	a.lastRequest = (a.lastRequest + 1) % a.totalRequests
}
