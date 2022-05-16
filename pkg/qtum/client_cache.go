package qtum

import (
	"encoding/json"
	"errors"
	"sync"
	"time"
)

const CACHABLE_METHOD_CACHE_TIMEOUT = time.Second * 10

const (
	QtumMethodGetblock             = "getblock"
	QtumMethodGetblockhash         = "getblockhash"
	QtumMethodGetblockheader       = "getblockheader"
	QtumMethodGetblockchaininfo    = "getblockchaininfo"
	QtumMethodGethexaddress        = "gethexaddress"
	QtumMethodGetrawtransaction    = "getrawtransaction"
	QtumMethodGettransaction       = "gettransaction"
	QtumMethodGettxout             = "gettxout"
	QtumMethodDecoderawtransaction = "decoderawtransaction"
)

var cachable_methods = []string{
	QtumMethodGetblock,
	// QtumMethodGetblockhash,
	// QtumMethodGetblockheader,
	// QtumMethodGetblockchaininfo,
	QtumMethodGethexaddress,
	QtumMethodGetrawtransaction,
	// QtumMethodGettransaction,
	QtumMethodGettxout,
	QtumMethodDecoderawtransaction,
}

// stores the rpc response for 'method' and 'params' in the cache
// 'methods' is a map where keys are method names and values are maps of rpc responses
type clientCache struct {
	mu      sync.RWMutex
	methods map[string]responses
}

// responses is a map where keys are rpc param bytes, and values are response bytes (for the given method)
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
		time.Sleep(CACHABLE_METHOD_CACHE_TIMEOUT)
		cache.mu.Lock()
		defer cache.mu.Unlock()
		delete(cache.methods[method], string(parambytes))
	}()
}
