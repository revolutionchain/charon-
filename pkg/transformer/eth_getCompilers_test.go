package transformer

import (
	"encoding/json"
	"testing"

	"github.com/qtumproject/janus/pkg/internal"
)

func TestGetCompilersReturnsEmptyArray(t *testing.T) {
	//preparing the request
	requestParams := []json.RawMessage{} //eth_getCompilers has no params
	request, err := internal.PrepareEthRPCRequest(1, requestParams)
	if err != nil {
		t.Fatal(err)
	}

	proxyEth := ETHGetCompilers{}
	got, jsonErr := proxyEth.Request(request, internal.NewEchoContext())
	if jsonErr != nil {
		t.Fatal(jsonErr)
	}

	want := []string{}

	internal.CheckTestResultDefault(want, got, t, false)
}
