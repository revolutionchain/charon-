package qtum

import (
	"testing"
	"time"
)

var (
	test_method = "getblock"

	test_params = struct {
		hash string
		idx  int
	}{
		"2d84b08cef430cf580539a4abee75326e1a0ca0c39f6a2c667e48a24ae0da5c4",
		1,
	}
	/*
	   {
	   	"error": null,
	   	"id": "2693",
	   	"result": {
	   		"bits": "1a0bbe94",
	   		"chainwork": "0000000000000000000000000000000000000000000002fd7a718803c8cdddaf",
	   		"confirmations": 381397,
	   		"difficulty": 1428501.632566092,
	   		"flags": "proof-of-stake",
	   		"hash": "2d84b08cef430cf580539a4abee75326e1a0ca0c39f6a2c667e48a24ae0da5c4",
	   		"hashStateRoot": "c1d92e104215f196fadab16de27a208fe3bd2ddacdfa7ad73cd8c0463b788699",
	   		"hashUTXORoot": "811f06ab8d39518ce86f0061ad2db70530526417486a8cdc1b28f960f48cddd4",
	   		"height": 1458070,
	   		"mediantime": 1639381704,
	   		"merkleroot": "e558c16bb1868139e0d5e4cc01e2d404cf89f8a8f7fd37d2c28027d86f9b0f44",
	   		"modifier": "2de0332fab8323d2b7c48b8be34790c34938ac1dfd3a652f138c64123bc9fc15",
	   		"nTx": 3,
	   		"nextblockhash": "ba0d74b5d2f8bd6ba80594b124ecc2771334876a5cd058aafd32123a6177285c",
	   		"nonce": 0,
	   		"previousblockhash": "c63ca0b83f0ae5fbb88ab181ada70cabcff59422c0a4c8b936325365d55d2b83",
	   		"prevoutStakeHash": "5a75351b9ba11f525a7d410db55ea2610cf90527805f83b985ce5d33fc46bfdd",
	   		"prevoutStakeVoutN": 883,
	   		"proofOfDelegation": "1f51175202748c96624a1f414d967d767b314fd259855c3bc88155426d5da8369140ec0ceb46de548b7f36c5c16863ee63a2cfd46df9c562fd35f627e5af1d3298",
	   		"proofhash": "0000000000000000000000000000000000000000000000000000000000000000",
	   		"signature": "20efd150a4d6d8d9b1e8ba1e0d70f3580a9599bfbf382eb436563f6315093da73f04d9bd54a3b02ac37720db5b5104e4aa7ff384f608593443cc019d88d446cc43",
	   		"size": 895,
	   		"strippedsize": 859,
	   		"time": 1639381828,
	   		"tx": [
	   		"78af41842e8329f6b4d7f37d821f0aa63042a27f59ebcddff2cde25c6df84465",
	   		"a3a2941152d33326ab9d8437b4b53f722b747b18d20d309ee31830e5cc2e41d5",
	   		"830b99f970ae51e70d6298f651a51a8f3f6679902534412c563e88ed621aafd5"
	   		],
	   		"version": 536870912,
	   		"versionHex": "20000000",
	   		"weight": 3472
	   	}
	   }
	*/

	test_expectedResult = []byte(`{"bits":"1a0bbe94","chainwork":"0000000000000000000000000000000000000000000002fd7a718803c8cdddaf","confirmations":381397,"difficulty":1428501.632566092,"flags":"proof-of-stake","hash":"2d84b08cef430cf580539a4abee75326e1a0ca0c39f6a2c667e48a24ae0da5c4","hashStateRoot":"c1d92e104215f196fadab16de27a208fe3bd2ddacdfa7ad73cd8c0463b788699","hashUTXORoot":"811f06ab8d39518ce86f0061ad2db70530526417486a8cdc1b28f960f48cddd4","height":1458070,"mediantime":1639381704,"merkleroot":"e558c16bb1868139e0d5e4cc01e2d404cf89f8a8f7fd37d2c28027d86f9b0f44","modifier":"2de0332fab8323d2b7c48b8be34790c34938ac1dfd3a652f138c64123bc9fc15","nTx":3,"nextblockhash":"ba0d74b5d2f8bd6ba80594b124ecc2771334876a5cd058aafd32123a6177285c","nonce":0,"previousblockhash":"c63ca0b83f0ae5fbb88ab181ada70cabcff59422c0a4c8b936325365d55d2b83","prevoutStakeHash":"5a75351b9ba11f525a7d410db55ea2610cf90527805f83b985ce5d33fc46bfdd","prevoutStakeVoutN":883,"proofOfDelegation":"1f51175202748c96624a1f414d967d767b314fd259855c3bc88155426d5da8369140ec0ceb46de548b7f36c5c16863ee63a2cfd46df9c562fd35f627e5af1d3298","proofhash":"0000000000000000000000000000000000000000000000000000000000000000","signature":"20efd150a4d6d8d9b1e8ba1e0d70f3580a9599bfbf382eb436563f6315093da73f04d9bd54a3b02ac37720db5b5104e4aa7ff384f608593443cc019d88d446cc43","size":895,"strippedsize":859,"time":1639381828,"tx":["78af41842e8329f6b4d7f37d821f0aa63042a27f59ebcddff2cde25c6df84465","a3a2941152d33326ab9d8437b4b53f722b747b18d20d309ee31830e5cc2e41d5","830b99f970ae51e70d6298f651a51a8f3f6679902534412c563e88ed621aafd5"],"version":536870912,"versionHex":"20000000","weight":3472}`)
)

