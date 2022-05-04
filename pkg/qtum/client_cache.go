package qtum

import (
	"encoding/json"
	"errors"
	"sync"
	"time"
)

/*
336   "method": "decoderawtransaction",
   1   "method": "eth_getBlockByNumber",
1001   "method": "getblock",
   1   "method": "getblockchaininfo",
   1   "method": "getblockhash",
   1   "method": "getblockheader",
 669   "method": "gethexaddress",
 673   "method": "getrawtransaction",
 336   "method": "gettransaction",
 672   "method": "gettxout",
*/
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

type clientCache struct {
	mu      sync.RWMutex
	methods map[string]responses
}

// holds the response for a given method where the map key is param #1
type responses map[string][]byte

func newClientCache() *clientCache {
	return &clientCache{
		methods: make(map[string]responses),
	}
}

func (cache *clientCache) IsCachable(method string) bool {
	for _, m := range cachable_methods {
		if m == method {
			return true
		}
	}
	return false
}

func (cache *clientCache) storeResponse(method string, param interface{}, response []byte) error {
	pb, err := json.Marshal(param)
	if err != nil {
		return errors.New("failed to marshal param")
	}
	cache.mu.RLock()
	defer cache.mu.RUnlock()
	resp, ok := cache.methods[method]
	if !ok {
		resp = make(map[string][]byte)
		cache.methods[method] = resp
		cache.setFlushResponseTimer(method, pb)
	}
	resp[string(pb)] = response
	return nil
}

func (cache *clientCache) getResponse(method string, param interface{}) ([]byte, error) {
	pb, err := json.Marshal(param)
	if err != nil {
		return nil, errors.New("failed to marshal param")
	}
	cache.mu.RLock()
	defer cache.mu.RUnlock()
	if resp, ok := cache.methods[method]; ok {
		if r, ok := resp[string(pb)]; ok {
			return r, nil
		}
	}
	return nil, nil
}

func (cache *clientCache) setFlushResponseTimer(method string, pb []byte) {
	go func() {
		time.Sleep(CACHABLE_METHOD_CACHE_TIMEOUT)
		cache.mu.Lock()
		defer cache.mu.Unlock()
		cache.methods[method][string(pb)] = nil
	}()
}
