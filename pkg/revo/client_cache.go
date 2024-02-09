package revo

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"sync"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

// sets the timeout for flushing out the cashed memory
const CACHABLE_METHOD_CACHE_TIMEOUT = time.Second * 15

const (
	RevoMethodGetblock             = "getblock"
	RevoMethodGetblockhash         = "getblockhash"
	RevoMethodGetblockheader       = "getblockheader"
	RevoMethodGetblockchaininfo    = "getblockchaininfo"
	RevoMethodGethexaddress        = "gethexaddress"
	RevoMethodGetrawtransaction    = "getrawtransaction"
	RevoMethodGettransaction       = "gettransaction"
	RevoMethodGettxout             = "gettxout"
	RevoMethodDecoderawtransaction = "decoderawtransaction"
)

var cachable_methods = []string{
	RevoMethodGetblock,
	// RevoMethodGetblockhash,
	// RevoMethodGetblockheader,
	// RevoMethodGetblockchaininfo,
	RevoMethodGethexaddress,
	RevoMethodGetrawtransaction,
	// RevoMethodGettransaction,
	RevoMethodGettxout,
	RevoMethodDecoderawtransaction,
}

// stores the rpc response for 'method' and 'params' in the cache
// 'methods' is a map where keys are method names and values are maps of rpc responses
type clientCache struct {
	mu        sync.RWMutex
	ctx       context.Context
	logger    log.Logger
	logWriter io.Writer
	debug     bool
	methods   map[string]responses
}

// 'responses' is a map where keys are rpc param bytes, and values are response bytes (for the given method)
type responses map[string][]byte

func newClientCache() *clientCache {
	return &clientCache{
		methods: make(map[string]responses),
	}
}

// checks if the method should be cached
func (cache *clientCache) isCachable(method string) bool {
	for _, m := range cachable_methods {
		if m == method {
			return true
		}
	}
	return false
}

// stores the rpc response for 'method' and 'params' in the cache
func (cache *clientCache) storeResponse(method string, params interface{}, response []byte) error {
	parambytes, err := json.Marshal(params)
	if err != nil {
		return errors.New("failed to marshal params")
	}
	cache.mu.Lock()
	defer cache.mu.Unlock()
	responses, ok := cache.methods[method]
	if !ok {
		responses = make(map[string][]byte)
		cache.methods[method] = responses
	}
	if _, ok := responses[string(parambytes)]; !ok {
		responses[string(parambytes)] = response
		cache.setFlushResponseTimer(method, parambytes)
	}
	return nil
}

// returns the cached rpc response for 'method' and 'params'
func (cache *clientCache) getResponse(method string, params interface{}) ([]byte, error) {
	parambytes, err := json.Marshal(params)
	if err != nil {
		return nil, errors.New("failed to marshal param")
	}
	cache.mu.RLock()
	defer cache.mu.RUnlock()
	if resp, ok := cache.methods[method]; ok {
		if r, ok := resp[string(parambytes)]; ok {
			return r, nil
		}
	}
	return nil, nil
}

// set a timer to flush the cached rpc response for 'method' and 'parambytes'
func (cache *clientCache) setFlushResponseTimer(method string, parambytes []byte) {
	go func() {
		// TODO check if this works as expected
		var done <-chan struct{}
		if cache.ctx != nil {
			done = cache.ctx.Done()
		} else {
			done = context.Background().Done()
		}
		select {
		case <-time.After(CACHABLE_METHOD_CACHE_TIMEOUT):
			cache.getDebugLogger().Log("msg", "flushing cache", "reason", "cache timeout", "method", method)
		case <-done:
			cache.getDebugLogger().Log("msg", "flushing cache", "reason", "context canceled", "method", method)
		}
		cache.mu.Lock()
		defer cache.mu.Unlock()
		delete(cache.methods[method], string(parambytes))
	}()
}

func (cache *clientCache) setContext(ctx context.Context) {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	if cache.ctx == nil {
		cache.ctx = ctx
	}
}

func (cache *clientCache) configLogger(logWriter io.Writer, debug bool) {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	if cache.logger == nil {
		cache.debug = debug
		cache.logWriter = logWriter
	}
}

func (cache *clientCache) getDebugLogger() log.Logger {
	if !cache.isDebugEnabled() {
		return log.NewNopLogger()
	}
	if cache.logger == nil {
		cache.logger = log.NewLogfmtLogger(cache.logWriter)
		cache.logger = log.With(level.Debug(cache.logger), "component", "clientCache")
	}
	return cache.logger
}

func (cache *clientCache) isDebugEnabled() bool {
	return cache.debug
}