func TestClientCache(t *testing.T) {
	cache := newClientCache()
	t.Run("Only cachable methods should be allowed", func(t *testing.T) {
		tests := []struct {
			method   string
			cachable bool
		}{
			{"getblock", true},
			{"getblockhash", false},
			{"getblockheader", false},
			{"getblockchaininfo", false},
			{"gethexaddress", true},
			{"getrawtransaction", true},
			{"gettransaction", false},
			{"gettxout", true},
			{"decoderawtransaction", true},
		}
		for _, test := range tests {
			if cache.isCachable(test.method) != test.cachable {
				t.Errorf("expected %v, got %v", test.cachable, cache.isCachable(test.method))
			}
		}
	})

	t.Run("Should allow to correctly store response for given method and params", func(t *testing.T) {

		err := cache.storeResponse(test_method, test_params, test_expectedResult)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		cachedResp, err := cache.getResponse(test_method, test_params)
		if err != nil {
			t.Fatal("expected to find cached response")
		}
		if string(cachedResp) != string(test_expectedResult) {
			t.Fatalf("expected to find %v, got %v", string(test_expectedResult), string(cachedResp))
		}
	})

	t.Run("Should return the cached response for given method and params", func(t *testing.T) {
		cachedResp, err := cache.getResponse(test_method, test_params)
		if err != nil {
			t.Fatal("expected to find cached response")
		}
		if string(cachedResp) != string(test_expectedResult) {
			t.Fatalf("expected to find %v, got %v", string(test_expectedResult), string(cachedResp))
		}
	})
	t.Run("cached response should be flushed after timeout", func(t *testing.T) {
		time.Sleep(CACHABLE_METHOD_CACHE_TIMEOUT + 5*time.Millisecond)
		cachedResp, err := cache.getResponse(test_method, test_params)
		if err != nil {
			t.Fatal("no error expected")
		}
		if cachedResp != nil {
			t.Fatalf("expected to find nil, got %v", cachedResp)
		}
	})

	t.Run("Should return nil when no response is cached", func(t *testing.T) {
		notCashedMethod := "gethexaddress"
		cachedResp, err := cache.getResponse(notCashedMethod, test_params)
		if err != nil {
			t.Fatal("no error expected")
		}
		if cachedResp != nil {
			t.Fatalf("expected to find nil, got %v", cachedResp)
		}
	})

}
