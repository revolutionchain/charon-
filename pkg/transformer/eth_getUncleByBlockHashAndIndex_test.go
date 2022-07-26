package transformer

import (
	"encoding/json"
	"testing"

	"github.com/qtumproject/janus/pkg/internal"
)

func TestGetUncleByBlockHashAndIndexReturnsNil(t *testing.T) {
	// request body doesn't matter, there is no QTUM object to proxy calls to
	requestParams := []json.RawMessage{}
	request, err := internal.PrepareEthRPCRequest(1, requestParams)
	if err != nil {
		t.Fatal(err)
	}

	proxyEth := ETHGetUncleByBlockHashAndIndex{}
	got, jsonErr := proxyEth.Request(request, internal.NewEchoContext())
	if jsonErr != nil {
		t.Fatal(jsonErr)
	}

	want := interface{}(nil)

	internal.CheckTestResultEthRequestRPC(*request, want, got, t, false)
}
