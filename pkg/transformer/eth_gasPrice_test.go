package transformer

import (
	"encoding/json"
	"testing"

	"github.com/qtumproject/janus/pkg/internal"
)

func TestGasPriceRequest(t *testing.T) {
	//preparing request
	requestParams := []json.RawMessage{}
	request, err := internal.PrepareEthRPCRequest(1, requestParams)
	if err != nil {
		t.Fatal(err)
	}

	mockedClientDoer := internal.NewDoerMappedMock()
	qtumClient, err := internal.CreateMockedClient(mockedClientDoer)
	if err != nil {
		t.Fatal(err)
	}

	//preparing proxy & executing request
	proxyEth := ProxyETHGasPrice{qtumClient}
	got, jsonErr := proxyEth.Request(request, internal.NewEchoContext())
	if jsonErr != nil {
		t.Fatal(jsonErr)
	}

	want := string("0x5d21dba000") //price is hardcoded inside the implement

	internal.CheckTestResultDefault(want, got, t, false)
}
