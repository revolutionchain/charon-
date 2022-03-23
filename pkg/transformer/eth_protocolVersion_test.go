package transformer

import (
	"encoding/json"
	"testing"

	"github.com/qtumproject/janus/pkg/internal"
)

func TestProtocolVersionReturnsHardcodedValue(t *testing.T) {
	//preparing the request
	requestParams := []json.RawMessage{} //eth_protocolVersion has no params
	request, err := internal.PrepareEthRPCRequest(1, requestParams)
	if err != nil {
		t.Fatal(err)
	}

	proxyEth := ETHProtocolVersion{}
	got, jsonErr := proxyEth.Request(request, nil)
	if jsonErr != nil {
		t.Fatal(jsonErr)
	}

	want := "0x41"

	internal.CheckTestResultDefault(want, got, t, false)
}
