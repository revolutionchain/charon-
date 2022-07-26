package transformer

import (
	"encoding/json"
	"testing"

	"github.com/qtumproject/janus/pkg/internal"
	"github.com/qtumproject/janus/pkg/qtum"
)

type ETHProxyInitializer = func(*qtum.Qtum) ETHProxy

func testETHProxyRequest(t *testing.T, initializer ETHProxyInitializer, requestParams []json.RawMessage, want interface{}) {
	request, err := internal.PrepareEthRPCRequest(1, requestParams)
	if err != nil {
		t.Fatal(err)
	}

	mockedClientDoer := internal.NewDoerMappedMock()
	qtumClient, err := internal.CreateMockedClient(mockedClientDoer)

	internal.SetupGetBlockByHashResponses(t, mockedClientDoer)

	//preparing proxy & executing request
	proxyEth := initializer(qtumClient)
	got, jsonErr := proxyEth.Request(request, internal.NewEchoContext())
	if jsonErr != nil {
		t.Fatalf("Failed to process request on %T.Request(%s): %s", proxyEth, requestParams, jsonErr)
	}

	internal.CheckTestResultEthRequestRPC(*request, want, got, t, false)
}
