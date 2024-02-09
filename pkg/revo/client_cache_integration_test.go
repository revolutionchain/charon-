package revo

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-kit/log"
)

var mockJsonRPCResponse = []byte(`{"error":null,"id":"2693","result":{"bits":"1a0bbe94","chainwork":"0000000000000000000000000000000000000000000002fd7a718803c8cdddaf","confirmations":381397,"difficulty":1428501.632566092,"flags":"proof-of-stake","hash":"2d84b08cef430cf580539a4abee75326e1a0ca0c39f6a2c667e48a24ae0da5c4","hashStateRoot":"c1d92e104215f196fadab16de27a208fe3bd2ddacdfa7ad73cd8c0463b788699","hashUTXORoot":"811f06ab8d39518ce86f0061ad2db70530526417486a8cdc1b28f960f48cddd4","height":1458070,"mediantime":1639381704,"merkleroot":"e558c16bb1868139e0d5e4cc01e2d404cf89f8a8f7fd37d2c28027d86f9b0f44","modifier":"2de0332fab8323d2b7c48b8be34790c34938ac1dfd3a652f138c64123bc9fc15","nTx":3,"nextblockhash":"ba0d74b5d2f8bd6ba80594b124ecc2771334876a5cd058aafd32123a6177285c","nonce":0,"previousblockhash":"c63ca0b83f0ae5fbb88ab181ada70cabcff59422c0a4c8b936325365d55d2b83","prevoutStakeHash":"5a75351b9ba11f525a7d410db55ea2610cf90527805f83b985ce5d33fc46bfdd","prevoutStakeVoutN":883,"proofOfDelegation":"1f51175202748c96624a1f414d967d767b314fd259855c3bc88155426d5da8369140ec0ceb46de548b7f36c5c16863ee63a2cfd46df9c562fd35f627e5af1d3298","proofhash":"0000000000000000000000000000000000000000000000000000000000000000","signature":"20efd150a4d6d8d9b1e8ba1e0d70f3580a9599bfbf382eb436563f6315093da73f04d9bd54a3b02ac37720db5b5104e4aa7ff384f608593443cc019d88d446cc43","size":895,"strippedsize":859,"time":1639381828,"tx":["78af41842e8329f6b4d7f37d821f0aa63042a27f59ebcddff2cde25c6df84465","a3a2941152d33326ab9d8437b4b53f722b747b18d20d309ee31830e5cc2e41d5","830b99f970ae51e70d6298f651a51a8f3f6679902534412c563e88ed621aafd5"],"version":536870912,"versionHex":"20000000","weight":3472}}`)
var logBuffer bytes.Buffer
var logWriter io.Writer = &logBuffer
var logger = log.NewLogfmtLogger(logWriter)

func TestCacheWithClient(t *testing.T) {
	revoMockServer := NewRevoMockServer(mockJsonRPCResponse)
	URL := "http://revouser:revopass@127.0.0.1:6969"
	defer revoMockServer.Close()
	client, err := NewClient(
		true,
		URL,
		SetDebug(true),
		SetLogWriter(logWriter),
		SetLogger(logger),
		SetContext(context.Background()),
	)
	client.URL = revoMockServer.URL
	if err != nil {
		t.Fatal(err)
	}

	var result interface{}
	t.Run("Revo RPC getblock call processed succesfully", func(t *testing.T) {
		err = client.Request(test_method, test_params, &result)
		if err != nil {
			t.Fatal(err)
		}
		assertResponseBody(t, result, test_expectedResult)

	})
	t.Run("Client cache is updated with getblock response", func(t *testing.T) {
		cached_response, err := client.cache.getResponse(test_method, test_params)
		if err != nil {
			t.Fatal("No error expected: ", err)
		}
		if !bytes.Equal(cached_response, test_expectedResult) {
			t.Errorf("\nexpected: %s\n\n, got: %s", string(test_expectedResult), string(cached_response))
		}
	})
	t.Run("Upon second getblock call, response is returned from cache", func(t *testing.T) {
		logBuffer.Reset()
		err = client.Request(test_method, test_params, &result)
		if err != nil {
			t.Fatal(err)
		}
		assertResponseBody(t, result, test_expectedResult)

		outputLog := logBuffer.String()
		if !strings.Contains(outputLog, "revo (CACHED) RPC response") {
			t.Errorf("\nexpected: %s\n\n, got: %s", "revo (CACHED) RPC response", outputLog)
		}

	})
	t.Run("Cache for getblock is flushed after timeout", func(t *testing.T) {
		logBuffer.Reset()
		// wait for the cache to be flushed
		time.Sleep(CACHABLE_METHOD_CACHE_TIMEOUT + time.Millisecond*100)
		cached_response, err := client.cache.getResponse(test_method, test_params)
		if err != nil {
			t.Fatal("No error expected: ", err)
		}
		if !bytes.Equal(cached_response, nil) {
			t.Errorf("\nexpected: nil\n\n, got: %s", string(cached_response))
		}
		outputLog := logBuffer.String()
		if !strings.Contains(outputLog, `msg="flushing cache" reason="cache timeout"`) {
			t.Errorf("expected log message not found: %s", outputLog)
		}
	})
	t.Run("After cache is flushed, it should be updated with a new getblock response", func(t *testing.T) {
		err = client.Request(test_method, test_params, &result)
		if err != nil {
			t.Fatal(err)
		}
		assertResponseBody(t, result, test_expectedResult)

		cached_response, err := client.cache.getResponse(test_method, test_params)
		if err != nil {
			t.Fatal("No error expected: ", err)
		}
		if !bytes.Equal(cached_response, test_expectedResult) {
			t.Errorf("\nexpected: %s\n\n, got: %s", string(test_expectedResult), string(cached_response))
		}
	})
	t.Run("Cache for getblock is flushed again after timeout", func(t *testing.T) {
		logBuffer.Reset()
		// wait for the cache to be flushed
		time.Sleep(CACHABLE_METHOD_CACHE_TIMEOUT + time.Millisecond*100)
		cached_response, err := client.cache.getResponse(test_method, test_params)
		if err != nil {
			t.Fatal("No error expected: ", err)
		}
		if !bytes.Equal(cached_response, nil) {
			t.Errorf("\nexpected: nil\n\n, got: %s", string(cached_response))
		}
		outputLog := logBuffer.String()
		if !strings.Contains(outputLog, `msg="flushing cache" reason="cache timeout"`) {
			t.Errorf("expected log message not found: %s", outputLog)
		}
	})

	// create http client with context.WithCancel to test cache flushing on cancelation
	ctx, canceFunc := context.WithCancel(context.Background())
	client, err = NewClient(
		true,
		URL,
		SetDebug(true),
		SetLogWriter(logWriter),
		SetLogger(logger),
		SetContext(ctx),
	)
	client.URL = revoMockServer.URL
	if err != nil {
		t.Fatal(err)
	}
	t.Run("Cache should be flushed when ctx is canceled", func(t *testing.T) {
		logBuffer.Reset()
		err = client.Request(test_method, test_params, &result)
		if err != nil {
			t.Fatal(err)
		}
		assertResponseBody(t, result, test_expectedResult)

		cached_response, err := client.cache.getResponse(test_method, test_params)
		if err != nil {
			t.Fatal("No error expected: ", err)
		}
		if !bytes.Equal(cached_response, test_expectedResult) {
			t.Errorf("\nexpected: %s\n\n, got: %s", string(test_expectedResult), string(cached_response))
		}

		canceFunc()
		// wait for the cache to be flushed after canceling the context
		time.Sleep(time.Millisecond * 100)
		cached_response, err = client.cache.getResponse(test_method, test_params)
		if err != nil {
			t.Fatal("No error expected: ", err)
		}
		if !bytes.Equal(cached_response, nil) {
			t.Errorf("\nexpected: %v\n, got: %s", nil, string(cached_response))
		}
		// wait for the 'flushing' go routine to write to the log
		time.Sleep(time.Millisecond * 100)
		outputLog := logBuffer.String()
		if !strings.Contains(outputLog, `msg="flushing cache" reason="context canceled"`) {
			t.Errorf("expected log message not found: %s", outputLog)
		}
	})

	// create http client with context.WithTimeout to test cache flushing on cancelation
	ctx, _ = context.WithTimeout(context.Background(), time.Second*2)
	client, err = NewClient(
		true,
		URL,
		SetDebug(true),
		SetLogWriter(logWriter),
		SetLogger(logger),
		SetContext(ctx),
	)
	client.URL = revoMockServer.URL
	if err != nil {
		t.Fatal(err)
	}
	t.Run("Cache should be flushed when ctx times out", func(t *testing.T) {
		logBuffer.Reset()
		err = client.Request(test_method, test_params, &result)
		if err != nil {
			t.Fatal(err)
		}
		assertResponseBody(t, result, test_expectedResult)

		cached_response, err := client.cache.getResponse(test_method, test_params)
		if err != nil {
			t.Fatal("No error expected: ", err)
		}
		if !bytes.Equal(cached_response, test_expectedResult) {
			t.Errorf("\nexpected: %s\n\n, got: %s\n", string(test_expectedResult), string(cached_response))
		}

		time.Sleep(time.Second * 3)

		cached_response, err = client.cache.getResponse(test_method, test_params)
		if err != nil {
			t.Fatal("No error expected: ", err)
		}
		if !bytes.Equal(cached_response, nil) {
			t.Errorf("\nexpected: %v\n, got: %s", nil, string(cached_response))
		}
		// wait for the 'flushing' go routine to write to the log
		time.Sleep(time.Millisecond * 200)
		outputLog := logBuffer.String()
		if !strings.Contains(outputLog, `msg="flushing cache" reason="context canceled"`) {
			t.Errorf("expected log message not found: %s", outputLog)
		}
	})

}

func NewRevoMockServer(body []byte) *httptest.Server {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(body)
	}))
	return server
}

func assertResponseBody(t testing.TB, result interface{}, expected []byte) {
	t.Helper()
	got, err := json.Marshal(result)
	if err != nil {
		t.Fatal("Error not expected: ", err)
	}

	want := expected
	if !bytes.Equal(got, want) {
		t.Errorf("\nexpected: %s\n\n, got: %s", string(want), string(got))
	}
}
